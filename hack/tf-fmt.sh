#!/bin/sh

# in prow, already in container, so no 'docker run'
if [ "$IS_CONTAINER" != "" ]; then
  set -x
  /terraform fmt -list -check -write=false
else
  docker run -e IS_CONTAINER='TRUE' --rm -v "$PWD":"$PWD":ro -v /tmp:/tmp:rw -w "$PWD" quay.io/coreos/terraform-alpine:v0.11.7 ./hack/tf-fmt.sh
fi
