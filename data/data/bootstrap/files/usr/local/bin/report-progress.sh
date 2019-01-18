#!/usr/bin/env bash

KUBECONFIG="${1}"
NAME="${2}"
MESSAGE="${3}"

wait_for_existance() {
	while [ ! -e "${1}" ]
	do
		sleep 5
	done
}

echo "Waiting for bootstrap to complete..."
wait_for_existance /opt/openshift/.bootkube.done
wait_for_existance /opt/openshift/.openshift.done

echo "Reporting install progress..."
timestamp="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
while ! oc --config="$KUBECONFIG" create -f - <<-EOF
	apiVersion: v1
	kind: Event
	metadata:
	  name: "${NAME}"
	  namespace: kube-system
	involvedObject:
	  namespace: kube-system
	message: "${MESSAGE}"
	firstTimestamp: "${timestamp}"
	lastTimestamp: "${timestamp}"
	count: 1
	source:
	  component: cluster
	  host: $(hostname)
EOF
do
	sleep 5
done
