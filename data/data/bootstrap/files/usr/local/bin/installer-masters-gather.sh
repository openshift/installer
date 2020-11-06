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
(
    journalctl -o json _PID=1 MESSAGE_ID=5eb03494b6584870a536b337290809b3 | jq -r .UNIT
    sed -n 's/^\* \([^ ]*\) .*/\1/p' < "${ARTIFACTS}/failed-units.txt"
) | sort -u | while read -r UNIT
do
    systemctl status --full "${UNIT}" >& "${ARTIFACTS}/unit-status/${UNIT}.txt"
    journalctl -o json-pretty -u "${UNIT}" > "${ARTIFACTS}/unit-status/${UNIT}.log.json"
done

echo "Gathering master journals ..."
mkdir -p "${ARTIFACTS}/journals"
for service in kubelet crio machine-config-daemon-host pivot openshift-azure-routes openshift-gcp-routes
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/journals/${service}.log"
done

echo "Gathering master containers ..."
mkdir -p "${ARTIFACTS}/containers"
for container in $(crictl ps --all --quiet)
do
    container_name=$(crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")
    crictl logs "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.log"
    crictl inspect "${container}" >& "${ARTIFACTS}/containers/${container_name}-${container}.inspect"
done
for container in $(podman ps --all --quiet)
do
    podman logs "${container}" >& "${ARTIFACTS}/containers/${container}.log"
    podman inspect "${container}" >& "${ARTIFACTS}/containers/${container}.inspect"
done

echo "Waiting for logs ..."
wait
