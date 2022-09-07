#!/usr/bin/env bash
set -euo pipefail

/usr/local/bin/release-image-download.sh

# shellcheck disable=SC1091
. /usr/local/bin/release-image.sh

# Store images in the environment file used by services and passed to assisted-service
# The agent image will be also retrieved when its script is run
cat <<EOF >/usr/local/share/assisted-service/agent-images.env
SERVICE_IMAGE=$(image_for agent-installer-api-server)
CONTROLLER_IMAGE=$(image_for agent-installer-csr-approver)
INSTALLER_IMAGE=$(image_for agent-installer-orchestrator)
AGENT_DOCKER_IMAGE=$(image_for agent-installer-node-agent)
EOF
