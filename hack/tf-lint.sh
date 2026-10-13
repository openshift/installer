#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  tflint
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/data:z" \
    --entrypoint sh \
    quay.io/coreos/tflint \
    ./hack/tf-lint.sh
fi;
