#!/bin/bash

set -euo pipefail

# Create a podman secret for the image-customization-server
base64 -w 0 /root/.docker/config.json | podman secret create pull-secret -
