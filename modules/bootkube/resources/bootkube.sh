#!/bin/bash
set -e

echo "Rendering Kubernetes core manifests..."

# shellcheck disable=SC2154
/usr/bin/podman run \
	--volume "$PWD:/assets:z" \
	--volume /etc/kubernetes:/etc/kubernetes:z \
	"${kube_core_renderer_image}" \
	--config=/assets/kco-config.yaml \
	--output=/assets

echo "Rendering MCO manifests..."

# shellcheck disable=SC2154
/usr/bin/podman run \
	--user 0 \
	--volume "$PWD:/assets:z" \
	"${machine_config_operator_image}" \
	bootstrap \
		--etcd-ca=/assets/tls/etcd-client-ca.crt \
		--root-ca=/assets/tls/root-ca.crt \
		--config-file=/assets/manifests/cluster-config.yaml \
		--dest-dir=/assets/mco-bootstrap \
		--images-json-configmap=/assets/manifests/machine-config-operator-01-images-configmap.yaml

mkdir -p /etc/kubernetes/manifests/
mkdir -p /etc/mcc/bootstrap/
mkdir -p /etc/ssl/mcs/

# Bootstrap MachineConfigController uses /etc/mcc/bootstrap/manifests/ dir to
# 1. read the controller config rendered by MachineConfigOperator
# 2. read the default MachineConfigPools rendered by MachineConfigOperator
# 3. read any additional MachineConfigs that are needed for the default MachineConfigPools.
cp -r "$PWD/mco-bootstrap/manifests" /etc/mcc/bootstrap/manifests

# /etc/ssl/mcs/tls.{crt, key} are locations for MachineConfigServer's tls assets.
cp "$PWD/tls/machine-config-server.crt" /etc/ssl/mcs/tls.crt
cp "$PWD/tls/machine-config-server.key" /etc/ssl/mcs/tls.key
cp "$PWD/mco-bootstrap/machineconfigoperator-bootstrap-pod.yaml" /etc/kubernetes/manifests/

# We originally wanted to run the etcd cert signer as
# a static pod, but kubelet could't remove static pod
# when API server is not up, so we have to run this as
# podman container.
# See https://github.com/kubernetes/kubernetes/issues/43292

echo "Starting etcd certificate signer..."

# shellcheck disable=SC2154
SIGNER=$(/usr/bin/podman run -d \
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
	--servercertdur=26280h)

echo "Waiting for etcd cluster..."

# Wait for the etcd cluster to come up.
i=0
while true
do
	set +e
	# shellcheck disable=SC2154,SC2086
	/usr/bin/podman run \
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
	status=$?
	set -e

	if [ "$status" -eq 0 ]
	then
		break
	fi

	i=$((i+1))
	[ $i -eq 10 ] && echo "etcdctl failed too many times." && exit 1

	echo "etcdctl failed. Retrying in 5 seconds..."
	sleep 5
done

echo "etcd cluster up. Killing etcd certificate signer..."

/usr/bin/podman kill "$SIGNER"
rm /etc/kubernetes/manifests/machineconfigoperator-bootstrap-pod.yaml

cp -r "$PWD/bootstrap-configs" /etc/kubernetes/bootstrap-configs

echo "Starting bootkube..."

# shellcheck disable=SC2154
/usr/bin/podman run \
	--volume "$PWD:/assets:z" \
	--volume /etc/kubernetes:/etc/kubernetes:z \
	--network=host \
	--entrypoint=/bootkube \
	"${bootkube_image}" \
	start --asset-dir=/assets
