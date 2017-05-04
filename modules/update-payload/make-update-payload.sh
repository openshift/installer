#!/bin/bash

# A script that helps generate the update payload for the given channel.
# It's using the manifests in the candidate assets dir.
# yaml2json(https://github.com/bronze1man/yaml2json) and jq are required.

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

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Invoking terraform to populate the templates..." >&2
pushd ${DIR}
while true; do echo "place_holder_for_terraform_input"; done | terraform apply ./
popd

ASSETS_DIR=${DIR}/generated

# TODO(yifan): Maybe put these files into separate dirs for each channel
# in the future.
deployments=(
  "tectonic-channel-operator.yaml"
  "kube-version-operator.yaml"
  "tectonic-prometheus-operator.yaml"
  "container-linux-update-operator.yaml"
)
appversions=(
  "app-version-kubernetes.yaml"
  "app-version-tectonic-monitoring.yaml"
)

echo "Creating update payload..." >&2
echo "Using deployments: [${deployments[*]}]" >&2
echo "Using app versions: [${appversions[*]}]" >&2

# Get the update payload version.
VERSION=$(yaml2json < ${ASSETS_DIR}/app-version-tectonic-cluster.yaml | jq .status.currentVersion)

# Get the deployments.
for f in ${deployments[*]}; do
  tmpfile=$(mktemp /tmp/deployment.XXXXXX)
  yaml2json < ${ASSETS_DIR}/${f} > ${tmpfile}
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
  name=$(yaml2json < ${ASSETS_DIR}/${f} | jq .metadata.name)
  desiredVersion=$(yaml2json < ${ASSETS_DIR}/${f} | jq .status.currentVersion)
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
cat <<EOF | jq . > ${DIR}/payload.json
{
  "version": ${VERSION},
  "deployments": ${DEPLOYMENTS},
  "desiredVersions": ${DESIRED_VERSIONS}
}
EOF

echo "Payload generated at [${DIR}/payload.json]" >&2
