package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	kcpapisv1alpha1 "github.com/kcp-dev/sdk/apis/apis/v1alpha1"
	kcpclientset "github.com/kcp-dev/sdk/client/clientset/versioned"
	"github.com/spf13/cobra"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"sigs.k8s.io/yaml"
)

func newCreateAPIResourceSchema(configFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	var (
		crdPath string
		prefix  string
		output  string
		dryRun  string
	)

	cmd := &cobra.Command{
		Use:          "apiresourceschema --crd <crd.yaml|->",
		Aliases:      []string{"apirs"},
		Short:        "Create an APIResourceSchema from a CRD",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dryRun != "" && dryRun != "none" && dryRun != "client" {
				return fmt.Errorf("invalid value %q for --dry-run; valid values are none, client", dryRun)
			}
			if output != "" && output != "yaml" && output != "json" {
				return fmt.Errorf("invalid value %q for --output; valid values are yaml, json", output)
			}

			crd, err := readCRD(crdPath, streams.In)
			if err != nil {
				return err
			}
			if prefix == "" {
				p, err := stablePrefix(crd)
				if err != nil {
					return fmt.Errorf("error computing stable prefix: %w", err)
				}
				prefix = p
			}

			ars, err := kcpapisv1alpha1.CRDToAPIResourceSchema(crd, prefix)
			if err != nil {
				return fmt.Errorf("converting CRD: %w", err)
			}
			ars.TypeMeta = metav1.TypeMeta{
				APIVersion: kcpapisv1alpha1.SchemeGroupVersion.String(),
				Kind:       "APIResourceSchema",
			}

			if dryRun != "client" {
				cfg, err := configFlags.ToRESTConfig()
				if err != nil {
					return fmt.Errorf("loading REST config: %w", err)
				}
				client, err := kcpclientset.NewForConfig(cfg)
				if err != nil {
					return fmt.Errorf("creating kcp client: %w", err)
				}
				ars, err = client.ApisV1alpha1().APIResourceSchemas().Create(cmd.Context(), ars, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("creating APIResourceSchema: %w", err)
				}
			}

			return printCreated(streams.Out, ars, output, dryRun == "client")
		},
	}
	cmd.Flags().StringVar(&crdPath, "crd", "", "Path to a file containing the CRD to convert, or - for stdin")
	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix for the APIResourceSchema name, before <plural>.<group>. Defaults to a stable hash of the CRD spec.")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output format. Valid values are 'json' and 'yaml'")
	cmd.Flags().StringVar(&dryRun, "dry-run", "none", "If 'client', only print the object that would be created")
	_ = cmd.MarkFlagRequired("crd")

	return cmd
}

// readCRD reads and decodes a CustomResourceDefinition manifest from path, or from in when path is "-".
func readCRD(path string, in io.Reader) (*apiextensionsv1.CustomResourceDefinition, error) {
	if path != "-" {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", path, err)
		}
		defer func() { _ = f.Close() }()
		in = f
	}

	raw, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("reading CRD: %w", err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := yaml.UnmarshalStrict(raw, crd); err != nil {
		return nil, fmt.Errorf("decoding CRD: %w", err)
	}

	wantAPIVersion := apiextensionsv1.SchemeGroupVersion.String()
	if crd.APIVersion != wantAPIVersion || crd.Kind != "CustomResourceDefinition" {
		return nil, fmt.Errorf("expected a %s CustomResourceDefinition, got %s %s", wantAPIVersion, crd.APIVersion, crd.Kind)
	}

	return crd, nil
}

// stablePrefix is stolen from api-syncagent to produce the same stable prefix.
func stablePrefix(crd *apiextensionsv1.CustomResourceDefinition) (string, error) {
	spec := crd.Spec.DeepCopy()
	spec.Conversion = nil

	h := sha256.New()
	if err := json.NewEncoder(h).Encode(spec); err != nil {
		return "", fmt.Errorf("hashing CRD spec: %w", err)
	}

	// prefix with v in case the first char is a digit
	return "v" + hex.EncodeToString(h.Sum(nil))[:8], nil
}

// printCreated prints the kubectl-style creation message, or the object itself when an output format is requested.
func printCreated(out io.Writer, ars *kcpapisv1alpha1.APIResourceSchema, output string, dryRun bool) error {
	switch output {
	case "yaml":
		raw, err := yaml.Marshal(ars)
		if err != nil {
			return fmt.Errorf("encoding APIResourceSchema: %w", err)
		}
		_, err = out.Write(raw)
		return err
	case "json":
		raw, err := json.MarshalIndent(ars, "", "  ")
		if err != nil {
			return fmt.Errorf("encoding APIResourceSchema: %w", err)
		}
		_, err = fmt.Fprintf(out, "%s\n", raw)
		return err
	default:
		suffix := ""
		if dryRun {
			suffix = " (dry run)"
		}
		_, err := fmt.Fprintf(out, "apiresourceschema.apis.kcp.io/%s created%s\n", ars.Name, suffix)
		return err
	}
}
