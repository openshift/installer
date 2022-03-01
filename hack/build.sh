#!/bin/sh

set -ex

# shellcheck disable=SC2068
version() { IFS="."; printf "%03d%03d%03d\\n" $@; unset IFS;}

# Copy the terraform binary and providers to the mirror to be embedded in the installer binary.
copy_terraform_to_mirror() {
  TARGET_OS_ARCH=$(go env GOOS)_$(go env GOARCH)

  # Clean the mirror, but preserve the README file.
  rm -rf "${PWD}"/pkg/terraform/providers/mirror/*/

  # Copy local terraform providers into data
  find "${PWD}/terraform/bin/" -maxdepth 1 -name "terraform-provider-*.zip" -exec bash -c '
      providerName="$(basename "$1" | cut -d '-' -f 3 | cut -d '.' -f 1)"
      targetOSArch="$2"
      dstDir="${PWD}/pkg/terraform/providers/mirror/openshift/local/$providerName"
      mkdir -p "$dstDir"
      echo "Copying $providerName provider to mirror"
      cp "$1" "$dstDir/terraform-provider-${providerName}_1.0.0_${targetOSArch}.zip"
    ' shell {} "${TARGET_OS_ARCH}" \;

  mkdir -p "${PWD}/pkg/terraform/providers/mirror/terraform/"
  cp "${PWD}/terraform/bin/terraform" "${PWD}/pkg/terraform/providers/mirror/terraform/"
}

minimum_go_version=1.17
current_go_version=$(go version | cut -d " " -f 3)

if [ "$(version "${current_go_version#go}")" -lt "$(version "$minimum_go_version")" ]; then
     echo "Go version should be greater or equal to $minimum_go_version"
     exit 1
fi

# build terraform binaries before setting environment variables since it messes up make
# XXX: need to figure out a good way of determining of dependencies have changed
if [ -z "${SKIP_TERRAFORM}" ]; then
	env TFBINDIR="${PWD}/terraform/bin" make -C terraform all
fi

# Copy terraform parts to embedded mirror.
copy_terraform_to_mirror

MODE="${MODE:-release}"
GIT_COMMIT="${SOURCE_GIT_COMMIT:-$(git rev-parse --verify 'HEAD^{commit}')}"
GIT_TAG="${BUILD_VERSION:-$(git describe --always --abbrev=40 --dirty)}"
DEFAULT_ARCH="${DEFAULT_ARCH:-amd64}"
GOFLAGS="${GOFLAGS:--mod=vendor}"
LDFLAGS="${LDFLAGS} -X github.com/openshift/installer/pkg/version.Raw=${GIT_TAG} -X github.com/openshift/installer/pkg/version.Commit=${GIT_COMMIT} -X github.com/openshift/installer/pkg/version.defaultArch=${DEFAULT_ARCH}"
TAGS="${TAGS:-}"
OUTPUT="${OUTPUT:-bin/openshift-install}"
export CGO_ENABLED=0

case "${MODE}" in
release)
	LDFLAGS="${LDFLAGS} -s -w"
	TAGS="${TAGS} release"
	if test "${SKIP_GENERATION}" != y
	then
		go generate ./data
	fi
	;;
dev)
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

if (echo "${TAGS}" | grep -q 'libvirt')
then
	export CGO_ENABLED=1
fi

# shellcheck disable=SC2086
go build ${GOFLAGS} -ldflags "${LDFLAGS}" -tags "${TAGS}" -o "${OUTPUT}" ./cmd/openshift-install
