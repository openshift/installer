# oVirt Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, ovirt-specific properties.

## Cluster-scoped properties

* `ovirt_cluster_id` (required string): The oVirt cluster where the VMs will be created.
* `ovirt_storage_domain_id` (required string): The storage domain ID where the VM disks will be created.
* `ovirt_network_name` (required string): The network name where the VM nics will be created.
* `vnicProfileID` (required string): The ID of the [vNic profile][vnic-profile] used for the VM network interfaces.
    This can be inferred if the cluster network has a single profile.
* `api_vip` (required string): An IP address on the machineNetwork that will be assigned to the API VIP.
* `dns_vip` (required string): An IP address on the machineNetwork that will be assigned to the DNS VIP.
* `ingress_vip` (required string): An IP address on the machineNetwork that will be assigned to the Ingress VIP.

## Machine pools

* `cpu` (optional object): Defines the CPU of the VM.
    * `cores` (required integer): The number of cores. Total vCPUS is cores * sockets.
    * `sockets` (required integer): The number of sockets. Total vCPUS is cores * sockets.
* `memoryMB` (optional integer): Memory of the VM in MiB.
* `instanceTypeID` (optional string): The VM [instance-type][instance-type].
* `osDisk` (optional string): Defines the first and bootable disk of the VM.
    * `sizeGB` (required number): Size of the disk in GiB.
* `vmType` (optional string): The VM workload type. One of [high-performance][high-perf], server or desktop.  


## Installing to Existing VPC & Subnetworks

The installer can use an existing VPC and subnets when provisioning an OpenShift cluster. A VPC will be inferred from the provided subnets. For a standard installation, a private and public subnet should be specified. ([see example below](#pre-existing-vpc--subnets)). Both of the subnets must be within the IP range specified in `networking.machineNetwork`. 

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal oVirt install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  ovirt:
    api_vip: 10.46.8.230
    dns_vip: 10.46.8.231
    ingress_vip: 10.46.8.232
    ovirt_cluster_id: 68833f9f-e89c-4891-b768-e2ba0815b76b
    ovirt_storage_domain_id: ed7b0f4e-0e96-492a-8fff-279213ee1468
    ovirt_network_name: ovirtmgmt
    vnicProfileID: 3fa86930-0be5-4052-b667-b79f0a729692
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom machine pools

An example oVirt install config with custom machine pools:

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    ovirt:
      cpu:
        cores: 4
        sockets: 2
      memoryMB: 65536
      osDisk:
        sizeGB: 100
      vmType: high_performance
  replicas: 3
compute:
- name: worker
  platform:
    ovirt:
      cpu:
        cores: 4
        sockets: 4
      memoryMB: 65536
      osDisk:
        sizeGB: 200
      vmType: high_performance
  replicas: 5
metadata:
  name: test-cluster
platform:
  ovirt:
    api_vip: 10.46.8.230
    dns_vip: 10.46.8.231
    ingress_vip: 10.46.8.232
    ovirt_cluster_id: 68833f9f-e89c-4891-b768-e2ba0815b76b
    ovirt_storage_domain_id: ed7b0f4e-0e96-492a-8fff-279213ee1468
    ovirt_network_name: ovirtmgmt
    vnicProfileID: 3fa86930-0be5-4052-b667-b79f0a729692
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

[instance-type]: https://www.ovirt.org/develop/release-management/features/virt/instance-types.html
[vnic-profile]: https://www.ovirt.org/develop/release-management/features/sla/vnic-profiles.html
[high-perf]: https://www.ovirt.org/develop/release-management/features/virt/high-performance-vm.html
