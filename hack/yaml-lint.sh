#!/bin/sh
exit 0  # temporarily disable while we work out whether to drop this
if [ "$IS_CONTAINER" != "" ]; then
  yamllint --config-data "{extends: default, rules: {line-length: {level: warning, max: 120}}}" .
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/workdir:z" \
    --entrypoint sh \
    quay.io/coreos/yamllint \
    ./hack/yaml-lint.sh
fi;
