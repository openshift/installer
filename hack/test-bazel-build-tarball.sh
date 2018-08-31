#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  set -x
  bazel --output_base=/tmp build "$@" tarball
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:${PWD}:z" \
    --workdir "${PWD}" \
    quay.io/coreos/tectonic-builder:bazel-v0.3 \
    ./hack/test-bazel-build-tarball.sh
fi
