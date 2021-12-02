#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  set -xe
  go generate ./pkg/types/installconfig.go
  # See https://github.com/openshift/installer/pull/5447#discussion_r762340594
  # go generate ./pkg/rhcos/ami.go
  set +ex
  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.16 \
    ./hack/verify-codegen.sh "${@}"
fi
