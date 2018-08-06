#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  tflint
else
  docker run -t --rm -v "$(pwd)":/data --env IS_CONTAINER='TRUE' --entrypoint sh quay.io/coreos/tflint ./hack/tf-lint.sh
fi;
