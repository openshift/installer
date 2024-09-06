#!/bin/sh
# Example:  ./hack/go-integration-test-nodejoiner.sh

go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
export KUBEBUILDER_ASSETS="$($GOPATH/bin/setup-envtest use 1.31.0 -p path)" 
go test -parallel 1 -p 1 -timeout 0 -run .Integration ./cmd/node-joiner/... "${@}"

