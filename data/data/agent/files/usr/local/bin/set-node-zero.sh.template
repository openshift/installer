#!/bin/bash

set -e

IS_NODE_ZERO=$(ip -j address | jq '[.[].addr_info] | flatten | map(.local=="{{.NodeZeroIP}}") | any')

if [ "${IS_NODE_ZERO}" = "true" ]; then
    echo "Node 0 IP {{.NodeZeroIP}} found on this host" 1>&2

    NODE0_PATH=/etc/assisted-service/node0
    mkdir -p "$(dirname "${NODE0_PATH}")"

    NODE_ZERO_MAC=$(ip -j address | jq -r '.[] | select(.addr_info | map(select(.local == "{{.NodeZeroIP}}")) | any).address')
    echo "MAC Address for Node 0: ${NODE_ZERO_MAC}"

    cat >"${NODE0_PATH}" <<EOF
BOOTSTRAP_HOST_MAC=${NODE_ZERO_MAC}
EOF

    echo "Created file ${NODE0_PATH}"
fi
