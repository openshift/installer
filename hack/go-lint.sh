#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... pkg/... tests/smoke

podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golangci/golangci-lint:v2.10.1 \
    golangci-lint run -c .golangci-lint-v2.yaml -v --new-from-rev=a6ba91c "${@}"
