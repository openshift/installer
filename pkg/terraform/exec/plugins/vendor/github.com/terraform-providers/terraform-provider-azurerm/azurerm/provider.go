package azurerm

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
			},

			// Client Certificate specific fields
			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
			},

			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},

			// Managed Service Identity specific fields
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
			},
			"msi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
			},

			// Managed Tracking GUID for User-agent
			"partner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_PARTNER_ID", ""),
				ValidateFunc: validate.UUIDOrEmpty,
			},

			// Advanced feature flags
			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_CREDENTIALS_VALIDATION", false),
			},

			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"azurerm_api_management":                         dataSourceApiManagementService(),
			"azurerm_api_management_api":                     dataSourceApiManagementApi(),
			"azurerm_api_management_group":                   dataSourceApiManagementGroup(),
			"azurerm_api_management_product":                 dataSourceApiManagementProduct(),
			"azurerm_api_management_user":                    dataSourceArmApiManagementUser(),
			"azurerm_app_service_plan":                       dataSourceAppServicePlan(),
			"azurerm_app_service":                            dataSourceArmAppService(),
			"azurerm_application_insights":                   dataSourceArmApplicationInsights(),
			"azurerm_application_security_group":             dataSourceArmApplicationSecurityGroup(),
			"azurerm_availability_set":                       dataSourceArmAvailabilitySet(),
			"azurerm_azuread_application":                    dataSourceArmAzureADApplication(),
			"azurerm_azuread_service_principal":              dataSourceArmActiveDirectoryServicePrincipal(),
			"azurerm_batch_account":                          dataSourceArmBatchAccount(),
			"azurerm_batch_certificate":                      dataSourceArmBatchCertificate(),
			"azurerm_batch_pool":                             dataSourceArmBatchPool(),
			"azurerm_builtin_role_definition":                dataSourceArmBuiltInRoleDefinition(),
			"azurerm_cdn_profile":                            dataSourceArmCdnProfile(),
			"azurerm_client_config":                          dataSourceArmClientConfig(),
			"azurerm_container_registry":                     dataSourceArmContainerRegistry(),
			"azurerm_cosmosdb_account":                       dataSourceArmCosmosDBAccount(),
			"azurerm_data_lake_store":                        dataSourceArmDataLakeStoreAccount(),
			"azurerm_dev_test_lab":                           dataSourceArmDevTestLab(),
			"azurerm_dns_zone":                               dataSourceArmDnsZone(),
			"azurerm_eventhub_namespace":                     dataSourceEventHubNamespace(),
			"azurerm_express_route_circuit":                  dataSourceArmExpressRouteCircuit(),
			"azurerm_firewall":                               dataSourceArmFirewall(),
			"azurerm_image":                                  dataSourceArmImage(),
			"azurerm_hdinsight_cluster":                      dataSourceArmHDInsightSparkCluster(),
			"azurerm_key_vault_access_policy":                dataSourceArmKeyVaultAccessPolicy(),
			"azurerm_key_vault_key":                          dataSourceArmKeyVaultKey(),
			"azurerm_key_vault_secret":                       dataSourceArmKeyVaultSecret(),
			"azurerm_key_vault":                              dataSourceArmKeyVault(),
			"azurerm_kubernetes_cluster":                     dataSourceArmKubernetesCluster(),
			"azurerm_lb":                                     dataSourceArmLoadBalancer(),
			"azurerm_lb_backend_address_pool":                dataSourceArmLoadBalancerBackendAddressPool(),
			"azurerm_log_analytics_workspace":                dataSourceLogAnalyticsWorkspace(),
			"azurerm_logic_app_workflow":                     dataSourceArmLogicAppWorkflow(),
			"azurerm_managed_disk":                           dataSourceArmManagedDisk(),
			"azurerm_management_group":                       dataSourceArmManagementGroup(),
			"azurerm_monitor_action_group":                   dataSourceArmMonitorActionGroup(),
			"azurerm_monitor_diagnostic_categories":          dataSourceArmMonitorDiagnosticCategories(),
			"azurerm_monitor_log_profile":                    dataSourceArmMonitorLogProfile(),
			"azurerm_network_interface":                      dataSourceArmNetworkInterface(),
			"azurerm_network_security_group":                 dataSourceArmNetworkSecurityGroup(),
			"azurerm_network_watcher":                        dataSourceArmNetworkWatcher(),
			"azurerm_notification_hub_namespace":             dataSourceNotificationHubNamespace(),
			"azurerm_notification_hub":                       dataSourceNotificationHub(),
			"azurerm_platform_image":                         dataSourceArmPlatformImage(),
			"azurerm_policy_definition":                      dataSourceArmPolicyDefinition(),
			"azurerm_public_ip":                              dataSourceArmPublicIP(),
			"azurerm_public_ips":                             dataSourceArmPublicIPs(),
			"azurerm_recovery_services_vault":                dataSourceArmRecoveryServicesVault(),
			"azurerm_recovery_services_protection_policy_vm": dataSourceArmRecoveryServicesProtectionPolicyVm(),
			"azurerm_resource_group":                         dataSourceArmResourceGroup(),
			"azurerm_role_definition":                        dataSourceArmRoleDefinition(),
			"azurerm_route_table":                            dataSourceArmRouteTable(),
			"azurerm_scheduler_job_collection":               dataSourceArmSchedulerJobCollection(),
			"azurerm_servicebus_namespace":                   dataSourceArmServiceBusNamespace(),
			"azurerm_shared_image_gallery":                   dataSourceArmSharedImageGallery(),
			"azurerm_shared_image_version":                   dataSourceArmSharedImageVersion(),
			"azurerm_shared_image":                           dataSourceArmSharedImage(),
			"azurerm_snapshot":                               dataSourceArmSnapshot(),
			"azurerm_stream_analytics_job":                   dataSourceArmStreamAnalyticsJob(),
			"azurerm_storage_account_sas":                    dataSourceArmStorageAccountSharedAccessSignature(),
			"azurerm_storage_account":                        dataSourceArmStorageAccount(),
			"azurerm_subnet":                                 dataSourceArmSubnet(),
			"azurerm_subscription":                           dataSourceArmSubscription(),
			"azurerm_subscriptions":                          dataSourceArmSubscriptions(),
			"azurerm_traffic_manager_geographical_location":  dataSourceArmTrafficManagerGeographicalLocation(),
			"azurerm_virtual_machine":                        dataSourceArmVirtualMachine(),
			"azurerm_virtual_network_gateway":                dataSourceArmVirtualNetworkGateway(),
			"azurerm_virtual_network":                        dataSourceArmVirtualNetwork(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"azurerm_api_management":                         resourceArmApiManagementService(),
			"azurerm_api_management_api":                     resourceArmApiManagementApi(),
			"azurerm_api_management_api_operation":           resourceArmApiManagementApiOperation(),
			"azurerm_api_management_api_version_set":         resourceArmApiManagementApiVersionSet(),
			"azurerm_api_management_authorization_server":    resourceArmApiManagementAuthorizationServer(),
			"azurerm_api_management_certificate":             resourceArmApiManagementCertificate(),
			"azurerm_api_management_group":                   resourceArmApiManagementGroup(),
			"azurerm_api_management_group_user":              resourceArmApiManagementGroupUser(),
			"azurerm_api_management_logger":                  resourceArmApiManagementLogger(),
			"azurerm_api_management_openid_connect_provider": resourceArmApiManagementOpenIDConnectProvider(),
			"azurerm_api_management_product":                 resourceArmApiManagementProduct(),
			"azurerm_api_management_product_api":             resourceArmApiManagementProductApi(),
			"azurerm_api_management_product_group":           resourceArmApiManagementProductGroup(),
			"azurerm_api_management_property":                resourceArmApiManagementProperty(),
			"azurerm_api_management_subscription":            resourceArmApiManagementSubscription(),
			"azurerm_api_management_user":                    resourceArmApiManagementUser(),
			"azurerm_app_service_active_slot":                resourceArmAppServiceActiveSlot(),
			"azurerm_app_service_custom_hostname_binding":    resourceArmAppServiceCustomHostnameBinding(),
			"azurerm_app_service_plan":                       resourceArmAppServicePlan(),
			"azurerm_app_service_slot":                       resourceArmAppServiceSlot(),
			"azurerm_app_service":                            resourceArmAppService(),
			"azurerm_application_gateway":                    resourceArmApplicationGateway(),
			"azurerm_application_insights_api_key":           resourceArmApplicationInsightsAPIKey(),
			"azurerm_application_insights":                   resourceArmApplicationInsights(),
			"azurerm_application_security_group":             resourceArmApplicationSecurityGroup(),
			"azurerm_automation_account":                     resourceArmAutomationAccount(),
			"azurerm_automation_credential":                  resourceArmAutomationCredential(),
			"azurerm_automation_dsc_configuration":           resourceArmAutomationDscConfiguration(),
			"azurerm_automation_dsc_nodeconfiguration":       resourceArmAutomationDscNodeConfiguration(),
			"azurerm_automation_module":                      resourceArmAutomationModule(),
			"azurerm_automation_runbook":                     resourceArmAutomationRunbook(),
			"azurerm_automation_schedule":                    resourceArmAutomationSchedule(),
			"azurerm_autoscale_setting":                      resourceArmAutoScaleSetting(),
			"azurerm_availability_set":                       resourceArmAvailabilitySet(),
			"azurerm_azuread_application":                    resourceArmActiveDirectoryApplication(),
			"azurerm_azuread_service_principal_password":     resourceArmActiveDirectoryServicePrincipalPassword(),
			"azurerm_azuread_service_principal":              resourceArmActiveDirectoryServicePrincipal(),
			"azurerm_batch_account":                          resourceArmBatchAccount(),
			"azurerm_batch_certificate":                      resourceArmBatchCertificate(),
			"azurerm_batch_pool":                             resourceArmBatchPool(),
			"azurerm_cdn_endpoint":                           resourceArmCdnEndpoint(),
			"azurerm_cdn_profile":                            resourceArmCdnProfile(),
			"azurerm_cognitive_account":                      resourceArmCognitiveAccount(),
			"azurerm_connection_monitor":                     resourceArmConnectionMonitor(),
			"azurerm_container_group":                        resourceArmContainerGroup(),
			"azurerm_container_registry":                     resourceArmContainerRegistry(),
			"azurerm_container_service":                      resourceArmContainerService(),
			"azurerm_cosmosdb_account":                       resourceArmCosmosDBAccount(),
			"azurerm_data_factory":                           resourceArmDataFactory(),
			"azurerm_data_factory_dataset_mysql":             resourceArmDataFactoryDatasetMySQL(),
			"azurerm_data_factory_dataset_postgresql":        resourceArmDataFactoryDatasetPostgreSQL(),
			"azurerm_data_factory_dataset_sql_server_table":  resourceArmDataFactoryDatasetSQLServerTable(),
			"azurerm_data_factory_linked_service_mysql":      resourceArmDataFactoryLinkedServiceMySQL(),
			"azurerm_data_factory_linked_service_postgresql": resourceArmDataFactoryLinkedServicePostgreSQL(),
			"azurerm_data_factory_linked_service_sql_server": resourceArmDataFactoryLinkedServiceSQLServer(),
			"azurerm_data_factory_pipeline":                  resourceArmDataFactoryPipeline(),
			"azurerm_data_lake_analytics_account":            resourceArmDataLakeAnalyticsAccount(),
			"azurerm_data_lake_analytics_firewall_rule":      resourceArmDataLakeAnalyticsFirewallRule(),
			"azurerm_data_lake_store_file":                   resourceArmDataLakeStoreFile(),
			"azurerm_data_lake_store_firewall_rule":          resourceArmDataLakeStoreFirewallRule(),
			"azurerm_data_lake_store":                        resourceArmDataLakeStore(),
			"azurerm_databricks_workspace":                   resourceArmDatabricksWorkspace(),
			"azurerm_ddos_protection_plan":                   resourceArmDDoSProtectionPlan(),
			"azurerm_dev_test_lab":                           resourceArmDevTestLab(),
			"azurerm_dev_test_linux_virtual_machine":         resourceArmDevTestLinuxVirtualMachine(),
			"azurerm_dev_test_policy":                        resourceArmDevTestPolicy(),
			"azurerm_dev_test_virtual_network":               resourceArmDevTestVirtualNetwork(),
			"azurerm_dev_test_windows_virtual_machine":       resourceArmDevTestWindowsVirtualMachine(),
			"azurerm_devspace_controller":                    resourceArmDevSpaceController(),
			"azurerm_dns_a_record":                           resourceArmDnsARecord(),
			"azurerm_dns_aaaa_record":                        resourceArmDnsAAAARecord(),
			"azurerm_dns_caa_record":                         resourceArmDnsCaaRecord(),
			"azurerm_dns_cname_record":                       resourceArmDnsCNameRecord(),
			"azurerm_dns_mx_record":                          resourceArmDnsMxRecord(),
			"azurerm_dns_ns_record":                          resourceArmDnsNsRecord(),
			"azurerm_dns_ptr_record":                         resourceArmDnsPtrRecord(),
			"azurerm_dns_srv_record":                         resourceArmDnsSrvRecord(),
			"azurerm_dns_txt_record":                         resourceArmDnsTxtRecord(),
			"azurerm_dns_zone":                               resourceArmDnsZone(),
			"azurerm_eventgrid_domain":                       resourceArmEventGridDomain(),
			"azurerm_eventgrid_event_subscription":           resourceArmEventGridEventSubscription(),
			"azurerm_eventgrid_topic":                        resourceArmEventGridTopic(),
			"azurerm_eventhub_authorization_rule":            resourceArmEventHubAuthorizationRule(),
			"azurerm_eventhub_consumer_group":                resourceArmEventHubConsumerGroup(),
			"azurerm_eventhub_namespace_authorization_rule":  resourceArmEventHubNamespaceAuthorizationRule(),
			"azurerm_eventhub_namespace":                     resourceArmEventHubNamespace(),
			"azurerm_eventhub":                               resourceArmEventHub(),
			"azurerm_express_route_circuit_authorization":    resourceArmExpressRouteCircuitAuthorization(),
			"azurerm_express_route_circuit_peering":          resourceArmExpressRouteCircuitPeering(),
			"azurerm_express_route_circuit":                  resourceArmExpressRouteCircuit(),
			"azurerm_firewall_application_rule_collection":   resourceArmFirewallApplicationRuleCollection(),
			"azurerm_firewall_network_rule_collection":       resourceArmFirewallNetworkRuleCollection(),
			"azurerm_firewall":                               resourceArmFirewall(),
			"azurerm_function_app":                           resourceArmFunctionApp(),
			"azurerm_hdinsight_hadoop_cluster":               resourceArmHDInsightHadoopCluster(),
			"azurerm_hdinsight_hbase_cluster":                resourceArmHDInsightHBaseCluster(),
			"azurerm_hdinsight_interactive_query_cluster":    resourceArmHDInsightInteractiveQueryCluster(),
			"azurerm_hdinsight_kafka_cluster":                resourceArmHDInsightKafkaCluster(),
			"azurerm_hdinsight_ml_services_cluster":          resourceArmHDInsightMLServicesCluster(),
			"azurerm_hdinsight_rserver_cluster":              resourceArmHDInsightRServerCluster(),
			"azurerm_hdinsight_spark_cluster":                resourceArmHDInsightSparkCluster(),
			"azurerm_hdinsight_storm_cluster":                resourceArmHDInsightStormCluster(),
			"azurerm_image":                                  resourceArmImage(),
			"azurerm_iothub_consumer_group":                  resourceArmIotHubConsumerGroup(),
			"azurerm_iothub":                                 resourceArmIotHub(),
			"azurerm_iothub_shared_access_policy":            resourceArmIotHubSharedAccessPolicy(),
			"azurerm_key_vault_access_policy":                resourceArmKeyVaultAccessPolicy(),
			"azurerm_key_vault_certificate":                  resourceArmKeyVaultCertificate(),
			"azurerm_key_vault_key":                          resourceArmKeyVaultKey(),
			"azurerm_key_vault_secret":                       resourceArmKeyVaultSecret(),
			"azurerm_key_vault":                              resourceArmKeyVault(),
			"azurerm_kubernetes_cluster":                     resourceArmKubernetesCluster(),
			"azurerm_lb_backend_address_pool":                resourceArmLoadBalancerBackendAddressPool(),
			"azurerm_lb_nat_pool":                            resourceArmLoadBalancerNatPool(),
			"azurerm_lb_nat_rule":                            resourceArmLoadBalancerNatRule(),
			"azurerm_lb_probe":                               resourceArmLoadBalancerProbe(),
			"azurerm_lb_outbound_rule":                       resourceArmLoadBalancerOutboundRule(),
			"azurerm_lb_rule":                                resourceArmLoadBalancerRule(),
			"azurerm_lb":                                     resourceArmLoadBalancer(),
			"azurerm_local_network_gateway":                  resourceArmLocalNetworkGateway(),
			"azurerm_log_analytics_solution":                 resourceArmLogAnalyticsSolution(),
			"azurerm_log_analytics_linked_service":           resourceArmLogAnalyticsLinkedService(),
			"azurerm_log_analytics_workspace_linked_service": resourceArmLogAnalyticsWorkspaceLinkedService(),
			"azurerm_log_analytics_workspace":                resourceArmLogAnalyticsWorkspace(),
			"azurerm_logic_app_action_custom":                resourceArmLogicAppActionCustom(),
			"azurerm_logic_app_action_http":                  resourceArmLogicAppActionHTTP(),
			"azurerm_logic_app_trigger_custom":               resourceArmLogicAppTriggerCustom(),
			"azurerm_logic_app_trigger_http_request":         resourceArmLogicAppTriggerHttpRequest(),
			"azurerm_logic_app_trigger_recurrence":           resourceArmLogicAppTriggerRecurrence(),
			"azurerm_logic_app_workflow":                     resourceArmLogicAppWorkflow(),
			"azurerm_managed_disk":                           resourceArmManagedDisk(),
			"azurerm_management_group":                       resourceArmManagementGroup(),
			"azurerm_management_lock":                        resourceArmManagementLock(),
			"azurerm_mariadb_database":                       resourceArmMariaDbDatabase(),
			"azurerm_mariadb_server":                         resourceArmMariaDbServer(),
			"azurerm_media_services_account":                 resourceArmMediaServicesAccount(),
			"azurerm_metric_alertrule":                       resourceArmMetricAlertRule(),
			"azurerm_monitor_autoscale_setting":              resourceArmMonitorAutoScaleSetting(),
			"azurerm_monitor_action_group":                   resourceArmMonitorActionGroup(),
			"azurerm_monitor_activity_log_alert":             resourceArmMonitorActivityLogAlert(),
			"azurerm_monitor_diagnostic_setting":             resourceArmMonitorDiagnosticSetting(),
			"azurerm_monitor_log_profile":                    resourceArmMonitorLogProfile(),
			"azurerm_monitor_metric_alert":                   resourceArmMonitorMetricAlert(),
			"azurerm_monitor_metric_alertrule":               resourceArmMonitorMetricAlertRule(),
			"azurerm_mssql_elasticpool":                      resourceArmMsSqlElasticPool(),
			"azurerm_mysql_configuration":                    resourceArmMySQLConfiguration(),
			"azurerm_mysql_database":                         resourceArmMySqlDatabase(),
			"azurerm_mysql_firewall_rule":                    resourceArmMySqlFirewallRule(),
			"azurerm_mysql_server":                           resourceArmMySqlServer(),
			"azurerm_mysql_virtual_network_rule":             resourceArmMySqlVirtualNetworkRule(),
			"azurerm_network_connection_monitor":             resourceArmNetworkConnectionMonitor(),
			"azurerm_network_ddos_protection_plan":           resourceArmNetworkDDoSProtectionPlan(),
			"azurerm_network_interface_application_gateway_backend_address_pool_association": resourceArmNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(),
			"azurerm_network_interface_application_security_group_association":               resourceArmNetworkInterfaceApplicationSecurityGroupAssociation(),
			"azurerm_network_interface_backend_address_pool_association":                     resourceArmNetworkInterfaceBackendAddressPoolAssociation(),
			"azurerm_network_interface_nat_rule_association":                                 resourceArmNetworkInterfaceNatRuleAssociation(),
			"azurerm_network_interface":                                                      resourceArmNetworkInterface(),
			"azurerm_network_packet_capture":                                                 resourceArmNetworkPacketCapture(),
			"azurerm_network_security_group":                                                 resourceArmNetworkSecurityGroup(),
			"azurerm_network_security_rule":                                                  resourceArmNetworkSecurityRule(),
			"azurerm_network_watcher":                                                        resourceArmNetworkWatcher(),
			"azurerm_notification_hub_authorization_rule":                                    resourceArmNotificationHubAuthorizationRule(),
			"azurerm_notification_hub_namespace":                                             resourceArmNotificationHubNamespace(),
			"azurerm_notification_hub":                                                       resourceArmNotificationHub(),
			"azurerm_packet_capture":                                                         resourceArmPacketCapture(),
			"azurerm_policy_assignment":                                                      resourceArmPolicyAssignment(),
			"azurerm_policy_definition":                                                      resourceArmPolicyDefinition(),
			"azurerm_policy_set_definition":                                                  resourceArmPolicySetDefinition(),
			"azurerm_postgresql_configuration":                                               resourceArmPostgreSQLConfiguration(),
			"azurerm_postgresql_database":                                                    resourceArmPostgreSQLDatabase(),
			"azurerm_postgresql_firewall_rule":                                               resourceArmPostgreSQLFirewallRule(),
			"azurerm_postgresql_server":                                                      resourceArmPostgreSQLServer(),
			"azurerm_postgresql_virtual_network_rule":                                        resourceArmPostgreSQLVirtualNetworkRule(),
			"azurerm_public_ip":                                                              resourceArmPublicIp(),
			"azurerm_public_ip_prefix":                                                       resourceArmPublicIpPrefix(),
			"azurerm_recovery_services_protected_vm":                                         resourceArmRecoveryServicesProtectedVm(),
			"azurerm_recovery_services_protection_policy_vm":                                 resourceArmRecoveryServicesProtectionPolicyVm(),
			"azurerm_recovery_services_vault":                                                resourceArmRecoveryServicesVault(),
			"azurerm_redis_cache":                                                            resourceArmRedisCache(),
			"azurerm_redis_firewall_rule":                                                    resourceArmRedisFirewallRule(),
			"azurerm_relay_namespace":                                                        resourceArmRelayNamespace(),
			"azurerm_resource_group":                                                         resourceArmResourceGroup(),
			"azurerm_role_assignment":                                                        resourceArmRoleAssignment(),
			"azurerm_role_definition":                                                        resourceArmRoleDefinition(),
			"azurerm_route_table":                                                            resourceArmRouteTable(),
			"azurerm_route":                                                                  resourceArmRoute(),
			"azurerm_scheduler_job_collection":                                               resourceArmSchedulerJobCollection(),
			"azurerm_scheduler_job":                                                          resourceArmSchedulerJob(),
			"azurerm_search_service":                                                         resourceArmSearchService(),
			"azurerm_security_center_contact":                                                resourceArmSecurityCenterContact(),
			"azurerm_security_center_subscription_pricing":                                   resourceArmSecurityCenterSubscriptionPricing(),
			"azurerm_security_center_workspace":                                              resourceArmSecurityCenterWorkspace(),
			"azurerm_service_fabric_cluster":                                                 resourceArmServiceFabricCluster(),
			"azurerm_servicebus_namespace_authorization_rule":                                resourceArmServiceBusNamespaceAuthorizationRule(),
			"azurerm_servicebus_namespace":                                                   resourceArmServiceBusNamespace(),
			"azurerm_servicebus_queue_authorization_rule":                                    resourceArmServiceBusQueueAuthorizationRule(),
			"azurerm_servicebus_queue":                                                       resourceArmServiceBusQueue(),
			"azurerm_servicebus_subscription_rule":                                           resourceArmServiceBusSubscriptionRule(),
			"azurerm_servicebus_subscription":                                                resourceArmServiceBusSubscription(),
			"azurerm_servicebus_topic_authorization_rule":                                    resourceArmServiceBusTopicAuthorizationRule(),
			"azurerm_servicebus_topic":                                                       resourceArmServiceBusTopic(),
			"azurerm_shared_image_gallery":                                                   resourceArmSharedImageGallery(),
			"azurerm_shared_image_version":                                                   resourceArmSharedImageVersion(),
			"azurerm_shared_image":                                                           resourceArmSharedImage(),
			"azurerm_signalr_service":                                                        resourceArmSignalRService(),
			"azurerm_snapshot":                                                               resourceArmSnapshot(),
			"azurerm_sql_active_directory_administrator":                                     resourceArmSqlAdministrator(),
			"azurerm_sql_database":                                                           resourceArmSqlDatabase(),
			"azurerm_sql_elasticpool":                                                        resourceArmSqlElasticPool(),
			"azurerm_sql_firewall_rule":                                                      resourceArmSqlFirewallRule(),
			"azurerm_sql_server":                                                             resourceArmSqlServer(),
			"azurerm_sql_virtual_network_rule":                                               resourceArmSqlVirtualNetworkRule(),
			"azurerm_storage_account":                                                        resourceArmStorageAccount(),
			"azurerm_storage_blob":                                                           resourceArmStorageBlob(),
			"azurerm_storage_container":                                                      resourceArmStorageContainer(),
			"azurerm_storage_queue":                                                          resourceArmStorageQueue(),
			"azurerm_storage_share":                                                          resourceArmStorageShare(),
			"azurerm_storage_table":                                                          resourceArmStorageTable(),
			"azurerm_stream_analytics_job":                                                   resourceArmStreamAnalyticsJob(),
			"azurerm_stream_analytics_function_javascript_udf":                               resourceArmStreamAnalyticsFunctionUDF(),
			"azurerm_stream_analytics_output_blob":                                           resourceArmStreamAnalyticsOutputBlob(),
			"azurerm_stream_analytics_output_eventhub":                                       resourceArmStreamAnalyticsOutputEventHub(),
			"azurerm_stream_analytics_output_servicebus_queue":                               resourceArmStreamAnalyticsOutputServiceBusQueue(),
			"azurerm_stream_analytics_stream_input_blob":                                     resourceArmStreamAnalyticsStreamInputBlob(),
			"azurerm_stream_analytics_stream_input_eventhub":                                 resourceArmStreamAnalyticsStreamInputEventHub(),
			"azurerm_stream_analytics_stream_input_iothub":                                   resourceArmStreamAnalyticsStreamInputIoTHub(),
			"azurerm_subnet_network_security_group_association":                              resourceArmSubnetNetworkSecurityGroupAssociation(),
			"azurerm_subnet_route_table_association":                                         resourceArmSubnetRouteTableAssociation(),
			"azurerm_subnet":                                                                 resourceArmSubnet(),
			"azurerm_template_deployment":                                                    resourceArmTemplateDeployment(),
			"azurerm_traffic_manager_endpoint":                                               resourceArmTrafficManagerEndpoint(),
			"azurerm_traffic_manager_profile":                                                resourceArmTrafficManagerProfile(),
			"azurerm_user_assigned_identity":                                                 resourceArmUserAssignedIdentity(),
			"azurerm_virtual_machine_data_disk_attachment":                                   resourceArmVirtualMachineDataDiskAttachment(),
			"azurerm_virtual_machine_extension":                                              resourceArmVirtualMachineExtensions(),
			"azurerm_virtual_machine_scale_set":                                              resourceArmVirtualMachineScaleSet(),
			"azurerm_virtual_machine":                                                        resourceArmVirtualMachine(),
			"azurerm_virtual_network_gateway_connection":                                     resourceArmVirtualNetworkGatewayConnection(),
			"azurerm_virtual_network_gateway":                                                resourceArmVirtualNetworkGateway(),
			"azurerm_virtual_network_peering":                                                resourceArmVirtualNetworkPeering(),
			"azurerm_virtual_network":                                                        resourceArmVirtualNetwork(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		builder := &authentication.Builder{
			SubscriptionID:     d.Get("subscription_id").(string),
			ClientID:           d.Get("client_id").(string),
			ClientSecret:       d.Get("client_secret").(string),
			TenantID:           d.Get("tenant_id").(string),
			Environment:        d.Get("environment").(string),
			MsiEndpoint:        d.Get("msi_endpoint").(string),
			ClientCertPassword: d.Get("client_certificate_password").(string),
			ClientCertPath:     d.Get("client_certificate_path").(string),

			// Feature Toggles
			SupportsClientCertAuth:         true,
			SupportsClientSecretAuth:       true,
			SupportsManagedServiceIdentity: d.Get("use_msi").(bool),
			SupportsAzureCliToken:          true,
		}

		config, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("Error building AzureRM Client: %s", err)
		}

		partnerId := d.Get("partner_id").(string)
		skipProviderRegistration := d.Get("skip_provider_registration").(bool)
		client, err := getArmClient(config, skipProviderRegistration, partnerId)

		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error {
			client.StopContext = p.StopContext()
			return nil
		}

		skipCredentialsValidation := d.Get("skip_credentials_validation").(bool)
		if !skipCredentialsValidation {
			// List all the available providers and their registration state to avoid unnecessary
			// requests. This also lets us check if the provider credentials are correct.
			ctx := client.StopContext
			providerList, err := client.providersClient.List(ctx, nil, "")
			if err != nil {
				return nil, fmt.Errorf("Unable to list provider registration status, it is possible that this is due to invalid "+
					"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
					"error: %s", err)
			}

			if !skipProviderRegistration {
				availableResourceProviders := providerList.Values()
				requiredResourceProviders := requiredResourceProviders()

				err := ensureResourceProvidersAreRegistered(ctx, client.providersClient, availableResourceProviders, requiredResourceProviders)
				if err != nil {
					return nil, fmt.Errorf("Error ensuring Resource Providers are registered: %s", err)
				}
			}
		}

		return client, nil
	}
}

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = mutexkv.NewMutexKV()

// Deprecated: use `suppress.CaseDifference` instead
func ignoreCaseDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}

// ignoreCaseStateFunc is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func ignoreCaseStateFunc(val interface{}) string {
	return strings.ToLower(val.(string))
}

func userDataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return userDataStateFunc(old) == new
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = base64Encode(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

func base64EncodedStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		return base64Encode(s)
	default:
		return ""
	}
}

// base64Encode encodes data if the input isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func base64Encode(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if isBase64Encoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func isBase64Encoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}
