# iSCSI Boot

### What is iSCSI?
iSCSI (Internet Small Computer Systems Interface) is a network storage protocol that allows clients (initiators) to send SCSI commands to storage devices (targets) over IP networks. iSCSI provides a centralized and cost-effective storage management to use remote storage as if it were local disk storage.

### Why iSCSI boot support?
Oracle Cloud Infrastructure (OCI) bare metal instances support iSCSI-attached boot volumes. To deploy OpenShift on these instances, the agent-based installer needs to establish and manage iSCSI sessions during installation so that the remote disk is available as the root device. While the primary motivation was OCI, the implementation is generic and also supports non-OCI bare metal environments that boot from iSCSI targets. See [MGMT-17556](https://redhat.atlassian.net/browse/MGMT-17556).

### Supported since version: OpenShift 4.18

For user-facing documentation, see [Preparing installation assets for iSCSI booting](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18/html/installing_an_on-premise_cluster_with_the_agent-based_installer/installing-using-iscsi).

## Prerequisites

- Only DHCP networking is supported when using iSCSI. Static networking is currently not supported because coreos-installer does not support `--copy-network` when booting off a separate iSCSI network, see https://github.com/coreos/coreos-installer/issues/1389.
- You must create an additional network for iSCSI that is separate from the machine network. The machine network is rebooted during cluster installation and cannot be used for the iSCSI session.
- iSCSI is currently only supported on x86 platforms. See https://issues.redhat.com/browse/MGMT-19307.
- A **minimal agent ISO** must be generated for any platform using iSCSI boot. By default for the external OCI platform, a minimal ISO is generated.

## Configuration

### `rootDeviceHints`

If the host on which you are booting the agent ISO image also has a local disk, it might be necessary to specify the iSCSI disk name in the `rootDeviceHints` parameter to ensure that it is chosen as the boot disk for the final RHCOS image. For example:

```yaml
rootDeviceHints:
  deviceName: /dev/sda
```

In a diskless environment (where the iSCSI disk is the only available disk), you do not need to set `rootDeviceHints`.

## Key Systemd Services for iSCSI Boot

The iSCSI boot flow on non-OCI platforms uses two services that work together:

1. `iscsistart.service` establishes the initial iSCSI session so the host can boot from the iSCSI target.
2. `iscsiadm.service` then takes over to manage the session using the full `iscsid` daemon, which is needed for the RHCOS image installation to the iSCSI disk.

### `iscsistart.service`
- **Path**: [`data/data/agent/systemd/units/iscsistart.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/iscsistart.service)
- **Purpose**: Initiates the iSCSI boot process by reading iSCSI boot firmware parameters (iBFT) and establishing a session to the iSCSI target. The `iscsistart` utility creates an initial iSCSI session without requiring the full `iscsid` daemon, allowing the host to access the remote iSCSI disk as a local block device for booting.
- **Supported Platforms**: non-OCI platforms such as baremetal, none

### `iscsiadm.service`
- **Path**: [`data/data/agent/systemd/units/iscsiadm.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/iscsiadm.service)
- **Purpose**: Starts the `iscsid` daemon and uses `iscsiadm` to take over management of the iSCSI session that was initially established by `iscsistart.service`. The `iscsid` daemon handles session persistence, error recovery, and re-authentication, which are required during the RHCOS image installation to the iSCSI disk by coreos-installer.
- **Supported Platforms**: non-OCI platforms such as baremetal, none

You can refer to the [agent installer services install workflow](https://github.com/openshift/installer/blob/main/docs/user/agent/agent_installer_services-install_workflow.png) for an example of the installation workflow on non-OCI platforms.

### `oci-eval-user-data.service`
- **Path**: [`data/data/agent/systemd/units/oci-eval-user-data.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/oci-eval-user-data.service)
- **Purpose**: Executes a script passed via OCI cloud-init user data. This is **required only for OCI**. The script is maintained by OCI and is available here: [`terraform-stacks/shared_modules/compute/userdata/iscsi-oci-configure-secondary-nic.sh`](https://github.com/oracle-quickstart/oci-openshift/blob/main/terraform-stacks/shared_modules/compute/userdata/iscsi-oci-configure-secondary-nic.sh)
- **Supported Platforms**: OCI
- **Note**: Without this service, network configuration will fail, causing the cluster installation to hang.

## iSCSI Boot on External OCI

OCI (Oracle Cloud Infrastructure) supports iSCSI as a boot volume attachment type. Bare metal instances on OCI use iSCSI attachments by default, and VM instances can also be configured to use iSCSI for better storage IOPS performance. When an OpenShift cluster on OCI is deployed with iSCSI-attached boot volumes, the iSCSI boot support in the agent-based installer is required. Because the iSCSI network is separate from the cluster network, a secondary NIC must be configured for OpenShift traffic. The secondary vNIC is automatically created in the Terraform stack provided by Oracle.

To configure iSCSI boot on OCI:
1. **Select iSCSI Boot Volume**:
   - When creating a custom image, choose **iSCSI** as the boot volume type.
2. **Add a Secondary NIC**:
   - Attach a secondary network interface card (NIC) when creating the instance.
3. **Assign Rendezvous IP**:
   - Ensure the **rendezvous IP** is configured on the **secondary NIC** (not the iSCSI NIC). The secondary NIC is the interface used for the cluster (machine) network. The rendezvous IP must be reachable on this network so that agent nodes can communicate with each other during bootstrap.
4. **Load Balancer (LB) Configuration**:
   - Attach the LB to the **secondary NIC**.
   - Ensure all OpenShift traffic is routed through secondary IPs, **not the iSCSI IPs**.
5. **Configure `oci-eval-user-data.service`**:
   - Include the [Oracle-maintained script](https://github.com/oracle-quickstart/oci-openshift/blob/beta-v1/custom_manifests/manifests/oci-eval-user-data-master.yaml) in the **cloud-init** section when creating the OCI instance.
   - **Important**: Without this script, network configuration fails.

### Subnet Usage
- **Caution**: The iSCSI subnet must not be used or referenced by OpenShift. The iSCSI network must remain separate from the machine network.

## Validation in Assisted Service

Assisted Service includes a validation to ensure that the iSCSI NIC is not part of the machine networks.

#### Validation Details:
- **ID**: `no-iscsi-nic-belongs-to-machine-cidr`
- **Status**: `success`
- **Message**: "Network interface connected to iSCSI disk does not belong to machine network CIDRs."

## Verify iSCSI Installation

Check the root device on the machine; it should be an iSCSI disk.
Use the following commands:
```bash
lsblk -S
ls -l /sys/block/sda
```
Example Output:
```bash
[core@abtest-cp-1-ad1 ~]$ lsblk -S
NAME HCTL       TYPE VENDOR   MODEL        REV SERIAL                           TRAN
sda  0:0:0:1    disk ORACLE   BlockVolume 1.0  600d1c8cec9e4882a55d99315005b858 iscsi

[core@abtest-cp-1-ad1 ~]$ ls -l /sys/block/sda
lrwxrwxrwx. 1 root root 0 Nov 18 18:09 /sys/block/sda -> ../devices/platform/host0/session1/target0:0:0/0:0:0:1/block/sda
```
