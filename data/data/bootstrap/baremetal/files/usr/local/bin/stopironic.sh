#!/bin/bash

set -x 

for name in ironic ironic-inspector httpd; do
    podman ps | grep -w "$name$" && podman kill $name
    podman ps --all | grep -w "$name$" && podman rm $name -f
done
