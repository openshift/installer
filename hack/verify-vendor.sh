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

  # Verify the sub-modules for the terraform providers.
  find terraform/providers -maxdepth 1 -mindepth 1 -print0 | while read -r -d '' dir
  do
    verify_module "$dir"
  done
  # Verify the terraform sub-module.
  verify_module "terraform/terraform"

  find cluster-api/providers -maxdepth 1 -mindepth 1 -print0 | while read -r -d '' dir
  do
    verify_module "$dir"
  done
  verify_module "cluster-api/cluster-api"

  git diff --exit-code
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/golang:1.22 \
    ./hack/verify-vendor.sh "${@}"
fi
