#!/bin/bash

set -e

AGENT_CONFIG_IMAGE_LABEL="agent_configimage"
AGENT_CONFIG_ARCHIVE_FILE="config.gz"
AGENT_CONFIG_MOUNT="/mnt/config_image"
RENDEZVOUS_HOST_ENV="/etc/assisted/rendezvous-host.env"
CLUSTER_IMAGE_SET="/etc/assisted/manifests/cluster-image-set.yaml"

copy_archive_contents() {
    
    # Mount device
    mkdir -p ${AGENT_CONFIG_MOUNT}
    mount ${device} ${AGENT_CONFIG_MOUNT}

    if [[ ! -f ${AGENT_CONFIG_MOUNT}/${AGENT_CONFIG_ARCHIVE_FILE} ]]; then
       echo "Could not find file ${AGENT_CONFIG_ARCHIVE_FILE} on ${device}"
       cleanup_files
       return 1 
    fi

    cp ${AGENT_CONFIG_MOUNT}/${AGENT_CONFIG_ARCHIVE_FILE} /tmp
    gunzip -f /tmp/${AGENT_CONFIG_ARCHIVE_FILE}
    unzipped_file=`echo /tmp/${AGENT_CONFIG_ARCHIVE_FILE} | cut -d'.' -f1`

    # Get the releaseImage in the archive and verify it matches the current cluster-image-set
    release_image=$(cat ${CLUSTER_IMAGE_SET} | grep releaseImage | sed -n -e 's/^.*releaseImage: //p')
    arch_release_image=$(cpio -icv --to-stdout ${CLUSTER_IMAGE_SET} < ${unzipped_file} | grep releaseImage | sed -n -e 's/^.*releaseImage: //p')
    if [[ ${release_image} != ${arch_release_image} ]]; then
       echo "The release $arch_release_image in archive does not match the current release ${release_image}"
       cleanup_files
       return 1 
    fi
    echo "Archive on ${device} contains release ${arch_release_image}"

    # Copy all files from archive
    # The rendezvousIP file must be copied last as it triggers set-node-zero.sh to continue configuration
    cpio -icvd --nonmatching ${AGENT_RENDEZVOUS_IP_FILE} < ${unzipped_file}
    cpio -icvd ${RENDEZVOUS_HOST_ENV} < ${unzipped_file}

    echo "Successfully copied contents of ${AGENT_CONFIG_ARCHIVE_FILE} on ${device}"

    cleanup_files
    return 0
}

cleanup_files() {

    if [[ -d  ${AGENT_CONFIG_MOUNT} ]]; then
       umount ${AGENT_CONFIG_MOUNT}
       rmdir ${AGENT_CONFIG_MOUNT}
    fi
    if [[ -f ${unzipped_file} ]]; then
       rm ${unzipped_file}
    fi
}

while true 
do
    # Look for devices matching config image label 
    device=$(lsblk -p -o NAME,LABEL | grep ${AGENT_CONFIG_IMAGE_LABEL} | cut -d ' ' -f1)

    if [[ ! -z ${device} ]]; then
       echo "Found ${device} matching label ${AGENT_CONFIG_IMAGE_LABEL}"

       if copy_archive_contents; then
          break
       fi
    fi

    echo "Retrying to find device matching label ${AGENT_CONFIG_IMAGE_LABEL}"
    sleep 5
done

