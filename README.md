# kcpctl

Experimental shim around kubectl and kcp plugins with some sugar.
Basically embeds kcp functionality into the kubectl cobra tree.

Adds a global `--workspace`/`-W` flag to target a specific kcp workspace, supports relative paths.
Adds completion work workspaces to both the workspace flag and the `kcpctl workspace` command.

"kubectl-native" Workspace and APIResourceSchema support:

```sh
kcpctl create workspace <name>
kcpctl create apiresourceschema --crd <crd.yaml>
```

The prefix for APIResourceSchema can be set explicitly and is otherwise a stable hash.

Otherwise embeds the kcp plugins:

```sh
kcpctl ws <path>
kcpctl ws tree
kcpctl bind apiexport <path>:<name>
kcpctl get apibinding <binding> --claims
kcpctl accept claims <binding> --resource ...
kcpctl reject claims <binding> --resource ...
```
