#!/bin/sh

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

export KUBECONFIG=/etc/kubernetes/kubeconfig

# Wait till the baremetalhosts are populated
until oc get baremetalhosts -n openshift-machine-api; do
   echo Waiting for BareMetalHosts to appear...
   sleep 20
done

# Wait for a master to appear.
while [ "$(curl -s http://localhost:6385/v1/nodes | jq '.nodes[] | .uuid' | wc -l)" -lt 1 ]; do
    echo waiting for a master node to show up
    sleep 20
done

# Wait for the nodes to become active after introspection.
# Probably don't need this but I want to be 100% sure.
while curl -s http://localhost:6385/v1/nodes | jq '.nodes[] | .provision_state' | grep -v active; do
    echo Waiting for nodes to become active
    sleep 20
done

echo Nodes are all active

BAREMETAL_OPERATOR_IMAGE=$(image_for baremetal-operator)

for node in $(curl -s http://localhost:6385/v1/nodes | jq -r '.nodes[] | .uuid'); do
    name=$(curl -H "X-OpenStack-Ironic-API-Version: 1.9" -s "http://localhost:6385/v1/nodes/$node" | jq -r .name)
    echo "Host $name, UUID: $node"
    # And use the baremetal operator tool to load the introspection data into
    # the BareMetalHost CRs as annotations, which BMO then picks up.
    HARDWARE_DETAILS=$(podman run --quiet --net=host \
        --rm \
        --name baremetal-operator \
        --entrypoint /get-hardware-details \
        "${BAREMETAL_OPERATOR_IMAGE}" \
        http://localhost:5050/v1 "$node" | jq '{hardware: .}')

     oc annotate --overwrite -n openshift-machine-api baremetalhosts "$name" 'baremetalhost.metal3.io/status'="$HARDWARE_DETAILS"
done
