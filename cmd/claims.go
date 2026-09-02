package cmd

import (
	kcpplugin "github.com/kcp-dev/cli/pkg/claims/plugin"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func newClaims(configFlags *genericclioptions.ConfigFlags, workspace *string, streams genericiooptions.IOStreams) *cobra.Command {
	root := &cobra.Command{
		Use:          "claims",
		Short:        "Operations related to viewing or updating permission claims",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	getCmd := &cobra.Command{
		Use:          "get",
		Short:        "Operations related to fetching APIs with respect to permission claims",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	getBindingOpts := kcpplugin.NewGetAPIBindingOptions(streams)
	getBindingOpts.OptOutOfDefaultKubectlFlags = true
	getBindingCmd := &cobra.Command{
		Use:          "apibinding <apibinding_name>",
		Short:        "Get claims related to apibinding",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(getBindingOpts.Options, configFlags, workspace)
			if err := getBindingOpts.Complete(args); err != nil {
				return err
			}
			if err := getBindingOpts.Validate(); err != nil {
				return err
			}
			return getBindingOpts.Run(cmd.Context())
		},
	}
	getBindingOpts.BindFlags(getBindingCmd)
	getCmd.AddCommand(getBindingCmd)
	root.AddCommand(getCmd)

	acceptOpts := kcpplugin.NewClaimsAcceptOptions(streams)
	acceptOpts.OptOutOfDefaultKubectlFlags = true
	acceptCmd := &cobra.Command{
		Use:          "accept <apibinding_name>",
		Short:        "Accept permission claims of an APIBinding",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(acceptOpts.Options, configFlags, workspace)
			if err := acceptOpts.Complete(args); err != nil {
				return err
			}
			if err := acceptOpts.Validate(); err != nil {
				return err
			}
			return acceptOpts.Run(cmd.Context())
		},
	}
	acceptOpts.BindFlags(acceptCmd)
	root.AddCommand(acceptCmd)

	rejectOpts := kcpplugin.NewClaimsRejectOptions(streams)
	rejectOpts.OptOutOfDefaultKubectlFlags = true
	rejectCmd := &cobra.Command{
		Use:          "reject <apibinding_name>",
		Short:        "Reject permission claims of an APIBinding",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			prepareBase(rejectOpts.Options, configFlags, workspace)
			if err := rejectOpts.Complete(args); err != nil {
				return err
			}
			if err := rejectOpts.Validate(); err != nil {
				return err
			}
			return rejectOpts.Run(cmd.Context())
		},
	}
	rejectOpts.BindFlags(rejectCmd)
	root.AddCommand(rejectCmd)

	return root
}
