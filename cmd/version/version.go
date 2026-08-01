// Package version implements "kubediag version".
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These are overridden at build time via:
//
//	go build -ldflags "-X github.com/jisinth/kubediag/cmd/version.Version=v1.0.0 \
//	  -X github.com/jisinth/kubediag/cmd/version.Commit=$(git rev-parse --short HEAD) \
//	  -X github.com/jisinth/kubediag/cmd/version.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewCommand builds the "kubediag version" command.
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the kubediag version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "kubediag %s (commit %s, built %s)\n", Version, Commit, Date)
			return nil
		},
	}
}
