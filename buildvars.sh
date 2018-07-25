#!/bin/bash -e

# Vars exported to the build info
echo OPENSHIFT_VERSION "${OPENSHIFT_VERSION}"
echo BUILD_TIME "$(date -u '+%Y-%m-%dT%H:%M:%S%z')"
