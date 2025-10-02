## 2.48.0 (February 18, 2021)

FEATURES:

* **New Data Source:** `azurerm_application_gateway` ([#10268](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10268))

ENHANCEMENTS:

* dependencies: updating to build using Go 1.16 which adds support for `darwin/arm64` (Apple Silicon) ([#10615](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10615))
* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v51.2.0` ([#10561](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10561))
* Data Source: `azurerm_bastion_host` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_point_to_site_vpn_gateway` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* Data Source: `azurerm_kubernetes_cluster_node_pool` - exposing the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* Data Source: `azurerm_route` - pdating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_subnet ` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_subscriptions` - adding the field `id` to the `subscriptions` block ([#10598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10598))
* Data Source: `azurerm_virtual_network` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_bastion_host` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_bastion_host` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_kubernetes_cluster` - support for configuring the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* `azurerm_kubernetes_cluster` - support for `automatic_channel_upgrade` ([#10530](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10530))
* `azurerm_kubernetes_cluster` - support for `skip_nodes_with_local_storage` within the `auto_scaler_profile` block ([#10531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10531))
* `azurerm_kubernetes_cluster` - support for `skip_nodes_with_system_pods` within the `auto_scaler_profile` block ([#10531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10531))
* `azurerm_kubernetes_cluster_node_pool` - support for configuring the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* `azurerm_lighthouse_definition` - add support for `principal_id_display_name` property ([#10613](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10613))
* `azurerm_log_analytics_workspace` - Support for `capacity_reservation_level` property and `CapacityReservation` SKU ([#10612](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10612))
* `azurerm_point_to_site_vpn_gateway` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_point_to_site_vpn_gateway` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_route` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_route` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_subnet` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_subnet` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `synapse_workspace_resource` - support for the `azure_devops_repo` and `github_repo` blocks ([#10157](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10157))
* `azurerm_virtual_network` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_virtual_network` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))

BUG FIXES:

* `azurerm_eventgrid_event_subscription` - change the number of possible `advanced_filter` items from `5` to `25` ([#10625](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10625))
* `azurerm_key_vault` - normalizing the casing on the `certificate_permissions`, `key_permissions`, `secret_permissions` and `storage_permissions` fields within the `access_policy` block ([#10593](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10593))
* `azurerm_key_vault_access_policy ` - normalizing the casing on the `certificate_permissions`, `key_permissions`, `secret_permissions` and `storage_permissions` fields ([#10593](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10593))
* `azurerm_mariadb_firewall_rule` - correctly validate the `name` property ([#10579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10579))
* `azurerm_postgresql_server` - correctly change `ssl_minimal_tls_version_enforced` on update ([#10606](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10606))
* `azurerm_private_endpoint` - only updating the associated Private DNS Zone Group when there's changes ([#10559](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10559))
* `azurerm_resource_group_template_deployment` - fixing an issue where the API version for nested items couldn't be found during deletion ([#10565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10565))

## 2.47.0 (February 11, 2021)

UPGRADE NOTES

* `azurerm_frontdoor` & `azurerm_frontdoor_custom_https_configuration` - the new fields `backend_pool_health_probes`, `backend_pool_load_balancing_settings`, `backend_pools`, `frontend_endpoints`, `routing_rules` have been added to the `azurerm_frontdoor` resource, which are a map of name-ID references. An upcoming version of the Azure Provider will change the blocks `backend_pool`, `backend_pool_health_probe`, `backend_pool_load_balancing`, `frontend_endpoint` and `routing_rule` from a List to a Set to work around an ordering issue within the Azure API - as such you should update your Terraform Configuration to reference these new Maps, rather than the Lists directly, due to the upcoming breaking change. For example, changing `azurerm_frontdoor.example.frontend_endpoint[1].id` to `azurerm_frontdoor.example.frontend_endpoints["exampleFrontendEndpoint2"]` ([#9357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9357))
* `azurerm_lb_backend_address_pool` - the field `backend_addresses` has been deprecated and is no longer functional - instead the `azurerm_lb_backend_address_pool_address` resource offers the same functionality. ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))
* `azurerm_linux_virtual_machine_scale_set` & `azurerm_windows_virtual_machine_scale_set` - the in-line `extension` block is now GA - the environment variable `ARM_PROVIDER_VMSS_EXTENSIONS_BETA` no longer has any effect and can be removed ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_data_factory_integration_runtime_managed` - this resource has been renamed/deprecated in favour of `azurerm_data_factory_integration_runtime_azure_ssis` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* The provider-block field `skip_credentials_validation` is now deprecated since this was non-functional and will be removed in 3.0 of the Azure Provider ([#10464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10464))

FEATURES:

* **New Data Source:** `azurerm_key_vault_certificate_data` ([#8184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8184))
* **New Resource:** `azurerm_application_insights_smart_detection_rule` ([#10539](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10539))
* **New Resource:** `azurerm_data_factory_integration_runtime_azure` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* **New Resource:** `azurerm_data_factory_integration_runtime_azure_ssis` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* **New Resource:** `azurerm_lb_backend_address_pool_address` ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))

ENHANCEMENTS:

* dependencies: updating `github.com/hashicorp/terraform-plugin-sdk` to `v1.16.0` ([#10521](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10521))
* `azurerm_frontdoor` - added the new fields `backend_pool_health_probes`, `backend_pool_load_balancing_settings`, `backend_pools`, `frontend_endpoints`, `routing_rules` which are a map of name-ID references ([#9357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9357))
* `azurerm_kubernetes_cluster` - updating the validation for the `log_analytics_workspace_id` field within the `oms_agent` block within the `addon_profile` block ([#10520](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10520))
* `azurerm_kubernetes_cluster` - support for configuring `only_critical_addons_enabled` ([#10307](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10307))
* `azurerm_kubernetes_cluster` - support for configuring `private_dns_zone_id` ([#10201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10201))
* `azurerm_linux_virtual_machine_scale_set` - the `extension` block is now GA and available without enabling the beta ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_media_streaming_endpoint` - exporting the field `host_name` ([#10527](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10527))
* `azurerm_mssql_virtual_machine` - support for `auto_backup` ([#10460](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10460))
* `azurerm_windows_virtual_machine_scale_set` - the `extension` block is now GA and available without enabling the beta ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_site_recovery_replicated_vm` - support for the `recovery_public_ip_address_id` property and changing `target_static_ip` or `target_static_ip` force a new resource to be created ([#10446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10446))

BUG FIXES:

* provider: the provider-block field `skip_credentials_validation` is now deprecated since this was non-functional. This will be removed in 3.0 of the Azure Provider ([#10464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10464))
* Data Source: `azurerm_shared_image_versions` - retrieving all versions of the image prior to filtering ([#10519](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10519))
* `azurerm_app_service` - the `ip_restriction.x.ip_address` propertynow accepts anything other than an empty string ([#10440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10440))
* `azurerm_cosmosdb_account` - validate the `key_vault_key_id` property is versionless ([#10420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10420))
* `azurerm_cosmosdb_account` - will no longer panic if the response is nil ([#10525](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10525))
* `azurerm_eventhub_namespace` - correctly downgrade to the `Basic` sku ([#10536](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10536))
* `azurerm_key_vault_key` - export the `versionless_id` attribute ([#10420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10420))
* `azurerm_lb_backend_address_pool` - the `backend_addresses` block is now deprecated and non-functional - use the `azurerm_lb_backend_address_pool_address` resource instead ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))
* `azurerm_linux_virtual_machine_scale_set` - fixing a bug when `protected_settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_linux_virtual_machine_scale_set` - fixing a bug when `settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_monitor_diagnostic_setting` - changing the `log_analytics_workspace_id` property no longer creates a new resource ([#10512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10512))
* `azurerm_storage_data_lake_gen2_filesystem` - do not set/retrieve ACLs when HNS is not enabled  ([#10470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10470))
* `azurerm_windows_virtual_machine_scale_set ` - fixing a bug when `protected_settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_windows_virtual_machine_scale_set ` - fixing a bug when `settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))

## 2.46.1 (February 05, 2021)

BUG FIXES:

* `azurerm_lb_backend_address_pool` - mark `backend_address` as computed ([#10481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10481))

## 2.46.0 (February 04, 2021)

FEATURES:

* **New Resource:** `azurerm_api_management_identity_provider_aadb2c` ([#10240](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10240))
* **New Resource:** `azurerm_cosmosdb_cassandra_table` ([#10328](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10328))

ENHANCEMENTS:

* dependencies: updating `recoveryservices` to API version `2018-07-10` ([#10373](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10373))
* `azurerm_api_management_diagnostic` - support for the `always_log_errors`, `http_correlation_protocol`, `log_client_ip`, `sampling_percentage` and `verbosity` properties ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_api_management_diagnostic` - support for the `frontend_request`, `frontend_response`, `backend_request` and `backend_response` blocks ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_kubernetes_cluster` - support for configuring the field `enable_host_encryption` within the `default_node_pool` block ([#10398](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10398))
* `azurerm_kubernetes_cluster` - added length validation to the `admin_password` field within the `windows_profile` block ([#10452](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10452))
* `azurerm_kubernetes_cluster_node_pool` - support for `enable_host_encryption` ([#10398](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10398))
* `azurerm_lb_backend_address_pool` - support for the `backend_address` block ([#10291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10291))
* `azurerm_redis_cache` - support for the `public_network_access_enabled` property ([#10410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10410))
* `azurerm_role_assignment` - adding validation for that the `scope` is either a Management Group, Subscription, Resource Group or Resource ID ([#10438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10438))
* `azurerm_service_fabric_cluster` - support for the `reverse_proxy_certificate_common_names` block ([#10367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10367))
* `azurerm_monitor_metric_alert` - support for the `skip_metric_validation` property ([#10422](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10422))

BUG FIXES:

* Data Source: `azurerm_api_management` fix an exception with User Assigned Managed Identities ([#10429](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10429))
* `azurerm_api_management_api_diagnostic` - fix a bug where specifying `log_client_ip = false` would not disable the setting ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_key_vault` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_certificate` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_key` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_secret` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_mssql_virtual_machine` - fixing a crash where the KeyVault was nil in the API response ([#10469](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10469))
* `azurerm_storage_account_datasource` - prevent panics from passing in an empty `name` ([#10370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10370))
* `azurerm_storage_data_lake_gen2_filesystem` - change the `ace` property to a TypeSet to ensure consistent ordering ([#10372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10372))
* `azurerm_storage_data_lake_gen2_path` - change the `ace` property to a TypeSet to ensure consistent ordering ([#10372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10372))

## 2.45.1 (January 28, 2021)

BUG FIXES:

* `azurerm_app_service_environment` - prevent a panic when the API returns a nil cluster settings ([#10365](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10365))

## 2.45.0 (January 28, 2021)

FEATURES:

* **New Data Source** `azurerm_search_service` ([#10181](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10181))
* **New Resource:** `azurerm_data_factory_linked_service_snowflake` ([#10239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10239))
* **New Resource:** `azurerm_data_factory_linked_service_azure_table_storage` ([#10305](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10305))
* **New Resource:** `azurerm_iothub_enrichment` ([#9239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9239))
* **New Resource:** `azurerm_iot_security_solution` ([#10034](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10034))
* **New Resource:** `azurerm_media_streaming_policy` ([#10133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10133))
* **New Resource:** `azurerm_spring_cloud_active_deployment` ([#9959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9959))
* **New Resource:** `azurerm_spring_cloud_java_deployment` ([#9959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9959))

IMPROVEMENTS:

* dependencies: updating to `v0.11.17` of `github.com/Azure/go-autorest/autorest` ([#10259](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10259))
* dependencies: updating the `firewall` resources to use the Networking API `2020-07-01` ([#10252](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10252))
* dependencies: updating the `load balancer` resources to use the Networking API version `2020-05-01` ([#10263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10263))
* Data Source: `azurerm_app_service_environment` - export the `cluster_setting` block ([#10303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10303))
* Data Source: `azurerm_key_vault_certificate` - support for the `certificate_data_base64` attribute ([#10275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10275))
* `azurerm_app_service` - support for the propety `number_of_workers` ([#10143](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10143))
* `azurerm_app_service_environment` - support for the `cluster_setting` block ([#10303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10303))
* `azurerm_data_factory_dataset_delimited_text` - support for the `compression_codec` property ([#10182](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10182))
* `azurerm_firewall_policy` - support for the `sku` property ([#10186](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10186))
* `azurerm_iothub` - support for the `enrichment` property ([#9239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9239))
* `azurerm_key_vault` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault` - support both ipv4 and cidr formats for the `network_acls.ip_rules` property ([#10266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10266))
* `azurerm_key_vault_certificate` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_key` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_secret` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_certificate` - support for the `certificate_data_base64` attribute ([#10275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10275))
* `azurerm_linux_virtual_machine` - skipping shutdown for a machine in a failed state ([#10189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10189))
* `azurerm_media_services_account` - support for setting the `storage_authentication_type` field to `System` ([#10133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10133))
* `azurerm_redis_cache` - support multiple availability zones ([#10283](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10283))
* `azurerm_storage_data_lake_gen2_filesystem` - support for the `ace` block ([#9917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9917))
* `azurerm_servicebus_namespace` - will now allow a capacity of `16` for the `Premium` SKU ([#10337](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10337))
* `azurerm_windows_virtual_machine` - skipping shutdown for a machine in a failed state ([#10189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10189))
* `azurerm_linux_virtual_machine_scale_set` - support for the `extensions_time_budget` property ([#10298](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10298))
* `azurerm_windows_virtual_machine_scale_set` - support for the `extensions_time_budget` property ([#10298](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10298))

BUG FIXES:

* `azurerm_iot_time_series_insights_reference_data_set` - the field `data_string_comparison_behavior` is now `ForceNew` ([#10343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10343))
* `azurerm_iot_time_series_insights_reference_data_set` - the `key_property` block is now `ForceNew` ([#10343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10343))
* `azurerm_linux_virtual_machine_scale_set` - fixing an issue where `protected_settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_linux_virtual_machine_scale_set` - fixing an issue where `settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_media_streaming_endpoint` - stopping the streaming endpoint prior to deletion if the endpoint is in a running state ([#10216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10216))
* `azurerm_role_definition` - don't add `scope` to `assignable_scopes` unless none are specified ([#8624](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8624))
* `azurerm_windows_virtual_machine_scale_set` - fixing an issue where `protected_settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_windows_virtual_machine_scale_set` - fixing an issue where `settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))

## 2.44.0 (January 21, 2021)

FEATURES:

* **New Data Source:** `azurerm_iothub` ([#10228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10228))
* **New Resource:** `azurerm_media_content_key_policy` ([#9971](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9971))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/go-autorest` to `v0.11.16` ([#10164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10164))
* dependencies: updating `appconfiguration` to API version `2020-06-01` ([#10176](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10176))
* dependencies: updating `appplatform` to API version `2020-07-01` ([#10175](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10175))
* dependencies: updating `containerservice` to API version `2020-12-01` ([#10171](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10171))
* dependencies: updating `msi` to API version `2018-11-30` ([#10174](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10174))
* Data Source: `azurerm_kubernetes_cluster` - support for the field `user_assigned_identity_id` within the `identity` block ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_api_management` - support additional TLS ciphers within the `security` block ([#9276](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9276))
* `azurerm_api_management_api_diagnostic` - support the `sampling_percentage` property ([#9321](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9321))
* `azurerm_container_group` - support for updating `tags` ([#10210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10210))
* `azurerm_kubernetes_cluster` - the field `type` within the `identity` block can now be set to `UserAssigned` ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_kubernetes_cluster` - support for the field `new_pod_scale_up_delay` within the `auto_scaler_profile` block ([#9291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9291))
* `azurerm_kubernetes_cluster` - support for the field `user_assigned_identity_id` within the `identity` block ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_monitor_autoscale_setting` - now supports the `dimensions` property ([#9795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9795))
* `azurerm_sentinel_alert_rule_scheduled` - now supports the `event_grouping_setting` property ([#10078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10078))

BUG FIXES:

* `azurerm_backup_protected_file_share` - updating to account for a breaking API change ([#9015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9015))
* `azurerm_key_vault_certificate` - fixing a crash when `subject` within the `certificate_policy` block was nil ([#10200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10200))
* `azurerm_user_assigned_identity` - adding a state migration to update the ID format ([#10196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10196))

## 2.43.0 (January 14, 2021)

FEATURES:

* **New Data Source:** `azurerm_sentinel_alert_rule_template` ([#7020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7020))

IMPROVEMENTS:

* Data Source: `azurerm_api_management` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_kubernetes_cluster` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_api_management` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service_slot` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_container_group` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_cosmosdb_account` - support for `analytical_storage_enabled property` ([#10055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10055))
* `azurerm_cosmosdb_gremlin_graph` - support the `default_ttl` property ([#10159](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10159))
* `azurerm_data_factory` - support for `public_network_enabled` ([#9605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9605))
* `azurerm_data_factory_dataset_delimited_text` - support for the `compression_type` property ([#10070](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10070))
* `azurerm_data_factory_linked_service_sql_server`: support for the `key_vault_password` block ([#10032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10032))
* `azurerm_eventgrid_domain` - support for the `public_network_access_enabled` and `inbound_ip_rule` properties  ([#9922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9922))
* `azurerm_eventgrid_topic` - support for the `public_network_access_enabled` and `inbound_ip_rule` properties  ([#9922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9922))
* `azurerm_eventhub_namespace` - support the `trusted_service_access_enabled` property ([#10169](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10169))
* `azurerm_function_app` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_function_app_slot` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kusto_cluster` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_security_center_automation` - the field `event_source` within the `source` block now supports `SecureScoreControls ` and `SecureScores` ([#10126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10126))
* `azurerm_synapse_workspace` - support for the `sql_identity_control_enabled` property ([#10033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10033))
* `azurerm_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))

BUG FIXES:

* Data Source: `azurerm_log_analytics_workspace` - returning the Resource ID in the correct casing ([#10162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10162))
* `azurerm_advanced_threat_protection` - fix a regression in the Resouce ID format ([#10190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10190))
* `azurerm_api_management` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service_slot` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_application_gateway` - ensuring the casing on `identity_ids` within the `identity` block ([#10031](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10031))
* `azurerm_blueprint_assignment` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_container_group` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_databricks_workspace` - changing the sku no longer always forces a new resource to be created ([#9541](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9541))
* `azurerm_function_app` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_function_app_slot` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kubernetes_cluster` - ensuring the casing of the `user_assigned_identity_id` field within the `kubelet_identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kusto_cluster` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_monitor_diagnostic_setting` - handling mixed casing of the EventHub Namespace Authorization Rule ID ([#10104](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10104))
* `azurerm_mssql_virtual_machine` - address persistent diff and use relative expiry for service principal password ([#10125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10125))
* `azurerm_role_assignment` - fix race condition in read after create ([#10134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10134))
* `azurerm_role_definition` - address eventual consistency issues in update and delete ([#10170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10170))
* `azurerm_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))

## 2.42.0 (January 08, 2021)

BREAKING CHANGES

* `azurerm_key_vault` - the field `soft_delete_enabled` is now defaulted to `true` to match the breaking change in the Azure API where Key Vaults now have Soft Delete enabled by default, which cannot be disabled. This property is now non-functional, defaults to `true` and will be removed in version 3.0 of the Azure Provider. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_key_vault` - the field `soft_delete_retention_days` is now defaulted to `90` days to match the Azure API behaviour, as the Azure API does not return a value for this field when not explicitly configured, so defaulting this removes a diff with `0`. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))

FEATURES:

* **New Data Source:** `azurerm_eventgrid_domain_topic` ([#10050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10050))
* **New Data Source:** `azurerm_ssh_public_key` ([#9842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9842))
* **New Resource:** `azurerm_data_factory_linked_service_synapse` ([#9928](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9928))
* **New Resource:** `azurerm_disk_access` ([#9889](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9889))
* **New Resource:** `azurerm_media_streaming_locator` ([#9992](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9992))
* **New Resource:** `azurerm_sentinel_alert_rule_fusion` ([#9829](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9829))
* **New Resource:** `azurerm_ssh_public_key` ([#9842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9842))

IMPROVEMENTS:

* batch: updating to API version `2020-03-01` ([#10036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10036))
* dependencies: upgrading to `v49.2.0` of `github.com/Azure/azure-sdk-for-go` ([#10042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10042))
* dependencies: upgrading to `v0.15.1` of `github.com/tombuildsstuff/giovanni` ([#10035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10035))
* Data Source: `azurerm_hdinsight_cluster` - support for the `kafka_rest_proxy_endpoint` property ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* Data Source: `azurerm_databricks_workspace` - support for the `tags` property ([#9933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9933))
* Data Source: `azurerm_subscription` - support for the `tags` property ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* `azurerm_app_service` - now supports  `detailed_error_mesage_enabled` and `failed_request_tracing_enabled ` logs settings ([#9162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9162))
* `azurerm_app_service` - now supports  `service_tag` in `ip_restriction` blocks ([#9609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9609))
* `azurerm_app_service_slot` - now supports  `detailed_error_mesage_enabled` and `failed_request_tracing_enabled ` logs settings ([#9162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9162))
* `azurerm_batch_pool` support for the `public_address_provisioning_type` property ([#10036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10036))
* `azurerm_api_management` - support `Consumption_0` for the `sku_name` property ([#6868](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6868))
* `azurerm_cdn_endpoint` - only send `content_types_to_compress` and `geo_filter` to the API when actually set ([#9902](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9902))
* `azurerm_cosmosdb_mongo_collection` - correctly read back the `_id` index when mongo 3.6 ([#8690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8690))
* `azurerm_container_group` - support for the `volume.empty_dir` property ([#9836](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9836))
* `azurerm_data_factory_linked_service_azure_file_storage` - support for the `file_share` property ([#9934](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9934))
* `azurerm_dedicated_host` - support for addtional `sku_name` values ([#9951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9951))
* `azurerm_devspace_controller` - deprecating since new DevSpace Controllers can no longer be provisioned, this will be removed in version 3.0 of the Azure Provider ([#10049](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10049))
* `azurerm_function_app` - make `pre_warmed_instance_count` computed to use azure's default ([#9069](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9069))
* `azurerm_function_app` - now supports  `service_tag` in `ip_restriction` blocks ([#9609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9609))
* `azurerm_hdinsight_hadoop_cluster` - allow the value `Standard_D4a_V4` for the `vm_type` property ([#10000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10000))
* `azurerm_hdinsight_kafka_cluster` - support for the `rest_proxy` and `kafka_management_node` blocks ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* `azurerm_key_vault` - the field `soft_delete_enabled` is now defaulted to `true` to match the Azure API behaviour where Soft Delete is force-enabled and can no longer be disabled. This field is deprecated, can be safely removed from your Terraform Configuration, and will be removed in version 3.0 of the Azure Provider. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_kubernetes_cluster` - add support for network_mode ([#8828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8828))
* `azurerm_log_analytics_linked_service` - add validation for resource ID type ([#9932](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9932))
* `azurerm_log_analytics_linked_service` - update validation to use generated validate functions ([#9950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9950))
* `azurerm_monitor_diagnostic_setting` - validation that `eventhub_authorization_rule_id` is a EventHub Namespace Authorization Rule ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_monitor_diagnostic_setting` - validation that `log_analytics_workspace_id` is a Log Analytics Workspace ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_monitor_diagnostic_setting` - validation that `storage_account_id` is a Storage Account ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_network_security_rule` - increase allowed the number of `application_security_group` blocks allowed ([#9884](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9884))
* `azurerm_sentinel_alert_rule_ms_security_incident` - support the `alert_rule_template_guid` and `display_name_exclude_filter` properties ([#9797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9797))
* `azurerm_sentinel_alert_rule_scheduled` - support for the `alert_rule_template_guid` property ([#9712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9712))
* `azurerm_sentinel_alert_rule_scheduled` - support for creating incidents ([#8564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8564))
* `azurerm_spring_cloud_app` - support the properties `https_only`, `is_public`, and `persistent_disk` ([#9957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9957))
* `azurerm_subscription` - support for the `tags` property ([#9047](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9047))
* `azurerm_synapse_workspace` - support for the `managed_resource_group_name` property ([#10017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10017))
* `azurerm_traffic_manager_profile` - support for the `traffic_view_enabled` property ([#10005](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10005))

BUG FIXES:

provider: will not correctly register the `Microsoft.Blueprint` and `Microsoft.HealthcareApis` RPs ([#10062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10062))
* `azurerm_application_gateway` - allow `750` for `file_upload_limit_mb` when the sku is `WAF_v2` ([#8753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8753))
* `azurerm_firewall_policy_rule_collection_group` - correctly validate the `network_rule_collection.destination_ports` property ([#9490](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9490))
* `azurerm_cdn_endpoint` - changing many `delivery_rule` condition `match_values` to optional ([#8850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8850))
* `azurerm_cosmosdb_account` - always include `key_vault_id` in update requests for azure policy enginer compatibility ([#9966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9966))
* `azurerm_cosmosdb_table` - do not call the throughput api when serverless ([#9749](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9749))
* `azurerm_key_vault` - the field `soft_delete_retention_days` is now defaulted to `90` days to match the Azure API behaviour. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_kubernetes_cluster` - parse oms `log_analytics_workspace_id` to ensure correct casing ([#9976](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9976))
* `azurerm_role_assignment` fix crash in retry logic ([#10051](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10051))
* `azurerm_storage_account` - allow hns when `account_tier` is `Premium` ([#9548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9548))
* `azurerm_storage_share_file` - allowing files smaller than 4KB to be uploaded ([#10035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10035))

## 2.41.0 (December 17, 2020)

UPGRADE NOTES:

* `azurerm_key_vault` - Azure will be introducing a breaking change on December 31st, 2020 by force-enabling Soft Delete on all new and existing Key Vaults. To workaround this, this release of the Azure Provider still allows you to configure Soft Delete on before this date (but once this is enabled this cannot be disabled). Since new Key Vaults will automatically be provisioned using Soft Delete in the future, and existing Key Vaults will be upgraded - a future release will deprecate the `soft_delete_enabled` field and default this to true early in 2021. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_certificate` - Terraform will now attempt to `purge` Certificates during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - Terraform will now attempt to `purge` Keys during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - Terraform will now attempt to `purge` Secrets during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))

FEATURES:

* **New Resource:** `azurerm_eventgrid_system_topic_event_subscription` ([#9852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9852))
* **New Resource:** `azurerm_media_job` ([#9859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9859))
* **New Resource:** `azurerm_media_streaming_endpoint` ([#9537](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9537))
* **New Resource:** `azurerm_subnet_service_endpoint_storage_policy` ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* **New Resource:** `azurerm_synapse_managed_private_endpoint` ([#9260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9260))

IMPROVEMENTS:

* `azurerm_app_service` - Add support for `outbound_ip_address_list` and `possible_outbound_ip_address_list` ([#9871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9871))
* `azurerm_disk_encryption_set` - support for updating `key_vault_key_id` ([#7913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7913))
* `azurerm_iot_time_series_insights_gen2_environment` - exposing `data_access_fqdn` ([#9848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9848))
* `azurerm_key_vault_certificate` - performing a "purge" of the Certificate during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - performing a "purge" of the Key during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - performing a "purge" of the Secret during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_linked_service` - Add new fields `workspace_id`, `read_access_id`, and `write_access_id` ([#9410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9410))
* `azurerm_linux_virtual_machine` - Normalise SSH keys to cover VM import cases ([#9897](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9897))
* `azurerm_subnet` - support for the `service_endpoint_policy` block ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* `azurerm_traffic_manager_profile` - support for new field `max_return` and support for `traffic_routing_method` to be `MultiValue` ([#9487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9487))

BUG FIXES:

* `azurerm_key_vault_certificate` - reading `dns_names` and `emails` within the `subject_alternative_names` block from the Certificate if not returned from the API ([#8631](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8631))
* `azurerm_key_vault_certificate` - polling until the Certificate is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - polling until the Key is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` -  polling until the Secret is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_workspace` - adding a state migration to correctly update the Resource ID ([#9853](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9853))

---

For information on changes between the v2.40.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
