package cmd

import (
	"errors"

	kcpplugin "github.com/kcp-dev/cli/pkg/claims/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// wireGetClaims adds --claims for `kcpctl get apibinding`
func wireGetClaims(getCmd *cobra.Command, configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) {
	opts := kcpplugin.NewGetAPIBindingOptions(streams)
	opts.OptOutOfDefaultKubectlFlags = true

	var showClaims bool
	getCmd.Flags().BoolVar(&showClaims, "claims", false, "Print the permission claims of the APIBinding(s). Only valid for apibindings.")

	kubectlRun := getCmd.Run
	getCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !showClaims {
			kubectlRun(cmd, args)
			return nil
		}

		if len(args) == 0 || (args[0] != "apibinding" && args[0] != "apibindings") {
			return errors.New("--claims is only valid for apibindings, e.g. kcpctl get apibindings --claims")
		}
		prepareBase(opts.Options, configFlags, workspace)
		if err := opts.Complete(args[1:]); err != nil {
			return err
		}
		if err := opts.Validate(); err != nil {
			return err
		}
		return opts.Run(cmd.Context())
	}
}

func newAccept(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	return newAcceptOrReject("accept", kcpplugin.NewClaimsAcceptOptions(streams), configFlags, workspace)
}

func newReject(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	return newAcceptOrReject("reject", kcpplugin.NewClaimsRejectOptions(streams), configFlags, workspace)
}

func newAcceptOrReject(action string, opts *kcpplugin.ClaimsAcceptOrRejectOptions, configFlags *genericclioptions.ConfigFlags, workspace *string) *cobra.Command {
	opts.OptOutOfDefaultKubectlFlags = true

	root := &cobra.Command{
		Use:          action,
		Short:        action + " kcp objects",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	claimsCmd := &cobra.Command{
		Use:          "claims <apibinding_name>",
		Short:        action + " permission claims of an APIBinding",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(opts.Options, configFlags, workspace)
			if err := opts.Complete(args); err != nil {
				return err
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			return opts.Run(cmd.Context())
		},
	}
	opts.BindFlags(claimsCmd)
	root.AddCommand(claimsCmd)

	return root
}
