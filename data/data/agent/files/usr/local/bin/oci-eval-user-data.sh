#!/usr/bin/env bash

set -euxo pipefail

chassis_asset_tag="$(dmidecode --string chassis-asset-tag)"
if [ "${chassis_asset_tag}" != "OracleCloud.com" ]
then
  echo "Not running in Oracle Cloud Infrastructure. Skipping."
  exit 0
fi

if ! user_data_b64=$(curl -H "Authorization: Bearer Oracle" --max-time 30 --retry 5 --retry-delay 30 --retry-connrefused -L http://169.254.169.254/opc/v2/instance/metadata/user_data)
then
  echo "Failed to retrieve user data from instance metadata"
  exit 1
fi

if grep -q "404 Not Found" <<< "${user_data_b64}"
then
  echo "No user data available"
  exit 0
fi

if ! user_data=$(echo -n "${user_data_b64}" | base64 --decode)
then
  echo "Failed to decode user data"
  exit 1
fi

bash <(echo "${user_data}")