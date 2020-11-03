#!/bin/bash
# This is the old image, see Makefile
# $ curl -O -L https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202007051540-0/x86_64/rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz
# $ mv rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz /tmp
# $ sudo gunzip /tmp/rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz

IGNITION_CONFIG="/var/lib/libvirt/images/sno.ign"
sudo cp "$1" "${IGNITION_CONFIG}"
sudo chown qemu:qemu "${IGNITION_CONFIG}"
sudo restorecon "${IGNITION_CONFIG}"

RHCOS_IMAGE="/tmp/rhcos-46.82.202008181646-0-qemu.x86_64.qcow2"
VM_NAME="sno-test"
OS_VARIANT="rhel8.1"
RAM_MB="16384"
DISK_GB="20"

virt-install \
    --connect qemu:///system \
    -n "${VM_NAME}" \
    -r "${RAM_MB}" \
    --os-variant="${OS_VARIANT}" \
    --import \
    --network=network:test-net,mac=52:54:00:ee:42:e1 \
    --graphics=none \
    --disk "size=${DISK_GB},backing_store=${RHCOS_IMAGE}" \
    --qemu-commandline="-fw_cfg name=opt/com.coreos/config,file=${IGNITION_CONFIG}"

