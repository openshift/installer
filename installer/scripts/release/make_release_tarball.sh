#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"
source "$DIR/../awsutil.sh"

check_aws_creds

echo "Retrieving matchbox release"
"$DIR/get_matchbox_release.sh"

echo "Retrieving TerraForm resources"
"$DIR/get_terraform_bins.sh"

echo "Retrieving Tectonic Installer binaries"
"$DIR/get_installer_bins.sh"

echo "Adding TerraForm sources"
cp -r $TERRAFORM_SOURCES "$TECTONIC_RELEASE_TOP_DIR"

echo "Building release tarball"
tar -cvzf "$ROOT/$TECTONIC_RELEASE_TARBALL_FILE" -C "$TECTONIC_RELEASE_DIR" .

echo "Release tarball is available at $ROOT/$TECTONIC_RELEASE_TARBALL_FILE"
