#!/bin/bash

check_host_config() {
    local host_config_issue="/etc/issue.d/80_host-config.issue"
    if [ -f /var/run/agent-installer/host-config-failures ]; then
        printf '\n\\e{lightred}Installation cannot proceed:\\e{reset}\n%s\n' "$(cat /var/run/agent-installer/host-config-failures)" >"${host_config_issue}"
    else
        rm -f "${host_config_issue}"
    fi
}

while true; do
    check_host_config
    sleep 5
done
