#!/bin/bash

if [ "${IPTABLES}" != "iptables" ] && [ "${IPTABLES}" != "ip6tables" ]; then
    echo "Environment variable \$IPTABLES must be set to iptables binary" >&2
    exit 1
fi

PORT="$2"

do_iptables() {
    if [ -z "${PORT}" ]; then
        echo "Port must be specified" >&2
        exit 1
    fi
    if [ -z "${PROVISIONING_INTERFACE}" ]; then
        echo "Environment variable \$PROVISIONING_INTERFACE must be specified" >&2
        exit 1
    fi

    "${IPTABLES}" "$1" INPUT -i "${PROVISIONING_INTERFACE}" -p tcp -m tcp --dport "${PORT}" -j ACCEPT
}

case "$1" in
    --enable-port)
        do_iptables -I
        ;;
    --disable-port)
        do_iptables -D
        ;;
    *)
        echo "Invalid command \"$1\"" >&2
        exit 1
        ;;
esac
