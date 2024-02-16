#!/bin/bash

set -euo pipefail

export KUBECONFIG=/opt/openshift/auth/kubeconfig-loopback

# Wait till the baremetalhosts are populated
until oc get baremetalhosts -n openshift-machine-api; do
   echo Waiting for BareMetalHosts to appear...
   sleep 20
done

AUTH_DIR=/opt/metal3/auth
set +x
ironic_url="$(printf 'http://%s:%s@localhost:6385/v1' "$(cat "${AUTH_DIR}/ironic/username")" "$(cat "${AUTH_DIR}/ironic/password")")"

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

# Shut down ironic containers so that the API VIP can fail over to the control
# plane.
echo "Stopping provisioning services..."
systemctl stop ironic.service

echo "Wait for control plane to fail over"
sleep 30

echo "Unpause all baremetal hosts"
oc annotate --overwrite -n openshift-machine-api baremetalhosts --all "baremetalhost.metal3.io/paused-"
