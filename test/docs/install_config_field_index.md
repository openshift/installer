# Install Config Field Index

This document maps install-config YAML fields to the features tracked in the
[Platform Feature Matrix](platform_feature_matrix.md). Use it to understand
which fields drive which features on each platform and where field-level test
coverage exists.

The index is organized by platform. Platform-agnostic fields that apply to all
platforms are listed first. Each platform section covers both the
`.platform.<name>` fields and the machine pool fields under
`.controlPlane.platform.<name>` and `.compute[].platform.<name>`.

## How to Read This Document

Each table uses five columns:

| Column | Description |
|--------|-------------|
| JSON Path | Dot-notation path in install-config YAML |
| Field | Human-readable name |
| Type | Go type (what the field accepts) |
| Feature(s) | Row(s) in the [Platform Feature Matrix](platform_feature_matrix.md) this field relates to |
| Test Coverage | AT (auto-tested), MT (manual), AT/MT (both), NT (not tested), NA (not applicable) |

JSON paths use `.platform.<name>.` for platform-level fields and
`.controlPlane.platform.<name>.` or `.compute[].platform.<name>.` for machine
pool fields. When a machine pool field applies to both control plane and
compute, the path is shown as `.<pool>.platform.<name>.` where `<pool>` is
either location.

## Platform-Agnostic Fields

These fields appear at the top level of install-config and apply regardless of
platform.

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.baseDomain` | Base Domain | `string` | IPI Install, UPI Install | AT |
| `.publish` | Publish Strategy | `External` / `Internal` / `Mixed` | Private Cluster | AT |
| `.fips` | FIPS Mode | `bool` | FIPS Mode | NT |
| `.proxy.httpProxy` | HTTP Proxy | `string` | Proxy Support | NT |
| `.proxy.httpsProxy` | HTTPS Proxy | `string` | Proxy Support | NT |
| `.proxy.noProxy` | No Proxy | `string` | Proxy Support | NT |
| `.networking.networkType` | Network Type | `string` | IPI Install | AT |
| `.networking.machineNetwork[].cidr` | Machine Network CIDR | `IPNet` | Dual Stack | AT |
| `.networking.clusterNetwork[].cidr` | Cluster Network CIDR | `IPNet` | Dual Stack | AT |
| `.networking.serviceNetwork[]` | Service Network CIDR | `IPNet` | Dual Stack | AT |
| `.networking.clusterNetworkMTU` | Cluster Network MTU | `uint32` | IPI Install | NT |
| `.controlPlane.replicas` | Control Plane Replicas | `int64` | IPI Install, Single Node | AT |
| `.controlPlane.architecture` | CPU Architecture | `string` | Default Machine Types | AT |
| `.compute[].replicas` | Compute Replicas | `int64` | IPI Install | AT |
| `.credentialsMode` | Credentials Mode | `Mint` / `Passthrough` / `Manual` | IPI Install | AT |
| `.featureSet` | Feature Set | `string` | Feature Gates | AT |
| `.featureGates` | Feature Gates | `[]string` | Feature Gates | AT |
| `.capabilities` | Capabilities | `object` | IPI Install | NT |
| `.additionalTrustBundle` | Additional Trust Bundle | `string` | Proxy Support | NT |
| `.sshKey` | SSH Key | `string` | IPI Install | AT |
| `.pullSecret` | Pull Secret | `string` | IPI Install | AT |
| `.cpuPartitioningMode` | CPU Partitioning | `None` / `AllNodes` | IPI Install | NT |

## AWS

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.aws.region` | Region | `string` | IPI Install | AT |
| `.platform.aws.vpc` | VPC | `object` | Existing VPC | AT |
| `.platform.aws.vpc.subnets` | VPC Subnets | `[]Subnet` | Existing VPC | AT |
| `.platform.aws.hostedZone` | Hosted Zone | `string` | User-Provisioned DNS | AT |
| `.platform.aws.hostedZoneRole` | Hosted Zone Role | `string` | User-Provisioned DNS | AT |
| `.platform.aws.userTags` | User Tags | `map[string]string` | User Tags / Labels | AT |
| `.platform.aws.serviceEndpoints` | Service Endpoints | `[]ServiceEndpoint` | Service Endpoints / PSC | AT |
| `.platform.aws.lbType` | Load Balancer Type | `Classic` / `NLB` | Load Balancer Config | NT |
| `.platform.aws.userProvisionedDNS` | User-Provisioned DNS | `Enabled` / `Disabled` | User-Provisioned DNS | AT |
| `.platform.aws.ipFamily` | IP Family | `IPv4` / `DualStack` | Dual Stack | AT |
| `.platform.aws.publicIpv4Pool` | Public IPv4 Pool | `string` | IPI Install | NT |
| `.platform.aws.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.aws.type` | Instance Type | `string` | Default Machine Types | AT |
| `.<pool>.platform.aws.zones` | Availability Zones | `[]string` | Failure Domains / Zones | AT |
| `.<pool>.platform.aws.amiID` | AMI ID | `string` | Custom OS Image | NT |
| `.<pool>.platform.aws.rootVolume.type` | Root Volume Type | `string` | Custom Disk Types | AT |
| `.<pool>.platform.aws.rootVolume.size` | Root Volume Size | `int` | Custom Disk Types | AT |
| `.<pool>.platform.aws.rootVolume.iops` | Root Volume IOPS | `int` | Custom Disk Types | AT |
| `.<pool>.platform.aws.rootVolume.throughput` | Root Volume Throughput | `int32` | Custom Disk Types | AT |
| `.<pool>.platform.aws.rootVolume.kmsKeyARN` | KMS Key ARN | `string` | KMS / Disk Encryption | AT |
| `.<pool>.platform.aws.iamRole` | IAM Role | `string` | IPI Install | AT |
| `.<pool>.platform.aws.iamProfile` | IAM Instance Profile | `string` | IPI Install | NT |
| `.<pool>.platform.aws.additionalSecurityGroupIDs` | Additional Security Groups | `[]string` | Existing VPC | AT |
| `.<pool>.platform.aws.cpuOptions` | CPU Options | `object` | Confidential Compute | AT |
| `.<pool>.platform.aws.hostPlacement` | Host Placement | `object` | Dedicated Hosts | AT |
| `.<pool>.platform.aws.metadataService` | Metadata Service | `object` | IPI Install | NT |

## Azure

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.azure.region` | Region | `string` | IPI Install | AT |
| `.platform.azure.baseDomainResourceGroupName` | DNS Resource Group | `string` | IPI Install | AT |
| `.platform.azure.networkResourceGroupName` | Network Resource Group | `string` | Existing VPC / VNet | AT |
| `.platform.azure.virtualNetwork` | Virtual Network | `string` | Existing VPC / VNet | AT |
| `.platform.azure.subnets` | Subnets | `[]SubnetSpec` | Existing VPC / VNet | AT |
| `.platform.azure.resourceGroupName` | Resource Group | `string` | IPI Install | AT |
| `.platform.azure.cloudName` | Cloud Name | `string` | Sovereign Cloud Support | NT |
| `.platform.azure.outboundType` | Outbound Type | `string` | IPI Install | AT |
| `.platform.azure.userTags` | User Tags | `map[string]string` | User Tags / Labels | AT |
| `.platform.azure.customerManagedKey` | Customer Managed Key | `object` | KMS / Disk Encryption | AT |
| `.platform.azure.userProvisionedDNS` | User-Provisioned DNS | `Enabled` / `Disabled` | User-Provisioned DNS | AT |
| `.platform.azure.ipFamily` | IP Family | `IPv4` / `DualStack` | Dual Stack | AT |
| `.platform.azure.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.azure.type` | VM Size | `string` | Default Machine Types | AT |
| `.<pool>.platform.azure.zones` | Availability Zones | `[]string` | Failure Domains / Zones | AT |
| `.<pool>.platform.azure.osDisk.diskSizeGB` | OS Disk Size | `int32` | Custom Disk Types | AT |
| `.<pool>.platform.azure.osDisk.diskType` | OS Disk Type | `string` | Custom Disk Types | AT |
| `.<pool>.platform.azure.osDisk.diskEncryptionSet` | Disk Encryption Set | `object` | KMS / Disk Encryption | AT |
| `.<pool>.platform.azure.osDisk.securityProfile` | Disk Security Profile | `object` | Confidential Compute | AT |
| `.<pool>.platform.azure.encryptionAtHost` | Encryption at Host | `bool` | KMS / Disk Encryption | AT |
| `.<pool>.platform.azure.osImage` | Custom OS Image | `object` | Custom OS Image | NT |
| `.<pool>.platform.azure.settings` | Security Settings | `object` | Secure Boot / UEFI, Confidential Compute | AT |
| `.<pool>.platform.azure.ultraSSDCapability` | Ultra SSD | `Enabled` / `Disabled` | Custom Disk Types | AT |
| `.<pool>.platform.azure.vmNetworkingType` | Accelerated Networking | `string` | IPI Install | NT |
| `.<pool>.platform.azure.identity` | VM Identity | `object` | IPI Install | NT |
| `.<pool>.platform.azure.dataDisks` | Data Disks | `[]DataDisk` | Custom Disk Types | NT |
| `.<pool>.platform.azure.bootDiagnostics` | Boot Diagnostics | `object` | IPI Install | NT |

## GCP

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.gcp.projectID` | Project ID | `string` | IPI Install, Sovereign Cloud Support | AT |
| `.platform.gcp.region` | Region | `string` | IPI Install, Sovereign Cloud Support | AT |
| `.platform.gcp.network` | VPC Network | `string` | Existing VPC | AT |
| `.platform.gcp.networkProjectID` | Network Project ID | `string` | Shared VPC / XPN | AT |
| `.platform.gcp.controlPlaneSubnet` | Control Plane Subnet | `string` | Existing VPC | AT |
| `.platform.gcp.computeSubnet` | Compute Subnet | `string` | Existing VPC | AT |
| `.platform.gcp.userLabels` | User Labels | `[]UserLabel` | User Tags / Labels | AT |
| `.platform.gcp.userTags` | User Tags | `[]UserTag` | User Tags / Labels | AT |
| `.platform.gcp.userProvisionedDNS` | User-Provisioned DNS | `Enabled` / `Disabled` | User-Provisioned DNS | AT |
| `.platform.gcp.endpoint` | PSC Endpoint | `object` | Service Endpoints / PSC | AT |
| `.platform.gcp.dns.privateZone` | Private DNS Zone | `object` | Shared VPC / XPN | AT |
| `.platform.gcp.firewallRulesManagement` | Firewall Rules | `Managed` / `Unmanaged` | IPI Install | NT |
| `.platform.gcp.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.gcp.type` | Instance Type | `string` | Default Machine Types | AT |
| `.<pool>.platform.gcp.zones` | Availability Zones | `[]string` | Failure Domains / Zones | AT |
| `.<pool>.platform.gcp.osDisk.diskType` | Disk Type | `string` | Custom Disk Types | AT |
| `.<pool>.platform.gcp.osDisk.DiskSizeGB` | Disk Size | `int64` | Custom Disk Types | AT |
| `.<pool>.platform.gcp.osDisk.encryptionKey` | KMS Encryption Key | `object` | KMS / Disk Encryption | AT |
| `.<pool>.platform.gcp.osImage` | Custom OS Image | `object` | Custom OS Image | AT |
| `.<pool>.platform.gcp.tags` | Network Tags | `[]string` | IPI Install | NT |
| `.<pool>.platform.gcp.secureBoot` | Secure Boot | `Enabled` / `Disabled` | Secure Boot / UEFI | NT |
| `.<pool>.platform.gcp.confidentialCompute` | Confidential Compute | `string` | Confidential Compute | AT |
| `.<pool>.platform.gcp.onHostMaintenance` | On Host Maintenance | `Migrate` / `Terminate` | Confidential Compute | AT |
| `.<pool>.platform.gcp.serviceAccount` | Service Account | `string` | Sovereign Cloud Support | AT |

## GCD (Google Cloud Dedicated)

GCD uses the same Go types as GCP (`pkg/types/gcp/`). The fields below
highlight sovereign-specific behavior: different defaults, additional
requirements, or restrictions that do not apply to public GCP. See the
[GCP section](#gcp) for the full field list.

### Platform fields with sovereign-specific behavior

| JSON Path | Field | Sovereign Behavior | Feature(s) | Test Coverage |
|-----------|-------|--------------------|------------|---------------|
| `.platform.gcp.projectID` | Project ID | Must be domain-scoped (contains `:`) | Sovereign Cloud Support | AT/MT |
| `.platform.gcp.region` | Region | Must have `u-` prefix | Sovereign Cloud Support | AT/MT |
| `.platform.gcp.endpoint` | PSC Endpoint | Forbidden (rejected by validation) | Service Endpoints / PSC | AT |

### Machine pool fields with sovereign-specific behavior

| JSON Path | Field | Sovereign Behavior | Feature(s) | Test Coverage |
|-----------|-------|--------------------|------------|---------------|
| `.<pool>.platform.gcp.type` | Instance Type | Defaults to `c3-standard-4` (C3 only) | Default Machine Types | AT/MT |
| `.<pool>.platform.gcp.osDisk.diskType` | Disk Type | Defaults to `hyperdisk-balanced` | Custom Disk Types | AT |
| `.<pool>.platform.gcp.osImage` | Custom OS Image | Required (name + project must be set) | Custom OS Image | AT/MT |
| `.<pool>.platform.gcp.serviceAccount` | Service Account | Domain-scoped email format | Sovereign Cloud Support | AT |

## Bare Metal

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.baremetal.hosts` | Host Definitions | `[]Host` | IPI Install, Agent-based Install | AT |
| `.platform.baremetal.apiVIPs` | API VIPs | `[]string` | IPI Install | AT |
| `.platform.baremetal.ingressVIPs` | Ingress VIPs | `[]string` | IPI Install | AT |
| `.platform.baremetal.provisioningNetwork` | Provisioning Network | `string` | IPI Install | AT |
| `.platform.baremetal.clusterProvisioningIP` | Cluster Provisioning IP | `string` | IPI Install | AT |
| `.platform.baremetal.bootstrapProvisioningIP` | Bootstrap Provisioning IP | `string` | IPI Install | AT |
| `.platform.baremetal.externalBridge` | External Bridge | `string` | IPI Install | AT |
| `.platform.baremetal.provisioningBridge` | Provisioning Bridge | `string` | IPI Install | NT |
| `.platform.baremetal.provisioningNetworkCIDR` | Provisioning CIDR | `IPNet` | IPI Install | AT |
| `.platform.baremetal.provisioningDHCPRange` | Provisioning DHCP Range | `string` | IPI Install | NT |
| `.platform.baremetal.provisioningNetworkInterface` | Provisioning NIC | `string` | IPI Install | NT |
| `.platform.baremetal.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.baremetal.dnsRecordsType` | DNS Records Type | `Internal` / `External` | IPI Install | NT |
| `.platform.baremetal.additionalNTPServers` | NTP Servers | `[]string` | IPI Install | NT |
| `.platform.baremetal.bmcVerifyCA` | BMC CA Verify | `string` | IPI Install | NT |
| `.platform.baremetal.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | IPI Install | NA |

Bare metal MachinePool has no platform-specific fields.

## IBM Cloud

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.ibmcloud.region` | Region | `string` | IPI Install | AT |
| `.platform.ibmcloud.resourceGroupName` | Resource Group | `string` | IPI Install | AT |
| `.platform.ibmcloud.networkResourceGroupName` | Network Resource Group | `string` | Existing VPC | AT |
| `.platform.ibmcloud.vpcName` | VPC Name | `string` | Existing VPC | AT |
| `.platform.ibmcloud.controlPlaneSubnets` | Control Plane Subnets | `[]string` | Existing VPC | AT |
| `.platform.ibmcloud.computeSubnets` | Compute Subnets | `[]string` | Existing VPC | AT |
| `.platform.ibmcloud.serviceEndpoints` | Service Endpoints | `[]ServiceEndpoint` | Service Endpoints / PSC | AT |
| `.platform.ibmcloud.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.ibmcloud.type` | Machine Profile | `string` | Default Machine Types | AT |
| `.<pool>.platform.ibmcloud.zones` | Availability Zones | `[]string` | IPI Install | AT |
| `.<pool>.platform.ibmcloud.bootVolume.encryptionKey` | Boot Encryption Key | `string` | KMS / Disk Encryption | AT |
| `.<pool>.platform.ibmcloud.dedicatedHosts` | Dedicated Hosts | `[]DedicatedHost` | Dedicated Hosts | AT |

## Nutanix

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.nutanix.prismCentral` | Prism Central | `object` | IPI Install | AT |
| `.platform.nutanix.prismElements` | Prism Elements | `[]PrismElement` | IPI Install | AT |
| `.platform.nutanix.subnetUUIDs` | Subnet UUIDs | `[]string` | Existing VPC | AT |
| `.platform.nutanix.apiVIPs` | API VIPs | `[]string` | IPI Install | AT |
| `.platform.nutanix.ingressVIPs` | Ingress VIPs | `[]string` | IPI Install | AT |
| `.platform.nutanix.failureDomains` | Failure Domains | `[]FailureDomain` | Failure Domains / Zones | AT |
| `.platform.nutanix.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.nutanix.dnsRecordsType` | DNS Records Type | `Internal` / `External` | IPI Install | NT |
| `.platform.nutanix.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | NT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.nutanix.cpus` | vCPU Count | `int64` | Default Machine Types | NT |
| `.<pool>.platform.nutanix.coresPerSocket` | Cores per Socket | `int64` | Default Machine Types | NT |
| `.<pool>.platform.nutanix.memoryMiB` | Memory (MiB) | `int64` | Default Machine Types | NT |
| `.<pool>.platform.nutanix.osDisk.diskSizeGiB` | OS Disk Size (GiB) | `int64` | Custom Disk Types | NT |
| `.<pool>.platform.nutanix.bootType` | Boot Type | `Legacy` / `UEFI` / `SecureBoot` | Secure Boot / UEFI | AT |
| `.<pool>.platform.nutanix.categories` | Prism Categories | `[]Category` | User Tags / Labels | NT |
| `.<pool>.platform.nutanix.project` | Prism Project | `object` | IPI Install | NT |
| `.<pool>.platform.nutanix.failureDomains` | Failure Domains | `[]string` | Failure Domains / Zones | AT |
| `.<pool>.platform.nutanix.gpus` | GPU Devices | `[]GPU` | IPI Install | NT |
| `.<pool>.platform.nutanix.dataDisks` | Data Disks | `[]DataDisk` | Custom Disk Types | NT |

## OpenStack

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.openstack.cloud` | Cloud Name | `string` | IPI Install | AT |
| `.platform.openstack.externalNetwork` | External Network | `string` | IPI Install | AT |
| `.platform.openstack.apiFloatingIP` | API Floating IP | `string` | IPI Install | AT |
| `.platform.openstack.ingressFloatingIP` | Ingress Floating IP | `string` | IPI Install | AT |
| `.platform.openstack.externalDNS` | External DNS Servers | `[]string` | IPI Install | NT |
| `.platform.openstack.clusterOSImage` | Custom OS Image | `string` | Custom OS Image | NT |
| `.platform.openstack.clusterOSImageProperties` | OS Image Properties | `map[string]string` | Custom OS Image | NT |
| `.platform.openstack.apiVIPs` | API VIPs | `[]string` | IPI Install | AT |
| `.platform.openstack.ingressVIPs` | Ingress VIPs | `[]string` | IPI Install | AT |
| `.platform.openstack.controlPlanePort` | Control Plane Port | `object` | Existing VPC | AT |
| `.platform.openstack.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.openstack.dnsRecordsType` | DNS Records Type | `Internal` / `External` | IPI Install | NT |
| `.platform.openstack.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.openstack.type` | Flavor | `string` | Default Machine Types | AT |
| `.<pool>.platform.openstack.zones` | Availability Zones | `[]string` | IPI Install | AT |
| `.<pool>.platform.openstack.rootVolume` | Root Volume | `object` | Custom Disk Types | NT |
| `.<pool>.platform.openstack.additionalNetworkIDs` | Additional Networks | `[]string` | Existing VPC | AT |
| `.<pool>.platform.openstack.additionalSecurityGroupIDs` | Additional Security Groups | `[]string` | Existing VPC | AT |
| `.<pool>.platform.openstack.serverGroupPolicy` | Server Group Policy | `string` | IPI Install | NT |

## PowerVC

PowerVC shares much of its type structure with OpenStack. It uses the same
CAPI provider.

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.powervc.cloud` | Cloud Name | `string` | IPI Install | AT |
| `.platform.powervc.externalNetwork` | External Network | `string` | IPI Install | NT |
| `.platform.powervc.apiFloatingIP` | API Floating IP | `string` | IPI Install | NT |
| `.platform.powervc.ingressFloatingIP` | Ingress Floating IP | `string` | IPI Install | NT |
| `.platform.powervc.externalDNS` | External DNS Servers | `[]string` | IPI Install | NT |
| `.platform.powervc.clusterOSImage` | Custom OS Image | `string` | Custom OS Image | NT |
| `.platform.powervc.apiVIPs` | API VIPs | `[]string` | IPI Install | NT |
| `.platform.powervc.ingressVIPs` | Ingress VIPs | `[]string` | IPI Install | NT |
| `.platform.powervc.controlPlanePort` | Control Plane Port | `object` | Existing VPC | NT |
| `.platform.powervc.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.powervc.imageName` | Image Name | `string` | Custom OS Image | NT |
| `.platform.powervc.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | NT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.powervc.type` | Flavor | `string` | Default Machine Types | NT |
| `.<pool>.platform.powervc.zones` | Availability Zones | `[]string` | IPI Install | NT |
| `.<pool>.platform.powervc.rootVolume` | Root Volume | `object` | Custom Disk Types | NT |
| `.<pool>.platform.powervc.additionalNetworkIDs` | Additional Networks | `[]string` | Existing VPC | NT |
| `.<pool>.platform.powervc.additionalSecurityGroupIDs` | Additional Security Groups | `[]string` | Existing VPC | NT |
| `.<pool>.platform.powervc.serverGroupPolicy` | Server Group Policy | `string` | IPI Install | NT |

## Power VS

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.powervs.powervsResourceGroup` | Resource Group | `string` | IPI Install | AT |
| `.platform.powervs.region` | Region | `string` | IPI Install | AT |
| `.platform.powervs.zone` | Zone | `string` | IPI Install | AT |
| `.platform.powervs.vpcRegion` | VPC Region | `string` | IPI Install | AT |
| `.platform.powervs.vpcName` | VPC Name | `string` | Existing VPC | AT |
| `.platform.powervs.vpcSubnets` | VPC Subnets | `[]string` | Existing VPC | AT |
| `.platform.powervs.clusterOSImage` | Custom OS Image | `string` | Custom OS Image | NT |
| `.platform.powervs.serviceInstanceGUID` | Service Instance GUID | `string` | IPI Install | AT |
| `.platform.powervs.serviceEndpoints` | Service Endpoints | `[]ServiceEndpoint` | Service Endpoints / PSC | AT |
| `.platform.powervs.tgName` | Transit Gateway | `string` | IPI Install | NT |
| `.platform.powervs.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.powervs.sysType` | System Type | `string` | Default Machine Types | AT |
| `.<pool>.platform.powervs.procType` | Processor Type | `string` | Default Machine Types | AT |
| `.<pool>.platform.powervs.processors` | Processors | `IntOrString` | Default Machine Types | AT |
| `.<pool>.platform.powervs.memoryGiB` | Memory (GiB) | `int32` | Default Machine Types | AT |
| `.<pool>.platform.powervs.smtLevel` | SMT Level | `string` | Default Machine Types | NT |

## vSphere

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.vsphere.vcenters` | vCenter Connections | `[]VCenter` | IPI Install | AT |
| `.platform.vsphere.failureDomains` | Failure Domains | `[]FailureDomain` | Failure Domains / Zones | AT |
| `.platform.vsphere.apiVIPs` | API VIPs | `[]string` | IPI Install | AT |
| `.platform.vsphere.ingressVIPs` | Ingress VIPs | `[]string` | IPI Install | AT |
| `.platform.vsphere.diskType` | Disk Provisioning | `thin` / `thick` / `eagerZeroedThick` | Custom Disk Types | NT |
| `.platform.vsphere.clusterOSImage` | Custom OS Image | `string` | Custom OS Image | NT |
| `.platform.vsphere.nodeNetworking` | Node Networking | `object` | IPI Install | AT |
| `.platform.vsphere.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.vsphere.dnsRecordsType` | DNS Records Type | `Internal` / `External` | IPI Install | NT |
| `.platform.vsphere.hosts` | Static IP Hosts | `[]Host` | IPI Install | AT |
| `.platform.vsphere.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | NT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.vsphere.cpus` | vCPU Count | `int32` | Default Machine Types | NT |
| `.<pool>.platform.vsphere.coresPerSocket` | Cores per Socket | `int32` | Default Machine Types | NT |
| `.<pool>.platform.vsphere.memoryMB` | Memory (MB) | `int64` | Default Machine Types | NT |
| `.<pool>.platform.vsphere.osDisk.diskSizeGB` | OS Disk Size (GB) | `int32` | Custom Disk Types | NT |
| `.<pool>.platform.vsphere.zones` | Failure Domain Zones | `[]string` | Failure Domains / Zones | AT |
| `.<pool>.platform.vsphere.dataDisks` | Data Disks | `[]DataDisk` | Custom Disk Types | NT |

## oVirt

oVirt is a legacy platform with destroy-only support. Fields are listed for
completeness.

### Platform fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.platform.ovirt.ovirt_cluster_id` | Cluster ID | `string` | IPI Install | NT |
| `.platform.ovirt.ovirt_storage_domain_id` | Storage Domain ID | `string` | IPI Install | NT |
| `.platform.ovirt.ovirt_network_name` | Network Name | `string` | IPI Install | NT |
| `.platform.ovirt.vnicProfileID` | VNIC Profile ID | `string` | IPI Install | NT |
| `.platform.ovirt.api_vips` | API VIPs | `[]string` | IPI Install | NT |
| `.platform.ovirt.ingress_vips` | Ingress VIPs | `[]string` | IPI Install | NT |
| `.platform.ovirt.affinityGroups` | Affinity Groups | `[]AffinityGroup` | IPI Install | NT |
| `.platform.ovirt.loadBalancer` | Load Balancer | `object` | Load Balancer Config | NT |
| `.platform.ovirt.defaultMachinePlatform` | Default Machine Platform | `MachinePool` | Default Machine Types | AT |

### Machine pool fields

| JSON Path | Field | Type | Feature(s) | Test Coverage |
|-----------|-------|------|------------|---------------|
| `.<pool>.platform.ovirt.instanceTypeID` | Instance Type ID | `string` | Default Machine Types | AT |
| `.<pool>.platform.ovirt.cpu` | CPU Config | `object` | Default Machine Types | AT |
| `.<pool>.platform.ovirt.memoryMB` | Memory (MB) | `int32` | Default Machine Types | AT |
| `.<pool>.platform.ovirt.osDisk.sizeGB` | OS Disk Size (GB) | `int64` | Custom Disk Types | NT |
| `.<pool>.platform.ovirt.vmType` | VM Type | `desktop` / `server` / `high_performance` | Default Machine Types | AT |
| `.<pool>.platform.ovirt.clone` | Clone Disks | `bool` | IPI Install | NT |
| `.<pool>.platform.ovirt.sparse` | Sparse Provisioning | `bool` | IPI Install | NT |

## Cross-Reference: Features to Fields

This section provides a reverse index from each feature in the
[Platform Feature Matrix](platform_feature_matrix.md) to the install-config
fields that drive it. Only features with direct field mappings are listed.

**IPI Install** - `.baseDomain`, `.sshKey`, `.pullSecret`, `.credentialsMode`,
`.platform.<name>.region`, `.platform.<name>.defaultMachinePlatform`, plus
platform-specific fields (VIPs, project IDs, cloud names, etc.)

**UPI Install** - `.baseDomain`, `.platform.<name>.region`, plus UPI templates
consume the generated install-config

**Agent-based Install** - `.platform.baremetal.hosts`,
`.platform.<name>.apiVIPs`, `.platform.<name>.ingressVIPs`

**Single Node (SNO)** - `.controlPlane.replicas` (set to 1),
`.compute[].replicas` (set to 0)

**Existing VPC / VNet** - `.platform.aws.vpc`, `.platform.azure.virtualNetwork`,
`.platform.azure.networkResourceGroupName`, `.platform.gcp.network`,
`.platform.gcp.controlPlaneSubnet`, `.platform.gcp.computeSubnet`,
`.platform.ibmcloud.vpcName`, `.platform.ibmcloud.controlPlaneSubnets`,
`.platform.ibmcloud.computeSubnets`, `.platform.openstack.controlPlanePort`,
`.platform.powervs.vpcName`, `.platform.powervs.vpcSubnets`,
`.platform.nutanix.subnetUUIDs`

**Shared VPC / XPN** - `.platform.gcp.networkProjectID`,
`.platform.gcp.dns.privateZone`

**Private Cluster** - `.publish` (set to `Internal`)

**Proxy Support** - `.proxy.httpProxy`, `.proxy.httpsProxy`, `.proxy.noProxy`,
`.additionalTrustBundle`

**Dual Stack (IPv4+IPv6)** - `.networking.machineNetwork`,
`.networking.clusterNetwork`, `.networking.serviceNetwork`,
`.platform.aws.ipFamily`, `.platform.azure.ipFamily`

**User-Provisioned DNS** - `.platform.<name>.userProvisionedDNS`,
`.platform.aws.hostedZone`, `.platform.aws.hostedZoneRole`

**Default Machine Types** - `.<pool>.platform.<name>.type` (or equivalent:
`.cpus`, `.memoryMiB` for Nutanix/vSphere; `.sysType`, `.procType` for Power VS)

**Custom Disk Types** - `.<pool>.platform.<name>.osDisk` (or `.rootVolume` for
AWS/OpenStack)

**Confidential Compute** - `.<pool>.platform.aws.cpuOptions`,
`.<pool>.platform.azure.settings`, `.<pool>.platform.azure.osDisk.securityProfile`,
`.<pool>.platform.gcp.confidentialCompute`,
`.<pool>.platform.gcp.onHostMaintenance`

**Secure Boot / UEFI** - `.<pool>.platform.azure.settings`,
`.<pool>.platform.gcp.secureBoot`, `.<pool>.platform.nutanix.bootType`

**Dedicated Hosts** - `.<pool>.platform.aws.hostPlacement`,
`.<pool>.platform.ibmcloud.dedicatedHosts`

**Custom OS Image** - `.<pool>.platform.aws.amiID`,
`.<pool>.platform.azure.osImage`, `.<pool>.platform.gcp.osImage`,
`.platform.openstack.clusterOSImage`, `.platform.powervs.clusterOSImage`,
`.platform.vsphere.clusterOSImage`

**KMS / Disk Encryption** - `.<pool>.platform.aws.rootVolume.kmsKeyARN`,
`.<pool>.platform.azure.osDisk.diskEncryptionSet`,
`.<pool>.platform.azure.encryptionAtHost`,
`.<pool>.platform.gcp.osDisk.encryptionKey`,
`.<pool>.platform.ibmcloud.bootVolume.encryptionKey`

**User Tags / Labels** - `.platform.aws.userTags`,
`.platform.azure.userTags`, `.platform.gcp.userLabels`,
`.platform.gcp.userTags`

**Service Endpoints / PSC** - `.platform.aws.serviceEndpoints`,
`.platform.gcp.endpoint`, `.platform.ibmcloud.serviceEndpoints`,
`.platform.powervs.serviceEndpoints`

**FIPS Mode** - `.fips`

**Load Balancer Config** - `.platform.<name>.loadBalancer` (available on
baremetal, nutanix, openstack, powervc, vsphere, ovirt)

**Failure Domains / Zones** - `.<pool>.platform.<name>.zones`,
`.platform.nutanix.failureDomains`, `.platform.vsphere.failureDomains`

**Feature Gates** - `.featureSet`, `.featureGates`

**Sovereign Cloud Support** - `.platform.gcp.projectID` (domain-scoped),
`.platform.gcp.region` (`u-` prefix), `.platform.azure.cloudName`,
plus GCD-specific defaults and restrictions on `.platform.gcp.endpoint`,
`.<pool>.platform.gcp.osImage`, `.<pool>.platform.gcp.type`

## References

- [Platform Feature Matrix](platform_feature_matrix.md)
- [Test documentation structure and conventions](README.md)
- [Install-config CRD definition](../../data/data/install.openshift.io_installconfigs.yaml)
- [InstallConfig Go type](../../pkg/types/installconfig.go)
