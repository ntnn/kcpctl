package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestNew(t *testing.T) {
	t.Parallel()

	root := New(genericiooptions.IOStreams{
		In:  &bytes.Buffer{},
		Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{},
	})
	require.NotNil(t, root)

	get, _, err := root.Find([]string{"get"})
	require.NoError(t, err, "embedded kubectl tree must expose get")

	assert.Equal(t, "kcpctl", root.Name())
	assert.Equal(t, "get", get.Name())
}
