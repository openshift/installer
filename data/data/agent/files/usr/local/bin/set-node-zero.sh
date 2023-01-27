#!/bin/bash

set -e

# shellcheck disable=SC1091
source common.sh
echo "NODE_ZERO_IP: $NODE_ZERO_IP"

timeout=$((SECONDS + 30))

while [[ $SECONDS -lt $timeout ]]
do
    if [[ $(is_node_zero) -eq 1 ]]; then
        IS_NODE_ZERO="true"
        break
    fi
    sleep 5
done

if [ "${IS_NODE_ZERO}" = "true" ]; then
    echo "Node 0 IP ${NODE_ZERO_IP} found on this host" 1>&2

    NODE0_PATH=/etc/assisted/node0
    mkdir -p "$(dirname "${NODE0_PATH}")"

    NODE_ZERO_MAC=$(ip -j address | jq -r ".[] | select(.addr_info | map(select(.local == \"$NODE_ZERO_IP\")) | any).address")
    echo "MAC Address for Node 0: ${NODE_ZERO_MAC}"

    cat >"${NODE0_PATH}" <<EOF
# This file exists if the agent-based installer has determined the host is node 0.
# The host is determined to be node 0 when one of its network interfaces has an 
# IP address matching NODE_ZERO_IP in /etc/assisted/agent-installer.env. 
# The MAC address of the network interface matching NODE_ZERO_IP is noted below 
# as BOOTSTRAP_HOST_MAC.
#
# BOOTSTRAP_HOST_MAC is read by assisted-service. The host with a MAC address
# matching this value in assisted-service is designated to be the bootstrap during
# cluster installation. In assisted-service.service, this file is included as a
# --env-file in the ExecStart command.
#
# This file is also a ConditionPathExists in the following systemd service
# definitions:
# apply-host-config.service
# assisted-service-pod.service
# create-cluster-and-infraenv.service
# install-status.service
# start-cluster-installation.service
BOOTSTRAP_HOST_MAC=${NODE_ZERO_MAC}
EOF

    echo "Created file ${NODE0_PATH}"

    rendezvousHostMessage="This host ${NODE_ZERO_IP} is the rendezvous host."

    cat <<EOF >/etc/motd
The primary service is assisted-service.service. To watch its status, run:

  journalctl -u assisted-service.service
EOF
else

    rendezvousHostMessage="This host is not the rendezvous host. The rendezvous host is at ${NODE_ZERO_IP}."
fi
mkdir -p /etc/motd.d/
echo "$rendezvousHostMessage" > /etc/motd.d/60-rendezvous-host
echo "$rendezvousHostMessage" > /etc/issue.d/60-rendezvous-host.issue
agetty --reload
