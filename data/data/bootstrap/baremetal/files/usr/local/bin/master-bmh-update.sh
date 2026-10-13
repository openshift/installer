#!/bin/bash

set -euo pipefail

export KUBECONFIG=/opt/openshift/auth/kubeconfig

# Wait till the baremetalhosts are populated
until oc get baremetalhosts -n openshift-machine-api; do
   echo Waiting for BareMetalHosts CRD to appear...
   sleep 20
done

echo "Waiting for $CONTROL_PLANE_REPLICA_COUNT masters to become provisioned"
while [ "$(oc get bmh -n openshift-machine-api -l installer.openshift.io/role=control-plane -o json | jq '.items[].status.provisioning.state' | grep provisioned -c)" -lt "$CONTROL_PLANE_REPLICA_COUNT"  ]; do
    echo "Waiting for $CONTROL_PLANE_REPLICA_COUNT masters to become provisioned"
    oc get bmh -A || true
    sleep 20
done

# Pause control-plane management via ironic for Two-Node OpenShift with Fencing (TNF), since this conflicts with how TNF achieves HA
TNF_TOPOLOGY="DualReplica"

echo "Looking up control-plane topology"
CONTROL_PLANE_TOPOLOGY=$(oc get infrastructures.config.openshift.io cluster -o json | jq -r '.status.controlPlaneTopology')
if [ -z "${CONTROL_PLANE_TOPOLOGY}" ]; then
    echo "Error: couldn't lookup control-plane topology" >&2
    exit 1
elif [ "${CONTROL_PLANE_TOPOLOGY}" = "${TNF_TOPOLOGY}" ]; then
    echo "Control-plane topology set to '${TNF_TOPOLOGY}'; setting the control-plane hosts to detached"
    while ! oc annotate --overwrite bmh -n openshift-machine-api -l installer.openshift.io/role=control-plane baremetalhost.metal3.io/detached="" ; do
        sleep 5
        echo "Setting control-plane nodes to detached failed, retrying"
    done
else
    echo "Control-plane topology set to '${CONTROL_PLANE_TOPOLOGY}'; no further actions required"
fi

# Shut down ironic containers so that the API VIP can fail over to the control
# plane.
echo "Stopping provisioning services..."
systemctl --no-block stop ironic.service
while systemctl is-active metal3-baremetal-operator.service; do
    sleep 10
done

echo "Unpause all baremetal hosts"
while ! oc annotate --overwrite -n openshift-machine-api baremetalhosts --all "baremetalhost.metal3.io/paused-" ; do
    sleep 5
    echo "Unpause failed, retrying"
done
