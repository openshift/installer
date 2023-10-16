#!/bin/sh
# Example:  ./hack/go-integration-test.sh

go test -parallel 1 -p 1 -timeout 0 -run .Integration ./cmd/... ./data/... ./pkg/... "${@}"
