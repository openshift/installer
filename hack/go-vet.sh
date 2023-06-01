#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  go vet "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.20 \
    ./hack/go-vet.sh "${@}"
fi;
