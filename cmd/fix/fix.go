// Package fix implements "kubediag fix": preview-confirm-apply-verify
// workflow for automatically remediating detected issues (e.g. "kubediag fix
// pending-pods"). Planned for v4.0 — see docs/roadmap.md. This is currently
// a stub that documents the intended interface without mutating the
// cluster.
package fix

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCommand builds the "kubediag fix" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fix [target]",
		Short: "Preview and apply an automatic fix for a detected issue (planned for v4.0)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "<target>"
			if len(args) == 1 {
				target = args[0]
			}
			fmt.Fprintf(cmd.OutOrStdout(),
				"kubediag fix %s: auto-fix is planned for v4.0 (preview -> confirm -> apply -> verify).\n"+
					"No changes have been made to the cluster. See docs/roadmap.md.\n", target)
			return nil
		},
	}
}
