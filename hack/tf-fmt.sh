#!/bin/sh

# in prow, already in container, so no 'podman run'
if [ "$IS_CONTAINER" != "" ]; then
  if [ "${#N}" -gt 1 ]; then
    set -- -list -check -write=false
  fi
  set -x
  /terraform fmt "${@}"  # FIXME: drop this slash after we update openshift/release
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:${PWD}:z" \
    --workdir "${PWD}" \
    quay.io/coreos/terraform-alpine:v0.11.8 \
    ./hack/tf-fmt.sh
fi
