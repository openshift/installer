#!/bin/bash
set -e

systemd_service="$1"

echo "Waiting for $systemd_service to be ready"

until systemctl is-active --quiet $systemd_service; do
    printf '.'
    sleep 5
done
