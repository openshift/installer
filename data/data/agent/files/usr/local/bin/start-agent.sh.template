#!/bin/bash

source common.sh

>&2 echo "Waiting for infra-env-id to be available"
INFRA_ENV_ID=""
until [[ $INFRA_ENV_ID != "" && $INFRA_ENV_ID != "null" ]]; do
    sleep 5
    >&2 echo "Querying assisted-service for infra-env-id..."
    INFRA_ENV_ID=$(curl '{{.ServiceBaseURL}}/api/assisted-install/v2/infra-envs' | jq .[0].id)
done
# trim quotes
INFRA_ENV_ID=${INFRA_ENV_ID//\"}
echo "Fetched infra-env-id and found: $INFRA_ENV_ID"

# use infra-env-id to have agent register this host with assisted-service
exec /usr/local/bin/agent --url '{{.ServiceBaseURL}}' --infra-env-id "$INFRA_ENV_ID" --agent-version quay.io/edge-infrastructure/assisted-installer-agent:latest --insecure=true