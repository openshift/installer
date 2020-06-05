#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  go mod vendor
  go mod verify
  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.14 \
    ./hack/verify-vendor.sh "${@}"
fi
