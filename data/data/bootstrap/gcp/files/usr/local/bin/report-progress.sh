#!/usr/bin/env bash

KUBECONFIG="${1}"

wait_for_existance_file() {
	while [ ! -e "${1}" ]
	do
		sleep 5
	done
}

wait_for_existance_k8s_object() {
	while ! oc --config="$KUBECONFIG" -n "${1}" get "${2}"
	do
		sleep 5
	done
}

echo "Waiting for bootstrap to complete..."
wait_for_existance_file /opt/openshift/.openshift.done
wait_for_existance_k8s_object kube-system event/bootstrap-success

echo "Reporting install progress..."
while ! oc --config="$KUBECONFIG" create -f - <<-EOF
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
