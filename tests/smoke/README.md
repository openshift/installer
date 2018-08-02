# Tectonic Smoke Tests

This directory contains all smoke tests for Tectonic.
The smoke tests are a set of Golang test files that perform minimal validation of a running Tectonic cluster.
This directory is further partitioned into platform-specific directories and should conform to the following directory hierarchy:

```
smoke/
├── aws        # Smoke tests for AWS
│   └── vars   # Terraform tfvars files for AWS smoke tests
├── *_test.go  # Smoke tests for all platforms
├── vendor     # Smoke test dependencies
└── ...
```

## Getting Started

The smoke tests assume a running Tectonic cluster, so before running any tests:
1. create a Tectonic cluster; and
2. download the cluster's kubeconfig to a known location.

## Running

The smoke tests only require one parameter: the file path of the cluster kubeconfig.
Export the following variable to parameterize the smoke tests:

```sh
export SMOKE_KUBECONFIG=/path/to/kubeconfig
```

Compile the smoke test binary from the root directory of the project:

```sh
bazel build tests/smoke
```

The tests can then be run by invoking the `smoke` binary in the `bazel-bin/tests/smoke/linux_amd64_stripped` directory.
This binary accepts the `--cluster` flag to specify which tests suites should be run, e.g.:

```sh
bazel-bin/tests/smoke/linux_amd64_stripped/smoke --cluster
```

*Note*: the `smoke` binary accepts several flags available to the `go test` command; to list them, invoke the `smoke` binary with the `--help` flag.
For example, to run the cluster suite verbosely, use the `--test.v` flag:

```sh
bazel-bin/tests/smoke/linux_amd64_stripped/smoke --cluster --test.v
```

## Cluster Suite

The cluster test suite runs a series of tests designed to verify the overall health of the cluster and ensure that all expected components are present.
The cluster test suite requires four additional parameters:

* whether or not Calico network policy support is enabled;
* the number of nodes in the cluster;
* the paths of manifests deployed on the cluster; and
* whether or not to test for experimental manifests.

Export the following environment variables to parameterize the cluster tests:

```sh
export SMOKE_NETWORKING=canal
export SMOKE_NODE_COUNT=3
export SMOKE_MANIFEST_PATHS=/path/to/kubernetes/manifests
export SMOKE_MANIFEST_EXPERIMENTAL=true
```

To run the cluster test suite, invoke the smoke test binary with the `--cluster` flag:
```sh
bazel-bin/tests/smoke/linux_amd64_stripped/smoke --cluster
```
