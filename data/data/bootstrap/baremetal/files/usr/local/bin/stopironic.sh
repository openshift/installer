#!/bin/bash

set -x 

for name in ironic-api ironic-conductor ironic-inspector ironic-deploy-ramdisk-logs ironic-inspector-ramdisk-logs dnsmasq httpd mariadb ipa-downloader coreos-downloader; do
    podman ps | grep -w "$name$" && podman kill $name
    podman ps --all | grep -w "$name$" && podman rm $name -f
done
