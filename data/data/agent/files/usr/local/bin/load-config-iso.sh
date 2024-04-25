#!/bin/bash
set -e

# shellcheck disable=SC1091
source issue_status.sh

status_issue="95_load_config_image"

AGENT_CONFIG_ARCHIVE_FILE="config.gz"
AGENT_CONFIG_MOUNT="/media/config-image"
CLUSTER_IMAGE_SET="/etc/assisted/manifests/cluster-image-set.yaml"
EXTRA_MANIFESTS="/etc/assisted/extra-manifests"
ASSISTED_NETWORK_DIR="/etc/assisted/network"
NM_CONNECTION_DIR="/etc/NetworkManager/system-connections"
PASSWORD_HASH="/opt/agent/tls/kubeadmin-password.hash"
OVERRIDE_PASSWORD_SET="/etc/assisted/appliance-override-password-set"

copy_archive_contents() {
    tmpdir=$(mktemp --tmpdir -d "config-image--XXXXXXXXXX")
    cp ${AGENT_CONFIG_MOUNT}/${AGENT_CONFIG_ARCHIVE_FILE} "${tmpdir}"
    gunzip -f "${tmpdir}"/${AGENT_CONFIG_ARCHIVE_FILE}
    unzipped_file="${tmpdir}/${AGENT_CONFIG_ARCHIVE_FILE%.gz}"
    filelist=$(cpio --list < "${unzipped_file}")
    echo -e "List of files to copy: ${filelist}"

    # Get the releaseImage in the archive and verify it matches the current cluster-image-set
    release_image=$(grep releaseImage ${CLUSTER_IMAGE_SET} | sed -n -e 's/^.*releaseImage: //p')
    if ! diff "${CLUSTER_IMAGE_SET}" <(cpio -icv --to-stdout "${CLUSTER_IMAGE_SET}" <"${unzipped_file}"); then
       echo "The cluster-image-set in archive does not match current release ${release_image}"
       printf '\\e{lightred}Installation cannot proceed:\\e{reset} cluster-image-set in archive does not match current release' | set_issue "${status_issue}"
       cleanup_files
       return 1
    fi
    echo "Archive on ${devname} contains release ${release_image}"

    # Get array from string
    IFS=',' read -r -a files <<< "${CONFIG_IMAGE_FILES}"

    # Copy expected files from archive, overwriting the existing file
    for file in "${files[@]}"
    do
       if [[ "${file}" = *.* || "${file}" == "/etc/issue" ]]; then
          cpio -icvdu "${file}" < "${unzipped_file}"
          if [[ -f ${file} ]]; then
             echo "Copied file ${file}"
          else
             echo "Expected file ${file} is not in archive"
             printf '\\e{lightred}Installation cannot proceed:\\e{reset} Failure copying files from config-image, expected file %s is not in archive' "${file}" | set_issue "${status_issue}"
	     return 1
          fi
       else
	  # copy all files in directory
          cpio -icvdu "${file}*" < "${unzipped_file}"

	  # directory may not contain files
          if [[ -d ${file} ]]; then
             echo "Copied files in ${file}"
	  fi
       fi
    done

    # assisted-service expects the extra-manifests dir to exist
    if [[ ! -d "${EXTRA_MANIFESTS}" ]]; then
       mkdir -p "/etc/assisted/extra-manifests"
    fi

    echo "Successfully copied contents of ${AGENT_CONFIG_ARCHIVE_FILE} on ${devname}"

    # Update core password with one from config-image if not overridden by appliance
    if [[ ! -f "${OVERRIDE_PASSWORD_SET}" ]]; then
       echo "Setting core password"
       usermod --password "$(cat ${PASSWORD_HASH})" core
    else
       echo "Setting of core password is overridden"
    fi

    # Enable any services which are not enabled by default
    declare -a servicelist=("start-cluster-installation.service")

    for service in "${servicelist[@]}"
    do
       is_enabled=$(systemctl is-enabled "$service")
       if [[ "${is_enabled}" == "disabled" ]]; then
          echo "Service ${service} is disabled, enabling it"
          systemctl --no-block --now enable "${service}"
       fi
    done

    if [[ -d "${ASSISTED_NETWORK_DIR}" ]]; then
       # Run script to generate NetworkManager keyfiles if network files exist
       /usr/local/bin/pre-network-manager-config.sh

       systemctl restart NetworkManager

       # Ensure networking is up for created files
       find ${NM_CONNECTION_DIR} -name "*.nmconnection" | while IFS= read -r nmconn_file; do
           filename=$(basename -- "$nmconn_file")
           interface="${filename%%.nmconnection*}"
           sudo nmcli conn up "${interface}"
       done
    fi

    cleanup_files
    return 0
}

cleanup_files() {

    if [[ -f "${unzipped_file}" ]]; then
       rm "${unzipped_file}"
    fi
    if [[ -d "${tmpdir}" ]]; then
       rmdir "${tmpdir}"
    fi
}

# This script will be invoked by a udev rule when it detects a device with the correct label
devname="$1"
systemd-mount --no-block --automount=yes --collect "$devname" "${AGENT_CONFIG_MOUNT}"
echo "Mounted ${devname} on ${AGENT_CONFIG_MOUNT}"

while true
do
    if [[ -f ${AGENT_CONFIG_MOUNT}/${AGENT_CONFIG_ARCHIVE_FILE} ]]; then
       # Copy contents from archive
       if copy_archive_contents; then
          break
       fi
    else
       echo "Could not find ${AGENT_CONFIG_ARCHIVE_FILE} in ${AGENT_CONFIG_MOUNT}"
    fi

    echo "Retrying to copy contents from ${AGENT_CONFIG_MOUNT}/${AGENT_CONFIG_ARCHIVE_FILE}"
    sleep 5
done

