#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

if [ -n "${CAPI_VERSION:-}" ]; then
  VERSION="$CAPI_VERSION"
else
  echo "CAPI_VERSION not set, querying latest release from kubernetes-sigs/cluster-api..."
  if ! VERSION=$(curl -sL "https://api.github.com/repos/kubernetes-sigs/cluster-api/releases/latest" 2>/dev/null | awk -F '"' '/tag_name/ { print $4; exit }'); then
    echo "Failed to query GitHub API for latest release." >&2
    exit 1
  fi
  if [ -z "$VERSION" ]; then
    echo "Failed to discover latest cluster-api release. Set CAPI_VERSION environment variable to continue." >&2
    exit 1
  fi
  echo "Using discovered version: $VERSION"
fi


pushd ./cluster-api

echo "Updating CAPI to $VERSION in cluster-api"
GO111MODULE=on go get "sigs.k8s.io/cluster-api@${VERSION}"
echo "Running: go mod tidy"
GO111MODULE=on go mod tidy
echo "Running: go mod vendor"
GO111MODULE=on go mod vendor
popd

echo "Updating CAPI to $VERSION in root directory"
pushd "$ROOT_DIR/cluster-api"
GO111MODULE=on go get "sigs.k8s.io/cluster-api@${VERSION}"
echo "Running: go mod tidy"
GO111MODULE=on go mod tidy
echo "Running: go mod vendor"
GO111MODULE=on go mod vendor
popd

echo "Done."
