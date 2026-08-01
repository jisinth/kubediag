// Package cluster handles connecting to a Kubernetes cluster and running
// basic connectivity/health checks against the API server.
package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps a connected Kubernetes clientset along with the REST config
// used to build it.
type Client struct {
	Clientset *kubernetes.Clientset
	Config    *rest.Config
}

// Connect loads the kubeconfig (explicit path, $KUBECONFIG, or
// ~/.kube/config, in that order) and returns a connected Client.
func Connect(kubeconfigPath, kubeContext string) (*Client, error) {
	resolved, err := resolveKubeconfigPath(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: resolved}
	overrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}

	restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig from %s: %w", resolved, err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Client{Clientset: clientset, Config: restConfig}, nil
}

func resolveKubeconfigPath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".kube", "config"), nil
}

// Version returns the Kubernetes API server's version string, verifying
// basic connectivity to the API server in the process.
func (c *Client) Version(ctx context.Context) (string, error) {
	v, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
