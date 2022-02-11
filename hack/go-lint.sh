#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... pkg/... tests/smoke

if [ "$IS_CONTAINER" != "" ]; then
  golint -set_exit_status "${@}"
else
  ENGINE="podman"
  if [ "$(uname)" = "Darwin" ]; then
    ENGINE="docker"
  fi
  "$ENGINE" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.16 \
    ./hack/go-lint.sh "${@}"
fi
