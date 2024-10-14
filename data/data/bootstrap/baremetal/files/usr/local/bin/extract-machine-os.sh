#!/bin/bash

set -eu

CID_FILE="$1"

while ! podman run --rm --name machine-os-extractor --cidfile="${CID_FILE}" --cgroups=no-conmon --log-driver=passthrough --env IP_OPTIONS="${PROVISIONING_IP_OPTIONS}" -v systemd-ironic:/shared:z "${MACHINE_OS_IMAGES_IMAGE}" /bin/copy-metal --all /shared/html/images/; do
    sleep 5
done
