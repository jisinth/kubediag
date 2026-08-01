// Package nodes implements "kubediag nodes": run only the node diagnostic
// checks.
package nodes

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jisinth/kubediag/internal/cluster"
	"github.com/jisinth/kubediag/internal/config"
	"github.com/jisinth/kubediag/internal/nodes"
	"github.com/jisinth/kubediag/internal/reporter"
)

// NewCommand builds the "kubediag nodes" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "nodes",
		Short: "Diagnose node health (Ready, MemoryPressure, DiskPressure, PIDPressure)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := cluster.Connect(config.Kubeconfig, config.Context)
			if err != nil {
				return err
			}

			issues, err := nodes.Check(cmd.Context(), client.Clientset)
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
