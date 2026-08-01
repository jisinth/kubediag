// Package diagnose orchestrates a full diagnostic run: connect to the
// cluster, run each diagnostic module, and assemble the results into a
// reporter.Report.
package diagnose

import (
	"context"
	"fmt"

	"github.com/jisinth/kubediag/internal/cluster"
	"github.com/jisinth/kubediag/internal/nodes"
	"github.com/jisinth/kubediag/internal/pods"
	"github.com/jisinth/kubediag/internal/reporter"
)

// Run connects to the cluster and executes every v1.0 diagnostic module
// (cluster health, nodes, pods), returning the aggregate report.
func Run(ctx context.Context, kubeconfigPath, kubeContext string) (reporter.Report, error) {
	client, err := cluster.Connect(kubeconfigPath, kubeContext)
	if err != nil {
		return reporter.Report{}, err
	}

	var checks []reporter.Check
	var issues []reporter.Issue

	version, err := client.Version(ctx)
	if err != nil {
		return reporter.Report{}, fmt.Errorf("failed to reach API server: %w", err)
	}
	checks = append(checks, reporter.Check{Category: "Cluster", Title: "API Server Healthy", Passed: true})

	nodeIssues, err := nodes.Check(ctx, client.Clientset)
	if err != nil {
		return reporter.Report{}, err
	}
	issues = append(issues, nodeIssues...)

	podIssues, err := pods.Check(ctx, client.Clientset)
	if err != nil {
		return reporter.Report{}, err
	}
	issues = append(issues, podIssues...)

	return reporter.NewReport(version, checks, issues), nil
}
