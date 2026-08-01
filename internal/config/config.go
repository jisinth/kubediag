// Package config holds CLI-wide settings populated from persistent flags on
// the root command and read by every subcommand.
package config

var (
	// Kubeconfig is the path to the kubeconfig file. Empty means fall back to
	// $KUBECONFIG or ~/.kube/config.
	Kubeconfig string

	// Context is the kubeconfig context to use. Empty means use the current
	// context.
	Context string

	// Output is the report format: table, json, or html.
	Output string
)
