#!/bin/sh

set -ex

# shellcheck disable=SC2068
version() { IFS="."; printf "%03d%03d%03d\\n" $@; unset IFS;}

minimum_go_version=1.14
current_go_version=$(go version | cut -d " " -f 3)

if [ "$(version "${current_go_version#go}")" -lt "$(version "$minimum_go_version")" ]; then
     echo "Go version should be greater or equal to $minimum_go_version"
     exit 1
fi

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
if test "$(go env GOARCH)" = "arm64"
then
	# https://github.com/golang/go/issues/40492
	LDFLAGS="${LDFLAGS} -linkmode external"
	export CGO_ENABLED=1
fi

# shellcheck disable=SC2086
go build ${GOFLAGS} -ldflags "${LDFLAGS}" -tags "${TAGS}" -o "${OUTPUT}" ./cmd/openshift-install
