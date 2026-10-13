#!/bin/bash
#
# Deprecated
#
# For OpenShift 4.17 and later, the node-joiner.sh and node-joiner-monitor.sh
# scripts cannot be used. The node-joiner scripts were created for
# OpenShift 4.16 as a development preview. They have been deprecated and
# replaced by the "oc adm node-image" command.
# See https://docs.openshift.com/container-platform/4.17/nodes/nodes/nodes-nodes-adding-node-iso.html
# for more details.
#

set -eu

if [ $# -eq 0 ]; then
    echo "At least one IP address must be provided"
    exit 1
fi

ipAddresses="$*"

# Setup a cleanup function to ensure to remove the temporary
# file when the script will be completed.
cleanup() {
  if [ -f "$pullSecretFile" ]; then
    echo "Removing temporary file $pullSecretFile"
    rm "$pullSecretFile"
  fi
}
trap cleanup EXIT TERM

# Retrieve the pullsecret and store it in a temporary file. 
pullSecretFile=$(mktemp -p "/tmp" -t "nodejoiner-XXXXXXXXXX")
oc get secret -n openshift-config pull-secret -o jsonpath='{.data.\.dockerconfigjson}' | base64 -d > "$pullSecretFile"

# Extract the baremetal-installer image pullspec from the current cluster.
nodeJoinerPullspec=$(oc adm release info --image-for=baremetal-installer --registry-config="$pullSecretFile")

# Use the same random temp file suffix for the namespace.
namespace=$(echo "openshift-node-joiner-${pullSecretFile#/tmp/nodejoiner-}" | tr '[:upper:]' '[:lower:]')

# Create the namespace to run the node-joiner-monitor, along with the required roles and bindings.
staticResources=$(cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: ${namespace}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-joiner-monitor
  namespace: ${namespace}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-joiner-monitor
rules:
- apiGroups:
  - certificates.k8s.io
  resources:
  - certificatesigningrequests
  verbs:
  - get
  - list
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  verbs:
  - get
  - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: node-joiner-monitor
subjects:
- kind: ServiceAccount
  name: node-joiner-monitor
  namespace: ${namespace}
roleRef:
  kind: ClusterRole
  name: node-joiner-monitor
  apiGroup: rbac.authorization.k8s.io
EOF
)
echo "$staticResources" | oc apply -f -

# Run the node-joiner-monitor to monitor node joining cluster
nodeJoinerPod=$(cat <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: node-joiner-monitor
  namespace: ${namespace}
  annotations:
    openshift.io/scc: anyuid
  labels:
    app: node-joiner-monitor    
spec:
  restartPolicy: Never
  serviceAccountName: node-joiner-monitor
  securityContext:
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: node-joiner-monitor
    imagePullPolicy: IfNotPresent
    image: $nodeJoinerPullspec
    command: ["/bin/sh", "-c", "node-joiner monitor-add-nodes $ipAddresses --dir /tmp --log-level=info; sleep 5"]
EOF
)
echo "$nodeJoinerPod" | oc apply -f -

oc project "${namespace}"

oc wait --for=condition=Ready=true --timeout=300s pod/node-joiner-monitor

oc logs -f -n "${namespace}" node-joiner-monitor
 
echo "Cleaning up"
oc delete namespace "${namespace}" --grace-period=0 >/dev/null 2>&1 &