#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... pkg/... tests/smoke

podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golangci/golangci-lint:v2.3.1 \
    golangci-lint run -v --new-from-rev=dcf8122 "${@}"
