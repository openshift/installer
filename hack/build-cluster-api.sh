#!/bin/sh

set -e

TARGET_OS_ARCH=$(go env GOOS)_$(go env GOARCH)
CLUSTER_API_BIN_DIR="${PWD}/cluster-api/bin/${TARGET_OS_ARCH}"
CLUSTER_API_MIRROR_DIR="${PWD}/pkg/clusterapi/mirror/"
ENVTEST_K8S_VERSION="1.35.0"
ENVTEST_ARCH=$(go env GOOS)-$(go env GOARCH)

copy_cluster_api_to_mirror() {
  mkdir -p "${CLUSTER_API_BIN_DIR}"
  mkdir -p "${CLUSTER_API_MIRROR_DIR}"

  # Clean the mirror, but preserve the README file.
  rm -f "${CLUSTER_API_MIRROR_DIR:?}"/*.zst "${CLUSTER_API_MIRROR_DIR:?}"/dict "${CLUSTER_API_MIRROR_DIR:?}"/*.zip

  if test "${SKIP_ENVTEST}" != y; then
    sync_envtest
  fi

  # Train a zstd dictionary on all binaries for cross-binary deduplication,
  # then compress each binary individually with the shared dictionary.
  zstd -q --train "${CLUSTER_API_BIN_DIR}"/* -o "${CLUSTER_API_MIRROR_DIR}/dict" --maxdict=262144
  for bin in "${CLUSTER_API_BIN_DIR}"/*; do
    name=$(basename "$bin")
    zstd -9 -D "${CLUSTER_API_MIRROR_DIR}/dict" "$bin" -o "${CLUSTER_API_MIRROR_DIR}/${name}.zst"
  done
}

sync_envtest() {
  if [ -f "${CLUSTER_API_BIN_DIR}/kube-apiserver" ]; then
    if [ "$(go env GOOS)" != "$(go env GOHOSTOS)" ] || [ "$(go env GOARCH)" != "$(go env GOHOSTARCH)" ]; then
      echo "Found cross-compiled artifact: skipping envtest binaries version check"
      return
    fi
    version=$( ("${CLUSTER_API_BIN_DIR}/kube-apiserver" --version || echo "v0.0.0") | sed 's/Kubernetes //' )
    echo "Found envtest binaries with version: ${version}"
    if printf '%s\n%s' v${ENVTEST_K8S_VERSION} "${version}" | sort -V -C; then
      return
    fi
  fi

  bucket="https://github.com/kubernetes-sigs/controller-tools/releases/download/envtest-v${ENVTEST_K8S_VERSION}"
  tar_file="envtest-v${ENVTEST_K8S_VERSION}-${ENVTEST_ARCH}.tar.gz"
  dst="${CLUSTER_API_BIN_DIR}/${tar_file}"
  if ! [ -f "${CLUSTER_API_BIN_DIR}/${tar_file}" ]; then
    echo "Downloading envtest binaries"
    curl -fL "${bucket}/${tar_file}" -o "${dst}"
  fi
  tar -C "${CLUSTER_API_BIN_DIR}" -xzf "${dst}" --strip-components=2
  rm "${dst}" # Remove the tar file.
  rm "${CLUSTER_API_BIN_DIR}/kubectl" # Remove kubectl since we don't need it.
}
