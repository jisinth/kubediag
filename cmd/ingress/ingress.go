// Package ingress implements "kubediag ingress". Planned for v2.0 — see
// docs/roadmap.md.
package ingress

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand builds the "kubediag ingress" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ingress",
		Short: "Diagnose Ingress resources, controllers, and TLS (planned for v2.0)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Ingress diagnostics are planned for v2.0 — see docs/roadmap.md.")
			return nil
		},
	}
}
