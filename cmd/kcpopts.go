package cmd

import (
	kcpbase "github.com/kcp-dev/cli/pkg/base"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// prepareBase wires kubectl flags into kcp plugin options.
func prepareBase(kcpopts *kcpbase.Options, configFlags *genericclioptions.ConfigFlags, workspace *string) {
	kcpopts.ClientConfig = configFlags.ToRawKubeConfigLoader()
	kcpopts.Workspace = *workspace
}
