#!/bin/bash
set -e

s3_clean() {
  # Delete Install assets from S3
  # shellcheck disable=SC2086,SC2154,SC2016
  /usr/bin/docker run \
    --volume /tmp:/tmp \
    --network=host \
    --env LOCATION="${assets_s3_location}" \
    --entrypoint=/bin/bash \
    ${awscli_image} \
    -c '
        REGION=$(wget -q -O - http://169.254.169.254/latest/meta-data/placement/availability-zone | sed '"'"'s/[a-zA-Z]$//'"'"')
        usr/bin/aws --region=$${REGION} s3 rm s3://$${LOCATION}
    '
}

# shellcheck disable=SC2086,SC2154
/usr/bin/docker run \
  --volume /run/metadata:/run/metadata \
  --volume /opt/detect-master.sh:/detect-master.sh:ro \
  --network=host \
  --env CLUSTER_NAME=${cluster_name} \
  --entrypoint=/detect-master.sh \
  ${awscli_image}

# Don't do anything if cluster is still in startup
STARTUP=$(cat /run/metadata/master)
if [ "$STARTUP" == "true" ]; then
    exit 0
fi

until s3_clean; do
  echo "failed to clean up S3 assets. retrying in 5 seconds."
  sleep 5
done
