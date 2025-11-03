#!/bin/bash

set -eu

TARGET_SIZE="$1"
BASE_DIR="/run/ephemeral_base"
LOOPBACK_FILE="${BASE_DIR}/loopfs"

function get_size_bytes() {
    printf "%d" "$(($(stat -f -c "%b * %s" "${BASE_DIR}")))"
}

TARGET_SIZE_BYTES="$(numfmt --from=iec "${TARGET_SIZE}")"
CURRENT_SIZE_BYTES="$(get_size_bytes)"
CURRENT_SIZE="$(numfmt --to=iec "${CURRENT_SIZE_BYTES}")"

if ((TARGET_SIZE_BYTES > CURRENT_SIZE_BYTES)); then
    echo "Expanding ephemeral base dir from ${CURRENT_SIZE} to ${TARGET_SIZE}"
    mount -o remount,size="${TARGET_SIZE}" "${BASE_DIR}"

    echo "Expanding ephemeral loopback"
    truncate -s "$(get_size_bytes)" "${LOOPBACK_FILE}"

    LOOPBACK_DEVICE="$(losetup -j "${LOOPBACK_FILE}" -O NAME -n)"
    losetup -c "${LOOPBACK_DEVICE}"

    echo "Expanding ephemeral filesystem"
    xfs_growfs -d "${LOOPBACK_DEVICE}"
else
    echo "Ephemeral base dir size ${CURRENT_SIZE} is already larger than ${TARGET_SIZE}; not expanding"
fi

df -h
