#!/bin/bash

set -euo pipefail

# It is possible machine-api-operator comes up while the bootstrap is
# online, meaning there could be two DHCP servers on the network. To
# avoid bootstrap responding to a worker, which would cause a failed
# deployment, we filter out requests from anyone else than the control
# plane.  We are using iptables instead of dnsmasq's dhcp-host because
# DHCPv6 wants to use DUID's instead of mac addresses.

install_filter() {
    $IPTABLES -t raw -N DHCP_IRONIC
    $IPTABLES -t raw -A PREROUTING -p udp --dport 67 -j DHCP_IRONIC
    $IPTABLES -t raw -A PREROUTING -p udp --dport 547 -j DHCP_IRONIC

    for mac in ${DHCP_ALLOW_MACS}; do
        $IPTABLES -t raw -A DHCP_IRONIC -m mac --mac-source "${mac}" -j ACCEPT
    done

    $IPTABLES -t raw -A DHCP_IRONIC -j DROP
}

remove_filter() {
    "$IPTABLES-save" -t raw | grep -v DHCP_IRONIC | "$IPTABLES-restore"
}

if [ -n "${DHCP_ALLOW_MACS}" ]; then
    case "$1" in
        install)
            install_filter
            ;;
        remove)
            remove_filter
            ;;
        *)
            echo "Unknown command \"$1\"" >&2
            exit 1
            ;;
    esac
fi
