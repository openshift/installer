#!/bin/bash

# shellcheck disable=SC1091
source "issue_status.sh"

inactive_services() {
    local services="assisted-service.service create-cluster-and-infraenv.service apply-host-config.service start-cluster-installation.service"
    for s in ${services}; do
        if ! systemctl is-active "${s}" >/dev/null; then
            printf "%s " "${s}"
        fi
    done
}

check_services() {
    local services_issue="70_agent-services"
    local not_started
    not_started="$(inactive_services)"
    if [ -z "${not_started}" ]; then
        clear_issue "${services_issue}"
    else
        read -ra show_services <<<"${not_started}"
        {
            printf '\\e{cyan}Waiting for services:\\e{reset}'
            systemctl --no-legend list-units --all "${show_services[@]}" | awk '{sub("dead", "not started", $4); printf("\n[\\e{cyan}%s\\e{reset}]", $4); for (i=5; i<=NF; i++) {if (i>5 || $i != "start") printf(" %s", $i)}}'
        } | set_issue "${services_issue}"
    fi
}

check_host_config() {
    local host_config_issue="80_host-config"
    if [ -f /var/run/agent-installer/host-config-failures ]; then
        {
            printf '\\e{lightred}Installation cannot proceed:\\e{reset}\n'
            cat /var/run/agent-installer/host-config-failures
        } | set_issue "${host_config_issue}"
    else
        clear_issue "${host_config_issue}"
    fi
}

while true; do
    check_services
    check_host_config
    sleep 5
done
