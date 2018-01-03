#!/bin/bash
set -e

mkdir -p /etc/kubernetes/manifests/
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
