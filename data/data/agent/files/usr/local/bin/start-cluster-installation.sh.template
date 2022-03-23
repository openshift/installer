#!/bin/bash
set -e

source common.sh

wait_for_assisted_service

BASE_URL="http://{{.NodeZeroIP}}:8090/api/assisted-install/v2"

# Get cluster id
cluster_id=$(curl -s -S ${BASE_URL}/clusters | jq -r .[].id)
# Get infra_env_id
infra_env_id=$(curl -s -S ${BASE_URL}/infra-envs| jq -r .[].id)

known_hosts=0

while [[ "$cluster_status" != "ready" && "${known_hosts}" != "${host_count}" ]]
do
    cluster_status=$(curl -s -S ${BASE_URL}/clusters | jq -r .[].status)
    host_status=$(curl -s -S ${BASE_URL}/infra-envs/${infra_env_id}/hosts | jq -r .[].status)
    if [[ -n ${host_status} ]]; then
	    host_count=$(curl -s -S ${BASE_URL}/infra-envs/${infra_env_id}/hosts | jq length)
        echo Total hosts: ${host_count}
        for status in ${host_status}; do
            echo "Host status is $status"
            if [[ "${status}" == "known" ]]; then
                ((known_hosts+=1))
                echo "Hosts known and ready for cluster installation (${known_hosts}/${host_count})"
                if [[ "${known_hosts}" == "${host_count}" ]]; then
                    break
                fi
            else
                echo "Waiting for hosts to become ready for cluster installation..."
                sleep 10
            fi
        done
    fi
done
echo "All ${host_count} hosts are ready."

while [[ "$cluster_status" != "ready" ]]
do
    echo "Setting api vip"
    curl -s -S -X PATCH ${BASE_URL}/clusters/${cluster_id} -H "Content-Type: application/json" -d '{"api_vip": "{{.APIVIP}}"}'
    cluster_status=$(curl -s -S ${BASE_URL}/clusters | jq -r .[].status)
    echo "Cluster status: ${cluster_status}"
    if [[ "$cluster_status" == "ready" ]]; then
        echo "Starting cluster installation..."
        curl -s -S -X POST ${BASE_URL}/clusters/${cluster_id}/actions/install \
            -H 'accept: application/json' \
            -d ''
        echo "Cluster installation started"
    fi
done
