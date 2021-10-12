package alicloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/go-homedir"
)

// Provider returns a schema.Provider for alicloud
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"source_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SOURCE_IP", os.Getenv("ALICLOUD_SOURCE_IP")),
				Description: descriptions["source_ip"],
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ECS_ROLE_NAME", os.Getenv("ALICLOUD_ECS_ROLE_NAME")),
				Description: descriptions["ecs_role_name"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_REGION", os.Getenv("ALICLOUD_REGION")),
				Description: descriptions["region"],
			},
			"ots_instance_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'ots_instance_name' has been deprecated from provider version 1.10.0. New field 'instance_name' of resource 'alicloud_ots_table' instead.",
			},
			"log_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'log_endpoint' has been deprecated from provider version 1.28.0. New field 'log' which in nested endpoints instead.",
			},
			"mns_endpoint": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'mns_endpoint' has been deprecated from provider version 1.28.0. New field 'mns' which in nested endpoints instead.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ACCOUNT_ID", os.Getenv("ALICLOUD_ACCOUNT_ID")),
				Description: descriptions["account_id"],
			},
			"assume_role": assumeRoleSchema(),
			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_SHARED_CREDENTIALS_FILE", ""),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_PROFILE", ""),
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTPS",
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"client_read_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_READ_TIMEOUT", 30000),
				Description: descriptions["client_read_timeout"],
			},
			"client_connect_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLIENT_CONNECT_TIMEOUT", 30000),
				Description: descriptions["client_connect_timeout"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{

			"alicloud_account":                dataSourceAlicloudAccount(),
			"alicloud_caller_identity":        dataSourceAlicloudCallerIdentity(),
			"alicloud_images":                 dataSourceAlicloudImages(),
			"alicloud_regions":                dataSourceAlicloudRegions(),
			"alicloud_zones":                  dataSourceAlicloudZones(),
			"alicloud_db_zones":               dataSourceAlicloudDBZones(),
			"alicloud_instance_type_families": dataSourceAlicloudInstanceTypeFamilies(),
			"alicloud_instance_types":         dataSourceAlicloudInstanceTypes(),
			"alicloud_instances":              dataSourceAlicloudInstances(),
			"alicloud_disks":                  dataSourceAlicloudEcsDisks(),
			"alicloud_network_interfaces":     dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_snapshots":              dataSourceAlicloudEcsSnapshots(),
			"alicloud_vpcs":                   dataSourceAlicloudVpcs(),
			"alicloud_vswitches":              dataSourceAlicloudVswitches(),
			"alicloud_eips":                   dataSourceAlicloudEipAddresses(),
			"alicloud_key_pairs":              dataSourceAlicloudEcsKeyPairs(),
			"alicloud_kms_keys":               dataSourceAlicloudKmsKeys(),
			"alicloud_kms_ciphertext":         dataSourceAlicloudKmsCiphertext(),
			"alicloud_kms_plaintext":          dataSourceAlicloudKmsPlaintext(),
			"alicloud_dns_resolution_lines":   dataSourceAlicloudDnsResolutionLines(),
			"alicloud_dns_domains":            dataSourceAlicloudAlidnsDomains(),
			"alicloud_dns_groups":             dataSourceAlicloudDnsGroups(),
			"alicloud_dns_records":            dataSourceAlicloudDnsRecords(),
			// alicloud_dns_domain_groups, alicloud_dns_domain_records have been deprecated.
			"alicloud_dns_domain_groups":  dataSourceAlicloudDnsGroups(),
			"alicloud_dns_domain_records": dataSourceAlicloudDnsRecords(),
			// alicloud_ram_account_alias has been deprecated
			"alicloud_ram_account_alias":                           dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_aliases":                         dataSourceAlicloudRamAccountAlias(),
			"alicloud_ram_groups":                                  dataSourceAlicloudRamGroups(),
			"alicloud_ram_users":                                   dataSourceAlicloudRamUsers(),
			"alicloud_ram_roles":                                   dataSourceAlicloudRamRoles(),
			"alicloud_ram_policies":                                dataSourceAlicloudRamPolicies(),
			"alicloud_security_groups":                             dataSourceAlicloudSecurityGroups(),
			"alicloud_security_group_rules":                        dataSourceAlicloudSecurityGroupRules(),
			"alicloud_slbs":                                        dataSourceAlicloudSlbLoadBalancers(),
			"alicloud_slb_attachments":                             dataSourceAlicloudSlbAttachments(),
			"alicloud_slb_backend_servers":                         dataSourceAlicloudSlbBackendServers(),
			"alicloud_slb_listeners":                               dataSourceAlicloudSlbListeners(),
			"alicloud_slb_rules":                                   dataSourceAlicloudSlbRules(),
			"alicloud_slb_server_groups":                           dataSourceAlicloudSlbServerGroups(),
			"alicloud_slb_master_slave_server_groups":              dataSourceAlicloudSlbMasterSlaveServerGroups(),
			"alicloud_slb_acls":                                    dataSourceAlicloudSlbAcls(),
			"alicloud_slb_server_certificates":                     dataSourceAlicloudSlbServerCertificates(),
			"alicloud_slb_ca_certificates":                         dataSourceAlicloudSlbCaCertificates(),
			"alicloud_slb_domain_extensions":                       dataSourceAlicloudSlbDomainExtensions(),
			"alicloud_slb_zones":                                   dataSourceAlicloudSlbZones(),
			"alicloud_oss_service":                                 dataSourceAlicloudOssService(),
			"alicloud_oss_bucket_objects":                          dataSourceAlicloudOssBucketObjects(),
			"alicloud_oss_buckets":                                 dataSourceAlicloudOssBuckets(),
			"alicloud_ons_instances":                               dataSourceAlicloudOnsInstances(),
			"alicloud_ons_topics":                                  dataSourceAlicloudOnsTopics(),
			"alicloud_ons_groups":                                  dataSourceAlicloudOnsGroups(),
			"alicloud_alikafka_consumer_groups":                    dataSourceAlicloudAlikafkaConsumerGroups(),
			"alicloud_alikafka_instances":                          dataSourceAlicloudAlikafkaInstances(),
			"alicloud_alikafka_topics":                             dataSourceAlicloudAlikafkaTopics(),
			"alicloud_alikafka_sasl_users":                         dataSourceAlicloudAlikafkaSaslUsers(),
			"alicloud_alikafka_sasl_acls":                          dataSourceAlicloudAlikafkaSaslAcls(),
			"alicloud_fc_functions":                                dataSourceAlicloudFcFunctions(),
			"alicloud_file_crc64_checksum":                         dataSourceAlicloudFileCRC64Checksum(),
			"alicloud_fc_services":                                 dataSourceAlicloudFcServices(),
			"alicloud_fc_triggers":                                 dataSourceAlicloudFcTriggers(),
			"alicloud_fc_custom_domains":                           dataSourceAlicloudFcCustomDomains(),
			"alicloud_fc_zones":                                    dataSourceAlicloudFcZones(),
			"alicloud_db_instances":                                dataSourceAlicloudDBInstances(),
			"alicloud_db_instance_engines":                         dataSourceAlicloudDBInstanceEngines(),
			"alicloud_db_instance_classes":                         dataSourceAlicloudDBInstanceClasses(),
			"alicloud_pvtz_zones":                                  dataSourceAlicloudPvtzZones(),
			"alicloud_pvtz_zone_records":                           dataSourceAlicloudPvtzZoneRecords(),
			"alicloud_router_interfaces":                           dataSourceAlicloudRouterInterfaces(),
			"alicloud_vpn_gateways":                                dataSourceAlicloudVpnGateways(),
			"alicloud_vpn_customer_gateways":                       dataSourceAlicloudVpnCustomerGateways(),
			"alicloud_vpn_connections":                             dataSourceAlicloudVpnConnections(),
			"alicloud_ssl_vpn_servers":                             dataSourceAlicloudSslVpnServers(),
			"alicloud_ssl_vpn_client_certs":                        dataSourceAlicloudSslVpnClientCerts(),
			"alicloud_mongo_instances":                             dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_instances":                           dataSourceAlicloudMongoDBInstances(),
			"alicloud_mongodb_zones":                               dataSourceAlicloudMongoDBZones(),
			"alicloud_gpdb_instances":                              dataSourceAlicloudGpdbInstances(),
			"alicloud_gpdb_zones":                                  dataSourceAlicloudGpdbZones(),
			"alicloud_kvstore_instances":                           dataSourceAlicloudKvstoreInstances(),
			"alicloud_kvstore_zones":                               dataSourceAlicloudKVStoreZones(),
			"alicloud_kvstore_permission":                          dataSourceAlicloudKVStorePermission(),
			"alicloud_kvstore_instance_classes":                    dataSourceAlicloudKVStoreInstanceClasses(),
			"alicloud_kvstore_instance_engines":                    dataSourceAlicloudKVStoreInstanceEngines(),
			"alicloud_cen_instances":                               dataSourceAlicloudCenInstances(),
			"alicloud_cen_bandwidth_packages":                      dataSourceAlicloudCenBandwidthPackages(),
			"alicloud_cen_bandwidth_limits":                        dataSourceAlicloudCenBandwidthLimits(),
			"alicloud_cen_route_entries":                           dataSourceAlicloudCenRouteEntries(),
			"alicloud_cen_region_route_entries":                    dataSourceAlicloudCenRegionRouteEntries(),
			"alicloud_cen_transit_router_route_entries":            dataSourceAlicloudCenTransitRouterRouteEntries(),
			"alicloud_cen_transit_router_route_table_associations": dataSourceAlicloudCenTransitRouterRouteTableAssociations(),
			"alicloud_cen_transit_router_route_table_propagations": dataSourceAlicloudCenTransitRouterRouteTablePropagations(),
			"alicloud_cen_transit_router_route_tables":             dataSourceAlicloudCenTransitRouterRouteTables(),
			"alicloud_cen_transit_router_vbr_attachments":          dataSourceAlicloudCenTransitRouterVbrAttachments(),
			"alicloud_cen_transit_router_vpc_attachments":          dataSourceAlicloudCenTransitRouterVpcAttachments(),
			"alicloud_cen_transit_routers":                         dataSourceAlicloudCenTransitRouters(),
			"alicloud_cs_kubernetes_clusters":                      dataSourceAlicloudCSKubernetesClusters(),
			"alicloud_cs_managed_kubernetes_clusters":              dataSourceAlicloudCSManagerKubernetesClusters(),
			"alicloud_cs_edge_kubernetes_clusters":                 dataSourceAlicloudCSEdgeKubernetesClusters(),
			"alicloud_cs_serverless_kubernetes_clusters":           dataSourceAlicloudCSServerlessKubernetesClusters(),
			"alicloud_cs_kubernetes_permissions":                   dataSourceAlicloudCSKubernetesPermissions(),
			"alicloud_cr_namespaces":                               dataSourceAlicloudCRNamespaces(),
			"alicloud_cr_repos":                                    dataSourceAlicloudCRRepos(),
			"alicloud_cr_ee_instances":                             dataSourceAlicloudCrEEInstances(),
			"alicloud_cr_ee_namespaces":                            dataSourceAlicloudCrEENamespaces(),
			"alicloud_cr_ee_repos":                                 dataSourceAlicloudCrEERepos(),
			"alicloud_cr_ee_sync_rules":                            dataSourceAlicloudCrEESyncRules(),
			"alicloud_mns_queues":                                  dataSourceAlicloudMNSQueues(),
			"alicloud_mns_topics":                                  dataSourceAlicloudMNSTopics(),
			"alicloud_mns_topic_subscriptions":                     dataSourceAlicloudMNSTopicSubscriptions(),
			"alicloud_api_gateway_service":                         dataSourceAlicloudApiGatewayService(),
			"alicloud_api_gateway_apis":                            dataSourceAlicloudApiGatewayApis(),
			"alicloud_api_gateway_groups":                          dataSourceAlicloudApiGatewayGroups(),
			"alicloud_api_gateway_apps":                            dataSourceAlicloudApiGatewayApps(),
			"alicloud_elasticsearch_instances":                     dataSourceAlicloudElasticsearch(),
			"alicloud_elasticsearch_zones":                         dataSourceAlicloudElaticsearchZones(),
			"alicloud_drds_instances":                              dataSourceAlicloudDRDSInstances(),
			"alicloud_nas_service":                                 dataSourceAlicloudNasService(),
			"alicloud_nas_access_groups":                           dataSourceAlicloudNasAccessGroups(),
			"alicloud_nas_access_rules":                            dataSourceAlicloudAccessRules(),
			"alicloud_nas_mount_targets":                           dataSourceAlicloudNasMountTargets(),
			"alicloud_nas_file_systems":                            dataSourceAlicloudFileSystems(),
			"alicloud_nas_protocols":                               dataSourceAlicloudNasProtocols(),
			"alicloud_cas_certificates":                            dataSourceAlicloudSslCertificatesServiceCertificates(),
			"alicloud_common_bandwidth_packages":                   dataSourceAlicloudCommonBandwidthPackages(),
			"alicloud_route_tables":                                dataSourceAlicloudRouteTables(),
			"alicloud_route_entries":                               dataSourceAlicloudRouteEntries(),
			"alicloud_nat_gateways":                                dataSourceAlicloudNatGateways(),
			"alicloud_snat_entries":                                dataSourceAlicloudSnatEntries(),
			"alicloud_forward_entries":                             dataSourceAlicloudForwardEntries(),
			"alicloud_ddoscoo_instances":                           dataSourceAlicloudDdoscooInstances(),
			"alicloud_ddosbgp_instances":                           dataSourceAlicloudDdosbgpInstances(),
			"alicloud_ess_alarms":                                  dataSourceAlicloudEssAlarms(),
			"alicloud_ess_notifications":                           dataSourceAlicloudEssNotifications(),
			"alicloud_ess_scaling_groups":                          dataSourceAlicloudEssScalingGroups(),
			"alicloud_ess_scaling_rules":                           dataSourceAlicloudEssScalingRules(),
			"alicloud_ess_scaling_configurations":                  dataSourceAlicloudEssScalingConfigurations(),
			"alicloud_ess_lifecycle_hooks":                         dataSourceAlicloudEssLifecycleHooks(),
			"alicloud_ess_scheduled_tasks":                         dataSourceAlicloudEssScheduledTasks(),
			"alicloud_ots_service":                                 dataSourceAlicloudOtsService(),
			"alicloud_ots_instances":                               dataSourceAlicloudOtsInstances(),
			"alicloud_ots_instance_attachments":                    dataSourceAlicloudOtsInstanceAttachments(),
			"alicloud_ots_tables":                                  dataSourceAlicloudOtsTables(),
			"alicloud_cloud_connect_networks":                      dataSourceAlicloudCloudConnectNetworks(),
			"alicloud_emr_instance_types":                          dataSourceAlicloudEmrInstanceTypes(),
			"alicloud_emr_disk_types":                              dataSourceAlicloudEmrDiskTypes(),
			"alicloud_emr_main_versions":                           dataSourceAlicloudEmrMainVersions(),
			"alicloud_sag_acls":                                    dataSourceAlicloudSagAcls(),
			"alicloud_yundun_dbaudit_instance":                     dataSourceAlicloudDbauditInstances(),
			"alicloud_yundun_bastionhost_instances":                dataSourceAlicloudBastionhostInstances(),
			"alicloud_bastionhost_instances":                       dataSourceAlicloudBastionhostInstances(),
			"alicloud_market_product":                              dataSourceAlicloudProduct(),
			"alicloud_market_products":                             dataSourceAlicloudProducts(),
			"alicloud_polardb_clusters":                            dataSourceAlicloudPolarDBClusters(),
			"alicloud_polardb_node_classes":                        dataSourceAlicloudPolarDBNodeClasses(),
			"alicloud_polardb_endpoints":                           dataSourceAlicloudPolarDBEndpoints(),
			"alicloud_polardb_accounts":                            dataSourceAlicloudPolarDBAccounts(),
			"alicloud_polardb_databases":                           dataSourceAlicloudPolarDBDatabases(),
			"alicloud_polardb_zones":                               dataSourceAlicloudPolarDBZones(),
			"alicloud_hbase_instances":                             dataSourceAlicloudHBaseInstances(),
			"alicloud_hbase_zones":                                 dataSourceAlicloudHBaseZones(),
			"alicloud_hbase_instance_types":                        dataSourceAlicloudHBaseInstanceTypes(),
			"alicloud_adb_clusters":                                dataSourceAlicloudAdbDbClusters(),
			"alicloud_adb_zones":                                   dataSourceAlicloudAdbZones(),
			"alicloud_cen_flowlogs":                                dataSourceAlicloudCenFlowlogs(),
			"alicloud_kms_aliases":                                 dataSourceAlicloudKmsAliases(),
			"alicloud_dns_domain_txt_guid":                         dataSourceAlicloudDnsDomainTxtGuid(),
			"alicloud_edas_service":                                dataSourceAlicloudEdasService(),
			"alicloud_fnf_service":                                 dataSourceAlicloudFnfService(),
			"alicloud_kms_service":                                 dataSourceAlicloudKmsService(),
			"alicloud_sae_service":                                 dataSourceAlicloudSaeService(),
			"alicloud_dataworks_service":                           dataSourceAlicloudDataWorksService(),
			"alicloud_mns_service":                                 dataSourceAlicloudMnsService(),
			"alicloud_cloud_storage_gateway_service":               dataSourceAlicloudCloudStorageGatewayService(),
			"alicloud_vs_service":                                  dataSourceAlicloudVsService(),
			"alicloud_pvtz_service":                                dataSourceAlicloudPvtzService(),
			"alicloud_cms_service":                                 dataSourceAlicloudCmsService(),
			"alicloud_maxcompute_service":                          dataSourceAlicloudMaxcomputeService(),
			"alicloud_brain_industrial_service":                    dataSourceAlicloudBrainIndustrialService(),
			"alicloud_iot_service":                                 dataSourceAlicloudIotService(),
			"alicloud_ack_service":                                 dataSourceAlicloudAckService(),
			"alicloud_cr_service":                                  dataSourceAlicloudCrService(),
			"alicloud_dcdn_service":                                dataSourceAlicloudDcdnService(),
			"alicloud_datahub_service":                             dataSourceAlicloudDatahubService(),
			"alicloud_ons_service":                                 dataSourceAlicloudOnsService(),
			"alicloud_fc_service":                                  dataSourceAlicloudFcService(),
			"alicloud_privatelink_service":                         dataSourceAlicloudPrivateLinkService(),
			"alicloud_edas_applications":                           dataSourceAlicloudEdasApplications(),
			"alicloud_edas_deploy_groups":                          dataSourceAlicloudEdasDeployGroups(),
			"alicloud_edas_clusters":                               dataSourceAlicloudEdasClusters(),
			"alicloud_resource_manager_folders":                    dataSourceAlicloudResourceManagerFolders(),
			"alicloud_dns_instances":                               dataSourceAlicloudAlidnsInstances(),
			"alicloud_resource_manager_policies":                   dataSourceAlicloudResourceManagerPolicies(),
			"alicloud_resource_manager_resource_groups":            dataSourceAlicloudResourceManagerResourceGroups(),
			"alicloud_resource_manager_roles":                      dataSourceAlicloudResourceManagerRoles(),
			"alicloud_resource_manager_policy_versions":            dataSourceAlicloudResourceManagerPolicyVersions(),
			"alicloud_alidns_domain_groups":                        dataSourceAlicloudAlidnsDomainGroups(),
			"alicloud_kms_key_versions":                            dataSourceAlicloudKmsKeyVersions(),
			"alicloud_alidns_records":                              dataSourceAlicloudAlidnsRecords(),
			"alicloud_resource_manager_accounts":                   dataSourceAlicloudResourceManagerAccounts(),
			"alicloud_resource_manager_resource_directories":       dataSourceAlicloudResourceManagerResourceDirectories(),
			"alicloud_resource_manager_handshakes":                 dataSourceAlicloudResourceManagerHandshakes(),
			"alicloud_waf_domains":                                 dataSourceAlicloudWafDomains(),
			"alicloud_kms_secrets":                                 dataSourceAlicloudKmsSecrets(),
			"alicloud_cen_route_maps":                              dataSourceAlicloudCenRouteMaps(),
			"alicloud_cen_private_zones":                           dataSourceAlicloudCenPrivateZones(),
			"alicloud_dms_enterprise_instances":                    dataSourceAlicloudDmsEnterpriseInstances(),
			"alicloud_cassandra_clusters":                          dataSourceAlicloudCassandraClusters(),
			"alicloud_cassandra_data_centers":                      dataSourceAlicloudCassandraDataCenters(),
			"alicloud_cassandra_zones":                             dataSourceAlicloudCassandraZones(),
			"alicloud_kms_secret_versions":                         dataSourceAlicloudKmsSecretVersions(),
			"alicloud_waf_instances":                               dataSourceAlicloudWafInstances(),
			"alicloud_eci_image_caches":                            dataSourceAlicloudEciImageCaches(),
			"alicloud_dms_enterprise_users":                        dataSourceAlicloudDmsEnterpriseUsers(),
			"alicloud_ecs_dedicated_hosts":                         dataSourceAlicloudEcsDedicatedHosts(),
			"alicloud_oos_templates":                               dataSourceAlicloudOosTemplates(),
			"alicloud_oos_executions":                              dataSourceAlicloudOosExecutions(),
			"alicloud_resource_manager_policy_attachments":         dataSourceAlicloudResourceManagerPolicyAttachments(),
			"alicloud_dcdn_domains":                                dataSourceAlicloudDcdnDomains(),
			"alicloud_mse_clusters":                                dataSourceAlicloudMseClusters(),
			"alicloud_actiontrail_trails":                          dataSourceAlicloudActiontrailTrails(),
			"alicloud_actiontrails":                                dataSourceAlicloudActiontrailTrails(),
			"alicloud_alidns_instances":                            dataSourceAlicloudAlidnsInstances(),
			"alicloud_alidns_domains":                              dataSourceAlicloudAlidnsDomains(),
			"alicloud_log_service":                                 dataSourceAlicloudLogService(),
			"alicloud_cen_instance_attachments":                    dataSourceAlicloudCenInstanceAttachments(),
			"alicloud_cdn_service":                                 dataSourceAlicloudCdnService(),
			"alicloud_cen_vbr_health_checks":                       dataSourceAlicloudCenVbrHealthChecks(),
			"alicloud_config_rules":                                dataSourceAlicloudConfigRules(),
			"alicloud_config_configuration_recorders":              dataSourceAlicloudConfigConfigurationRecorders(),
			"alicloud_config_delivery_channels":                    dataSourceAlicloudConfigDeliveryChannels(),
			"alicloud_cms_alarm_contacts":                          dataSourceAlicloudCmsAlarmContacts(),
			"alicloud_kvstore_connections":                         dataSourceAlicloudKvstoreConnections(),
			"alicloud_cms_alarm_contact_groups":                    dataSourceAlicloudCmsAlarmContactGroups(),
			"alicloud_enhanced_nat_available_zones":                dataSourceAlicloudEnhancedNatAvailableZones(),
			"alicloud_cen_route_services":                          dataSourceAlicloudCenRouteServices(),
			"alicloud_kvstore_accounts":                            dataSourceAlicloudKvstoreAccounts(),
			"alicloud_cms_group_metric_rules":                      dataSourceAlicloudCmsGroupMetricRules(),
			"alicloud_fnf_flows":                                   dataSourceAlicloudFnfFlows(),
			"alicloud_fnf_schedules":                               dataSourceAlicloudFnfSchedules(),
			"alicloud_ros_change_sets":                             dataSourceAlicloudRosChangeSets(),
			"alicloud_ros_stacks":                                  dataSourceAlicloudRosStacks(),
			"alicloud_ros_stack_groups":                            dataSourceAlicloudRosStackGroups(),
			"alicloud_ros_templates":                               dataSourceAlicloudRosTemplates(),
			"alicloud_privatelink_vpc_endpoint_services":           dataSourceAlicloudPrivatelinkVpcEndpointServices(),
			"alicloud_privatelink_vpc_endpoints":                   dataSourceAlicloudPrivatelinkVpcEndpoints(),
			"alicloud_privatelink_vpc_endpoint_connections":        dataSourceAlicloudPrivatelinkVpcEndpointConnections(),
			"alicloud_privatelink_vpc_endpoint_service_resources":  dataSourceAlicloudPrivatelinkVpcEndpointServiceResources(),
			"alicloud_privatelink_vpc_endpoint_service_users":      dataSourceAlicloudPrivatelinkVpcEndpointServiceUsers(),
			"alicloud_resource_manager_resource_shares":            dataSourceAlicloudResourceManagerResourceShares(),
			"alicloud_privatelink_vpc_endpoint_zones":              dataSourceAlicloudPrivatelinkVpcEndpointZones(),
			"alicloud_ga_accelerators":                             dataSourceAlicloudGaAccelerators(),
			"alicloud_eci_container_groups":                        dataSourceAlicloudEciContainerGroups(),
			"alicloud_resource_manager_shared_resources":           dataSourceAlicloudResourceManagerSharedResources(),
			"alicloud_resource_manager_shared_targets":             dataSourceAlicloudResourceManagerSharedTargets(),
			"alicloud_ga_listeners":                                dataSourceAlicloudGaListeners(),
			"alicloud_tsdb_instances":                              dataSourceAlicloudTsdbInstances(),
			"alicloud_tsdb_zones":                                  dataSourceAlicloudTsdbZones(),
			"alicloud_ga_bandwidth_packages":                       dataSourceAlicloudGaBandwidthPackages(),
			"alicloud_ga_endpoint_groups":                          dataSourceAlicloudGaEndpointGroups(),
			"alicloud_brain_industrial_pid_organizations":          dataSourceAlicloudBrainIndustrialPidOrganizations(),
			"alicloud_ga_ip_sets":                                  dataSourceAlicloudGaIpSets(),
			"alicloud_ga_forwarding_rules":                         dataSourceAlicloudGaForwardingRules(),
			"alicloud_eipanycast_anycast_eip_addresses":            dataSourceAlicloudEipanycastAnycastEipAddresses(),
			"alicloud_brain_industrial_pid_projects":               dataSourceAlicloudBrainIndustrialPidProjects(),
			"alicloud_cms_monitor_groups":                          dataSourceAlicloudCmsMonitorGroups(),
			"alicloud_ram_saml_providers":                          dataSourceAlicloudRamSamlProviders(),
			"alicloud_quotas_quotas":                               dataSourceAlicloudQuotasQuotas(),
			"alicloud_quotas_application_infos":                    dataSourceAlicloudQuotasQuotaApplications(),
			"alicloud_cms_monitor_group_instanceses":               dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_cms_monitor_group_instances":                 dataSourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_quotas_quota_alarms":                         dataSourceAlicloudQuotasQuotaAlarms(),
			"alicloud_ecs_commands":                                dataSourceAlicloudEcsCommands(),
			"alicloud_cloud_storage_gateway_storage_bundles":       dataSourceAlicloudCloudStorageGatewayStorageBundles(),
			"alicloud_ecs_hpc_clusters":                            dataSourceAlicloudEcsHpcClusters(),
			"alicloud_brain_industrial_pid_loops":                  dataSourceAlicloudBrainIndustrialPidLoops(),
			"alicloud_quotas_quota_applications":                   dataSourceAlicloudQuotasQuotaApplications(),
			"alicloud_ecs_auto_snapshot_policies":                  dataSourceAlicloudEcsAutoSnapshotPolicies(),
			"alicloud_rds_parameter_groups":                        dataSourceAlicloudRdsParameterGroups(),
			"alicloud_ecs_launch_templates":                        dataSourceAlicloudEcsLaunchTemplates(),
			"alicloud_resource_manager_control_policies":           dataSourceAlicloudResourceManagerControlPolicies(),
			"alicloud_resource_manager_control_policy_attachments": dataSourceAlicloudResourceManagerControlPolicyAttachments(),
			"alicloud_rds_accounts":                                dataSourceAlicloudRdsAccounts(),
			"alicloud_havips":                                      dataSourceAlicloudHavips(),
			"alicloud_ecs_snapshots":                               dataSourceAlicloudEcsSnapshots(),
			"alicloud_ecs_key_pairs":                               dataSourceAlicloudEcsKeyPairs(),
			"alicloud_adb_db_clusters":                             dataSourceAlicloudAdbDbClusters(),
			"alicloud_vpc_flow_logs":                               dataSourceAlicloudVpcFlowLogs(),
			"alicloud_network_acls":                                dataSourceAlicloudNetworkAcls(),
			"alicloud_ecs_disks":                                   dataSourceAlicloudEcsDisks(),
			"alicloud_ddoscoo_domain_resources":                    dataSourceAlicloudDdoscooDomainResources(),
			"alicloud_ddoscoo_ports":                               dataSourceAlicloudDdoscooPorts(),
			"alicloud_slb_load_balancers":                          dataSourceAlicloudSlbLoadBalancers(),
			"alicloud_ecs_network_interfaces":                      dataSourceAlicloudEcsNetworkInterfaces(),
			"alicloud_config_aggregators":                          dataSourceAlicloudConfigAggregators(),
			"alicloud_config_aggregate_config_rules":               dataSourceAlicloudConfigAggregateConfigRules(),
			"alicloud_config_aggregate_compliance_packs":           dataSourceAlicloudConfigAggregateCompliancePacks(),
			"alicloud_config_compliance_packs":                     dataSourceAlicloudConfigCompliancePacks(),
			"alicloud_eip_addresses":                               dataSourceAlicloudEipAddresses(),
			"alicloud_direct_mail_receiverses":                     dataSourceAlicloudDirectMailReceiverses(),
			"alicloud_log_projects":                                dataSourceAlicloudLogProjects(),
			"alicloud_log_stores":                                  dataSourceAlicloudLogStores(),
			"alicloud_event_bridge_service":                        dataSourceAlicloudEventBridgeService(),
			"alicloud_event_bridge_event_buses":                    dataSourceAlicloudEventBridgeEventBuses(),
			"alicloud_amqp_virtual_hosts":                          dataSourceAlicloudAmqpVirtualHosts(),
			"alicloud_amqp_queues":                                 dataSourceAlicloudAmqpQueues(),
			"alicloud_amqp_exchanges":                              dataSourceAlicloudAmqpExchanges(),
			"alicloud_cassandra_backup_plans":                      dataSourceAlicloudCassandraBackupPlans(),
			"alicloud_cen_transit_router_peer_attachments":         dataSourceAlicloudCenTransitRouterPeerAttachments(),
			"alicloud_amqp_instances":                              dataSourceAlicloudAmqpInstances(),
			"alicloud_hbr_vaults":                                  dataSourceAlicloudHbrVaults(),
			"alicloud_ssl_certificates_service_certificates":       dataSourceAlicloudSslCertificatesServiceCertificates(),
			"alicloud_arms_alert_contacts":                         dataSourceAlicloudArmsAlertContacts(),
			"alicloud_event_bridge_rules":                          dataSourceAlicloudEventBridgeRules(),
			"alicloud_cloud_firewall_control_policies":             dataSourceAlicloudCloudFirewallControlPolicies(),
			"alicloud_sae_namespaces":                              dataSourceAlicloudSaeNamespaces(),
			"alicloud_sae_config_maps":                             dataSourceAlicloudSaeConfigMaps(),
			"alicloud_alb_security_policies":                       dataSourceAlicloudAlbSecurityPolicies(),
			"alicloud_event_bridge_event_sources":                  dataSourceAlicloudEventBridgeEventSources(),
			"alicloud_ecd_policy_groups":                           dataSourceAlicloudEcdPolicyGroups(),
			"alicloud_ecp_key_pairs":                               dataSourceAlicloudEcpKeyPairs(),
			"alicloud_hbr_ecs_backup_plans":                        dataSourceAlicloudHbrEcsBackupPlans(),
			"alicloud_hbr_nas_backup_plans":                        dataSourceAlicloudHbrNasBackupPlans(),
			"alicloud_hbr_oss_backup_plans":                        dataSourceAlicloudHbrOssBackupPlans(),
			"alicloud_scdn_domains":                                dataSourceAlicloudScdnDomains(),
			"alicloud_alb_server_groups":                           dataSourceAlicloudAlbServerGroups(),
			"alicloud_data_works_folders":                          dataSourceAlicloudDataWorksFolders(),
			"alicloud_arms_alert_contact_groups":                   dataSourceAlicloudArmsAlertContactGroups(),
			"alicloud_express_connect_access_points":               dataSourceAlicloudExpressConnectAccessPoints(),
			"alicloud_cloud_storage_gateway_gateways":              dataSourceAlicloudCloudStorageGatewayGateways(),
			"alicloud_lindorm_instances":                           dataSourceAlicloudLindormInstances(),
			"alicloud_express_connect_physical_connection_service": dataSourceAlicloudExpressConnectPhysicalConnectionService(),
			"alicloud_cddc_dedicated_host_groups":                  dataSourceAlicloudCddcDedicatedHostGroups(),
			"alicloud_hbr_ecs_backup_clients":                      dataSourceAlicloudHbrEcsBackupClients(),
			"alicloud_msc_sub_contacts":                            dataSourceAlicloudMscSubContacts(),
			"alicloud_express_connect_physical_connections":        dataSourceAlicloudExpressConnectPhysicalConnections(),
			"alicloud_alb_load_balancers":                          dataSourceAlicloudAlbLoadBalancers(),
			"alicloud_alb_zones":                                   dataSourceAlicloudAlbZones(),
			"alicloud_sddp_rules":                                  dataSourceAlicloudSddpRules(),
			"alicloud_bastionhost_user_groups":                     dataSourceAlicloudBastionhostUserGroups(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"alicloud_instance":                           resourceAliyunInstance(),
			"alicloud_image":                              resourceAliCloudImage(),
			"alicloud_reserved_instance":                  resourceAliCloudReservedInstance(),
			"alicloud_copy_image":                         resourceAliCloudImageCopy(),
			"alicloud_image_export":                       resourceAliCloudImageExport(),
			"alicloud_image_copy":                         resourceAliCloudImageCopy(),
			"alicloud_image_import":                       resourceAliCloudImageImport(),
			"alicloud_image_share_permission":             resourceAliCloudImageSharePermission(),
			"alicloud_ram_role_attachment":                resourceAlicloudRamRoleAttachment(),
			"alicloud_disk":                               resourceAlicloudEcsDisk(),
			"alicloud_disk_attachment":                    resourceAlicloudEcsDiskAttachment(),
			"alicloud_network_interface":                  resourceAlicloudEcsNetworkInterface(),
			"alicloud_network_interface_attachment":       resourceAlicloudEcsNetworkInterfaceAttachment(),
			"alicloud_snapshot":                           resourceAlicloudEcsSnapshot(),
			"alicloud_snapshot_policy":                    resourceAlicloudEcsAutoSnapshotPolicy(),
			"alicloud_launch_template":                    resourceAlicloudEcsLaunchTemplate(),
			"alicloud_security_group":                     resourceAliyunSecurityGroup(),
			"alicloud_security_group_rule":                resourceAliyunSecurityGroupRule(),
			"alicloud_db_database":                        resourceAlicloudDBDatabase(),
			"alicloud_db_account":                         resourceAlicloudRdsAccount(),
			"alicloud_db_account_privilege":               resourceAlicloudDBAccountPrivilege(),
			"alicloud_db_backup_policy":                   resourceAlicloudDBBackupPolicy(),
			"alicloud_db_connection":                      resourceAlicloudDBConnection(),
			"alicloud_db_read_write_splitting_connection": resourceAlicloudDBReadWriteSplittingConnection(),
			"alicloud_db_instance":                        resourceAlicloudDBInstance(),
			"alicloud_mongodb_instance":                   resourceAlicloudMongoDBInstance(),
			"alicloud_mongodb_sharding_instance":          resourceAlicloudMongoDBShardingInstance(),
			"alicloud_gpdb_instance":                      resourceAlicloudGpdbInstance(),
			"alicloud_gpdb_elastic_instance":              resourceAlicloudGpdbElasticInstance(),
			"alicloud_gpdb_connection":                    resourceAlicloudGpdbConnection(),
			"alicloud_db_readonly_instance":               resourceAlicloudDBReadonlyInstance(),
			"alicloud_auto_provisioning_group":            resourceAlicloudAutoProvisioningGroup(),
			"alicloud_ess_scaling_group":                  resourceAlicloudEssScalingGroup(),
			"alicloud_ess_scaling_configuration":          resourceAlicloudEssScalingConfiguration(),
			"alicloud_ess_scaling_rule":                   resourceAlicloudEssScalingRule(),
			"alicloud_ess_schedule":                       resourceAlicloudEssScheduledTask(),
			"alicloud_ess_scheduled_task":                 resourceAlicloudEssScheduledTask(),
			"alicloud_ess_attachment":                     resourceAlicloudEssAttachment(),
			"alicloud_ess_lifecycle_hook":                 resourceAlicloudEssLifecycleHook(),
			"alicloud_ess_notification":                   resourceAlicloudEssNotification(),
			"alicloud_ess_alarm":                          resourceAlicloudEssAlarm(),
			"alicloud_ess_scalinggroup_vserver_groups":    resourceAlicloudEssScalingGroupVserverGroups(),
			"alicloud_vpc":                                resourceAlicloudVpc(),
			"alicloud_nat_gateway":                        resourceAlicloudNatGateway(),
			"alicloud_nas_file_system":                    resourceAlicloudNasFileSystem(),
			"alicloud_nas_mount_target":                   resourceAlicloudNasMountTarget(),
			"alicloud_nas_access_group":                   resourceAlicloudNasAccessGroup(),
			"alicloud_nas_access_rule":                    resourceAlicloudNasAccessRule(),
			// "alicloud_subnet" aims to match aws usage habit.
			"alicloud_subnet":                        resourceAlicloudVswitch(),
			"alicloud_vswitch":                       resourceAlicloudVswitch(),
			"alicloud_route_entry":                   resourceAliyunRouteEntry(),
			"alicloud_route_table":                   resourceAlicloudRouteTable(),
			"alicloud_route_table_attachment":        resourceAliyunRouteTableAttachment(),
			"alicloud_snat_entry":                    resourceAlicloudSnatEntry(),
			"alicloud_forward_entry":                 resourceAlicloudForwardEntry(),
			"alicloud_eip":                           resourceAlicloudEipAddress(),
			"alicloud_eip_association":               resourceAliyunEipAssociation(),
			"alicloud_slb":                           resourceAlicloudSlbLoadBalancer(),
			"alicloud_slb_listener":                  resourceAliyunSlbListener(),
			"alicloud_slb_attachment":                resourceAliyunSlbAttachment(),
			"alicloud_slb_backend_server":            resourceAliyunSlbBackendServer(),
			"alicloud_slb_domain_extension":          resourceAlicloudSlbDomainExtension(),
			"alicloud_slb_server_group":              resourceAliyunSlbServerGroup(),
			"alicloud_slb_master_slave_server_group": resourceAliyunSlbMasterSlaveServerGroup(),
			"alicloud_slb_rule":                      resourceAliyunSlbRule(),
			"alicloud_slb_acl":                       resourceAlicloudSlbAcl(),
			"alicloud_slb_ca_certificate":            resourceAlicloudSlbCaCertificate(),
			"alicloud_slb_server_certificate":        resourceAlicloudSlbServerCertificate(),
			"alicloud_oss_bucket":                    resourceAlicloudOssBucket(),
			"alicloud_oss_bucket_object":             resourceAlicloudOssBucketObject(),
			"alicloud_ons_instance":                  resourceAlicloudOnsInstance(),
			"alicloud_ons_topic":                     resourceAlicloudOnsTopic(),
			"alicloud_ons_group":                     resourceAlicloudOnsGroup(),
			"alicloud_alikafka_consumer_group":       resourceAlicloudAlikafkaConsumerGroup(),
			"alicloud_alikafka_instance":             resourceAlicloudAlikafkaInstance(),
			"alicloud_alikafka_topic":                resourceAlicloudAlikafkaTopic(),
			"alicloud_alikafka_sasl_user":            resourceAlicloudAlikafkaSaslUser(),
			"alicloud_alikafka_sasl_acl":             resourceAlicloudAlikafkaSaslAcl(),
			"alicloud_dns_record":                    resourceAlicloudDnsRecord(),
			"alicloud_dns":                           resourceAlicloudDns(),
			"alicloud_dns_group":                     resourceAlicloudDnsGroup(),
			"alicloud_key_pair":                      resourceAlicloudEcsKeyPair(),
			"alicloud_key_pair_attachment":           resourceAlicloudEcsKeyPairAttachment(),
			"alicloud_kms_key":                       resourceAlicloudKmsKey(),
			"alicloud_kms_ciphertext":                resourceAlicloudKmsCiphertext(),
			"alicloud_ram_user":                      resourceAlicloudRamUser(),
			"alicloud_ram_account_password_policy":   resourceAlicloudRamAccountPasswordPolicy(),
			"alicloud_ram_access_key":                resourceAlicloudRamAccessKey(),
			"alicloud_ram_login_profile":             resourceAlicloudRamLoginProfile(),
			"alicloud_ram_group":                     resourceAlicloudRamGroup(),
			"alicloud_ram_role":                      resourceAlicloudRamRole(),
			"alicloud_ram_policy":                    resourceAlicloudRamPolicy(),
			// alicloud_ram_alias has been deprecated
			"alicloud_ram_alias":                                  resourceAlicloudRamAccountAlias(),
			"alicloud_ram_account_alias":                          resourceAlicloudRamAccountAlias(),
			"alicloud_ram_group_membership":                       resourceAlicloudRamGroupMembership(),
			"alicloud_ram_user_policy_attachment":                 resourceAlicloudRamUserPolicyAtatchment(),
			"alicloud_ram_role_policy_attachment":                 resourceAlicloudRamRolePolicyAttachment(),
			"alicloud_ram_group_policy_attachment":                resourceAlicloudRamGroupPolicyAtatchment(),
			"alicloud_container_cluster":                          resourceAlicloudCSSwarm(),
			"alicloud_cs_application":                             resourceAlicloudCSApplication(),
			"alicloud_cs_swarm":                                   resourceAlicloudCSSwarm(),
			"alicloud_cs_kubernetes":                              resourceAlicloudCSKubernetes(),
			"alicloud_cs_managed_kubernetes":                      resourceAlicloudCSManagedKubernetes(),
			"alicloud_cs_edge_kubernetes":                         resourceAlicloudCSEdgeKubernetes(),
			"alicloud_cs_serverless_kubernetes":                   resourceAlicloudCSServerlessKubernetes(),
			"alicloud_cs_kubernetes_autoscaler":                   resourceAlicloudCSKubernetesAutoscaler(),
			"alicloud_cs_kubernetes_node_pool":                    resourceAlicloudCSKubernetesNodePool(),
			"alicloud_cs_kubernetes_permissions":                  resourceAlicloudCSKubernetesPermissions(),
			"alicloud_cs_autoscaling_config":                      resourceAlicloudCSAutoscalingConfig(),
			"alicloud_cr_namespace":                               resourceAlicloudCRNamespace(),
			"alicloud_cr_repo":                                    resourceAlicloudCRRepo(),
			"alicloud_cr_ee_instance":                             resourceAlicloudCrEEInstance(),
			"alicloud_cr_ee_namespace":                            resourceAlicloudCrEENamespace(),
			"alicloud_cr_ee_repo":                                 resourceAlicloudCrEERepo(),
			"alicloud_cr_ee_sync_rule":                            resourceAlicloudCrEESyncRule(),
			"alicloud_cdn_domain":                                 resourceAlicloudCdnDomain(),
			"alicloud_cdn_domain_new":                             resourceAlicloudCdnDomainNew(),
			"alicloud_cdn_domain_config":                          resourceAlicloudCdnDomainConfig(),
			"alicloud_router_interface":                           resourceAlicloudRouterInterface(),
			"alicloud_router_interface_connection":                resourceAlicloudRouterInterfaceConnection(),
			"alicloud_ots_table":                                  resourceAlicloudOtsTable(),
			"alicloud_ots_instance":                               resourceAlicloudOtsInstance(),
			"alicloud_ots_instance_attachment":                    resourceAlicloudOtsInstanceAttachment(),
			"alicloud_cms_alarm":                                  resourceAlicloudCmsAlarm(),
			"alicloud_cms_site_monitor":                           resourceAlicloudCmsSiteMonitor(),
			"alicloud_pvtz_zone":                                  resourceAlicloudPvtzZone(),
			"alicloud_pvtz_zone_attachment":                       resourceAlicloudPvtzZoneAttachment(),
			"alicloud_pvtz_zone_record":                           resourceAlicloudPvtzZoneRecord(),
			"alicloud_log_alert":                                  resourceAlicloudLogAlert(),
			"alicloud_log_audit":                                  resourceAlicloudLogAudit(),
			"alicloud_log_dashboard":                              resourceAlicloudLogDashboard(),
			"alicloud_log_etl":                                    resourceAlicloudLogETL(),
			"alicloud_log_machine_group":                          resourceAlicloudLogMachineGroup(),
			"alicloud_log_oss_shipper":                            resourceAlicloudLogOssShipper(),
			"alicloud_log_project":                                resourceAlicloudLogProject(),
			"alicloud_log_store":                                  resourceAlicloudLogStore(),
			"alicloud_log_store_index":                            resourceAlicloudLogStoreIndex(),
			"alicloud_logtail_config":                             resourceAlicloudLogtailConfig(),
			"alicloud_logtail_attachment":                         resourceAlicloudLogtailAttachment(),
			"alicloud_fc_service":                                 resourceAlicloudFCService(),
			"alicloud_fc_function":                                resourceAlicloudFCFunction(),
			"alicloud_fc_trigger":                                 resourceAlicloudFCTrigger(),
			"alicloud_fc_alias":                                   resourceAlicloudFCAlias(),
			"alicloud_fc_custom_domain":                           resourceAlicloudFCCustomDomain(),
			"alicloud_fc_function_async_invoke_config":            resourceAlicloudFCFunctionAsyncInvokeConfig(),
			"alicloud_vpn_gateway":                                resourceAliyunVpnGateway(),
			"alicloud_vpn_customer_gateway":                       resourceAliyunVpnCustomerGateway(),
			"alicloud_vpn_route_entry":                            resourceAliyunVpnRouteEntry(),
			"alicloud_vpn_connection":                             resourceAliyunVpnConnection(),
			"alicloud_ssl_vpn_server":                             resourceAliyunSslVpnServer(),
			"alicloud_ssl_vpn_client_cert":                        resourceAliyunSslVpnClientCert(),
			"alicloud_cen_instance":                               resourceAlicloudCenInstance(),
			"alicloud_cen_instance_attachment":                    resourceAlicloudCenInstanceAttachment(),
			"alicloud_cen_bandwidth_package":                      resourceAlicloudCenBandwidthPackage(),
			"alicloud_cen_bandwidth_package_attachment":           resourceAlicloudCenBandwidthPackageAttachment(),
			"alicloud_cen_bandwidth_limit":                        resourceAlicloudCenBandwidthLimit(),
			"alicloud_cen_route_entry":                            resourceAlicloudCenRouteEntry(),
			"alicloud_cen_instance_grant":                         resourceAlicloudCenInstanceGrant(),
			"alicloud_cen_transit_router":                         resourceAlicloudCenTransitRouter(),
			"alicloud_cen_transit_router_route_entry":             resourceAlicloudCenTransitRouterRouteEntry(),
			"alicloud_cen_transit_router_route_table":             resourceAlicloudCenTransitRouterRouteTable(),
			"alicloud_cen_transit_router_route_table_association": resourceAlicloudCenTransitRouterRouteTableAssociation(),
			"alicloud_cen_transit_router_route_table_propagation": resourceAlicloudCenTransitRouterRouteTablePropagation(),
			"alicloud_cen_transit_router_vbr_attachment":          resourceAlicloudCenTransitRouterVbrAttachment(),
			"alicloud_cen_transit_router_vpc_attachment":          resourceAlicloudCenTransitRouterVpcAttachment(),
			"alicloud_kvstore_instance":                           resourceAlicloudKvstoreInstance(),
			"alicloud_kvstore_backup_policy":                      resourceAlicloudKVStoreBackupPolicy(),
			"alicloud_kvstore_account":                            resourceAlicloudKvstoreAccount(),
			"alicloud_datahub_project":                            resourceAlicloudDatahubProject(),
			"alicloud_datahub_subscription":                       resourceAlicloudDatahubSubscription(),
			"alicloud_datahub_topic":                              resourceAlicloudDatahubTopic(),
			"alicloud_mns_queue":                                  resourceAlicloudMNSQueue(),
			"alicloud_mns_topic":                                  resourceAlicloudMNSTopic(),
			"alicloud_havip":                                      resourceAlicloudHavip(),
			"alicloud_mns_topic_subscription":                     resourceAlicloudMNSSubscription(),
			"alicloud_havip_attachment":                           resourceAliyunHaVipAttachment(),
			"alicloud_api_gateway_api":                            resourceAliyunApigatewayApi(),
			"alicloud_api_gateway_group":                          resourceAliyunApigatewayGroup(),
			"alicloud_api_gateway_app":                            resourceAliyunApigatewayApp(),
			"alicloud_api_gateway_app_attachment":                 resourceAliyunApigatewayAppAttachment(),
			"alicloud_api_gateway_vpc_access":                     resourceAliyunApigatewayVpc(),
			"alicloud_common_bandwidth_package":                   resourceAlicloudCommonBandwidthPackage(),
			"alicloud_common_bandwidth_package_attachment":        resourceAliyunCommonBandwidthPackageAttachment(),
			"alicloud_drds_instance":                              resourceAlicloudDRDSInstance(),
			"alicloud_elasticsearch_instance":                     resourceAlicloudElasticsearch(),
			"alicloud_cas_certificate":                            resourceAlicloudSslCertificatesServiceCertificate(),
			"alicloud_ddoscoo_instance":                           resourceAlicloudDdoscooInstance(),
			"alicloud_ddosbgp_instance":                           resourceAlicloudDdosbgpInstance(),
			"alicloud_network_acl":                                resourceAlicloudNetworkAcl(),
			"alicloud_network_acl_attachment":                     resourceAliyunNetworkAclAttachment(),
			"alicloud_network_acl_entries":                        resourceAliyunNetworkAclEntries(),
			"alicloud_emr_cluster":                                resourceAlicloudEmrCluster(),
			"alicloud_cloud_connect_network":                      resourceAlicloudCloudConnectNetwork(),
			"alicloud_cloud_connect_network_attachment":           resourceAlicloudCloudConnectNetworkAttachment(),
			"alicloud_cloud_connect_network_grant":                resourceAlicloudCloudConnectNetworkGrant(),
			"alicloud_sag_acl":                                    resourceAlicloudSagAcl(),
			"alicloud_sag_acl_rule":                               resourceAlicloudSagAclRule(),
			"alicloud_sag_qos":                                    resourceAlicloudSagQos(),
			"alicloud_sag_qos_policy":                             resourceAlicloudSagQosPolicy(),
			"alicloud_sag_qos_car":                                resourceAlicloudSagQosCar(),
			"alicloud_sag_snat_entry":                             resourceAlicloudSagSnatEntry(),
			"alicloud_sag_dnat_entry":                             resourceAlicloudSagDnatEntry(),
			"alicloud_sag_client_user":                            resourceAlicloudSagClientUser(),
			"alicloud_yundun_dbaudit_instance":                    resourceAlicloudDbauditInstance(),
			"alicloud_yundun_bastionhost_instance":                resourceAlicloudBastionhostInstance(),
			"alicloud_bastionhost_instance":                       resourceAlicloudBastionhostInstance(),
			"alicloud_polardb_cluster":                            resourceAlicloudPolarDBCluster(),
			"alicloud_polardb_backup_policy":                      resourceAlicloudPolarDBBackupPolicy(),
			"alicloud_polardb_database":                           resourceAlicloudPolarDBDatabase(),
			"alicloud_polardb_account":                            resourceAlicloudPolarDBAccount(),
			"alicloud_polardb_account_privilege":                  resourceAlicloudPolarDBAccountPrivilege(),
			"alicloud_polardb_endpoint":                           resourceAlicloudPolarDBEndpoint(),
			"alicloud_polardb_endpoint_address":                   resourceAlicloudPolarDBEndpointAddress(),
			"alicloud_hbase_instance":                             resourceAlicloudHBaseInstance(),
			"alicloud_market_order":                               resourceAlicloudMarketOrder(),
			"alicloud_adb_cluster":                                resourceAlicloudAdbDbCluster(),
			"alicloud_adb_backup_policy":                          resourceAlicloudAdbBackupPolicy(),
			"alicloud_adb_account":                                resourceAlicloudAdbAccount(),
			"alicloud_adb_connection":                             resourceAlicloudAdbConnection(),
			"alicloud_cen_flowlog":                                resourceAlicloudCenFlowlog(),
			"alicloud_kms_secret":                                 resourceAlicloudKmsSecret(),
			"alicloud_maxcompute_project":                         resourceAlicloudMaxcomputeProject(),
			"alicloud_kms_alias":                                  resourceAlicloudKmsAlias(),
			"alicloud_dns_instance":                               resourceAlicloudAlidnsInstance(),
			"alicloud_dns_domain_attachment":                      resourceAlicloudAlidnsDomainAttachment(),
			"alicloud_alidns_domain_attachment":                   resourceAlicloudAlidnsDomainAttachment(),
			"alicloud_edas_application":                           resourceAlicloudEdasApplication(),
			"alicloud_edas_deploy_group":                          resourceAlicloudEdasDeployGroup(),
			"alicloud_edas_application_scale":                     resourceAlicloudEdasInstanceApplicationAttachment(),
			"alicloud_edas_slb_attachment":                        resourceAlicloudEdasSlbAttachment(),
			"alicloud_edas_cluster":                               resourceAlicloudEdasCluster(),
			"alicloud_edas_instance_cluster_attachment":           resourceAlicloudEdasInstanceClusterAttachment(),
			"alicloud_edas_application_deployment":                resourceAlicloudEdasApplicationPackageAttachment(),
			"alicloud_dns_domain":                                 resourceAlicloudAlidnsDomain(),
			"alicloud_dms_enterprise_instance":                    resourceAlicloudDmsEnterpriseInstance(),
			"alicloud_waf_domain":                                 resourceAlicloudWafDomain(),
			"alicloud_cen_route_map":                              resourceAlicloudCenRouteMap(),
			"alicloud_resource_manager_role":                      resourceAlicloudResourceManagerRole(),
			"alicloud_resource_manager_resource_group":            resourceAlicloudResourceManagerResourceGroup(),
			"alicloud_resource_manager_folder":                    resourceAlicloudResourceManagerFolder(),
			"alicloud_resource_manager_handshake":                 resourceAlicloudResourceManagerHandshake(),
			"alicloud_cen_private_zone":                           resourceAlicloudCenPrivateZone(),
			"alicloud_resource_manager_policy":                    resourceAlicloudResourceManagerPolicy(),
			"alicloud_resource_manager_account":                   resourceAlicloudResourceManagerAccount(),
			"alicloud_waf_instance":                               resourceAlicloudWafInstance(),
			"alicloud_resource_manager_resource_directory":        resourceAlicloudResourceManagerResourceDirectory(),
			"alicloud_alidns_domain_group":                        resourceAlicloudAlidnsDomainGroup(),
			"alicloud_resource_manager_policy_version":            resourceAlicloudResourceManagerPolicyVersion(),
			"alicloud_kms_key_version":                            resourceAlicloudKmsKeyVersion(),
			"alicloud_alidns_record":                              resourceAlicloudAlidnsRecord(),
			"alicloud_ddoscoo_scheduler_rule":                     resourceAlicloudDdoscooSchedulerRule(),
			"alicloud_cassandra_cluster":                          resourceAlicloudCassandraCluster(),
			"alicloud_cassandra_data_center":                      resourceAlicloudCassandraDataCenter(),
			"alicloud_cen_vbr_health_check":                       resourceAlicloudCenVbrHealthCheck(),
			"alicloud_eci_openapi_image_cache":                    resourceAlicloudEciImageCache(),
			"alicloud_eci_image_cache":                            resourceAlicloudEciImageCache(),
			"alicloud_dms_enterprise_user":                        resourceAlicloudDmsEnterpriseUser(),
			"alicloud_ecs_dedicated_host":                         resourceAlicloudEcsDedicatedHost(),
			"alicloud_oos_template":                               resourceAlicloudOosTemplate(),
			"alicloud_edas_k8s_cluster":                           resourceAlicloudEdasK8sCluster(),
			"alicloud_oos_execution":                              resourceAlicloudOosExecution(),
			"alicloud_resource_manager_policy_attachment":         resourceAlicloudResourceManagerPolicyAttachment(),
			"alicloud_dcdn_domain":                                resourceAlicloudDcdnDomain(),
			"alicloud_mse_cluster":                                resourceAlicloudMseCluster(),
			"alicloud_actiontrail_trail":                          resourceAlicloudActiontrailTrail(),
			"alicloud_actiontrail":                                resourceAlicloudActiontrailTrail(),
			"alicloud_alidns_domain":                              resourceAlicloudAlidnsDomain(),
			"alicloud_alidns_instance":                            resourceAlicloudAlidnsInstance(),
			"alicloud_edas_k8s_application":                       resourceAlicloudEdasK8sApplication(),
			"alicloud_config_rule":                                resourceAlicloudConfigRule(),
			"alicloud_config_configuration_recorder":              resourceAlicloudConfigConfigurationRecorder(),
			"alicloud_config_delivery_channel":                    resourceAlicloudConfigDeliveryChannel(),
			"alicloud_cms_alarm_contact":                          resourceAlicloudCmsAlarmContact(),
			"alicloud_cen_route_service":                          resourceAlicloudCenRouteService(),
			"alicloud_kvstore_connection":                         resourceAlicloudKvstoreConnection(),
			"alicloud_cms_alarm_contact_group":                    resourceAlicloudCmsAlarmContactGroup(),
			"alicloud_cms_group_metric_rule":                      resourceAlicloudCmsGroupMetricRule(),
			"alicloud_fnf_flow":                                   resourceAlicloudFnfFlow(),
			"alicloud_fnf_schedule":                               resourceAlicloudFnfSchedule(),
			"alicloud_ros_change_set":                             resourceAlicloudRosChangeSet(),
			"alicloud_ros_stack":                                  resourceAlicloudRosStack(),
			"alicloud_ros_stack_group":                            resourceAlicloudRosStackGroup(),
			"alicloud_ros_template":                               resourceAlicloudRosTemplate(),
			"alicloud_privatelink_vpc_endpoint_service":           resourceAlicloudPrivatelinkVpcEndpointService(),
			"alicloud_privatelink_vpc_endpoint":                   resourceAlicloudPrivatelinkVpcEndpoint(),
			"alicloud_privatelink_vpc_endpoint_connection":        resourceAlicloudPrivatelinkVpcEndpointConnection(),
			"alicloud_privatelink_vpc_endpoint_service_resource":  resourceAlicloudPrivatelinkVpcEndpointServiceResource(),
			"alicloud_privatelink_vpc_endpoint_service_user":      resourceAlicloudPrivatelinkVpcEndpointServiceUser(),
			"alicloud_resource_manager_resource_share":            resourceAlicloudResourceManagerResourceShare(),
			"alicloud_privatelink_vpc_endpoint_zone":              resourceAlicloudPrivatelinkVpcEndpointZone(),
			"alicloud_ga_accelerator":                             resourceAlicloudGaAccelerator(),
			"alicloud_eci_container_group":                        resourceAlicloudEciContainerGroup(),
			"alicloud_resource_manager_shared_resource":           resourceAlicloudResourceManagerSharedResource(),
			"alicloud_resource_manager_shared_target":             resourceAlicloudResourceManagerSharedTarget(),
			"alicloud_ga_listener":                                resourceAlicloudGaListener(),
			"alicloud_tsdb_instance":                              resourceAlicloudTsdbInstance(),
			"alicloud_ga_bandwidth_package":                       resourceAlicloudGaBandwidthPackage(),
			"alicloud_ga_endpoint_group":                          resourceAlicloudGaEndpointGroup(),
			"alicloud_brain_industrial_pid_organization":          resourceAlicloudBrainIndustrialPidOrganization(),
			"alicloud_ga_bandwidth_package_attachment":            resourceAlicloudGaBandwidthPackageAttachment(),
			"alicloud_ga_ip_set":                                  resourceAlicloudGaIpSet(),
			"alicloud_ga_forwarding_rule":                         resourceAlicloudGaForwardingRule(),
			"alicloud_eipanycast_anycast_eip_address":             resourceAlicloudEipanycastAnycastEipAddress(),
			"alicloud_brain_industrial_pid_project":               resourceAlicloudBrainIndustrialPidProject(),
			"alicloud_cms_monitor_group":                          resourceAlicloudCmsMonitorGroup(),
			"alicloud_eipanycast_anycast_eip_address_attachment":  resourceAlicloudEipanycastAnycastEipAddressAttachment(),
			"alicloud_ram_saml_provider":                          resourceAlicloudRamSamlProvider(),
			"alicloud_quotas_application_info":                    resourceAlicloudQuotasQuotaApplication(),
			"alicloud_cms_monitor_group_instances":                resourceAlicloudCmsMonitorGroupInstances(),
			"alicloud_quotas_quota_alarm":                         resourceAlicloudQuotasQuotaAlarm(),
			"alicloud_ecs_command":                                resourceAlicloudEcsCommand(),
			"alicloud_cloud_storage_gateway_storage_bundle":       resourceAlicloudCloudStorageGatewayStorageBundle(),
			"alicloud_ecs_hpc_cluster":                            resourceAlicloudEcsHpcCluster(),
			"alicloud_vpc_flow_log":                               resourceAlicloudVpcFlowLog(),
			"alicloud_brain_industrial_pid_loop":                  resourceAlicloudBrainIndustrialPidLoop(),
			"alicloud_quotas_quota_application":                   resourceAlicloudQuotasQuotaApplication(),
			"alicloud_ecs_auto_snapshot_policy":                   resourceAlicloudEcsAutoSnapshotPolicy(),
			"alicloud_rds_parameter_group":                        resourceAlicloudRdsParameterGroup(),
			"alicloud_ecs_launch_template":                        resourceAlicloudEcsLaunchTemplate(),
			"alicloud_resource_manager_control_policy":            resourceAlicloudResourceManagerControlPolicy(),
			"alicloud_resource_manager_control_policy_attachment": resourceAlicloudResourceManagerControlPolicyAttachment(),
			"alicloud_rds_account":                                resourceAlicloudRdsAccount(),
			"alicloud_ecs_snapshot":                               resourceAlicloudEcsSnapshot(),
			"alicloud_ecs_key_pair":                               resourceAlicloudEcsKeyPair(),
			"alicloud_ecs_key_pair_attachment":                    resourceAlicloudEcsKeyPairAttachment(),
			"alicloud_adb_db_cluster":                             resourceAlicloudAdbDbCluster(),
			"alicloud_ecs_disk":                                   resourceAlicloudEcsDisk(),
			"alicloud_ecs_disk_attachment":                        resourceAlicloudEcsDiskAttachment(),
			"alicloud_ecs_auto_snapshot_policy_attachment":        resourceAlicloudEcsAutoSnapshotPolicyAttachment(),
			"alicloud_ddoscoo_domain_resource":                    resourceAlicloudDdoscooDomainResource(),
			"alicloud_ddoscoo_port":                               resourceAlicloudDdoscooPort(),
			"alicloud_slb_load_balancer":                          resourceAlicloudSlbLoadBalancer(),
			"alicloud_ecs_network_interface":                      resourceAlicloudEcsNetworkInterface(),
			"alicloud_ecs_network_interface_attachment":           resourceAlicloudEcsNetworkInterfaceAttachment(),
			"alicloud_config_aggregator":                          resourceAlicloudConfigAggregator(),
			"alicloud_config_aggregate_config_rule":               resourceAlicloudConfigAggregateConfigRule(),
			"alicloud_config_aggregate_compliance_pack":           resourceAlicloudConfigAggregateCompliancePack(),
			"alicloud_config_compliance_pack":                     resourceAlicloudConfigCompliancePack(),
			"alicloud_direct_mail_receivers":                      resourceAlicloudDirectMailReceivers(),
			"alicloud_eip_address":                                resourceAlicloudEipAddress(),
			"alicloud_event_bridge_event_bus":                     resourceAlicloudEventBridgeEventBus(),
			"alicloud_amqp_virtual_host":                          resourceAlicloudAmqpVirtualHost(),
			"alicloud_amqp_queue":                                 resourceAlicloudAmqpQueue(),
			"alicloud_amqp_exchange":                              resourceAlicloudAmqpExchange(),
			"alicloud_cassandra_backup_plan":                      resourceAlicloudCassandraBackupPlan(),
			"alicloud_cen_transit_router_peer_attachment":         resourceAlicloudCenTransitRouterPeerAttachment(),
			"alicloud_amqp_instance":                              resourceAlicloudAmqpInstance(),
			"alicloud_hbr_vault":                                  resourceAlicloudHbrVault(),
			"alicloud_ssl_certificates_service_certificate":       resourceAlicloudSslCertificatesServiceCertificate(),
			"alicloud_arms_alert_contact":                         resourceAlicloudArmsAlertContact(),
			"alicloud_event_bridge_slr":                           resourceAlicloudEventBridgeSlr(),
			"alicloud_event_bridge_rule":                          resourceAlicloudEventBridgeRule(),
			"alicloud_cloud_firewall_control_policy":              resourceAlicloudCloudFirewallControlPolicy(),
			"alicloud_sae_namespace":                              resourceAlicloudSaeNamespace(),
			"alicloud_sae_config_map":                             resourceAlicloudSaeConfigMap(),
			"alicloud_alb_security_policy":                        resourceAlicloudAlbSecurityPolicy(),
			"alicloud_kvstore_audit_log_config":                   resourceAlicloudKvstoreAuditLogConfig(),
			"alicloud_event_bridge_event_source":                  resourceAlicloudEventBridgeEventSource(),
			"alicloud_cloud_firewall_control_policy_order":        resourceAlicloudCloudFirewallControlPolicyOrder(),
			"alicloud_ecd_policy_group":                           resourceAlicloudEcdPolicyGroup(),
			"alicloud_ecp_key_pair":                               resourceAlicloudEcpKeyPair(),
			"alicloud_hbr_ecs_backup_plan":                        resourceAlicloudHbrEcsBackupPlan(),
			"alicloud_hbr_nas_backup_plan":                        resourceAlicloudHbrNasBackupPlan(),
			"alicloud_hbr_oss_backup_plan":                        resourceAlicloudHbrOssBackupPlan(),
			"alicloud_scdn_domain":                                resourceAlicloudScdnDomain(),
			"alicloud_alb_server_group":                           resourceAlicloudAlbServerGroup(),
			"alicloud_data_works_folder":                          resourceAlicloudDataWorksFolder(),
			"alicloud_arms_alert_contact_group":                   resourceAlicloudArmsAlertContactGroup(),
			"alicloud_dcdn_domain_config":                         resourceAlicloudDcdnDomainConfig(),
			"alicloud_scdn_domain_config":                         resourceAlicloudScdnDomainConfig(),
			"alicloud_cloud_storage_gateway_gateway":              resourceAlicloudCloudStorageGatewayGateway(),
			"alicloud_lindorm_instance":                           resourceAlicloudLindormInstance(),
			"alicloud_cddc_dedicated_host_group":                  resourceAlicloudCddcDedicatedHostGroup(),
			"alicloud_hbr_ecs_backup_client":                      resourceAlicloudHbrEcsBackupClient(),
			"alicloud_msc_sub_contact":                            resourceAlicloudMscSubContact(),
			"alicloud_express_connect_physical_connection":        resourceAlicloudExpressConnectPhysicalConnection(),
			"alicloud_alb_load_balancer":                          resourceAlicloudAlbLoadBalancer(),
			"alicloud_sddp_rule":                                  resourceAlicloudSddpRule(),
			"alicloud_bastionhost_user_group":                     resourceAlicloudBastionhostUserGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var providerConfig map[string]interface{}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var getProviderConfig = func(str string, key string) string {
		if str == "" {
			value, err := getConfigFromProfile(d, key)
			if err == nil && value != nil {
				str = value.(string)
			}
		}
		return str
	}

	accessKey := getProviderConfig(d.Get("access_key").(string), "access_key_id")
	secretKey := getProviderConfig(d.Get("secret_key").(string), "access_key_secret")
	region := getProviderConfig(d.Get("region").(string), "region_id")
	if region == "" {
		region = DEFAULT_REGION
	}

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	config := &connectivity.Config{
		SourceIp:             strings.TrimSpace(d.Get("source_ip").(string)),
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
		ClientReadTimeout:    d.Get("client_read_timeout").(int),
		ClientConnectTimeout: d.Get("client_connect_timeout").(int),
	}
	token := getProviderConfig(d.Get("security_token").(string), "sts_token")
	config.SecurityToken = strings.TrimSpace(token)

	config.RamRoleArn = getProviderConfig("", "ram_role_arn")
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
	expiredSeconds, err := getConfigFromProfile(d, "expired_seconds")
	if err == nil && expiredSeconds != nil {
		config.RamRoleSessionExpiration = (int)(expiredSeconds.(float64))
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if assumeRole["role_arn"].(string) != "" {
			config.RamRoleArn = assumeRole["role_arn"].(string)
		}
		if assumeRole["session_name"].(string) != "" {
			config.RamRoleSessionName = assumeRole["session_name"].(string)
		}
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("ALICLOUD_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
				if expiredSeconds, err := strconv.Atoi(v); err == nil {
					config.RamRoleSessionExpiration = expiredSeconds
				}
			}
			if config.RamRoleSessionExpiration == 0 {
				config.RamRoleSessionExpiration = 3600
			}
		} else {
			config.RamRoleSessionExpiration = assumeRole["session_expiration"].(int)
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration)
	}

	if err := config.MakeConfigByEcsRoleName(); err != nil {
		return nil, err
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)
	endpointInit := make(map[string]interface{})
	config.Endpoints = endpointInit

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.Endpoints = endpoints
		config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
		config.RdsEndpoint = strings.TrimSpace(endpoints["rds"].(string))
		config.SlbEndpoint = strings.TrimSpace(endpoints["slb"].(string))
		config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
		config.EssEndpoint = strings.TrimSpace(endpoints["ess"].(string))
		config.OssEndpoint = strings.TrimSpace(endpoints["oss"].(string))
		config.OnsEndpoint = strings.TrimSpace(endpoints["ons"].(string))
		config.AlikafkaEndpoint = strings.TrimSpace(endpoints["alikafka"].(string))
		config.DnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		config.RamEndpoint = strings.TrimSpace(endpoints["ram"].(string))
		config.CsEndpoint = strings.TrimSpace(endpoints["cs"].(string))
		config.CrEndpoint = strings.TrimSpace(endpoints["cr"].(string))
		config.CdnEndpoint = strings.TrimSpace(endpoints["cdn"].(string))
		config.KmsEndpoint = strings.TrimSpace(endpoints["kms"].(string))
		config.OtsEndpoint = strings.TrimSpace(endpoints["ots"].(string))
		config.CmsEndpoint = strings.TrimSpace(endpoints["cms"].(string))
		config.PvtzEndpoint = strings.TrimSpace(endpoints["pvtz"].(string))
		config.StsEndpoint = strings.TrimSpace(endpoints["sts"].(string))
		config.LogEndpoint = strings.TrimSpace(endpoints["log"].(string))
		config.DrdsEndpoint = strings.TrimSpace(endpoints["drds"].(string))
		config.DdsEndpoint = strings.TrimSpace(endpoints["dds"].(string))
		config.GpdbEnpoint = strings.TrimSpace(endpoints["gpdb"].(string))
		config.KVStoreEndpoint = strings.TrimSpace(endpoints["kvstore"].(string))
		config.PolarDBEndpoint = strings.TrimSpace(endpoints["polardb"].(string))
		config.FcEndpoint = strings.TrimSpace(endpoints["fc"].(string))
		config.ApigatewayEndpoint = strings.TrimSpace(endpoints["apigateway"].(string))
		config.DatahubEndpoint = strings.TrimSpace(endpoints["datahub"].(string))
		config.MnsEndpoint = strings.TrimSpace(endpoints["mns"].(string))
		config.LocationEndpoint = strings.TrimSpace(endpoints["location"].(string))
		config.ElasticsearchEndpoint = strings.TrimSpace(endpoints["elasticsearch"].(string))
		config.NasEndpoint = strings.TrimSpace(endpoints["nas"].(string))
		config.ActiontrailEndpoint = strings.TrimSpace(endpoints["actiontrail"].(string))
		config.BssOpenApiEndpoint = strings.TrimSpace(endpoints["bssopenapi"].(string))
		config.DdoscooEndpoint = strings.TrimSpace(endpoints["ddoscoo"].(string))
		config.DdosbgpEndpoint = strings.TrimSpace(endpoints["ddosbgp"].(string))
		config.EmrEndpoint = strings.TrimSpace(endpoints["emr"].(string))
		config.CasEndpoint = strings.TrimSpace(endpoints["cas"].(string))
		config.MarketEndpoint = strings.TrimSpace(endpoints["market"].(string))
		config.AdbEndpoint = strings.TrimSpace(endpoints["adb"].(string))
		config.CbnEndpoint = strings.TrimSpace(endpoints["cbn"].(string))
		config.MaxComputeEndpoint = strings.TrimSpace(endpoints["maxcompute"].(string))
		config.DmsEnterpriseEndpoint = strings.TrimSpace(endpoints["dms_enterprise"].(string))
		config.WafOpenapiEndpoint = strings.TrimSpace(endpoints["waf_openapi"].(string))
		config.ResourcemanagerEndpoint = strings.TrimSpace(endpoints["resourcemanager"].(string))
		config.EciEndpoint = strings.TrimSpace(endpoints["eci"].(string))
		config.OosEndpoint = strings.TrimSpace(endpoints["oos"].(string))
		config.DcdnEndpoint = strings.TrimSpace(endpoints["dcdn"].(string))
		config.MseEndpoint = strings.TrimSpace(endpoints["mse"].(string))
		config.ConfigEndpoint = strings.TrimSpace(endpoints["config"].(string))
		config.RKvstoreEndpoint = strings.TrimSpace(endpoints["r_kvstore"].(string))
		config.FnfEndpoint = strings.TrimSpace(endpoints["fnf"].(string))
		config.RosEndpoint = strings.TrimSpace(endpoints["ros"].(string))
		config.PrivatelinkEndpoint = strings.TrimSpace(endpoints["privatelink"].(string))
		config.ResourcesharingEndpoint = strings.TrimSpace(endpoints["resourcesharing"].(string))
		config.GaEndpoint = strings.TrimSpace(endpoints["ga"].(string))
		config.HitsdbEndpoint = strings.TrimSpace(endpoints["hitsdb"].(string))
		config.BrainIndustrialEndpoint = strings.TrimSpace(endpoints["brain_industrial"].(string))
		config.EipanycastEndpoint = strings.TrimSpace(endpoints["eipanycast"].(string))
		config.ImsEndpoint = strings.TrimSpace(endpoints["ims"].(string))
		config.QuotasEndpoint = strings.TrimSpace(endpoints["quotas"].(string))
		config.SgwEndpoint = strings.TrimSpace(endpoints["sgw"].(string))
		config.DmEndpoint = strings.TrimSpace(endpoints["dm"].(string))
		config.EventbridgeEndpoint = strings.TrimSpace(endpoints["eventbridge"].(string))
		config.OnsproxyEndpoint = strings.TrimSpace(endpoints["onsproxy"].(string))
		config.CdsEndpoint = strings.TrimSpace(endpoints["cds"].(string))
		config.HbrEndpoint = strings.TrimSpace(endpoints["hbr"].(string))
		config.ArmsEndpoint = strings.TrimSpace(endpoints["arms"].(string))
		config.ServerlessEndpoint = strings.TrimSpace(endpoints["serverless"].(string))
		config.AlbEndpoint = strings.TrimSpace(endpoints["alb"].(string))
		config.RedisaEndpoint = strings.TrimSpace(endpoints["redisa"].(string))
		config.GwsecdEndpoint = strings.TrimSpace(endpoints["gwsecd"].(string))
		config.CloudphoneEndpoint = strings.TrimSpace(endpoints["cloudphone"].(string))
		config.ScdnEndpoint = strings.TrimSpace(endpoints["scdn"].(string))
		config.DataworkspublicEndpoint = strings.TrimSpace(endpoints["dataworkspublic"].(string))
		config.HcsSgwEndpoint = strings.TrimSpace(endpoints["hcs_sgw"].(string))
		config.CddcEndpoint = strings.TrimSpace(endpoints["cddc"].(string))
		config.MscopensubscriptionEndpoint = strings.TrimSpace(endpoints["mscopensubscription"].(string))
		config.SddpEndpoint = strings.TrimSpace(endpoints["sddp"].(string))

		config.BastionhostEndpoint = strings.TrimSpace(endpoints["bastionhost"].(string))
		if endpoint, ok := endpoints["alidns"]; ok {
			config.AlidnsEndpoint = strings.TrimSpace(endpoint.(string))
		} else {
			config.AlidnsEndpoint = strings.TrimSpace(endpoints["dns"].(string))
		}
		config.CassandraEndpoint = strings.TrimSpace(endpoints["cassandra"].(string))
	}

	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config.AccessKey, config.SecretKey, config.SecurityToken, region, config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.StsEndpoint)
		if err != nil {
			return nil, err
		}
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(ots_instance_name.(string))
	}

	if logEndpoint, ok := d.GetOk("log_endpoint"); ok && logEndpoint.(string) != "" {
		config.LogEndpoint = strings.TrimSpace(logEndpoint.(string))
	}
	if mnsEndpoint, ok := d.GetOk("mns_endpoint"); ok && mnsEndpoint.(string) != "" {
		config.MnsEndpoint = strings.TrimSpace(mnsEndpoint.(string))
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}

	if fcEndpoint, ok := d.GetOk("fc"); ok && fcEndpoint.(string) != "" {
		config.FcEndpoint = strings.TrimSpace(fcEndpoint.(string))
	}

	if config.ConfigurationSource == "" {
		sourceName := fmt.Sprintf("Default/%s:%s", config.AccessKey, strings.Trim(uuid.New().String(), "-"))
		if len(sourceName) > 64 {
			sourceName = sourceName[:64]
		}
		config.ConfigurationSource = sourceName
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// This is a global MutexKV for use within this plugin.
var alicloudMutexKV = mutexkv.NewMutexKV()

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"source_ip": "The access key for API operations. You can retrieve this from the 'Security Management' section of the Alibaba Cloud console.",

		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the Alibaba Cloud console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the Alibaba Cloud console.",

		"ecs_role_name": "The RAM Role Name attached on a ECS instance for API operations. You can retrieve this from the 'Access Control' section of the Alibaba Cloud console.",

		"region": "The region where Alibaba Cloud operations will take place. Examples are cn-beijing, cn-hangzhou, eu-central-1, etc.",

		"security_token": "security token. A security token is only required if you are using Security Token Service.",

		"account_id": "The account ID for some service API operations. You can retrieve this from the 'Security Settings' section of the Alibaba Cloud console.",

		"profile": "The profile for API operations. If not set, the default profile created with `aliyun configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set this defaults to ~/.aliyun/config.json",

		"assume_role_role_arn": "The ARN of a RAM role to assume prior to making API calls.",

		"assume_role_session_name": "The session name to use when assuming the role. If omitted, `terraform` is passed to the AssumeRole call as session name.",

		"assume_role_policy": "The permissions applied when assuming a role. You cannot use, this policy to grant further permissions that are in excess to those of the, role that is being assumed.",

		"assume_role_session_expiration": "The time after which the established session for assuming role expires. Valid value range: [900-3600] seconds. Default to 0 (in this case Alicloud use own default value).",

		"skip_region_validation": "Skip static validation of region ID. Used by users of alternative AlibabaCloud-like APIs or users w/ access to regions that are not public (yet).",

		"configuration_source": "Use this to mark a terraform configuration file source.",

		"client_read_timeout":    "The maximum timeout of the client read request.",
		"client_connect_timeout": "The maximum timeout of the client connection server.",
		"ecs_endpoint":           "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.",

		"rds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.",

		"slb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.",

		"vpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.",

		"ess_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.",

		"oss_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.",

		"ons_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ONS endpoints.",

		"alikafka_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ALIKAFKA endpoints.",

		"dns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DNS endpoints.",

		"ram_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RAM endpoints.",

		"cs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Service endpoints.",

		"cr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Container Registry endpoints.",

		"cdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CDN endpoints.",

		"kms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom KMS endpoints.",

		"ots_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Table Store endpoints.",

		"cms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Cloud Monitor endpoints.",

		"pvtz_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Private Zone endpoints.",

		"sts_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom STS endpoints.",

		"log_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Log Service endpoints.",

		"drds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DRDS endpoints.",

		"dds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MongoDB endpoints.",

		"polardb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom PolarDB endpoints.",

		"gpdb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom GPDB endpoints.",

		"kvstore_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom R-KVStore endpoints.",

		"fc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Function Computing endpoints.",

		"apigateway_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Api Gateway endpoints.",

		"datahub_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Datahub endpoints.",

		"mns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MNS endpoints.",

		"location_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Location Service endpoints.",

		"elasticsearch_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Elasticsearch endpoints.",

		"nas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom NAS endpoints.",

		"actiontrail_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Actiontrail endpoints.",

		"cas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CAS endpoints.",

		"bssopenapi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BSSOPENAPI endpoints.",

		"ddoscoo_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DDOSCOO endpoints.",

		"ddosbgp_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DDOSBGP endpoints.",

		"emr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom EMR endpoints.",

		"market_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Market Place endpoints.",

		"hbase_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom HBase endpoints.",

		"adb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom AnalyticDB endpoints.",

		"cbn_endpoint":        "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cbn endpoints.",
		"maxcompute_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom MaxCompute endpoints.",

		"dms_enterprise_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dms_enterprise endpoints.",

		"waf_openapi_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom waf_openapi endpoints.",

		"resourcemanager_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom resourcemanager endpoints.",

		"alidns_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom alidns endpoints.",

		"cassandra_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cassandra endpoints.",

		"eci_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eci endpoints.",

		"oos_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom oos endpoints.",

		"dcdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dcdn endpoints.",

		"mse_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom mse endpoints.",

		"config_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom config endpoints.",

		"r_kvstore_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom r_kvstore endpoints.",

		"fnf_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom fnf endpoints.",

		"ros_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ros endpoints.",

		"privatelink_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom privatelink endpoints.",

		"resourcesharing_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom resourcesharing endpoints.",

		"ga_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ga endpoints.",

		"hitsdb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hitsdb endpoints.",

		"brain_industrial_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom brain_industrial endpoints.",

		"eipanycast_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eipanycast endpoints.",

		"ims_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ims endpoints.",

		"quotas_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom quotas endpoints.",

		"sgw_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom sgw endpoints.",

		"dm_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dm endpoints.",

		"eventbridge_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom eventbridge_share endpoints.",

		"onsproxy_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom onsproxy endpoints.",

		"cds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cds endpoints.",

		"hbr_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hbr endpoints.",

		"arms_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom arms endpoints.",

		"serverless_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom serverless endpoints.",

		"alb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom alb endpoints.",

		"redisa_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom redisa endpoints.",

		"gwsecd_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom gwsecd endpoints.",

		"cloudphone_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cloudphone endpoints.",

		"scdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom scdn endpoints.",

		"dataworkspublic_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom dataworkspublic endpoints.",

		"hcs_sgw_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom hcs_sgw endpoints.",

		"cddc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom cddc endpoints.",

		"mscopensubscription_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom mscopensubscription endpoints.",

		"sddp_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom sddp endpoints.",

		"bastionhost_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom bastionhost endpoints.",
	}
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_role_arn"],
					DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ASSUME_ROLE_ARN", ""),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.EnvDefaultFunc("ALICLOUD_ASSUME_ROLE_SESSION_NAME", ""),
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
				"session_expiration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  descriptions["assume_role_session_expiration"],
					ValidateFunc: intBetween(900, 3600),
				},
			},
		},
	}
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bastionhost": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bastionhost_endpoint"],
				},

				"cddc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cddc_endpoint"],
				},
				"sddp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sddp_endpoint"],
				},

				"mscopensubscription": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mscopensubscription_endpoint"],
				},

				"dataworkspublic": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dataworkspublic_endpoint"],
				},

				"hcs_sgw": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hcs_sgw_endpoint"],
				},

				"cloudphone": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cloudphone_endpoint"],
				},

				"alb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alb_endpoint"],
				},
				"redisa": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["redisa_endpoint"],
				},
				"gwsecd": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gwsecd_endpoint"],
				},
				"scdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["scdn_endpoint"],
				},

				"arms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["arms_endpoint"],
				},
				"serverless": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["serverless_endpoint"],
				},

				"hbr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hbr_endpoint"],
				},

				"onsproxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["onsproxy_endpoint"],
				},
				"cds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cds_endpoint"],
				},

				"dm": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dm_endpoint"],
				},

				"eventbridge": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eventbridge_endpoint"],
				},

				"sgw": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sgw_endpoint"],
				},

				"quotas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["quotas_endpoint"],
				},

				"ims": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ims_endpoint"],
				},

				"brain_industrial": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["brain_industrial_endpoint"],
				},

				"resourcesharing": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["resourcesharing_endpoint"],
				},
				"ga": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ga_endpoint"],
				},

				"hitsdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["hitsdb_endpoint"],
				},

				"privatelink": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["privatelink_endpoint"],
				},

				"eipanycast": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eipanycast_endpoint"],
				},

				"fnf": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fnf_endpoint"],
				},

				"ros": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ros_endpoint"],
				},

				"r_kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["r_kvstore_endpoint"],
				},

				"config": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["config_endpoint"],
				},

				"dcdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dcdn_endpoint"],
				},

				"mse": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mse_endpoint"],
				},

				"oos": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oos_endpoint"],
				},

				"eci": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eci_endpoint"],
				},

				"alidns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alidns_endpoint"],
				},

				"resourcemanager": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["resourcemanager_endpoint"],
				},

				"waf_openapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["waf_openapi_endpoint"],
				},

				"dms_enterprise": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dms_enterprise_endpoint"],
				},

				"cassandra": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cassandra_endpoint"],
				},

				"cbn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cbn_endpoint"],
				},

				"ecs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ecs_endpoint"],
				},
				"rds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"slb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["slb_endpoint"],
				},
				"vpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpc_endpoint"],
				},
				"ess": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ess_endpoint"],
				},
				"oss": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oss_endpoint"],
				},
				"ons": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ons_endpoint"],
				},
				"alikafka": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alikafka_endpoint"],
				},
				"dns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dns_endpoint"],
				},
				"ram": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ram_endpoint"],
				},
				"cs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cs_endpoint"],
				},
				"cr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cr_endpoint"],
				},
				"cdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cdn_endpoint"],
				},

				"kms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kms_endpoint"],
				},

				"ots": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ots_endpoint"],
				},

				"cms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cms_endpoint"],
				},

				"pvtz": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["pvtz_endpoint"],
				},

				"sts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sts_endpoint"],
				},
				// log service is sls service
				"log": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["log_endpoint"],
				},
				"drds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["drds_endpoint"],
				},
				"dds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dds_endpoint"],
				},
				"polardb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["polardb_endpoint"],
				},
				"gpdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gpdb_endpoint"],
				},
				"kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kvstore_endpoint"],
				},
				"fc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fc_endpoint"],
				},
				"apigateway": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["apigateway_endpoint"],
				},
				"datahub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["datahub_endpoint"],
				},
				"mns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["location_endpoint"],
				},
				"elasticsearch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elasticsearch_endpoint"],
				},
				"nas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["nas_endpoint"],
				},
				"actiontrail": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["actiontrail_endpoint"],
				},
				"cas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cas_endpoint"],
				},
				"bssopenapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bssopenapi_endpoint"],
				},
				"ddoscoo": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddoscoo_endpoint"],
				},
				"ddosbgp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddosbgp_endpoint"],
				},
				"emr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["emr_endpoint"],
				},
				"market": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["market_endpoint"],
				},
				"adb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["adb_endpoint"],
				},
				"maxcompute": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["maxcompute_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["ecs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["slb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ess"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oss"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ons"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alikafka"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ram"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ots"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["pvtz"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sts"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["log"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["drds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gpdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["polardb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["apigateway"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["datahub"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["location"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elasticsearch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["nas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["actiontrail"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bssopenapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddoscoo"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddosbgp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["emr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["market"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["adb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cbn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["maxcompute"].(string)))

	buf.WriteString(fmt.Sprintf("%s-", m["dms_enterprise"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["waf_openapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["resourcemanager"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alidns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cassandra"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eci"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oos"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dcdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mse"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["config"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["r_kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fnf"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ros"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["privatelink"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["resourcesharing"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ga"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hitsdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["brain_industrial"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eipanycast"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ims"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["quotas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sgw"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dm"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eventbridge"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["onsproxy"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hbr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["arms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["serverless"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["redisa"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gwsecd"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cloudphone"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["scdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dataworkspublic"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["hcs_sgw"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cddc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mscopensubscription"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sddp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bastionhost"].(string)))
	return hashcode.String(buf.String())
}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {

	if providerConfig == nil {
		if v, ok := d.GetOk("profile"); !ok && v.(string) == "" {
			return nil, nil
		}
		current := d.Get("profile").(string)
		// Set CredsFilename, expanding home directory
		profilePath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
		if err != nil {
			return nil, WrapError(err)
		}
		if profilePath == "" {
			profilePath = fmt.Sprintf("%s/.aliyun/config.json", os.Getenv("HOME"))
			if runtime.GOOS == "windows" {
				profilePath = fmt.Sprintf("%s/.aliyun/config.json", os.Getenv("USERPROFILE"))
			}
		}
		providerConfig = make(map[string]interface{})
		_, err = os.Stat(profilePath)
		if !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(profilePath)
			if err != nil {
				return nil, WrapError(err)
			}
			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, WrapError(err)
			}
			for _, v := range config["profiles"].([]interface{}) {
				if current == v.(map[string]interface{})["name"] {
					providerConfig = v.(map[string]interface{})
				}
			}
		}
	}

	mode := ""
	if v, ok := providerConfig["mode"]; ok {
		mode = v.(string)
	} else {
		return v, nil
	}
	switch ProfileKey {
	case "access_key_id", "access_key_secret":
		if mode == "EcsRamRole" {
			return "", nil
		}
	case "ram_role_name":
		if mode != "EcsRamRole" {
			return "", nil
		}
	case "sts_token":
		if mode != "StsToken" {
			return "", nil
		}
	case "ram_role_arn", "ram_session_name":
		if mode != "RamRoleArn" {
			return "", nil
		}
	case "expired_seconds":
		if mode != "RamRoleArn" {
			return float64(0), nil
		}
	}

	return providerConfig[ProfileKey], nil
}

func getAssumeRoleAK(accessKey, secretKey, stsToken, region, roleArn, sessionName, policy string, sessionExpiration int, stsEndpoint string) (string, string, string, error) {
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = roleArn
	request.RoleSessionName = sessionName
	request.DurationSeconds = requests.NewInteger(sessionExpiration)
	request.Policy = policy
	request.Scheme = "https"
	request.Domain = stsEndpoint

	var client *sts.Client
	var err error
	if stsToken == "" {
		client, err = sts.NewClientWithAccessKey(region, accessKey, secretKey)
	} else {
		client, err = sts.NewClientWithStsToken(region, accessKey, secretKey, stsToken)
	}

	if err != nil {
		return "", "", "", err
	}

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}
