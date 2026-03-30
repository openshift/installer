#!/bin/bash
set -euo pipefail

ostree_checkout=/ostree/repo/tmp/node-image
if [ ! -d "${ostree_checkout}" ]; then
    ostree_checkout=/var/ostreecontainer/checkout
fi

echo "Overlaying node image content"

# keep /usr/lib/modules from the booted deployment for kernel modules
mount -o bind,ro "/usr/lib/modules" "${ostree_checkout}/usr/lib/modules"
mount -o rbind,ro "${ostree_checkout}/usr" /usr
rsync -a "${ostree_checkout}/usr/etc/" /etc

# reload the new policy
echo "Reloading SELinux policy"
semodule -R
