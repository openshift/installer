#!/bin/bash
set -e
set -o pipefail

# Download the assets from S3.
# shellcheck disable=SC2154
/usr/bin/bash /opt/s3-puller.sh "${assets_s3_location}" /var/tmp/tectonic.zip
unzip -o -d /var/tmp/tectonic/ /var/tmp/tectonic.zip
rm /var/tmp/tectonic.zip
# make files in /opt/tectonic available atomically
mv /var/tmp/tectonic /opt/tectonic

exit 0
