#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck disable=SC1090
source "$DIR/common.env.sh"
GIT_SHA=$("$ROOT/../git-version")

# if AWS credentials are not set, then we should not continue
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ];then
    echo "Must export both \$AWS_ACCESS_KEY_ID and \$AWS_SECRET_ACCESS_KEY"
    return 1
fi

#aws s3 cp "$ROOT/bin/windows/installer.exe" "s3://$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$GIT_SHA/bin/windows/installer.exe"
aws s3 cp "$ROOT/bin/darwin/installer" "s3://$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$GIT_SHA/bin/darwin/installer"
aws s3 cp "$ROOT/bin/linux/installer" "s3://$TECTONIC_BINARY_BUCKET/build-artifacts/installer/$GIT_SHA/bin/linux/installer"
