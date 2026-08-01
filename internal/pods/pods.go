// Package pods implements the pod diagnostic checks described in
// docs/architecture.md: scheduling failures, crash loops, image pull
// failures, and OOM kills.
package pods

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jisinth/kubediag/internal/recommendations"
	"github.com/jisinth/kubediag/internal/reporter"
)

const category = "Pods"

// List returns all pods across all namespaces.
func List(ctx context.Context, clientset kubernetes.Interface) ([]corev1.Pod, error) {
	list, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return list.Items, nil
}

// Check runs pod health checks across all namespaces and returns one issue
// per unhealthy pod condition found.
func Check(ctx context.Context, clientset kubernetes.Interface) ([]reporter.Issue, error) {
	podList, err := List(ctx, clientset)
	if err != nil {
		return nil, err
	}

	var issues []reporter.Issue
	for _, pod := range podList {
		issues = append(issues, checkPod(pod)...)
	}
	return issues, nil
}

func checkPod(pod corev1.Pod) []reporter.Issue {
	var issues []reporter.Issue

	switch pod.Status.Phase {
	case corev1.PodPending:
		if issue, ok := pendingIssue(pod); ok {
			issues = append(issues, issue)
		}
	case corev1.PodFailed:
		issues = append(issues, newIssue(pod, "Failed", reporter.SeverityCritical, recommendations.PodFailed,
			fmt.Sprintf("Pod %s/%s is in Failed phase: %s", pod.Namespace, pod.Name, pod.Status.Reason)))
	}

	for _, cs := range append(append([]corev1.ContainerStatus{}, pod.Status.InitContainerStatuses...), pod.Status.ContainerStatuses...) {
		if cs.State.Waiting != nil {
			switch cs.State.Waiting.Reason {
			case "CrashLoopBackOff":
				issues = append(issues, newIssue(pod, "CrashLoopBackOff", reporter.SeverityCritical, recommendations.PodCrashLoopBackOff,
					fmt.Sprintf("Container %s in pod %s/%s is crash-looping: %s", cs.Name, pod.Namespace, pod.Name, cs.State.Waiting.Message)))
			case "ImagePullBackOff", "ErrImagePull":
				issues = append(issues, newIssue(pod, "ImagePullBackOff", reporter.SeverityCritical, recommendations.PodImagePullBackOff,
					fmt.Sprintf("Container %s in pod %s/%s cannot pull its image: %s", cs.Name, pod.Namespace, pod.Name, cs.State.Waiting.Message)))
			}
		}
		if cs.State.Terminated != nil && cs.State.Terminated.Reason == "OOMKilled" {
			issues = append(issues, newIssue(pod, "OOMKilled", reporter.SeverityCritical, recommendations.PodOOMKilled,
				fmt.Sprintf("Container %s in pod %s/%s was OOMKilled", cs.Name, pod.Namespace, pod.Name)))
		}
		if cs.LastTerminationState.Terminated != nil && cs.LastTerminationState.Terminated.Reason == "OOMKilled" && cs.State.Waiting == nil {
			issues = append(issues, newIssue(pod, "OOMKilled", reporter.SeverityWarning, recommendations.PodOOMKilled,
				fmt.Sprintf("Container %s in pod %s/%s was previously OOMKilled", cs.Name, pod.Namespace, pod.Name)))
		}
	}

	return issues
}

// pendingIssue inspects the PodScheduled condition to classify why a pod is
// stuck Pending: insufficient CPU, insufficient memory, or otherwise
// unschedulable.
func pendingIssue(pod corev1.Pod) (reporter.Issue, bool) {
	for _, cond := range pod.Status.Conditions {
		if cond.Type != corev1.PodScheduled || cond.Status == corev1.ConditionTrue {
			continue
		}

		message := strings.ToLower(cond.Message)
		switch {
		case strings.Contains(message, "insufficient cpu"):
			return newIssue(pod, "Pending", reporter.SeverityWarning, recommendations.PodPendingInsufficientCPU,
				fmt.Sprintf("Pod %s/%s is Pending: insufficient CPU", pod.Namespace, pod.Name)), true
		case strings.Contains(message, "insufficient memory"):
			return newIssue(pod, "Pending", reporter.SeverityWarning, recommendations.PodPendingInsufficientMemory,
				fmt.Sprintf("Pod %s/%s is Pending: insufficient memory", pod.Namespace, pod.Name)), true
		default:
			return newIssue(pod, "Pending", reporter.SeverityWarning, recommendations.PodPendingUnschedulable,
				fmt.Sprintf("Pod %s/%s is Pending: %s", pod.Namespace, pod.Name, cond.Message)), true
		}
	}
	return reporter.Issue{}, false
}

func newIssue(pod corev1.Pod, title string, severity reporter.Severity, code recommendations.Code, reason string) reporter.Issue {
	issue := reporter.Issue{
		Category: category,
		Severity: severity,
		Title:    fmt.Sprintf("Pod %s/%s: %s", pod.Namespace, pod.Name, title),
		Reason:   reason,
	}
	if rec, ok := recommendations.For(code); ok {
		issue.Recommendation = rec.Text
		issue.DocsLink = rec.DocsLink
	}
	return issue
}
