#!/bin/bash
set -euo pipefail

/usr/local/bin/release-image-download.sh

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

IMAGE=$(image_for agent-installer-node-agent)

echo "Using agent image: ${IMAGE} to copy bin"

/usr/bin/podman run --privileged --rm -v /usr/local/bin:/hostbin "${IMAGE}" cp /usr/bin/agent /hostbin
