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

# Config file
nodesConfigFile=${1:-"nodes-config.yaml"}
if [ ! -f "$nodesConfigFile" ]; then
  echo "Cannot find the config file $nodesConfigFile"
  exit 1
fi

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

# Create the namespace to run the node-joiner, along with the required roles and bindings.
staticResources=$(cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: ${namespace}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-joiner
  namespace: ${namespace}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-joiner
rules:
- apiGroups:
  - config.openshift.io
  resources:
  - clusterversions
  - proxies
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - secrets
  - configmaps
  - nodes
  verbs:
  - get
  - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: node-joiner
subjects:
- kind: ServiceAccount
  name: node-joiner
  namespace: ${namespace}
roleRef:
  kind: ClusterRole
  name: node-joiner
  apiGroup: rbac.authorization.k8s.io
EOF
)
echo "$staticResources" | oc apply -f -

# Generate a configMap to store the user configuration
oc create configmap nodes-config --from-file=nodes-config.yaml="${nodesConfigFile}" -n "${namespace}" -o yaml --dry-run=client | oc apply -f -

# Run the node-joiner pod to generate the ISO
nodeJoinerPod=$(cat <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: node-joiner
  namespace: ${namespace}
  annotations:
    openshift.io/scc: anyuid
  labels:
    app: node-joiner    
spec:
  restartPolicy: Never
  serviceAccountName: node-joiner
  securityContext:
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: node-joiner
    imagePullPolicy: IfNotPresent
    image: $nodeJoinerPullspec
    volumeMounts:
    - name: nodes-config
      mountPath: /config
    - name: assets
      mountPath: /assets
    command: ["/bin/sh", "-c", "cp /config/nodes-config.yaml /assets; HOME=/assets node-joiner add-nodes --dir=/assets --log-level=debug; sleep 600"]    
  volumes:
  - name: nodes-config
    configMap: 
      name: nodes-config
      namespace: ${namespace}
  - name: assets
    emptyDir: 
      sizeLimit: "4Gi"
EOF
)
echo "$nodeJoinerPod" | oc apply -f -

while true; do 
  if oc exec node-joiner -n "${namespace}" -- test -e /assets/exit_code >/dev/null 2>&1; then
    break
  else 
    echo "Waiting for node-joiner pod to complete..."
    sleep 10s
  fi
done

res=$(oc exec node-joiner -n "${namespace}" -- cat /assets/exit_code)
if [ "$res" = 0 ]; then
  echo "node-joiner successfully completed, extracting ISO image..."
  oc cp -n "${namespace}" node-joiner:/assets/node.x86_64.iso node.x86_64.iso
else
  oc logs node-joiner -n "${namespace}"
  echo "node-joiner failed"
fi

echo "Cleaning up"
oc delete namespace "${namespace}" --grace-period=0 >/dev/null 2>&1 &