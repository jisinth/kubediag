// Package security implements "kubediag security". Planned for v3.0 — see
// docs/roadmap.md.
package security

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand builds the "kubediag security" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Diagnose RBAC, pod security, and secret exposure (planned for v3.0)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Security diagnostics are planned for v3.0 — see docs/roadmap.md.")
			return nil
		},
	}
}
