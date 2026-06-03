# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

@AGENTS.md

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
