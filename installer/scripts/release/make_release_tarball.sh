#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"
source "$DIR/../awsutil.sh"

check_aws_creds

echo "Retrieving matchbox release"
"$DIR/get_matchbox_release.sh"

echo "Retrieving Tectonic Installer binaries"
"$DIR/get_installer_bins.sh"

tar -cvzf "$ROOT/$TECTONIC_RELEASE_TARBALL_FILE" -C "$TECTONIC_RELEASE_DIR" .

echo "Release tarball is available at $ROOT/$TECTONIC_RELEASE_TARBALL_FILE"
