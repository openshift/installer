#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  set -x
  bazel --output_base=/tmp build tarball
else
  docker run -e IS_CONTAINER='TRUE' --rm -v "$PWD":"$PWD" -v /tmp:/tmp:rw -w "$PWD" quay.io/coreos/tectonic-builder:bazel-v0.3 ./hack/test-bazel-build-tarball.sh
fi
