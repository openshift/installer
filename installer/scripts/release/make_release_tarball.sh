#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/common.env.sh"

echo "Retrieving TerraForm resources"
"$DIR/get_terraform_bins.sh"

echo "Copying Tectonic Installer binaries"
#cp "$ROOT/bin/windows/installer.exe" "$INSTALLER_RELEASE_DIR/windows/installer.exe"
cp "$ROOT/bin/darwin/installer"      "$INSTALLER_RELEASE_DIR/darwin/installer"
cp "$ROOT/bin/linux/installer"       "$INSTALLER_RELEASE_DIR/linux/installer"

echo "Adding TerraForm sources"
for p in "${TERRAFORM_SOURCES[@]}"
do
  cp -r "$p" "$TECTONIC_RELEASE_TOP_DIR"
done

echo "Building release tarball"
tar -cvzf "$ROOT/$TECTONIC_RELEASE_TARBALL_FILE" -C "$TECTONIC_RELEASE_DIR" "${TECTONIC_RELEASE_TOP_DIR#$TECTONIC_RELEASE_DIR/}"

echo "Release tarball is available at $ROOT/$TECTONIC_RELEASE_TARBALL_FILE"
