#!/bin/bash

# shellcheck disable=SC2086,SC2154
/usr/bin/rkt run \
  --trust-keys-from-https \
  --volume assets,kind=host,source="$(pwd)" \
  --mount volume=assets,target=/assets \
  ${hyperkube_image} \
  --net=host \
  --dns=host \
  --exec=/bin/bash -- /assets/tectonic.sh /assets/auth/kubeconfig /assets ${experimental}
