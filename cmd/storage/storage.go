// Package storage implements "kubediag storage". Planned for v2.0 — see
// docs/roadmap.md.
package storage

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand builds the "kubediag storage" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "storage",
		Short: "Diagnose PVCs, PVs, and StorageClasses (planned for v2.0)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Storage diagnostics are planned for v2.0 — see docs/roadmap.md.")
			return nil
		},
	}
}
