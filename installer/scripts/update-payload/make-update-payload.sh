#!/bin/bash

# A script that helps generate the update payload for the given channel.
# It's using the manifests in the candidate assets dir.
# yaml2json(https://github.com/bronze1man/yaml2json) and jq are required.

if [[ $# != 1 ]]; then
    echo "Usage: $0 [ASSETS_DIR]" >&2
    exit 1
fi

which yaml2json > /dev/null
if [[ $? != 0 ]]; then
    echo "Require yaml2json (https://github.com/bronze1man/yaml2json)" >&2
    exit 1
fi

which jq > /dev/null
if [[ $? != 0 ]]; then
    echo "Require jq" >&2
    exit 1
fi

set -e

ASSETS_DIR=${1}

# TODO(yifan): Maybe put these files into separate dirs for each channel
# in the future.
deployments=(
  "tectonic-channel-operator-deployment.yaml"
  "kube-version-operator-deployment.yaml"
)
appversions=(
  "app-version-kubernetes.json"
)

echo "Creating update payload..." >&2
echo "Using deployments: ${deployments[*]}" >&2
echo "Using app versions: ${appversions[*]}" >&2

# Get the update payload version.
VERSION=$(cat ${ASSETS_DIR}/app-version-tectonic-cluster.json | jq .status.currentVersion)

# Get the deployments.
for f in ${deployments[*]}; do
  tmpfile=$(mktemp /tmp/deployment.XXXXXX)
  cat ${ASSETS_DIR}/${f} | yaml2json > ${tmpfile}
  tmpfiles+=(${tmpfile})
done

if [[ ${#tmpfiles[*]} > 0 ]]; then
    DEPLOYMENTS=$(jq -s . ${tmpfiles[*]})
    rm ${tmpfiles[*]}
fi

unset tmpfiles

# Get the desired versions.
for f in ${appversions[*]}; do
  tmpfile=$(mktemp /tmp/desiredVersion.XXXXXX)
  name=$(jq .metadata.name ${ASSETS_DIR}/${f})
  desiredVersion=$(jq .status.currentVersion ${ASSETS_DIR}/${f})
  cat <<EOF > ${tmpfile}
{
  "name": ${name},
  "version": ${desiredVersion}
}
EOF
  tmpfiles+=(${tmpfile})
done

if [[ ${#tmpfiles[*]} > 0 ]]; then
    DESIRED_VERSIONS=$(jq -s . ${tmpfiles[*]})
    rm ${tmpfiles[*]}
fi

# Create the final payload.
cat <<EOF | jq . 
{
  "version": ${VERSION},
  "deployments": ${DEPLOYMENTS},
  "desiredVersions": ${DESIRED_VERSIONS}
}
EOF
