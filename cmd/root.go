// Package cmd assembles the kcpctl command tree.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/rest"
	kubectlcmd "k8s.io/kubectl/pkg/cmd"

	"github.com/ntnn/kcpctl/pkg/kcpurl"
)

// New returns the kcpctl root command writing to streams.
func New(streams genericiooptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true).WithWarningPrinter(streams)

	// wokrspaceHost is compute in PersistentPreRunE because WrapConfigFn cannot return errors.
	var workspaceHost string
	configFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
		if workspaceHost != "" {
			c.Host = workspaceHost
		}
		return c
	}

	root := kubectlcmd.NewKubectlCommand(kubectlcmd.KubectlOptions{
		Arguments:   os.Args,
		ConfigFlags: configFlags,
		IOStreams:   streams,
	})
	root.Use = "kcpctl"
	root.Short = "kcpctl controls kcp"
	root.Long = "kcpctl is a kubectl-like CLI for kcp."

	// -W isntead -w as that is taken by kubectl get --watch
	var workspace string
	root.PersistentFlags().StringVarP(&workspace, "workspace", "W", "", "Workspace path to target.")
	_ = root.RegisterFlagCompletionFunc("workspace", completeWorkspacePath(configFlags))

	for _, c := range root.Commands() {
		switch c.Name() {
		case "create":
			c.AddCommand(newCreateWorkspace(configFlags, &workspace, streams))
			c.AddCommand(newCreateAPIResourceSchema(configFlags, streams))
		case "get":
			wireGetClaims(c, configFlags, &workspace, streams)
		}
	}

	root.AddCommand(newWorkspace(configFlags, &workspace, streams))
	root.AddCommand(newBind(configFlags, &workspace, streams))
	root.AddCommand(newAccept(configFlags, &workspace, streams))
	root.AddCommand(newReject(configFlags, &workspace, streams))

	kubectlPreRunE := root.PersistentPreRunE
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// shell completion must not trip over half-typed flag values
		if cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd {
			return kubectlPreRunE(cmd, args)
		}
		if workspace != "" {
			if configFlags.APIServer != nil && *configFlags.APIServer != "" {
				return errors.New("--workspace cannot be combined with --server")
			}
			cfg, err := configFlags.ToRawKubeConfigLoader().ClientConfig()
			if err != nil {
				return fmt.Errorf("loading kubeconfig: %w", err)
			}
			host, err := kcpurl.WithWorkspace(cfg.Host, workspace)
			if err != nil {
				return err
			}
			workspaceHost = host
		}
		return kubectlPreRunE(cmd, args)
	}

	return root
}

// Main runs kcpctl with OS streams and exits nonzero on error.
func Main() {
	streams := genericiooptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	if err := New(streams).Execute(); err != nil {
		os.Exit(1)
	}
}
