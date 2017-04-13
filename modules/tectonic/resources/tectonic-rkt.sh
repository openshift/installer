#!/bin/bash

/usr/bin/rkt run \
  --trust-keys-from-https \
  --volume assets,kind=host,source=$(pwd) \
  --mount volume=assets,target=/assets \
  ${hyperkube_image} \
  --net=host \
  --dns=host \
  --working-dir=/assets/tectonic \
  --exec=/bin/bash -- /assets/tectonic.sh /assets/kubeconfig /assets ${experimental}
