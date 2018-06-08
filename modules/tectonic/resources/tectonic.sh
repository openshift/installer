#!/bin/sh
set -e

KUBECONFIG="$1"
ASSETS_PATH="$2"

# Setup API Authentication
KUBECTL="oc --config=$KUBECONFIG"

# Setup helper functions

kubectl() {
  i=0

  echo "Executing kubectl" "$@"
  while true; do
    i=$((i+1))
    [ $i -eq 100 ] && echo "kubectl failed, giving up" && exit 1

    set +e
    out=$($KUBECTL "$@" 2>&1)
    status=$?
    set -e

    if echo "$out" | grep -q "AlreadyExists"; then
      echo "$out, skipping"
      return
    fi

    echo "$out"
    if [ "$status" -eq 0 ]; then
      return
    fi

    echo "kubectl failed, retrying in 5 seconds"
    sleep 5
  done
}

wait_for_crd() {
  set +e
  i=0

  echo "Waiting for CRD $2"
  until $KUBECTL -n "$1" get customresourcedefinition "$2"; do
    i=$((i+1))
    echo "CRD $2 not available yet, retrying in 5 seconds ($i)"
    sleep 5
  done
  set -e
}

wait_for_tpr() {
  set +e
  i=0

  echo "Waiting for TPR $2"
  until $KUBECTL -n "$1" get thirdpartyresources "$2"; do
    i=$((i+1))
    echo "TPR $2 not available yet, retrying in 5 seconds ($i)"
    sleep 5
  done
  set -e
}

wait_for_pods() {
  set +e
  echo "Waiting for pods in namespace $1"
  while true; do
  
    out=$($KUBECTL -n "$1" get po -o custom-columns=STATUS:.status.phase,NAME:.metadata.name)
    status=$?
    echo "$out"
  
    if [ "$status" -ne "0" ]; then
      echo "kubectl command failed, retrying in 5 seconds"
      sleep 5
      continue
    fi
  
    # make sure kubectl does not return "no resources found"
    if [ "$(echo "$out" | tail -n +2 | grep -c '^')" -eq 0 ]; then
      echo "no resources were found, retrying in 5 seconds"
      sleep 5
      continue
    fi
  
    stat=$(echo "$out"| tail -n +2 | grep -v '^Running')
    if [ -z "$stat" ]; then
      return
    fi
  
    echo "Pods not available yet, waiting for 5 seconds"
    sleep 5
  done
  set -e
}

asset_cleanup() {
  echo "Cleaning up installation assets"

  # shellcheck disable=SC2034
  for d in "manifests" "auth" "tectonic" "tls"; do
      rm -rf "$${ASSETS_PATH:?}/$${d:?}/"*
  done

  # shellcheck disable=SC2034
  for f in "bootkube.sh" "tectonic.sh" "tectonic-wrapper.sh"; do
      rm -f "$${ASSETS_PATH:?}/$${f:?}"
  done
}

# chdir into the assets path directory
cd "$ASSETS_PATH/tectonic"

# Wait for Kubernetes to be in a proper state
set +e
i=0
echo "Waiting for Kubernetes API..."
until $KUBECTL status; do
  i=$((i+1))
  echo "Cluster not available yet, waiting for 5 seconds ($i)"
  sleep 5
done
set -e

# wait for Kubernetes pods
wait_for_pods kube-system

echo "Creating Initial Roles"
kubectl delete -f rbac/role-admin.yaml

kubectl create -f ingress/svc-account.yaml
kubectl create -f rbac/role-admin.yaml
kubectl create -f rbac/role-user.yaml
kubectl create -f rbac/binding-admin.yaml
kubectl create -f rbac/binding-discovery.yaml

echo "Creating Cluster Config For Tectonic"
kubectl create -f cluster-config.yaml
kubectl create -f ingress/cluster-config.yaml

echo "Creating Tectonic Secrets"
kubectl create -f secrets/pull.json
kubectl create -f secrets/license.json
kubectl create -f secrets/ingress-tls.yaml
kubectl create -f secrets/ca-cert.yaml
kubectl create -f ingress/pull.json

echo "Creating Operators"
kubectl create -f security/priviledged-scc-tectonic.yaml
kubectl create -f updater/tectonic-channel-operator-kind.yaml
kubectl create -f updater/app-version-kind.yaml
kubectl create -f updater/migration-status-kind.yaml

wait_for_crd tectonic-system channeloperatorconfigs.tco.coreos.com
kubectl create -f updater/tectonic-channel-operator-config.yaml

kubectl create -f updater/operators/kube-core-operator.yaml
kubectl create -f updater/operators/tectonic-channel-operator.yaml
kubectl create -f updater/operators/kube-addon-operator.yaml
kubectl create -f updater/operators/tectonic-alm-operator.yaml
kubectl create -f updater/operators/tectonic-utility-operator.yaml
kubectl create -f updater/operators/tectonic-ingress-controller-operator.yaml

wait_for_crd tectonic-system appversions.tco.coreos.com
kubectl create -f updater/app_versions/app-version-tectonic-cluster.yaml
kubectl create -f updater/app_versions/app-version-kube-core.yaml
kubectl create -f updater/app_versions/app-version-kube-addon.yaml
kubectl create -f updater/app_versions/app-version-tectonic-alm.yaml
kubectl create -f updater/app_versions/app-version-tectonic-utility.yaml
kubectl create -f updater/app_versions/app-version-tectonic-ingress.yaml

# wait for Tectonic pods
wait_for_pods tectonic-system
asset_cleanup

echo "Tectonic installation is done"
exit 0
