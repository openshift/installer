## 1.46.0 (Unreleased)

FEATURES

* __New Resource__: `blockstorage_qos_v3` ([#1325](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1325))
* __New Resource__: `blockstorage_qos_association_v3` ([#1331](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1331))
* __New Data Source__: `blockstorage_quotaset_v3` ([#1319](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1319))
* __New Data Source__: `networking_quota_v2` ([#1318](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1318))

IMPROVEMENTS

* Added `region` argument to `compute_aggregate_v2` resource ([#1276](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1276))
* Fixed default `0` value in skipped arguments of `networking_quota_v2` resource ([#1316](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1316))
* Added `tags` to `lb_loadbalancer_v2` resource ([#1301](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1301))
* Use Otavia API for `lb_loadbalancer_v2` resource by default ([#1326](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1326))
* Updated `images_image_v2` resource to not recreate image `on min_disk_gb`, `min_ram_mb`, `protected` attributes changes ([#1299](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1299))
* Updated `gophercloud` to `v0.23.0` ([#1315](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1315))
* Updated `terraform-plugin-sdk` to `v2.10.0` ([#1333](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1333))

## 1.45.0 (4 November, 2021)

FEATURES

* __New Data Source__: `openstack_compute_quotaset_v2` ([#1302](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1302))

IMPROVEMENTS

* Added retries reading `dns_zone_v2` and `compute_instance_v2` state after creation in case of 502, 504 HTTP errors ([#1303](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1303))

BUG FIXES

* Improved removal of `networking_router_interface_v2` resource so it will delete only needed port on a router ([#1297](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1297))
* Flagged `url` attribute of `objectstorage_tempurl_v1` resource as sensitive ([#1305](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1305))
* Fixed not specified quota values are being set to 0 in `compute_quotaset_v2` resource ([#1304](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1304))

## 1.44.0 (2 October, 2021)

NOTES

* This release updates major version of `terraform-plugin-sdk` from `v1` to `v2` and that caused lots of changes in the code. If you experiencing new bugs after updating the provider please create an issue with a description of how to reproduce them.

FEATURES

* Updated `terraform-plugin-sdk` to `v2.7.1` ([#1139](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1139))
* Updated Go to `1.17` ([#1295](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1295))

## 1.43.1 (21 September, 2021)

BUG FIXES

* Fixed panics when a token doesn't have a project scope ([#1282](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1282))

## 1.43.0 (16 July, 2021)

FEATURES

* __New Resource__: `dns_transfer_request_v2` ([#1268](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1268))
* __New Resource__: `dns_transfer_accept_v2` ([#1268](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1268))

IMPROVEMENTS

* Added `SCTP`, `PROXYV2` protocols for `lb_pool_v2` resource ([#1251](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1251))
* Added `project_id` argument for `dns_recordset_v2` resource ([#1254](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1254))
* Added support for `shelved_offloaded` power state of `compute_instance_v2` resource ([#1259](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1259))
* Added `cidr` argument input check for `networking_subnet_v2` resource ([#1267](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1267))
* Removed Octavia microversions and added explanation about minor version usage ([#1249](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1249))
* Fixed `endpoints` argument for `vpnaas_endpoint_group_v2` resource in that way so endpoints order is not relevant anymore ([#1247](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1247))
* Added `addresses` argument for `db_instance_v1` resource ([#1260](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1260))
* Better formatted documentation for some resources and data sources ([#1252](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1252)), ([#1255](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1255)), ([#1256](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1256))
* Updated issues links so they point to the right repo ([#1272](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1272))

BUG FIXES

* Fixed `nil` panic in `compute_instance_v2` resource that could be caught while trying to unassign a server group from an instance ([#1248](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1248))

## 1.42.0 (12 May, 2021)

IMPROVEMENTS

* Added `SCTP` protocol support for `lb_listener_v2` resource, note that will work only in Octavia ([#1236](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1236))
* Added support for `HEALTHY` status of `db_instance_v1` resource ([#1241](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1241))
* Added `address_group` as `object_type` for `networking_rbac_policy_v2` resource ([#1243](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1243))
* * Updated `terraform-plugin-sdk` to `v1.17.2` ([#1244](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1244))

## 1.41.0 (23 April, 2021)

FEATURES

* __New Resource__: `blockstorage_volume_type_access_v3` ([#1223](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1223))

IMPROVEMENTS

* Added `disable_status_check` argument for `dns_recordset_v2` resource ([#1221](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1221))
* Added `availability_zone` argument for `lb_loadbalancer_v2` resource ([#1225](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1225))
* Added `backup` argument for `lb_members_v2` resource ([#1227](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1227))

## 1.40.0 (23 March, 2021)

FEATURES

* __New Resource__: `networking_portforwarding_v2` ([#940](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/940))
* __New Resource__: `blockstorage_volume_type_v3` ([#1204](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1204))

IMPROVEMENTS

* Go version is updated to `1.16` and we're providing `darwin/arm64` binaries starting from this release ([#1206](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1206))

BUG FIXES

* Fixed Bad request API error while updating `images_image_v2` resource because old OpenStack released don't have `hidden` argument ([#1209](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1209))
* Fixed Bad request API error while updating `blockstorage_quotaset_v2`, `blockstorage_quotaset_v3` ([#1200](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1200))

## 1.39.0 (6 March, 2021)

IMPROVEMENTS

* Added ability to manage `blockstorage_quotaset_v2` for the same project across several regions with a single resource ([#1182](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1182))
* Added ability to manage `blockstorage_quotaset_v3` for the same project across several regions with a single resource ([#1183](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1183))
* Added ability to manage `openstack_compute_quotaset_v2` for the same project across several regions with a single resource ([#1181](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1181))
* Added `volume_type_quota` argument for `blockstorage_quotaset_v2` resource ([#1187](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1187))
* Added `volume_type_quota` argument for `blockstorage_quotaset_v3` resource ([#1185](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1185))
* Added `hidden` argument for `openstack_images_image_v2` resource and datasource ([#1186](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1186))

BUG FIXES

* Fixed error updating `networking_quota_v2` when it was created with the version older than `1.38.0` ([#1180](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1180))

## 1.38.0 (February 24, 2021)

FEATURES

* __New Resource__: `openstack_lb_quota_v2` ([#1169](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1169))

IMPROVEMENTS

* Updated gophercloud/utils, which now recognizes `clouds.yml` in addition to `clouds.yaml` and correctly applies per-region value overrides ([#1172](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1172))
* Added `vip_port_id` for `lb_loadbalancer_v2` resource. It can be used only with Octavia ([#1164](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1164))
* Added `service_catalog` attribute for `identity_auth_scope_v3` data source ([#1167](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1167))
* Set `2.15` microversion for any type of `server_group_v2` policy except `affinity` and `anti-affinity` since they don't need any microversion ([#1141](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1141))
* Add a note about using names in `security_groups` in `compute_instance_v2` resource in docs ([#1178](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1178))
* Added ability to manage `networking_quota_v2` for the same project across several regions with a single resource ([#1177](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1177))

## 1.37.0 (February 8, 2021)

IMPROVEMENTS

* Added `image_source_username`, `image_source_password` arguments to `images_image_v2` resource ([#1157](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1157))
* Updated `networking_floatingip_v2` resource to retry subnets on floating IP creation, when a subnet is exhausted ([#1163](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1163))
* Updated security notices for sensitive arguments and attributes in documentation ([#1161](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1161))

BUG FIXES

* Fixed multiple `networking_router_v2` resource creation while using `external_subnet_ids` argument ([#1163](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1163))

## 1.36.0 (February 2, 2021)

NOTES

* The `dhcp_disabled` argument in `networking_subnet_v2` data source is deprecated. Use the `dhcp_enabled = false` argument value instead. ([#1153](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1153))
* The `max_retries` provider parameter now honors the `429` code and uses the `Retry-After` header to extend the retry function ([#1159](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1159))

FEATURES

* __New Resource__: `openstack_identity_user_membership_v3` ([#1149](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1149))
* __New Data Source__: `openstack_networking_subnet_ids_v2` ([#1153](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1153))

IMPROVEMENTS

* Updated `zone` argument to be `Optional` instead of `Required` in `compute_aggregate_v2` resource ([#1133](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1133))
* Updated local provider block in docs ([#1135](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1135))
* Updated Go version to `1.15` ([#1137](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1137))
* Updated `networking_router_v2` resource to retry external subnets on router creation, when a subnet is exhausted ([#1151](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1151))
* Added `subnets` attribute to `networking_network_v2` data source ([#1152](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1152))
* Extended `max_retries` provider parameter to use the `Retry-After` header ([#1159](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1159))

BUG FIXES

* Fixed copying `sync.Locker` by updating `gophercloud/utils` with the fix ([#1144](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1144))
* Fixed recreation of `lb_loadbalancer_v2` resource if `flavor_id` haven't been specified ([#1147](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1147))
* Fixed `networking_port_v2` resource update if `binding.profile` is not set ([#1154](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1154))

## 1.35.0 (January 15, 2021)

FEATURES

* __New Resource__: `openstack_compute_aggregate_v2` ([#1121](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1121))
* __New Data Source__: `openstack_compute_aggregate_v2` ([#1121](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1121))
* __New Data Source__: `openstack_compute_hypervisor_v2` ([#1126](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1126))

IMPROVEMENTS

* Added valid handling of the read-only `stores` property of the `images_image_v2` resource ([#1124](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1124))
* Added `image_id` argument for the `images_image_v2` resource ([#1125](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1125))
* Added `vendor_options.ignore_volume_confirmation` argument for the `compute_volume_attach_v2` resource to control whether to ignore volume status confirmation of the attached volume. ([#1127](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1127))
* Updated Gophercloud to `1.15.0` with utils package that now uses `imageservice` instead of `compute` to resolve image IDs ([#1128](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1128))

## 1.34.1 (December 21, 2020)

BUG FIXES

* Fixed an issue when empty a `flavor_id` argument in `compute_flavor_v2` resource could create plan changes ([#1120](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1120))

## 1.34.0 (December 20, 2020)

IMPROVEMENTS

* Added `flavor_id` to `compute_flavor_v2` resource creation options ([#1107](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1107))
* Updated `compute_flavor_v2` resource docs with `ephemeral` argument ([#1113](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1113))
* Updated `compute_instance_v2` resource docs with `guest_format` argument and added example with `swap` ([#1113](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1113))
* Added volume status check in `compute_volume_attach_v2` resource create function ([#1106](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1106))
* Added `disable_status_check` argument for `dns_zone_v2` resource ([#1114](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1114))
* Removed mention of `floating_ip` argument in `compute_instance_v2` from the documentation of `compute_floatingip_associate_v2` resource ([#1117](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1117))

BUG FIXES

* Fixed an issue when updating a `networking_router_v2` resource deleted extra routes on the router by upgrading Gophercloud to `1.14.0` ([#1109](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1109))
* Fixed an issue when changing `domain_id`, `is_domain` or `parent_id` arguments of `identity_project_v3` resource caused errors ([#1101](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1101))
* Fixed an issue when `fixed_ip` wasn't updated on read of `compute_interface_attach_v2` resource ([#1118](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1118))

## 1.33.0 (November 11, 2020)

IMPROVEMENTS

* Add `address_scope`, `security_group` and `subnetpool` RBAC types to `networking_rbac_policy_v2.go` resource ([#1086](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1086))
* Add `project_id` for `dns_zone_v2` resource, `project_id`, `all_projects` arguments for `dns_zone_v2` datasource, allow importing resource by `<zone_uuid>:<project_id>` value ([#1087](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1087))
* Add `different_cell` scheduler hint for `compute_instance_v2` resource ([#1070](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1070))
* Update and cleanup `blockstorage_quotaset_v2`, `blockstorage_quotaset_v3`, `compute_quotaset_v2`, `networking_quota_v2` resource docs ([#1095](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1095)), ([#1096](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1096))
* Updated `terraform-plugin-sdk` to `v1.16.0` ([#1092](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1092))

BUG FIXES

* Fixed an issue when `binding.host_id` was set to `null` in case of using any other `binding` parameters in `networking_port_v2` resource ([#1084](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1084))
* Fixed an issue with unnecessary server rebuild while using two default networks in `compute_instance_v2` resource ([#1073](https://github.com/terraform-provider-openstack/terraform-provider-openstack/pull/1073))

## 1.32.0 (September 15, 2020)

NOTES

* This is the first release that is available from [registry.terraform.io](https://registry.terraform.io)

IMPROVEMENTS

* Fixed documentation for `identity_ec2_credential_v3` resource ([#1052](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1052))
* Added `network_mode` argument for `compute_instance_v2` resource ([#1054](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1054))

## 1.31.0 (August 28, 2020)

FEATURES

* __New Resource__: `identity_ec2_credential_v3` ([#1033](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1033))

IMPROVEMENTS

* Reduced Identity requests across some `identity` resources and data sources by reusing functions to get the current token scope details ([#1044](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1044))
* Added `floating_ip_enabled` argument into `containerinfra_cluster_v1` datasource ([#1043](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1043))
* Updated Rackspace compatibility notes in documentation ([#1049](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1049))
* Updated `terraform-plugin-sdk` to `v1.15.0` ([#1051](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1051))
* Updated Go version to `1.14.7` ([#1051](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1051))

BUG FIXES

* Fixed backward compatibility issue with empty value in `merge_labels` argument of `containerinfra_cluster_v1` ([#1039](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1039))
* Fixed errors while creating `keymanager_container_v1` resource with the `certificate` type ([#1046](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1046))

## 1.30.0 (August 05, 2020)

FEATURES

* __New Resource__: `openstack_identity_group_v3` ([#1028](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1028))
* __New Data Source__: `openstack_images_image_ids_v2` ([#139](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/139))

IMPROVEMENTS

* Added `floating_ip_enabled` argument/attribute and `merge_labels` argument for `containerinfra_cluster_v1` resource ([#1024](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1024))
* Added `allowed_cidrs` argument/attribute for `lb_listener_v2` resource ([#1034](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1034))

## 1.29.0 (June 29, 2020)

FEATURES

* __New Data Source__: `compute_instance_v2` ([#984](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/984))

IMPROVEMENTS

* Added `vip_network_id` argument to `openstack_lb_loadbalancer_v2` resource. It can be used only with Octavia ([#948](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/948))
* Allowed to use `is_public` as argument in `compute_flavor_v2` datasource ([#1017](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1017))
* Updated `gophercloud` to `v0.12.0` to fix goroutine leaks during reauthentication ([#1020](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1020))
* Updated `terraform-plugin-sdk` to `v1.14.0` ([#1021](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1021))
* Updated Go version to `1.14.4` ([#1022](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1022))

BUG FIXES

* Fixed documentation bug for the `binding` argument of the `networking_port_v2` resource ([#1009](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1009))

## 1.28.0 (May 04, 2020)

NOTES

* This release sets `delayed_auth` and `allow_reauth` to `true` so Terraform provider won't request a new Identity token for every request against OpenStack API. We're happy to see you feedback about this change in our provider repo.

IMPROVEMENTS

* Provider parameters `delayed_auth` and `allow_reauth` are set to `true` by default ([#996](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/996))
* Added support to import `objectstorage_container_v1` resource. Some attributes can't be imported yet: `force_destroy`, `content_type`, `metadata`, `container_sync_to`, `container_sync_key` ([#998](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/998))
* Added `availability_zone_hints` parameter to `compute_instance_v2` resource ([#985](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/985))
* Added `SOURCE_IP_PORT` load balancing method for `lb_pool_v2` resource. It's only available in Octavia LoadBalancer service ([#993](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/993))
* Added `tags` for `identity_project_v3` resource and data source ([#978](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/978))
* Added `scheduler_hints` for `blockstorage_volume_v2`, `blockstorage_volume_v3` resources ([#983](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/983))
* Added `kubeconfig` attribute for `containerinfra_cluster_v1` resource ([#937](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/937))
* Updated the existing manifests in the `examples` directory and added new manifests with attaching and using a new volume ([#892](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/892))
* Updated Go version to `1.14.2` ([#1001](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1001))
* Updated `terraform-plugin-sdk` to `v1.11.0` ([#1001](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/1001))

BUG FIXES

* Fixed race conditions for `networking_secgroup_rule_v2` resource on some OpenStack environments ([#994](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/994))
* Fixed error logs for `keymanager_secret_v1` resource ([#997](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/997))

## 1.27.0 (April 13, 2020)

FEATURES

* __New Resource__: `openstack_keymanager_order_v1` ([#992](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/992))
* __New Resource__: `openstack_lb_members_v2` ([#898](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/898))

IMPROVEMENTS

* Added `detach_ports_before_destroy` argument for `compute_instance_v2` resource that allows to detach all instance ports prior trying to delete the instance ([#866](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/866))
* Added `web-download` import method to `openstack_images_image_v2` resource ([#888](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/888))
* Updated object URL in documentation for `versioning.type` attribute of `objectstorage_container_v1` resource ([#986](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/986))
* Added ACL examples in documentation of `objectstorage_container_v1` resource ([#987](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/987))

BUG FIXES

* Fixed `master_addresses`, `node_addresses` types to `schema.TypeList` since they are lists of strings instead of just strings ([#981](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/981))

## 1.26.0 (February 25, 2020)

IMPROVEMENTS

* Added `acl` argument and attribute to `openstack_keymanager_secret_v1`, `openstack_keymanager_container_v1` resources and datasources ([#956](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/956))
* Added `insert_headers` argument to `openstack_lb_listener_v2` resource ([#959](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/959))
* Added `block_device.volume_type` argument to `openstack_compute_instance_v2` resource ([#963](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/963))
* Updated `terraform-plugin-sdk` to `v1.7.0` ([#970](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/970))

BUG FIXES

* Fixed documentation bug for the `id` attribute of the `lb_policy_v2` resource ([#957](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/957))

NOTES

* This release drops Ubuntu Trusty and OpenStack Mitaka from testing CI environment.

## 1.25.0 (December 25, 2019)

FEATURES

* __New Resource__: `openstack_orchestration_stack_v1` ([#944](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/944))
* __New Data Source__: `openstack_blockstorage_volume_v2` ([#928](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/928))
* __New Data Source__: `openstack_blockstorage_volume_v3` ([#947](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/947))

IMPROVEMENTS

* Added `allow_reauth` optional boolean flag to the provided configuration block. This flag allows to automatically re-issue a new auth token if the initial token was expired ([#918](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/918))
* Added `fixed_network` and `fixed_subnet` arguments and attributes to `openstack_containerinfra_cluster_v1` resource and datasource ([#933](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/933))
* Added `access_rules` argument into `openstack_identity_application_credential_v3` resource ([#920](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/920))
* Support `SHELVE_OFFLOADED` status for `openstack_compute_instance_v2` resource ([#942](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/942))
* Added `max_retries_down` to `lb_monitor_v2` resource ([#945](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/945))
* Updated `terraform-plugin-sdk` to `v1.4.1` ([#936](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/936))

BUG FIXES

* Fixed the bug where empty `external_fixed_ips.ip_address` of `openstack_networking_router_v2` caused errors ([#628](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/628))
* Fixed documentation example for `openstack_identity_user_v3.extra` ([#923](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/923))
* Fixed documentation link for `clouds.yaml` ([#943](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/943))

## 1.24.0 (October 22, 2019)

FEATURES

* __New Resource__: `openstack_networking_quota_v2` ([#915](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/915))
* __New Resource__: `openstack_compute_quotaset_v2` ([#914](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/914))

IMPROVEMENTS

* Added `tags` argument/attribute and `all_tags` for `openstack_compute_instance_v2` resource ([#899](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/899))
* Added `UDP` protocol support for `openstack_lb_pool_v2`, `openstack_lb_monitor_v2`, `openstack_lb_listener_v2` resources. It is available only when `use_octavia` is set to `true` ([#896](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/896))
* Added ability to reuse the existing token when scope parameters are not defined ([#912](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/912))
* Migrated from Terraform in-tree `helper/*` SDK to the separate `terraform-plugin-sdk v1.1.1` ([#880](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/880)), ([#909](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/909))
* Migrated to use the common JSON debugging implementation from the upstream `gophercloud/utils` library ([#910](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/910))

BUG FIXES

* Fixed the bug with unchecked errors in initialization of Identity V3 client in `identity_auth_scope_v3` data source [[#878](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/878)] 
* Fixed the bug with unchecked errors in initialization of Compute V2 client in `compute_floatingip_associate_v2` resource [[#878](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/878)] 
* Fixed the bug with 404 errors handling while getting statuses tree in `openstack_lb_loadbalancer_v2` resource ([#883](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/883))
* Fixed the bug where is was unable to remove TLS references in `openstack_lb_listener_v2` resource ([#891](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/891))
* Fixed the bug where empty `scheduler_hints` list caused a panic in `openstack_compute_instance_v2` resource ([#885](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/885))
* Fixed the bug with usage of the wrong `flavor` argument instead of `flavor_id` for `openstack_lb_loadbalancer_v2` resource. Old argument has never worked. ([#904](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/904))
* Fixed the documentation bug with usage of `type` and `name` of the `persistence` of the `lb_pool_v2` resource ([#908](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/908))

## 1.23.0 (September 20, 2019)

FEATURES

* __New Resource__: `openstack_images_image_access_accept_v2` ([#872](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/872))
* __New Resource__: `openstack_images_image_access_v2` ([#872](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/872))

IMPROVEMENTS

* Added ability to reduce auth requests against the Identity service. This behaviour can be enabled via `delayed_auth` config flag ([#861](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/861))
* Added `Cache-Control: no-cache` header by default in all requests. This behaviour can be disabled via `disable_no_cache_header` config flag ([#849](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/849))
* Added `timeout_client_data`, `timeout_member_connect`, `timeout_member_data`, `timeout_tcp_inspect` arguments to the `openstack_lb_listener_v2` resource. Those arguments available if `use_octavia` is set to `true` ([#876](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/876)], [[#877](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/877))
* Added `domain_id`, `domain_name` attributes to the `openstack_identity_auth_scope_v3` data source ([#871](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/871))
* Added `description` attributes to the `openstack_identity_group_v3`, `openstack_identity_user_v3` data sources ([#874](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/874))
* Updated Terraform SDK to `v0.12.8` ([#859](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/859))
* Refactored headers formatting functions to not use external libraries and nested loops ([#865](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/865))

BUG FIXES

* Fixed the bug where `openstack_identity_auth_scope_v3` caused a panic within the domain-scope ([#851](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/851))
* Fixed the bug where `openstack_compute_flavor_access_v2` resource wasn't removed from the Terraform state when it has been deleted in the OpenStack ([#856](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/856))
* Fixed the bug where `openstack_identity_role_assignment_v3` resource wasn't removed from the Terraform state when it has been deleted in the OpenStack ([#856](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/856))
* Fixed the bug where `ephemeral` argument wasn't set for `openstack_compute_flavor_v2` while reading this resource from the API ([#855](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/855))

## 1.22.0 (September 05, 2019)

FEATURES

* __New Data Source__: `openstack_keymanager_container_v1` ([#846](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/846))

IMPROVEMENTS

* Added workaround for cases when the Neutron API doesn't provide the status for some load-balancer resources ([#839](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/839))
* Added workaround for cases when the OpenContrail API doesn't provide the ID for some load-balancer resources ([#840](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/840))
* Set computed attribute to `dns_name` and `dns_domain` for the `openstack_networking_network_v2` and  `openstack_networking_floatingip_v2` resources ([#837](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/837))
* Fixed code highlighting in website documentation for the `openstack_compute_instance_v2` resource ([#834](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/834))

BUG FIXES

* Fixed the bug where project info wasn't accessible to non-admin users ([#833](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/833))
* Fixed the bug where role assignments weren't accessible to non-admin users ([#845](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/845))

## 1.21.1 (August 08, 2019)

BUG FIXES

* Fixed the bug where OpenStack Networking V2 resources and data sources didn't work in old OpenStack environments because of different time format used for `created_at` and `updated_at` fields ([#831](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/831))

## 1.21.0 (August 06, 2019)

FEATURES

* __New Resource__: `openstack_keymanager_secret_v1` ([#650](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/650)), ([#807](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/807))
* __New Resource__: `openstack_keymanager_container_v1` ([#808](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/808))
* __New Resource__: `openstack_identity_service_v3` ([#821](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/821))
* __New Resource__: `openstack_identity_endpoint_v3` ([#823](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/823))
* __New Resource__: `openstack_networking_rbac_policy_v2` ([#811](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/811))
* __New Resource__: `openstack_blockstorage_quotaset_v2` ([#806](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/806))
* __New Resource__: `openstack_blockstorage_quotaset_v3` ([#828](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/828))
* __New Data Source__: `openstack_keymanager_secret_v1` ([#815](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/815))
* __New Data Source__: `openstack_identity_service_v3` ([#819](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/819))

IMPROVEMENTS

* Enabled the `openstack_compute_instance_v2` resource import ([#768](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/768))
* Added ability to update metadata of the `openstack_sharedfilesystem_share_v2` resource ([#825](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/825))
* Added ability to filter `openstack_identity_endpoint_v3` datasource by `service_type`, `endpoint_region` and `name` arguments ([#817](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/817))
* Updated the website documentation to formalize inline HCL code to canonical format according to Terraform v0.12 style conventions ([#797](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/797))
* Updated the website documentation to use `openstack_compute_volume_attach_v2` instead of `openstack_compute_volume_attach_v3` that doesn't exist ([#800](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/800))
* Updated the website documentation for the `security_groups` argument of the `openstack_compute_instance_v2` resource ([#826](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/826))

BUG FIXES

* Fixed the bug where `openstack_vpnaas_site_connection` resource set `admin_state_up` argument to `false` istead of `true` by default ([#799](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/799))
* Fixed the bug where `openstack_networking_subnet_v2` resource could cause a panic if `dns_nameservers` argument set to an empty list ([#726](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/726))
* Fixed the bug where `openstack_lb_pool_v2` resource could cause a panic because of passing a struct instead of a flattened list into the `persistence` attribute ([#725](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/725))
* Fixed the bug where `openstack_networking_port_v2` resource built an invalid request against the API with the empty `binding:profile` parameter ([#759](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/759))

## 1.20.0 (July 09, 2019)

FEATURES

* __New Resource__: `openstack_networking_qos_policy_v2` ([#774](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/774))
* __New Resource__: `openstack_networking_qos_bandwidth_limit_rule_v2` ([#783](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/783))
* __New Resource__: `openstack_networking_qos_dscp_marking_rule_v2` ([#784](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/784))
* __New Resource__: `openstack_networking_qos_minimum_bandwidth_rule_v2` ([#790](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/790))
* __New Data Source__: `openstack_networking_qos_policy_v2`([#779](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/779))
* __New Data Source__: `openstack_networking_qos_bandwidth_limit_rule_v2` ([#788](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/788))
* __New Data Source__: `openstack_networking_qos_dscp_marking_rule_v2` ([#789](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/789))
* __New Data Source__: `openstack_networking_qos_minimum_bandwidth_rule_v2` ([#793](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/793))

IMPROVEMENTS

* Updated documentation and Travis CI configuration with newer versions of Go and Terraform ([#777](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/777))
* Added `qos_policy_id` to `openstack_networking_network_v2` ([#780](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/780))
* Added `qos_policy_id` to `openstack_networking_port_v2` ([#781](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/781))
* Updated Terraform SDK to `v0.12.2` ([#795](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/795))

BUG FIXES

* Fixed bug preventing a floating IP from being re-associated with an instance when using `create_before_destroy` ([#761](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/761))
* Fixed bug preventing `openstack_compute_instance_v2` scheduler hint queries from working ([#771](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/771))

## 1.19.0 (May 22, 2019)

IMPROVEMENTS

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.
* `openstack_compute_instance_v2.stop_before_destroy` is now configurable by the `delete` timeout ([#750](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/750))

BUG FIXES

* Fixed bug where `openstack_dns_recordset_v2.ttl` was being cleared ([#752](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/752))
* Fixed an out of memory issue when running in debug mode ([#755](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/755))
* Fixed printing of clear text password in case of `v2` auth ([#757](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/757))

## 1.18.0 (May 08, 2019)

NOTES

* The `openstack_networking_subnet_v2` argument `allocation_pools` has been deprecated in favor of `allocation_pool`.

FEATURES

* __New Data Source__: `openstack_networking_addressscope_v2` ([#741](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/741))

BUG FIXES

* Fixed bug where `master_flavor` was being ignored in `openstack_containerinfra_cluster_v1` ([#730](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/730))
* Fixed case-sensitivity for validation on `access_type` and `access_level` in `openstack_sharedfilesystem_share_access_v2` ([#730](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/730))
* The `openstack_networking_subnet_v2` argument `allocation_pools` has been deprecated in favor of `allocation_pool`. This deprecation helps resolve an issue where multiple allocation pools in a single subnet were being returned out of order ([#739](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/739))
* Fixed a bug where `dns_nameservers` could not be cleared in `openstack_networking_subnet_v2` ([#728](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/728))
* Fixed a bug where a port's `dns_name` was being unset by Terraform ([#748](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/748))


## 1.17.0 (April 01, 2019)

NOTES

* `extra_dhcp_option` in the `openstack_networking_port_v2` data source has been changed to a List. This is to resolve a bug where multiple DHCP options were not being rendered.


FEATURES

* __New Resource__: `openstack_identity_application_credential_v3` ([#660](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/660))
* __New Data Source__: `openstack_blockstorage_availability_zones_v3` ([#652](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/652))
* __New Data Source__: `openstack_sharedfilesystem_availability_zones_v2` ([#652](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/652))
* __New Data Source__: `openstack_networking_trunk_v2` ([#626](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/626))

IMPROVEMENTS

* Reduced API calls when updating `extra_dhcp_option` in `openstack_networking_port_v2` ([#689](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/689))
* Added `port_security_enabled` to `openstack_networking_network_v2` ([#681](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/681))
* Added `port_security_enabled` to `openstack_networking_port_v2` ([#682](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/682))
* Added `prefix_length` to `openstack_networking_subnet_v2` ([#705](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/705))
* Added `binding` to `openstack_networking_port_v2` ([#693](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/693))
* Added `binding` to `openstack_networking_port_v2` data source ([#693](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/693))
* Added `mtu` to `openstack_networking_network_v2` ([#708](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/708))
* Added `mtu` to `openstack_networking_network_v2` data source ([#708](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/708))
* Added `dns_name` and `dns_domain` to `openstack_networking_floatingip_v2` ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_domain` to `openstack_networking_floatingip_v2` data source ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `dns_domain` to `openstack_networking_network_v2` ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `dns_domain` to `openstack_networking_network_v2` data source ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_assignment` to `openstack_networking_port_v2` ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `dns_name` and `dns_assignment` to `openstack_networking_port_v2` data source ([#706](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/706))
* Added `fixed_ip` to `openstack_networking_floatingip_associate_v2` ([#709](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/709))
* Enable `fixed_ip` to be updated in `openstack_networking_floatingip_v2` ([#709](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/709))
* Added ability to specify `cephx` as `access_type` and to retrieve the `access_key` in `openstack_sharedfilesystem_share_access_v2` ([#715](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/715))

BUG FIXES

* Fixed bug in `openstack_identity_auth_scope_v3` data source where the `user_id` attribute was being set to the user's Name and not ID ([#660](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/660))
* Fixed bug in Load Balancer resources for Contrail-based load balancers ([#691](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/691))
* `extra_dhcp_option` in the `openstack_networking_port_v2` data source has been changed to a List. This is to resolve a bug where multiple DHCP options were not being rendered ([#695](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/695))

## 1.16.0 (February 14, 2019)

NOTES

* The `openstack_networking_subnet_v2.host_routes` argument has been marked as deprecated. Please use the dedicated `openstack_networking_subnet_route_v2` resource instead.

FEATURES

* __New Data Source__: `openstack_compute_availability_zones_v2` ([#655](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/655))

BUG FIXES

* The `openstack_networking_subnet_v2.host_routes` argument has been deprecated due to schema issues and conflicts with `openstack_networking_subnet_route_v2` ([#668](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/668))
* A previously added validation to `openstack_networking_port_v2.fixed_ip.ip_address` was removed as it was causing problems for prior behavior of using an empty string ([#678](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/678))

## 1.15.1 (February 08, 2019)

BUG FIXES

* Fixed issue where volume multiattachments would not be retried ([#540](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/540))
* Reverted an incorrect schema validation for `openstack_networking_port_v2.allowed_address_pairs` ([#661](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/661))

## 1.15.0 (February 06, 2019)

NOTES

* The `openstack_images_image_v2.update_at` attribute has been deprecated in favor of `updated_at` ([#617](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/617))

FEATURES

* __New Resource__: `openstack_networking_addressscope_v2` ([#634](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/634))
* __New Resource__: `openstack_networking_port_secgroup_associate_v2` ([#574](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/574))

IMPROVEMENTS

* Added `flavor_id` to the `openstack_compute_flavor_v2` data source so flavors can be queried by ID ([#587](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/587))
* `openstack_networking_port_ids_v2` data source can now return an empty set of results ([#631](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/631))
* Added `description` to `openstack_networking_trunk_v2` resource ([#625](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/625))
* Added `tags` to the networking data source to query by tags and `all_tags` to see a full list of tags ([#624](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/624))
* `openstack_compute_instance_v2.admin_pass` is now a "sensitive" attribute ([#647](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/647))
* Added support to authenticate with Application Credentials ([#642](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/642))
* Added ability to specify region in `openstack_sharedfilesystem_share_access_v2` ([#654](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/654))
* Added ability to specify region in `openstack_sharedfilesystem_share_v2` ([#654](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/654))
* Added `all_tags` attribute to Networking resources to set tags provided by the OpenStack backend automatically ([#623](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/623))

BUG FIXES

* Fixed `created_at`, `updated_at`, and `tag` fields in the `openstack_images_image_v2` data source ([#615](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/615))
* Fixed `created_at` and `updated_at` fields in the `openstack_networking_subnetpool_v2` resource ([#619](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/619))
* Fixed `created_at` and `updated_at` fields in the `openstack_networking_subnetpool_v2` data source ([#616](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/616))
* Fixed issue where updating the description of a floating IP would cause the port to disassociate ([#606](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/606))
* `admin_state_up` and `shared` fields of `openstack_networking_network_v2` are now correct boolean fields ([#593](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/593))
* `external` field of `openstack_networking_network_v2` field will now show an actual value ([#593](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/593))
* Fixed issue where `status` was being used as the query value for `network_id` in `openstack_networking_port_v2` data source ([#631](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/631))
* Fixed issue where `status` was being used as the query value for `network_id` in `openstack_networking_port_ids_v2` data source ([#631](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/631))
* The `openstack_images_image_v2` fields `update_at`, `updated_at`, and `created_at` all now set correctly ([#617](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/617))
* Fixed issue with `openstack_dns_recordset_v2` where `records` would be returned out of order ([#636](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/636))
* Fixed issue where `openstack_compute_volume_attach_v2` and `openstack_blockstorage_volume_v2` were trying to detach volumes at the same time ([#640](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/640))
* Fixed a regression bug where destroying networks was failing on a 409 code ([#644](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/644))
* Fixed an issue with `openstack_compute_instance_v2` where a 404 was triggering an error ([#647](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/647))
* Fixed an issue where `all_fixed_ips` was not being set in `openstack_networking_port_v2` data source ([#649](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/649))
* Fixed an issue where `openstack_networking_port_v2` would cause an API error ([#649](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/649))
* Fixed an issue where Blockstorage volume resources couldn't be detached because they had been removed ([#641](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/641))

## 1.14.0 (January 15, 2019)

NOTES

* The Load Balancer v2 resources have been updated to provide more efficient status checks. If you encounter any problems due to this, please report them and we will make it a priority to resolve.
* `openstack_networking_port_v2` will now set the `admin_state_up` to `true/UP` if it is left omitted from the resource configuration. This now correctly conforms to the OpenStack API. This should be a transparent change, but let us know if this causes you problems.

FEATURES

* __New Resource__: `openstack_lb_l7policy_v2` ([#527](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/527))
* __New Resource__: `openstack_lb_l7rule_v2` ([#522](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/522))
* __New Resource__: `openstack_sharedfilesystem_share_v2` ([#525](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/525))
* __New Resource__: `openstack_sharedfilesystem_share_access_v2` ([#526](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/526))
* __New Data Source__: `openstack_sharedfilesystem_share_v2` ([#564](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/564))
* __New Data Source__: `openstack_networking_port_v2` ([#567](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/567))
* __New Data Source__: `openstack_sharedfilesystem_sharenetwork_v2` ([#576](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/576))
* __New Data Source__: `openstack_networking_port_ids_v2` ([#569](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/569))
* __New Data Source__: `openstack_sharedfilesystem_snapshot_v2` ([#577](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/577))

IMPROVEMENTS

* Provider options `swauth` and `use_octavia` will correctly use a default value of `false` when they are not specified. This is to help with compatibility for v0.12 ([#494](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/494))
* Enhanced the pending status checks of the Load Balancer v2 resources ([#550](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/550))
* Prioritized the status of Load Balancer v2 resources to first use the Load Balancer's master status ([#556](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/556))
* Fix flavor detection in `openstack_compute_instance_v2` and `openstack_containerinfra_cluster_v1` for Terraform v0.12 ([#551](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/551))
* Added the ability to import `openstack_lb_loadbalancer_v2` ([#524](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_listener_v2` ([#524](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_pool_v2` ([#524](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_member_v2` ([#524](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/524))
* Added the ability to import `openstack_lb_monitor_v2` ([#524](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/524))
* Added `device_type` and `disk_bus` to `openstack_compute_instance_v2` block device ([#558](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/558))
* Added `transparent_vlan` to `openstack_networking_network_v2` ([#513](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/513))
* Added `transparent_vlan` to `openstack_networking_network_v2` data source ([#538](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/538))
* Added `max_retries` to the provider options ([#413](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/413))
* Added the ability to override catalog endpoints ([#501](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/501))
* Changed the `segments` attribute of the `openstack_networking_network_v2` to `TypeSet` [[#578](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/578)] 

BUG FIXES

* `openstack_compute_interface_attach_v2` now correctly sets the `instance_id` [[#557](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/557)] 
* `openstack_networking_port_v2` will now correctly set the `admin_state_up` to `true/UP` if left omitted ([#594](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/594))
* Fixed out of range panic in `openstack_compute_instance_v2` when no IP addresses were detected ([#539](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/539))

## 1.13.0 (December 18, 2018)

FEATURES

* __New Resource__: `openstack_sharedfilesystem_securityservice_v2` ([#515](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/515))
* __New Resource__: `openstack_sharedfilesystem_sharenetwork_v2` ([#515](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/515))
* __New Data Source__: `openstack_containerinfra_cluster_v1` ([#488](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/488))
* __New Data Source__: `openstack_blockstorage_snapshot_v2` ([#448](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/448))
* __New Data Source__: `openstack_blockstorage_snapshot_v3` ([#448](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/448))

IMPROVEMENTS

* Added object versioning to `openstack_objectstorage_container_v1` ([#465](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/465))
* Added support for soft affinities in `openstack_compute_servergroup_v2` ([#490](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/490))
* Allow `default_pool_id` to be updated in `openstack_lb_listener_v2` ([#516](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/516))
* Added `description` to `openstack_networking_router_v2` ([#529](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/529))
* Added `description` to `openstack_networking_port_v2` ([#531](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/531))
* Added `description` to `openstack_networking_subnet_v2` ([#533](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/533))
* Added `description` to `openstack_networking_floatingip_v2` ([#534](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/534))
* Added `description` to `openstack_networking_secgroup_v2` data source ([#535](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/535))
* Added `description` to `openstack_networking_network_v2` ([#532](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/532))
* Added `description` to `openstack_networking_subnet_v2` data source ([#528](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/528))
* Added `description` to `openstack_networking_router_v2` data source ([#530](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/530))
* Added `description` to `openstack_networking_network_v2` data source ([#536](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/536))
* Added `description` to `openstack_networking_floatingip_v2` data source ([#523](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/523))

BUG FIXES

* Allow instances to be in a state of `migrating` when performing a plan/refresh ([#496](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/496))
* Fix issue when `openstack_networking_floatingip_v2`, `openstack_networking_router_v2`, `openstack_networking_subnet_v2`, and `openstack_networking_subnetpool_v2` tag updates send empty updates for the resource. ([#519](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/519))

## 1.12.0 (November 13, 2018)

FEATURES

* __New Resource__: `openstack_compute_interface_attach_v2` ([#470](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/470))

IMPROVEMENTS

* Added `tags` to `openstack_networking_network_v2` ([#454](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/454))
* Added `tags` to `openstack_networking_subnet_v2` ([#459](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/459))
* Added `tags` to `openstack_networking_subnetpool_v2` ([#460](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/460))
* Added `tags` to `openstack_networking_port_v2` ([#461](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/461))
* Added `tags` to `openstack_networking_secgroup_v2` ([#463](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/463))
* Added `tags` to `openstack_networking_floatingip_v2` ([#466](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/466))
* Added `tags` to `openstack_networking_router_v2` ([#467](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/467))
* Added `extra_dhcp_options` to `openstack_networking_port_v2` ([#258](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/258))
* Added `fingerprint` to `openstack_compute_keypair_v2` data source ([#481](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/481))
* Added `extra_specs` to `openstack_compute_flavor_v2` data source ([#480](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/480))

BUG FIXES

* Fixed issue with nova-network based environments having the `tenantnetworks` API disabled ([#485](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/485))


## 1.11.0 (October 29, 2018)

FEATURES

* __New Resource__: `openstack_networking_trunk_v2` ([#446](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/446))
* __New Resource__: `openstack_compute_flavor_access_v2` ([#447](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/447))

IMPROVEMENTS

* Added `multiattach` argument and attribute for the `openstack_blockstorage_volume_v3` resource ([#431](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/431))
* `openstack_dns_recordset_v2` can now accept IPv6 addresses with and without brackets ([#443](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/443))
* Added `multiattach` argument for the `openstack_compute_volume_attach_v2` resource ([#442](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/442))
* `openstack_lb_member_v2` resources can now use a weight of 0 ([#451](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/451))

BUG FIXES

* Fixed an issue where environment variables were overwriting specified arguments ([#436](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/436))
* Fixed an issue where security group rule descriptions were not working with older verisons of OpenStack ([#438](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/438))

## 1.10.0 (October 01, 2018)

FEATURES

* __New Resource__: `openstack_containerinfra_cluster_v1` ([#421](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/421))
* __New Data Source__: `openstack_containerinfra_clustertemplate_v1` ([#415](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/415))

IMPROVEMENTS

* Added `description` argument for the `openstack_networking_secgroup_rule_v2` resource ([#416](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/416))
* Added a vendor option of `ignore_resize_confirmation` to `openstack_compute_instance_v2` ([#422](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/422))
* `openstack_compute_instance_v2` IP addresses are now visible in Rackspace. This provider still does not officially support Rackspace, though. ([#426](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/426))
* Added `no_fixed_ip` argument to `openstack_networking_port_v2` which allows the port to not have an IP address ([#433](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/433))

BUG FIXES

* Enabled instances to be in an `ERROR` state so they can be cleanly deleted ([#428](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/428))

## 1.9.0 (September 05, 2018)

FEATURES

* __New Resource__: `openstack_objectstorage_tempurl_v1` ([#379](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/379))
* __New Resource__: `openstack_containerinfra_clustertemplate_v1` ([#403](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/403))
* __New Data Source__: `openstack_fw_policy_v1` ([#398](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/398))
* __New Data Source__: `openstack_networking_router_v2` ([#401](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/401))

IMPROVEMENTS

* The `openstack_images_image_v2` resource can now finally update properties. This update has been in progress over the last two release cycles. Please let us know if you encounter any problems ([#409](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/409))

## 1.8.0 (August 08, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Support for `default_domain` has been added. This should not cause any issues, but please report any issues encountered.
* `openstack_images_image_v2.properties` has been set to `ForceNew`. If properties are modified, the image will be recreated. Previously, updates to the properties were only happening in the Terraform state and not actually reflected on the image itself.

FEATURES

* __New Data Source__: `openstack_identity_group_v3` ([#385](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/385))
* __New Data Source__: `openstack_networking_floatingip_v2` ([#387](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/387))

IMPROVEMENTS

* Added support for `default_domain` during authentication ([#329](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/329))
* The upcoming OpenStack Rocky release will be automatically adding additional properties to the `openstack_images_image_v2` resource. This resource has been patched to account for this and to reconcile these server-provided properties with the user-provided properties. In addition, `openstack_images_image_v2.properties` has been set to `ForceNew` and will recreate the image when properties have been modified. Previously, any updates to the properties were only happening in the state and not actually reflected on the image itself. ([#390](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/390))

BUG FIXES

* The addition of the `openstack_networking_network_v2.external` data source argument caused unintended behavior of results only containing external or non-external networks. This bug has been fixed and we apologize for the inconvenience ([#384](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/384))
* The addition of the `openstack_compute_floatingip_associate_v2.wait_until_associated` argument caused the floating IP association to be recreated when updating to a later release of this provider. This was unintended and this has been resolved ([#395](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/395))

## 1.7.0 (August 01, 2018)

FEATURES

* __New Data Source__: `openstack_identity_endpoint_v3` ([#377](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/377))

IMPROVEMENTS

* Allow resize for stopped instances ([#348](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/348))
* Added `power_state` to `openstack_compute_instance_v2` ([#350](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/350))
* Added `external` to `openstack_networking_network_v2` resource ([#357](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/357))
* Added `external` to `openstack_networking_network_v2` data source ([#358](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/358))
* Return the default network uuid for `openstack_compute_instance_v2` ([#365](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/365))
* Allow a specific floating IP to be specified in `openstack_networking_floatingip_v2` ([#371](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/371))
* Allow `PROXY` protocol for `openstack_lb_pool_v2` ([#375](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/375))

BUG FIXES

* Allow explicit values of `0` for `min_disk_gb` and `min_ram_mb` in the `openstack_images_image_v2` resource ([#351](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/351))
* Make `peer_ep_group_id` optional in `openstack_vpnaas_site_connection` ([#353](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/353))

## 1.6.0 (June 20, 2018)

FEATURES

* __New Resource__: `openstack_vpnaas_site_connection_v2` ([#330](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/330))

IMPROVEMENTS

* Added `wait_until_associated` to `openstack_compute_floatingip_associate_v2` ([#310](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/310))
* Added support for SSL settings in a `clouds.yaml` file ([#340](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/340))

## 1.5.0 (May 15, 2018)

FEATURES

* __New Resource__: `openstack_blockstorage_volume_v3` ([#324](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/324))
* __New Resource__: `openstack_blockstorage_volume_attach_v3` ([#324](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/324))
* __New Resource__: `openstack_networking_subnet_route_v2` ([#314](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/314))
* __New Resource__: `openstack_networking_floatingip_associate_v2` ([#313](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/313))
* __New Resource__: `openstack_vpnaas_ipsec_policy_v2` ([#270](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/270))
* __New Resource__: `openstack_vpnaas_service_v2` ([#300](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/300))
* __New Resource__: `openstack_vpnaas_ike_policy_v2` ([#316](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/316))
* __New Resource__: `openstack_vpnaas_endpoint_group_v2` ([#321](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/321))
* __New Data Source__: `openstack_compute_keypair_v2` ([#307](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/307))
* __New Data Source__: `openstack_identity_auth_scope_v3` ([#204](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/204))

IMPROVEMENTS

* Added `verify_checksum` to `openstack_images_image_v2` resource so that checksum verification can be disabled ([#305](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/305))
* The LBaaS v2 resources have lower "delay" times when waiting for state changes. This should speed up creation of a Load Balancing stack ([#297](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/297))

BUG FIXES

* Fixed issue where `OS_IDENTITY_API_VERSION=2` was not recognized ([#315](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/315))
* Fixed issue when using Identity v3 resources when an Identity v2 endpoint is published ([#320](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/320))
* `openstack_networking_router_v2.distributed` will now pass `false` correctly ([#308](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/308))
* `openstack_networking_router_v2.enable_snat` will now pass `false` correctly ([#309](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/309))

## 1.4.0 (May 01, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The OpenStack provider now has additional authentication options for `project_domain_name`, `project_domain_id`, `user_domain_name`, and `user_domain_id`. This will allow for more fine-grainted authentication scoping. This should not cause any problems with existing deployments, but please report any authentication issues after upgrading.

FEATURES

* __New Resource__: `openstack_identity_role_assignment_v3` ([#265](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/265))
* __New Data Source__: `openstack_identity_project_v3` ([#251](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/251))
* __New Data Source__: `openstack_identity_user_v3` ([#252](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/252))

IMPROVEMENTS

* Added `member_status` to `openstack_images_image_v2` data source ([#269](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/269))
* Add support for `OS_TOKEN` environment variable ([#272](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/272))
* Added `force_destroy` to `openstack_objectstorage_container_v1` which will cause all objects in the container to be deleted when the container is deleted ([#276](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/276))
* CIDR is now optional in `openstack_networking_subnet_v2` allowing a CIDR to be allocated from a subnet pool ([#294](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/294))
* Added additional authentication options for domain scoping ([#290](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/290))
* `openstack_images_image_v2` can now support OVA format ([#302](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/302))

BUG FIXES

* `openstack_compute_instance_v2` resources can handle Availability Zones in the format of `az:host:node` ([#291](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/291))

## 1.3.0 (March 14, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* `openstack_compute_keypair_v2` can now generate a private key, however the private key will be stored in your Terraform state. Please use caution.
* The MAC addresses in `openstack_networking_port_v2.allowed_address_pairs` is no longer computed. This should not cause an issue for users since if an `allowed_address_pairs` MAC address was not specified, the AAP MAC will match `openstack_networking_port_v2.mac_address`.

FEATURES

* __New Resource:__ `openstack_networking_subnetpool_v2` ([#243](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/243))
* __New Resource:__ `openstack_identity_role_v3` ([#250](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/250))
* __New Data Source:__ `openstack_networking_subnetpool_v2` ([#243](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/243))
* __New Data Source:__ `openstack_identity_role_v3` ([#250](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/250))

IMPROVEMENTS

* Added `additional_properties` to `openstack_compute_instance_v2` scheduler hints ([#230](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/230))
* `openstack_compute_keypair_v2` can now generate a private key ([#217](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/217))
* `openstack_networking_router_v2` can now optionally set a default gateway after it has been created ([#209](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/209))
* Added `subnetpool_id` to `openstack_networking_subnet_v2` resource and data source ([#249](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/249))
* Added `extra_specs` to `openstack_compute_flavor_v2` ([#241](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/241))
* Added `subnet_id` to `openstack_networking_floatingip_v2` ([#240](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/240))

BUG FIXES

* Fixed bug with `openstack_networking_network_v2` and `openstack_networking_subnet_v2` where the `OS_TENANT_ID` was incorrectly being used as a default value ([#254](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/254))
* Correctly detect if an object storage container is deleted ([#261](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/261))
* Fixed a few small bugs with `openstack_fw_rule_v1` updating ([#224](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/224))
* Fixed an issue with `openstack_networking_port_v2` `allowed_address_pairs` and MAC addresses ([#244](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/244))

## 1.2.0 (January 18, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The way IP addresses for `allowed_address_pairs` in the `openstack_networking_port_v2` resource are stored in the Terraform state has changed. 
* The `external_gateway` argument in the `openstack_networking_router_v2` has been deprecated in favor of the more appropriately named `external_network_id`.

FEATURES

* __New Resource:__ `openstack_db_database_v1` ([#179](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/179))
* __New Resource:__ `openstack_db_user_v1` ([#180](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/180))
* __New Resource:__ `openstack_db_configuration_v1` ([#185](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/185))
* __New Data Source:__ `openstack_compute_flavor_v2` ([#190](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/190))


IMPROVEMENTS

* Added `external_fixed_ips` to the `openstack_networking_router_v2` resource ([#178](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/178))
* Added `ipv6_address_mode` and `ipv6_ra_mode` to the `openstack_networking_subnet_v2` resource and data source ([#193](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/193))
* Several new `openstack_networking_subnet_v2` attributes are now accessible in the data source ([#199](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/199))
* Added `availability_zone_hints` to the `openstack_networking_network_v2` resource and data source ([#196](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/196))
* Added `availability_zone_hints` to the `openstack_networking_router_v2` resource ([#203](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/203))
* User's password field in `openstack_db_instance_v2` resource has been marked sensitive ([#220](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/220))
* `openstack_db_instance_v1` now supports setting a `configuration_id` ([#221](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/221))

BUG FIXES

* Allow the same `ip_address` with a different `mac_address` to be specified multiple times in the `openstack_networking_port_v2` resource ([#168](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/168))
* Fixed unhandled error checks which were causing crashes in `openstack_networking_secgroup_v2` and `openstack_networking_network_v2` data sources ([#201](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/201))
* Fixed unhandled error check when creating `openstack_networking_floatingip_v2` ([#206](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/206))
* Fixed region detection when using `clouds.yaml` ([#216](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/216))
* Make `subnet_id` optional for `openstack_lb_member_v2` ([#189](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/189))
* Fix ordering of DNS servers in `openstack_networking_subnet_v2` ([#226](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/226))

## 1.1.0 (December 04, 2017)

FEATURES

* __New Resource:__ `openstack_objectstorage_object_v1` ([#146](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/146))
* __New Resource:__ `openstack_db_instance_v1` ([#155](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/155))

IMPROVEMENTS

* Better handling of mutually exclusive options `no_gateway` and `gateway_ip` in the `openstack_networking_subnet_v2` resource ([#136](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/136))
* Can now authenticate with a `clouds.yaml` file ([#154](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/154))

BUG FIXES

* Fixed issue with automatic detection of an Octavia client and Networking client ([#172](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/172))
* Fixed issue with creating public flavors ([#177](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/177))

## 1.0.0 (November 08, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* If your OpenStack cloud supports the Octavia Load Balancing service, you can now use it by setting the provider-level `use_octavia` argument to `true`. The `openstack_lb_*_v2` resources will then seamlessly use Octavia.

FEATURES

* __New Data Source:__ `openstack_networking_subnet_v2` ([#135](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/135))
* __New Data Source:__ `openstack_dns_zone_v2` ([#145](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/145))

IMPROVEMENTS

* `openstack_networking_router_v2`: Added `enable_snat` argument ([#140](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/140))
* Added provider-level option of `use_octavia` to use the Octavia load balancing service ([#149](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/149))

## 0.3.0 (October 23, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* The `openstack_networking_port_v2` resource had a significant update to how it handles security groups. If you have not explicitly defined security groups in the port resource, any security groups which were automatically applied by OpenStack (such as the `default` security group) will be removed upon the next apply. To prevent this from happening, add the ID of the security groups to the `security_group_ids` argument. If you are already explicitly specifying security groups, you should see no change in behavior.

IMPROVEMENTS

 * `openstack_networking_router_interface_v2` will now set `subnet_id` when importing ([#119](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/119))
 * `openstack_networking_router_route_v2` can now be imported ([#120](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/120))
 * `openstack_images_image_v2` resource and data source now supports reading and setting properties ([#113](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/113))

BUG FIXES

  * `openstack_networking_port_v2`: Fixed issues with how security groups and allowed address pairs are applied and updated [[#114](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/114)].

## 0.2.2 (September 15, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Unused `id` fields in the LBaaS v2 resources were removed. This should not cause any issues, but please report if you find otherwise.

FEATURES:

* __New Data Source:__ `openstack_networking_secgroup_v2` ([#86](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/86))
* __New Resource:__: `openstack_compute_flavor_v2` ([#83](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/83))

IMPROVEMENTS
 * Added `status` field to `openstack_networking_network_v2` data source ([#105](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/105))
 * `openstack_networking_router_v2` can now be imported ([#111](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/111))
 * `openstack_networking_router_interface_v2` can now be imported ([#112](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/112))
 
BUG FIXES

* `openstack_lb_listener_v2`: Don't send `connection_limit` unless it has been set ([#90](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/90))
* `openstack_lb_pool_v2`: Find Load Balancer via Listener ([#97](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/97))
* LBaaS v2: Removed unused `id` fields ([#93](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/93))
* `openstack_lb_monitor_v2`: Check if a monitor was successfully created before proceeding ([#102](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/102))
* `openstack_networking_router_v2`: Fix region parameter ([#107](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/107))
* `openstack_compute_instance_v2`: Fix regression bug with NIC detection ([#117](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/117))

## 0.2.1 (August 23, 2017)

IMPROVEMENTS:

* `openstack_lb_loadbalancer_v2` timeouts have been lowered to 10 and 5 minutes ([#74](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/74))

BUG FIXES:

* `openstack_images_image_v2` data source now sorts images by `CreatedAt` instead of `UpdatedAt` ([#78](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/78))
* `openstack_networking_secgroup_v2` now re-reads security group before deleteing rules when `delete_default_rules => true` ([#82](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/82))
* Fixed `openstack_compute_instance_v2` access IP address detection in dual-stack environments ([#85](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/85))

## 0.2.0 (August 14, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* Network detection in the `openstack_compute_instance_v2` resource was cleaned up and updated. There should be no incompatibilities, but you should do a `plan` before `apply` just to be safe.
* The `openstack_lb_loadbalancer_v2.provider` argument has been removed entirely. This was an erroneous argument from the beginning, so it should not be in use. However, if you do have it set in your configurations, please rename it to `loadbalancer_provider`.

FEATURES:

* __New Resource:__ `openstack_identity_project_v3` ([#50](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/50))
* __New Resource:__ `openstack_identity_user_v3` ([#52](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/52))

IMPROVEMENTS:

* `openstack_compute_instance_v2` now supports Neutron for network detection ([#39](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/39))
* `openstack_compute_instance_v2` support for multiple NICs on the same network ([#39](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/39))
* Added support for `TERMINATED_HTTPS` protocol in `openstack_lb_listener_v2` ([#49](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/49))
* Improvements to LBaaS v2 resource coordination ([#59](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/59))
* `openstack_lb_loadbalancer_v2.provider` has been removed. See notes above. ([#65](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/65))

BUG FIXES:
* `openstack_lb_pool_v2` handling of `persistence` updated, `cookie_name` is now optional. ([#57](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/57))
* `openstack_fw_firewall_v1.associated_routers` is now computed. ([#53](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/53))
* All `openstack_fw_rule_v1` attributes are now passed during an update phase. ([#53](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/53))
* `openstack_networking_secgroup_v2` now correctly updates description. ([#60](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/60))
* `openstack_fw_firewall_v1` now correctly translates `value_specs` on create. ([#66](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/66))

## 0.1.0 (June 21, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

* You can now specify `region` in the provider block. All resources will inherit this region setting, or you can override it in the resource-level `region`. Make sure to do a `plan` before an `apply` to make sure the resource is not destroyed due to incorrectly determining the region! If you see this happening, either explicitly set the `region` in the resource or use `lifecycle.ignore_changes`. 
* `floating_ip` has been removed from `openstack_compute_instance_v2`. You must now use `openstack_compute_floatingip_associate_v2` to associate a Floating IP with an Instance.
* `volume` has been removed from `openstack_compute_instance_v2`. You must now use `openstack_compute_volume_attach_v2` to attach a Volume with an Instance.
* `member` has been removed from `openstack_lb_pool_v1`. You must now use `openstack_lb_member_v1` to add a LBaaS v1 Member to a Pool.


IMPROVEMENTS:

* Can specify `region` in the provider ([#25](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/25))

BUG FIXES

* Wait for LoadBalancer to be active before creating Pools and Monitors ([#29](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/29))
* Choose first network found with a matching name for compute instances ([#36](https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/36))
