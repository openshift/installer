#!/bin/bash

set -x 

for name in ironic ironic-inspector ironic-ramdisk-logs dnsmasq httpd coreos-downloader image-customization; do
    podman ps | grep -w "$name$" && podman kill $name
    podman ps --all | grep -w "$name$" && podman rm $name -f
done
