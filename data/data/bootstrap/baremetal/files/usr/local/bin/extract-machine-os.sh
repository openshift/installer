#!/bin/bash

set -eu

CID_FILE="$1"

podman run --rm --name machine-os-extractor --cidfile="${CID_FILE}" --cgroups=no-conmon --sdnotify=ignore --log-driver=passthrough --env IP_OPTIONS="${PROVISIONING_IP_OPTIONS}" -v systemd-ironic:/shared:z "${MACHINE_OS_IMAGES_IMAGE}" /bin/copy-metal --all /shared/html/images/

systemd-notify --ready
