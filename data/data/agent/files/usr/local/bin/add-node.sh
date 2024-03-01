#!/bin/bash
set -e

# shellcheck disable=SC1091
source issue_status.sh

BASE_URL="${SERVICE_BASE_URL}api/assisted-install/v2"

cluster_id=""
while [[ "${cluster_id}" = "" ]]
do
    # Get cluster id
    cluster_id=$(curl -s -S "${BASE_URL}/clusters" | jq -r .[].id)
    if [[ "${cluster_id}" = "" ]]; then
        sleep 2
    fi
done

printf '\nInfra env id is %s\n' "${INFRA_ENV_ID}" 1>&2

status_issue="90_add-node"

# For some reason the initial role patching doesn't seem to work properly
echo "Patching host..."
export HOST_ID=$(curl -s ${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts | jq -r '.[].id')
curl -X PATCH -d '{"host_role":"worker"}' -H "Content-Type: application/json" ${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts/${HOST_ID}

# Wait for the current host to be ready
host_ready=false
while [[ host_ready == false ]]
do
    host_status=$(curl -s -S "${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts/${HOST_ID}" | jq -r .[].status)
    if [[ "${host_status}" != "known" ]]; then
        printf '\\e{yellow}Waiting for the host to be ready' | set_issue "${status_issue}"
        sleep 10
    else
        host_ready=true
    fi
done

clear_issue "${status_issue}"

sleep 1m

# Add the current host to the cluster
curl -X POST -s -S ${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts/${HOST_ID}/actions/install
echo "Host installation started" 1>&2
