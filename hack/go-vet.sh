#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  go vet "${@}"
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
    ./hack/go-vet.sh "${@}"
fi;
