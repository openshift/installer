#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  yamllint .
else
  ENGINE="podman"
  if [ "$(uname)" = "Darwin" ]; then
    ENGINE="docker"
  fi
  "$ENGINE" run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/workdir:z" \
    --entrypoint sh \
    quay.io/coreos/yamllint \
    ./hack/yaml-lint.sh
fi;
