#!/bin/sh

# in prow, already in container, so no 'docker run'
if [ "$IS_CONTAINER" != "" ]; then
  set -x
  /terraform fmt -list -check -write=false
else
  docker run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:${PWD}:ro,z" \
    --workdir "${PWD}" \
    quay.io/coreos/terraform-alpine:v0.11.7 \
    ./hack/tf-fmt.sh
fi
