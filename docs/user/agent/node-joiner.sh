#!/bin/bash

if [ $# -lt 1 ]; then
    echo "./node-joiner.sh <pull secret path>"
    echo "Usage example:"
    echo "$ ./node-joiner.sh ~/config/my-pull-secret"

    exit 1
fi
pullSecret=$1

# Extract the installer image pullspec and release version.
releaseImage=$(oc get clusterversion version -o=jsonpath='{.status.history[?(@.state == "Completed")].image}')
nodeJoinerPullspec=$(oc adm release info -a "$pullSecret" --image-for=installer "$releaseImage")

# Create the namespace to run the node-joiner, along with the required roles and bindings.
staticResources=$(cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: openshift-node-joiner
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: node-joiner
  namespace: openshift-node-joiner
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
  namespace: openshift-node-joiner
roleRef:
  kind: ClusterRole
  name: node-joiner
  apiGroup: rbac.authorization.k8s.io
EOF
)
echo "$staticResources" | oc apply -f -

# Generate a configMap to store the user configuration
oc create configmap nodes-config --from-file=nodes-config.yaml -n openshift-node-joiner -o yaml --dry-run=client | oc apply -f -

# Runt the node-joiner pod to generate the ISO
nodeJoinerPod=$(cat <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: node-joiner
  namespace: openshift-node-joiner
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
    command: ["/bin/sh", "-c", "cp /config/nodes-config.yaml /assets; HOME=/assets node-joiner add-nodes --dir=/assets --log-level=debug; echo \$? > /assets/completed; sleep 600"]    
  volumes:
  - name: nodes-config
    configMap: 
      name: nodes-config
      namespace: openshift-node-joiner
  - name: assets
    emptyDir: 
      sizeLimit: "4Gi"
EOF
)
echo "$nodeJoinerPod" | oc apply -f -

# Wait until the node-joiner was completed.
while true; do 
  if oc exec node-joiner -n openshift-node-joiner -- test -e /assets/completed >/dev/null 2>&1; then
    break
  else 
    echo "Waiting for node-joiner pod to complete..."
    sleep 10s
  fi
done

# In case of success, let's extract the ISO, otherwise the logs are shown for troubleshooting the error.
completed=$(oc exec node-joiner -n openshift-node-joiner -- cat /assets/completed)
if [ "$completed" = 0 ]; then
  echo "node-joiner successfully completed, extracting ISO image..."
  oc cp -n openshift-node-joiner node-joiner:/assets/agent-addnodes.x86_64.iso agent-addnodes.x86_64.iso
else
  oc logs node-joiner -n openshift-node-joiner 
  echo "node-joiner failed"
fi

# Remove all the resources previously created.
echo "Cleaning up"
oc delete namespace openshift-node-joiner --grace-period=0 >/dev/null 2>&1 &