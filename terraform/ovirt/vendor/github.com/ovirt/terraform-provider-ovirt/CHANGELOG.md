## v4.4.5-r.0

This version is a feature release and **breaks backwards compatibility**. It contains two changes:

- [251: Introduce auto-pinning policy](https://github.com/oVirt/terraform-provider-ovirt/pull/251)
- [254: Add certificate verification](https://github.com/oVirt/terraform-provider-ovirt/pull/253)

Please check the above links for details on these changes. *This release is the first that matches the oVirt version it is built for.*

## v0.99.0

Dear community,

We are in the process of catching up on open pull requests. In order to make that happen we are now tagging the version before merging PR's as 0.99.

After this release we will be switching to a versioning system that aligns with the oVirt release it supports. If you have any questions please feel free to raise them as a GitHub issue.

## [v0.5.0](https://github.com/oVirt/terraform-provider-ovirt/tree/v0.5.0) (2021-01-07)

[Full Changelog](https://github.com/oVirt/terraform-provider-ovirt/compare/v0.4.2...v0.5.0)

**Implemented enhancements:**

- Request for New Data Source: ovirt\_template  [\#182](https://github.com/oVirt/terraform-provider-ovirt/issues/182)
- Get IP address of the new VM [\#172](https://github.com/oVirt/terraform-provider-ovirt/issues/172)
- Feature request: support for sysprep to initialize windows VMs [\#158](https://github.com/oVirt/terraform-provider-ovirt/issues/158)
- Drop vendor dir [\#148](https://github.com/oVirt/terraform-provider-ovirt/issues/148)
- Request for New Resource: ovirt\_domain [\#94](https://github.com/oVirt/terraform-provider-ovirt/issues/94)

**Fixed bugs:**

- Update for libgo.so.11-\>libgo.so.13? [\#125](https://github.com/oVirt/terraform-provider-ovirt/issues/125)

**Closed issues:**

- VM resource delete leaving disks [\#223](https://github.com/oVirt/terraform-provider-ovirt/issues/223)
- Could not locate Gemfile during `make website` [\#211](https://github.com/oVirt/terraform-provider-ovirt/issues/211)
- Cannot use ovirt\_image\_transfer with FCOS image [\#197](https://github.com/oVirt/terraform-provider-ovirt/issues/197)
- Support OS property on VM [\#191](https://github.com/oVirt/terraform-provider-ovirt/issues/191)
- vm boot devices [\#189](https://github.com/oVirt/terraform-provider-ovirt/issues/189)
- ovirt\_image\_transfer: The string '' isnt a valid value for the DiskFormat  [\#187](https://github.com/oVirt/terraform-provider-ovirt/issues/187)
- RESTEASY002010: Failed to execute: javax.ws.rs.WebApplicationException: HTTP 404 Not Found [\#180](https://github.com/oVirt/terraform-provider-ovirt/issues/180)
- Unable to set network on RHEV VM [\#179](https://github.com/oVirt/terraform-provider-ovirt/issues/179)

**Merged pull requests:**

- WIP: Add priority to affinity group resource [\#240](https://github.com/oVirt/terraform-provider-ovirt/pull/240) ([Gal-Zaidman](https://github.com/Gal-Zaidman))
- Fix throwing a 404 error if we cannot find a disk by ID [\#238](https://github.com/oVirt/terraform-provider-ovirt/pull/238) ([jake2184](https://github.com/jake2184))
- Add 'latest' as a flag to data\_source\_template [\#237](https://github.com/oVirt/terraform-provider-ovirt/pull/237) ([jake2184](https://github.com/jake2184))
- Revert "increase template and image transfer timeout" [\#234](https://github.com/oVirt/terraform-provider-ovirt/pull/234) ([emesika](https://github.com/emesika))
- increase template and image transfer timeout [\#233](https://github.com/oVirt/terraform-provider-ovirt/pull/233) ([Gal-Zaidman](https://github.com/Gal-Zaidman))
- moving docs into the new location the registry prefers [\#232](https://github.com/oVirt/terraform-provider-ovirt/pull/232) ([attachmentgenie](https://github.com/attachmentgenie))
- import of release workflow from provider scafolding [\#231](https://github.com/oVirt/terraform-provider-ovirt/pull/231) ([attachmentgenie](https://github.com/attachmentgenie))
- Add auto\_start flag to ovirt\_vm resource [\#227](https://github.com/oVirt/terraform-provider-ovirt/pull/227) ([Gal-Zaidman](https://github.com/Gal-Zaidman))
- Add Affinity Group resource [\#225](https://github.com/oVirt/terraform-provider-ovirt/pull/225) ([jake2184](https://github.com/jake2184))
- Set detachOnly to false when deleting template clones [\#224](https://github.com/oVirt/terraform-provider-ovirt/pull/224) ([jake2184](https://github.com/jake2184))
- Added static hardware address configuration for vm nics [\#221](https://github.com/oVirt/terraform-provider-ovirt/pull/221) ([s-all-kin](https://github.com/s-all-kin))
- Allow OS to be updated rather than recreating VM [\#219](https://github.com/oVirt/terraform-provider-ovirt/pull/219) ([jake2184](https://github.com/jake2184))
- Added option storage\_domain parameter. Related only for build VM from… [\#216](https://github.com/oVirt/terraform-provider-ovirt/pull/216) ([pavel-z1](https://github.com/pavel-z1))
- Disk Alias configuration during deploy VM from Template [\#214](https://github.com/oVirt/terraform-provider-ovirt/pull/214) ([pavel-z1](https://github.com/pavel-z1))
- image\_transfer: failed OPTIONS request will panic [\#206](https://github.com/oVirt/terraform-provider-ovirt/pull/206) ([rgolangh](https://github.com/rgolangh))
- Handle empty instance type Creating a Vm with empty instance type is allowed, and the resource read need to handle a case where its missing. [\#203](https://github.com/oVirt/terraform-provider-ovirt/pull/203) ([rgolangh](https://github.com/rgolangh))
- Additions to VM resource [\#202](https://github.com/oVirt/terraform-provider-ovirt/pull/202) ([rgolangh](https://github.com/rgolangh))
- Override disk attachment on VM resource create [\#201](https://github.com/oVirt/terraform-provider-ovirt/pull/201) ([rgolangh](https://github.com/rgolangh))
- Initialization additions to VM [\#200](https://github.com/oVirt/terraform-provider-ovirt/pull/200) ([rgolangh](https://github.com/rgolangh))
- Update VM parameters feature and VM running Status [\#198](https://github.com/oVirt/terraform-provider-ovirt/pull/198) ([pavel-z1](https://github.com/pavel-z1))
- Fix immediate EOF on reads for image transfers via http\(s\) [\#195](https://github.com/oVirt/terraform-provider-ovirt/pull/195) ([r0ci](https://github.com/r0ci))
- implement boot devices for vms [\#194](https://github.com/oVirt/terraform-provider-ovirt/pull/194) ([mathianasj](https://github.com/mathianasj))
- Update go-ovirt dependency to v4.4.1 [\#190](https://github.com/oVirt/terraform-provider-ovirt/pull/190) ([imjoey](https://github.com/imjoey))
- Migrate to terraform-plugin-sdk [\#188](https://github.com/oVirt/terraform-provider-ovirt/pull/188) ([LorbusChris](https://github.com/LorbusChris))
- Update main.tf [\#185](https://github.com/oVirt/terraform-provider-ovirt/pull/185) ([Oddly](https://github.com/Oddly))
- Add Image Transfer resources [\#184](https://github.com/oVirt/terraform-provider-ovirt/pull/184) ([rgolangh](https://github.com/rgolangh))
- Add support for data source ovirt\_template [\#183](https://github.com/oVirt/terraform-provider-ovirt/pull/183) ([imjoey](https://github.com/imjoey))
- doc: Improve the expressions [\#181](https://github.com/oVirt/terraform-provider-ovirt/pull/181) ([imjoey](https://github.com/imjoey))
- docs: Update to using go version 1.12+ [\#178](https://github.com/oVirt/terraform-provider-ovirt/pull/178) ([imjoey](https://github.com/imjoey))
- Replaces apache thrift module source [\#177](https://github.com/oVirt/terraform-provider-ovirt/pull/177) ([beremtl](https://github.com/beremtl))

## 0.4.2 (Sep 10, 2019)

FEATURES:

* **New Data Source:** `ovirt_nics` ([#173](https://github.com/oVirt/terraform-provider-ovirt/pull/173))

IMPROVEMENTS:

* resource/ovirt_datacenter: Make the status field exportable ([#170](https://github.com/oVirt/terraform-provider-ovirt/pull/170))
* data/ovirt_vms: Export IP configurations of VM ([#174](https://github.com/oVirt/terraform-provider-ovirt/pull/174))

## 0.4.1 (Jul 31, 2019)

BUG FIXES:

* resource/ovirt_vm: Do not try to start a VM after updating attributes ([#167](https://github.com/oVirt/terraform-provider-ovirt/pull/167))
* resource/ovirt_disk_attachment: Fix failed to check if a disk attachment exists ([#162](https://github.com/oVirt/terraform-provider-ovirt/pull/162))

FEATURES:

* **New Resource:** `ovirt_snapshot` ([#157](https://github.com/oVirt/terraform-provider-ovirt/pull/157))

IMPROVEMENTS:

* doc: Format inline HCL codes in docs ([#164](https://github.com/oVirt/terraform-provider-ovirt/pull/164))
* provider: Add more general method for parsing composite resource ID ([#163](https://github.com/oVirt/terraform-provider-ovirt/pull/163))
* provider: Format the HCL codes definied in acceptance tests ([#160](https://github.com/oVirt/terraform-provider-ovirt/pull/160))

## 0.4.0 (Jul 8, 2019)

BACKWARDS INCOMPATIBILITIES / NOTES:

* provider: This is the first release since it has been transferred to oVirt community under incubation. Please access to the provider with the new ([oVirt/terraform-provider-ovirt](https://github.com/oVirt/terraform-provider-ovirt)).

IMPROVEMENTS:

* provider: Update to Terraform v0.12.2 ([#145](https://github.com/oVirt/terraform-provider-ovirt/pull/145))
* provider: Remove serveral unnecessary scripts in CI process ([#153](https://github.com/oVirt/terraform-provider-ovirt/pull/153))
* provider: Set `GOFLAGS` in CI environment to force `go mod` to use packages under vendor directory ([#155](https://github.com/oVirt/terraform-provider-ovirt/pull/155))

## 0.3.1 (Jun 10, 2019)

BUG FIXES:

* resource/ovirt_vm: Prevent reading VM failure in case of the `original_template` attribute is unavaliable ([#140](https://github.com/imjoey/terraform-provider-ovirt/pull/140))

FEATURES:

* **New Data Source:** `ovirt_hosts` ([#138](https://github.com/imjoey/terraform-provider-ovirt/pull/138))

IMPROVEMENTS:

* provider: Update to Terraform v0.12.1 ([#141](https://github.com/imjoey/terraform-provider-ovirt/pull/141))

## 0.3.0 (May 29, 2019)

BACKWARDS INCOMPATIBILITIES / NOTES:

* provider: This release contains only a Terraform SDK upgrade for compatibility with Terraform v0.12. The provider should remains backwards compatible with Terraform v0.11. This update should have no significant changes in behavior for the provider. Please report any unexpected behavior in new GitHub issues (Terraform oVirt Provider: https://github.com/imjoey/terraform-provider-ovirt/issues) ([#133](https://github.com/imjoey/terraform-provider-ovirt/pull/133))

## 0.2.2 (May 27, 2019)

BUG FIXES:

* resource/ovirt_vm: Prevent creating VM failure and mistaken state diffs due to `memory` attribute

## 0.2.1 (May 22, 2019)

FEATURES:

* **New Resource:** `ovirt_storage_domain` ([#92](https://github.com/imjoey/terraform-provider-ovirt/pull/92))
* **New Resource:** `ovirt_user` ([#98](https://github.com/imjoey/terraform-provider-ovirt/pull/98))
* **New Resource:** `ovirt_cluster` ([#103](https://github.com/imjoey/terraform-provider-ovirt/pull/103))
* **New Resource:** `ovirt_mac_pool` ([#107](https://github.com/imjoey/terraform-provider-ovirt/pull/107))
* **New Resource:** `ovirt_tag` ([#107](https://github.com/imjoey/terraform-provider-ovirt/pull/114))
* **New Resource:** `ovirt_host` ([#121](https://github.com/imjoey/terraform-provider-ovirt/pull/121))
* **New Data Source:** `ovirt_authzs` ([#97](https://github.com/imjoey/terraform-provider-ovirt/pull/97))
* **New Data Source:** `ovirt_users` ([#102](https://github.com/imjoey/terraform-provider-ovirt/pull/102))
* **New Data Source:** `ovirt_mac_pools` ([#109](https://github.com/imjoey/terraform-provider-ovirt/pull/109))
* **New Data Source:** `ovirt_vms` ([#118](https://github.com/imjoey/terraform-provider-ovirt/pull/118))

IMPROVEMENTS:

* provider: Add `header` params support for connection settings ([#72](https://github.com/imjoey/terraform-provider-ovirt/pull/72))
* resource/ovirt_disk: Add `quota_id` attribute support ([#80](https://github.com/imjoey/terraform-provider-ovirt/pull/80))
* doc: Add webswebsite infrastructure and provider documantations ([#81](https://github.com/imjoey/terraform-provider-ovirt/pull/81))
* resource/ovirt_vnic: Add acceptance tests ([#90](https://github.com/imjoey/terraform-provider-ovirt/pull/90))
* resource/ovirt_network: Add acceptance tests ([#91](https://github.com/imjoey/terraform-provider-ovirt/pull/91))
* resource/ovirt_vm: Add `clone` support ([#131](https://github.com/imjoey/terraform-provider-ovirt/pull/131))

## 0.2.0 (September 26, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* provider: All the new or existing resources and data sources have been refactored with the [oVirt Go SDK](https://github.com/imjoey/go-ovirt) to access the oVirt engine API

FEATURES:

* **New Resource:** `ovirt_disk_attachment` ([#1](https://github.com/imjoey/terraform-provider-ovirt/pull/1))
* **New Resource:** `ovirt_datacenter` ([#3](https://github.com/imjoey/terraform-provider-ovirt/pull/3))
* **New Resource:** `ovirt_network` ([#6](https://github.com/imjoey/terraform-provider-ovirt/pull/6))
* **New Resource:** `ovirt_vnic_profile` ([#41](https://github.com/imjoey/terraform-provider-ovirt/pull/41))
* **New Resource:** `ovirt_vnic` ([#56](https://github.com/imjoey/terraform-provider-ovirt/pull/56))
* **New Data Source:** `ovirt_datacenters` ([#4](https://github.com/imjoey/terraform-provider-ovirt/pull/4))
* **New Data Source:** `ovirt_networks` ([#13](https://github.com/imjoey/terraform-provider-ovirt/pull/13))
* **New Data Source:** `ovirt_clusters` ([#26](https://github.com/imjoey/terraform-provider-ovirt/pull/26))
* **New Data Source:** `ovirt_storagedomains` ([#27](https://github.com/imjoey/terraform-provider-ovirt/pull/27))
* **New Data Source:** `ovirt_vnic_profiles` ([#51](https://github.com/imjoey/terraform-provider-ovirt/pull/51))

IMPROVEMENTS:

* provider: Add GNU make integration: ([#15](https://github.com/imjoey/terraform-provider-ovirt/pull/15))
* provider: Add acceptance tests for provider ([#8](https://github.com/imjoey/terraform-provider-ovirt/pull/8))
* provider: Add acceptance tests for all the resources and data sources
* provider: Add travis CI support ([#47](https://github.com/imjoey/terraform-provider-ovirt/pull/47))
* provider: Add missing attributes and processing logic for the existing `ovirt_vm`, `ovirt_disk` resources and `ovirt_disk` data source defined in v0.1.0

## 0.1.0 (March 13, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Release by [EMSL-MSC](https://github.com/EMSL-MSC/terraform-provider-ovirt/commits) Orgnization, please see [here](https://github.com/EMSL-MSC/terraform-provider-ovirt/releases/tag/0.1.0) for details.

FEATURES:

* **New Resource:** `ovirt_vm`
* **New Resource:** `ovirt_disk`
* **New Data Source:** `ovirt_disk`
