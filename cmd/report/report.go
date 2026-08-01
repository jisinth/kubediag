// Package report implements "kubediag report": run a full diagnosis and
// export it to a file. HTML dashboard export (charts, tables) is planned
// for v4.0 — see docs/roadmap.md; this currently supports table and JSON.
package report

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/jisinth/kubediag/internal/config"
	"github.com/jisinth/kubediag/internal/diagnose"
)

// NewCommand builds the "kubediag report" command.
func NewCommand() *cobra.Command {
	var outPath string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Run a full diagnosis and export the report to a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := diagnose.Run(cmd.Context(), config.Kubeconfig, config.Context)
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			if outPath != "" {
				f, err := os.Create(outPath)
				if err != nil {
					return err
				}
				defer f.Close()
				w = f
			}

			switch config.Output {
			case "html":
				return report.WriteHTML(w)
			case "table":
				report.WriteTable(w)
				return nil
			default:
				return report.WriteJSON(w)
			}
		},
	}

	cmd.Flags().StringVar(&outPath, "out", "", "file to write the report to (defaults to stdout)")
	return cmd
}
