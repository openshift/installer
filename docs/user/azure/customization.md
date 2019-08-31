# Azure Platform Customization

The following options are available when using Azure:

- `machines.platform.azure.type` - the VM instance type
- `machines.platform.azure.osDisk.diskSizeGB` - The Azure OS disk size in Gigabytes
- `platform.azure.region` - the Azure region (location) that the installer will use when creating resource group and resources
- `platform.azure.baseDomainResourceGroupName` - the Azure Resource Group that has the public DNS zone for base domain

## Examples

An example `install-config.yaml` is shown below. This configuration has been modified to show the customization that is possible via the install config.

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    azure:
      type: Standard_DS4_v2
      osDisk:
        diskSizeGB: 512
  replicas: 3
compute:
- name: worker
  platform:
    azure:
      type: Standard_DS4_v2
      osDisk:
        diskSizeGB: 512
  replicas: 5
metadata:
  name: test-cluster
networking:
  clusterNetworks:
  - cidr: 10.128.0.0/14
    hostSubnetLength: 9
  machineCIDR: 10.0.0.0/16
  serviceCIDR: 172.30.0.0/16
  type: OpenShiftSDN
platform:
  azure:
    region: centralus
    baseDomainResourceGroupName: os4-common
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
