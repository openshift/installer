#!/usr/bin/env bash

ARTIFACTS="${1:-/tmp/artifacts}"
mkdir -p "${ARTIFACTS}"

echo "Gathering master journals ..."
mkdir -p "${ARTIFACTS}/journals"
for service in kubelet crio
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/journals/${service}.log"
done

echo "Gathering master containers ..."
mkdir -p "${ARTIFACTS}/containers"
for container in $(crictl ps --all --quiet)
do
    container_name=$(crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")
    crictl logs "${container}" >& "${ARTIFACTS}/containers/${container_name}.log"
    crictl inspect "${container}" >& "${ARTIFACTS}/containers/${container_name}.inspect"
done
for container in $(podman ps --all --quiet)
do
    podman logs "${container}" >& "${ARTIFACTS}/containers/${container}.log"
    podman inspect "${container}" >& "${ARTIFACTS}/containers/${container}.inspect"
done

echo "Waiting for logs ..."
wait
