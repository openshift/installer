#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/common.env.sh"

mkdir -p "${TMP_DIR}" "${TERRAFORM_BIN_TMP_DIR}"

set -x

#for os in "linux" "darwin" "windows"; do
for os in "linux" "darwin"; do
    curl -L "${TERRAFORM_BIN_BASE_URL}_${os}_amd64.zip" -o "${TERRAFORM_BIN_TMP_DIR}/terraform_${os}.zip"
    mkdir -p "${INSTALLER_RELEASE_DIR}/${os}/"
    unzip -o "${TERRAFORM_BIN_TMP_DIR}/terraform_${os}.zip" -d "${INSTALLER_RELEASE_DIR}/${os}/"
    chmod +x "${INSTALLER_RELEASE_DIR}/${os}"/terraform*
done
