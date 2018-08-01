#!/bin/sh

if [ "$IS_CONTAINER" != "" ]; then
  cd ..
  SOURCE_DIR="$(pwd)"
  cd .. || exit
  ROOT_DIR="$(pwd)"
  # if running in a local prow instance, with repo ../not-openshift/installer, have to modify here
  TARGET_DIR="${ROOT_DIR}/openshift"
  if [ "$SOURCE_DIR" != "$TARGET_DIR" ]; then
    mv "$SOURCE_DIR" "$TARGET_DIR"
  fi;
  cd "${TARGET_DIR}/installer/" || exit
  CGO_ENABLED=0 go test ./installer/...
else
  docker run -e IS_CONTAINER='TRUE' -v "$PWD":/go/src/github.com/openshift/installer -w /go/src/github.com/openshift/installer quay.io/coreos/golang-testing ./hack/gotest.sh
fi

