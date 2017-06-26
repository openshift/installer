#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

# shellcheck disable=SC1090
source "$DIR/common.env.sh"

# if AWS credentials are not set, then we should not continue
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ];then
    echo "Must export both \$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY"
    return 1
fi

aws s3 cp "$ROOT/$TECTONIC_RELEASE_TARBALL_FILE" "s3://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE_TARBALL_FILE"
echo "Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL"

aws s3 cp "$INSTALLER_RELEASE_DIR/linux/installer" "s3://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE-linux"
aws s3 cp "$INSTALLER_RELEASE_DIR/darwin/installer" "s3://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE-darwin"
#aws s3 cp "$INSTALLER_RELEASE_DIR/windows/installer.exe" "s3://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE-windows.exe"
