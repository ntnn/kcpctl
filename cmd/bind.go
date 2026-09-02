package cmd

import (
	kcpplugin "github.com/kcp-dev/cli/pkg/bind/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func newBind(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	root := &cobra.Command{
		Use:          "bind",
		Short:        "Bind different types into current workspace.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	bindOpts := kcpplugin.NewBindOptions(streams)
	bindOpts.OptOutOfDefaultKubectlFlags = true
	bindCmd := &cobra.Command{
		Use:          "apiexport <workspace_path:apiexport-name>",
		Short:        "Bind to an APIExport",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(bindOpts.Options, configFlags, workspace)
			if err := bindOpts.Complete(args); err != nil {
				return err
			}
			if err := bindOpts.Validate(); err != nil {
				return err
			}
			return bindOpts.Run(cmd.Context())
		},
	}
	bindOpts.BindFlags(bindCmd)
	root.AddCommand(bindCmd)

	return root
}
