#!/usr/bin/env bash

set -e
set -x

chassis_asset_tag="$(dmidecode --string chassis-asset-tag)"
if [ "${chassis_asset_tag}" != "OracleCloud.com" ]
then
  echo "Not running in Oracle Cloud Infrastructure. Skipping."
  exit 0
fi

if [ ! -d "/sys/firmware/ibft" ]
then
  echo "No IBFT configuration found. Skipping."
  exit 0
fi

MTU=9000

function get_if_name_from_mac_address {
  mac_address="${1}"
  ip -json link | jq --raw-output --arg mac_address "${mac_address}" '. | map(select(.address==($mac_address|ascii_downcase))) | .[].ifname'
}

# /opc/v2/vnics endpoint returns something that will look like the following
# structure:
# [
#   {
#     "macAddr": "00:10:e0:ec:72:fc",
#     "nicIndex": 0,
#     "privateIp": "10.0.29.201",
#     "subnetCidrBlock": "10.0.16.0/20",
#     "virtualRouterIp": "10.0.16.1",
#     "vlanTag": 0,
#     "vnicId": "ocid1.vnic.oc1.us-sanjose-1.abzwuljrppq34sbvgltddp7wujxwqw6xb7zjkwg54oaewx5mc4wr5cgtdzna"
#   },
#   {
#     "macAddr": "00:10:e0:ec:72:fd",
#     "nicIndex": 1,
#     "privateIp": "10.0.32.210",
#     "subnetCidrBlock": "10.0.32.0/20",
#     "virtualRouterIp": "10.0.32.1",
#     "vlanTag": 0,
#     "vnicId": "ocid1.vnic.oc1.us-sanjose-1.abzwuljrsndaptsyq5mppfsoaqoun3gbvnpngcaspybo2nbpcmrozx25jenq"
#   }
# ]

vnics=$(curl --silent -H "Authorization: Bearer Oracle" -L http://169.254.169.254/opc/v2/vnics/)
secondary_if_mac_address=$(jq -r '.[1].macAddr' <<< "${vnics}")
secondary_if_name=$(get_if_name_from_mac_address "${secondary_if_mac_address}")
secondary_if_ip_address=$(jq -r '.[1].privateIp' <<< "${vnics}")
secondary_if_default_gateway=$(jq -r '.[1].virtualRouterIp' <<< "${vnics}")
secondary_if_subnet=$(jq -r '.[1].subnetCidrBlock' <<< "${vnics}")
secondary_if_subnet_size=$(cut -f 2 -d '/' <<< "${secondary_if_subnet}")

if [ ! -f "/etc/NetworkManager/system-connections/${secondary_if_name}.nmconnection" ]
then
  nmcli connection add con-name "${secondary_if_name}" ifname "${secondary_if_name}" type ethernet ip4 "${secondary_if_ip_address}/${secondary_if_subnet_size}" gw4 "${secondary_if_default_gateway}"
  nmcli connection modify "${secondary_if_name}" ethernet.mtu ${MTU}
  nmcli connection modify "${secondary_if_name}" ipv4.route-metric 0 # make this interface the default interface
  nmcli connection modify "${secondary_if_name}" connection.autoconnect true

  nmcli connection reload
  nmcli connection up "${secondary_if_name}"
fi