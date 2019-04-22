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

for file in $(find . -maxdepth 1 -type f | sort)
do
	echo "Creating object from file: $file ..."
	kubectl create --filename "$file"
	echo "Done creating object from file: $file ..."
done

# Workaround for https://github.com/opencontainers/runc/pull/1807
touch /opt/openshift/.openshift.done

echo "OpenShift installation is done"
