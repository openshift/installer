# Bare Metal Platform Customization

## Cluster-scoped properties

### Advanced Configuration Parameters

| Parameter | Default | Description |
| --- | --- | --- |
`libvirtURI` | `qemu://localhost/system` | The location of the hypervisor for running the bootstrap VM. See [Using a Remote Hypervisor](using-a-remote-hypervisor) for more details. |
`clusterProvisioningIP` | The third address on the provisioning network. `172.22.0.3` | The IP within the cluster where the provisioning services run. |
`bootstrapProvisioningIP` | The second address on the provisioning network. `172.22.0.2` | The IP on the bootstrap VM where the provisioning services run while the control plane is being deployed. |
`externalBridge` | `baremetal` | The name of the bridge of the hypervisor attached to the external network. |
`provisioningBridge` | `provisioning` | The name of the bridge on the hypervisor attached to the provisioning network. |
`provisioningNetworkCIDR` | `172.22.0.0/24` | The CIDR for the network to use for provisioning. |
`provisioningDHCPExternal` | `false` | Flag indicating that DHCP for the provisioning network is managed outside of the cluster by existing infrastructure services. |
`provisioningDHCPRange` | The tenth through the second last IP on the provisioning network. `172.22.0.10,172.22.0.254` | The IP range to use for hosts on the provisioning network. |
`defaultMachinePlatform` | | The default configuration used for machine pools without a platform configuration. |
`bootstrapOSImage` | *based on the release image* | A URL to override the default operating system image for the bootstrap node. The URL must contain a sha256 hash of the image. Example `https://mirror.example.com/images/qemu.qcow2.gz?sha256=a07bd...` |
`clusterOSImage` | *based on the release image* | A URL to override the default operating system for cluster nodes. The URL must include a sha256 hash of the image. Example `https://mirror.example.com/images/metal.qcow2.gz?sha256=3b5a8...` |

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
  the 10th to the second from last address. To use a different range, specify this
  using the provisioningDHCPRange option in the baremetal platform. This
  should be a comma-separated list indicating the start and end range.

Example:

```yaml
platform:
  baremetal:
    provisioningDHCPRange: "172.23.0.10,172.23.0.254"
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

## Using a Remote Hypervisor

The IPI installation process requires access to a libvirt-based
hypervisor host on which to run a bootstrap VM. The VM is removed
after the control plane is up and fully functional, so the hypervisor
is not needed to operate the cluster. When running the installer by
hand, it is most common to use the local host as the hypervisor. When
network topology requires, it is possible to use a separate host.

The `libvirtURI` can be used to specify the location of the remote
hypervisor. For example
`qemu+ssh://hyperuser@hypervisor.example.com/system` tells the
installer to connect to `hypervisor.example.com` over ssh as the
`hyperuser` user and create the bootstrap VM there.

The user on the host running the installer must be able to connect via
ssh to the hypervisor using the username given in the URI, without
being prompted for a password.

The user on the hypervisor must be in the `libvirt` group and have
permission to communicate with the libvirt services.

The hypervisor must meet the network requirements described in
the [Prerequisites](install_ipi.md#prerequisites) section.

Example:

```yaml
platform:
  baremetal:
    libvirtURI: qemu+ssh://hyperuser@hypervisor.example.com/system
```

## Disabling Certificate Verification for BMCs

By default TLS clients communicating with BMCs will require valid
certificates signed by a known certificate authority. In environments
where certificates are signed by unknown authorities, this behavior
can be disabled by setting `disableCertificateVerification` to `true`
for each `bmc` entry.

## Customizing ironic for provisioning

Should you need to adjust any of the config options in ironic you can
set ironicExtraConf in the platform section of install-config.yaml. Each
config option should be expressed as a key/value pair with the format

OS_<section>_\_<name>=<value> - where `section` and `name` are the
reprepresent the config option in ironic.conf e.g. to set a IPA ssh key and
set the number of ironic API workers

```yaml
platform:
  baremetal:
    ironicExtraConf: {"OS_PXE__PXE_APPEND_PARAMS":'nofb nomodeset vga=normal sshkey="ssh-rsa AAAA..."', "OS_API__API_WORKERS":"8"}
```
