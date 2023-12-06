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

minimum_go_version=1.20
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
  make -C terraform all
  copy_terraform_to_mirror # Copy terraform parts to embedded mirror.
fi

# build cluster-api binaries
if [ -n "${OPENSHIFT_INSTALL_CLUSTER_API}" ]; then
  make -C cluster-api all
  copy_cluster_api_to_mirror
fi

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

if (echo "${TAGS}" | grep -q 'libvirt')
then
	export CGO_ENABLED=1
fi

echo "building openshift-install"

# shellcheck disable=SC2086
go build ${GOFLAGS} -gcflags "${GCFLAGS}" -ldflags "${LDFLAGS}" -tags "${TAGS}" -o "${OUTPUT}" ./cmd/openshift-install
