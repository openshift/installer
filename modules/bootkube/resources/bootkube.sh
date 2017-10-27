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
/usr/bin/docker run \
    --volume "$(pwd)":/assets \
    --volume /etc/kubernetes:/etc/kubernetes \
    --network=host \
    --entrypoint=/bootkube \
    "${bootkube_image}" \
    start --asset-dir=/assets
