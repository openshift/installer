#!/usr/bin/env bash
set -e

mkdir --parents /etc/kubernetes/manifests/

if [ ! -d kco-bootstrap ]
then
	echo "Rendering Kubernetes core manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		--volume /etc/kubernetes:/etc/kubernetes:z \
		"${kube_core_renderer_image}" \
		--config=/assets/kco-config.yaml \
		--output=/assets/kco-bootstrap

	cp --recursive kco-bootstrap/bootstrap-configs /etc/kubernetes/bootstrap-configs
	cp --recursive kco-bootstrap/bootstrap-manifests .
	cp --recursive kco-bootstrap/manifests .
fi

if [ ! -d "mco-bootstrap" ]
then
	echo "Rendering MCO manifests..."

	# shellcheck disable=SC2154
	podman run \
		--user 0 \
		--volume "$PWD:/assets:z" \
		"${machine_config_operator_image}" \
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
	--tmpfs /tmp \
	--volume /opt/tectonic/tls:/opt/tectonic/tls:ro,z \
	--network host \
	"${etcd_cert_signer_image}" \
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
		"${etcdctl_image}" \
		/usr/local/bin/etcdctl \
		--dial-timeout=10m \
		--cacert=/opt/tectonic/tls/etcd-client-ca.crt \
		--cert=/opt/tectonic/tls/etcd-client.crt \
		--key=/opt/tectonic/tls/etcd-client.key \
		--endpoints=${etcd_cluster} \
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
	"${bootkube_image}" \
	start --asset-dir=/assets
