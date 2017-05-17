#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"
source "$DIR/../awsutil.sh"

# if $VERSION is not set, then we should not continue
if [ -z "${VERSION}" ]; then
    echo "VERSION needs to be set. Exiting."
    exit 1
else
    echo "Found tag ${VERSION}, retrieving binaries from S3"
fi

# if AWS credentials are not set, then we should not continue
check_aws_creds

# get_bin <bucket> <remote> <dest>
function get_bin() {
    mkdir -p `dirname $3`
    aws_download_file "$2" "$3" "$1" "application/octet-stream"
    chmod +x "$3"
}

#get_bin "$TECTONIC_BINARY_BUCKET" "build-artifacts/installer/$VERSION/bin/windows/installer.exe" "$INSTALLER_RELEASE_DIR/windows/installer.exe"
get_bin "$TECTONIC_BINARY_BUCKET" "build-artifacts/installer/$VERSION/bin/darwin/installer"      "$INSTALLER_RELEASE_DIR/darwin/installer"
get_bin "$TECTONIC_BINARY_BUCKET" "build-artifacts/installer/$VERSION/bin/linux/installer"       "$INSTALLER_RELEASE_DIR/linux/installer"
