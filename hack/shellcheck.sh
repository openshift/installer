#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  find "${1:-.}" -type f -name '*.sh' -exec shellcheck --format=gcc {} \+
else
  docker run -e IS_CONTAINER='TRUE' --rm -v "$(pwd)":/workdir:ro --entrypoint sh quay.io/coreos/shellcheck-alpine:v0.5.0 /workdir/hack/shellcheck.sh /workdir;
fi;
