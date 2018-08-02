#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  yamllint --config-data "{extends: default, rules: {line-length: {level: warning, max: 120}}}" ./examples/ ./installer/
else
  docker run -t --rm -v "$(pwd)":/workdir --env IS_CONTAINER='TRUE' --entrypoint sh quay.io/coreos/yamllint ./hack/yaml-lint.sh
fi;
