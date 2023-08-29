#!/bin/bash
#
# restart-agent.sh 
# 
# Restarts agent.service and itself if it detects the assisted-service
# REST-API is offline. Being offline could mean the assisted-service
# pod or the rendezvous host was restarted. Because restarting 
# assisted-service clears its internal database. Any existing host/agent
# that had previously registered with assisted-service needs to restart
# its agent.service to reregister.

>&2 echo "Waiting for infra-env-id to be available"
INFRA_ENV_ID=""
until [[ $INFRA_ENV_ID != "" && $INFRA_ENV_ID != "null" ]]; do
    sleep 5
    >&2 echo "Querying assisted-service for infra-env-id..."
    INFRA_ENV_ID=$(curl -s -S "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs" | jq -r .[0].id)
done
echo "Fetched infra-env-id and found: $INFRA_ENV_ID"


>&2 echo "Restart agent.service if REST-API becomes unavailable"
INFRA_ENV_ID="tbd"
until [[ $INFRA_ENV_ID == "" || $INFRA_ENV_ID == "null" ]]; do
    sleep 5
    >&2 echo "Checking assisted-service is still online"
    INFRA_ENV_ID=$(curl -s -S "${SERVICE_BASE_URL}/api/assisted-install/v2/infra-envs" | jq -r .[0].id)
    echo "assisted-service is online with infra-env-id: $INFRA_ENV_ID"
done
echo "assisted-service interrupted, restarting agent.service"

sudo systemctl restart agent.service
sudo systemctl restart agent-restart.service