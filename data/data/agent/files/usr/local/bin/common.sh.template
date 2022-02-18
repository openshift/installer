#!/bin/bash

get_host() {
    local default_gateway
    local host_ip

    default_gateway=$(ip -j route show default | jq -r '.[0].gateway')
    host_ip=$(ip -j route get "${default_gateway}" | jq -r '.[0].prefsrc')

    local host_fmt="%s"
    if [[ ${host_ip} =~ : ]]; then
        host_fmt="[%s]"
    fi

    # shellcheck disable=SC2059
    printf "${host_fmt}" "${host_ip}"
}
