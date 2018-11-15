#!/usr/bin/env bash
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

for file in $(find . -maxdepth 1 -type f | sort)
do
	echo "Creating object from file: $file ..."
	kubectl create --filename "$file"
	echo "Done creating object from file: $file ..."
done

# Workaround for https://github.com/opencontainers/runc/pull/1807
touch /opt/tectonic/.tectonic.done

echo "Tectonic installation is done"
