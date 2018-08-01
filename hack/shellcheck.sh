#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  find ./ ! -name "$(printf "*\\n*")" -name '*.sh' > tmp-search
  while IFS= read -r file
  do
    if ! shellcheck --format=gcc "$file"; then
      export FAILED=true;
    fi;
  done < tmp-search
  rm tmp-search

  if [ "$FAILED" != "" ]; then
    exit 1;
  fi;
  echo shellcheck passed;
else
  docker run -e IS_CONTAINER='TRUE' --rm -v "$(pwd)":/workdir:ro --entrypoint sh quay.io/coreos/shellcheck-alpine:v0.5.0 '/workdir/hack/shellcheck.sh';
fi;
