package cmd

import (
	kcpplugin "github.com/kcp-dev/cli/pkg/workspace/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// newCreateWorkspace wraps kcp's create-workspace plugin as `kcpctl create workspace`.
func newCreateWorkspace(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	kcpopts := kcpplugin.NewCreateWorkspaceOptions(streams)
	kcpopts.OptOutOfDefaultKubectlFlags = true

	cmd := &cobra.Command{
		Use:          "workspace <name>",
		Aliases:      []string{"ws"},
		Short:        "Create a kcp workspace",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(kcpopts.Options, configFlags, workspace)
			if err := kcpopts.Validate(); err != nil {
				return err
			}
			if err := kcpopts.Complete(args); err != nil {
				return err
			}
			return kcpopts.Run(cmd.Context())
		},
	}
	kcpopts.BindFlags(cmd)

	return cmd
}
