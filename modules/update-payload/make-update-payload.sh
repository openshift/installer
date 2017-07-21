#!/bin/bash

# A script that helps generate the update payload for the given channel.
# It's using the manifests in the candidate assets dir.
# yaml2json(https://github.com/bronze1man/yaml2json) and jq are required.

if ! which yaml2json > /dev/null; then
    echo "Require yaml2json (https://github.com/bronze1man/yaml2json)" >&2
    exit 1
fi

if ! which jq > /dev/null; then
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

echo "Creating update payload..." >&2
echo "Using deployments:" >&2
for f in ${ASSETS_DIR}/operators/*.yaml; do
  basename "${f}" >&2
done

echo "Using app versions:" >&2
for f in ${ASSETS_DIR}/app_versions/*.yaml; do
  basename "${f}" >&2
done

# Get the update payload version.
# shellcheck disable=SC2086
VERSION=$(yaml2json < ${ASSETS_DIR}/app_versions/app-version-tectonic-cluster.yaml | jq .status.currentVersion)
# Don't include the meta app-version for sub-component's desired version.
rm "${ASSETS_DIR}/app_versions/app-version-tectonic-cluster.yaml"

# Get the deployments.
for f in ${ASSETS_DIR}/operators/*.yaml; do
  tmpfile=$(mktemp /tmp/deployment.XXXXXX)
  # shellcheck disable=SC2086
  yaml2json < ${f} > ${tmpfile}
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
for f in ${ASSETS_DIR}/app_versions/*.yaml; do
  tmpfile=$(mktemp /tmp/desiredVersion.XXXXXX)
  # shellcheck disable=SC2086
  name=$(yaml2json < ${f} | jq .metadata.name)
  # shellcheck disable=SC2086
  desiredVersion=$(yaml2json < ${f} | jq .status.currentVersion)
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

echo "Is this payload enforced? (y/N)" >&2

read -r payload_enforced

PAYLOAD_ENFORCED=false
if [[ ${payload_enforced} == "y" ]]; then
    echo "Payload will be enforced!" >&2
    PAYLOAD_ENFORCED=true
else
    echo "Payload will NOT be enforced" >&2
fi

# Create the final payload.
# shellcheck disable=SC2086
cat <<EOF | jq . > ${DIR}/payload.json
{
  "version": ${VERSION},
  "deployments": ${DEPLOYMENTS},
  "desiredVersions": ${DESIRED_VERSIONS},
  "enforced": ${PAYLOAD_ENFORCED}
}
EOF

echo "Payload generated at [${DIR}/payload.json]" >&2
