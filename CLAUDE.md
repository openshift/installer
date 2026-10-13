# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the OpenShift Installer, a tool that deploys OpenShift clusters across multiple cloud platforms (AWS, Azure, GCP, vSphere, bare metal, etc.). The installer generates Ignition configs for bootstrap, control plane, and worker nodes, and can optionally provision the underlying infrastructure.

## Quick Reference Documentation

- **Getting Started**: See [README.md](README.md) for quick start guide
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution workflow, linting, testing, and commit message format
- **Build Dependencies**: See [docs/dev/dependencies.md](docs/dev/dependencies.md) for required system packages and Go version

## Build and Development Commands

### Building the Installer

```sh
# Build the openshift-install binary
hack/build.sh

# Skip Terraform build (faster)
SKIP_TERRAFORM=y hack/build.sh

# Development build (with debugging symbols)
MODE=dev hack/build.sh
```

The binary is output to `bin/openshift-install`.

### Testing

```sh
# Run unit tests
hack/go-test.sh

# Run specific tests with additional arguments
hack/go-test.sh -v -run TestSpecificTest

# Run integration tests
hack/go-integration-test.sh

# Run node joiner integration tests
hack/go-integration-test-nodejoiner.sh
```

### Linting and Formatting

See [CONTRIBUTING.md](CONTRIBUTING.md#contribution-flow) for the complete list of linters to run before submitting a PR. Quick reference:

```sh
# Format Go code and organize imports
hack/go-fmt.sh .

# Run Go linter
hack/go-lint.sh $(go list -f '{{ .ImportPath }}' ./...)

# Run Go vet
hack/go-vet.sh ./...

# Check shell scripts
hack/shellcheck.sh

# Format Terraform files
hack/tf-fmt.sh -list -check

# Lint Terraform
hack/tf-lint.sh

# Lint YAML files
hack/yaml-lint.sh
```

### Generating Code

```sh
# Regenerate mocks for unit tests
hack/go-genmock.sh

# Update install config CRD (after bumping github.com/openshift/api)
go generate ./pkg/types/installconfig.go
```

## Architecture

### Asset-Based Architecture

The installer uses a dependency-graph architecture where everything is an "Asset". See [docs/design/assetgeneration.md](docs/design/assetgeneration.md) for complete details.

Key points:
- **Asset**: Interface with `Dependencies()`, `Generate()`, and `Name()` methods
- **WritableAsset**: Assets that can be written to disk and loaded
- Main assets in `pkg/asset/`: install-config, manifests, ignition-configs, cluster

### Cluster API Integration

The installer uses Cluster API (CAPI) controllers running in a local control plane. See [docs/dev/cluster-api.md](docs/dev/cluster-api.md) for complete details.

Key points:
- Local `kube-apiserver` and `etcd` run via envtest
- Platform-specific infrastructure providers in `cluster-api/providers/`
- Build CAPI binaries with `hack/build-cluster-api.sh` (called automatically by `hack/build.sh`)

### Platform Types

Platform-specific logic lives in `pkg/types/<platform>/`:
- Platform type definitions
- Validation logic in `validation/`
- Default values in `defaults/`

Supported platforms: AWS, Azure, GCP, vSphere, OpenStack, IBM Cloud, Power VS, Nutanix, bare metal.

### Bootstrap Process

The installer creates a temporary bootstrap machine that:
1. Hosts resources for control plane machines to boot
2. Forms initial etcd cluster with control plane nodes
3. Starts temporary Kubernetes control plane
4. Schedules production control plane on control plane machines
5. Injects OpenShift components
6. Shuts down once cluster is self-hosting

## Dependency Management

See [docs/dev/dependencies.md](docs/dev/dependencies.md) for complete dependency management instructions including:
- Adding/updating Go dependencies with `go get`, `go mod tidy`, `go mod vendor`
- Updating CAPI provider dependencies (detailed multi-step process)
- Special case: updating after bumping `github.com/openshift/api`

**Important**: Always commit vendored code in a separate commit from functional changes.

## Commit Message Format

See [CONTRIBUTING.md](CONTRIBUTING.md#commit-message-format) for the complete format specification.

Quick reference:
```
<subsystem>: <what changed>

<why this change was made>

Fixes #<issue-number>
```

Common subsystems:
- `baremetal`, `vsphere`, `aws`, `azure`, `gcp`, etc. - for platform-specific changes
- `agent`, `ibi` (image-based installer) - for installation method changes
- `terraform`, `cluster-api` - for infrastructure provider changes
- `docs` - for documentation changes
- `unit tests`, `integration tests` - for test-only changes (makes it clear the change doesn't affect the actual installer)

## Testing Approach

The installer has different types of tests with varying external requirements:

### Pure Unit Tests

Most tests in `pkg/` are pure unit tests that test Go code logic without external dependencies. These can be run with:

```sh
go test ./pkg/...
```

These tests should pass in any environment with Go installed.

### Integration Tests with External Requirements

Some test files have external dependencies and will fail without specific tools installed:

#### Node Joiner Integration Tests
- **Location**: `cmd/node-joiner/*_integration_test.go`
- **Requirements**:
  - Kubernetes test binaries (etcd, kube-apiserver) via `setup-envtest`
  - Uses `sigs.k8s.io/controller-runtime/pkg/envtest` to run a local control plane
  - The test script automatically downloads the required binaries
- **Run with**: `hack/go-integration-test-nodejoiner.sh` (handles setup automatically)
- **Example test**: `TestNodeJoinerIntegration`
- **Note**: Running `go test` directly without the script will fail with "fork/exec .../etcd: no such file or directory"

#### Agent Integration Tests
- **Location**: `cmd/openshift-install/*_integration_test.go` (tests with "Agent" in name)
- **Requirements**:
  - `oc` binary (OpenShift CLI) in `$PATH`
  - `nmstatectl` binary (for network state validation)
  - Registry credentials for `registry.ci.openshift.org` (for full test pass)
- **Run with**: `hack/go-integration-test.sh`
- **Example test**: `TestAgentIntegration`

**Installing oc**: Download and extract the OpenShift client tools from the official mirror:

```sh
curl -L https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz -o /tmp/oc.tar.gz
mkdir -p ~/.local/bin
tar -xzf /tmp/oc.tar.gz -C ~/.local/bin oc kubectl
rm /tmp/oc.tar.gz
```

Make sure `~/.local/bin` is in your `$PATH`.

**Installing nmstatectl**: Some tests validate network configuration using nmstate. This requires running dnf outside the sandbox:

```sh
sudo dnf install -y nmstate
```

Without `nmstatectl`, network configuration tests will fail with:

```
failed to validate network yaml for host 0, install nmstate package, exec: "nmstatectl": executable file not found in $PATH
```

**Note on registry credentials**: Many agent integration tests query the CI registry (`registry.ci.openshift.org`) to extract release image information. Without credentials, tests will fail with:

```
error: image "registry.ci.openshift.org/origin/release:4.21" not found: manifest unknown: manifest unknown
```

In CI environments, credentials are provided via the `AUTH_FILE` environment variable.

#### General Integration Tests
- **Location**: `test/`
- **Requirements**: Various depending on the test (cloud credentials, network access, etc.)
- **Run with**: `hack/go-integration-test.sh`

### Running Tests

```sh
# Run all unit tests (via podman container with all dependencies)
hack/go-test.sh

# Run unit tests directly (may skip integration tests if dependencies missing)
go test ./...

# Run specific package tests
go test ./pkg/asset/...

# Run integration tests (requires full environment setup)
hack/go-integration-test.sh

# Run node joiner integration tests
hack/go-integration-test-nodejoiner.sh
```

### Test Environment Notes

- **Preferred method**: Use `hack/go-test.sh` which runs tests in a podman container with all dependencies
- **Direct execution**: Running `go test` directly may skip integration tests if tools are missing
- Integration test failures due to missing tools (nmstatectl, kubebuilder, etc.) are expected in minimal environments
- All code in `./cmd/...`, `./data/...`, `./pkg/...` must have unit tests
- Use `hack/go-genmock.sh` to regenerate mocks when interfaces change

## Important Notes

- The installer consumes state from a directory (default: current directory)
- Pass `--dir` to specify asset directory for cluster creation/destruction
- Install config can be pre-created and reused across multiple clusters
