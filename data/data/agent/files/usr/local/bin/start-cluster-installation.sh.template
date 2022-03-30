#!/bin/bash
set -e

source common.sh

wait_for_assisted_service

BASE_URL="http://{{.NodeZeroIP}}:8090/api/assisted-install/v2"

# Get cluster id
cluster_id=$(curl -s -S "${BASE_URL}/clusters" | jq -r .[].id)
# Get infra_env_id
infra_env_id=$(curl -s -S "${BASE_URL}/infra-envs"| jq -r .[].id)

required_hosts={{.ControlPlaneAgents}}
echo "Number of required hosts: ${required_hosts}"

known_hosts=0

while [[ "${required_hosts}" != "${known_hosts}" ]]
do
    host_status=$(curl -s -S "${BASE_URL}/infra-envs/${infra_env_id}/hosts" | jq -r .[].status)
    if [[ -n ${host_status} ]]; then
        for status in ${host_status}; do
            if [[ "${status}" == "known" ]]; then
                ((known_hosts+=1))
                echo "Hosts known and ready for cluster installation (${known_hosts}/${required_hosts})"
            else
                echo "Waiting for hosts to become ready for cluster installation..."
                sleep 30
            fi
        done
    fi
done

echo "All ${required_hosts} hosts are ready."

api_vip=$(curl -s -S "${BASE_URL}/clusters" | jq -r .[].api_vip)
if [ "${api_vip}" == null ]; then
    echo "Setting api vip"
    curl -s -S -X PATCH "${BASE_URL}/clusters/${cluster_id}" -H "Content-Type: application/json" -d '{"api_vip": "{{.APIVIP}}"}'
fi

while [[ "${cluster_status}" != "ready" ]]
do
    cluster_status=$(curl -s -S "${BASE_URL}/clusters" | jq -r .[].status)
    echo "Cluster status: ${cluster_status}"
    sleep 5
    if [[ "${cluster_status}" == "ready" ]]; then
        echo "Starting cluster installation..."
        curl -s -S -X POST "${BASE_URL}/clusters/${cluster_id}/actions/install" \
            -H 'accept: application/json' \
            -d ''
        echo "Cluster installation started"
    fi
done
