# Tectonic Sanity Tests

Tectonic sanity tests perform minimal validation on a running AWS cluster. They require you've check out the "tectonic" project in a valid GOPATH.

```
git clone git@github.com:coreos/tectonic-installer $GOPATH/src/github.com/coreos/tectonic-installer
```

Arguments are defined as environment variables (and all required):

```
export TEST_KUBECONFIG=/path/to/kubeconfig
export NODE_COUNT=3
```

Tests can then be run using `go test`.

```
# Install the large number of Kubernetes test dependencies before running.
go test -v -i github.com/coreos/tectonic-installer/installer/tests/sanity
go test -v github.com/coreos/tectonic-installer/installer/tests/sanity
```
