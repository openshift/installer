#!/bin/bash

# shellcheck disable=SC1091
source /usr/local/share/assisted-service/assisted-service.env 
source /etc/assisted/agent-installer.env

wait_for_assisted_service() {
    echo "Waiting for assisted-service to be ready"
    until curl --output /dev/null --silent --fail "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs"; do
        printf '.'
        sleep 5
    done
}

is_node_zero() {
    local is_rendezvous_host
    is_rendezvous_host=$(ip -j address | jq "[.[].addr_info] | flatten | map(.local==\"$NODE_ZERO_IP\") | any")
    if [[ "${is_rendezvous_host}" == "true" ]]; then
        echo 1
    else
        echo 0
    fi
}