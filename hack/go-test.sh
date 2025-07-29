#!/bin/sh
# Example:  ./hack/go-test.sh

if [ "$IS_CONTAINER" != "" ]; then
  go test -short ./cmd/... ./data/... ./pkg/... "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.24 \
    ./hack/go-test.sh "${@}"
fi
