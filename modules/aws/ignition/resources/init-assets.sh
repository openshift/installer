#!/bin/bash
set -e

# Defer cleanup rkt containers and images.
trap "{ /usr/bin/rkt gc --grace-period=0; /usr/bin/rkt image gc --grace-period 0; } &> /dev/null" EXIT

mkdir -p /run/metadata
/usr/bin/rkt run \
    --dns=host --net=host --trust-keys-from-https --interactive \
    --volume=metadata,kind=host,source=/run/metadata,readOnly=false \
    --mount=volume=metadata,target=/run/metadata \
    ${awscli_image} \
    --exec=/bin/bash -- -c '
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

        # Fetch the Instance IDs of the ASG, sorted lexically, and write them
        # to disk, alongside with the ID of this instance.
        echo "$INSTANCE_ID" > /run/metadata/ec2_instance_id
        echo "$ASG_INSTANCE_IDS" > /run/metadata/ec2_asg_instances_ids
    '

# Exit if this host is not the first of its ASG.
INSTANCE_ID=$(cat /run/metadata/ec2_instance_id)
BOOTKUBE_MASTER=$(cat /run/metadata/ec2_asg_instances_ids | head -n1)
if [ "$BOOTKUBE_MASTER" != "$INSTANCE_ID" ]; then
    echo "This instance is not the bootkube master, '$BOOTKUBE_MASTER' is. Exiting."
    exit 0
fi

# Download the assets from S3.
/usr/bin/bash /opt/s3-puller.sh ${assets_s3_location} /opt/tectonic/tectonic.zip
unzip -o -d /opt/tectonic/ /opt/tectonic/tectonic.zip
rm /opt/tectonic/tectonic.zip

# Populate the kubelet.env file.
mkdir -p /etc/kubernetes
echo "KUBELET_IMAGE_URL=${kubelet_image_url}" > /etc/kubernetes/kubelet.env
echo "KUBELET_IMAGE_TAG=${kubelet_image_tag}" >> /etc/kubernetes/kubelet.env

exit 0
