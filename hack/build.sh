#!/bin/sh

set -ex

cd "$(dirname "$0")/.."

MODE="${MODE:-release}"
TAGS=

case "${MODE}" in
release)
	TAGS=release
	GOPATH="${PWD}/vendor:${GOPATH}" go generate ./data
	;;
dev)
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

CGO_ENABLED=0 go build -tags "${TAGS}" -o ./bin/openshift-install ./cmd/openshift-install
