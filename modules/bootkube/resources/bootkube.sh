#!/bin/bash

# When self-hosted etcd is enabled, bootkube places an static pod manifest in
# /etc/kubernetes/manifests for Kubelet to boot a temporary etcd instance.
# However, Kubelet might not have started yet and therefore the folder might
# be missing for now, making bootkube crash.
mkdir -p /etc/kubernetes/manifests/

/usr/bin/rkt run \
  --trust-keys-from-https \
  --volume assets,kind=host,source=$(pwd) \
  --mount volume=assets,target=/assets \
  --volume etc-kubernetes,kind=host,source=/etc/kubernetes \
  --mount volume=etc-kubernetes,target=/etc/kubernetes \
  ${bootkube_image} \
  --net=host \
  --dns=host \
  --exec=/bootkube -- start --asset-dir=/assets
