// Package diagnose implements "kubediag diagnose": the main entry point
// that runs every v1.0 diagnostic module and prints a cluster report.
package diagnose

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jisinth/kubediag/internal/config"
	"github.com/jisinth/kubediag/internal/diagnose"
)

// NewCommand builds the "kubediag diagnose" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diagnose",
		Short: "Run a full cluster diagnosis and print a health report",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := diagnose.Run(cmd.Context(), config.Kubeconfig, config.Context)
			if err != nil {
				return err
			}

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
