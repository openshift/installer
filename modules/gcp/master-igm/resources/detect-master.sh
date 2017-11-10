#!/bin/bash
# shellcheck disable=SC2034,SC2154
set -e

# This script assumes it'll be run inside google/cloud-sdk:alpine
# If running from host you can make gcloud available by uncommenting:
# alias gcloud="(docker images $DOCKER_IMAGE || docker pull $DOCKER_IMAGE) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config -v /var/run/docker.sock:/var/run/doker.sock -v /usr/bin/docker:/usr/bin/docker $DOCKER_IMAGE gcloud"
# shopt -s expand_aliases
# See /etc/profile.d/google-cloud-sdk.sh

apk update
apk add jq
# Wait for the IGM to run at the expected scale.
REGION=${region}
INSTANCE_GROUP_NAME=${instance_group_name}
CURRENT_INSTANCE_ID=$(curl -s -H 'Metadata-Flavor: Google' http://metadata.google.internal/computeMetadata/v1/instance/id)
while true; do
  INSTACE_GROUP_INFO=$(gcloud compute instance-groups managed describe "$${INSTANCE_GROUP_NAME}" --region "$${REGION}" --format json)
  INSTANCES_INFO=$(gcloud compute instance-groups managed list-instances "$${INSTANCE_GROUP_NAME}" --region "$${REGION}"  --format json)
  INSTANCE_ID_LIST=$(echo "$INSTANCES_INFO" | jq -r  "sort_by(.id) | .[].id")
  TARGET_SIZE=$(echo "$${INSTACE_GROUP_INFO}" | jq -r ".targetSize")
  CURRENT_SIZE=$(echo "$${INSTACE_GROUP_INFO}" | jq -r ".currentActions.none")
  if [ "$TARGET_SIZE" == "$CURRENT_SIZE" ]; then
    break
  fi

  echo "Waiting for the IGM to be at desired capacity (Desired: $TARGET_SIZE, Current: $CURRENT_SIZE)"
  sleep 15
done

BOOTKUBE_MASTER=$(echo "$INSTANCE_ID_LIST" | head -n1)
if [ "$BOOTKUBE_MASTER" != "$CURRENT_INSTANCE_ID" ]; then
    echo "This instance '$CURRENT_INSTANCE_ID' is not the bootkube master, '$BOOTKUBE_MASTER' is."
    echo -n "false" >/run/metadata/master
    exit 0
fi

echo -n "true" >/run/metadata/master
