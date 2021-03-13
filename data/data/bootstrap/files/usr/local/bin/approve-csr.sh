#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/bootstrap-service-record.sh

KUBECONFIG="${1}"

echo "Approving all CSR requests until bootstrapping is complete..."
while [ ! -f /opt/openshift/.bootkube.done ]
do
    oc --kubeconfig="$KUBECONFIG" get csr --no-headers | grep Pending | \
        awk '{print $1}' | \
        xargs --no-run-if-empty oc --kubeconfig="$KUBECONFIG" adm certificate approve
	sleep 20
done
