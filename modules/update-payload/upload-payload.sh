#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# A script that uploads the payload.json to aws s3 bucket, and
# create a new package on the core update server.
# jq, docker and updateservicectl (https://github.com/coreos/updateservicectl)
# are required.
# Assume the payload.json is in the current working directory.

function print_usage() {
  echo "Usage:"
  echo "export AWS_ACCESS_KEY_ID=<id>"
  echo "export AWS_SECRET_ACCESS_KEY=<key>"
  echo "export COREUPDATE_USR=<user@coreos.com>"
  echo "export COREUPDATE_KEY=<coreupdate_key>"
  echo "$0"
  exit 1
}

# main function

if [[ ${AWS_ACCESS_KEY_ID} == "" || ${AWS_SECRET_ACCESS_KEY} == "" || ${COREUPDATE_USR} == ""  || ${COREUPDATE_KEY} == "" ]]; then
    print_usage
fi

if ! which jq > /dev/null; then
    echo "Require jq"
    exit 1
fi

if which updateservicectl > /dev/null; then
    export UPDATESERVICECTL
    UPDATESERVICECTL=$(which updateservicectl)
fi

if [[ ${UPDATESERVICECTL} == "" ]]; then
    echo "Require updateservicectl (https://github.com/coreos/updateservicectl)"
    exit 1
fi

AWSCLI=${AWSCLI:-"aws"}

set -e

payload=${DIR}/payload.json

for f in "${payload}" "${payload}.sig"; do
    if [[ ! -f "${f}" ]]; then
        echo "Expecting ${f} in the current directory" >&2
        exit 1
    fi
done

# shellcheck disable=SC2002,SC2086
VERSION=${VERSION:-$(cat ${payload} | jq -r .version)}

if [[ ${VERSION} == "" ]]; then
    echo "Invalid payload format"
    exit 1
fi

DESTINATION=${DESTINATION:-"${VERSION}.json"}
BUCKET=${BUCKET:-"tectonic-update-payload"}
PAYLOAD_URL="https://s3-us-west-2.amazonaws.com/${BUCKET}/${DESTINATION}"

echo "Uploading payload to \"${PAYLOAD_URL}\", version: \"${VERSION}\""

${AWSCLI} s3 cp "${payload}" "s3://${BUCKET}/${DESTINATION}"
${AWSCLI} s3 cp "${payload}.sig" "s3://${BUCKET}/${DESTINATION}.sig"

SERVER=${SERVER:-"https://tectonic.update.core-os.net"}
APPID=${APPID:-"6bc7b986-4654-4a0f-94b3-84ce6feb1db4"}

echo "Payload successfully uploaded"

echo "Creating package ${VERSION} on Core Update server ${SERVER} for ${APPID}"

# shellcheck disable=SC2086,SC2154
${UPDATESERVICECTL} --server ${SERVER} \
                    --key ${COREUPDATE_KEY} \
                    --user ${COREUPDATE_USR} \
                    package create \
                    --app-id ${APPID} \
                    --url ${PAYLOAD_URL} \
                    --version ${VERSION} \
                    --file ${payload}

echo "Packaged successfully created"
