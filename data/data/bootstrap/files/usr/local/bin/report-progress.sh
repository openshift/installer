#!/usr/bin/env bash

KUBECONFIG="${1}"

wait_for_existance() {
	while [ ! -e "${1}" ]
	do
		sleep 5
	done
}

echo "Waiting for bootstrap to complete..."
wait_for_existance /opt/openshift/.bootkube.done

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
