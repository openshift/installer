#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  cd ..
  SOURCE_DIR="$(pwd)"
  cd .. || exit
  ROOT_DIR="$(pwd)"
  TARGET_DIR="${ROOT_DIR}/openshift"
  if [ "$SOURCE_DIR" != "$TARGET_DIR" ]; then
    mv "$SOURCE_DIR" "$TARGET_DIR"
  fi;
  cd "${TARGET_DIR}/installer/" || exit
  go vet "$1"
else
  for dir in "$@"
  do
    docker run --rm --env IS_CONTAINER='TRUE' -v "$PWD":/go/src/github.com/openshift/installer -w /go/src/github.com/openshift/installer quay.io/coreos/golang-testing ./hack/go-vet.sh "$dir"
  done;
fi;
