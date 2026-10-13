#!/bin/bash

set -e

status_name=60-rendezvous-host
set_rendezvous_message() {
    mkdir -p /etc/motd.d/
    tee "/etc/issue.d/${status_name}.issue" | sed -e 's/\\e[{][^}]*[}]//g' | tee "/etc/motd.d/${status_name}" 1>&2
    agetty --reload
}

rendezvous_host_env="/etc/assisted/rendezvous-host.env"
while [ ! -f "${rendezvous_host_env}" ]; do
    printf '\\e{lightred}Not configured - no Rendezvous IP set\\e{reset}\n' | set_rendezvous_message
    sleep 30
done
rm -f "/etc/issue.d/${status_name}.issue" "/etc/motd.d/${status_name}"
agetty --reload

# shellcheck disable=SC1090
source "${rendezvous_host_env}"
echo "NODE_ZERO_IP: $NODE_ZERO_IP"

is_node_zero() {
    local is_rendezvous_host
    is_rendezvous_host=$(ip -j address | jq "[.[].addr_info] | flatten | map(.local==\"$NODE_ZERO_IP\") | any")
    if [[ "${is_rendezvous_host}" == "true" ]]; then
        echo 1
    else
        echo 0
    fi
}

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

    # Create tls certs, if they don't exist, via the installer command.
    # This allows the certs to be created at run-time, e.g. when installed via the UI
    AGENT_TLS_DIR=/opt/agent/tls
    if [ -z "$(ls -A "$AGENT_TLS_DIR")" ]; then
       # shellcheck disable=SC1091
       . /usr/local/bin/release-image.sh
       IMAGE=$(image_for installer)
       /usr/bin/podman run --privileged -v /tmp:/assets --rm "${IMAGE}" agent create certificates --dir=/assets
       cp /tmp/tls/* $AGENT_TLS_DIR
    fi

    NODE0_PATH=/etc/assisted/node0
    mkdir -p "$(dirname "${NODE0_PATH}")"

    NODE_ZERO_MAC=$(ip -j address | jq -r ".[] | select(.addr_info | map(select(.local == \"$NODE_ZERO_IP\")) | any).address")
    echo "MAC Address for Node 0: ${NODE_ZERO_MAC}"

    cat >"${NODE0_PATH}" <<EOF
# This file exists if the agent-based installer has determined the host is node 0.
# The host is determined to be node 0 when one of its network interfaces has an 
# IP address matching NODE_ZERO_IP in /etc/assisted/rendezvous-host.env. 
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
# agent-register-cluster.service 
# agent-register-infraenv.service
# install-status.service
# start-cluster-installation.service
BOOTSTRAP_HOST_MAC=${NODE_ZERO_MAC}
EOF

    echo "Created file ${NODE0_PATH}"

    printf 'This host (%s) is the rendezvous host.\n' "${NODE_ZERO_IP}" | set_rendezvous_message

    cat <<EOF >/etc/motd
The primary service is assisted-service.service. To watch its status, run:

  journalctl -u assisted-service.service
EOF
else

    printf 'This host is not the rendezvous host. The rendezvous host is at %s.\n' "${NODE_ZERO_IP}" | set_rendezvous_message
fi
