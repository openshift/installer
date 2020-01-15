# Bare Metal Platform Customization

## Cluster-scoped properties

### Image Overrides

When doing a disconnected installation, the baremetal platform has the
additional requirement that we have locations to download the RHCOS
images. The installer downloads these from a location described in
[/data/data/rhcos.json](/data/data/rhcos.json), but they can be
overridden to point to a local mirror.

The SHA256 parameter in the URLs are required, and should match the
uncompressed SHA256 from rhcos.json.


* `bootstrapOSImage` (optional string): Override the image used for the
    bootstrap virtual machine.
* `clusterOSImage` (optional string): Override the image used for
    cluster machines.

Example:

```yaml
platform:
  baremetal:
      bootstrapOSImage: http://mirror.example.com/images/rhcos-43.81.201912131630.0-qemu.x86_64.qcow2.gz?sha256=f40e826ac4a6c5c073416a7bc0039ec8726a338885d2031e7607cec8783e580e
      clusterOSImage: http://mirror.example.com/images/rhcos-43.81.201912131630.0-openstack.x86_64.qcow2.gz?sha256=ffebbd68e8a1f2a245ca19522c16c86f67f9ac8e4e0c1f0a812b068b16f7265d
```

### Networking customization

By default, the baremetal IPI environment uses a provisioning network of
`172.22.0.0/24`, picks the 2nd and 3rd address of that subnet for the
bootstrap and cluster provisioning IP's, and operates an internal DHCP
and TFTP server in the cluster to support provisioning. Much of this can
be customized.


* `provisioningNetorkCIDR` (optional string): Override the default provisioning network.
* `bootstrapProvisioningIP` (optional string): Override the bootstrap
    provisioning IP. If unspecified, uses the 2nd address in the
    provisioning network's subnet.
* `provisioningHostIP` (optional string): Override the IP used by the
    cluster's provisioning infrastructure. If unspecified, uses the 3rd
    address in the provisioning network's subnet.

Example:

```yaml
platform:
  baremetal:
    provisioningNetworkCIDR: 172.23.0.0/16
    bootstrapProvisioningIP: 172.23.0.2
    provisioningHostIP: 172.23.0.3
```

* `provisioningDHCPRange` (optional string): By default, the installer picks a range from
  the 10th to 100th addresses. To use a different range, specify this
  using the provisioingDHCPRange option in the baremetal platform. This
  should be a comma-separated list indicating the start and end range.

Example:

```yaml
platform:
  baremetal:
    provisioningDHCPRange: "172.23.0.10,172.23.0.100"
```

* `provisioningDHCPExternal` (optional boolean): If you would prefer to
use an external DHCP server, you can specify provisioningDHCPExternal,
in which case the cluster will only run TFTP.  When using PXE boot for
the control plane and workers, your DHCP server needs to specify the
next-server as `bootstrapProvisioningIP` for the control plane, and
`provisioningHostIP` for the workers.

Example:

```yaml
platform:
  baremetal:
    provisioningDHCPExternal: true
```

