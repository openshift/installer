#!/bin/bash

set -euo pipefail

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

export KUBECONFIG=/opt/openshift/auth/kubeconfig-loopback

build_metal3_env() {
    printf 'METAL3_BAREMETAL_OPERATOR_IMAGE="%s"\n' "$(image_for baremetal-operator)"
}

build_metal3_env | tee -a /etc/metal3.env
