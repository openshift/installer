#!/bin/bash

source common.sh

HOST=$(get_host)
echo Using hostname ${HOST} 1>&2

if [[ ${HOST} == {{.NodeZeroIP}} ]] ;then

    NODE0_PATH=/etc/assisted-service/node0
    mkdir -p "$(dirname "${NODE0_PATH}")"

    NODE_ZERO_MAC=$(ip -j address | jq -r '.[] | select(.addr_info | map(select(.local == "{{.NodeZeroIP}}")) | any).address')
    echo "MAC Address for Node 0: ${NODE_ZERO_MAC}"

    cat >"${NODE0_PATH}" <<EOF
BOOTSTRAP_HOST_MAC=${NODE_ZERO_MAC}
EOF

    echo "Created file ${NODE0_PATH}"
fi
