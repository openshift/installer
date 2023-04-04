#!/bin/bash
set -e

# shellcheck disable=SC1091
source /etc/assisted/rendezvous-host.env

echo "Waiting for assisted-service to be ready"
until curl --output /dev/null --silent --fail "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs"; do
    printf '.'
    sleep 5
done
