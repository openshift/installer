#!/bin/bash

set -e # Halt on error

seed_image=${1:-$SEED_IMAGE}
authfile=${AUTH_FILE:-"/var/tmp/backup-secret.json"}
ibi_config=${IBI_CONFIGURATION_FILE:-"/var/tmp/ibi-configuration.json"}

# Copy the lca-cli binary to the host
podman create --authfile "${authfile}" --name lca-cli "${seed_image}" lca-cli
podman cp lca-cli:lca-cli /usr/local/bin/lca-cli
podman rm lca-cli

/usr/local/bin/lca-cli ibi -f "${ibi_config}"
