package integration

import (
	"os"
	"path/filepath"
	"testing"

	kcptc "github.com/ntnn/kcp-testcontainer"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/ntnn/kcpctl/cmd"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"kcpctl": cmd.Main,
	})
}

func TestScript(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	inst, err := kcptc.Single(ctx, "")
	require.NoError(t, err)
	tc.CleanupContainer(t, inst)

	raw, err := inst.Kubeconfig(ctx)
	require.NoError(t, err)

	testscript.Run(t, testscript.Params{
		Dir:                 "testdata",
		RequireExplicitExec: true,
		UpdateScripts:       os.Getenv("UPDATE") != "",
		Setup: func(env *testscript.Env) error {
			path := filepath.Join(env.WorkDir, "admin.kubeconfig")
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				return err
			}
			env.Setenv("KUBECONFIG", path)
			return nil
		},
	})
}
