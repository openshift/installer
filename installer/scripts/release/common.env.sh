#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../.."
REPOSITORY_ROOT="$ROOT/.."
WORKSPACE_DIR="$ROOT/.workspace"
TMP_DIR="$WORKSPACE_DIR/tmpdir"

export VERSION=${VERSION:-$("$REPOSITORY_ROOT/git-version")}

export TECTONIC_RELEASE_BUCKET=releases.tectonic.com
export TECTONIC_BINARY_BUCKET=tectonic-release
export TECTONIC_RELEASE="tectonic-$VERSION"
export TECTONIC_RELEASE_TARBALL_FILE="$TECTONIC_RELEASE.tar.gz"
export TECTONIC_RELEASE_TARBALL_URL="https://$TECTONIC_RELEASE_BUCKET/$TECTONIC_RELEASE_TARBALL_FILE"
export TECTONIC_RELEASE_DIR="$WORKSPACE_DIR/$VERSION"
export TECTONIC_RELEASE_TOP_DIR="$TECTONIC_RELEASE_DIR/tectonic"
export INSTALLER_RELEASE_DIR="$TECTONIC_RELEASE_TOP_DIR/tectonic-installer"

export TERRAFORM_BIN_TMP_DIR="$TMP_DIR/terraform-bin"
export TERRAFORM_BIN_VERSION=0.10.4
export TERRAFORM_BIN_BASE_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_BIN_VERSION}/terraform_${TERRAFORM_BIN_VERSION}"
export TERRAFORM_LICENSE_URL="https://raw.githubusercontent.com/hashicorp/terraform/v${TERRAFORM_BIN_VERSION}/LICENSE"
export TERRAFORM_SOURCES=(
  "${REPOSITORY_ROOT}/modules"
  "${REPOSITORY_ROOT}/platforms"
  "${REPOSITORY_ROOT}/config.tf"
  "${REPOSITORY_ROOT}/examples"
)

export PROVIDER_MATCHBOX_VERSION=v0.2.2
export PROVIDER_MATCHBOX_BASE_URL="https://github.com/coreos/terraform-provider-matchbox/releases/download/${PROVIDER_MATCHBOX_VERSION}/"
