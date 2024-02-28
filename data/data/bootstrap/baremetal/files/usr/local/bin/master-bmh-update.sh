#!/bin/bash

set -euo pipefail

export KUBECONFIG=/opt/openshift/auth/kubeconfig-loopback

# Wait till the baremetalhosts are populated
until oc get baremetalhosts -n openshift-machine-api; do
   echo Waiting for BareMetalHosts CRD to appear...
   sleep 20
done

while [ "$(oc get bmh -n openshift-machine-api -o name | wc -l)" -lt 1  ]; do
    echo "Waiting for bmh"
    sleep 20
done

while [ "$(oc get bmh -n openshift-machine-api -l installer.openshift.io/role=master -o json | jq '.items[].status.provisioning.state' | grep -v provisioned -c)" -gt 0  ]; do
    echo "Waiting for masters to become provisioned"
    oc get bmh -A
    sleep 20
done

# Shut down ironic containers so that the API VIP can fail over to the control
# plane.
echo "Stopping provisioning services..."
systemctl --no-block stop ironic.service
while systemctl is-active metal3-baremetal-operator.service; do
    sleep 10
done

echo "Unpause all baremetal hosts"
oc annotate --overwrite -n openshift-machine-api baremetalhosts --all "baremetalhost.metal3.io/paused-"
