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
pushd "${DIR}"
while true; do echo "place_holder_for_terraform_input"; done | terraform apply ./
popd

ASSETS_DIR=${DIR}/generated

# TODO(yifan): Maybe put these files into separate dirs for each channel
# in the future.
deployments=(
  "tectonic-channel-operator.yaml"
  "kube-version-operator.yaml"
  "tectonic-prometheus-operator.yaml"
  "tectonic-etcd-operator.yaml"
  "container-linux-update-operator.yaml"
)
appversions=(
  "app-version-kubernetes.yaml"
  "app-version-tectonic-monitoring.yaml"
  "app-version-tectonic-etcd.yaml"
)

echo "Creating update payload..." >&2
echo "Using deployments: [${deployments[*]}]" >&2
echo "Using app versions: [${appversions[*]}]" >&2

# Get the update payload version.
# shellcheck disable=SC2086
VERSION=$(yaml2json < ${ASSETS_DIR}/app-version-tectonic-cluster.yaml | jq .status.currentVersion)

# Get the deployments.
for f in ${deployments[*]}; do
  tmpfile=$(mktemp /tmp/deployment.XXXXXX)
  # shellcheck disable=SC2086
  yaml2json < ${ASSETS_DIR}/${f} > ${tmpfile}
  tmpfiles+=(${tmpfile})
done

# TODO: (ggreer) I'm pretty sure ">" is *not* what we want here
# shellcheck disable=SC2071
if [[ ${#tmpfiles[*]} > 0 ]]; then
    # shellcheck disable=SC2086
    DEPLOYMENTS=$(jq -s . ${tmpfiles[*]})
    # shellcheck disable=SC2086
    rm ${tmpfiles[*]}
fi

unset tmpfiles

# Get the desired versions.
for f in ${appversions[*]}; do
  tmpfile=$(mktemp /tmp/desiredVersion.XXXXXX)
  # shellcheck disable=SC2086
  name=$(yaml2json < ${ASSETS_DIR}/${f} | jq .metadata.name)
  # shellcheck disable=SC2086
  desiredVersion=$(yaml2json < ${ASSETS_DIR}/${f} | jq .status.currentVersion)
  # shellcheck disable=SC2086
  cat <<EOF > ${tmpfile}
{
  "name": ${name},
  "version": ${desiredVersion}
}
EOF
  tmpfiles+=(${tmpfile})
done

# TODO: (ggreer) I'm pretty sure ">" is *not* what we want here
# shellcheck disable=SC2071
if [[ ${#tmpfiles[*]} > 0 ]]; then
    # shellcheck disable=SC2086
    DESIRED_VERSIONS=$(jq -s . ${tmpfiles[*]})
    # shellcheck disable=SC2086
    rm ${tmpfiles[*]}
fi

# Create the final payload.
# shellcheck disable=SC2086
cat <<EOF | jq . > ${DIR}/payload.json
{
  "version": ${VERSION},
  "deployments": ${DEPLOYMENTS},
  "desiredVersions": ${DESIRED_VERSIONS}
}
EOF

echo "Payload generated at [${DIR}/payload.json]" >&2
