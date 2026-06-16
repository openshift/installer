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
# Exclude crypto-policies as it may be bind-mounted by fips-crypto-policy-overlay
# during FIPS boot, causing "Device or resource busy" and "Read-only file system"
# errors when rsync attempts to rename/symlink files under it.
rsync -a --exclude crypto-policies "${ostree_checkout}/usr/etc/" /etc

# reload the new policy
echo "Reloading SELinux policy"
semodule -R

# handle upgrade of sshd between RHEL 9.6 and 9.8
systemctl --no-block try-restart sshd.service
