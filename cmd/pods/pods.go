// Package pods implements "kubediag pods": run only the pod diagnostic
// checks.
package pods

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jisinth/kubediag/internal/cluster"
	"github.com/jisinth/kubediag/internal/config"
	"github.com/jisinth/kubediag/internal/pods"
	"github.com/jisinth/kubediag/internal/reporter"
)

// NewCommand builds the "kubediag pods" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pods",
		Short: "Diagnose pod health (CrashLoopBackOff, ImagePullBackOff, Pending, Failed, OOMKilled)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := cluster.Connect(config.Kubeconfig, config.Context)
			if err != nil {
				return err
			}

			issues, err := pods.Check(cmd.Context(), client.Clientset)
			if err != nil {
				return err
			}

			report := reporter.NewReport("", nil, issues)
			switch config.Output {
			case "json":
				return report.WriteJSON(os.Stdout)
			case "html":
				return report.WriteHTML(os.Stdout)
			default:
				report.WriteTable(os.Stdout)
				return nil
			}
		},
	}
}
