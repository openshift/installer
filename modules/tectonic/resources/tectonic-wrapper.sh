#!/bin/bash
set -e

# shellcheck disable=SC2086,SC2154
/usr/bin/docker run \
    --volume "$(pwd)":/assets \
    --network=host \
    --entrypoint=/bin/sh \
    ${hyperkube_image} \
    /assets/tectonic.sh /assets/auth/kubeconfig /assets

# remove TNC static pod
rm -f /etc/kubernetes/manifests/tectonic-node-controller-pod.yaml
