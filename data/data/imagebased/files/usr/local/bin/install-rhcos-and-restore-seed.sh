#!/bin/bash

set -e # Halt on error

seed_image=${1:-$SEED_IMAGE}
authfile=${PULL_SECRET_FILE:-"/var/tmp/pull-secret.json"}
ibi_config=${IBI_CONFIGURATION_FILE:-"/var/tmp/ibi-configuration.json"}

# Copy the lca-cli binary to the host, pulling the seed image can sometimes fail
until podman create --authfile "${authfile}" --name lca-cli "${seed_image}" lca-cli ; do
    sleep 10
done
podman cp lca-cli:lca-cli /usr/local/bin/lca-cli
podman rm lca-cli

/usr/local/bin/lca-cli ibi -f "${ibi_config}"
