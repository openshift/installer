#!/bin/sh
# Example:  ./hack/go-genmock.sh

if [ "$IS_CONTAINER" != "" ]; then
  go generate ./pkg/asset/installconfig/... "${@}"
else
  podman build -t openshift-install-mock ./images/mock
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    openshift-install-mock \
    ./hack/go-genmock.sh "${@}"
fi
