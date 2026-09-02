package cmd

import (
	"context"
	"strings"
	"time"

	kcptenancyv1alpha1 "github.com/kcp-dev/sdk/apis/tenancy/v1alpha1"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ntnn/kcpctl/pkg/kcpurl"
)

func init() {
	utilruntime.Must(kcptenancyv1alpha1.AddToScheme(scheme.Scheme))
}

// completeWorkspacePath provides completion for workspaces
func completeWorkspacePath(configFlags *genericclioptions.ConfigFlags) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		parent := "."
		// toComplete=:root:my:work:space
		if i := strings.LastIndex(toComplete, ":"); i >= 0 {
			// parent=:root:my:work
			parent = toComplete[:i]
			if parent == "" {
				// only root sits at the top of the tree
				return []string{":root"}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cfg, err := configFlags.ToRawKubeConfigLoader().ClientConfig()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		cfg = rest.CopyConfig(cfg)
		cfg.Host, err = kcpurl.WithWorkspace(cfg.Host, parent)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		cl, err := client.New(cfg, client.Options{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		list := &kcptenancyv1alpha1.WorkspaceList{}
		if err := cl.List(ctx, list); err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		candidates := make([]string, 0, len(list.Items))
		for _, ws := range list.Items {
			candidates = append(candidates, parent+":"+ws.Name)
		}
		return candidates, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
