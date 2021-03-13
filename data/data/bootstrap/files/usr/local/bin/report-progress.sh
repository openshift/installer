#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/bootstrap-service-record.sh

KUBECONFIG="${1}"

wait_for_existence() {
	while [ ! -e "${1}" ]
	do
		sleep 5
	done
}

record_service_stage_start "wait-for-bootstrap-complete"
echo "Waiting for bootstrap to complete..."
wait_for_existence /opt/openshift/.bootkube.done
record_service_stage_success

record_service_stage_start "report-bootstrap-complete"
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
record_service_stage_success
