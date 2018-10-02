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

for file in $(find . -type f -maxdepth 1)
do
	echo "Creating object from file: $file ..."
	kubectl create --filename "$file"
	echo "Done creating object from file: $file ..."
done

# Wait for Tectonic pods
wait_for_pods tectonic-system

echo "Tectonic installation is done"`
)
