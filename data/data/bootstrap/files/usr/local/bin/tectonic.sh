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

# This is a temporary hack to set up development clusters with allowAll idp, please remove once the default idp is worked out for 4.0
set_up_allow_all_idp() {
    echo "Configuring cluster with allowAll identity-provider..."
    set -x
    url=$(oc --config="$KUBECONFIG" status | grep -o 'https.*')
    cat > /tmp/idp-config-patch << EOF
{"spec":{"unsupportedConfigOverrides":{"oauthConfig":{"identityProviders":[{"challenge":true,"login":true,"name":"anypassword","provider":{"apiVersion":"v1","kind":"AllowAllPasswordIdentityProvider"}}],"masterCA":"/etc/kubernetes/static-pod-resources/configmaps/kubelet-serving-ca/ca-bundle.crt","masterPublicURL":"${url}","masterURL":"${url}"}}}}
EOF
    oc --config="$KUBECONFIG" patch clusterversions.config.openshift.io/version -p '{"spec":{"overrides": [{"kind": "Deployment","name": "openshift-cluster-kube-apiserver-operator","unmanaged": true}]}}' --type merge
    oc --config="$KUBECONFIG" patch kubeapiserveroperatorconfig/instance -p "$(cat /tmp/idp-config-patch)" --type merge
}

# Wait for Kubernetes pods
wait_for_pods kube-system

for file in $(find . -maxdepth 1 -type f | sort)
do
	echo "Creating object from file: $file ..."
	kubectl create --filename "$file"
	echo "Done creating object from file: $file ..."
done

# Again, temporary hack, please remove once official default idp is configured
set_up_allow_all_idp

# Workaround for https://github.com/opencontainers/runc/pull/1807
touch /opt/tectonic/.tectonic.done

echo "Tectonic installation is done"
