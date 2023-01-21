#!/bin/bash

echo "agent-interactive-console start"
systemctl status pre-network-manager-config.service

# TODO: Execute tui and remove sleep
echo "sleeping 60 seconds"
sleep 60

echo "agent-interactive-console end"