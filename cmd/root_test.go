package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestWorkspaceServerFlagConflict(t *testing.T) {
	t.Parallel()

	root := New(discardStreams())
	root.SetArgs([]string{"get", "configmaps", "-W", ":root", "--server", "https://example.com"})
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)

	require.ErrorContains(t, root.Execute(), "--workspace cannot be combined with --server")
}

func discardStreams() genericiooptions.IOStreams {
	return genericiooptions.IOStreams{
		In:     &bytes.Buffer{},
		Out:    io.Discard,
		ErrOut: io.Discard,
	}
}
