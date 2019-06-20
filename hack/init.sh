#!/bin/sh

eval "$(go env | grep -e "GOHOSTOS" -e "GOHOSTARCH")"

export GOOS=${GOOS:-${GOHOSTOS}}
export GOARCH=${GOACH:-${GOHOSTARCH}}
OUT_DIR=${OUT_DIR:-_output}

function setup_env() {
  local init_source

  init_source="$( dirname "${BASH_SOURCE[0]}" )/.."
  INSTALLER_ROOT="$( absolute_path "${init_source}" )"
  export INSTALLER_ROOT
  INSTALLER_GO_PKG="github.com/openshift/installer"
  export INSTALLER_GO_PKG

  if [[ -z "${GOPATH+a}" ]]; then
    unset GOBIN
    # create a local GOPATH in _output
    GOPATH="${INSTALLER_ROOT}/${OUT_DIR}/go"
    local go_pkg_dir="${GOPATH}/src/${INSTALLER_GO_PKG}"

    mkdir -p "$(dirname "${go_pkg_dir}")"
    rm -f "${go_pkg_dir}"
    ln -s "${INSTALLER_ROOT}" "${go_pkg_dir}"

    export GOPATH
  fi

  if [[ -z "${INSTALLER_BIN_PATH+a}" ]]; then
    export INSTALLER_BIN_PATH="${INSTALLER_ROOT}/${OUT_DIR}/${GOOS}/${GOARCH}"
  fi
  mkdir -p "${INSTALLER_BIN_PATH}"
}
readonly -f setup_env

# absolute_path returns the absolute path to the directory provided
function absolute_path() {
        local relative_path="$1"
        local absolute_path

        pushd "${relative_path}" >/dev/null || exit
        relative_path="$( pwd )"
        if [[ -h "${relative_path}" ]]; then
                absolute_path="$( readlink "${relative_path}" )"
        else
                absolute_path="${relative_path}"
        fi
        popd >/dev/null || exit

	echo "${absolute_path}"
}
readonly -f absolute_path
