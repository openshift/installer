#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  set -x
  go generate ./pkg/types/installconfig.go
  go generate ./pkg/rhcos/ami.go
  set +x
  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.13 \
    ./hack/verify-codegen.sh "${@}"
fi
