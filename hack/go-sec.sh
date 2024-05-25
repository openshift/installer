#!/bin/sh
# Example:  ./hack/gosec.sh
set -x

if [ "$IS_CONTAINER" != "" ]; then
  if [ ! "$(command -v gosec >/dev/null)" ]; then
      go get github.com/securego/gosec/cmd/gosec
  fi
  gosec -severity high -confidence high -exclude G304 ./cmd/... ./data/... ./pkg/... "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.22 \
    ./hack/go-sec.sh "${@}"
fi
