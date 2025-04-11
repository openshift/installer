#!/bin/bash
set -e

# shellcheck disable=SC1091
source "common.sh"
# shellcheck disable=SC1091
source "issue_status.sh"

cluster_id=""
while [[ "${cluster_id}" = "" ]]
do
    # Get cluster id
    cluster_id=$(curl_assisted_service "/clusters" GET | jq -r .[].id)
    if [[ "${cluster_id}" = "" ]]; then
        sleep 2
    fi
done

printf '\nInfra env id is %s\n' "${INFRA_ENV_ID}" 1>&2

total_required_nodes=$(( REQUIRED_MASTER_NODES + REQUIRED_WORKER_NODES ))
echo "Number of required master nodes: ${REQUIRED_MASTER_NODES}" 1>&2
echo "Number of required worker nodes: ${REQUIRED_WORKER_NODES}" 1>&2
echo "Total number of required nodes: ${total_required_nodes}" 1>&2

status_issue="90_start-install"

num_known_hosts() {
    local known_hosts=0
    local insufficient_hosts=0
    host_status=$(curl_assisted_service "/infra-envs/${INFRA_ENV_ID}/hosts" GET | jq -r .[].status)
    if [[ -n ${host_status} ]]; then
        for status in ${host_status}; do
            if [[ "${status}" == "known" ]]; then
                ((known_hosts+=1))
            fi
            if [ "${status}" == "insufficient" ]; then
                ((insufficient_hosts+=1))
            fi
        done
        echo "Hosts known and ready for cluster installation (${known_hosts}/${total_required_nodes})" 1>&2
    fi
    if (( known_hosts != total_required_nodes )); then
        printf '\\e{yellow}Waiting for all hosts to be ready:\\e{reset}\n%d hosts expected\n%d hosts ready, %d hosts not validated' "${total_required_nodes}" "${known_hosts}" "${insufficient_hosts}" | set_issue "${status_issue}"
    fi
    echo "${known_hosts}"
}

while [[ "${total_required_nodes}" != $(num_known_hosts) ]]
do
    echo "Waiting for hosts to become ready for cluster installation..." 1>&2
    sleep 10
done

echo "All ${total_required_nodes} hosts are ready." 1>&2
clear_issue "${status_issue}"

while [[ "${cluster_status}" != "installed" ]]
do
    sleep 5
    cluster_info="$(curl_assisted_service "/clusters" GET)"
    cluster_status=$(printf '%s' "${cluster_info}" | jq -r .[].status)
    echo "Cluster status: ${cluster_status}" 1>&2
    # Start the cluster install, if it transitions back to Ready due to a failure,
    # then it will be restarted
    case "${cluster_status}" in
        "ready")
            echo "Starting cluster installation..." 1>&2
            res=$(curl_assisted_service "/clusters/${cluster_id}/actions/install" POST -w "%{http_code}" -o /dev/null)
            if [[ $res = "202" ]]; then 
                printf '\nCluster installation started\n' 1>&2
            fi
            ;&
        "installed" | "preparing-for-installation" | "installing")
            printf '\\e{lightgreen}Cluster installation in progress\\e{reset}' | set_issue "${status_issue}"
            ;;
        *)
            printf '\\e{lightred}Installation cannot proceed:\\e{reset} Cluster status: %s' "${cluster_status}" | set_issue "${status_issue}"
            printf '%s\n' "${cluster_info}"
            ;;
    esac
done
