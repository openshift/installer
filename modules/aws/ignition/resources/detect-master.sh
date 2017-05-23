#!/bin/bash

REGION=$(wget -q -O - http://169.254.169.254/latest/meta-data/placement/availability-zone | sed s'/[a-zA-Z]$//')
INSTANCE_ID=$(wget -qO- http://169.254.169.254/latest/meta-data/instance-id)
ASG_NAME=$(aws autoscaling describe-auto-scaling-instances --region="$REGION" --instance-ids="$INSTANCE_ID" | jq -r ".AutoScalingInstances[0] .AutoScalingGroupName")

# Wait for the ASG to run at the expected scale.
while true; do
  ASG_DESCRIPTION=$(aws autoscaling describe-auto-scaling-groups --region="$REGION" --auto-scaling-group-names="$ASG_NAME")
  ASG_DESIRED_CAP=$(echo $ASG_DESCRIPTION | jq ".AutoScalingGroups[0] .DesiredCapacity")
  ASG_INSTANCE_IDS=$(echo $ASG_DESCRIPTION | jq -r ".AutoScalingGroups[0] .Instances | sort_by(.InstanceId) | .[].InstanceId")
  ASG_CURRENT_CAP=$(echo -e "$ASG_INSTANCE_IDS" | wc -l)

  if [ "$ASG_CURRENT_CAP" == "$ASG_DESIRED_CAP" ]; then
    break
  fi

  echo "Waiting for the ASG to be at desired capacity (Desired: $ASG_DESIRED_CAP, Current: $ASG_CURRENT_CAP)"
  sleep 15
done

API_HEALTHY=$(aws elb describe-instance-health --region="$REGION" --load-balancer-name "$CLUSTER_NAME-api-internal" | jq -r '[ .InstanceStates[] | select(.State | contains("InService")) ] | length > 1')

if [ "$API_HEALTHY" == "true" ]; then
    echo "Healthy API instances found, cluster is already installed."
    echo -n "false" >/run/metadata/master
    exit 0
fi

BOOTKUBE_MASTER=$(echo "$ASG_INSTANCE_IDS" | head -n1)

if [ "$BOOTKUBE_MASTER" != "$INSTANCE_ID" ]; then
    echo "This instance is not the bootkube master, '$BOOTKUBE_MASTER' is."
    echo -n "false" >/run/metadata/master
    exit 0
fi

echo -n "true" >/run/metadata/master
