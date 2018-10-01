#!/bin/sh

set -ex

cd "$(dirname "$0")/.."

MODE="${MODE:-release}"
TAGS="${TAGS:-}"
export CGO_ENABLED=0

case "${MODE}" in
release)
	TAGS="${TAGS} release"
	GOPATH="${PWD}/vendor:${GOPATH}" go generate ./data
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

go build -tags "${TAGS}" -o ./bin/openshift-install ./cmd/openshift-install
