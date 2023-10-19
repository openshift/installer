#!/bin/bash

# verify_module verifies the vendor of a module
# $1: directory of the module
# $2: golang compatability requirement (optional)
verify_module() {
  pushd "$1"
  tidy_args=()
  [ -n "${2:-}" ] && tidy_args=("-compat=$2")
  go mod tidy "${tidy_args[@]}"
  go mod vendor
  go mod verify
  popd
}

if [ "$IS_CONTAINER" != "" ]; then
  set -eux

  # Verify the main installer module.
  verify_module "${PWD}"

  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.20 \
    ./hack/verify-vendor.sh "${@}"
fi
