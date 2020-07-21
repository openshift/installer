#!/usr/bin/env bash

KUBECONFIG="${1}"

echo "Approving all CSR requests until bootstrapping is complete..."
while [ ! -f /opt/openshift/.bootkube.done ]
do
    oc --kubeconfig="$KUBECONFIG" get csrs --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --kubeconfig="$KUBECONFIG" adm certificate approve
    oc --kubeconfig="$KUBECONFIG" get csrs.v1beta1.certificates.k8s.io --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --kubeconfig="$KUBECONFIG" adm certificate approve
	sleep 20
done
