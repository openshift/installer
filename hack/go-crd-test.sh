#!/bin/sh
# Run CRD schema and behavioral validation tests.
# These tests require envtest (kube-apiserver + etcd binaries).
# Example:  ./hack/go-crd-test.sh

set -e

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

go install -mod=mod sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.19
# shellcheck disable=SC2086
KUBEBUILDER_ASSETS="$($(go env GOPATH)/bin/setup-envtest use 1.31.0 -p path --bin-dir /tmp)" \
    go test -timeout 5m -v -C "${REPO_ROOT}/tests/crd" ./... "${@}"
