#!/bin/sh
# Example:  ./hack/go-test.sh

GIT_TAG="${BUILD_VERSION:-$(git describe --always --abbrev=40 --dirty)}"
DEFAULT_ARCH="${DEFAULT_ARCH:-amd64}"
LDFLAGS="${LDFLAGS} -X github.com/openshift/installer/pkg/version.Raw=${GIT_TAG} -X github.com/openshift/installer/pkg/version.defaultArch=${DEFAULT_ARCH}"

if [ "$IS_CONTAINER" != "" ]; then
  go test -ldflags "${LDFLAGS}" -short ./cmd/... ./data/... ./pkg/... "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --env LDFLAGS="${LDFLAGS}" \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.25 \
    ./hack/go-test.sh "${@}"
fi
