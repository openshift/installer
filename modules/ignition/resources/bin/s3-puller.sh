#!/bin/bash
set -e
set -o pipefail

if [ "$#" -ne "2" ]; then
    echo "Usage: $0 location destination"
    exit 1
fi

# shellcheck disable=SC2034
LOCATION=$1
# shellcheck disable=SC1083,SC1001
LOCATION_ESCAPED=$${1//\//+}
DESTINATION=$2

s3_pull() {
    # shellcheck disable=SC2086,SC2154,SC2016
    /usr/bin/docker run \
        --volume /tmp:/tmp \
        --network=host \
        --env LOCATION=$LOCATION \
        --env LOCATION_ESCAPED=$LOCATION_ESCAPED \
        --entrypoint=/bin/bash \
        ${awscli_image} \
        -c '
            set -e
            set -o pipefail
            REGION=$(wget -q -O - http://169.254.169.254/latest/meta-data/placement/availability-zone | sed '"'"'s/[a-zA-Z]$//'"'"')
            /usr/bin/aws --region="$REGION" s3 cp s3://"$LOCATION" /tmp/"$LOCATION_ESCAPED"
        '
}

until s3_pull; do
    echo "failed to pull from S3; retrying in 5 seconds"
    sleep 5
done

/usr/bin/sudo mv /tmp/"$LOCATION_ESCAPED" "$DESTINATION"
