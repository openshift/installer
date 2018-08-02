#!/bin/sh
DIR="${1:-.}"
if [ "$IS_CONTAINER" != "" ]; then
  if  find "${DIR}" -type f -name '*.sh' -exec shellcheck --format=gcc {} \+
  then
    echo "Shell Check Passed"
  fi
else
  docker run -e IS_CONTAINER='TRUE' --rm -v "$(realpath "${DIR}")":/workdir:ro --entrypoint sh quay.io/coreos/shellcheck-alpine:v0.5.0 /workdir/hack/shellcheck.sh /workdir;
fi;
