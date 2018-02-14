#!/bin/bash
set -e

s3_clean() {
    # instead of simply removing the remote assets.zip,
    # overwrite it with a zero byte file, such that terraform doesn't
    # detect deletion but rather a simple change which it can ignore.
    touch /tmp/assets.zip

    # shellcheck disable=SC2086,SC2154,SC2016
    /usr/bin/docker run \
        --volume /tmp:/tmp \
        --network=host \
        --env LOCATION="${assets_s3_location}" \
        --entrypoint=/bin/bash \
        ${awscli_image} \
        -c '
            set -e
            set -o pipefail
            REGION=$(wget -q -O - http://169.254.169.254/latest/meta-data/placement/availability-zone | sed '"'"'s/[a-zA-Z]$//'"'"')
            /usr/bin/aws --region="$REGION" s3 cp /tmp/assets.zip s3://"$LOCATION"
        '
}

until s3_clean; do
  echo "failed to clean up S3 assets. retrying in 5 seconds."
  sleep 5
done
