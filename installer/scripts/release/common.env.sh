#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../.."
REPOSITORY_ROOT="$ROOT/.."
WORKSPACE_DIR="$ROOT/.workspace"
TMP_DIR="$WORKSPACE_DIR/tmpdir"

export TECTONIC_RELEASE_BUCKET=releases.tectonic.com
export TECTONIC_BINARY_BUCKET=tectonic-release
export TECTONIC_RELEASE="tectonic-$VERSION"
export TECTONIC_RELEASE_TARBALL_FILE="$TECTONIC_RELEASE.tar.gz"
export TECTONIC_RELEASE_TARBALL_URL="https://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE_TARBALL_FILE"
export TECTONIC_RELEASE_DIR="$WORKSPACE_DIR/$VERSION"
export TECTONIC_RELEASE_TOP_DIR="$TECTONIC_RELEASE_DIR/tectonic"
export INSTALLER_RELEASE_DIR="$TECTONIC_RELEASE_TOP_DIR/tectonic-installer"

export TERRAFORM_BIN_TMP_DIR="$TMP_DIR/terraform-bin"
export TERRAFORM_BIN_VERSION=0.9.6
export TERRAFORM_BIN_BASE_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_BIN_VERSION}/terraform_${TERRAFORM_BIN_VERSION}"
export TERRAFORM_SOURCES=(
  "${REPOSITORY_ROOT}/modules"
  "${REPOSITORY_ROOT}/platforms"
  "${REPOSITORY_ROOT}/config.tf"
  "${REPOSITORY_ROOT}/terraformrc.example"
  "${REPOSITORY_ROOT}/examples"
)
