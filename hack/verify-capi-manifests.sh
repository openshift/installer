#!/bin/bash

MANIFESTS_DIR="/go/src/github.com/openshift/installer/data/data/cluster-api"

generate_capi_manifest() {
	provider="$(basename "$1")"

	echo "Generating ${provider} manifest"
	pushd "$1"
	# Parse provider module URL and revision
	# Workaround the import path for azure-service-operator & openstackorc being different from the module path
	provider_go_module="$(grep _ tools.go | awk '{ print $2 }' | sed -En 's/"//g;s#/cmd/(controller|manager)$##; p')"
	mod_info="$(go mod download -json "${provider_go_module}")"
	popd
	version="$(echo "${mod_info}" | jq '.Version' | sed 's|"||g')"
	info_path="$(echo "${mod_info}" | jq '.Info' | sed 's|"||g')"
	repo_origin="$(jq '.Origin.URL' "${info_path}" | sed 's|"||g')"
	revision="$(jq '.Origin.Hash' "${info_path}" | sed 's|"||g')"

	if [ "${provider}" = "azureaso" ]; then
		# Copy the operator YAML and filtered CRDs from upstream
		# List of allowed ASO CRDs (matching cluster-api-provider-azure)
		aso_crds=(
			"resourcegroups.resources.azure.com"
			"natgateways.network.azure.com"
			"managedclusters.containerservice.azure.com"
			"managedclustersagentpools.containerservice.azure.com"
			"bastionhosts.network.azure.com"
			"virtualnetworks.network.azure.com"
			"virtualnetworkssubnets.network.azure.com"
			"privateendpoints.network.azure.com"
			"fleetsmembers.containerservice.azure.com"
			"extensions.kubernetesconfiguration.azure.com"
		)

		# Build the yq filter for allowed CRD names
		set +x
		crd_filter=""
		for crd_name in "${aso_crds[@]}"; do
			crd_filter="${crd_filter}.metadata.name == \"${crd_name}\" or "
		done
		crd_filter="${crd_filter}false"  # Add false at the end to close the OR chain
		set -x

		# Download and filter CRDs (keeping webhooks and other non-CRD resources)
		# We filter by selecting: (not a CRD) OR (CRD with allowed name)
		curl -fSsL "https://github.com/Azure/azure-service-operator/releases/download/${version}/azureserviceoperator_customresourcedefinitions_${version}.yaml" | \
		    yq e ". | select(.kind != \"CustomResourceDefinition\" or (.kind == \"CustomResourceDefinition\" and (${crd_filter})))" - \
			>>"${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
	elif [ "${provider}" = "openstackorc" ]; then
		# Just copy the CRD from upstream
		curl -fSsL "https://github.com/k-orc/openstack-resource-controller/releases/download/${version}/install.yaml" -o "${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
	elif [ "${provider}" = "nutanix" ]; then
		# Download pre-built infrastructure components from GitHub releases
		curl -fSsL "https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/releases/download/${version}/infrastructure-components.yaml" -o "${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
	else
		# Generate provider manifest from specified revision
		clone_path="$(mktemp -d)"
		git clone "${repo_origin}" "${clone_path}"
		pushd "${clone_path}"
		git checkout "${revision}"

		case "${provider}" in
		vsphere)
			make release-manifests-all
			;;
		*)
			make release-manifests
			;;
		esac

		case "${provider}" in
		cluster-api)
			cp out/cluster-api-components.yaml "${MANIFESTS_DIR}/core-components.yaml"
			;;
		*)
			cp out/infrastructure-components.yaml "${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
			;;
		esac
		popd
		rm -rf "${clone_path}"
	fi
}

if [ "$IS_CONTAINER" != "" ]; then
	set -eux

	# Install `jq` if not present
	if ! command -v jq; then
		curl -L https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 -o /usr/bin/jq
		chmod u+x /usr/bin/jq
	fi

	# Install `yq` if not present
	if ! command -v yq; then
		curl -L https://github.com/mikefarah/yq/releases/download/v4.44.6/yq_linux_amd64 -o /usr/bin/yq
		chmod u+x /usr/bin/yq
	fi

	# Silence git hints and advices
	git config --global init.defaultBranch master
	git config --global advice.detachedHead false

	if [ $# -gt 0 ]; then
		for target in "${@}"; do
			generate_capi_manifest "${target}"
		done
	else
		find cluster-api/providers -maxdepth 1 -mindepth 1 -type d -print0 | while read -r -d '' dir; do
			generate_capi_manifest "${dir}"
		done
		generate_capi_manifest "cluster-api/cluster-api"
	fi

	git diff --exit-code
else
	podman run --rm \
		--env IS_CONTAINER=TRUE \
		--volume "${PWD}:/go/src/github.com/openshift/installer:z" \
		--workdir /go/src/github.com/openshift/installer \
		docker.io/golang:1.24 \
		./hack/verify-capi-manifests.sh "${@}"
fi
