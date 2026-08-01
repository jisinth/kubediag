// Package cmd wires together the kubediag CLI: the root command and every
// subcommand.
package cmd

import (
	"github.com/spf13/cobra"

	cmddiagnose "github.com/jisinth/kubediag/cmd/diagnose"
	cmdfix "github.com/jisinth/kubediag/cmd/fix"
	cmdingress "github.com/jisinth/kubediag/cmd/ingress"
	cmdnetwork "github.com/jisinth/kubediag/cmd/network"
	cmdnodes "github.com/jisinth/kubediag/cmd/nodes"
	cmdpods "github.com/jisinth/kubediag/cmd/pods"
	cmdreport "github.com/jisinth/kubediag/cmd/report"
	cmdsecurity "github.com/jisinth/kubediag/cmd/security"
	cmdstorage "github.com/jisinth/kubediag/cmd/storage"
	cmdversion "github.com/jisinth/kubediag/cmd/version"
	"github.com/jisinth/kubediag/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "kubediag",
	Short: "Automatically diagnose Kubernetes clusters",
	Long: `kubediag connects to a Kubernetes cluster, collects information about
its health, runs diagnostic checks across nodes, pods, deployments, networking,
storage and security, and produces actionable recommendations.`,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.Kubeconfig, "kubeconfig", "", "path to kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)")
	rootCmd.PersistentFlags().StringVar(&config.Context, "context", "", "kubeconfig context to use")
	rootCmd.PersistentFlags().StringVarP(&config.Output, "output", "o", "table", "output format: table, json, html")

	rootCmd.AddCommand(cmddiagnose.NewCommand())
	rootCmd.AddCommand(cmdnodes.NewCommand())
	rootCmd.AddCommand(cmdpods.NewCommand())
	rootCmd.AddCommand(cmdingress.NewCommand())
	rootCmd.AddCommand(cmdstorage.NewCommand())
	rootCmd.AddCommand(cmdnetwork.NewCommand())
	rootCmd.AddCommand(cmdsecurity.NewCommand())
	rootCmd.AddCommand(cmdreport.NewCommand())
	rootCmd.AddCommand(cmdfix.NewCommand())
	rootCmd.AddCommand(cmdversion.NewCommand())
}
