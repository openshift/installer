# Minimal ISO configurations

The agent-based installer will generate a minimal ISO at all times when the platform type is `External`. In addition, it can be generated for other platform types when configured.  This document provides comprehensive details for configuring minimal ISOs using the agent-based installer for OpenShift. It covers the minimum version requirement, configuration options, and the method of retrieving the rootfs.

## Minimum OpenShift Version
The agent-based installer for OpenShift requires a minimum OpenShift version of 4.14 to support external platforms. To support minimal ISOs for other platform types, a minimum Openshift version of 4.18 is required.

## Configuring the External Platform
The external platform can be specified using either the install-config.yaml or by utilizing ZTP (Zero Touch Provisioning) manifests in agent-cluster-install.yaml. 
When configuring the external platform, ensure that the platformName field is set to `oci`.

__install-config.yaml__
```
apiVersion: v1
baseDomain: test.metalkube.org
metadata:
  name: ostest
  namespace: cluster0
................
................
................
................
platform:
    external:
      platformName: oci
```

__agent-cluster-install.yaml__
```
kind: AgentClusterInstall
metadata:
  name: ostest
  namespace: cluster0
spec:
  external:
    platformName: oci
................
................
................
................
```

## Configuring a Minimal ISO on other platform types

For platform types other than External, the agent-based installer can be configured to generate a minimal ISO via the `minimalISO` field in agent-config.yaml. By default this field is False; when set to True the agent-based installer will generate a minimal ISO.

## Generation of Minimal ISOs
When the agent-based installer generates a minimal ISO, it may or may not generate the rootfs file explicitly depending whether the rootFS URL is specified via the `bootArtifactsBaseURL` field in `agent-config.yaml`. A minimal ISO is similar to the full ISO generated for other platforms, with the distinction that it does not contain the rootfs file within it, but rather embedded ignition in the minimal ISO has the details from where to pull the appropriate rootfs file.

When generating the minimal ISO, the agent-based installer follows these steps:

1. Deletes the rootfs image file from the RHCOS (Red Hat CoreOS) base ISO.
2. Updates the grub configuration parameter `coreos.live.rootfs_url=` with the URL for the rootfs image file location.
3. Users can specify the rootfs URL via an optional field named `bootArtifactsBaseURL` in the `agent-config.yaml`.

# Downloading Rootfs Image
When agent nodes are booted with the minimal ISO, the actual rootfs image file is dynamically downloaded into memory from the URL provided internally by the agent-based installer in the grub configuration.

## RootFS URL Configuration
When running `openshift-install agent create image`

- If the rootFS URL is specified via the `bootArtifactsBaseURL` field in `agent-config.yaml`, the agent-based installer embeds the specified URL into the grub configuration. It also generates a minimal ISO along with the rootfs.img file in the `boot-artifacts` directory. For disconnected cluster installations, ensure that the agent-installer generated rootfs image file is uploaded to the URL specified in `bootArtifactsBaseURL` before booting the nodes with the minimal ISO.

- If the rootFS URL is not specified via `bootArtifactsBaseURL` in `agent-config.yaml`, the agent-based installer embeds the default rootfs URL from the RHCOS streams file into the grub configuration. In this case, only a minimal ISO is generated. This is particularly useful for connected cluster installations, as the default rootfs URL from the RHCOS streams is readily accessible in connected environments.
