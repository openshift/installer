#!/bin/bash

# shellcheck disable=SC1091
source /etc/assisted/rendezvous-host.env

is_node_zero() {
    local is_rendezvous_host
    is_rendezvous_host=$(ip -j address | jq "[.[].addr_info] | flatten | map(.local==\"$NODE_ZERO_IP\") | any")
    if [[ "${is_rendezvous_host}" == "true" ]]; then
        echo 1
    else
        echo 0
    fi
}
