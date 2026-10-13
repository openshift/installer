#!/bin/sh
#
# Prepare for a release.  Usage:
#
#   $ hack/release.sh v0.1.0

set -ex

cd "$(dirname "$0")"

TAG="${1}"

git tag -sm "version ${TAG}" "${TAG}"
./build.sh # ensure freshly-generated data
for GOOS in darwin linux
do
	GOARCH=amd64
	OUTPUT="bin/openshift-install-${GOOS}-${GOARCH}"
	GOOS="${GOOS}" GOARCH="${GOARCH}" OUTPUT="${OUTPUT}" SKIP_GENERATION=y ./build.sh
done
(
	cd ../bin
	sha256sum openshift-install-* >release.sha256
	gpg --output release.sha256.sig --detach-sig release.sha256
)
