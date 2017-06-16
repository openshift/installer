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

The smoke tests require two parameters: the file path of the cluster kubeconfig and the number of nodes in the cluster.
Export the following variables to parameterize the smoke tests:

```
export TEST_KUBECONFIG=/path/to/kubeconfig
export NODE_COUNT=3
```

Tests can then be run using the `go test` tool:

```
go test -v github.com/coreos/tectonic-installer/installer/tests/sanity
```
