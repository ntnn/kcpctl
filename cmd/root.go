// Package cmd assembles the kcpctl command tree.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	kubectlcmd "k8s.io/kubectl/pkg/cmd"
)

// New returns the kcpctl root command writing to streams.
func New(streams genericiooptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true).WithWarningPrinter(streams)

	root := kubectlcmd.NewKubectlCommand(kubectlcmd.KubectlOptions{
		Arguments:   os.Args,
		ConfigFlags: configFlags,
		IOStreams:   streams,
	})
	root.Use = "kcpctl"
	root.Short = "kcpctl controls kcp"
	root.Long = "kcpctl is a kubectl-like CLI for kcp."

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
