#!/bin/bash
cd /workspace

# Add the specific files
git add pkg/types/openstack/defaults/platform.go pkg/types/openstack/defaults/platform_test.go

# Show status
git status

# Commit with the provided message
git commit -m 'openstack: implement BootstrapFlavor default value inheritance

Add default value logic in SetPlatformDefaults() to inherit BootstrapFlavor
from DefaultMachinePlatform.FlavorName when not explicitly set:

1. If BootstrapFlavor is already set, preserve it
2. Otherwise, inherit from DefaultMachinePlatform.FlavorName if available
3. If neither is set, BootstrapFlavor remains empty

Add unit tests covering all inheritance scenarios.'
