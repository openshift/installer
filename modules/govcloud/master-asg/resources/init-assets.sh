#!/bin/bash
set -e
set -o pipefail

detect_master() {
    mkdir -p /run/metadata
    # shellcheck disable=SC2086,SC2154
    /usr/bin/docker run \
        --volume /run/metadata:/run/metadata \
        --volume /opt/detect-master.sh:/detect-master.sh:ro \
        --network=host \
        --env CLUSTER_NAME=${cluster_name} \
        --entrypoint=/detect-master.sh \
        ${awscli_image}
}

until detect_master; do
    echo "failed to detect master; retrying in 5 seconds"
    sleep 5
done

MASTER=$(cat /run/metadata/master)
if [ "$MASTER" != "true" ]; then
    exit 0
fi

# Download the assets from S3.
# shellcheck disable=SC2154
/usr/bin/bash /opt/s3-puller.sh "${assets_s3_location}" /var/tmp/tectonic.zip
unzip -o -d /var/tmp/tectonic/ /var/tmp/tectonic.zip
rm /var/tmp/tectonic.zip
# make files in /opt/tectonic available atomically
mv /var/tmp/tectonic /opt/tectonic

exit 0
