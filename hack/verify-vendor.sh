#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  if [ ! "$(command -v dep >/dev/null)" ]; then
    go get -u github.com/golang/dep/cmd/dep
  fi

  dep check
  dep ensure
  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.12 \
    ./hack/verify-vendor.sh "${@}"
fi