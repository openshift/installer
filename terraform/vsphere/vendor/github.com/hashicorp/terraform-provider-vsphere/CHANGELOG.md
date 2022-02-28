## 1.24.3 (December 14, 2020)

BUG FIXES:
* `resource/vm`: Support for no disks in config ([#1241](https://github.com/hashicorp/terraform-provider-vsphere/1241))
* `resource/vm`: Make API timeout configurable when building VMs ([#1278](https://github.com/hashicorp/terraform-provider-vsphere/1278))

## 1.24.2 (October 16, 2020)

BUG FIXES:
* `resource/vm`: Prevent guest_id nil condition. ([#1234](https://github.com/hashicorp/terraform-provider-vsphere/pull/1234))

## 1.24.1 (October 07, 2020)

IMPROVEMENTS:
* `data/content_library_item`: Add `type` to content library item data source. ([#1184](https://github.com/hashicorp/terraform-provider-vsphere/pull/1184))
* `resource/virtual_switch`: Fix port group resource to enable LACP in virtual switch only. ([#1214](https://github.com/hashicorp/terraform-provider-vsphere/pull/1214))
* `resource/distributed_port_group`: Import distributed port group using MOID. ([#1208](https://github.com/hashicorp/terraform-provider-vsphere/pull/1208))
* `resource/host_port_group`: Add support for importing. ([#1194](https://github.com/hashicorp/terraform-provider-vsphere/pull/1194))
* `resource/VM`: Allow more config options to be changed from OVF. ([#1218](https://github.com/hashicorp/terraform-provider-vsphere/pull/1218))
* `resource/VM`: Convert folder path to MOID. ([#1207](https://github.com/hashicorp/terraform-provider-vsphere/pull/1207))

BUG FIXES:
* `resource/datastore_cluster`: Fix missing field in import. ([#1203](https://github.com/hashicorp/terraform-provider-vsphere/pull/1203))
* `resource/VM`: Change default OS method on bare VMs. ([#1217](https://github.com/hashicorp/terraform-provider-vsphere/pull/1217))
* `resource/VM`: Read virtual machine after clone and OVF/OVA deploy. ([#1221](https://github.com/hashicorp/terraform-provider-vsphere/pull/1221))

## 1.24.0 (September 02, 2020)

BUG FIXES:
* `resource/vm`: Skip SCSI controller check when empty. ([#1179](https://github.com/hashicorp/terraform-provider-vsphere/pull/1179))
* `resource/vm`: Make storage_policy_id computed to prevent flapping when unset. ([#1185](https://github.com/hashicorp/terraform-provider-vsphere/pull/1185))
* `resource/vm`: Ignore nil objects in host network on read. ([#1186](https://github.com/hashicorp/terraform-provider-vsphere/pull/1186))
* `resource/vm`: Keep progress channel open when deploying an OVF. ([#1187](https://github.com/hashicorp/terraform-provider-vsphere/pull/1187))
* `resource/vm`: Set SCSI controller type to unknown when nil. ([#1188](https://github.com/hashicorp/terraform-provider-vsphere/pull/1188))

IMPROVEMENTS:
* `resource/content_library_item`: Add local upload, OVA, and vm-template 
  sources. ([#1196](https://github.com/hashicorp/terraform-provider-vsphere/pull/1196))
* `resource/content_library`: Subscription and publication support. ([#1197](https://github.com/hashicorp/terraform-provider-vsphere/pull/1197))
* `resource/vm`: Content library vm-template, disk type, and vApp property 
  support. ([#1198](https://github.com/hashicorp/terraform-provider-vsphere/pull/1198))

## 1.23.0 (August 21, 2020)

BUG FIXES:
* `resource/vnic`: Fix missing fields on vnic import. ([#1162](https://github.com/hashicorp/terraform-provider-vsphere/pull/1162))
* `resource/virtual_machine`: Ignore thin_provisioned and eagerly_scrub during DiskPostCloneOperation. ([#1161](https://github.com/hashicorp/terraform-provider-vsphere/pull/1161))
* `resource/virtual_machine`: Fix SetHardwareOptions to fetch the hardware version from QueryConfigOption. ([#1159](https://github.com/hashicorp/terraform-provider-vsphere/pull/1159))

IMPROVEMENTS:
* `resource/virtual_machine`: Allow performing a linked-clone from a template. ([#1158](https://github.com/hashicorp/terraform-provider-vsphere/pull/1158))
* `data/virtual_machine`: Merge the virtual machine configuration schema. ([#1157](https://github.com/hashicorp/terraform-provider-vsphere/pull/1157))

## 1.22.0 (August 07, 2020)

FEATURES:
* `resource/compute_cluster`: Basic vSAN support on compute clusters. ([#1151](https://github.com/hashicorp/terraform-provider-vsphere/pull/1151))
* `resource/role`: Resource and data source to create and manage vSphere roles. ([#1144](https://github.com/hashicorp/terraform-provider-vsphere/pull/1144))
* `resource/entity_permission`: Resource to create and manage vSphere permissions. ([#1144](https://github.com/hashicorp/terraform-provider-vsphere/pull/1144))
* `data/entity_permission`: Data source to acquire ESXi host thumbprints . ([#1142](https://github.com/hashicorp/terraform-provider-vsphere/pull/1142))

## 1.21.1 (July 20, 2020)
BUG FIXES:
* `resource/vm`: Set guest_id before customization. ([#1139](https://github.com/hashicorp/terraform-provider-vsphere/pull/1139))

## 1.21.0 (June 30, 2020)
FEATURES:
* `resource/vm`: Support for SATA and IDE disks. ([#1118](https://github.com/hashicorp/terraform-provider-vsphere/pull/1118))

## 1.20.0 (June 23, 2020)

FEATURES:
* `resource/vm`: Add support for OVA deployment. ([#1105](https://github.com/hashicorp/terraform-provider-vsphere/pull/1105))

BUG FIXES:
* `resource/vm`: Delete disks on destroy when deployed from OVA/OVF. ([#1106](https://github.com/hashicorp/terraform-provider-vsphere/pull/1106))
* `resource/vm`: Skip PCI passthrough operations if there are no changes. ([#1112](https://github.com/hashicorp/terraform-provider-vsphere/pull/1112))

## 1.19.0 (June 16, 2020)

FEATURES:
* `data/dynamic`: Data source which can be used to match any tagged managed object. ([#1103](https://github.com/hashicorp/terraform-provider-vsphere/pull/1103))
* `resource/vm_storage_policy_profile`: A resource for tag based storage placement.
  policies management. ([#1094](https://github.com/hashicorp/terraform-provider-vsphere/pull/1094))
* `resource/virtual_machine`: Add support for PCI passthrough devices on virtual
  machines. ([#1099](https://github.com/hashicorp/terraform-provider-vsphere/pull/1099))
* `data/host_pci_device`: Data source which will locate the address of a PCI
  device on a host. ([#1099](https://github.com/hashicorp/terraform-provider-vsphere/pull/1099))

## 1.18.3 (June 01, 2020)

IMPROVEMENTS:
* `resource/custom_attribute`: Fix id in error message when category is 
  missing. ([#1088](https://github.com/hashicorp/terraform-provider-vsphere/pull/1088))
* `resource/virtual_machine`: Add vApp properties with OVF deployment. ([#1082](https://github.com/hashicorp/terraform-provider-vsphere/pull/1082))

## 1.18.2 (May 22, 2020)

IMPROVEMENTS:
* `resource/host` & `resource/compute_cluster`: Add arguments for specifying 
  if cluster management should be handled in `host` or `compute_cluster` 
  resource. ([#1085](https://github.com/hashicorp/terraform-provider-vsphere/pull/1085))
* `resource/virtual_machine`: Handle OVF argument validation during VM 
  creation. ([#1084](https://github.com/hashicorp/terraform-provider-vsphere/pull/1084))
* `resource/host`: Disconnect rather than entering maintenance mode when 
  deleting. ([#1083](https://github.com/hashicorp/terraform-provider-vsphere/pull/1083)) 


## 1.18.1 (May 12, 2020)

BUG FIXES:
* `resource/virtual_machine`: Skip unexpected NIC entries. ([#1067](https://github.com/hashicorp/terraform-provider-vsphere/pull/1067))
* Respect `session_persistence` for REST sessions. ([#1077](https://github.com/hashicorp/terraform-provider-vsphere/pull/1077))

## 1.18.0 (May 04, 2020)

FEATURES:
* `resource/virtual_machine`: Allow users to deploy OVF templates from both 
  from local system and remote URL. ([#1052](https://github.com/hashicorp/terraform-provider-vsphere/pull/1052))

## 1.17.4 (April 29, 2020)

IMPROVEMENTS:
* `resource/virtual_machine`: Mark `product_key` as sensitive. ([#1045](https://github.com/hashicorp/terraform-provider-vsphere/pull/1045))
* `resource/virtual_machine`: Increase max `hardware_version` for vSphere v7.0. ([#1056](https://github.com/hashicorp/terraform-provider-vsphere/pull/1056))

BUG FIXES:
* `resource/virtual_machine`: Fix to disk bus sorting. ([#1039](https://github.com/hashicorp/terraform-provider-vsphere/pull/1039))
* `resource/virtual_machine`: Only include `hardware_version` in CreateSpecs. ([#1055](https://github.com/hashicorp/terraform-provider-vsphere/pull/1055))

## 1.17.3 (April 22, 2020)


IMPROVEMENTS:
* Use built in session persistence in govmomi. ([#1050](https://github.com/hashicorp/terraform-provider-vsphere/pull/1050))

## 1.17.2 (April 13, 2020)

IMPROVEMENTS:
* `resource/virtual_disk`: Support VMDK files. ([#987](https://github.com/hashicorp/terraform-provider-vsphere/pull/987))

BUG FIXES:
* `resource/virtual_machine`: Fix disk controller sorting. ([#1032](https://github.com/hashicorp/terraform-provider-vsphere/pull/1032))

## 1.17.1 (April 07, 2020)

IMPROVEMENTS:
* `resource/virtual_machine`: Add support for hardware version tracking and
  upgrading. ([#1020](https://github.com/hashicorp/terraform-provider-vsphere/pull/1020))
* `data/vsphere_network`: Handle cases of network port groups with same name
  using `distributed_virtual_switch_uuid`. ([#1001](https://github.com/hashicorp/terraform-provider-vsphere/pull/1001))

BUG FIXES:
* `resource/virtual_machine`: Fix working with orphaned devices. ([#1005](https://github.com/hashicorp/terraform-provider-vsphere/pull/1005))
* `resource/virtual_machine`: Ignore `guest_id` with content library. ([#1014](https://github.com/hashicorp/terraform-provider-vsphere/pull/1014))

## 1.17.0 (March 23, 2020)

FEATURES:
* **New Data Source:** `content_library` ([#985](https://github.com/hashicorp/terraform-provider-vsphere/pull/985))
* **New Data Source:** `content_library_item` ([#985](https://github.com/hashicorp/terraform-provider-vsphere/pull/985))
* **New Resource:** `content_library` ([#985](https://github.com/hashicorp/terraform-provider-vsphere/pull/985))
* **New Resource:** `content_library_item` ([#985](https://github.com/hashicorp/terraform-provider-vsphere/pull/985))

IMPROVEMENTS:
* `resource/virtual_machine`: Add `poweron_timeout` option for the amount of 
  time to give a VM to power on. ([#990](https://github.com/hashicorp/terraform-provider-vsphere/pull/990))

## 1.16.2 (March 04, 2020)

IMPROVEMENTS:
* `resource/virtual_machine`: Optimize OS family query. ([#959](https://github.com/hashicorp/terraform-provider-vsphere/pull/959))
* Migrate provider to Terraform plugin SDK. ([#982](https://github.com/hashicorp/terraform-provider-vsphere/pull/982))

## 1.16.1 (February 06, 2020)

BUG FIXES:
* `resource/virtual_machine`: Set `storage_policy_id` based off of VM rather
  than template. ([#970](https://github.com/hashicorp/terraform-provider-vsphere/pull/970))

## 1.16.0 (February 04, 2020)

FEATURES:
* **New Data Source:** `storage_policy` ([#881](https://github.com/hashicorp/terraform-provider-vsphere/pull/881))

IMPROVEMENTS:
* Switch to govmomi REST client ([#955](https://github.com/hashicorp/terraform-provider-vsphere/pull/955))
* Add storage policy to `virtual_machine` resource. ** Requires `profile-driven 
  storage` permissions for Terraform user. ([#881](https://github.com/hashicorp/terraform-provider-vsphere/pull/881))

## 1.15.0 (January 23, 2020)

IMPROVEMENTS:
* `resource/virtual_machine`: Do not throw error when disk path is not known
  yet. ([#944](https://github.com/hashicorp/terraform-provider-vsphere/pull/944))

BUG FIXES:
* `resource/virtual_machine`: Do not set datastoreID in RelocateSpec when
  datastore_cluster is set. ([#933](https://github.com/hashicorp/terraform-provider-vsphere/pull/933))
* `resource/vapp_container`: Fix handling of child vApp containers. ([#941](https://github.com/hashicorp/terraform-provider-vsphere/pull/941))
* `resource/virtual_disk`: Enforce .vmdk suffix on `vmdk_path`. ([#942](https://github.com/hashicorp/terraform-provider-vsphere/pull/942))

## 1.14.0 (December 18, 2019)

IMPROVEMENTS
* `resource/host` Add details to error messages. ([#850](https://github.com/hashicorp/terraform-provider-vsphere/pull/850))
* `resource/virtual_machine`: Pick default datastore for extra disks. ([#897](https://github.com/hashicorp/terraform-provider-vsphere/pull/897))
* `resource/virtual_machine`: Extend `ignored_guest_ips` to support CIDR. ([#841](https://github.com/hashicorp/terraform-provider-vsphere/pull/841))

FEATURES:
* **New Resource:** `vsphere_vnic` ([#876](https://github.com/hashicorp/terraform-provider-vsphere/pull/876))

BUG FIXES:
* `resource/virtual_machine`: Allow blank networkID in order to support cloning
 into clusters that do not include the source network. ([#787](https://github.com/hashicorp/terraform-provider-vsphere/pull/787))
* `resource/host`: Properly handle situation where NIC teaming policy is `nil`. ([#889](https://github.com/hashicorp/terraform-provider-vsphere/pull/889))
* Limit scope when listing network interfaces. ([#840](https://github.com/hashicorp/terraform-provider-vsphere/pull/840))
* `resource/compute_cluster`: Set HA Admission Control Failure to `off` before deleting. ([#891](https://github.com/hashicorp/terraform-provider-vsphere/pull/891))
* `resource/virtual_machine_snapshot`: Fix typo in error condition. ([#906](https://github.com/hashicorp/terraform-provider-vsphere/pull/906))
* `tags`: Return matched tag rather than last tag in list. ([#910](https://github.com/hashicorp/terraform-provider-vsphere/pull/910))
* `resource/virtual_machine`: Unmount ISO when switching CDROM backends. ([#920](https://github.com/hashicorp/terraform-provider-vsphere/pull/920))
* `resource/virtual_machine`: Migrate VM when moving to different root resource pool. ([#931](https://github.com/hashicorp/terraform-provider-vsphere/pull/931))

## 1.13.0 (October 01, 2019)

IMPROVEMENTS:
* Add `vim_keep_alive` which sets a keepalive interval for VIM session. ([#792](https://github.com/hashicorp/terraform-provider-vsphere/pull/792))
* `resource/virtual_machine`: Mark `windows_sysprep_text` as sensitive. ([#802](https://github.com/hashicorp/terraform-provider-vsphere/pull/802))

FEATURES:
* **New Resource:** `vsphere_host` ([#836](https://github.com/hashicorp/terraform-provider-vsphere/pull/836))

BUG FIXES:
* `resource/virtual_machine`: Change the way we detect if a VM is in a vApp. ([#825](https://github.com/hashicorp/terraform-provider-vsphere/pull/825))
* Delete tags and tag_categories when they are removed. ([#801](https://github.com/hashicorp/terraform-provider-vsphere/pull/801))

## 1.12.0 (June 19, 2019)

IMPROVEMENTS:
* `resource/virtual_machine`: Allow cloning of powered on virtual machines. ([#785](https://github.com/hashicorp/terraform-provider-vsphere/pull/785))
* Add keep alive timer for VIM sessions. ([#792](https://github.com/hashicorp/terraform-provider-vsphere/pull/792))

BUG FIXES:
* `resource/virtual_machine`: Ignore validation when interpolation is not
  available. ([#784](https://github.com/hashicorp/terraform-provider-vsphere/pull/784))
* `resource/virtual_machine`: Only set vApp properties that are
  UserConfigurable. ([#751](https://github.com/hashicorp/terraform-provider-vsphere/pull/751))
* `resource/virtual_machine`: Set `network_id` to empty string when cloning a
  `virtual_machine` to a cluster that is not part of source DVS. ([#787](https://github.com/hashicorp/terraform-provider-vsphere/pull/787))

## 1.11.0 (May 09, 2019)

IMPROVEMENTS:

* Add support for importing datacenters. ([#737](https://github.com/hashicorp/terraform-provider-vsphere/pull/737))
* Document max character limit on `run_once_command_list`. ([#748](https://github.com/hashicorp/terraform-provider-vsphere/pull/748))
* Add missing ENV variable checks for acceptance tests. ([#758](https://github.com/hashicorp/terraform-provider-vsphere/pull/758))
* Switch to Terraform 0.12 SDK which is required for Terraform 0.12 support.
  This is the first release to use the 0.12 SDK required for Terraform 0.12
  support. Some provider behaviour may have changed as a result of changes made
  by the new SDK version. ([#760](https://github.com/hashicorp/terraform-provider-vsphere/pull/760)) 

## 1.10.0 (March 15, 2019)

FEATURES:

* **New Data Source:** `vsphere_folder` ([#709](https://github.com/hashicorp/terraform-provider-vsphere/pull/709))

IMPROVEMENTS:

* Update tf-vsphere-devrc.mk.example to include all environment variables ([#707](https://github.com/hashicorp/terraform-provider-vsphere/pull/707))
* Add Go Modules support ([#705](https://github.com/hashicorp/terraform-provider-vsphere/pull/705))
* Fix assorted typos in documentation
* `resource/virtual_machine`: Add support for using guest.ipAddress for older
  versions of VM Tools. ([#684](https://github.com/hashicorp/terraform-provider-vsphere/issues/684))

BUG FIXES:

* `resource/virtual_machine`: Do not set optional `ignored_guest_ips` on read ([#726](https://github.com/hashicorp/terraform-provider-vsphere/pull/726))

## 1.9.1 (January 10, 2019)

IMPROVEMENTS:

* `resource/virtual_machine`: Increase logging after old config expansion during
  diff checking ([#661](https://github.com/hashicorp/terraform-provider-vsphere/issues/661))
* `resource/virtual_machine`: Unlock `memory_reservation` from maximum when 
  `memory_reservation` is not equal to `memory`. ([#680](https://github.com/hashicorp/terraform-provider-vsphere/issues/680))

BUG FIXES:

* `resource/virtual_machine`: Return zero instead of nil for memory allocation
  and reservation values ([#655](https://github.com/hashicorp/terraform-provider-vsphere/issues/655))
* Ignore nil interfaces when converting a slice of interfaces into a slice
  of strings. ([#666](https://github.com/hashicorp/terraform-provider-vsphere/issues/666))
* `resource/virtual_machine`: Use schema for `properties` elem definition in `vapp` schema. ([#678](https://github.com/hashicorp/terraform-provider-vsphere/issues/678))

## 1.9.0 (October 31, 2018)

FEATURES:

* **New Resource:** `vsphere_vapp_entity` ([#640](https://github.com/hashicorp/terraform-provider-vsphere/issues/640))
* `resource/host_virtual_switch`: Add support for importing ([#625](https://github.com/hashicorp/terraform-provider-vsphere/issues/625))

IMPROVEMENTS:

* `resource/virtual_disk`: Update existing and add additional tests ([#635](https://github.com/hashicorp/terraform-provider-vsphere/issues/635))

BUG FIXES:

* `resource/virtual_disk`: Ignore "already exists" errors when creating
  directories on vSAN. ([#639](https://github.com/hashicorp/terraform-provider-vsphere/issues/639))
* Find tag changes when first tag is changed. ([#632](https://github.com/hashicorp/terraform-provider-vsphere/issues/632))
* `resource/virtual_machine`: Do not ForceNew when clone `timeout` is changed. ([#631](https://github.com/hashicorp/terraform-provider-vsphere/issues/631))
* `resource/virtual_machine_snapshot`: Raise error on snapshot create task
  error. ([#628](https://github.com/hashicorp/terraform-provider-vsphere/pull/628))

## 1.8.1 (September 11, 2018)

IMPROVEMENTS:

* `data/vapp_container`: Re-add `data_source_vapp_container`. ([#617](https://github.com/hashicorp/terraform-provider-vsphere/issues/617))

## 1.8.0 (September 10, 2018)

FEATURES:

* **New Data Source:** `vsphere_vapp_container` ([#610](https://github.com/hashicorp/terraform-provider-vsphere/issues/610))

BUG FIXES:

* `resource/virtual_machine`: Only relocate after create if `host_system_id` is
  set and does not match host the VM currently resides on. ([#609](https://github.com/hashicorp/terraform-provider-vsphere/issues/609))
* `resource/compute_cluster`: Return empty policy instead of trying to read
  `nil` variable when `ha_admission_control_policy` is set to `disabled`. ([#611](https://github.com/hashicorp/terraform-provider-vsphere/issues/611))
* `resource/virtual_machine`: Skip reading latency sensitivity parameters when
  LatencySensitivity is `nil`. ([#612](https://github.com/hashicorp/terraform-provider-vsphere/issues/612))
* `resource/compute_cluster`: Unset ID when the resource is not found. ([#613](https://github.com/hashicorp/terraform-provider-vsphere/issues/613))
* `resource/virtual_machine`: Skip OS specific customization checks when 
  `resource_pool_id` is not set. ([#614](https://github.com/hashicorp/terraform-provider-vsphere/issues/614))

## 1.7.0 (August 24, 2018)

FEATURES:

* **New Resource:** `vsphere_vapp_container` ([#566](https://github.com/hashicorp/terraform-provider-vsphere/issues/566))
* `resource/vsphere_virtual_machine`: Added support for bus sharing on SCSI
  adapters. ([#574](https://github.com/hashicorp/terraform-provider-vsphere/issues/574))

IMPROVEMENTS:

* `resource/vsphere_datacenter`: Added `moid` to expose the managed object ID
  because the datacenter's name is currently being used as the `id`.
  ([#575](https://github.com/hashicorp/terraform-provider-vsphere/issues/575))
* `resource/vsphere_virtual_machine`: Check if relocation is necessary after
  creation. ([#583](https://github.com/hashicorp/terraform-provider-vsphere/issues/583))

BUG FIXES:

* `resource/vsphere_virtual_machine`: The resource no longer attempts to set
  ResourceAllocation on virtual ethernet cards when the vSphere version is under 6.0. ([#579](https://github.com/hashicorp/terraform-provider-vsphere/issues/579))
* `resource/vsphere_resource_pool`: The read function is now called at the end
  of resource creation.
  ([#560](https://github.com/hashicorp/terraform-provider-vsphere/issues/560))
* Updated govmomi to v0.18. ([#600](https://github.com/hashicorp/terraform-provider-vsphere/issues/600))

## 1.6.0 (May 31, 2018)

FEATURES:

* **New Resource:** `vsphere_resource_pool` ([#535](https://github.com/hashicorp/terraform-provider-vsphere/issues/535))

IMPROVEMENTS:

* `data/vsphere_host`: Now exports the `resource_pool_id` attribute, which
  points to the root resource pool of either the standalone host, or the
  cluster's root resource pool in the event the host is a member of a cluster.
  ([#535](https://github.com/hashicorp/terraform-provider-vsphere/issues/535))

BUG FIXES:

* `resource/vsphere_virtual_machine`: Scenarios that force a new resource will
  no longer create diff mismatches when external disks are attached with the
  `attach` parameter. ([#528](https://github.com/hashicorp/terraform-provider-vsphere/issues/528))

## 1.5.0 (May 11, 2018)

FEATURES:

* **New Data Source:** `vsphere_compute_cluster` ([#492](https://github.com/hashicorp/terraform-provider-vsphere/issues/492))
* **New Resource:** `vsphere_compute_cluster` ([#487](https://github.com/hashicorp/terraform-provider-vsphere/issues/487))
* **New Resource:** `vsphere_drs_vm_override` ([#498](https://github.com/hashicorp/terraform-provider-vsphere/issues/498))
* **New Resource:** `vsphere_ha_vm_override` ([#501](https://github.com/hashicorp/terraform-provider-vsphere/issues/501))
* **New Resource:** `vsphere_dpm_host_override` ([#503](https://github.com/hashicorp/terraform-provider-vsphere/issues/503))
* **New Resource:** `vsphere_compute_cluster_vm_group` ([#506](https://github.com/hashicorp/terraform-provider-vsphere/issues/506))
* **New Resource:** `vsphere_compute_cluster_host_group` ([#508](https://github.com/hashicorp/terraform-provider-vsphere/issues/508))
* **New Resource:** `vsphere_compute_cluster_vm_host_rule` ([#511](https://github.com/hashicorp/terraform-provider-vsphere/issues/511))
* **New Resource:** `vsphere_compute_cluster_vm_dependency_rule` ([#513](https://github.com/hashicorp/terraform-provider-vsphere/issues/513))
* **New Resource:** `vsphere_compute_cluster_vm_affinity_rule` ([#515](https://github.com/hashicorp/terraform-provider-vsphere/issues/515))
* **New Resource:** `vsphere_compute_cluster_vm_anti_affinity_rule` ([#515](https://github.com/hashicorp/terraform-provider-vsphere/issues/515))
* **New Resource:** `vsphere_datastore_cluster_vm_anti_affinity_rule` ([#520](https://github.com/hashicorp/terraform-provider-vsphere/issues/520))

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: Exposed `latency_sensitivity`, which can
  be used to adjust the scheduling priority of the virtual machine for
  low-latency applications. ([#490](https://github.com/hashicorp/terraform-provider-vsphere/issues/490))
* `resource/vsphere_virtual_disk`: Introduced the `create_directories` setting,
  which tells this resource to create any parent directories in the VMDK path.
  ([#512](https://github.com/hashicorp/terraform-provider-vsphere/issues/512))

## 1.4.1 (April 23, 2018)

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: Introduced the
  `wait_for_guest_net_routable` setting, which controls whether or not the guest
  network waiter waits on an address that matches the virtual machine's
  configured default gateway. ([#470](https://github.com/hashicorp/terraform-provider-vsphere/issues/470))

BUG FIXES:

* `resource/vsphere_virtual_machine`: The resource now correctly blocks `clone`
  workflows on direct ESXi connections, where cloning is not supported. ([#476](https://github.com/hashicorp/terraform-provider-vsphere/issues/476))
* `resource/vsphere_virtual_machine`: Corrected an issue that was preventing VMs
  from being migrated from one cluster to another. ([#474](https://github.com/hashicorp/terraform-provider-vsphere/issues/474))
* `resource/vsphere_virtual_machine`: Corrected an issue where changing
  datastore information and cloning/customization parameters (which forces a new
  resource) at the same time was creating a diff mismatch after destroying the
  old virtual machine. ([#469](https://github.com/hashicorp/terraform-provider-vsphere/issues/469))
* `resource/vsphere_virtual_machine`: Corrected a crash that can come up from an
  incomplete lookup of network information during network device management.
  ([#456](https://github.com/hashicorp/terraform-provider-vsphere/issues/456))
* `resource/vsphere_virtual_machine`: Corrected some issues where some
  post-clone configuration errors were leaving the resource half-completed and
  irrecoverable without direct modification of the state. ([#467](https://github.com/hashicorp/terraform-provider-vsphere/issues/467))
* `resource/vsphere_virtual_machine`: Corrected a crash that can come up when a
  retrieved virtual machine has no lower-level configuration object in the API.
  ([#463](https://github.com/hashicorp/terraform-provider-vsphere/issues/463))
* `resource/vsphere_virtual_machine`: Fixed an issue where disk sub-resource
  configurations were not being checked for newly created disks.
  ([#481](https://github.com/hashicorp/terraform-provider-vsphere/issues/481))

## 1.4.0 (April 10, 2018)

FEATURES:

* **New Resource:** `vsphere_storage_drs_vm_override` ([#450](https://github.com/hashicorp/terraform-provider-vsphere/issues/450))
* **New Resource:** `vsphere_datastore_cluster` ([#436](https://github.com/hashicorp/terraform-provider-vsphere/issues/436))
* **New Data Source:** `vsphere_datastore_cluster` ([#437](https://github.com/hashicorp/terraform-provider-vsphere/issues/437))

IMPROVEMENTS:

* The provider now has the ability to persist sessions to disk, which can help
  when running large amounts of consecutive or concurrent Terraform operations
  at once. See the [provider
  documentation](https://www.terraform.io/docs/providers/vsphere/index.html) for
  more details. ([#422](https://github.com/hashicorp/terraform-provider-vsphere/issues/422))
* `resource/vsphere_virtual_machine`: This resource now supports import of
  resources or migrations from legacy versions of the provider (provider version
  0.4.2 or earlier) into configurations that have the `clone` block specified.
  See [Additional requirements and notes for
  importing](https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html#additional-requirements-and-notes-for-importing)
  in the resource documentation for more details. ([#460](https://github.com/hashicorp/terraform-provider-vsphere/issues/460))
* `resource/vsphere_virtual_machine`: Now supports datastore clusters. Virtual
  machines placed in a datastore cluster will use Storage DRS recommendations
  for initial placement, virtual disk creation, and migration between datastore
  clusters. Migrations made by Storage DRS outside of Terraform will no longer
  create diffs when datastore clusters are in use. ([#447](https://github.com/hashicorp/terraform-provider-vsphere/issues/447))
* `resource/vsphere_virtual_machine`: Added support for ISO transport of vApp
  properties. The resource should now behave better with virtual machines cloned
  from OVF/OVA templates that use the ISO transport to supply configuration
  settings. ([#381](https://github.com/hashicorp/terraform-provider-vsphere/issues/381))
* `resource/vsphere_virtual_machine`: Added support for client mapped CDROM
  devices. ([#421](https://github.com/hashicorp/terraform-provider-vsphere/issues/421))
* `resource/vsphere_virtual_machine`: Destroying a VM that currently has
  external disks attached should now function correctly and not give a duplicate
  UUID error. ([#442](https://github.com/hashicorp/terraform-provider-vsphere/issues/442))
* `resource/vsphere_nas_datastore`: Now supports datastore clusters. ([#439](https://github.com/hashicorp/terraform-provider-vsphere/issues/439))
* `resource/vsphere_vmfs_datastore`: Now supports datastore clusters. ([#439](https://github.com/hashicorp/terraform-provider-vsphere/issues/439))

## 1.3.3 (March 01, 2018)

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: The `moid` attribute has now be re-added
  to the resource, exporting the managed object ID of the virtual machine.
  ([#390](https://github.com/hashicorp/terraform-provider-vsphere/issues/390))

BUG FIXES:

* `resource/vsphere_virtual_machine`: Fixed a crash scenario that can happen
  when a virtual machine is deployed to a cluster that does not have any hosts,
  or under certain circumstances such an expired vCenter license. ([#414](https://github.com/hashicorp/terraform-provider-vsphere/issues/414))
* `resource/vsphere_virtual_machine`: Corrected an issue reading disk capacity
  values after a vCenter or ESXi upgrade. ([#405](https://github.com/hashicorp/terraform-provider-vsphere/issues/405))
* `resource/vsphere_virtual_machine`: Opaque networks, such as those coming from
  NSX, should now be able to be correctly added as networks for virtual
  machines. ([#398](https://github.com/hashicorp/terraform-provider-vsphere/issues/398))

## 1.3.2 (February 07, 2018)

BUG FIXES:

* `resource/vsphere_virtual_machine`: Changed the update implemented in ([#377](https://github.com/hashicorp/terraform-provider-vsphere/issues/377))
  to use a local filter implementation. This corrects situations where virtual
  machines in inventory with orphaned or otherwise corrupt configurations were
  interfering with UUID searches, creating erroneous duplicate UUID errors. This
  fix applies to vSphere 6.0 and lower only. vSphere 6.5 was not affected.
  ([#391](https://github.com/hashicorp/terraform-provider-vsphere/issues/391))

## 1.3.1 (February 01, 2018)

BUG FIXES:

* `resource/vsphere_virtual_machine`: Looking up templates by their UUID now
  functions correctly for vSphere 6.0 and earlier. ([#377](https://github.com/hashicorp/terraform-provider-vsphere/issues/377))

## 1.3.0 (January 26, 2018)

BREAKING CHANGES:

* The `vsphere_virtual_machine` resource now has a new method of identifying
  virtual disk sub-resources, via the `label` attribute. This replaces the
  `name` attribute, which has now been marked as deprecated and will be removed
  in the next major version (2.0.0). Further to this, there is a `path`
  attribute that now must also be supplied for external disks. This has lifted
  several virtual disk-related cloning and migration restrictions, in addition
  to changing requirements for importing. See the [resource
  documentation](https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html)
  for usage details.

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: Fixed an issue where certain changes
  happening at the same time (such as a disk resize along with a change of SCSI
  type) were resulting in invalid device change operations. ([#371](https://github.com/hashicorp/terraform-provider-vsphere/issues/371))
* `resource/vsphere_virtual_machine`: Introduced the `label` argument, which
  allows one to address a virtual disk independent of its VMDK file name and
  position on the SCSI bus. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Introduced the `path` argument, which
  replaces the `name` attribute for supplying the path for externally attached
  disks supplied with `attach = true`, and is otherwise a computed attribute
  pointing to the current path of any specific virtual disk. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Introduced the `uuid` attribute, a new
  computed attribute that allows the tracking of a disk independent of its
  current position on the SCSI bus. This is used in all scenarios aside from
  freshly-created or added virtual disks. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: The virtual disk `name` argument is now
  deprecated and will be removed from future releases. It no longer dictates the
  name of non-attached VMDK files and serves as an alias to the now-split `label`
  and `path` attributes. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Cloning no longer requires you to choose a
  disk label (name) that matches the name of the VM. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Storage vMotion can now be performed on
  renamed virtual machines. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Storage vMotion no longer cares what your
  disks are labeled (named), and will not block migrations based on the naming
  criteria added after 1.1.1. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Storage vMotion now works on linked
  clones. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: The import restrictions for virtual disks
  have changed, and rather than ensuring that disk `name` arguments match a
  certain convention, `label` is now expected to match a convention of `diskN`,
  where N is the disk number, ordered by the disk's position on the SCSI bus.
  Importing to a configuration still using `name` to address disks is no longer
  supported. ([#363](https://github.com/hashicorp/terraform-provider-vsphere/issues/363))
* `resource/vsphere_virtual_machine`: Now supports setting vApp properties that
  usually come from an OVF/OVA template or virtual appliance. ([#303](https://github.com/hashicorp/terraform-provider-vsphere/issues/303))

## 1.2.0 (January 11, 2018)

FEATURES:

* **New Resource:** `vsphere_custom_attribute` ([#229](https://github.com/hashicorp/terraform-provider-vsphere/issues/229))
* **New Data Source:** `vsphere_custom_attribute` ([#229](https://github.com/hashicorp/terraform-provider-vsphere/issues/229))

IMPROVEMENTS:

* All vSphere provider resources that are capable of doing so now support custom
  attributes. Check the documentation of any specific resource for more details!
  ([#229](https://github.com/hashicorp/terraform-provider-vsphere/issues/229))
* `resource/vsphere_virtual_machine`: The resource will now disallow a disk's
  `name` coming from a value that is still unavailable at plan time (such as a
  computed value from a resource). ([#329](https://github.com/hashicorp/terraform-provider-vsphere/issues/329))

BUG FIXES:

* `resource/vsphere_virtual_machine`: Fixed an issue that was causing crashes
  when working with virtual machines or templates when no network interface was
  occupying the first available device slot on the PCI bus. ([#344](https://github.com/hashicorp/terraform-provider-vsphere/issues/344))

## 1.1.1 (December 14, 2017)

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: Network interface resource allocation
  options are now restricted to vSphere 6.0 and higher, as they are unsupported
  on vSphere 5.5. ([#322](https://github.com/hashicorp/terraform-provider-vsphere/issues/322))
* `resource/vsphere_virtual_machine`: Resources that were deleted outside of
  Terraform will now be marked as gone in the state, causing them to be
  re-created during the next apply. ([#321](https://github.com/hashicorp/terraform-provider-vsphere/issues/321))
* `resource/vsphere_virtual_machine`: Added some restrictions to storage vMotion
  to cover some currently un-supported scenarios that were still allowed,
  leading to potentially dangerous situations or invalid post-application
  states. ([#319](https://github.com/hashicorp/terraform-provider-vsphere/issues/319))
* `resource/vsphere_virtual_machine`: The resource now treats disks that it does
  not recognize at a known device address as orphaned, and will set
  `keep_on_remove` to safely remove them. ([#317](https://github.com/hashicorp/terraform-provider-vsphere/issues/317))
* `resource/vsphere_virtual_machine`: The resource now attempts to detect unsafe
  disk deletion scenarios that can happen from the renaming of a virtual machine
  in situations where the VM and disk names may share a common variable. The
  provider will block such operations from proceeding. ([#305](https://github.com/hashicorp/terraform-provider-vsphere/issues/305))

## 1.1.0 (December 07, 2017)

BREAKING CHANGES:

* The `vsphere_virtual_machine` _data source_ has a new sub-resource attribute
  for disk information, named `disks`. This takes the place of `disk_sizes`,
  which has been moved to a `size` attribute within this new sub-resource, and
  also contains information about the discovered disks' `eagerly_scrub` and
  `thin_provisioned` settings. This is to facilitate the ability to discover all
  settings that could cause issues when cloning virtual machines.

To transition to the new syntax, any `disk` sub-resource in a
`vsphere_virtual_machine` resource that depends on a syntax such as:

```
resource "vsphere_virtual_machine" "vm" {
  ...

  disk {
    name = "terraform-test.vmdk"
    size = "${data.vsphere_virtual_machine.template.disk_sizes[0]}"
  }
}
```

Should be changed to:

```
resource "vsphere_virtual_machine" "vm" {
  ...

  disk {
    name = "terraform-test.vmdk"
    size = "${data.vsphere_virtual_machine.template.disks.0.size}"
  }
}
```

If you are using `linked_clone`, add the new settings for `eagerly_scrub` and
`thin_provisioned`:

```
resource "vsphere_virtual_machine" "vm" {
  ...

  disk {
    name             = "terraform-test.vmdk"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }
}
```

For a more complete example, see the [cloning and customization
example](https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html#cloning-and-customization-example)
in the documentation.

BUG FIXES:

* `resource/vsphere_virtual_machine`: Fixed a bug with NIC device assignment
  logic that was causing a crash when adding more than 3 NICs to a VM. ([#280](https://github.com/hashicorp/terraform-provider-vsphere/issues/280))
* `resource/vsphere_virtual_machine`: CDROM devices on cloned virtual machines
  are now connected properly on power on. ([#278](https://github.com/hashicorp/terraform-provider-vsphere/issues/278))
* `resource/vsphere_virtual_machine`: Tightened the pre-clone checks for virtual
  disks to ensure that the size and disk types are the same between the template
  and the created virtual machine's configuration. ([#277](https://github.com/hashicorp/terraform-provider-vsphere/issues/277))

## 1.0.3 (December 06, 2017)

BUG FIXES:

* `resource/vsphere_virtual_machine`: Fixed an issue in the post-clone process
  when a CDROM device exists in configuration. ([#276](https://github.com/hashicorp/terraform-provider-vsphere/issues/276))

## 1.0.2 (December 05, 2017)

BUG FIXES:

* `resource/vsphere_virtual_machine`: Fixed issues related to correct processing
  VM templates with no network interfaces, or fewer network interfaces than the
  amount that will ultimately end up in configuration. ([#269](https://github.com/hashicorp/terraform-provider-vsphere/issues/269))
* `resource/vsphere_virtual_machine`: Version comparison logic now functions
  correctly to properly disable certain features when using older versions of
  vSphere. ([#272](https://github.com/hashicorp/terraform-provider-vsphere/issues/272))

## 1.0.1 (December 02, 2017)

BUG FIXES:

* `resource/vsphere_virtual_machine`: Corrected an issue that was preventing the
  use of this resource on standalone ESXi. ([#263](https://github.com/hashicorp/terraform-provider-vsphere/issues/263))
* `data/vsphere_resource_pool`: This data source now works as documented on
  standalone ESXi. ([#263](https://github.com/hashicorp/terraform-provider-vsphere/issues/263))

## 1.0.0 (December 01, 2017)

BREAKING CHANGES:

* The `vsphere_virtual_machine` resource has received a major update and change
  to its interface. See the documentation for the resource for full details,
  including information on things to consider while migrating the new version of
  the resource.

FEATURES:

* **New Data Source:** `vsphere_resource_pool` ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* **New Data Source:** `vsphere_datastore` ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* **New Data Source:** `vsphere_virtual_machine` ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))

IMPROVEMENTS:

* `resource/vsphere_virtual_machine`: The distinct VM workflows are now better
  defined: all cloning options are now contained within a `clone` sub-resource,
  with customization being a `customize` sub-resource off of that. Absence of
  the `clone` sub-resource means no cloning or customization will occur.
  ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: Nearly all customization options have now
  been exposed. Magic values such as hostname and DNS defaults have been
  removed, with some of these options now being required values depending on the
  OS being customized. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: Device management workflows have been
  greatly improved, exposing more options and fixing several bugs. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: Added support for CPU and memory hot-plug.
  Several other VM reconfiguration operations are also supported while the VM is
  powered on, guest type and VMware tools permitting in some cases. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: The resource now supports both host and
  storage vMotion. Virtual machines can now be moved between hosts, clusters,
  resource pools, and datastores. Individual disks can be pinned to a single
  datastore with a VM located in another. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: The resource now supports import. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_virtual_machine`: Several other minor improvements, see
  documentation for more details. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))

BUG FIXES:

* `resource/vsphere_virtual_machine`: Several long-standing issues have been fixed,
  namely surrounding virtual disk and network device management. ([#244](https://github.com/hashicorp/terraform-provider-vsphere/issues/244))
* `resource/vsphere_host_virtual_switch`: This resource now correctly supports a
  configuration with no NICs. ([#256](https://github.com/hashicorp/terraform-provider-vsphere/issues/256))
* `data/vsphere_network`: No longer restricted to being used on vCenter. ([#248](https://github.com/hashicorp/terraform-provider-vsphere/issues/248))

## 0.4.2 (October 13, 2017)

FEATURES:

* **New Data Source:** `vsphere_network` ([#201](https://github.com/hashicorp/terraform-provider-vsphere/issues/201))
* **New Data Source:** `vsphere_distributed_virtual_switch` ([#170](https://github.com/hashicorp/terraform-provider-vsphere/issues/170))
* **New Resource:** `vsphere_distributed_port_group` ([#189](https://github.com/hashicorp/terraform-provider-vsphere/issues/189))
* **New Resource:** `vsphere_distributed_virtual_switch` ([#188](https://github.com/hashicorp/terraform-provider-vsphere/issues/188))

IMPROVEMENTS:

* resource/vsphere_virtual_machine: The customization waiter is now tunable
  through the `wait_for_customization_timeout` argument. The timeout can be
  adjusted or the waiter can be disabled altogether. ([#199](https://github.com/hashicorp/terraform-provider-vsphere/issues/199))
* resource/vsphere_virtual_machine: `domain` now acts as a default for
  `dns_suffixes` if the latter is not defined, setting the value in `domain` as
  a search domain in the customization specification. `vsphere.local` is not
  used as a last resort only. ([#185](https://github.com/hashicorp/terraform-provider-vsphere/issues/185))
* resource/vsphere_virtual_machine: Expose the `adapter_type` parameter to allow
  the control of the network interface type. This is currently restricted to
  `vmxnet3` and `e1000` but offers more control than what was available before,
  and more interface types will follow in later versions of the provider.
  ([#193](https://github.com/hashicorp/terraform-provider-vsphere/issues/193))

BUG FIXES:

* resource/vsphere_virtual_machine: Fixed a regression with network discovery
  that was causing Terraform to crash while the VM was in a powered off state.
  ([#198](https://github.com/hashicorp/terraform-provider-vsphere/issues/198))
* All resources that can use tags will now properly remove their tags completely
  (or remove any out-of-band added tags) when the `tags` argument is not present
  in configuration. ([#196](https://github.com/hashicorp/terraform-provider-vsphere/issues/196))

## 0.4.1 (October 02, 2017)

BUG FIXES:

* resource/vsphere_folder: Migration of state from a version of this resource
  before v0.4.0 now works correctly. ([#187](https://github.com/hashicorp/terraform-provider-vsphere/issues/187))

## 0.4.0 (September 29, 2017)

BREAKING CHANGES:

* The `vsphere_folder` resource has been re-written, and its configuration is
  significantly different. See the [resource
  documentation](https://www.terraform.io/docs/providers/vsphere/r/folder.html)
  for more details. Existing state will be migrated. ([#179](https://github.com/hashicorp/terraform-provider-vsphere/issues/179))

FEATURES:

* **New Data Source:** `vsphere_tag` ([#171](https://github.com/hashicorp/terraform-provider-vsphere/issues/171))
* **New Data Source:** `vsphere_tag_category` ([#167](https://github.com/hashicorp/terraform-provider-vsphere/issues/167))
* **New Resoruce:** `vsphere_tag` ([#171](https://github.com/hashicorp/terraform-provider-vsphere/issues/171))
* **New Resoruce:** `vsphere_tag_category` ([#164](https://github.com/hashicorp/terraform-provider-vsphere/issues/164))

IMPROVEMENTS:

* resource/vsphere_folder: You can now create any kind of folder with this
  resource, not just virtual machine folders. ([#179](https://github.com/hashicorp/terraform-provider-vsphere/issues/179))
* resource/vsphere_folder: Now supports tags. ([#179](https://github.com/hashicorp/terraform-provider-vsphere/issues/179))
* resource/vsphere_folder: Now supports import. ([#179](https://github.com/hashicorp/terraform-provider-vsphere/issues/179))
* resource/vsphere_datacenter: Tags can now be applied to datacenters. ([#177](https://github.com/hashicorp/terraform-provider-vsphere/issues/177))
* resource/vsphere_nas_datastore: Tags can now be applied to NAS datastores.
  ([#176](https://github.com/hashicorp/terraform-provider-vsphere/issues/176))
* resource/vsphere_vmfs_datastore: Tags can now be applied to VMFS datastores.
  ([#176](https://github.com/hashicorp/terraform-provider-vsphere/issues/176))
* resource/vsphere_virtual_machine: Tags can now be applied to virtual machines.
  ([#175](https://github.com/hashicorp/terraform-provider-vsphere/issues/175))
* resource/vsphere_virtual_machine: Adjusted the customization timeout to 10
  minutes ([#168](https://github.com/hashicorp/terraform-provider-vsphere/issues/168))

BUG FIXES:

* resource/vsphere_virtual_machine: This resource can now be used with networks
  with unescaped slashes in its network name. ([#181](https://github.com/hashicorp/terraform-provider-vsphere/issues/181))
* resource/vsphere_virtual_machine: Fixed a crash where virtual NICs were
  created with networks backed by a 3rd party hardware VDS. ([#181](https://github.com/hashicorp/terraform-provider-vsphere/issues/181))
* resource/vsphere_virtual_machine: Fixed crashes and spurious diffs that were
  caused by errors in the code that associates the default gateway with its
  correct network device during refresh. ([#180](https://github.com/hashicorp/terraform-provider-vsphere/issues/180))

## 0.3.0 (September 14, 2017)

BREAKING CHANGES:

* `vsphere_virtual_machine` now waits on a _routeable_ IP address by default,
  and does not wait when running `terraform plan`, `terraform refresh`, or
  `terraform destroy`. There is also now a timeout of 5 minutes, after which
  `terraform apply` will fail with an error. Note that the apply may not fail
  exactly on the 5 minute mark. The network waiter can be disabled completely by
  setting `wait_for_guest_net` to `false`. ([#158](https://github.com/hashicorp/terraform-provider-vsphere/issues/158))

FEATURES:

* **New Resource:** `vsphere_virtual_machine_snapshot` ([#107](https://github.com/hashicorp/terraform-provider-vsphere/issues/107))

IMPROVEMENTS:

* resource/vsphere_virtual_machine: Virtual machine power state is now enforced.
  Terraform will trigger a diff if the VM is powered off or suspended, and power
  it back on during the next apply. ([#152](https://github.com/hashicorp/terraform-provider-vsphere/issues/152))

BUG FIXES:

* resource/vsphere_virtual_machine: Fixed customization behavior to watch
  customization events for success, rather than returning immediately when the
  `CustomizeVM` task returns. This is especially important during Windows
  customization where a large part of the customization task involves
  out-of-band configuration through Sysprep. ([#158](https://github.com/hashicorp/terraform-provider-vsphere/issues/158))

## 0.2.2 (September 07, 2017)

FEATURES:

* **New Resource:** `vsphere_nas_datastore` ([#149](https://github.com/hashicorp/terraform-provider-vsphere/issues/149))
* **New Resource:** `vsphere_vmfs_datastore` ([#142](https://github.com/hashicorp/terraform-provider-vsphere/issues/142))
* **New Data Source:** `vsphere_vmfs_disks` ([#141](https://github.com/hashicorp/terraform-provider-vsphere/issues/141))

## 0.2.1 (August 31, 2017)

FEATURES:

* **New Resource:** `vsphere_host_port_group` ([#139](https://github.com/hashicorp/terraform-provider-vsphere/issues/139))
* **New Resource:** `vsphere_host_virtual_switch` ([#138](https://github.com/hashicorp/terraform-provider-vsphere/issues/138))
* **New Data Source:** `vsphere_datacenter` ([#144](https://github.com/hashicorp/terraform-provider-vsphere/issues/144))
* **New Data Source:** `vsphere_host` ([#146](https://github.com/hashicorp/terraform-provider-vsphere/issues/146))

IMPROVEMENTS:

* resource/vsphere_virtual_machine: Allow customization of hostname ([#79](https://github.com/hashicorp/terraform-provider-vsphere/issues/79))

BUG FIXES:

* resource/vsphere_virtual_machine: Fix IPv4 address mapping issues causing
  spurious diffs, in addition to IPv6 normalization issues that can lead to spurious
  diffs as well. ([#128](https://github.com/hashicorp/terraform-provider-vsphere/issues/128))

## 0.2.0 (August 23, 2017)

BREAKING CHANGES:

* resource/vsphere_virtual_disk: Default adapter type is now `lsiLogic`,
  changed from `ide`. ([#94](https://github.com/hashicorp/terraform-provider-vsphere/issues/94))

FEATURES:

* **New Resource:** `vsphere_datacenter` ([#126](https://github.com/hashicorp/terraform-provider-vsphere/issues/126))
* **New Resource:** `vsphere_license` ([#110](https://github.com/hashicorp/terraform-provider-vsphere/issues/110))

IMPROVEMENTS:

* resource/vsphere_virtual_machine: Add annotation argument ([#111](https://github.com/hashicorp/terraform-provider-vsphere/issues/111))

BUG FIXES:

* Updated [govmomi](https://github.com/vmware/govmomi) to v0.15.0 ([#114](https://github.com/hashicorp/terraform-provider-vsphere/issues/114))
* Updated network interface discovery behaviour in refresh. [[#129](https://github.com/hashicorp/terraform-provider-vsphere/issues/129)]. This fixes
  several reported bugs - see the PR for references!

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
