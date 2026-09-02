package kcpurl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kcp-dev/logicalcluster/v3"
)

const clustersSegment = "/clusters/"

// WithWorkspace rewrites base to target workspace.
// workspace may be an absolute or relative workspace path.
func WithWorkspace(base, workspace string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parsing server URL %q: %w", base, err)
	}

	prefix, current := splitClusters(u.Path)
	resolved, err := resolve(current, workspace)
	if err != nil {
		return "", err
	}
	u.Path = prefix + resolved.RequestPath()
	return u.String(), nil
}

// splitClusters splits path into everything before "/clusters/" and after it.
func splitClusters(path string) (string, logicalcluster.Path) {
	prefix, cluster, found := strings.Cut(path, clustersSegment)
	if !found {
		return strings.TrimSuffix(path, "/"), logicalcluster.Path{}
	}
	if i := strings.Index(cluster, "/"); i >= 0 {
		cluster = cluster[:i]
	}
	return prefix, logicalcluster.NewPath(cluster)
}

// resolve turns input into an absolute workspace path, resolving relative
// input against current.
func resolve(current logicalcluster.Path, input string) (logicalcluster.Path, error) {
	if absolute, ok := strings.CutPrefix(input, ":"); ok {
		p, valid := logicalcluster.NewValidatedPath(absolute)
		if !valid {
			return logicalcluster.Path{}, fmt.Errorf("invalid workspace path %q", input)
		}
		return p, nil
	}

	if current.Empty() {
		return logicalcluster.Path{}, fmt.Errorf("relative workspace path %q requires the current server URL to contain a workspace path", input)
	}

	p := current
	for segment := range strings.SplitSeq(input, ":") {
		switch segment {
		case ".":
		case "..":
			parent, ok := p.Parent()
			if !ok {
				return logicalcluster.Path{}, fmt.Errorf("workspace path %q has no parent", p)
			}
			p = parent
		default:
			p = p.Join(segment)
		}
	}
	if !p.IsValid() {
		return logicalcluster.Path{}, fmt.Errorf("invalid workspace path %q", p)
	}
	return p, nil
}
