#!/bin/bash
set -e

# shellcheck disable=SC1091
source issue_status.sh

BASE_URL="${SERVICE_BASE_URL}api/assisted-install/v2"

cluster_id=""
while [[ "${cluster_id}" = "" ]]
do
    # Get cluster id
    cluster_id=$(curl -s -S "${BASE_URL}/clusters" -H "Authorization: ${AGENT_AUTH_TOKEN}" | jq -r .[].id)
    if [[ "${cluster_id}" = "" ]]; then
        sleep 2
    fi
done

printf '\nInfra env id is %s\n' "${INFRA_ENV_ID}" 1>&2

status_issue="90_add-node"

# Wait for the current host to be ready
host_ready=false
while [[ $host_ready == false ]]
do
    host_status=$(curl -s -S "${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts" -H "Authorization: ${AGENT_AUTH_TOKEN}" | jq -r ".[].status")
    if [[ "${host_status}" != "known" ]]; then
        printf '\\e{yellow}Waiting for the host to be ready' | set_issue "${status_issue}"
        sleep 10
    else
        host_ready=true
    fi
done

HOST_ID=$(curl -s "${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts" -H "Authorization: ${AGENT_AUTH_TOKEN}" | jq -r '.[].id')
printf '\nHost %s is ready for installation\n' "${HOST_ID}" 1>&2
clear_issue "${status_issue}"

# Add the current host to the cluster
res=$(curl -X POST -s -S -w "%{http_code}\\n" -o /dev/null "${BASE_URL}/infra-envs/${INFRA_ENV_ID}/hosts/${HOST_ID}/actions/install" -H "Authorization:  ${AGENT_AUTH_TOKEN}")
if [[ $res = "202" ]]; then 
    printf '\nHost installation started\n' 1>&2
else
    printf '\nHost installation failed\n' 1>&2
    exit 1
fi
