package kcpurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithWorkspace(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		base      string
		workspace string
		want      string
		wantErr   string
	}{
		"absolute replaces existing cluster": {
			base:      "https://kcp.example.com:6443/clusters/root",
			workspace: ":root:team-a",
			want:      "https://kcp.example.com:6443/clusters/root:team-a",
		},
		"absolute appended without clusters segment": {
			base:      "https://kcp.example.com:6443",
			workspace: ":root",
			want:      "https://kcp.example.com:6443/clusters/root",
		},
		"trailing slash without clusters segment": {
			base:      "https://kcp.example.com:6443/",
			workspace: ":root",
			want:      "https://kcp.example.com:6443/clusters/root",
		},
		"path prefix before clusters preserved": {
			base:      "https://front.example.com/proxy/clusters/root",
			workspace: ":root:team-a",
			want:      "https://front.example.com/proxy/clusters/root:team-a",
		},
		"relative child": {
			base:      "https://kcp.example.com:6443/clusters/root",
			workspace: "team-a",
			want:      "https://kcp.example.com:6443/clusters/root:team-a",
		},
		"relative parent": {
			base:      "https://kcp.example.com:6443/clusters/root:a:b",
			workspace: "..",
			want:      "https://kcp.example.com:6443/clusters/root:a",
		},
		"relative current": {
			base:      "https://kcp.example.com:6443/clusters/root:a",
			workspace: ".",
			want:      "https://kcp.example.com:6443/clusters/root:a",
		},
		"chained relative": {
			base:      "https://kcp.example.com:6443/clusters/root:a:b",
			workspace: "..:..:c",
			want:      "https://kcp.example.com:6443/clusters/root:c",
		},
		"relative without clusters segment": {
			base:      "https://kcp.example.com:6443",
			workspace: "..",
			wantErr:   "requires the current server URL to contain a workspace path",
		},
		"parent of top-level workspace": {
			base:      "https://kcp.example.com:6443/clusters/root",
			workspace: "..",
			wantErr:   "has no parent",
		},
		"invalid absolute path": {
			base:      "https://kcp.example.com:6443/clusters/root",
			workspace: ":Team_A",
			wantErr:   "invalid workspace path",
		},
		"invalid relative result": {
			base:      "https://kcp.example.com:6443/clusters/root",
			workspace: "Team_A",
			wantErr:   "invalid workspace path",
		},
		"unparsable base URL": {
			base:      "https://kcp.example.com:bad-port",
			workspace: ":root",
			wantErr:   "parsing server URL",
		},
	}

	for title, cas := range cases {
		t.Run(title, func(t *testing.T) {
			t.Parallel()

			got, err := WithWorkspace(cas.base, cas.workspace)
			if cas.wantErr != "" {
				require.ErrorContains(t, err, cas.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, cas.want, got)
		})
	}
}
