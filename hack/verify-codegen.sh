#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  set -xe
  go generate ./pkg/types/installconfig.go
  set +ex
  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.23 \
    ./hack/verify-codegen.sh "${@}"
fi
