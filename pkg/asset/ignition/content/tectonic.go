package content

const (
	// TectonicSystemdContents is a service that runs tectonic on the masters.
	TectonicSystemdContents = `[Unit]
Description=Bootstrap a Tectonic cluster
Wants=bootkube.service
After=bootkube.service

[Service]
WorkingDirectory=/opt/tectonic/tectonic

ExecStart=/opt/tectonic/tectonic.sh /opt/tectonic/auth/kubeconfig

Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target`

	// TectonicShFileContents is a script file for running tectonic on bootstrap
	// nodes.
	TectonicShFileContents = `#!/usr/bin/env bash
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
kubectl create --filename updater/app-version-kind.yaml
kubectl create --filename updater/migration-status-kind.yaml

kubectl create --filename updater/operators/kube-core-operator.yaml
kubectl create --filename updater/operators/kube-addon-operator.yaml
kubectl create --filename updater/operators/tectonic-ingress-controller-operator.yaml

kubectl --namespace=tectonic-system get customresourcedefinition appversions.tco.coreos.com
kubectl create --filename updater/app_versions/app-version-tectonic-cluster.yaml
kubectl create --filename updater/app_versions/app-version-kube-core.yaml
kubectl create --filename updater/app_versions/app-version-kube-addon.yaml
kubectl create --filename updater/app_versions/app-version-tectonic-ingress.yaml

# Wait for Tectonic pods
wait_for_pods tectonic-system

echo "Tectonic installation is done"`
)
