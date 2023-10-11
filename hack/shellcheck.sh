#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  TOP_DIR="${1:-.}"
  find "${TOP_DIR}" \
    -path "${TOP_DIR}/vendor" -prune \
    -o -path "${TOP_DIR}/cluster-api/*/vendor" -prune \
    -o -path "${TOP_DIR}/terraform/*/vendor" -prune \
    -o -type f -name '*.sh' -exec shellcheck --format=gcc {} \+
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/workdir:ro,z" \
    --entrypoint sh \
    --workdir /workdir \
    quay.io/coreos/shellcheck-alpine:v0.5.0 \
    /workdir/hack/shellcheck.sh "${@}"
fi;
