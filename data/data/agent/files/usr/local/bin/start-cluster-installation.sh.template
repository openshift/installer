#!/bin/bash
set -e

source common.sh

wait_for_assisted_service

BASE_URL="{{.ServiceBaseURL}}api/assisted-install/v2"

# Get cluster id
cluster_id=$(curl -s -S "${BASE_URL}/clusters" | jq -r .[].id)
# Get infra_env_id
infra_env_id=$(curl -s -S "${BASE_URL}/infra-envs"| jq -r .[].id)

required_master_nodes={{.ControlPlaneAgents}}
required_worker_nodes={{.WorkerAgents}}
total_required_nodes=$(( ${required_master_nodes}+${required_worker_nodes} ))
echo "Number of required master nodes: ${required_master_nodes}"
echo "Number of required worker nodes: ${required_worker_nodes}"
echo "Total number of required nodes: ${total_required_nodes}"

known_hosts=0

while [[ "${total_required_nodes}" != "${known_hosts}" ]]
do
    host_status=$(curl -s -S "${BASE_URL}/infra-envs/${infra_env_id}/hosts" | jq -r .[].status)
    if [[ -n ${host_status} ]]; then
        for status in ${host_status}; do
            if [[ "${status}" == "known" ]]; then
                ((known_hosts+=1))
                echo "Hosts known and ready for cluster installation (${known_hosts}/${total_required_nodes})"
            else
                echo "Waiting for hosts to become ready for cluster installation..."
                sleep 30
            fi
        done
    fi
done

echo "All ${total_required_nodes} hosts are ready."

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
