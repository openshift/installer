#!/bin/bash

set -e

status_name=60-check-config-image
set_mount_config_image_message() {
    mkdir -p /etc/motd.d/
    tee "/etc/issue.d/${status_name}.issue" | sed -e 's/\\e[{][^}]*[}]//g' | tee "/etc/motd.d/${status_name}" 1>&2
    agetty --reload
}

rendezvous_host_env="/etc/assisted/rendezvous-host.env"
while [ ! -f "${rendezvous_host_env}" ]; do
    printf '\\e{lightred}Insert or mount config image to start cluster installation\\e{reset}\n' | set_mount_config_image_message
    sleep 30
done
rm -f "/etc/issue.d/${status_name}.issue" "/etc/motd.d/${status_name}"
agetty --reload
