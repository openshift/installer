#!/bin/bash
set -e

if [ "$#" -ne "2" ]; then
    echo "Usage: $0 kubeconfig assets_path"
    exit 1
fi

KUBECONFIG=$1
ASSETS_PATH=$2

# Setup API Authentication
K8S_API=$(grep "server" $KUBECONFIG | cut -f 2- -d ":" | tr -d " ")
K8S_API_CA=$(mktemp); grep "certificate-authority-data" $KUBECONFIG | cut -f 2- -d ":" | tr -d " " | base64 -d > $K8S_API_CA
K8S_API_CERT=$(mktemp); grep "client-certificate-data" $KUBECONFIG | cut -f 2- -d ":" | tr -d " " | base64 -d > $K8S_API_CERT
K8S_API_KEY=$(mktemp); grep "client-key-data" $KUBECONFIG | cut -f 2- -d ":" | tr -d " " | base64 -d > $K8S_API_KEY
CURL="curl -sNL --cacert $K8S_API_CA --cert $K8S_API_CERT --key $K8S_API_KEY"
trap "rm -f $K8S_API_CA $K8S_API_CERT $K8S_API_KEY" EXIT

# Setup helper functions
function wait_for_tpr() {
  echo "Waiting for third-party resource definitions..."
  until $CURL -f "$K8S_API/$1" &> /dev/null; do
    sleep 5
  done
}

function create_resource() {
  STATUS=$($CURL -o /dev/null --write-out '%{http_code}\n' -H "Content-Type: application/$1" -d"$(cat $ASSETS_PATH/$2)" "$K8S_API/$3")
  if [ "$STATUS" != "200" ] && [ "$STATUS" != "201" ] && [ "$STATUS" != "409" ]; then
    echo -e "Failed to create $2 (got $STATUS): " >&2
    $CURL -H "Content-Type: application/$1" -d"$(cat $ASSETS_PATH/$2)" "$K8S_API/$3" >&2
    exit 1
  fi
}

function delete_resource() {
  $CURL -H "Content-Type: application/$1" -XDELETE "$K8S_API/$2" &> /dev/null
}

# Wait for Kubernetes to be in a proper state
echo "Waiting for Kubernetes API..."
until $CURL -f "$K8S_API/version" &> /dev/null; do
  sleep 5
done

echo "Waiting for Kubernetes components..."
while $CURL "$K8S_API/api/v1/namespaces/kube-system/pods" 2>/dev/null | grep Pending > /dev/null; do
  sleep 5
done
sleep 10

# Creating resources
echo "Creating Tectonic Namespace"
create_resource yaml namespace.yaml api/v1/namespaces

echo "Creating Initial Roles"
delete_resource yaml apis/rbac.authorization.k8s.io/v1alpha1/clusterroles/admin
create_resource yaml rbac/role-admin.yaml        apis/rbac.authorization.k8s.io/v1alpha1/clusterroles
create_resource yaml rbac/role-user.yaml         apis/rbac.authorization.k8s.io/v1alpha1/clusterroles
create_resource yaml rbac/binding-admin.yaml     apis/rbac.authorization.k8s.io/v1alpha1/clusterrolebindings
create_resource yaml rbac/binding-discovery.yaml apis/rbac.authorization.k8s.io/v1alpha1/clusterrolebindings

echo "Creating Tectonic ConfigMaps"
create_resource yaml config.yaml api/v1/namespaces/tectonic-system/configmaps

echo "Creating Tectonic Secrets"
create_resource json secrets/pull.json                 api/v1/namespaces/tectonic-system/secrets
create_resource json secrets/license.json              api/v1/namespaces/tectonic-system/secrets
create_resource yaml secrets/ingress-tls.yaml          api/v1/namespaces/tectonic-system/secrets
create_resource yaml secrets/ca-cert.yaml              api/v1/namespaces/tectonic-system/secrets
create_resource yaml secrets/identity-grpc-client.yaml api/v1/namespaces/tectonic-system/secrets
create_resource yaml secrets/identity-grpc-server.yaml api/v1/namespaces/tectonic-system/secrets

echo "Creating Tectonic Identity"
create_resource yaml identity/configmap.yaml  api/v1/namespaces/tectonic-system/configmaps
create_resource yaml identity/services.yaml   api/v1/namespaces/tectonic-system/services
create_resource yaml identity/deployment.yaml apis/extensions/v1beta1/namespaces/tectonic-system/deployments

echo "Creating Tectonic Console"
create_resource yaml console/service.yaml    api/v1/namespaces/tectonic-system/services
create_resource yaml console/deployment.yaml apis/extensions/v1beta1/namespaces/tectonic-system/deployments

echo "Creating Tectonic Monitoring"
create_resource yaml monitoring/prometheus-operator-service-account.yaml      api/v1/namespaces/tectonic-system/serviceaccounts
create_resource yaml monitoring/prometheus-operator-cluster-role.yaml         apis/rbac.authorization.k8s.io/v1alpha1/clusterroles
create_resource yaml monitoring/prometheus-operator-cluster-role-binding.yaml apis/rbac.authorization.k8s.io/v1alpha1/clusterrolebindings
create_resource yaml monitoring/prometheus-k8s-service-account.yaml           api/v1/namespaces/tectonic-system/serviceaccounts
create_resource yaml monitoring/prometheus-k8s-cluster-role.yaml              apis/rbac.authorization.k8s.io/v1alpha1/clusterroles
create_resource yaml monitoring/prometheus-k8s-cluster-role-binding.yaml      apis/rbac.authorization.k8s.io/v1alpha1/clusterrolebindings
create_resource yaml monitoring/prometheus-k8s-config.yaml                    api/v1/namespaces/tectonic-system/configmaps
create_resource yaml monitoring/prometheus-k8s-rules.yaml                     api/v1/namespaces/tectonic-system/configmaps
create_resource yaml monitoring/prometheus-svc.yaml                           api/v1/namespaces/tectonic-system/services
create_resource yaml monitoring/node-exporter-svc.yaml                        api/v1/namespaces/tectonic-system/services
create_resource yaml monitoring/node-exporter-ds.yaml                         apis/extensions/v1beta1/namespaces/tectonic-system/daemonsets
create_resource yaml monitoring/prometheus-operator.yaml                      apis/extensions/v1beta1/namespaces/tectonic-system/deployments
wait_for_tpr apis/monitoring.coreos.com/v1alpha1/prometheuses
create_resource json monitoring/prometheus-k8s.json                           apis/monitoring.coreos.com/v1alpha1/namespaces/tectonic-system/prometheuses

echo "Creating Ingress"
create_resource yaml ingress/default-backend/configmap.yaml  api/v1/namespaces/tectonic-system/configmaps
create_resource yaml ingress/default-backend/service.yaml    api/v1/namespaces/tectonic-system/services
create_resource yaml ingress/default-backend/deployment.yaml apis/extensions/v1beta1/namespaces/tectonic-system/deployments
create_resource yaml ingress/ingress.yaml                    apis/extensions/v1beta1/namespaces/tectonic-system/ingresses

if [ "${ingress_kind}" = "HostPort" ]; then
  create_resource yaml ingress/hostport.yaml apis/extensions/v1beta1/namespaces/tectonic-system/daemonsets
elif [ "${ingress_kind}" = "NodePort" ]; then
  create_resource yaml ingress/nodeport/configmap.yaml  api/v1/namespaces/tectonic-system/configmaps
  create_resource yaml ingress/nodeport/service.yaml    api/v1/namespaces/tectonic-system/services
  create_resource yaml ingress/nodeport/deployment.yaml apis/extensions/v1beta1/namespaces/tectonic-system/deployments
else
  echo "Unrecognized Ingress Kind: ${ingress_kind}"
fi

echo "Creating Heapster / Stats Emitter"
create_resource yaml heapster/service.yaml    api/v1/namespaces/kube-system/services
create_resource yaml heapster/deployment.yaml apis/extensions/v1beta1/namespaces/kube-system/deployments
create_resource yaml stats-emitter.yaml       apis/extensions/v1beta1/namespaces/tectonic-system/deployments

echo "Creating Tectonic Updater"
create_resource yaml updater/tectonic-channel-operator-kind.yaml        apis/extensions/v1beta1/thirdpartyresources
create_resource yaml updater/app-version-kind.yaml                      apis/extensions/v1beta1/thirdpartyresources
create_resource yaml updater/migration-status-kind.yaml                 apis/extensions/v1beta1/thirdpartyresources
create_resource yaml updater/node-agent.yaml                            apis/extensions/v1beta1/namespaces/tectonic-system/daemonsets
create_resource yaml updater/kube-version-operator.yaml                 apis/extensions/v1beta1/namespaces/tectonic-system/deployments
create_resource yaml updater/tectonic-channel-operator.yaml             apis/extensions/v1beta1/namespaces/tectonic-system/deployments
wait_for_tpr apis/coreos.com/v1/channeloperatorconfigs
create_resource json updater/tectonic-channel-operator-config.json      apis/coreos.com/v1/namespaces/tectonic-system/channeloperatorconfigs
wait_for_tpr apis/coreos.com/v1/appversions
create_resource json updater/app-version-tectonic-cluster.json          apis/coreos.com/v1/namespaces/tectonic-system/appversions
create_resource json updater/app-version-kubernetes.json                apis/coreos.com/v1/namespaces/tectonic-system/appversions

echo "Tectonic installation is done"
exit 0