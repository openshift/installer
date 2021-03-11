#!/usr/bin/env bash

if test "x${1}" = 'x--id'
then
	GATHER_ID="${2}"
	shift 2
fi

ARTIFACTS="/tmp/artifacts-${GATHER_ID}"
mkdir -p "${ARTIFACTS}"

echo "Gathering master systemd summary ..."
LANG=POSIX systemctl list-units --state=failed >& "${ARTIFACTS}/failed-units.txt"

echo "Gathering master failed systemd unit status ..."
mkdir -p "${ARTIFACTS}/unit-status"
sed -n 's/^\* \([^ ]*\) .*/\1/p' < "${ARTIFACTS}/failed-units.txt" | while read -r UNIT
do
    systemctl status --full "${UNIT}" >& "${ARTIFACTS}/unit-status/${UNIT}.txt"
    journalctl -u "${UNIT}" > "${ARTIFACTS}/unit-status/${UNIT}.log"
done

echo "Gathering master journals ..."
mkdir -p "${ARTIFACTS}/journals"
for service in kubelet crio machine-config-daemon-host pivot openshift-azure-routes openshift-gcp-routes
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/journals/${service}.log"
done

echo "Gathering master networking ..."
mkdir -p "${ARTIFACTS}/network"
ip addr >& "${ARTIFACTS}/network/ip-addr.txt"
ip route >& "${ARTIFACTS}/network/ip-route.txt"
hostname >& "${ARTIFACTS}/network/hostname.txt"
cp -r /etc/resolv.conf "${ARTIFACTS}/network/"

echo "Gathering master containers ..."
mkdir -p "${ARTIFACTS}/containers"
for container in $(crictl ps --all --quiet)
do
    container_name=$(crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")
    crictl logs "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.log"
    crictl inspect "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.inspect"
done

podman ps --all --format "{{ .ID }} {{ .Names }}" | while read -r container_id container_name
do
    podman logs "${container_id}" >& "${ARTIFACTS}/containers/${container_name}-${container_id}.log"
    podman inspect "${container_id}" >& "${ARTIFACTS}/containers/${container_name}-${container_id}.inspect"
done

echo "Waiting for logs ..."
while wait -n; do jobs; done
