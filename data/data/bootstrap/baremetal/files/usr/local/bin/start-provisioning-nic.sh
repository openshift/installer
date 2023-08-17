#!/bin/bash

set -euo pipefail

# Start the provisioning nic if not already started

if ! nmcli -t device | grep -q "${PROVISIONING_INTERFACE}:ethernet:connected"; then
    nmcli c add type ethernet ifname "${PROVISIONING_INTERFACE}" con-name provisioning "${PROV_IP_FAMILY}" "${PROVISIONING_SUBNET}"
    nmcli c up provisioning
else
    connection=$(nmcli -t device show "${PROVISIONING_INTERFACE}" | grep '^GENERAL\.CONNECTION' | cut -d: -f2)
    nmcli con modify "${connection}" ifname "${PROVISIONING_INTERFACE}" "${PROV_IP_FAMILY}" "${PROVISIONING_SUBNET}"
    nmcli con reload "${connection}"
    nmcli con up "${connection}"
fi

# Wait for the interface to come up
# This is how the ironic container currently detects IRONIC_IP, this could
# probably be improved by using nmcli show provisioning there instead, but we
# need to confirm that works with the static-ip-manager
while ! ip -o addr show dev "${PROVISIONING_INTERFACE}" | grep -qv link; do
    echo "Waiting for interface ${PROVISIONING_INTERFACE} to come up"
    sleep 1
done
