#!/bin/sh

set -e

TARGET_OS_ARCH=$(go env GOOS)_$(go env GOARCH)
CLUSTER_API_BIN_DIR="${PWD}/cluster-api/bin/${TARGET_OS_ARCH}"
CLUSTER_API_MIRROR_DIR="${PWD}/pkg/clusterapi/mirror/"
ENVTEST_K8S_VERSION="1.30.0"
ENVTEST_ARCH=$(go env GOOS)-$(go env GOARCH)

copy_cluster_api_to_mirror() {
  mkdir -p "${CLUSTER_API_BIN_DIR}"
  mkdir -p "${CLUSTER_API_MIRROR_DIR}"

  # Clean the mirror, but preserve the README file.
  rm -rf "${CLUSTER_API_MIRROR_DIR:?}/*.zip"

  if test "${SKIP_ENVTEST}" != y; then
    sync_envtest
  fi

  # Zip every binary in the folder into a single zip file.
  zip -j1 "${CLUSTER_API_MIRROR_DIR}/cluster-api.zip" "${CLUSTER_API_BIN_DIR}"/*
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
