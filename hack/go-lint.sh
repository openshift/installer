#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... pkg/... tests/smoke

if [ "$IS_CONTAINER" != "" ]; then
  GOFLAGS="" go install golang.org/x/lint/golint@latest
  "${GOPATH}"/bin/golint -set_exit_status "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.18 \
    ./hack/go-lint.sh "${@}"
fi
