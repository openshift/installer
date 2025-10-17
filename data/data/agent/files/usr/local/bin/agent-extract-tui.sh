#!/bin/bash

set -xeuo pipefail

/usr/local/bin/get-container-images.sh
# shellcheck disable=SC1091
source /usr/local/share/assisted-service/agent-images.env

echo "Extracting agent-tui and libnmstate from agent-installer-utils image $AGENT_INSTALLER_UTILS_IMAGE"

container_id=$(podman create "$AGENT_INSTALLER_UTILS_IMAGE")
mnt=$(podman mount "$container_id")

cp "${mnt}/usr/bin/agent-tui" /usr/local/bin
find "${mnt}/usr/lib64" -name "libnmstate*" -exec cp -a -t /usr/local/lib {} +

podman unmount "$container_id"
podman rm "$container_id"

restorecon -FRv /usr/local/bin
