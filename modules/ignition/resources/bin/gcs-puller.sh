#!/bin/bash
set -x

if [ "$#" -ne "2" ]; then
    echo "Usage: $0 location destination"
    exit 1
fi

docker pull google/cloud-sdk > /dev/null
# shellcheck disable=SC2034,SC1083
gsutil="docker run -t --net=host -v /tmp:/gs google/cloud-sdk gsutil"
# shellcheck disable=SC2034,SC1083
assets=$(basename $${1})
# shellcheck disable=SC2034,SC1083
$${gsutil} cp gs://$${1} /gs/$${assets}
# shellcheck disable=SC2034,SC1083
/usr/bin/sudo mv /tmp/$${assets} $${2}
