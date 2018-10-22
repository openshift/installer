#!/usr/bin/env bash
set -e

KUBECONFIG="${1}"
NAME="${2}"
MESSAGE="${3}"
TIMESTAMP="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

echo "Reporting install progress..."

oc --config="$KUBECONFIG" create -f - <<EOF
apiVersion: v1
kind: Event
metadata:
  name: "${NAME}"
  namespace: kube-system
involvedObject:
  namespace: kube-system
message: "${MESSAGE}"
firstTimestamp: "${TIMESTAMP}"
lastTimestamp: "${TIMESTAMP}"
count: 1
source:
  component: cluster
  host: $(hostname)
EOF
