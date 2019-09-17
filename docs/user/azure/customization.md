# Azure Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, Azure-specific properties.

## Cluster-scoped properties

The following options are available when using Azure:

* `region` (required string): The Azure region where the cluster will be created.
* `baseDomainResourceGroupName` (required string): The resource group where the Azure DNS zone for the base domain is found.
* `defaultMachinePlatform` (optional object): Default [Azure-specific machine pool properties](#machine-pools) which applies to [machine pools](../customization.md#machine-pools) that do not define their own Azure-specific properties.

## Machine pools

* `osDisk` (optional object):
    * `diskSizeGB` (optional integer): The size of the disk in gigabytes (GB).
* `type` (optional string): The Azure instance type.
* `zones` (optional string slice): List of Azure availability zones that can be used (for example, `["1", "2", "3"]`).

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal Azure install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  azure:
    region: centralus
    baseDomainResourceGroupName: os4-common
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom machine pools

An example Azure install config with custom machine pools:

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
      zones:
      - "1"
      - "2"
      - "3"
  replicas: 5
metadata:
  name: test-cluster
platform:
  azure:
    region: centralus
    baseDomainResourceGroupName: os4-common
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
