#!/bin/bash

cat <<EOF >/etc/motd
The primary service is agent.service. To watch its status, run:

  journalctl -u agent.service

To view the agent log, run:

  journalctl TAG=agent
EOF
echo "Waiting for network to determine if this is the rendezvous host." > /etc/motd.d/60-rendezvous-host

#
# The hostnames defined in agent-config.yaml are written out
# to files at /etc/assisted/hostnames/<MAC-address>.
#
# If a host has multiple interfaces, the host's first network
# interface's MAC address is used.
#
# This script compares the MAC addresses on the current host
# with the addresses in /etc/assisted/hostnames/.
#
# If a match is found, then the hostname in the file is set
# as this host's hostname.
#

HOSTNAMES_PATH=/etc/assisted/hostnames
FILES=$(ls $HOSTNAMES_PATH)
for filename in ${FILES}
do
    MATCHED_MAC_ADDRESS_WITH_HOST=$(ip address | grep "${filename}")
    if [ "$MATCHED_MAC_ADDRESS_WITH_HOST" != "" ]; then
        HOSTNAME="$(cat "${HOSTNAMES_PATH}/${filename}")"
        echo "Host has matching MAC address: ${filename}" 1>&2
        echo "Setting hostname to ${HOSTNAME}" 1>&2
        hostnamectl set-hostname "${HOSTNAME}"
    else
        echo "MAC address, ${filename}, does not exist on this host" 1>&2
    fi
done
