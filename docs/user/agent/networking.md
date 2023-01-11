# Networking configurations

This document describes the various networking configurations supported by the agent-based installer for OpenShift.
- DHCP
- Static Networking

## DHCP
The agent-based installer for Openshift can be deployed using the below 2 ways with DHCP:

### 1. install config + agent config (with only the rendezvous IP)
For this option, provide the usual `install-config.yaml` based on the cluster topology you want to create a cluster for i.e. SNO, HA or Compact. Provide the `agent-config.yaml` with as little content as below (of course, if you want to add role or rootDeviceHints config  you may)

___agent-config.yaml___
```
apiVersion: v1alpha1
kind: AgentConfig
metadata:
  name: ostest
  namespace: cluster0
rendezvousIP: <NODE_ZERO_IP>
```
With this option, when you create an agent iso, the agent-based installer for OpenShift, first reads the provided `agent-config.yaml` and  `install-config.yaml` and then programmatically creates the low level ZTP manifests such as `agent-cluster-install.yaml`, `cluster-image-set.yaml`, `cluster-deployment.yaml`, etc. intentionally skipping `nmstateconfig.yaml`

### 2. ZTP manifests without the nmstateconfig.yaml
With this approach, you can manually provide all the ZTP manifests such as `agent-cluster-install.yaml`, `cluster-image-set.yaml`, `cluster-deployment.yaml`, etc. **except** the `nmstateconfig.yaml` in the default `cluster-manifests` directory, for example, and then create the agent iso.

## Note:
To create the iso, run `openshift-install agent create image` command.

## Static Networking
TBD
