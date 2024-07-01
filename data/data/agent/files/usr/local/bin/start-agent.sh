#!/bin/bash

>&2 echo "Waiting for infra-env-id to be available"
INFRA_ENV_ID=""
until [[ $INFRA_ENV_ID != "" && $INFRA_ENV_ID != "null" ]]; do
    sleep 5
    >&2 echo "Querying assisted-service for infra-env-id..."
    INFRA_ENV_ID=$(curl -s -S "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs" -H "Authorization: ${AGENT_AUTH_TOKEN}" | jq -r .[0].id)
done
echo "Fetched infra-env-id and found: $INFRA_ENV_ID"

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

IMAGE=$(image_for agent-installer-node-agent)

echo "Using agent image: ${IMAGE} to start agent"

# use infra-env-id to have agent register this host with assisted-service
exec /usr/local/bin/agent --url "${SERVICE_BASE_URL}" --infra-env-id "${INFRA_ENV_ID}" --agent-version "${IMAGE}" --insecure=true
