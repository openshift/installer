#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  go vet "${@}"
else
  docker run --rm --env IS_CONTAINER='TRUE' -v "$PWD":/go/src/github.com/openshift/installer -w /go/src/github.com/openshift/installer quay.io/coreos/golang-testing ./hack/go-vet.sh "${@}"
fi;
