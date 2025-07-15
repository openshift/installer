#!/bin/bash

MANIFESTS_DIR="/go/src/github.com/openshift/installer/data/data/cluster-api"

# Generate provider manifest from released assets
generate_capi_manifest_from_released_assets() {
	echo "Generating ${provider} manifest from released assets"
	provider="$1"
	repo_origin="$2"
	version="$3"

	# Not a version, but a revision
	if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$  ]]; then
		return 1
	fi

	# Core CAPI generates cluster-api-components.yaml
	# while provider generates infrastructure-components.yaml
	case "${provider}" in
	cluster-api)
		asset_name="cluster-api-components.yaml"
		saved_asset_name="${MANIFESTS_DIR}/core-components.yaml"
		;;
	*)
		asset_name="infrastructure-components.yaml"
		saved_asset_name="${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
		;;
	esac

	if ! curl -fSsL -o "${saved_asset_name}" "$repo_origin/releases/download/${version}/${asset_name}"; then
		echo "Failed generating ${provider} manifest from released assets. Falling back to generate from specified revision"
		return 1
	fi
}

# Generate provider manifest from specified revision
generate_capi_manifest_from_revision() {
	echo "Generating ${provider} manifest from specified revision"
	provider="$1"
	repo_origin="$2"
	revision="$3"

	clone_path="$(mktemp -d)"
	git clone "${repo_origin}" "${clone_path}"
	pushd "${clone_path}"
	git fetch "${repo_origin}" "${revision}"
	git checkout "${revision}"

	# Provider-specific generate command
	case "${provider}" in
	vsphere)
		make release-manifests-all
		;;
	*)
		make release-manifests
		;;
	esac

	# Core CAPI generates cluster-api-components.yaml
	# while provider generates infrastructure-components.yaml
	# except azureaso that needs combining 2 manifests.
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
}

generate_capi_manifest() {
	provider="$(basename "$1")"

	echo "Generating ${provider} manifest"
	pushd "$1"
	# Parse provider module URL and revision
	# Workaround the import path for azure-service-operator being different from the module path
	provider_go_module="$(grep _ tools.go | awk '{ print $2 }' | sed 's|"||g' | sed 's|/cmd/controller$||g')"
	mod_info="$(go mod download -json "${provider_go_module}")"
	popd
	version="$(echo "${mod_info}" | jq '.Version' | sed 's|"||g')"
	info_path="$(echo "${mod_info}" | jq '.Info' | sed 's|"||g')"
	repo_origin="$(jq '.Origin.URL' "${info_path}" | sed 's|"||g')"
	revision="$(jq '.Origin.Hash' "${info_path}" | sed 's|"||g')"

	case "${provider}" in
		azurestack)
		    # skip this for now--until unforked
			;;
		azureaso)
			# Just copy the CRD from upstream release assets
			curl -fSsL "https://github.com/Azure/azure-service-operator/releases/download/${version}/azureserviceoperator_${version}.yaml" -o "${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
			echo "---" >>"${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
			curl -fSsL "https://github.com/Azure/azure-service-operator/releases/download/${version}/azureserviceoperator_customresourcedefinitions_${version}.yaml" >>"${MANIFESTS_DIR}/${provider}-infrastructure-components.yaml"
			;;
		*)
			# Attempt to find the infrastructure manifest in the released assets
			# If none is found, generate the infrastucture manifests from the pinned revision
			generate_capi_manifest_from_released_assets "$provider" "$repo_origin" "$version" || \
			generate_capi_manifest_from_revision "$provider" "$repo_origin" "$revision"
			;;
	esac
}

if [ "$IS_CONTAINER" != "" ]; then
	set -eux

	# Install `jq` if not present
	if ! command -v jq >/dev/null 2>&1; then
		curl -L https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 -o /usr/bin/jq
		chmod u+x /usr/bin/jq
	fi

	# Install `controller-gen` & `kustomize`, which are needed by nutanix, if not present
	if ! command -v controller-gen >/dev/null 2>&1; then
		go install sigs.k8s.io/controller-tools/cmd/controller-gen
	fi

	if ! command -v kustomize >/dev/null 2>&1; then
		go install sigs.k8s.io/kustomize/kustomize/v5@latest
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
		docker.io/golang:1.23 \
		./hack/verify-capi-manifests.sh "${@}"
fi
