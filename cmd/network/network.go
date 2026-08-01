// Package network implements "kubediag network". Planned for v2.0 — see
// docs/roadmap.md.
package network

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand builds the "kubediag network" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "network",
		Short: "Diagnose DNS, CoreDNS, Services, and network policies (planned for v2.0)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Network diagnostics are planned for v2.0 — see docs/roadmap.md.")
			return nil
		},
	}
}
