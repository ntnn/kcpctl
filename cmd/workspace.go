package cmd

import (
	"errors"

	kcpbase "github.com/kcp-dev/cli/pkg/base"
	kcpplugin "github.com/kcp-dev/cli/pkg/workspace/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// newWorkspace builds kcp's workspace plugin as a cobra.Command.
// Not mounted directly because otherwise `kcpctl --workspace` is ignored.
// Only adding root, use, current and tree because I'm lazy.
func newWorkspace(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	prepare := func(o *kcpbase.Options) {
		prepareBase(o, configFlags, workspace)
	}
	optOut := func(o *kcpbase.Options) {
		o.OptOutOfDefaultKubectlFlags = true
	}

	useOpts := kcpplugin.NewUseWorkspaceOptions(streams)
	optOut(useOpts.Options)
	interactiveTreeOpts := kcpplugin.NewTreeOptions(streams)
	optOut(interactiveTreeOpts.Options)
	var interactive bool

	root := &cobra.Command{
		ValidArgsFunction: completeWorkspacePath(configFlags),
		Aliases:           []string{"ws", "workspaces"},
		Use:               "workspace [use|current|tree|<workspace>|..|.|-|~|<root:absolute:workspace>] [-i|--interactive]",
		Short:             "Manages kcp workspaces",
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				if len(args) != 0 {
					return errors.New("interactive mode does not accept arguments")
				}
				prepare(interactiveTreeOpts.Options)
				interactiveTreeOpts.Interactive = true
				if err := interactiveTreeOpts.Validate(); err != nil {
					return err
				}
				if err := interactiveTreeOpts.Complete(); err != nil {
					return err
				}
				return interactiveTreeOpts.Run(cmd.Context())
			}

			if len(args) != 1 {
				return cmd.Help()
			}
			prepare(useOpts.Options)
			if err := useOpts.Complete(args); err != nil {
				return err
			}
			if err := useOpts.Validate(); err != nil {
				return err
			}
			return useOpts.Run(cmd.Context())
		},
	}
	root.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive workspace tree browser")
	useOpts.BindFlags(root)

	useCmdOpts := kcpplugin.NewUseWorkspaceOptions(streams)
	optOut(useCmdOpts.Options)
	useCmd := &cobra.Command{
		ValidArgsFunction: completeWorkspacePath(configFlags),
		Aliases:           []string{"cd"},
		Use:               "use <workspace>|..|.|-|~|<:root:absolute:workspace>|<relative:workspace>",
		Short:             "Uses the given workspace as the current workspace. Using - means previous workspace, .. means parent workspace, . means current, ~ means home workspace",
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}
			prepare(useCmdOpts.Options)
			if err := useCmdOpts.Complete(args); err != nil {
				return err
			}
			if err := useCmdOpts.Validate(); err != nil {
				return err
			}
			return useCmdOpts.Run(cmd.Context())
		},
	}
	useCmdOpts.BindFlags(useCmd)
	root.AddCommand(useCmd)

	currentOpts := kcpplugin.NewCurrentWorkspaceOptions(streams)
	optOut(currentOpts.Options)
	currentCmd := &cobra.Command{
		Use:          "current [--short]",
		Short:        "Print the current workspace. Same as 'kcpctl ws .'.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return cmd.Help()
			}
			prepare(currentOpts.Options)
			if err := currentOpts.Validate(); err != nil {
				return err
			}
			if err := currentOpts.Complete(); err != nil {
				return err
			}
			return currentOpts.Run(cmd.Context())
		},
	}
	currentOpts.BindFlags(currentCmd)
	root.AddCommand(currentCmd)

	treeOpts := kcpplugin.NewTreeOptions(streams)
	optOut(treeOpts.Options)
	treeCmd := &cobra.Command{
		Use:          "tree",
		Short:        "Print the current workspace tree.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return cmd.Help()
			}
			prepare(treeOpts.Options)
			if err := treeOpts.Validate(); err != nil {
				return err
			}
			if err := treeOpts.Complete(); err != nil {
				return err
			}
			return treeOpts.Run(cmd.Context())
		},
	}
	treeOpts.BindFlags(treeCmd)
	root.AddCommand(treeCmd)

	return root
}
