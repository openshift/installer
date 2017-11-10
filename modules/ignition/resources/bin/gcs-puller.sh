#!/bin/bash
# shellcheck disable=SC2139,SC2154,SC2034,SC1083
set -x

if [ "$#" -ne "2" ]; then
  echo "Usage: $0 location destination"
  exit 1
fi

# Overriding /etc/profile.d/google-cloud-sdk.sh
DOCKER_IMAGE=${gcloudsdk_image}
alias gsutil="(docker images $DOCKER_IMAGE  || docker pull $DOCKER_IMAGE) > /dev/null; \
docker run -i --net=host -v $HOME/.config:/.config -v /tmp:/gs $DOCKER_IMAGE gsutil"
shopt -s expand_aliases

assets=$(basename $${1})
gsutil cp gs://$${1} /gs/$${assets}
gsutil rm gs://$${1}
/usr/bin/sudo mv /tmp/$${assets} $${2}
