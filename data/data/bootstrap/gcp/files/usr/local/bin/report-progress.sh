#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/wait-for-ha-api.sh

KUBECONFIG="${1}"

wait_for_existance() {
	while [ ! -e "${1}" ]
	do
		sleep 5
	done
}

echo "Waiting for cb-bootstrap to complete..."
wait_for_existance /opt/openshift/cb-bootstrap.done

## remove the routes setup so that we can open up the blackhole
systemctl stop gcp-routes.service

echo "Waiting for bootstrap to complete..."
wait_for_existance /opt/openshift/.bootkube.done

## wait for API to be available
wait_for_ha_api

echo "Reporting install progress..."
while ! oc --kubeconfig="$KUBECONFIG" create -f - <<-EOF
	apiVersion: v1
	kind: ConfigMap
	metadata:
	  name: bootstrap
	  namespace: kube-system
	data:
	  status: complete
EOF
do
	sleep 5
done
