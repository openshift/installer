# iSCSI Boot

### What is iSCSI?
iSCSI (Internet Small Computer Systems Interface) is a network storage protocol that allows clients (initiators) to send SCSI commands to storage devices (targets) over IP networks. iSCSI provides a centralized and cost-effective storage management to use remote storage as if it were local disk storage.


### Supported since version: OpenShift 4.18

## Note: 
- Static networking is currently not supported when using iSCSI. This is because coreos-installer does not support the `--copy-network` when booting off a separate iSCSI network, see https://github.com/coreos/coreos-installer/issues/1389. 
- ISCSI currently only supported on x86 platforms. See https://issues.redhat.com/browse/MGMT-19307 

## Key Systemd Services for iSCSI Boot

### `iscsistart.service`
- **Path**: [`data/data/agent/systemd/units/iscsistart.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/iscsistart.service)
- **Purpose**: Initiates the iSCSI boot process by preparing and connecting to the iSCSI target. Calls the iscsistart utility which allows an iSCSI session to be started for booting off an iSCSI target
- **Platform**: non OCI platforms such as baremetal, none

### `iscsiadm.service`
- **Path**: [`data/data/agent/systemd/units/iscsiadm.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/iscsiadm.service)
- **Purpose**: Calls the iscsiadm utility that will start iscsid to manage the iSCSI session needed to install the final rhcos image
- **Platform**: non OCI platforms such as baremetal, none

You can refer to the agent installer services for an example of the installation workflow, particularly for non-OCI platforms, here. https://github.com/openshift/installer/blob/main/docs/user/agent/agent_installer_services-install_workflow.png

### `oci-eval-user-data.service`
- **Path**: [`data/data/agent/systemd/units/oci-eval-user-data.service`](https://github.com/openshift/installer/blob/main/data/data/agent/systemd/units/oci-eval-user-data.service)
- **Purpose**: Executes a script passed via OCI cloud-init user data. This is **required only for OCI**. The script is maintained by OCI and is available here [`terraform-stacks/shared_modules/compute/userdata
/iscsi-oci-configure-secondary-nic.sh`](https://github.com/oracle-quickstart/oci-openshift/blob/main/terraform-stacks/shared_modules/compute/userdata/iscsi-oci-configure-secondary-nic.sh)
- **Platform**: OCI
- **Note**: Without this service, network configuration will fail, causing the cluster installations to hang.


## iSCSI Boot on External OCI

To test and configure iSCSI boot on OCI:
1. **Select iSCSI Boot Volume**: 
   - When creating a custom image, choose **iSCSI** as the boot volume type.
2. **Add a Secondary NIC**: 
   - Attach a secondary network interface card (NIC) when creating the instance.
3. **Assign Rendezvous IP**: 
   - Ensure the **rendezvous IP** is configured on the **secondary NIC**.
4. **Load Balancer (LB) Configuration**: 
   - Attach the LB to the **secondary NIC**.
   - Ensure all OpenShift traffic is routed through secondary IPs, **not the iSCSI IPs**.
5. **Configure `oci-eval-user-data.service`**:
   - Include the [Oracle-maintained script](https://github.com/oracle-quickstart/oci-openshift/blob/beta-v1/custom_manifests/manifests/oci-eval-user-data-master.yaml) in the **cloud-init** section when creating the OCI instance.
   - **Important**: Without this script, network configuration fails.

### Subnet Usage
- **Caution**: The iSCSI subnet must not be used or referenced by OpenShift.

## Validation in Assisted Service

Assisted Service includes a validation to ensure that the iSCSI NIC is not part of the machine networks.

#### Validation Details:
- **ID**: `no-iscsi-nic-belongs-to-machine-cidr`
- **Status**: `success`
- **Message**: "Network interface connected to iSCSI disk does not belong to machine network CIDRs."

## Testing and Verification

### Minimal Agent ISO
- For any platform using iSCSI boot, a **minimal agent ISO** must be generated. By default for external OCI platform, a minimal ISO is generated.

### Verify iSCSI Installation

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
