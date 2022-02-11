#!/bin/bash

HOST=$(ip -4 -j address | jq -r '[.[].addr_info]|flatten(1)|map(select(.scope=="global"))|.[0].local')
echo Using hostname ${HOST} 1>&2

podman play kube --configmap <(sed -e "/SERVICE_BASE_URL/ s/127\.0\.0\.1/${HOST}/" /usr/local/share/assisted-service/configmap.yml) /usr/local/share/assisted-service/pod.yml
