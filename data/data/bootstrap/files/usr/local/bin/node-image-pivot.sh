#!/bin/bash
set -euo pipefail

ostree_repo=/var/ostree-container/repo
if [ ! -d "${ostree_repo}" ]; then
    ostree_repo=/ostree/repo
fi

checkout="${ostree_repo}/tmp/node-image"

mount -o bind,ro "${checkout}/usr" /usr
rsync -a "${checkout}/usr/etc/" /etc

# reload the new policy
semodule -R
