# Tectonic Smoke Tests

This directory contains all smoke tests for Tectonic.
The smoke tests are a set of Golang test files that perform minimal validation of a running Tectonic cluster.
This directory is further partitioned into platform-specific directories and should conform to the following directory hierarchy:

```
smoke/
├── aws        # Smoke tests for AWS
│   └── vars   # Terraform tfvars files for AWS smoke tests
├── azure
│   └── vars   # Terraform tfvars files for Azure smoke tests
├── bare-metal
│   └── vars   # Terraform tfvars files for bare-metal smoke tests
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
make bin/smoke
```

The tests can then be run by invoking the `smoke` binary in the `bin` directory.
This binary accepts `--cluster` and `--qa` flags to specify which tests suites should be run, e.g.:

```sh
bin/smoke --cluster --qa
```

*Note*: the `smoke` binary accepts several flags available to the `go test` command; to list them, invoke the `smoke` binary with the `--help` flag.
For example, to run the cluster and QA test suites verbosely, use the `--test.v` flag:

```sh
bin/smoke --cluster --qa --test.v
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
export SMOKE_NETWORKING=flannel
export SMOKE_NODE_COUNT=3
export SMOKE_MANIFEST_PATHS=/path/to/kubernetes/manifests
export SMOKE_MANIFEST_EXPERIMENTAL=true
```

To run the cluster test suite, invoke the smoke test binary with the `--cluster` flag:
```sh
bin/smoke --cluster
```

## QA Suite

The QA test suite is designed to automate several QA checklist items that verify that Tectonic cluster integrations work as expected.
The QA test suite requires two additional parameters:

* the spec of the BigQuery database to test for cluster metrics; and
* the path to a file containing Google application credentials with job permissions on BigQuery database.

Export the following environment variables to parameterize the QA tests:

```sh
export SMOKE_BIGQUERY_SPEC=bigquery://<project>.<dataset>.<table>
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/google/application/credentials
```

To run the QA suite, invoke the smoke test binary with the `--qa` flag:
```sh
bin/smoke --qa
```
