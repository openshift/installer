---
name: nodejoiner-command
description: Describes the acceptance criteria, BDD tests and specific implementation details when developing a new feature (or fixing an issue) related to the node-joiner command.
user-invocable: false
---

# Implementation details
The nodejoiner is a internal command (not meant to be used directly by the user) wrapped by the `oc adm nodeimage` command, to allow the user adding a new node to an existing target cluster by building an ISO. The command reuses most of the code of the Agent-based Installer, but since it follows a slightly different workflow
some differences have been implemented in the code (through the usage of the `workflow.AgentWorkflowTypeAddNodes` enum value).

The most important asset is the `joiner.ClusterInfo`, as it contains the logic to inspect the current target cluster and to extract the input configuration required by
the ABI assets to generate the add node ISO. So usually the target cluster is considered the authoritative source of truth for the configuration.

# Testing

- To run just the unit tests, you must use `hack/go-test.sh`
- To run the specific node-joiner integration test, you must use only `hack/go-integration-test-nodejoiner.sh`
- **Integration test prerequisites**:
  - `oc` must be installed and in `$PATH`. Install with:
    `curl -L https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz -o /tmp/oc.tar.gz && mkdir -p ~/.local/bin && tar -xzf /tmp/oc.tar.gz -C ~/.local/bin oc kubectl && rm /tmp/oc.tar.gz`
  - Override `GOPATH` to a writable location (the default `/tmp/claude/go` is not writable in containers):
    `PATH=$HOME/.local/bin:$PATH GOPATH=$HOME/go hack/go-integration-test-nodejoiner.sh`

# Useful references
- https://github.com/openshift/enhancements/blob/master/enhancements/oc/day2-add-nodes.md. This the official `oc` enhancement proposal to 
  add a new `oc adm nodeimage` command for letting the user easily add a new node by building an ISO. The `oc` command is a thin wrapper of the `node-joiner` command.
- https://github.com/openshift/oc/tree/main/pkg/cli/admin/nodeimage. The code repository for the `oc adm nodeimage` commands.

# Supporting files
- `acceptance/coreos-bootimages-data.md`: how to fetch the coreos bootimages data required for building the ISO.