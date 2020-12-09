#!/bin/bash
# This is the old image, see Makefile
# $ curl -O -L https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202007051540-0/x86_64/rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz
# $ mv rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz /tmp
# $ sudo gunzip /tmp/rhcos-46.82.202007051540-0-qemu.x86_64.qcow2.gz

RHCOS_ISO="/tmp/images/installer-SNO-image.iso"
VM_NAME="sno-test"
OS_VARIANT="rhel8.1"
RAM_MB="16384"
DISK_GB="40"
CPU_CORE="6"

virt-install \
    --connect qemu:///system \
    -n "${VM_NAME}" \
    -r "${RAM_MB}" \
    --vcpus "${CPU_CORE}" \
    --os-variant="${OS_VARIANT}" \
    --import \
    --network=network:test-net,mac=52:54:00:ee:42:e1 \
    --graphics=none \
    --events on_reboot=restart \
    --cdrom "${RHCOS_ISO}" \
    --disk pool=default,size="${DISK_GB}" \
    --boot hd,cdrom
