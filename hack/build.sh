#!/bin/sh

set -ex

# Source the Cluster API build script.
# shellcheck source=hack/build-cluster-api.sh
. "$(dirname "$0")/build-cluster-api.sh"

# shellcheck disable=SC2068
version() { IFS="."; printf "%03d%03d%03d\\n" $@; unset IFS;}

# Copy the terraform binary and providers to the mirror to be embedded in the installer binary.
copy_terraform_to_mirror() {
  TARGET_OS_ARCH=$(go env GOOS)_$(go env GOARCH)

  # Clean the mirror, but preserve the README file.
  rm -rf "${PWD}"/pkg/terraform/providers/mirror/*/

  # Copy local terraform providers into data
  find "${PWD}/terraform/bin/${TARGET_OS_ARCH}/" -maxdepth 1 -name "terraform-provider-*.zip" -exec bash -c '
      providerName="$(basename "$1" | cut -d '-' -f 3 | cut -d '.' -f 1)"
      targetOSArch="$2"
      dstDir="${PWD}/pkg/terraform/providers/mirror/openshift/local/$providerName"
      mkdir -p "$dstDir"
      echo "Copying $providerName provider to mirror"
      cp "$1" "$dstDir/terraform-provider-${providerName}_1.0.0_${targetOSArch}.zip"
    ' shell {} "${TARGET_OS_ARCH}" \;

  mkdir -p "${PWD}/pkg/terraform/providers/mirror/terraform/"
  cp "${PWD}/terraform/bin/${TARGET_OS_ARCH}/terraform" "${PWD}/pkg/terraform/providers/mirror/terraform/"
}

# Check if a provider has changed based on the version stored in the binary
check_module_changes() {
	binpath="$1"
	srcpath="$2"
	# Check if a provider has changed based on its go.mod git hash
	version_info="$(go version -m "${binpath}" | grep -Eo 'main.builtGoModHash=[a-fA-F0-9]+' | cut -f2 -d'=')"
	test "${version_info}" == "$(git hash-object "${srcpath}/go.mod")"
}

# Build terraform and providers only if needed
build_terraform_and_providers() {
	TARGET_OS_ARCH=$(go env GOOS)_$(go env GOARCH)
	bindir="${PWD}/terraform/bin/${TARGET_OS_ARCH}"
	find "${PWD}/terraform/providers/" -maxdepth 1 -mindepth 1 -type d | while read -r dir; do
		provider="$(basename "${dir}")"
		binpath="${bindir}/terraform-provider-${provider}"
		if test -s "${binpath}" && test -s "${binpath}.zip" && check_module_changes "${binpath}" "terraform/providers/${provider}"; then
			echo "${provider} is up-to-date"
		else
			echo "Rebuilding ${provider}"
			make -C terraform "go-build.${provider}"
		fi
	done
	if test -s "${bindir}/terraform" && check_module_changes "${bindir}/terraform" "terraform/terraform"; then
		echo "terraform is up-to-date"
	else
		echo "Rebuilding terraform"
		make -C terraform go-build-terraform
	fi
}

minimum_go_version=1.22
current_go_version=$(go version | cut -d " " -f 3)

if [ "$(version "${current_go_version#go}")" -lt "$(version "$minimum_go_version")" ]; then
     echo "Go version should be greater or equal to $minimum_go_version"
     exit 1
fi

export CGO_ENABLED=0
MODE="${MODE:-release}"

# Build terraform binaries before setting environment variables since it messes up make
if test "${SKIP_TERRAFORM}" != y && ! (echo "${TAGS}" | grep -q -e 'aro' -e 'altinfra')
then
	build_terraform_and_providers
	copy_terraform_to_mirror # Copy terraform parts to embedded mirror.
fi

# build cluster-api binaries
make -C cluster-api all
copy_cluster_api_to_mirror

GIT_COMMIT="${SOURCE_GIT_COMMIT:-$(git rev-parse --verify 'HEAD^{commit}')}"
GIT_TAG="${BUILD_VERSION:-$(git describe --always --abbrev=40 --dirty)}"
DEFAULT_ARCH="${DEFAULT_ARCH:-amd64}"
GOFLAGS="${GOFLAGS:--mod=vendor}"
GCFLAGS=""
LDFLAGS="${LDFLAGS} -X github.com/openshift/installer/pkg/version.Raw=${GIT_TAG} -X github.com/openshift/installer/pkg/version.Commit=${GIT_COMMIT} -X github.com/openshift/installer/pkg/version.defaultArch=${DEFAULT_ARCH}"
TAGS="${TAGS:-}"
OUTPUT="${OUTPUT:-bin/openshift-install}"

case "${MODE}" in
release)
	LDFLAGS="${LDFLAGS} -s -w"
	TAGS="${TAGS} release"
	;;
dev)
    GCFLAGS="${GCFLAGS} all=-N -l"
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

if test "${SKIP_GENERATION}" != y
then
	# this step has to be run natively, even when cross-compiling
	GOOS='' GOARCH='' go generate ./data
fi

if (echo "${TAGS}" | grep -q '\bfipscapable\b')
then
	export CGO_ENABLED=1
fi

echo "building openshift-install"

# shellcheck disable=SC2086
go build ${GOFLAGS} -gcflags "${GCFLAGS}" -ldflags "${LDFLAGS}" -tags "${TAGS}" -o "${OUTPUT}" ./cmd/openshift-install
