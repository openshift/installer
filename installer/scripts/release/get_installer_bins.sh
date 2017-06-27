#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck disable=SC1090
source "$DIR/common.env.sh"

# if $VERSION is not set, then we should not continue
if [ -z "${VERSION}" ]; then
    echo "VERSION needs to be set. Exiting."
    exit 1
else
    echo "Found tag ${VERSION}, retrieving binaries from S3"
fi

# if AWS credentials are not set, then we should not continue
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ];then
    echo "Must export both \$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY"
    return 1
fi

# get_bin <bucket> <remote> <dest>
function get_bin() {
    mkdir -p "$(dirname "$2")"
    aws s3 cp "s3://$1" "$2"
    chmod +x "$2"
}

#get_bin "$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$VERSION/bin/windows/installer.exe" "$INSTALLER_RELEASE_DIR/windows/installer.exe"
get_bin "$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$VERSION/bin/darwin/installer"      "$INSTALLER_RELEASE_DIR/darwin/installer"
get_bin "$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$VERSION/bin/linux/installer"       "$INSTALLER_RELEASE_DIR/linux/installer"
