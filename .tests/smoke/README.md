# OpenShift Installer Smoke Tests

This directory contains all smoke tests for OpenShift Installer.
The smoke tests are a set of Golang test files that perform minimal validation of a running OpenShift cluster.

## Getting Started

The smoke tests assume a running OpenShift cluster, so before running any tests:
1. create a OpenShift cluster; and
2. download the cluster's kubeconfig to a known location.

## Running

The smoke tests only require one parameter: the file path of the cluster kubeconfig.
Export the following variable to parameterize the smoke tests:

```sh
export SMOKE_KUBECONFIG=/path/to/kubeconfig
```

Compile the smoke test binary from the root directory of the project:

```sh
bazel build smoke_tests
```

The tests can then be run by invoking the `go_default_test` binary in the `bazel-bin/tests/smoke/linux_amd64_pure_stripped` directory.

*Note*: the `go_default_test` binary accepts several flags available to the `go test` command; to list them, invoke the `go_default_test` binary with the `-help` flag.

## Cluster Suite

The cluster test suite runs a series of tests designed to verify the overall health of the cluster and ensure that all expected components are present.
The cluster test suite requires additional parameters:

* the number of nodes in the cluster;
* the paths of manifests deployed on the cluster; and
* whether or not to test for experimental manifests.

Export the following environment variables to parameterize the cluster tests:

```sh
export SMOKE_NODE_COUNT=3
export SMOKE_MANIFEST_PATHS=/path/to/kubernetes/manifests
export SMOKE_MANIFEST_EXPERIMENTAL=true
```

To run the cluster test suite, invoke the smoke test binary with the `-cluster` flag.
To run the cluster suite verbosely, add the `-test.v` flag:

```sh
bazel-bin/tests/smoke/linux_amd64_pure_stripped/go_default_test -cluster -test.v
```
