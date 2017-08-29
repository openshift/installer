#!/bin/bash

if [ "$#" -ne "2" ]; then
    echo "Usage: $0 location destination"
    exit 1
fi

# shellcheck disable=SC2086,SC2154
/usr/bin/sudo /usr/bin/rkt run \
    --net=host --dns=host \
    --trust-keys-from-https ${awscli_image} \
    --volume=tmp,kind=host,source=/tmp --mount=volume=tmp,target=/tmp \
    --set-env="LOCATION=$1" \
    --exec=/bin/bash -- -c '
        REGION=$(wget -q -O - http://169.254.169.254/latest/meta-data/placement/availability-zone | sed s'/[a-zA-Z]$//')
        until /usr/bin/aws --region=$${REGION} s3 cp s3://$${LOCATION} /tmp/$${LOCATION//\//+}; do
          echo "Could not pull from S3, retrying in 5 seconds"
          sleep 5
        done
    '

# shellcheck disable=SC1001,SC1083,SC2086
/usr/bin/sudo mv /tmp/$${1//\//+} $2
