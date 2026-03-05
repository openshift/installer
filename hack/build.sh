#!/bin/sh

set -ex

# Source the Cluster API build script.
# shellcheck source=hack/build-cluster-api.sh
. "$(dirname "$0")/build-cluster-api.sh"

# shellcheck disable=SC2068
version() { IFS="."; printf "%03d%03d%03d\\n" $@; unset IFS;}

minimum_go_version=1.23
current_go_version=$(go version | cut -d " " -f 3)

if [ "$(version "${current_go_version#go}")" -lt "$(version "$minimum_go_version")" ]; then
     echo "Go version should be greater or equal to $minimum_go_version"
     exit 1
fi

export CGO_ENABLED=0
MODE="${MODE:-release}"

# if -j [::numeric::] is in MAKEFLAGS, parse it out and pass it along to go build
MAKEJOBS="$(echo "'${MAKEFLAGS}'" | sed -nE 's|.*-j\s+?([1-9]+).*|\1|p')"
# shellcheck disable=SC2086
if [ -n "${MAKEJOBS}" ] && [ ${MAKEJOBS} -gt 0 ]; then
        GOBUILDFLAGS="${GOBUILDFLAGS} -p ${MAKEJOBS}"
fi

# pass along to any forked processes (sub make files and go build processes)
export GOBUILDFLAGS
export MAKEFLAGS

# build cluster-api binaries
# shellcheck disable=SC2086
make ${MAKEFLAGS} -C cluster-api all
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
go build ${GOBUILDFLAGS} ${GOFLAGS} -gcflags "${GCFLAGS}" -ldflags "${LDFLAGS}" -tags "${TAGS}" -o "${OUTPUT}" ./cmd/openshift-install
