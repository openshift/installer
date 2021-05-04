#!/bin/bash

set -euo pipefail

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

export KUBECONFIG=/opt/openshift/auth/kubeconfig-loopback

# Wait till the baremetalhosts are populated
until oc get baremetalhosts -n openshift-machine-api; do
   echo Waiting for BareMetalHosts to appear...
   sleep 20
done

AUTH_DIR=/opt/metal3/auth
ironic_url="$(printf 'http://%s:%s@localhost:6385/v1' "$(cat "${AUTH_DIR}/ironic/username")" "$(cat "${AUTH_DIR}/ironic/password")")"
inspector_url="$(printf 'http://%s:%s@localhost:5050/v1' "$(cat "${AUTH_DIR}/ironic-inspector/username")" "$(cat "${AUTH_DIR}/ironic-inspector/password")")"

# Wait for a master to appear.
while [ "$(curl -s "${ironic_url}/nodes" | jq '.nodes[] | .uuid' | wc -l)" -lt 1 ]; do
    echo waiting for a master node to show up
    sleep 20
done

# Wait for the nodes to become active after introspection.
# Probably don't need this but I want to be 100% sure.
while curl -s "${ironic_url}/nodes" | jq '.nodes[] | .provision_state' | grep -v active; do
    echo Waiting for nodes to become active
    sleep 20
done

echo Nodes are all active

BAREMETAL_OPERATOR_IMAGE=$(image_for baremetal-operator)

for node in $(curl -s "${ironic_url}/nodes" | jq -r '.nodes[] | .uuid'); do
    name=$(curl -H "X-OpenStack-Ironic-API-Version: 1.9" -s "${ironic_url}/nodes/${node}" | jq -r .name)
    echo "Host $name, UUID: $node"
    # And use the baremetal operator tool to load the introspection data into
    # the BareMetalHost CRs as annotations, which BMO then picks up.
    HARDWARE_DETAILS=$(podman run --quiet --net=host \
        --rm \
        --entrypoint /get-hardware-details \
        "${BAREMETAL_OPERATOR_IMAGE}" \
        "${inspector_url}" "$node" | jq '{hardware: .}')

     oc annotate --overwrite -n openshift-machine-api baremetalhosts "$name" 'baremetalhost.metal3.io/status'="$HARDWARE_DETAILS" 'baremetalhost.metal3.io/paused-'
done

touch /opt/openshift/.master-bmh-update.done