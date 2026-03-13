#!/bin/bash
set -e

ASSISTED_SERVICE_ENV_FILE="/usr/local/share/assisted-service/assisted-service.env"

# default
MASTER_DISK_SIZE=100
SNO_DISK_SIZE=100

# Increase disk size requirements for NoRegistryClusterInstall
# Check if any file in extra-manifests contains kind: InternalReleaseImage
if [ -d /etc/assisted/extra-manifests ] && \
   grep -q "^kind:[[:space:]]*InternalReleaseImage[[:space:]]*$" /etc/assisted/extra-manifests/* 2>/dev/null; then
    MASTER_DISK_SIZE=220
    SNO_DISK_SIZE=220
fi

# Build requirements with variables
HW_VALIDATOR_REQUIREMENTS="[{\"version\":\"default\",\"master\":{\"cpu_cores\":4,\"ram_mib\":16384,\"disk_size_gb\":${MASTER_DISK_SIZE},\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":100,\"packet_loss_percentage\":0},\"arbiter\":{\"cpu_cores\":2,\"ram_mib\":8192,\"disk_size_gb\":50,\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":1000,\"packet_loss_percentage\":0},\"worker\":{\"cpu_cores\":2,\"ram_mib\":8192,\"disk_size_gb\":100,\"installation_disk_speed_threshold_ms\":10,\"network_latency_threshold_ms\":1000,\"packet_loss_percentage\":10},\"sno\":{\"cpu_cores\":8,\"ram_mib\":16384,\"disk_size_gb\":${SNO_DISK_SIZE},\"installation_disk_speed_threshold_ms\":10}}]"

# Replace the final value in the env file
sed -i "s|^HW_VALIDATOR_REQUIREMENTS=.*|HW_VALIDATOR_REQUIREMENTS=$HW_VALIDATOR_REQUIREMENTS|" "$ASSISTED_SERVICE_ENV_FILE"