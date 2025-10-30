#!/usr/bin/env bash
set -euo pipefail -o errtrace -x

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$ROOT_DIR"

if [ -n "${CAPA_VERSION:-}" ]; then
  REF="$CAPA_VERSION"
else
  echo "CAPA_VERSION not set, querying latest commit SHA from kubernetes-sigs/cluster-api-provider-aws main..."
  set +o pipefail
  REF=$(curl -sL "https://api.github.com/repos/kubernetes-sigs/cluster-api-provider-aws/commits/main" | awk -F '"' '/"sha"/ { print $4; exit }')
  set -o pipefail
  if [ -z "$REF" ]; then
    echo "Failed to determine latest CAPA commit. Set CAPA_VERSION to continue." >&2
    exit 1
  fi
  echo "Using discovered CAPA ref: $REF"
fi

P_MODULE_DIR="$ROOT_DIR/cluster-api/providers/aws"
pushd "$P_MODULE_DIR"

echo "Upgrading CAPA to commit ${REF} in providers/aws"
GO111MODULE=on go get "sigs.k8s.io/cluster-api-provider-aws/v2@${REF}"
echo "Running go mod tidy"
GO111MODULE=on go mod tidy
echo "Running go mod vendor"
GO111MODULE=on go mod vendor

GO_VERSION=$(grep 'require sigs.k8s.io/cluster-api-provider-aws/v2' go.mod | awk '{print $NF}')
echo "Go generated pseudo-version: $GO_VERSION"

echo "Reading existing replace directives from cluster-api/providers/aws/go.mod"
EXISTING_REPLACES=$(grep -E '^replace ' go.mod | awk '{print $2}' || true)

if [ -n "$EXISTING_REPLACES" ]; then
  echo "Found existing replace directives for packages:"
  echo "$EXISTING_REPLACES"
  
  CAPI_CORE_GOMOD="$ROOT_DIR/cluster-api/cluster-api/go.mod"
  if [ ! -f "$CAPI_CORE_GOMOD" ]; then
    echo "Warning: cluster-api/cluster-api/go.mod not found at $CAPI_CORE_GOMOD" >&2
  else
    echo "Looking up versions in cluster-api/cluster-api/go.mod"
    while IFS= read -r pkg; do
      if [ -n "$pkg" ]; then
        VERSION=$(grep -E "^\s*${pkg}\s+" "$CAPI_CORE_GOMOD" | awk '{print $2}' | head -1 || true)
        if [ -z "$VERSION" ]; then
          VERSION=$(grep -E "^replace.*=>\s*${pkg}\s+" "$CAPI_CORE_GOMOD" | awk '{print $NF}' || true)
        fi
        
        if [ -n "$VERSION" ]; then
          echo "Applying replace: ${pkg} => ${pkg} ${VERSION}"
          go mod edit -replace="${pkg}=${pkg}@${VERSION}"
        else
          echo "Warning: Could not find version for ${pkg} in cluster-api/cluster-api/go.mod"
        fi
      fi
    done <<< "$EXISTING_REPLACES"
  fi
fi

echo "Running go mod vendor in cluster-api/providers/aws"
GO111MODULE=on go mod vendor
popd

echo "Upgrading CAPA to ${GO_VERSION} in root directory"
pushd "$ROOT_DIR"
GO111MODULE=on go get "sigs.k8s.io/cluster-api-provider-aws/v2@${REF}"

echo "Running go mod tidy"
GO111MODULE=on go mod tidy
echo "Running go mod vendor"
GO111MODULE=on go mod vendor
popd

echo "Done."
