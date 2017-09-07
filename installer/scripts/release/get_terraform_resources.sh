#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/common.env.sh"

mkdir -p "${TMP_DIR}" "${TERRAFORM_BIN_TMP_DIR}"

set -x

#for os in "linux" "darwin" "windows"; do
for os in "linux" "darwin"; do
    mkdir -p "${INSTALLER_RELEASE_DIR}/${os}/"
    curl -L "${PROVIDER_MATCHBOX_BASE_URL}/terraform-provider-matchbox-${PROVIDER_MATCHBOX_VERSION}-${os}-amd64.tar.gz" -o "${TERRAFORM_BIN_TMP_DIR}/terraform-provider-matchbox-${PROVIDER_MATCHBOX_VERSION}-${os}-amd64.tar.gz"
    pushd "${TERRAFORM_BIN_TMP_DIR}"
    tar xvf "terraform-provider-matchbox-${PROVIDER_MATCHBOX_VERSION}-${os}-amd64.tar.gz"
    cp "./terraform-provider-matchbox-${PROVIDER_MATCHBOX_VERSION}-${os}-amd64/terraform-provider-matchbox" "${INSTALLER_RELEASE_DIR}/${os}/"
    popd
    curl -L "${TERRAFORM_BIN_BASE_URL}_${os}_amd64.zip" -o "${TERRAFORM_BIN_TMP_DIR}/terraform_${os}.zip"
    unzip -o "${TERRAFORM_BIN_TMP_DIR}/terraform_${os}.zip" -d "${INSTALLER_RELEASE_DIR}/${os}/"
    chmod +x "${INSTALLER_RELEASE_DIR}/${os}"/terraform*
done

curl --fail --silent -L "${TERRAFORM_LICENSE_URL}" -o "${TECTONIC_RELEASE_TOP_DIR}"/TERRAFORM_LICENSE
