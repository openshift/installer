#!/bin/sh
# Example:  ./hack/go-lint.sh installer/... tests/smoke

if [ "$IS_CONTAINER" != "" ]; then
  golint -set_exit_status "${@}"
else
  docker run --rm --env IS_CONTAINER='TRUE' \
    -v "$PWD":/go/src/github.com/openshift/installer \
    -w /go/src/github.com/openshift/installer \
    --entrypoint sh quay.io/coreos/golang-testing \
    ./hack/go-lint.sh "${@}"
fi
