#!/bin/bash
set -e

echo "Waiting for assisted-service to be ready"
until curl --output /dev/null --silent --fail "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs" -H "Authorization: ${AGENT_AUTH_TOKEN}"; do
    printf '.'
    sleep 5
done
