#!/bin/bash
set -euo pipefail

DATA_DIR="/usr/local/share/assisted-service"

validator_query() {
    printf '.'
    # Increase disk size requirements for InternalReleaseImage
    if [ -d /etc/assisted/extra-manifests ] && \
        grep -q "^kind:[[:space:]]*InternalReleaseImage[[:space:]]*$" /etc/assisted/extra-manifests/* 2>/dev/null; then
        printf '|.[].master.disk_size_gb += %d' 120
        printf '|.[].sno.disk_size_gb += %d' 120
    fi
}

hw_requirements() {
    # The env file is used both with Podman via --env-file and systemd via
    # EnvironmentFile, and unfortunately these have different escaping
    # requirements. We use unescaped data as that seems to work for both, even
    # though it is not valid bash.
    jq -f <(validator_query) <"${DATA_DIR}/default_hw_requirements.json" | \
        tee >(cat >&2) | \
        jq -r '"HW_VALIDATOR_REQUIREMENTS=\(.|tostring)"'
}

# Replace the final value in the env file
sed -i "s|^HW_VALIDATOR_REQUIREMENTS=.*|$(hw_requirements)|" "${DATA_DIR}/assisted-service.env"
