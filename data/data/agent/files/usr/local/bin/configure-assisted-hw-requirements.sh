#!/bin/bash
set -e

ASSISTED_SERVICE_ENV_FILE="/usr/local/share/assisted-service/assisted-service.env"
DISK_SIZE_GB=100

# Increase disk size requirement for NoRegistryClusterInstall aka OVE
if [ -f /etc/assisted/extra-manifests/internalreleaseimage.yaml ]; then
    DISK_SIZE_GB=220
fi

HW_VALIDATOR_REQUIREMENTS="[{\"version\":\"default\",\"master\":{\"cpu_cores\":4,\"ram_mib\":16384,\"disk_size_gb\":${DISK_SIZE_GB},\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":100,\"packet_loss_percentage\":0},\"arbiter\":{\"cpu_cores\":2,\"ram_mib\":8192,\"disk_size_gb\":50,\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":1000,\"packet_loss_percentage\":0},\"worker\":{\"cpu_cores\":2,\"ram_mib\":8192,\"disk_size_gb\":${DISK_SIZE_GB},\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":1000,\"packet_loss_percentage\":10},\"sno\":{\"cpu_cores\":8,\"ram_mib\":16384,\"disk_size_gb\":${DISK_SIZE_GB},\"installation_disk_speed_threshold_ms\":10}}]"

sed -i "s|^HW_VALIDATOR_REQUIREMENTS=.*|HW_VALIDATOR_REQUIREMENTS=$HW_VALIDATOR_REQUIREMENTS|" "$ASSISTED_SERVICE_ENV_FILE"