#!/bin/bash
set -e

# When self-hosted etcd is enabled, bootkube places an static pod manifest in
# /etc/kubernetes/manifests for Kubelet to boot a temporary etcd instance.
# However, Kubelet might not have started yet and therefore the folder might
# be missing for now, making bootkube crash.
mkdir -p /etc/kubernetes/manifests/

# Move optional self hosted etcd manifests into bootkube friendly locations
if [ -d /opt/tectonic/etcd ]; then
    mv /opt/tectonic/etcd/manifests/* /opt/tectonic/manifests/
    rm -r /opt/tectonic/etcd/manifests
    mv /opt/tectonic/etcd/bootstrap-manifests/* /opt/tectonic/bootstrap-manifests/
    rm -r /opt/tectonic/etcd/bootstrap-manifests
fi

# Move network related manifests into bootkube friendly locations
if [ -d /opt/tectonic/net-manifests ]; then
    mv /opt/tectonic/net-manifests/* /opt/tectonic/manifests/
    rm -r /opt/tectonic/net-manifests
fi

# shellcheck disable=SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets \
    --volume /etc/kubernetes:/etc/kubernetes \
    --network=host \
    --entrypoint=/bootkube \
    "${bootkube_image}" \
    start --asset-dir=/assets
