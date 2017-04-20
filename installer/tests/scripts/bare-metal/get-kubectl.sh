#!/usr/bin/env bash
# USAGE: ./get-kubectl bin
# Get the kubectl client
set -eu

DEST=${1:-"bin"}
VERSION="v1.5.6"

URL="https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/amd64/kubectl"

if [[ ! -f $DEST/kubectl ]]; then
  mkdir -p ${DEST}
  curl -L -o ${DEST}/kubectl ${URL}
  chmod +x ${DEST}/kubectl
fi
