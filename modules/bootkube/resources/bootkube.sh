#!/bin/bash
set -e

# shellcheck disable=SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets \
    --volume /etc/kubernetes:/etc/kubernetes \
    "${kube_core_renderer_image}" \
    --config=/assets/kco-config.yaml \
    --output=/assets

# shellcheck disable=SC2154
/usr/bin/docker run \
    --user 0 \
    --volume "$(pwd)":/assets \
    "${tnc_operator_image}" \
    --config=/assets/tnco-config.yaml \
    --render-bootstrap=true \
    --render-output=/assets/tnc-bootstrap

mkdir -p /etc/kubernetes/manifests/
cp $(pwd)/tnc-bootstrap/tectonic-node-controller-config.yaml /etc/kubernetes/tnc-config
cp $(pwd)/tnc-bootstrap/tectonic-node-controller-pod.yaml $(pwd)/bootstrap-manifests/

# shellcheck disable=SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets \
    --volume /etc/kubernetes:/etc/kubernetes \
    --network=host \
    --entrypoint=/bootkube \
    "${bootkube_image}" \
    start --asset-dir=/assets
