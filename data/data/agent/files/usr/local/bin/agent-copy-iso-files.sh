#!/bin/bash

# ignition.img is a compressed cpio archive usually containing just the
# config.ign file. In case of agent-based installation, it could be 
# enriched with additional files

AGENT_FILES_TEMP="$(mktemp -d)"

cd "${AGENT_FILES_TEMP}" || { echo "Temp folder creation failed"; exit 1; }
zcat /run/media/iso/images/ignition.img | cpio -idmv

# agent-tui is required by the agent-interactive-console.service
cp agent-tui /usr/local/bin/

rm -rf "${AGENT_FILES_TEMP}"