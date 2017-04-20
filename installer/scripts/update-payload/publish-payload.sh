#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/.."

# A script that make an Core Update channel to point to a target package.
# updateservicectl(https://github.com/coreos/updateservicectl) is required.
# Assume the payload.json is in the current working directory.

function print_usage() {
  echo "Usage:"
  echo "export COREUPDATE_USR=<user@coreos.com>"
  echo "export COREUPDATE_KEY=<coreupdate_key>"
  echo "$0 [tectonic-1.5|tectonic-1.5-staging]"
  exit 1
}

# main function

if [[ ${COREUPDATE_USR} == ""  || ${COREUPDATE_KEY} == "" || $# != 1 ]]; then
    print_usage
fi

which updateservicectl > /dev/null
if [[ $? == 0 ]]; then
    export UPDATESERVICECTL=$(which updateservicectl)
fi

if [[ ${UPDATESERVICECTL} == "" ]]; then
    echo "Require updateservicectl (https://github.com/coreos/updateservicectl)"
    exit 1
fi

set -e

payload=${DIR}/payload.json

if [ ! -f ${payload} ]; then
    echo "Expecting payload.json in the current directory"
    exit 1
fi

VERSION=${VERSION:-$(cat ${payload} | jq -r .version)}
CHANNEL=$1

echo "Publishing new payload ${VERSION} to channel ${CHANNEL}"

SERVER=${SERVER:-"https://tectonic.update.core-os.net"}
APPID=${APPID:-"6bc7b986-4654-4a0f-94b3-84ce6feb1db4"}

${UPDATESERVICECTL} --server ${SERVER} \
                    --key ${COREUPDATE_KEY} \
                    --user ${COREUPDATE_USR} \
                    channel update \
                    --channel ${CHANNEL} \
                    --app-id ${APPID} \
                    --version ${VERSION}
