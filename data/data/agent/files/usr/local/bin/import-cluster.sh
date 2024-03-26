#!/bin/bash
set -e

echo "Waiting for assisted-service to be ready"
until curl --output /dev/null --silent --fail "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs"; do
    printf '.'
    sleep 5
done

echo "Importing cluster from ${API_VIP_DNSNAME}"
curl -d '{"name":"abi", "api_vip_dnsname":"'"${API_VIP_DNSNAME}"'", "openshift_cluster_id":"764eb48f-8ffb-4126-b7ed-0ca746365654"}' -H "Content-Type: application/json" -X POST "${SERVICE_BASE_URL}/api/assisted-install/v2/clusters/import"