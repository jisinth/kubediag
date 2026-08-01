// Package nodes implements the node diagnostic checks described in
// docs/architecture.md: readiness and resource-pressure conditions.
package nodes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jisinth/kubediag/internal/recommendations"
	"github.com/jisinth/kubediag/internal/reporter"
)

const category = "Nodes"

// List returns all nodes in the cluster.
func List(ctx context.Context, clientset kubernetes.Interface) ([]corev1.Node, error) {
	list, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}
	return list.Items, nil
}

// Check runs readiness and pressure-condition checks across all nodes and
// returns one issue per unhealthy condition found.
func Check(ctx context.Context, clientset kubernetes.Interface) ([]reporter.Issue, error) {
	nodeList, err := List(ctx, clientset)
	if err != nil {
		return nil, err
	}

	var issues []reporter.Issue
	for _, node := range nodeList {
		issues = append(issues, checkNode(node)...)
	}
	return issues, nil
}

func checkNode(node corev1.Node) []reporter.Issue {
	var issues []reporter.Issue

	conditions := map[corev1.NodeConditionType]corev1.ConditionStatus{}
	for _, cond := range node.Status.Conditions {
		conditions[cond.Type] = cond.Status
	}

	if conditions[corev1.NodeReady] != corev1.ConditionTrue {
		issues = append(issues, newIssue(node.Name, "NotReady", reporter.SeverityCritical, recommendations.NodeNotReady,
			fmt.Sprintf("Node %s is not in Ready state", node.Name)))
	}
	if conditions[corev1.NodeMemoryPressure] == corev1.ConditionTrue {
		issues = append(issues, newIssue(node.Name, "MemoryPressure", reporter.SeverityWarning, recommendations.NodeMemoryPressure,
			fmt.Sprintf("Node %s is under memory pressure", node.Name)))
	}
	if conditions[corev1.NodeDiskPressure] == corev1.ConditionTrue {
		issues = append(issues, newIssue(node.Name, "DiskPressure", reporter.SeverityWarning, recommendations.NodeDiskPressure,
			fmt.Sprintf("Node %s is under disk pressure", node.Name)))
	}
	if conditions[corev1.NodePIDPressure] == corev1.ConditionTrue {
		issues = append(issues, newIssue(node.Name, "PIDPressure", reporter.SeverityWarning, recommendations.NodePIDPressure,
			fmt.Sprintf("Node %s is under PID pressure", node.Name)))
	}

	return issues
}

func newIssue(nodeName, title string, severity reporter.Severity, code recommendations.Code, reason string) reporter.Issue {
	issue := reporter.Issue{
		Category: category,
		Severity: severity,
		Title:    fmt.Sprintf("Node %s: %s", nodeName, title),
		Reason:   reason,
	}
	if rec, ok := recommendations.For(code); ok {
		issue.Recommendation = rec.Text
		issue.DocsLink = rec.DocsLink
	}
	return issue
}
