#!/usr/bin/env bash

KUBECONFIG="${1}"

echo "Approving all CSR requests until bootstrapping is complete..."
while [ ! -f /opt/openshift/bootkube.done ]
do
    oc --config="$KUBECONFIG" get csr --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --config="$KUBECONFIG" adm certificate approve
	sleep 20
done
