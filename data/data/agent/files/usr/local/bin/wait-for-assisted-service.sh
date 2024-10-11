#!/bin/bash
set -e

# shellcheck disable=SC1091
source "common.sh"

echo "Waiting for assisted-service to be ready"

until curl_assisted_service "/infra-envs" GET "$USER_AUTH_TOKEN" -o /dev/null --silent --fail; do
    printf '.'
    sleep 5
done
