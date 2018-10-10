package content

import (
	"text/template"
)

const (
	// BootkubeSystemdContents is a service for running bootkube on the bootstrap
	// nodes
	BootkubeSystemdContents = `[Unit]
Description=Bootstrap a Kubernetes cluster
Wants=kubelet.service
After=kubelet.service

[Service]
WorkingDirectory=/opt/tectonic

ExecStart=/opt/tectonic/bootkube.sh

Restart=on-failure
RestartSec=5s`
)

var (
	// BootkubeShFileTemplate is a script file for running bootkube on the
	// bootstrap nodes.
	BootkubeShFileTemplate = template.Must(template.New("bootkube.sh").Parse(`#!/usr/bin/env bash
set -e

mkdir --parents /etc/kubernetes/manifests/

MACHINE_CONFIG_OPERATOR_IMAGE=$(podman run --rm {{.ReleaseImage}} image machine-config-operator)
MACHINE_CONFIG_CONTROLLER_IMAGE=$(podman run --rm {{.ReleaseImage}} image machine-config-controller)
MACHINE_CONFIG_SERVER_IMAGE=$(podman run --rm {{.ReleaseImage}} image machine-config-server)
MACHINE_CONFIG_DAEMON_IMAGE=$(podman run --rm {{.ReleaseImage}} image machine-config-daemon)

KUBE_APISERVER_OPERATOR_IMAGE="{{.ClusterKubeApiserverOperatorImage}}"
KUBE_CONTROLLER_MANAGER_OPERATOR_IMAGE="{{.ClusterKubeControllerManagerOperatorImage}}"
KUBE_SCHEDULER_OPERATOR_IMAGE="{{.ClusterKubeSchedulerOperatorImage}}"

KUBE_APISERVER_OPERATOR_IMAGE=${KUBE_APISERVER_OPERATOR_IMAGE:-$(podman run --rm {{.ReleaseImage}} image cluster-kube-apiserver-operator)}
KUBE_CONTROLLER_MANAGER_OPERATOR_IMAGE=${KUBE_CONTROLLER_MANAGER_OPERATOR_IMAGE:-$(podman run --rm {{.ReleaseImage}} image cluster-kube-controller-manager-operator)}
KUBE_SCHEDULER_OPERATOR_IMAGE=${KUBE_SCHEDULER_OPERATOR_IMAGE:-$(podman run --rm {{.ReleaseImage}} image cluster-kube-scheduler-operator)}

if [ ! -d cvo-bootstrap ]
then
	echo "Rendering Cluster Version Operator Manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		"{{.ReleaseImage}}" \
		render \
			--output-dir=/assets/cvo-bootstrap \
			--release-image="{{.ReleaseImage}}"

	cp --recursive cvo-bootstrap/manifests .
	cp --recursive cvo-bootstrap/bootstrap/bootstrap-pod.yaml /etc/kubernetes/manifests/
fi

if [ ! -d kco-bootstrap ]
then
	echo "Rendering Kubernetes core manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		--volume /etc/kubernetes:/etc/kubernetes:z \
		"{{.KubeCoreRenderImage}}" \
		--config=/assets/kco-config.yaml \
		--output=/assets/kco-bootstrap

	cp --recursive kco-bootstrap/bootstrap-configs /etc/kubernetes/bootstrap-configs
	cp --recursive kco-bootstrap/bootstrap-manifests .
	cp --recursive kco-bootstrap/manifests .
fi

if [ ! -d kube-apiserver-bootstrap ]
then
	echo "Rendering Kubernetes API server core manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		"${KUBE_APISERVER_OPERATOR_IMAGE}" \
		/usr/bin/cluster-kube-apiserver-operator render \
		--asset-input-dir=/assets/tls \
		--asset-output-dir=/assets/kube-apiserver-bootstrap \
		--config-override-file=/usr/share/bootkube/manifests/config/config-overrides.yaml \
		--config-output-file=/assets/kube-apiserver-bootstrap/config

	# TODO: copy the bootstrap manifests to replace kube-core-operator
	cp --recursive kube-apiserver-bootstrap/manifests/00_openshift-kube-apiserver-ns.yaml manifests/00_openshift-kube-apiserver-ns.yaml
	cp --recursive kube-apiserver-bootstrap/manifests/secret-* manifests/
	cp --recursive kube-apiserver-bootstrap/manifests/configmap-* manifests/
fi

if [ ! -d kube-controller-manager-bootstrap ]
then
	echo "Rendering Kubernetes Controller Manager core manifests..."

	# shellcheck disable=SC2154
	podman run \
		--volume "$PWD:/assets:z" \
		"${KUBE_CONTROLLER_MANAGER_OPERATOR_IMAGE}" \
		/usr/bin/cluster-kube-controller-manager-operator render \
		--asset-input-dir=/assets/tls \
		--asset-output-dir=/assets/kube-controller-manager-bootstrap \
		--config-override-file=/usr/share/bootkube/manifests/config/config-overrides.yaml \
		--config-output-file=/assets/kube-controller-manager-bootstrap/config

	# TODO: copy the bootstrap manifests to replace kube-core-operator
	cp --recursive kube-controller-manager-bootstrap/manifests/00_openshift-kube-controller-manager-ns.yaml manifests/00_openshift-kube-controller-manager-ns.yaml
	cp --recursive kube-controller-manager-bootstrap/manifests/secret-* manifests/
	cp --recursive kube-controller-manager-bootstrap/manifests/configmap-* manifests/
fi

if [ ! -d kube-scheduler-bootstrap ]
then
        echo "Rendering Kubernetes Scheduler core manifests..."

        # shellcheck disable=SC2154
        podman run \
                --volume "$PWD:/assets:z" \
                "${KUBE_SCHEDULER_OPERATOR_IMAGE}" \
                /usr/bin/cluster-kube-scheduler-operator render \
                --asset-input-dir=/assets/tls \
                --asset-output-dir=/assets/kube-scheduler-bootstrap \
                --config-override-file=/usr/share/bootkube/manifests/config/config-overrides.yaml \
                --config-output-file=/assets/kube-scheduler-bootstrap/config

        # TODO: copy the bootstrap manifests to replace kube-core-operator
        cp --recursive kube-scheduler-bootstrap/manifests/00_openshift-kube-scheduler-ns.yaml manifests/00_openshift-kube-scheduler-ns.yaml
        cp --recursive kube-scheduler-bootstrap/manifests/secret-* manifests/
        cp --recursive kube-scheduler-bootstrap/manifests/configmap-* manifests/
fi

if [ ! -d mco-bootstrap ]
then
	echo "Rendering MCO manifests..."

	# shellcheck disable=SC2154
	podman run \
		--user 0 \
		--volume "$PWD:/assets:z" \
		"${MACHINE_CONFIG_OPERATOR_IMAGE}" \
		bootstrap \
			--etcd-ca=/assets/tls/etcd-client-ca.crt \
			--root-ca=/assets/tls/root-ca.crt \
			--config-file=/assets/manifests/cluster-config.yaml \
			--dest-dir=/assets/mco-bootstrap \
			--machine-config-controller-image=${MACHINE_CONFIG_CONTROLLER_IMAGE} \
			--machine-config-server-image=${MACHINE_CONFIG_SERVER_IMAGE} \
			--machine-config-daemon-image=${MACHINE_CONFIG_DAEMON_IMAGE} \

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
	"{{.EtcdCertSignerImage}}" \
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
		"{{.EtcdctlImage}}" \
		/usr/local/bin/etcdctl \
		--dial-timeout=10m \
		--cacert=/opt/tectonic/tls/etcd-client-ca.crt \
		--cert=/opt/tectonic/tls/etcd-client.crt \
		--key=/opt/tectonic/tls/etcd-client.key \
		--endpoints={{.EtcdCluster}} \
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
	"{{.BootkubeImage}}" \
	start --asset-dir=/assets`))
)
