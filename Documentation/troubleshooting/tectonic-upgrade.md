# Troubleshooting Tectonic upgrades

This document describes how to troubleshoot issues encountered when upgrading Tectonic.

## Upgrading to Tectonic 1.7.1

To update to Tectonic 1.7.1, first update to Tectonic 1.6.7-tectonic.2. Updates to Tectonic 1.7.1 from channels previous to 1.6.7-tectonic.2 will fail.

### Switching to 1.7 channel before updating to v1.6.7_tectonic.2.

If Tectonic Console was used to switch to the 1.7.1 channel from v1.6.7-tectonic.1 or previous, first revert to the channel listed before update. Then wait for the next update check. When Tectonic Console lists the option, switch to v1.6.7-tectonic.2. Once that update is complete, use the Console to update to 1.7.1.

### Updating to 1.7 before updating to v1.6.7_tectonic.2.

If Tectonic has been updated to 1.7.1 before updating to v1.6.7_tectonic.2, reset the TPR to clear the update error.

You will receive the error:

```
Updates are not possible : Upgrade is not supported: minor version upgrade is not supported, desired: "1.7.2-tectonic.1", current: "1.6.7-tectonic.1"
```

First, use `kubectl replace` to reset to the desired version:

```
cat<<EOF | kubectl replace -f -
apiVersion: coreos.com/v1
kind: AppVersion
metadata:
name: tectonic-cluster
namespace: tectonic-system
labels:
managed-by-channel-operator: "true"
status:
currentVersion: 1.6.7-tectonic.1
paused: false
spec:
desiredVersion: 1.6.7-tectonic.1
paused: false
EOF
```
Then, use Tectonic Console to switch the channel back to `tectonic-1.6`. Click `Check for Updates`, then click `Start Upgrade`.

After upgrading to `1.6.7-tectonic.2`, switch to the `tectonic-1.7` channel and upgrade from there.
