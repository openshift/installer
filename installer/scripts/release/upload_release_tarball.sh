#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."
source "$DIR/../awsutil.sh"
source "$DIR/common.env.sh"

check_aws_creds

aws_upload_file "$ROOT/$TECTONIC_RELEASE_TARBALL_FILE" "$TECTONIC_RELEASE_TARBALL_FILE" "$TECTONIC_RELEASE_BUCKET" application/x-gzip
echo "Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL"
aws_upload_file "$INSTALLER_RELEASE_DIR/linux/installer" "$TECTONIC_RELEASE-linux" "$TECTONIC_RELEASE_BUCKET" application/octet-stream
aws_upload_file "$INSTALLER_RELEASE_DIR/darwin/installer" "$TECTONIC_RELEASE-darwin" "$TECTONIC_RELEASE_BUCKET" application/octet-stream
#aws_upload_file "$INSTALLER_RELEASE_DIR/windows/installer.exe" "$TECTONIC_RELEASE-windows.exe" "$TECTONIC_RELEASE_BUCKET" application/octet-stream
