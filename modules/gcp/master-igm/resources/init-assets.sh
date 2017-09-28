#!/bin/bash
set -e

# Download the assets from GCS
# shellcheck disable=SC2086,SC2154
/usr/bin/bash /opt/gcs-puller.sh ${assets_gcs_location} /var/tmp/tectonic.zip
unzip -o -d /var/tmp/tectonic/ /var/tmp/tectonic.zip
rm /var/tmp/tectonic.zip
# make files in /opt/tectonic available atomically
mv /var/tmp/tectonic /opt/tectonic

# Populate the kubelet.env file.
mkdir -p /etc/kubernetes
# shellcheck disable=SC2154
echo "KUBELET_IMAGE_URL=${kubelet_image_url}" > /etc/kubernetes/kubelet.env
# shellcheck disable=SC2154
echo "KUBELET_IMAGE_TAG=${kubelet_image_tag}" >> /etc/kubernetes/kubelet.env

exit 0
