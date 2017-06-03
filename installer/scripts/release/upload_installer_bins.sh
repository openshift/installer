#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/common.env.sh"
# shellcheck disable=SC1090
source "$DIR/../awsutil.sh"
GIT_SHA=$("$ROOT/../git-version")

# if AWS credentials are not set, then we should not continue
check_aws_creds

#aws_upload_file "$ROOT/bin/windows/installer.exe" "build-artifacts/installer/$GIT_SHA/bin/windows/installer.exe" "$TECTONIC_BINARY_BUCKET" "application/octet-stream"
aws_upload_file "$ROOT/bin/darwin/installer" "build-artifacts/installer/$GIT_SHA/bin/darwin/installer" "$TECTONIC_BINARY_BUCKET" "application/octet-stream"
aws_upload_file "$ROOT/bin/linux/installer" "build-artifacts/installer/$GIT_SHA/bin/linux/installer" "$TECTONIC_BINARY_BUCKET" "application/octet-stream"
