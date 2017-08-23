#!/bin/sh
set -e

if [ "$#" -ne "3" ]; then
    echo "Usage: $0 kubeconfig assets_path experimental"
    exit 1
fi

KUBECONFIG="$1"
ASSETS_PATH="$2"
EXPERIMENTAL="$3"

# Setup API Authentication
KUBECTL="/kubectl --kubeconfig=$KUBECONFIG"

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

# chdir into the assets path directory
cd "$ASSETS_PATH/tectonic"

# Wait for Kubernetes to be in a proper state
set +e
i=0
echo "Waiting for Kubernetes API..."
until $KUBECTL cluster-info; do
  i=$((i+1))
  echo "Cluster not available yet, waiting for 5 seconds ($i)"
  sleep 5
done
set -e

# wait for Kubernetes pods
wait_for_pods kube-system

# Creating resources
echo "Creating Tectonic Namespace"
kubectl create -f namespace.yaml

echo "Creating Initial Roles"
kubectl delete -f rbac/role-admin.yaml

kubectl create -f rbac/role-admin.yaml
kubectl create -f rbac/role-user.yaml
kubectl create -f rbac/binding-admin.yaml
kubectl create -f rbac/binding-discovery.yaml

echo "Creating Tectonic ConfigMaps"
kubectl create -f config.yaml

echo "Creating Tectonic Secrets"
kubectl create -f secrets/pull.json
kubectl create -f secrets/license.json
kubectl create -f secrets/ingress-tls.yaml
kubectl create -f secrets/ca-cert.yaml
kubectl create -f secrets/identity-grpc-client.yaml
kubectl create -f secrets/identity-grpc-server.yaml

echo "Creating Ingress"
kubectl create -f ingress/default-backend/configmap.yaml
kubectl create -f ingress/default-backend/service.yaml
kubectl create -f ingress/default-backend/deployment.yaml
kubectl create -f ingress/ingress.yaml

# shellcheck disable=SC2154
if [ "${ingress_kind}" = "HostPort" ]; then
  kubectl create -f ingress/hostport/service.yaml
  kubectl create -f ingress/hostport/daemonset.yaml
elif [ "${ingress_kind}" = "NodePort" ]; then
  kubectl create -f ingress/nodeport/service.yaml
  kubectl create -f ingress/nodeport/deployment.yaml
else
  echo "Unrecognized Ingress Kind: ${ingress_kind}"
fi

echo "Creating Tectonic Identity"
kubectl create -f identity/configmap.yaml
kubectl create -f identity/services.yaml
kubectl create -f identity/deployment.yaml

echo "Creating Tectonic Console"
kubectl create -f console/service.yaml
kubectl create -f console/deployment.yaml

echo "Creating Tectonic Monitoring"
kubectl create -f monitoring/prometheus-operator-cluster-role-binding.yaml
kubectl create -f monitoring/prometheus-operator-cluster-role.yaml
kubectl create -f monitoring/prometheus-operator-service-account.yaml
kubectl create -f monitoring/prometheus-operator-svc.yaml
kubectl create -f monitoring/prometheus-operator.yaml

wait_for_tpr tectonic-system prometheus.monitoring.coreos.com
wait_for_tpr tectonic-system alertmanager.monitoring.coreos.com
wait_for_tpr tectonic-system service-monitor.monitoring.coreos.com

kubectl create -f monitoring/alertmanager-config.yaml
kubectl create -f monitoring/alertmanager-service.yaml
kubectl create -f monitoring/alertmanager.yaml
kubectl create -f monitoring/kube-controller-manager-svc.yaml
kubectl create -f monitoring/kube-scheduler-svc.yaml
kubectl create -f monitoring/kube-state-metrics-cluster-role-binding.yaml
kubectl create -f monitoring/kube-state-metrics-cluster-role.yaml
kubectl create -f monitoring/kube-state-metrics-deployment.yaml
kubectl create -f monitoring/kube-state-metrics-service-account.yaml
kubectl create -f monitoring/kube-state-metrics-service.yaml
kubectl create -f monitoring/node-exporter-ds.yaml
kubectl create -f monitoring/node-exporter-svc.yaml
kubectl create -f monitoring/prometheus-k8s-cluster-role-binding.yaml
kubectl create -f monitoring/prometheus-k8s-cluster-role.yaml
kubectl create -f monitoring/prometheus-k8s-rules.yaml
kubectl create -f monitoring/prometheus-k8s-service-account.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-alertmanager.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-apiserver.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-kube-controller-manager.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-kube-scheduler.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-kube-state-metrics.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-kubelet.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-node-exporter.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-prometheus.yaml
kubectl create -f monitoring/prometheus-k8s-service-monitor-prometheus-operator.yaml
kubectl create -f monitoring/prometheus-k8s.yaml
kubectl create -f monitoring/prometheus-svc.yaml

kubectl create -f monitoring/tectonic-monitoring-auth-secret.yaml
kubectl create -f monitoring/tectonic-monitoring-auth-alertmanager-deployment.yaml
kubectl create -f monitoring/tectonic-monitoring-auth-alertmanager-svc.yaml
kubectl create -f monitoring/tectonic-monitoring-auth-prometheus-deployment.yaml
kubectl create -f monitoring/tectonic-monitoring-auth-prometheus-svc.yaml
kubectl create -f monitoring/tectonic-monitoring-ingress.yaml

echo "Creating Etcd Operator"
# Operator in the tectonic-system namespace used for etcd as a service
kubectl create -f etcd/etcd-operator.yaml

echo "Creating Heapster / Stats Emitter"
kubectl create -f heapster/service.yaml
kubectl create -f heapster/deployment.yaml
kubectl create -f stats-emitter.yaml

echo "Creating Operators"
kubectl create -f updater/tectonic-channel-operator-kind.yaml
kubectl create -f updater/app-version-kind.yaml
kubectl create -f updater/migration-status-kind.yaml
kubectl create -f updater/node-agent.yaml
kubectl create -f updater/tectonic-monitoring-config.yaml

wait_for_crd tectonic-system channeloperatorconfigs.tco.coreos.com
kubectl create -f updater/tectonic-channel-operator-config.yaml

kubectl create -f updater/operators/kube-version-operator.yaml
kubectl create -f updater/operators/tectonic-channel-operator.yaml
kubectl create -f updater/operators/tectonic-prometheus-operator.yaml
kubectl create -f updater/operators/tectonic-cluo-operator.yaml

wait_for_crd tectonic-system appversions.tco.coreos.com
kubectl create -f updater/app_versions/app-version-tectonic-cluster.yaml
kubectl create -f updater/app_versions/app-version-kubernetes.yaml
kubectl create -f updater/app_versions/app-version-tectonic-monitoring.yaml
kubectl create -f updater/app_versions/app-version-tectonic-cluo.yaml

if [ "$EXPERIMENTAL" = "true" ]; then
  echo "Creating Experimental resources"
  kubectl apply -f updater/cluster-config.yaml
  kubectl create -f updater/app_versions/app-version-tectonic-etcd.yaml
  kubectl create -f updater/operators/tectonic-etcd-operator.yaml
fi

# wait for Tectonic pods
wait_for_pods tectonic-system

echo "Tectonic installation is done"
exit 0
