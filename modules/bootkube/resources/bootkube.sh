#!/bin/bash

# When self-hosted etcd is enabled, bootkube places an static pod manifest in
# /etc/kubernetes/manifests for Kubelet to boot a temporary etcd instance.
# However, Kubelet might not have started yet and therefore the folder might
# be missing for now, making bootkube crash.
mkdir -p /etc/kubernetes/manifests/

# Move optional experimental manifests into bootkube friendly locations
[ -d /opt/tectonic/experimental ] && mv /opt/tectonic/experimental/* /opt/tectonic/manifests/ && rm -r /opt/tectonic/experimental
[ -d /opt/tectonic/bootstrap-experimental ] && mv /opt/tectonic/bootstrap-experimental/* /opt/tectonic/bootstrap-manifests/ && rm -r /opt/tectonic/bootstrap-experimental
# Move network related manifests into bootkube friendly locations
[ -d /opt/tectonic/net-manifests ] && mv /opt/tectonic/net-manifests/* /opt/tectonic/manifests/ && rm -r /opt/tectonic/net-manifests

# shellcheck disable=SC2154
/usr/bin/rkt run \
  --trust-keys-from-https \
  --volume assets,kind=host,source="$(pwd)" \
  --mount volume=assets,target=/assets \
  --volume etc-kubernetes,kind=host,source=/etc/kubernetes \
  --mount volume=etc-kubernetes,target=/etc/kubernetes \
  "${bootkube_image}" \
  --net=host \
  --dns=host \
  --exec=/bootkube -- start --asset-dir=/assets
