# Azure Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, Azure-specific properties.

## Cluster-scoped properties

The following options are available when using Azure:

* `region` (required string): The Azure region where the cluster will be created.
* `baseDomainResourceGroupName` (required string): The resource group where the Azure DNS zone for the base domain is found.
* `defaultMachinePlatform` (optional object): Default [Azure-specific machine pool properties](#machine-pools) which applies to [machine pools](../customization.md#machine-pools) that do not define their own Azure-specific properties.
* `resourceGroupName` (optional string):  The name of an already existing resource group where the cluster should be installed. If empty, a new resource group will created for the cluster.
* `networkResourceGroupName` (optional string): The resource group where the Azure VNet is found.
* `virtualNetwork` (optional string): The name of an existing VNet where the cluster infrastructure should be provisioned.
* `controlPlaneSubnet` (optional string): An existing subnet which should be used for the cluster control plane.
* `computeSubnet` (optional string): An existing subnet which should be used by cluster nodes.
* `outboundType` (optional string):  OutboundType is a strategy for how egress from cluster is achieved. Valid values are `Loadbalancer` or `UserDefinedRouting`
    * `Loadbalancer` (default): LoadbalancerOutboundType uses Standard loadbalancer for egress from the cluster, see [docs][azure-lb-outbound]
    * `UserDefinedRouting`: UserDefinedRoutingOutboundType uses user defined routing for egress from the cluster, see [docs][azure-udr-outbound]. User defined routing for egress can only be used when deploying clusters to pre-existing virtual networks.

## Machine pools

* `osDisk` (optional object):
    * `diskSizeGB` (optional integer): The size of the disk in gigabytes (GB).
    * `diskType` (optional string): The type of disk (allowed values are: `Premium_LRS`, `Standard_LRS`, and `StandardSSD_LRS`).
* `type` (optional string): The Azure instance type.
* `zones` (optional string slice): List of Azure availability zones that can be used (for example, `["1", "2", "3"]`).

## Installing to Existing Resource Group

The installer can use an existing resource group when provisioning an OpenShift cluster. This resource group should only be used for this specific cluster and the cluster components will assume assume ownership of all resources in the resource group. Destroying the cluster using installer will delete this resource group. This resource group must be empty with no other resources when trying to use it for creating a cluster.

If you're limiting the installer's Service Principal scope to the Resource Group defined with `resourceGroupName`, you will also need to ensure proper permissions for any other resource used by the installer in your environment such as Public DNS Zone, VNet, etc.

## Installing to Existing Networks & Subnetworks

The installer can use an existing VNet and subnets when provisioning an OpenShift cluster. If one of `networkResourceGroupName`, `virtualNetwork`, `controlPlaneSubnet`, or `computeSubnet`is specified, all must be specified [(see example below)](#existing-vnet). The installer will use these existing networks when creating infrastructure such as virtual machines, load balancers, and DNS zones.

### Cluster Isolation

When pre-existing subnets are provided, the installer will not create a network security group (NSG) or alter an existing one attached to the subnet. This restriction means that no security rules are created. If multiple clusters are installed to the same VNet and isolation is desired, it must be enforced through an administrative task after the cluster is installed.

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
        diskType: Premium_LRS
  replicas: 3
compute:
- name: worker
  platform:
    azure:
      type: Standard_DS4_v2
      osDisk:
        diskSizeGB: 512
        diskType: Standard_LRS
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
    osDisk:
        diskSizeGB: 512
        diskType: Premium_LRS
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Existing Resource Group

An example Azure install config to use a pre-existing resource group:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  creationTimestamp: null
  name: test-cluster
platform:
  azure:
    baseDomainResourceGroupName: os4-common
    resourceGroupName: example-rg
    cloudName: AzurePublicCloud
    outboundType: Loadbalancer
    region: centralus
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Existing VNet

An example Azure install config to use a pre-existing VNet and subnets:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  azure:
    region: centralus
    baseDomainResourceGroupName: os4-common
    networkResourceGroupName: example_vnet_rg
    virtualNetwork: example_vnet
    controlPlaneSubnet: example_master_subnet
    computeSubnet: example_worker_subnet
    osDisk:
        diskSizeGB: 512
        diskType: Premium_LRS
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

[azure-lb-outbound]: https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-outbound-connections#lb
[azure-udr-outbound]: https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-udr-overview
