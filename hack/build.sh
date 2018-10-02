#!/bin/sh

set -ex

cd "$(dirname "$0")/.."

MODE="${MODE:-release}"
LDFLAGS="${LDFLAGS} -X main.version=$(git describe --always --abbrev=0 --dirty)"
TAGS="${TAGS:-}"
export CGO_ENABLED=0

case "${MODE}" in
release)
	TAGS="${TAGS} release"
	GOPATH="${PWD}/vendor:$(go env GOPATH)" go generate ./data
	;;
dev)
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

if (echo "${TAGS}" | grep -q 'libvirt_destroy')
then
	export CGO_ENABLED=1
fi

go build -ldflags "${LDFLAGS}" -tags "${TAGS}" -o ./bin/openshift-install ./cmd/openshift-install
