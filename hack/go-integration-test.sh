#!/bin/sh
# Example:  ./hack/go-integration-test.sh

go test -run .Integration ./cmd/... ./data/... ./pkg/... "${@}"
