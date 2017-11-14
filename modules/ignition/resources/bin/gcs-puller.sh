#!/bin/bash
# shellcheck disable=SC2139,SC2154,SC2034,SC1083
set -x

if [ "$#" -ne "2" ]; then
  echo "Usage: $0 location destination"
  exit 1
fi

GSUTIL="docker run --net=host -v $HOME/.config:/.config -v /tmp:/gs ${gcloudsdk_image} gsutil"
ASSETS_ADDRESS=$${1}
ASSETS_NAME=$(basename $${1})
TARGET_FOLDER=$${2}

gcs_copy_assets(){
  $${GSUTIL} cp gs://$${ASSETS_ADDRESS} /gs/$${ASSETS_NAME}
}

gcs_remove_assets(){
  $${GSUTIL} rm gs://$${ASSETS_ADDRESS}
}

until gcs_copy_assets; do
    echo "Could not pull from GCS, retrying in 5 seconds"
    sleep 5
done

until gcs_remove_assets; do
  echo "Could not delete assets from GCS, retrying in 5 seconds"
  sleep 5
done

/usr/bin/sudo mv /tmp/$${ASSETS_NAME} $${TARGET_FOLDER}
