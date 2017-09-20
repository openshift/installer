#!/bin/bash
set -e

# Defer cleanup rkt containers and images.
trap "{ /usr/bin/rkt gc --grace-period=0; /usr/bin/rkt image gc --grace-period 0; } &> /dev/null" EXIT

mkdir -p /run/metadata
# shellcheck disable=SC2086,SC2154
/usr/bin/rkt run \
    --dns=host --net=host --trust-keys-from-https --interactive \
    \
    --set-env=CLUSTER_NAME=${cluster_name} \
    \
    --volume=metadata,kind=host,source=/run/metadata,readOnly=false \
    --mount=volume=metadata,target=/run/metadata \
    \
    --volume=detect-master,kind=host,source=/opt/detect-master.sh,readOnly=true \
    --mount=volume=detect-master,target=/detect-master.sh \
    \
    ${awscli_image} \
    --exec=/detect-master.sh

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
