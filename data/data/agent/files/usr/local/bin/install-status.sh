#!/bin/bash

inactive_services() {
    local services="assisted-service.service create-cluster-and-infraenv.service apply-host-config.service start-cluster-installation.service"
    for s in ${services}; do
        if ! systemctl is-active "${s}" >/dev/null; then
            printf "%s " "${s}"
        fi
    done
}

check_services() {
    local services_issue="/etc/issue.d/70_agent-services.issue"
    not_started="$(inactive_services)"
    if [ -z "${not_started}" ]; then
        rm -f "${services_issue}"
    else
        read -ra show_services <<<"${not_started}"
        {
            printf '\\e{cyan}Waiting for services:\\e{reset}\n'
            systemctl --no-legend list-units --all "${show_services[@]}" | awk '{printf("[\\e{cyan}%s\\e{reset}]", $4); for (i=5; i<=NF; i++) {if (i>5 || $i != "start") printf(" %s", $i)}; print ""}'
        } >"${services_issue}"
    fi
}

check_host_config() {
    local host_config_issue="/etc/issue.d/80_host-config.issue"
    if [ -f /var/run/agent-installer/host-config-failures ]; then
        printf '\n\\e{lightred}Installation cannot proceed:\\e{reset}\n%s\n' "$(cat /var/run/agent-installer/host-config-failures)" >"${host_config_issue}"
    else
        rm -f "${host_config_issue}"
    fi
}

while true; do
    check_services
    check_host_config
    sleep 5
done
