package ignition

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/content"
)

// TestBootstrapGenerate tests generating the bootstrap asset.
func TestBootstrapGenerate(t *testing.T) {
	installConfig := `
metadata:
  name: test-cluster
admin:
  email: test-admin-email
  sshKey: test-admin-ssh-key
baseDomain: test-domain
networking:
  ServiceCIDR: 10.0.1.0/24
platform:
  aws:
    region: us-east
machines:
- name: master
  replicas: 3
`
	installConfigAsset := &testAsset{"install-config"}
	rootCAAsset := &testAsset{"rootCA"}
	etcdCAAsset := &testAsset{"etcdCA"}
	ingressCertKeyAsset := &testAsset{"ingress-ca"}
	kubeCAAsset := &testAsset{"kubeCA"}
	aggregatorCAAsset := &testAsset{"aggregator-ca"}
	serviceServingCAAsset := &testAsset{"service-serving-ca"}
	clusterAPIServerCertKeyAsset := &testAsset{"cluster-apiserver-ca"}
	etcdClientCertKeyAsset := &testAsset{"etcd-client-ca"}
	apiServerCertKeyAsset := &testAsset{"apiserver-ca"}
	openshiftAPIServerCertKeyAsset := &testAsset{"openshift-apiserver-ca"}
	apiServerProxyCertKeyAsset := &testAsset{"apiserver-proxy-ca"}
	adminCertKeyAsset := &testAsset{"admin-ca"}
	kubeletCertKeyAsset := &testAsset{"kubelet-ca"}
	mcsCertKeyAsset := &testAsset{"mcs-ca"}
	serviceAccountKeyPairAsset := &testAsset{"service-account-ca"}
	kubeconfigAsset := &testAsset{"kubeconfig"}
	kubeconfigKubeletAsset := &testAsset{"kubeconfig-kubelet"}
	bootstrap := &bootstrap{
		directory:                 "test-directory",
		installConfig:             installConfigAsset,
		rootCA:                    rootCAAsset,
		etcdCA:                    etcdCAAsset,
		ingressCertKey:            ingressCertKeyAsset,
		kubeCA:                    kubeCAAsset,
		aggregatorCA:              aggregatorCAAsset,
		serviceServingCA:          serviceServingCAAsset,
		clusterAPIServerCertKey:   clusterAPIServerCertKeyAsset,
		etcdClientCertKey:         etcdClientCertKeyAsset,
		apiServerCertKey:          apiServerCertKeyAsset,
		openshiftAPIServerCertKey: openshiftAPIServerCertKeyAsset,
		apiServerProxyCertKey:     apiServerProxyCertKeyAsset,
		adminCertKey:              adminCertKeyAsset,
		kubeletCertKey:            kubeletCertKeyAsset,
		mcsCertKey:                mcsCertKeyAsset,
		serviceAccountKeyPair:     serviceAccountKeyPairAsset,
		kubeconfig:                kubeconfigAsset,
		kubeconfigKubelet:         kubeconfigKubeletAsset,
	}
	dependencies := map[asset.Asset]*asset.State{
		installConfigAsset:             stateWithContentsData(installConfig),
		rootCAAsset:                    stateWithContentsData("test-rootCA-priv", "test-rootCA-pub"),
		etcdCAAsset:                    stateWithContentsData("test-etcdCA-priv", "test-etcdCA-pub"),
		ingressCertKeyAsset:            stateWithContentsData("test-ingress-ca-priv", "test-ingress-ca-pub"),
		kubeCAAsset:                    stateWithContentsData("test-kubeCA-priv", "test-kubeCA-pub"),
		aggregatorCAAsset:              stateWithContentsData("test-aggregator-ca-priv", "test-aggregator-ca-pub"),
		serviceServingCAAsset:          stateWithContentsData("test-service-serving-ca-priv", "test-service-serving-ca-pub"),
		clusterAPIServerCertKeyAsset:   stateWithContentsData("test-cluster-apiserver-cert-priv", "test-cluster-apiserver-cert-pub"),
		etcdClientCertKeyAsset:         stateWithContentsData("test-etcd-client-cert-priv", "test-etcd-client-cert-pub"),
		apiServerCertKeyAsset:          stateWithContentsData("test-apiserver-cert-priv", "test-apiserver-cert-pub"),
		openshiftAPIServerCertKeyAsset: stateWithContentsData("test-openshift-apiserver-cert-priv", "test-openshift-apiserver-cert-pub"),
		apiServerProxyCertKeyAsset:     stateWithContentsData("test-apiserver-proxy-cert-priv", "test-apiserver-proxy-cert-pub"),
		adminCertKeyAsset:              stateWithContentsData("test-admin-cert-priv", "test-admin-cert-pub"),
		kubeletCertKeyAsset:            stateWithContentsData("test-kubelet-cert-priv", "test-kubelet-cert-pub"),
		mcsCertKeyAsset:                stateWithContentsData("test-mcs-cert-priv", "test-mcs-cert-pub"),
		serviceAccountKeyPairAsset:     stateWithContentsData("test-service-account-cert-priv", "test-service-account-cert-pub"),
		kubeconfigAsset:                stateWithContentsData("test-kubeconfig"),
		kubeconfigKubeletAsset:         stateWithContentsData("test-kubeconfig-kubelet"),
	}
	bootstrapState, err := bootstrap.Generate(dependencies)
	assert.NoError(t, err, "unexpected error generating bootstrap asset")
	assert.Equal(t, 1, len(bootstrapState.Contents), "unexpected number of contents in bootstrap state")
	assert.Equal(t, "test-directory/bootstrap.ign", bootstrapState.Contents[0].Name, "unexpected name for bootstrap ignition config")

	assertFilesInIgnitionConfig(
		t,
		bootstrapState.Contents[0].Data,
		fileAssertion{
			path: "/etc/ssl/etcd/ca.crt",
			data: "test-etcdCA-pub",
		},
		fileAssertion{
			path: "/etc/ssl/etcd/root-ca.crt",
			data: "test-rootCA-pub",
		},
		fileAssertion{
			path: "/etc/ssl/certs/root_ca.pem",
			data: "test-rootCA-priv",
		},
		fileAssertion{
			path: "/etc/ssl/certs/ingress_ca.pem",
			data: "test-ingress-ca-priv",
		},
		fileAssertion{
			path: "/etc/ssl/certs/etcd_ca.pem",
			data: "test-etcdCA-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/auth/kubeconfig",
			data: "test-kubeconfig",
		},
		fileAssertion{
			path: "/opt/tectonic/auth/kubeconfig-kubelet",
			data: "test-kubeconfig-kubelet",
		},
		fileAssertion{
			path: "/opt/tectonic/bootkube.sh",
			data: expectedBootkubeSh,
		},
		fileAssertion{
			path: "/opt/tectonic/tectonic.sh",
			data: expectedTectonicSh,
		},
		fileAssertion{
			path: "/opt/tectonic/tls/root-ca.crt",
			data: "test-rootCA-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/kube-ca.key",
			data: "test-kubeCA-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/kube-ca.crt",
			data: "test-kubeCA-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/aggregator-ca.key",
			data: "test-aggregator-ca-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/aggregator-ca.crt",
			data: "test-aggregator-ca-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/service-serving-ca.key",
			data: "test-service-serving-ca-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/service-serving-ca.crt",
			data: "test-service-serving-ca-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/etcd-client-ca.key",
			data: "test-etcdCA-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/etcd-client-ca.crt",
			data: "test-etcdCA-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/cluster-apiserver-ca.key",
			data: "test-cluster-apiserver-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/cluster-apiserver-ca.crt",
			data: "test-cluster-apiserver-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/etcd-client.key",
			data: "test-etcd-client-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/etcd-client.crt",
			data: "test-etcd-client-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/apiserver.key",
			data: "test-apiserver-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/apiserver.crt",
			data: "test-apiserver-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/openshift-apiserver.key",
			data: "test-openshift-apiserver-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/openshift-apiserver.crt",
			data: "test-openshift-apiserver-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/apiserver-proxy.key",
			data: "test-apiserver-proxy-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/apiserver-proxy.crt",
			data: "test-apiserver-proxy-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/admin.key",
			data: "test-admin-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/admin.crt",
			data: "test-admin-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/kubelet.key",
			data: "test-kubelet-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/kubelet.crt",
			data: "test-kubelet-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/mcs.key",
			data: "test-mcs-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/mcs.crt",
			data: "test-mcs-cert-pub",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/service-account.key",
			data: "test-service-account-cert-priv",
		},
		fileAssertion{
			path: "/opt/tectonic/tls/service-account.crt",
			data: "test-service-account-cert-pub",
		},
	)

	assertSystemdUnitsInIgnitionConfig(
		t,
		bootstrapState.Contents[0].Data,
		systemdUnitAssertion{
			name:     "bootkube.service",
			contents: content.BootkubeSystemdContents,
		},
		systemdUnitAssertion{
			name:     "tectonic.service",
			contents: content.TectonicSystemdContents,
		},
		systemdUnitAssertion{
			name:     "kubelet.service",
			contents: expectedKubeletService,
		},
	)

	assertUsersInIgnitionConfig(
		t,
		bootstrapState.Contents[0].Data,
		userAssertion{
			name:   "core",
			sshKey: "test-admin-ssh-key",
		},
	)
}

const (
	expectedBootkubeSh = `#!/usr/bin/env bash
set -e

mkdir --parents /etc/kubernetes/manifests/

if [ ! -d kco-bootstrap ]
then
	echo "Rendering Kubernetes core manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		--volume /etc/kubernetes:/etc/kubernetes:z \
		"quay.io/coreos/kube-core-renderer-dev:436b1b4395ae54d866edc88864c9b01797cebac1" \
		--config=/assets/kco-config.yaml \
		--output=/assets/kco-bootstrap

	cp --recursive kco-bootstrap/bootstrap-configs /etc/kubernetes/bootstrap-configs
	cp --recursive kco-bootstrap/bootstrap-manifests .
	cp --recursive kco-bootstrap/manifests .
fi

if [ ! -d mco-bootstrap ]
then
	echo "Rendering MCO manifests..."

	# shellcheck disable=SC2154
	podman run \
		--user 0 \
		--volume "$PWD:/assets:z" \
		"docker.io/openshift/origin-machine-config-operator:v4.0.0" \
		bootstrap \
			--etcd-ca=/assets/tls/etcd-client-ca.crt \
			--root-ca=/assets/tls/root-ca.crt \
			--config-file=/assets/manifests/cluster-config.yaml \
			--dest-dir=/assets/mco-bootstrap \
			--images-json-configmap=/assets/manifests/machine-config-operator-01-images-configmap.yaml

	# Bootstrap MachineConfigController uses /etc/mcc/bootstrap/manifests/ dir to
	# 1. read the controller config rendered by MachineConfigOperator
	# 2. read the default MachineConfigPools rendered by MachineConfigOperator
	# 3. read any additional MachineConfigs that are needed for the default MachineConfigPools.
	mkdir --parents /etc/mcc/bootstrap/
	cp --recursive mco-bootstrap/manifests /etc/mcc/bootstrap/manifests
	cp mco-bootstrap/machineconfigoperator-bootstrap-pod.yaml /etc/kubernetes/manifests/

	# /etc/ssl/mcs/tls.{crt, key} are locations for MachineConfigServer's tls assets.
	mkdir --parents /etc/ssl/mcs/
	cp tls/machine-config-server.crt /etc/ssl/mcs/tls.crt
	cp tls/machine-config-server.key /etc/ssl/mcs/tls.key
fi

# We originally wanted to run the etcd cert signer as
# a static pod, but kubelet could't remove static pod
# when API server is not up, so we have to run this as
# podman container.
# See https://github.com/kubernetes/kubernetes/issues/43292

echo "Starting etcd certificate signer..."

trap "podman rm --force etcd-signer" ERR

# shellcheck disable=SC2154
podman run \
	--name etcd-signer \
	--detach \
	--volume /opt/tectonic/tls:/opt/tectonic/tls:ro,z \
	--network host \
	"quay.io/coreos/kube-etcd-signer-server:678cc8e6841e2121ebfdb6e2db568fce290b67d6" \
	serve \
	--cacrt=/opt/tectonic/tls/etcd-client-ca.crt \
	--cakey=/opt/tectonic/tls/etcd-client-ca.key \
	--servcrt=/opt/tectonic/tls/apiserver.crt \
	--servkey=/opt/tectonic/tls/apiserver.key \
	--address=0.0.0.0:6443 \
	--csrdir=/tmp \
	--peercertdur=26280h \
	--servercertdur=26280h

echo "Waiting for etcd cluster..."

# Wait for the etcd cluster to come up.
set +e
# shellcheck disable=SC2154,SC2086
until podman run \
		--rm \
		--network host \
		--name etcdctl \
		--env ETCDCTL_API=3 \
		--volume /opt/tectonic/tls:/opt/tectonic/tls:ro,z \
		"quay.io/coreos/etcd:v3.2.14" \
		/usr/local/bin/etcdctl \
		--dial-timeout=10m \
		--cacert=/opt/tectonic/tls/etcd-client-ca.crt \
		--cert=/opt/tectonic/tls/etcd-client.crt \
		--key=/opt/tectonic/tls/etcd-client.key \
		--endpoints=https://test-cluster-etcd-0.test-domain:2379,https://test-cluster-etcd-1.test-domain:2379,https://test-cluster-etcd-2.test-domain:2379 \
		endpoint health
do
	echo "etcdctl failed. Retrying in 5 seconds..."
	sleep 5
done
set -e

echo "etcd cluster up. Killing etcd certificate signer..."

podman rm --force etcd-signer
rm --force /etc/kubernetes/manifests/machineconfigoperator-bootstrap-pod.yaml

echo "Starting bootkube..."

# shellcheck disable=SC2154
podman run \
	--rm \
	--volume "$PWD:/assets:z" \
	--volume /etc/kubernetes:/etc/kubernetes:z \
	--network=host \
	--entrypoint=/bootkube \
	"quay.io/coreos/bootkube:v0.10.0" \
	start --asset-dir=/assets`

	expectedTectonicSh = `#!/usr/bin/env bash
set -e

KUBECONFIG="$1"

kubectl() {
	echo "Executing kubectl $*" >&2
	while true
	do
		set +e
		out=$(oc --config="$KUBECONFIG" "$@" 2>&1)
		status=$?
		set -e

		if grep --quiet "AlreadyExists" <<< "$out"
		then
			echo "$out, skipping" >&2
			return
		fi

		echo "$out"
		if [ "$status" -eq 0 ]
		then
			return
		fi

		echo "kubectl $* failed. Retrying in 5 seconds..." >&2
		sleep 5
	done
}

wait_for_pods() {
	echo "Waiting for pods in namespace $1..."
	while true
	do
		out=$(kubectl --namespace "$1" get pods --output custom-columns=STATUS:.status.phase,NAME:.metadata.name --no-headers=true)
		echo "$out"

		# make sure kubectl returns at least one status
		if [ "$(wc --lines <<< "$out")" -eq 0 ]
		then
			echo "No pods were found. Waiting for 5 seconds..."
			sleep 5
			continue
		fi

		if ! grep --invert-match '^Running' <<< "$out"
		then
			return
		fi

		echo "Not all pods available yet. Waiting for 5 seconds..."
		sleep 5
	done
	set -e
}

# Wait for Kubernetes pods
wait_for_pods kube-system

echo "Creating initial roles..."
kubectl delete --filename rbac/role-admin.yaml

kubectl create --filename ingress/svc-account.yaml
kubectl create --filename rbac/role-admin.yaml
kubectl create --filename rbac/role-user.yaml
kubectl create --filename rbac/binding-admin.yaml
kubectl create --filename rbac/binding-discovery.yaml

echo "Creating cluster config for Tectonic..."
kubectl create --filename cluster-config.yaml
kubectl create --filename ingress/cluster-config.yaml

echo "Creating Tectonic secrets..."
kubectl create --filename secrets/pull.json
kubectl create --filename secrets/ingress-tls.yaml
kubectl create --filename secrets/ca-cert.yaml
kubectl create --filename ingress/pull.json

echo "Creating operators..."
kubectl create --filename security/priviledged-scc-tectonic.yaml
kubectl create --filename updater/tectonic-channel-operator-kind.yaml
kubectl create --filename updater/app-version-kind.yaml
kubectl create --filename updater/migration-status-kind.yaml

kubectl --namespace=tectonic-system get customresourcedefinition channeloperatorconfigs.tco.coreos.com
kubectl create --filename updater/tectonic-channel-operator-config.yaml

kubectl create --filename updater/operators/kube-core-operator.yaml
kubectl create --filename updater/operators/tectonic-channel-operator.yaml
kubectl create --filename updater/operators/kube-addon-operator.yaml
kubectl create --filename updater/operators/tectonic-alm-operator.yaml
kubectl create --filename updater/operators/tectonic-utility-operator.yaml
kubectl create --filename updater/operators/tectonic-ingress-controller-operator.yaml

kubectl --namespace=tectonic-system get customresourcedefinition appversions.tco.coreos.com
kubectl create --filename updater/app_versions/app-version-tectonic-cluster.yaml
kubectl create --filename updater/app_versions/app-version-kube-core.yaml
kubectl create --filename updater/app_versions/app-version-kube-addon.yaml
kubectl create --filename updater/app_versions/app-version-tectonic-alm.yaml
kubectl create --filename updater/app_versions/app-version-tectonic-utility.yaml
kubectl create --filename updater/app_versions/app-version-tectonic-ingress.yaml

# Wait for Tectonic pods
wait_for_pods tectonic-system

echo "Tectonic installation is done"`

	// TODO (staebler): --cluster-dns should be 10.0.1.10 instead of 0:a::ffff:a00:100.
	expectedKubeletService = `[Unit]
Description=Kubernetes Kubelet
Wants=rpc-statd.service

[Service]
ExecStartPre=/bin/mkdir --parents /etc/kubernetes/manifests
ExecStartPre=/usr/bin/bash -c "gawk '/certificate-authority-data/ {print $2}' /etc/kubernetes/kubeconfig | base64 --decode > /etc/kubernetes/ca.crt"

ExecStart=/usr/bin/hyperkube \
  kubelet \
    --bootstrap-kubeconfig=/etc/kubernetes/kubeconfig \
    --kubeconfig=/var/lib/kubelet/kubeconfig \
    --rotate-certificates \
    --cni-conf-dir=/etc/kubernetes/cni/net.d \
    --cni-bin-dir=/var/lib/cni/bin \
    --network-plugin=cni \
    --lock-file=/var/run/lock/kubelet.lock \
    --exit-on-lock-contention \
    --pod-manifest-path=/etc/kubernetes/manifests \
    --allow-privileged \
    --node-labels=node-role.kubernetes.io/bootstrap \
    --register-with-taints=node-role.kubernetes.io/bootstrap=:NoSchedule \
    --minimum-container-ttl-duration=6m0s \
    --cluster-dns=0:a::ffff:a00:100 \
    --cluster-domain=cluster.local \
    --client-ca-file=/etc/kubernetes/ca.crt \
    --cloud-provider=aws \
    --anonymous-auth=false \
    --cgroup-driver=systemd \
     \
     \

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target`
)
