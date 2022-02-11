#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  tflint
else
  ENGINE="podman"
  if [ "$(uname)" = "Darwin" ]; then
    ENGINE="docker"
  fi
  "$ENGINE" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/data:z" \
    --entrypoint sh \
    quay.io/coreos/tflint \
    ./hack/tf-lint.sh
fi;
