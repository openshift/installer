#!/bin/bash
# shellcheck disable=SC2086,SC2154
set -e

detect_master() {
  mkdir -p /run/metadata
  /usr/bin/docker run \
      --volume /run/metadata:/run/metadata \
      --volume /opt/detect-master.sh:/detect-master.sh:ro \
      --network=host \
      --entrypoint=/detect-master.sh \
      ${gcloudsdk_image}
}

until detect_master; do
  echo "failed to detect master; retrying in 5 seconds"
  sleep 5
done

MASTER=$(cat /run/metadata/master)
if [ "$MASTER" != "true" ]; then
    exit 0
fi
# Download the assets from GCS
/usr/bin/bash /opt/gcs-puller.sh ${assets_gcs_location} /var/tmp/tectonic.zip
unzip -o -d /var/tmp/tectonic/ /var/tmp/tectonic.zip
rm /var/tmp/tectonic.zip
# make files in /opt/tectonic available atomically
mv /var/tmp/tectonic /opt/tectonic

exit 0
