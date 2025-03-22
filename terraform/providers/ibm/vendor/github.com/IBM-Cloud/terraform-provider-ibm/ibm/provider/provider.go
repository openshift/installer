// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/apigateway"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/appconfiguration"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/appid"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/atracker"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/backuprecovery"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/catalogmanagement"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cdtektonpipeline"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cdtoolchain"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cis"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/classicinfrastructure"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cloudant"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cloudfoundry"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cloudshell"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/codeengine"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/configurationaggregator"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/contextbasedrestrictions"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/cos"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/database"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/db2"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/directlink"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/dnsservices"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/enterprise"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/eventnotification"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/eventstreams"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/functions"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/globaltagging"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/hpcs"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/iamaccessgroup"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/iamidentity"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/iampolicy"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/kms"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/kubernetes"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/logs"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/logsrouting"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/metricsrouter"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/mqcloud"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/pag"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/partnercentersell"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/power"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/project"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/pushnotification"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/registry"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcemanager"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/satellite"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/scc"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/schematics"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/sdsaas"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/secretsmanager"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/transitgateway"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/usagereports"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vmware"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	provider := schema.Provider{
		Schema: map[string]*schema.Schema{
			"bluemix_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_API_KEY", "BLUEMIX_API_KEY"}, nil),
				Deprecated:  "This field is deprecated please use ibmcloud_api_key",
			},
			"bluemix_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any Bluemix API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_TIMEOUT", "BLUEMIX_TIMEOUT"}, nil),
				Deprecated:  "This field is deprecated please use ibmcloud_timeout",
			},
			"ibmcloud_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM Cloud API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_API_KEY", "IBMCLOUD_API_KEY"}, nil),
			},
			"ibmcloud_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any IBM Cloud API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_TIMEOUT", "IBMCLOUD_TIMEOUT"}, 60),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM cloud Region (for example 'us-south').",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_REGION", "IBMCLOUD_REGION", "BM_REGION", "BLUEMIX_REGION"}, "us-south"),
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM cloud Region zone (for example 'us-south-1') for power resources.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_ZONE", "IBMCLOUD_ZONE"}, ""),
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Resource group id.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_RESOURCE_GROUP", "IBMCLOUD_RESOURCE_GROUP", "BM_RESOURCE_GROUP", "BLUEMIX_RESOURCE_GROUP"}, ""),
			},
			"softlayer_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SoftLayer API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_API_KEY", "SOFTLAYER_API_KEY"}, nil),
				Deprecated:  "This field is deprecated please use iaas_classic_api_key",
			},
			"softlayer_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SoftLayer user name",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_USERNAME", "SOFTLAYER_USERNAME"}, nil),
				Deprecated:  "This field is deprecated please use iaas_classic_username",
			},
			"softlayer_endpoint_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Softlayer Endpoint",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_ENDPOINT_URL", "SOFTLAYER_ENDPOINT_URL"}, nil),
				Deprecated:  "This field is deprecated please use iaas_classic_endpoint_url",
			},
			"softlayer_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any SoftLayer API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_TIMEOUT", "SOFTLAYER_TIMEOUT"}, nil),
				Deprecated:  "This field is deprecated please use iaas_classic_timeout",
			},
			"iaas_classic_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Classic Infrastructure API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAAS_CLASSIC_API_KEY"}, nil),
			},
			"iaas_classic_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Classic Infrastructure API user name",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAAS_CLASSIC_USERNAME"}, nil),
			},
			"iaas_classic_endpoint_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Classic Infrastructure Endpoint",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAAS_CLASSIC_ENDPOINT_URL"}, "https://api.softlayer.com/rest/v3"),
			},
			"iaas_classic_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any Classic Infrastructure API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IAAS_CLASSIC_TIMEOUT"}, 60),
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The retry count to set for API calls.",
				DefaultFunc: schema.EnvDefaultFunc("MAX_RETRIES", 10),
			},
			"function_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM Cloud Function namespace",
				DefaultFunc: schema.EnvDefaultFunc("FUNCTION_NAMESPACE", nil),
				Deprecated:  "This field will be deprecated soon",
			},
			"riaas_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The next generation infrastructure service endpoint url.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"RIAAS_ENDPOINT"}, nil),
				Deprecated:  "This field is deprecated use generation",
			},
			"generation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Generation of Virtual Private Cloud. Default is 2",
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_GENERATION", "IBMCLOUD_GENERATION"}, nil),
				Deprecated: "The generation field is deprecated and will be removed after couple of releases",
			},
			"iam_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IAM Trusted Profile Authentication token",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_IAM_PROFILE_ID", "IBMCLOUD_IAM_PROFILE_ID"}, nil),
			},
			"iam_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IAM Authentication token",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_IAM_TOKEN", "IBMCLOUD_IAM_TOKEN"}, nil),
			},
			"iam_refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IAM Authentication refresh token",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_IAM_REFRESH_TOKEN", "IBMCLOUD_IAM_REFRESH_TOKEN"}, nil),
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"public", "private", "public-and-private"}),
				Description:  "Visibility of the provider if it is private or public.",
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"IC_VISIBILITY", "IBMCLOUD_VISIBILITY"}, "public"),
			},
			"endpoints_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the file that contains private and public regional endpoints mapping",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_ENDPOINTS_FILE_PATH", "IBMCLOUD_ENDPOINTS_FILE_PATH"}, nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ibm_api_gateway":                      apigateway.DataSourceIBMApiGateway(),
			"ibm_account":                          cloudfoundry.DataSourceIBMAccount(),
			"ibm_app":                              cloudfoundry.DataSourceIBMApp(),
			"ibm_app_domain_private":               cloudfoundry.DataSourceIBMAppDomainPrivate(),
			"ibm_app_domain_shared":                cloudfoundry.DataSourceIBMAppDomainShared(),
			"ibm_app_route":                        cloudfoundry.DataSourceIBMAppRoute(),
			"ibm_config_aggregator_configurations": configurationaggregator.AddConfigurationAggregatorInstanceFields(configurationaggregator.DataSourceIbmConfigAggregatorConfigurations()),
			"ibm_config_aggregator_settings":       configurationaggregator.AddConfigurationAggregatorInstanceFields(configurationaggregator.DataSourceIbmConfigAggregatorSettings()),
			"ibm_config_aggregator_resource_collection_status": configurationaggregator.AddConfigurationAggregatorInstanceFields(configurationaggregator.DataSourceIbmConfigAggregatorResourceCollectionStatus()),

			// // BackupAndRecovery
			"ibm_backup_recovery_agent_upgrade_tasks":      backuprecovery.DataSourceIbmBackupRecoveryAgentUpgradeTasks(),
			"ibm_backup_recovery_download_agent":           backuprecovery.DataSourceIbmBackupRecoveryDownloadAgent(),
			"ibm_backup_recovery_search_indexed_object":    backuprecovery.DataSourceIbmBackupRecoverySearchIndexedObject(),
			"ibm_backup_recovery_object_snapshots":         backuprecovery.DataSourceIbmBackupRecoveryObjectSnapshots(),
			"ibm_backup_recovery_connectors_metadata":      backuprecovery.DataSourceIbmBackupRecoveryConnectorsMetadata(),
			"ibm_backup_recovery_data_source_connections":  backuprecovery.DataSourceIbmBackupRecoveryDataSourceConnections(),
			"ibm_backup_recovery_data_source_connectors":   backuprecovery.DataSourceIbmBackupRecoveryDataSourceConnectors(),
			"ibm_backup_recovery_search_objects":           backuprecovery.DataSourceIbmBackupRecoverySearchObjects(),
			"ibm_backup_recovery_search_protected_objects": backuprecovery.DataSourceIbmBackupRecoverySearchProtectedObjects(),
			"ibm_backup_recovery_protection_group":         backuprecovery.DataSourceIbmBackupRecoveryProtectionGroup(),
			"ibm_backup_recovery_protection_groups":        backuprecovery.DataSourceIbmBackupRecoveryProtectionGroups(),
			"ibm_backup_recovery_protection_group_runs":    backuprecovery.DataSourceIbmBackupRecoveryProtectionGroupRuns(),
			"ibm_backup_recovery_protection_policies":      backuprecovery.DataSourceIbmBackupRecoveryProtectionPolicies(),
			"ibm_backup_recovery_protection_policy":        backuprecovery.DataSourceIbmBackupRecoveryProtectionPolicy(),
			"ibm_backup_recovery":                          backuprecovery.DataSourceIbmBackupRecovery(),
			"ibm_backup_recoveries":                        backuprecovery.DataSourceIbmBackupRecoveries(),
			"ibm_backup_recovery_download_files":           backuprecovery.DataSourceIbmBackupRecoveryDownloadFiles(),
			"ibm_backup_recovery_source_registrations":     backuprecovery.DataSourceIbmBackupRecoverySourceRegistrations(),
			"ibm_backup_recovery_source_registration":      backuprecovery.DataSourceIbmBackupRecoverySourceRegistration(),
			"ibm_backup_recovery_download_indexed_files":   backuprecovery.DataSourceIbmBackupRecoveryDownloadIndexedFiles(),
			"ibm_backup_recovery_protection_sources":       backuprecovery.DataSourceIbmBackupRecoveryProtectionSources(),

			// // AppID
			"ibm_appid_action_url":               appid.DataSourceIBMAppIDActionURL(),
			"ibm_appid_apm":                      appid.DataSourceIBMAppIDAPM(),
			"ibm_appid_application":              appid.DataSourceIBMAppIDApplication(),
			"ibm_appid_application_scopes":       appid.DataSourceIBMAppIDApplicationScopes(),
			"ibm_appid_application_roles":        appid.DataSourceIBMAppIDApplicationRoles(),
			"ibm_appid_applications":             appid.DataSourceIBMAppIDApplications(),
			"ibm_appid_audit_status":             appid.DataSourceIBMAppIDAuditStatus(),
			"ibm_appid_cloud_directory_template": appid.DataSourceIBMAppIDCloudDirectoryTemplate(),
			"ibm_appid_cloud_directory_user":     appid.DataSourceIBMAppIDCloudDirectoryUser(),
			"ibm_appid_idp_cloud_directory":      appid.DataSourceIBMAppIDIDPCloudDirectory(),
			"ibm_appid_idp_custom":               appid.DataSourceIBMAppIDIDPCustom(),
			"ibm_appid_idp_facebook":             appid.DataSourceIBMAppIDIDPFacebook(),
			"ibm_appid_idp_google":               appid.DataSourceIBMAppIDIDPGoogle(),
			"ibm_appid_idp_saml":                 appid.DataSourceIBMAppIDIDPSAML(),
			"ibm_appid_idp_saml_metadata":        appid.DataSourceIBMAppIDIDPSAMLMetadata(),
			"ibm_appid_languages":                appid.DataSourceIBMAppIDLanguages(),
			"ibm_appid_mfa":                      appid.DataSourceIBMAppIDMFA(),
			"ibm_appid_mfa_channel":              appid.DataSourceIBMAppIDMFAChannel(),
			"ibm_appid_password_regex":           appid.DataSourceIBMAppIDPasswordRegex(),
			"ibm_appid_token_config":             appid.DataSourceIBMAppIDTokenConfig(),
			"ibm_appid_redirect_urls":            appid.DataSourceIBMAppIDRedirectURLs(),
			"ibm_appid_role":                     appid.DataSourceIBMAppIDRole(),
			"ibm_appid_roles":                    appid.DataSourceIBMAppIDRoles(),
			"ibm_appid_theme_color":              appid.DataSourceIBMAppIDThemeColor(),
			"ibm_appid_theme_text":               appid.DataSourceIBMAppIDThemeText(),
			"ibm_appid_user_roles":               appid.DataSourceIBMAppIDUserRoles(),

			"ibm_function_action":                           functions.DataSourceIBMFunctionAction(),
			"ibm_function_package":                          functions.DataSourceIBMFunctionPackage(),
			"ibm_function_rule":                             functions.DataSourceIBMFunctionRule(),
			"ibm_function_trigger":                          functions.DataSourceIBMFunctionTrigger(),
			"ibm_function_namespace":                        functions.DataSourceIBMFunctionNamespace(),
			"ibm_cis":                                       cis.DataSourceIBMCISInstance(),
			"ibm_cis_dns_records":                           cis.DataSourceIBMCISDNSRecords(),
			"ibm_cis_certificates":                          cis.DataSourceIBMCISCertificates(),
			"ibm_cis_global_load_balancers":                 cis.DataSourceIBMCISGlbs(),
			"ibm_cis_origin_pools":                          cis.DataSourceIBMCISOriginPools(),
			"ibm_cis_healthchecks":                          cis.DataSourceIBMCISHealthChecks(),
			"ibm_cis_domain":                                cis.DataSourceIBMCISDomain(),
			"ibm_cis_firewall":                              cis.DataSourceIBMCISFirewallsRecord(),
			"ibm_cis_cache_settings":                        cis.DataSourceIBMCISCacheSetting(),
			"ibm_cis_waf_packages":                          cis.DataSourceIBMCISWAFPackages(),
			"ibm_cis_range_apps":                            cis.DataSourceIBMCISRangeApps(),
			"ibm_cis_custom_certificates":                   cis.DataSourceIBMCISCustomCertificates(),
			"ibm_cis_rate_limit":                            cis.DataSourceIBMCISRateLimit(),
			"ibm_cis_ip_addresses":                          cis.DataSourceIBMCISIP(),
			"ibm_cis_waf_groups":                            cis.DataSourceIBMCISWAFGroups(),
			"ibm_cis_alerts":                                cis.DataSourceIBMCISAlert(),
			"ibm_cis_origin_auths":                          cis.DataSourceIBMCISOriginAuthPull(),
			"ibm_cis_mtlss":                                 cis.DataSourceIBMCISMtls(),
			"ibm_cis_mtls_apps":                             cis.DataSourceIBMCISMtlsApp(),
			"ibm_cis_bot_managements":                       cis.DataSourceIBMCISBotManagement(),
			"ibm_cis_bot_analytics":                         cis.DataSourceIBMCISBotAnalytics(),
			"ibm_cis_rulesets":                              cis.DataSourceIBMCISRulesets(),
			"ibm_cis_ruleset_versions":                      cis.DataSourceIBMCISRulesetVersions(),
			"ibm_cis_ruleset_rules_by_tag":                  cis.DataSourceIBMCISRulesetRulesByTag(),
			"ibm_cis_ruleset_entrypoint_versions":           cis.DataSourceIBMCISRulesetEntrypointVersions(),
			"ibm_cis_webhooks":                              cis.DataSourceIBMCISWebhooks(),
			"ibm_cis_logpush_jobs":                          cis.DataSourceIBMCISLogPushJobs(),
			"ibm_cis_edge_functions_actions":                cis.DataSourceIBMCISEdgeFunctionsActions(),
			"ibm_cis_edge_functions_triggers":               cis.DataSourceIBMCISEdgeFunctionsTriggers(),
			"ibm_cis_custom_pages":                          cis.DataSourceIBMCISCustomPages(),
			"ibm_cis_page_rules":                            cis.DataSourceIBMCISPageRules(),
			"ibm_cis_waf_rules":                             cis.DataSourceIBMCISWAFRules(),
			"ibm_cis_filters":                               cis.DataSourceIBMCISFilters(),
			"ibm_cis_firewall_rules":                        cis.DataSourceIBMCISFirewallRules(),
			"ibm_cis_origin_certificates":                   cis.DataSourceIBMCISOriginCertificateOrder(),
			"ibm_cloudant":                                  cloudant.DataSourceIBMCloudant(),
			"ibm_cloudant_database":                         cloudant.DataSourceIBMCloudantDatabase(),
			"ibm_database":                                  database.DataSourceIBMDatabaseInstance(),
			"ibm_database_connection":                       database.DataSourceIBMDatabaseConnection(),
			"ibm_database_point_in_time_recovery":           database.DataSourceIBMDatabasePointInTimeRecovery(),
			"ibm_database_remotes":                          database.DataSourceIBMDatabaseRemotes(),
			"ibm_database_task":                             database.DataSourceIBMDatabaseTask(),
			"ibm_database_tasks":                            database.DataSourceIBMDatabaseTasks(),
			"ibm_database_backup":                           database.DataSourceIBMDatabaseBackup(),
			"ibm_database_backups":                          database.DataSourceIBMDatabaseBackups(),
			"ibm_db2":                                       db2.DataSourceIBMDb2Instance(),
			"ibm_db2_connection_info":                       db2.DataSourceIbmDb2ConnectionInfo(),
			"ibm_db2_whitelist_ip":                          db2.DataSourceIbmDb2Whitelist(),
			"ibm_db2_allowlist_ip":                          db2.DataSourceIbmDb2Allowlist(),
			"ibm_db2_autoscale":                             db2.DataSourceIbmDb2Autoscale(),
			"ibm_db2_backup":                                db2.DataSourceIbmDb2Backup(),
			"ibm_db2_tuneable_param":                        db2.DataSourceIbmDb2TuneableParam(),
			"ibm_compute_bare_metal":                        classicinfrastructure.DataSourceIBMComputeBareMetal(),
			"ibm_compute_image_template":                    classicinfrastructure.DataSourceIBMComputeImageTemplate(),
			"ibm_compute_placement_group":                   classicinfrastructure.DataSourceIBMComputePlacementGroup(),
			"ibm_compute_reserved_capacity":                 classicinfrastructure.DataSourceIBMComputeReservedCapacity(),
			"ibm_compute_ssh_key":                           classicinfrastructure.DataSourceIBMComputeSSHKey(),
			"ibm_compute_vm_instance":                       classicinfrastructure.DataSourceIBMComputeVmInstance(),
			"ibm_container_addons":                          kubernetes.DataSourceIBMContainerAddOns(),
			"ibm_container_alb":                             kubernetes.DataSourceIBMContainerALB(),
			"ibm_container_alb_cert":                        kubernetes.DataSourceIBMContainerALBCert(),
			"ibm_container_ingress_instance":                kubernetes.DataSourceIBMContainerIngressInstance(),
			"ibm_container_ingress_secret_tls":              kubernetes.DataSourceIBMContainerIngressSecretTLS(),
			"ibm_container_ingress_secret_opaque":           kubernetes.DataSourceIBMContainerIngressSecretOpaque(),
			"ibm_container_bind_service":                    kubernetes.DataSourceIBMContainerBindService(),
			"ibm_container_cluster":                         kubernetes.DataSourceIBMContainerCluster(),
			"ibm_container_cluster_config":                  kubernetes.DataSourceIBMContainerClusterConfig(),
			"ibm_container_cluster_versions":                kubernetes.DataSourceIBMContainerClusterVersions(),
			"ibm_container_cluster_worker":                  kubernetes.DataSourceIBMContainerClusterWorker(),
			"ibm_container_nlb_dns":                         kubernetes.DataSourceIBMContainerNLBDNS(),
			"ibm_container_vpc_cluster_alb":                 kubernetes.DataSourceIBMContainerVPCClusterALB(),
			"ibm_container_vpc_alb":                         kubernetes.DataSourceIBMContainerVPCClusterALB(),
			"ibm_container_vpc_cluster":                     kubernetes.DataSourceIBMContainerVPCCluster(),
			"ibm_container_vpc_cluster_worker":              kubernetes.DataSourceIBMContainerVPCClusterWorker(),
			"ibm_container_vpc_cluster_worker_pool":         kubernetes.DataSourceIBMContainerVpcClusterWorkerPool(),
			"ibm_container_vpc_worker_pool":                 kubernetes.DataSourceIBMContainerVpcClusterWorkerPool(),
			"ibm_container_worker_pool":                     kubernetes.DataSourceIBMContainerWorkerPool(),
			"ibm_container_storage_attachment":              kubernetes.DataSourceIBMContainerVpcWorkerVolumeAttachment(),
			"ibm_container_dedicated_host_pool":             kubernetes.DataSourceIBMContainerDedicatedHostPool(),
			"ibm_container_dedicated_host_flavor":           kubernetes.DataSourceIBMContainerDedicatedHostFlavor(),
			"ibm_container_dedicated_host_flavors":          kubernetes.DataSourceIBMContainerDedicatedHostFlavors(),
			"ibm_container_dedicated_host":                  kubernetes.DataSourceIBMContainerDedicatedHost(),
			"ibm_cr_namespaces":                             registry.DataIBMContainerRegistryNamespaces(),
			"ibm_cloud_shell_account_settings":              cloudshell.DataSourceIBMCloudShellAccountSettings(),
			"ibm_cos_bucket":                                cos.DataSourceIBMCosBucket(),
			"ibm_cos_backup_vault":                          cos.DataSourceIBMCosBackupVault(),
			"ibm_cos_backup_policy":                         cos.DataSourceIBMCosBackupPolicy(),
			"ibm_cos_bucket_object":                         cos.DataSourceIBMCosBucketObject(),
			"ibm_dns_domain_registration":                   classicinfrastructure.DataSourceIBMDNSDomainRegistration(),
			"ibm_dns_domain":                                classicinfrastructure.DataSourceIBMDNSDomain(),
			"ibm_dns_secondary":                             classicinfrastructure.DataSourceIBMDNSSecondary(),
			"ibm_event_streams_topic":                       eventstreams.DataSourceIBMEventStreamsTopic(),
			"ibm_event_streams_schema":                      eventstreams.DataSourceIBMEventStreamsSchema(),
			"ibm_event_streams_schema_global_rule":          eventstreams.DataSourceIBMEventStreamsSchemaGlobalCompatibilityRule(),
			"ibm_event_streams_quota":                       eventstreams.DataSourceIBMEventStreamsQuota(),
			"ibm_event_streams_mirroring_config":            eventstreams.DataSourceIBMEventStreamsMirroringConfig(),
			"ibm_hpcs":                                      hpcs.DataSourceIBMHPCS(),
			"ibm_hpcs_managed_key":                          hpcs.DataSourceIbmManagedKey(),
			"ibm_hpcs_key_template":                         hpcs.DataSourceIbmKeyTemplate(),
			"ibm_hpcs_keystore":                             hpcs.DataSourceIbmKeystore(),
			"ibm_hpcs_vault":                                hpcs.DataSourceIbmVault(),
			"ibm_iam_access_group":                          iamaccessgroup.DataSourceIBMIAMAccessGroup(),
			"ibm_iam_access_group_policy":                   iampolicy.DataSourceIBMIAMAccessGroupPolicy(),
			"ibm_iam_access_group_template_versions":        iamaccessgroup.DataSourceIBMIAMAccessGroupTemplateVersions(),
			"ibm_iam_access_group_template_assignment":      iamaccessgroup.DataSourceIBMIAMAccessGroupTemplateAssignment(),
			"ibm_iam_account_settings":                      iamidentity.DataSourceIBMIAMAccountSettings(),
			"ibm_iam_effective_account_settings":            iamidentity.DataSourceIBMIamEffectiveAccountSettings(),
			"ibm_iam_auth_token":                            iamidentity.DataSourceIBMIAMAuthToken(),
			"ibm_iam_role_actions":                          iampolicy.DataSourceIBMIAMRoleAction(),
			"ibm_iam_users":                                 iamidentity.DataSourceIBMIAMUsers(),
			"ibm_iam_roles":                                 iampolicy.DataSourceIBMIAMRole(),
			"ibm_iam_user_policy":                           iampolicy.DataSourceIBMIAMUserPolicy(),
			"ibm_iam_authorization_policies":                iampolicy.DataSourceIBMIAMAuthorizationPolicies(),
			"ibm_iam_user_profile":                          iamidentity.DataSourceIBMIAMUserProfile(),
			"ibm_iam_service_id":                            iamidentity.DataSourceIBMIAMServiceID(),
			"ibm_iam_service_policy":                        iampolicy.DataSourceIBMIAMServicePolicy(),
			"ibm_iam_api_key":                               iamidentity.DataSourceIBMIamApiKey(),
			"ibm_iam_trusted_profile":                       iamidentity.DataSourceIBMIamTrustedProfile(),
			"ibm_iam_trusted_profile_identity":              iamidentity.DataSourceIBMIamTrustedProfileIdentity(),
			"ibm_iam_trusted_profile_identities":            iamidentity.DataSourceIBMIamTrustedProfileIdentities(),
			"ibm_iam_trusted_profile_claim_rule":            iamidentity.DataSourceIBMIamTrustedProfileClaimRule(),
			"ibm_iam_trusted_profile_link":                  iamidentity.DataSourceIBMIamTrustedProfileLink(),
			"ibm_iam_trusted_profile_claim_rules":           iamidentity.DataSourceIBMIamTrustedProfileClaimRules(),
			"ibm_iam_trusted_profile_links":                 iamidentity.DataSourceIBMIamTrustedProfileLinks(),
			"ibm_iam_trusted_profiles":                      iamidentity.DataSourceIBMIamTrustedProfiles(),
			"ibm_iam_trusted_profile_policy":                iampolicy.DataSourceIBMIAMTrustedProfilePolicy(),
			"ibm_iam_user_mfa_enrollments":                  iamidentity.DataSourceIBMIamUserMfaEnrollments(),
			"ibm_iam_account_settings_template":             iamidentity.DataSourceIBMAccountSettingsTemplate(),
			"ibm_iam_trusted_profile_template":              iamidentity.DataSourceIBMTrustedProfileTemplate(),
			"ibm_iam_account_settings_template_assignment":  iamidentity.DataSourceIBMAccountSettingsTemplateAssignment(),
			"ibm_iam_trusted_profile_template_assignment":   iamidentity.DataSourceIBMTrustedProfileTemplateAssignment(),
			"ibm_iam_policy_template":                       iampolicy.DataSourceIBMIAMPolicyTemplate(),
			"ibm_iam_policy_template_version":               iampolicy.DataSourceIBMIAMPolicyTemplateVersion(),
			"ibm_iam_policy_assignments":                    iampolicy.DataSourceIBMIAMPolicyAssignments(),
			"ibm_iam_policy_assignment":                     iampolicy.DataSourceIBMIAMPolicyAssignment(),
			"ibm_iam_account_settings_external_interaction": iampolicy.DataSourceIBMIAMAccountSettingsExternalInteraction(),

			// backup as Service
			"ibm_is_backup_policy":       vpc.DataSourceIBMIsBackupPolicy(),
			"ibm_is_backup_policies":     vpc.DataSourceIBMIsBackupPolicies(),
			"ibm_is_backup_policy_plan":  vpc.DataSourceIBMIsBackupPolicyPlan(),
			"ibm_is_backup_policy_plans": vpc.DataSourceIBMIsBackupPolicyPlans(),
			"ibm_is_backup_policy_job":   vpc.DataSourceIBMIsBackupPolicyJob(),
			"ibm_is_backup_policy_jobs":  vpc.DataSourceIBMIsBackupPolicyJobs(),

			// bare_metal_server
			"ibm_is_bare_metal_server_disk":                           vpc.DataSourceIBMIsBareMetalServerDisk(),
			"ibm_is_bare_metal_server_disks":                          vpc.DataSourceIBMIsBareMetalServerDisks(),
			"ibm_is_bare_metal_server_initialization":                 vpc.DataSourceIBMIsBareMetalServerInitialization(),
			"ibm_is_bare_metal_server_network_attachment":             vpc.DataSourceIBMIsBareMetalServerNetworkAttachment(),
			"ibm_is_bare_metal_server_network_attachments":            vpc.DataSourceIBMIsBareMetalServerNetworkAttachments(),
			"ibm_is_bare_metal_server_network_interface_floating_ip":  vpc.DataSourceIBMIsBareMetalServerNetworkInterfaceFloatingIP(),
			"ibm_is_bare_metal_server_network_interface_floating_ips": vpc.DataSourceIBMIsBareMetalServerNetworkInterfaceFloatingIPs(),
			"ibm_is_bare_metal_server_network_interface_reserved_ip":  vpc.DataSourceIBMISBareMetalServerNICReservedIP(),
			"ibm_is_bare_metal_server_network_interface_reserved_ips": vpc.DataSourceIBMISBareMetalServerNICReservedIPs(),
			"ibm_is_bare_metal_server_network_interface":              vpc.DataSourceIBMIsBareMetalServerNetworkInterface(),
			"ibm_is_bare_metal_server_network_interfaces":             vpc.DataSourceIBMIsBareMetalServerNetworkInterfaces(),
			"ibm_is_bare_metal_server_profile":                        vpc.DataSourceIBMIsBareMetalServerProfile(),
			"ibm_is_bare_metal_server_profiles":                       vpc.DataSourceIBMIsBareMetalServerProfiles(),
			"ibm_is_bare_metal_server":                                vpc.DataSourceIBMIsBareMetalServer(),
			"ibm_is_bare_metal_servers":                               vpc.DataSourceIBMIsBareMetalServers(),

			// cluster
			"ibm_is_cluster_network":                      vpc.DataSourceIBMIsClusterNetwork(),
			"ibm_is_cluster_networks":                     vpc.DataSourceIBMIsClusterNetworks(),
			"ibm_is_cluster_network_interface":            vpc.DataSourceIBMIsClusterNetworkInterface(),
			"ibm_is_cluster_network_interfaces":           vpc.DataSourceIBMIsClusterNetworkInterfaces(),
			"ibm_is_cluster_network_profile":              vpc.DataSourceIBMIsClusterNetworkProfile(),
			"ibm_is_cluster_network_profiles":             vpc.DataSourceIBMIsClusterNetworkProfiles(),
			"ibm_is_cluster_network_subnet":               vpc.DataSourceIBMIsClusterNetworkSubnet(),
			"ibm_is_cluster_network_subnets":              vpc.DataSourceIBMIsClusterNetworkSubnets(),
			"ibm_is_cluster_network_subnet_reserved_ip":   vpc.DataSourceIBMIsClusterNetworkSubnetReservedIP(),
			"ibm_is_cluster_network_subnet_reserved_ips":  vpc.DataSourceIBMIsClusterNetworkSubnetReservedIps(),
			"ibm_is_instance_cluster_network_attachment":  vpc.DataSourceIBMIsInstanceClusterNetworkAttachment(),
			"ibm_is_instance_cluster_network_attachments": vpc.DataSourceIBMIsInstanceClusterNetworkAttachments(),

			"ibm_is_dedicated_host":                  vpc.DataSourceIbmIsDedicatedHost(),
			"ibm_is_dedicated_hosts":                 vpc.DataSourceIbmIsDedicatedHosts(),
			"ibm_is_dedicated_host_profile":          vpc.DataSourceIbmIsDedicatedHostProfile(),
			"ibm_is_dedicated_host_profiles":         vpc.DataSourceIbmIsDedicatedHostProfiles(),
			"ibm_is_dedicated_host_group":            vpc.DataSourceIbmIsDedicatedHostGroup(),
			"ibm_is_dedicated_host_groups":           vpc.DataSourceIbmIsDedicatedHostGroups(),
			"ibm_is_dedicated_host_disk":             vpc.DataSourceIbmIsDedicatedHostDisk(),
			"ibm_is_dedicated_host_disks":            vpc.DataSourceIbmIsDedicatedHostDisks(),
			"ibm_is_placement_group":                 vpc.DataSourceIbmIsPlacementGroup(),
			"ibm_is_placement_groups":                vpc.DataSourceIbmIsPlacementGroups(),
			"ibm_is_floating_ip":                     vpc.DataSourceIBMISFloatingIP(),
			"ibm_is_floating_ips":                    vpc.DataSourceIBMIsFloatingIps(),
			"ibm_is_flow_log":                        vpc.DataSourceIBMIsFlowLog(),
			"ibm_is_flow_logs":                       vpc.DataSourceIBMISFlowLogs(),
			"ibm_is_image":                           vpc.DataSourceIBMISImage(),
			"ibm_is_images":                          vpc.DataSourceIBMISImages(),
			"ibm_is_image_export_job":                vpc.DataSourceIBMIsImageExport(),
			"ibm_is_image_export_jobs":               vpc.DataSourceIBMIsImageExports(),
			"ibm_is_endpoint_gateway_targets":        vpc.DataSourceIBMISEndpointGatewayTargets(),
			"ibm_is_instance_group":                  vpc.DataSourceIBMISInstanceGroup(),
			"ibm_is_instance_groups":                 vpc.DataSourceIBMISInstanceGroups(),
			"ibm_is_instance_group_memberships":      vpc.DataSourceIBMISInstanceGroupMemberships(),
			"ibm_is_instance_group_membership":       vpc.DataSourceIBMISInstanceGroupMembership(),
			"ibm_is_instance_group_manager":          vpc.DataSourceIBMISInstanceGroupManager(),
			"ibm_is_instance_group_managers":         vpc.DataSourceIBMISInstanceGroupManagers(),
			"ibm_is_instance_group_manager_policies": vpc.DataSourceIBMISInstanceGroupManagerPolicies(),
			"ibm_is_instance_group_manager_policy":   vpc.DataSourceIBMISInstanceGroupManagerPolicy(),
			"ibm_is_instance_group_manager_action":   vpc.DataSourceIBMISInstanceGroupManagerAction(),
			"ibm_is_instance_group_manager_actions":  vpc.DataSourceIBMISInstanceGroupManagerActions(),
			"ibm_is_virtual_endpoint_gateways":       vpc.DataSourceIBMISEndpointGateways(),
			"ibm_is_virtual_endpoint_gateway_ips":    vpc.DataSourceIBMISEndpointGatewayIPs(),
			"ibm_is_virtual_endpoint_gateway":        vpc.DataSourceIBMISEndpointGateway(),
			"ibm_is_instance_template":               vpc.DataSourceIBMISInstanceTemplate(),
			"ibm_is_instance_templates":              vpc.DataSourceIBMISInstanceTemplates(),
			"ibm_is_instance_profile":                vpc.DataSourceIBMISInstanceProfile(),
			"ibm_is_instance_profiles":               vpc.DataSourceIBMISInstanceProfiles(),
			"ibm_is_instance":                        vpc.DataSourceIBMISInstance(),
			"ibm_is_instances":                       vpc.DataSourceIBMISInstances(),
			"ibm_is_instance_network_attachment":     vpc.DataSourceIBMIsInstanceNetworkAttachment(),
			"ibm_is_instance_network_attachments":    vpc.DataSourceIBMIsInstanceNetworkAttachments(),
			"ibm_is_instance_network_interface":      vpc.DataSourceIBMIsInstanceNetworkInterface(),
			"ibm_is_instance_network_interfaces":     vpc.DataSourceIBMIsInstanceNetworkInterfaces(),
			"ibm_is_instance_disk":                   vpc.DataSourceIbmIsInstanceDisk(),
			"ibm_is_instance_disks":                  vpc.DataSourceIbmIsInstanceDisks(),

			// reserved ips
			"ibm_is_instance_network_interface_reserved_ip":  vpc.DataSourceIBMISInstanceNICReservedIP(),
			"ibm_is_instance_network_interface_reserved_ips": vpc.DataSourceIBMISInstanceNICReservedIPs(),

			"ibm_is_instance_volume_attachment":                    vpc.DataSourceIBMISInstanceVolumeAttachment(),
			"ibm_is_instance_volume_attachments":                   vpc.DataSourceIBMISInstanceVolumeAttachments(),
			"ibm_is_ipsec_policy":                                  vpc.DataSourceIBMIsIpsecPolicy(),
			"ibm_is_ipsec_policies":                                vpc.DataSourceIBMIsIpsecPolicies(),
			"ibm_is_ike_policies":                                  vpc.DataSourceIBMIsIkePolicies(),
			"ibm_is_ike_policy":                                    vpc.DataSourceIBMIsIkePolicy(),
			"ibm_is_lb":                                            vpc.DataSourceIBMISLB(),
			"ibm_is_lb_listener":                                   vpc.DataSourceIBMISLBListener(),
			"ibm_is_lb_listeners":                                  vpc.DataSourceIBMISLBListeners(),
			"ibm_is_lb_listener_policies":                          vpc.DataSourceIBMISLBListenerPolicies(),
			"ibm_is_lb_listener_policy":                            vpc.DataSourceIBMISLBListenerPolicy(),
			"ibm_is_lb_listener_policy_rule":                       vpc.DataSourceIBMISLBListenerPolicyRule(),
			"ibm_is_lb_listener_policy_rules":                      vpc.DataSourceIBMISLBListenerPolicyRules(),
			"ibm_is_lb_pool":                                       vpc.DataSourceIBMISLBPool(),
			"ibm_is_lb_pools":                                      vpc.DataSourceIBMISLBPools(),
			"ibm_is_lb_pool_member":                                vpc.DataSourceIBMIBLBPoolMember(),
			"ibm_is_lb_pool_members":                               vpc.DataSourceIBMISLBPoolMembers(),
			"ibm_is_lb_profile":                                    vpc.DataSourceIBMISLbProfile(),
			"ibm_is_lb_profiles":                                   vpc.DataSourceIBMISLbProfiles(),
			"ibm_is_lbs":                                           vpc.DataSourceIBMISLBS(),
			"ibm_is_private_path_service_gateway":                  vpc.DataSourceIBMIsPrivatePathServiceGateway(),
			"ibm_is_private_path_service_gateway_account_policy":   vpc.DataSourceIBMIsPrivatePathServiceGatewayAccountPolicy(),
			"ibm_is_private_path_service_gateway_account_policies": vpc.DataSourceIBMIsPrivatePathServiceGatewayAccountPolicies(),
			"ibm_is_private_path_service_gateways":                 vpc.DataSourceIBMIsPrivatePathServiceGateways(),
			"ibm_is_private_path_service_gateway_endpoint_gateway_binding":  vpc.DataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBinding(),
			"ibm_is_private_path_service_gateway_endpoint_gateway_bindings": vpc.DataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindings(),
			"ibm_is_public_gateway":              vpc.DataSourceIBMISPublicGateway(),
			"ibm_is_public_gateways":             vpc.DataSourceIBMISPublicGateways(),
			"ibm_is_region":                      vpc.DataSourceIBMISRegion(),
			"ibm_is_regions":                     vpc.DataSourceIBMISRegions(),
			"ibm_is_reservation":                 vpc.DataSourceIBMIsReservation(),
			"ibm_is_reservations":                vpc.DataSourceIBMIsReservations(),
			"ibm_is_ssh_key":                     vpc.DataSourceIBMISSSHKey(),
			"ibm_is_ssh_keys":                    vpc.DataSourceIBMIsSshKeys(),
			"ibm_is_subnet":                      vpc.DataSourceIBMISSubnet(),
			"ibm_is_subnets":                     vpc.DataSourceIBMISSubnets(),
			"ibm_is_subnet_reserved_ip":          vpc.DataSourceIBMISReservedIP(),
			"ibm_is_subnet_reserved_ips":         vpc.DataSourceIBMISReservedIPs(),
			"ibm_is_security_group":              vpc.DataSourceIBMISSecurityGroup(),
			"ibm_is_security_groups":             vpc.DataSourceIBMIsSecurityGroups(),
			"ibm_is_security_group_rule":         vpc.DataSourceIBMIsSecurityGroupRule(),
			"ibm_is_security_group_rules":        vpc.DataSourceIBMIsSecurityGroupRules(),
			"ibm_is_security_group_target":       vpc.DataSourceIBMISSecurityGroupTarget(),
			"ibm_is_security_group_targets":      vpc.DataSourceIBMISSecurityGroupTargets(),
			"ibm_is_snapshot_clone":              vpc.DataSourceSnapshotClone(),
			"ibm_is_snapshot_clones":             vpc.DataSourceSnapshotClones(),
			"ibm_is_snapshot":                    vpc.DataSourceSnapshot(),
			"ibm_is_snapshot_consistency_group":  vpc.DataSourceIBMIsSnapshotConsistencyGroup(),
			"ibm_is_snapshot_consistency_groups": vpc.DataSourceIBMIsSnapshotConsistencyGroups(),
			"ibm_is_snapshots":                   vpc.DataSourceSnapshots(),
			"ibm_is_share":                       vpc.DataSourceIbmIsShare(),
			"ibm_is_source_share":                vpc.DataSourceIbmIsSourceShare(),
			"ibm_is_shares":                      vpc.DataSourceIbmIsShares(),
			"ibm_is_share_profile":               vpc.DataSourceIbmIsShareProfile(),
			"ibm_is_share_profiles":              vpc.DataSourceIbmIsShareProfiles(),
			"ibm_is_share_accessor_bindings":     vpc.DataSourceIBMIsShareAccessorBindings(),
			"ibm_is_share_accessor_binding":      vpc.DataSourceIBMIsShareAccessorBinding(),
			"ibm_is_share_snapshot":              vpc.DataSourceIBMIsShareSnapshot(),
			"ibm_is_share_snapshots":             vpc.DataSourceIBMIsShareSnapshots(),
			"ibm_is_virtual_network_interface":   vpc.DataSourceIBMIsVirtualNetworkInterface(),
			"ibm_is_virtual_network_interfaces":  vpc.DataSourceIBMIsVirtualNetworkInterfaces(),

			// vni

			"ibm_is_virtual_network_interface_floating_ip":  vpc.DataSourceIBMIsVirtualNetworkInterfaceFloatingIP(),
			"ibm_is_virtual_network_interface_floating_ips": vpc.DataSourceIBMIsVirtualNetworkInterfaceFloatingIPs(),
			"ibm_is_virtual_network_interface_ip":           vpc.DataSourceIBMIsVirtualNetworkInterfaceIP(),
			"ibm_is_virtual_network_interface_ips":          vpc.DataSourceIBMIsVirtualNetworkInterfaceIPs(),

			"ibm_is_share_mount_target":          vpc.DataSourceIBMIsShareTarget(),
			"ibm_is_share_mount_targets":         vpc.DataSourceIBMIsShareTargets(),
			"ibm_is_volume":                      vpc.DataSourceIBMISVolume(),
			"ibm_is_volumes":                     vpc.DataSourceIBMIsVolumes(),
			"ibm_is_volume_profile":              vpc.DataSourceIBMISVolumeProfile(),
			"ibm_is_volume_profiles":             vpc.DataSourceIBMISVolumeProfiles(),
			"ibm_is_vpc":                         vpc.DataSourceIBMISVPC(),
			"ibm_is_vpc_dns_resolution_binding":  vpc.DataSourceIBMIsVPCDnsResolutionBinding(),
			"ibm_is_vpc_dns_resolution_bindings": vpc.DataSourceIBMIsVPCDnsResolutionBindings(),
			"ibm_is_vpcs":                        vpc.DataSourceIBMISVPCs(),
			"ibm_is_vpn_gateway":                 vpc.DataSourceIBMISVPNGateway(),
			"ibm_is_vpn_gateways":                vpc.DataSourceIBMISVPNGateways(),
			"ibm_is_vpc_address_prefixes":        vpc.DataSourceIbmIsVpcAddressPrefixes(),
			"ibm_is_vpc_address_prefix":          vpc.DataSourceIBMIsVPCAddressPrefix(),
			"ibm_is_vpn_gateway_connection":      vpc.DataSourceIBMISVPNGatewayConnection(),
			"ibm_is_vpn_gateway_connections":     vpc.DataSourceIBMISVPNGatewayConnections(),

			"ibm_is_vpn_gateway_connection_local_cidrs": vpc.DataSourceIBMIsVPNGatewayConnectionLocalCidrs(),
			"ibm_is_vpn_gateway_connection_peer_cidrs":  vpc.DataSourceIBMIsVPNGatewayConnectionPeerCidrs(),

			"ibm_is_vpc_default_routing_table":       vpc.DataSourceIBMISVPCDefaultRoutingTable(),
			"ibm_is_vpc_routing_table":               vpc.DataSourceIBMIBMIsVPCRoutingTable(),
			"ibm_is_vpc_routing_tables":              vpc.DataSourceIBMISVPCRoutingTables(),
			"ibm_is_vpc_routing_table_route":         vpc.DataSourceIBMIBMIsVPCRoutingTableRoute(),
			"ibm_is_vpc_routing_table_routes":        vpc.DataSourceIBMISVPCRoutingTableRoutes(),
			"ibm_is_vpn_server":                      vpc.DataSourceIBMIsVPNServer(),
			"ibm_is_vpn_servers":                     vpc.DataSourceIBMIsVPNServers(),
			"ibm_is_vpn_server_client":               vpc.DataSourceIBMIsVPNServerClient(),
			"ibm_is_vpn_server_client_configuration": vpc.DataSourceIBMIsVPNServerClientConfiguration(),
			"ibm_is_vpn_server_clients":              vpc.DataSourceIBMIsVPNServerClients(),
			"ibm_is_vpn_server_route":                vpc.DataSourceIBMIsVPNServerRoute(),
			"ibm_is_vpn_server_routes":               vpc.DataSourceIBMIsVPNServerRoutes(),
			"ibm_is_zone":                            vpc.DataSourceIBMISZone(),
			"ibm_is_zones":                           vpc.DataSourceIBMISZones(),
			"ibm_is_operating_system":                vpc.DataSourceIBMISOperatingSystem(),
			"ibm_is_operating_systems":               vpc.DataSourceIBMISOperatingSystems(),
			"ibm_is_network_acls":                    vpc.DataSourceIBMIsNetworkAcls(),
			"ibm_is_network_acl":                     vpc.DataSourceIBMIsNetworkACL(),
			"ibm_is_network_acl_rule":                vpc.DataSourceIBMISNetworkACLRule(),
			"ibm_is_network_acl_rules":               vpc.DataSourceIBMISNetworkACLRules(),
			"ibm_lbaas":                              classicinfrastructure.DataSourceIBMLbaas(),
			"ibm_network_vlan":                       classicinfrastructure.DataSourceIBMNetworkVlan(),
			"ibm_org":                                cloudfoundry.DataSourceIBMOrg(),
			"ibm_org_quota":                          cloudfoundry.DataSourceIBMOrgQuota(),
			"ibm_kms_instance_policies":              kms.DataSourceIBMKmsInstancePolicies(),
			"ibm_kp_key":                             kms.DataSourceIBMkey(),
			"ibm_kms_key_rings":                      kms.DataSourceIBMKMSkeyRings(),
			"ibm_kms_key_policies":                   kms.DataSourceIBMKMSkeyPolicies(),
			"ibm_kms_keys":                           kms.DataSourceIBMKMSkeys(),
			"ibm_kms_key":                            kms.DataSourceIBMKMSkey(),
			"ibm_kms_kmip_adapter":                   kms.DataSourceIBMKMSKmipAdapter(),
			"ibm_kms_kmip_adapters":                  kms.DataSourceIBMKMSKmipAdapters(),
			"ibm_kms_kmip_client_cert":               kms.DataSourceIBMKmsKMIPClientCertificate(),
			"ibm_kms_kmip_client_certs":              kms.DataSourceIBMKmsKMIPClientCertificates(),
			"ibm_kms_kmip_object":                    kms.DataSourceIBMKMSKMIPObject(),
			"ibm_kms_kmip_objects":                   kms.DataSourceIBMKMSKMIPObjects(),
			"ibm_pn_application_chrome":              pushnotification.DataSourceIBMPNApplicationChrome(),
			"ibm_app_config_environment":             appconfiguration.DataSourceIBMAppConfigEnvironment(),
			"ibm_app_config_environments":            appconfiguration.DataSourceIBMAppConfigEnvironments(),
			"ibm_app_config_collection":              appconfiguration.DataSourceIBMAppConfigCollection(),
			"ibm_app_config_collections":             appconfiguration.DataSourceIBMAppConfigCollections(),
			"ibm_app_config_feature":                 appconfiguration.DataSourceIBMAppConfigFeature(),
			"ibm_app_config_features":                appconfiguration.DataSourceIBMAppConfigFeatures(),
			"ibm_app_config_property":                appconfiguration.DataSourceIBMAppConfigProperty(),
			"ibm_app_config_properties":              appconfiguration.DataSourceIBMAppConfigProperties(),
			"ibm_app_config_segment":                 appconfiguration.DataSourceIBMAppConfigSegment(),
			"ibm_app_config_segments":                appconfiguration.DataSourceIBMAppConfigSegments(),
			"ibm_app_config_snapshot":                appconfiguration.DataSourceIBMAppConfigSnapshot(),
			"ibm_app_config_snapshots":               appconfiguration.DataSourceIBMAppConfigSnapshots(),

			"ibm_resource_quota":    resourcecontroller.DataSourceIBMResourceQuota(),
			"ibm_resource_group":    resourcemanager.DataSourceIBMResourceGroup(),
			"ibm_resource_instance": resourcecontroller.DataSourceIBMResourceInstance(),
			"ibm_resource_key":      resourcecontroller.DataSourceIBMResourceKey(),
			"ibm_security_group":    classicinfrastructure.DataSourceIBMSecurityGroup(),
			"ibm_service_instance":  cloudfoundry.DataSourceIBMServiceInstance(),
			"ibm_service_key":       cloudfoundry.DataSourceIBMServiceKey(),
			"ibm_service_plan":      cloudfoundry.DataSourceIBMServicePlan(),
			"ibm_space":             cloudfoundry.DataSourceIBMSpace(),

			// Added for Schematics
			"ibm_schematics_workspace":      schematics.DataSourceIBMSchematicsWorkspace(),
			"ibm_schematics_output":         schematics.DataSourceIBMSchematicsOutput(),
			"ibm_schematics_state":          schematics.DataSourceIBMSchematicsState(),
			"ibm_schematics_action":         schematics.DataSourceIBMSchematicsAction(),
			"ibm_schematics_job":            schematics.DataSourceIBMSchematicsJob(),
			"ibm_schematics_inventory":      schematics.DataSourceIBMSchematicsInventory(),
			"ibm_schematics_resource_query": schematics.DataSourceIBMSchematicsResourceQuery(),
			"ibm_schematics_policies":       schematics.DataSourceIbmSchematicsPolicies(),
			"ibm_schematics_policy":         schematics.DataSourceIbmSchematicsPolicy(),
			"ibm_schematics_agents":         schematics.DataSourceIbmSchematicsAgents(),
			"ibm_schematics_agent":          schematics.DataSourceIbmSchematicsAgent(),
			"ibm_schematics_agent_prs":      schematics.DataSourceIbmSchematicsAgentPrs(),
			"ibm_schematics_agent_deploy":   schematics.DataSourceIbmSchematicsAgentDeploy(),
			"ibm_schematics_agent_health":   schematics.DataSourceIbmSchematicsAgentHealth(),

			// Added for Power Resources
			"ibm_pi_available_hosts":                        power.DataSourceIBMPIAvailableHosts(),
			"ibm_pi_catalog_images":                         power.DataSourceIBMPICatalogImages(),
			"ibm_pi_cloud_connection":                       power.DataSourceIBMPICloudConnection(),
			"ibm_pi_cloud_connections":                      power.DataSourceIBMPICloudConnections(),
			"ibm_pi_cloud_instance":                         power.DataSourceIBMPICloudInstance(),
			"ibm_pi_console_languages":                      power.DataSourceIBMPIInstanceConsoleLanguages(),
			"ibm_pi_datacenter":                             power.DataSourceIBMPIDatacenter(),
			"ibm_pi_datacenters":                            power.DataSourceIBMPIDatacenters(),
			"ibm_pi_dhcp":                                   power.DataSourceIBMPIDhcp(),
			"ibm_pi_dhcps":                                  power.DataSourceIBMPIDhcps(),
			"ibm_pi_disaster_recovery_location":             power.DataSourceIBMPIDisasterRecoveryLocation(),
			"ibm_pi_disaster_recovery_locations":            power.DataSourceIBMPIDisasterRecoveryLocations(),
			"ibm_pi_host_group":                             power.DataSourceIBMPIHostGroup(),
			"ibm_pi_host_groups":                            power.DataSourceIBMPIHostGroups(),
			"ibm_pi_host":                                   power.DataSourceIBMPIHost(),
			"ibm_pi_hosts":                                  power.DataSourceIBMPIHosts(),
			"ibm_pi_image":                                  power.DataSourceIBMPIImage(),
			"ibm_pi_images":                                 power.DataSourceIBMPIImages(),
			"ibm_pi_instance_ip":                            power.DataSourceIBMPIInstanceIP(),
			"ibm_pi_instance_snapshot":                      power.DataSourceIBMPIInstanceSnapshot(),
			"ibm_pi_instance_snapshots":                     power.DataSourceIBMPIInstanceSnapshots(),
			"ibm_pi_instance_volumes":                       power.DataSourceIBMPIInstanceVolumes(),
			"ibm_pi_instance":                               power.DataSourceIBMPIInstance(),
			"ibm_pi_instances":                              power.DataSourceIBMPIInstances(),
			"ibm_pi_key":                                    power.DataSourceIBMPIKey(),
			"ibm_pi_keys":                                   power.DataSourceIBMPIKeys(),
			"ibm_pi_network_address_group":                  power.DataSourceIBMPINetworkAddressGroup(),
			"ibm_pi_network_address_groups":                 power.DataSourceIBMPINetworkAddressGroups(),
			"ibm_pi_network_interface":                      power.DataSourceIBMPINetworkInterface(),
			"ibm_pi_network_interfaces":                     power.DataSourceIBMPINetworkInterfaces(),
			"ibm_pi_network_peers":                          power.DataSourceIBMPINetworkPeers(),
			"ibm_pi_network_port":                           power.DataSourceIBMPINetworkPort(),
			"ibm_pi_network_security_group":                 power.DataSourceIBMPINetworkSecurityGroup(),
			"ibm_pi_network_security_groups":                power.DataSourceIBMPINetworkSecurityGroups(),
			"ibm_pi_network":                                power.DataSourceIBMPINetwork(),
			"ibm_pi_networks":                               power.DataSourceIBMPINetworks(),
			"ibm_pi_placement_group":                        power.DataSourceIBMPIPlacementGroup(),
			"ibm_pi_placement_groups":                       power.DataSourceIBMPIPlacementGroups(),
			"ibm_pi_public_network":                         power.DataSourceIBMPIPublicNetwork(),
			"ibm_pi_pvm_snapshots":                          power.DataSourceIBMPIPVMSnapshot(),
			"ibm_pi_sap_profile":                            power.DataSourceIBMPISAPProfile(),
			"ibm_pi_sap_profiles":                           power.DataSourceIBMPISAPProfiles(),
			"ibm_pi_shared_processor_pool":                  power.DataSourceIBMPISharedProcessorPool(),
			"ibm_pi_shared_processor_pools":                 power.DataSourceIBMPISharedProcessorPools(),
			"ibm_pi_spp_placement_group":                    power.DataSourceIBMPISPPPlacementGroup(),
			"ibm_pi_spp_placement_groups":                   power.DataSourceIBMPISPPPlacementGroups(),
			"ibm_pi_storage_pool_capacity":                  power.DataSourceIBMPIStoragePoolCapacity(),
			"ibm_pi_storage_pools_capacity":                 power.DataSourceIBMPIStoragePoolsCapacity(),
			"ibm_pi_storage_tiers":                          power.DataSourceIBMPIStorageTiers(),
			"ibm_pi_storage_type_capacity":                  power.DataSourceIBMPIStorageTypeCapacity(),
			"ibm_pi_storage_types_capacity":                 power.DataSourceIBMPIStorageTypesCapacity(),
			"ibm_pi_system_pools":                           power.DataSourceIBMPISystemPools(),
			"ibm_pi_tenant":                                 power.DataSourceIBMPITenant(),
			"ibm_pi_virtual_serial_number":                  power.DataSourceIBMPIVirtualSerialNumber(),
			"ibm_pi_virtual_serial_numbers":                 power.DataSourceIBMPIVirtualSerialNumbers(),
			"ibm_pi_volume_clone":                           power.DataSourceIBMPIVolumeClone(),
			"ibm_pi_volume_flash_copy_mappings":             power.DataSourceIBMPIVolumeFlashCopyMappings(),
			"ibm_pi_volume_group_details":                   power.DataSourceIBMPIVolumeGroupDetails(),
			"ibm_pi_volume_group_remote_copy_relationships": power.DataSourceIBMPIVolumeGroupRemoteCopyRelationships(),
			"ibm_pi_volume_group_storage_details":           power.DataSourceIBMPIVolumeGroupStorageDetails(),
			"ibm_pi_volume_group":                           power.DataSourceIBMPIVolumeGroup(),
			"ibm_pi_volume_groups_details":                  power.DataSourceIBMPIVolumeGroupsDetails(),
			"ibm_pi_volume_groups":                          power.DataSourceIBMPIVolumeGroups(),
			"ibm_pi_volume_onboarding":                      power.DataSourceIBMPIVolumeOnboarding(),
			"ibm_pi_volume_onboardings":                     power.DataSourceIBMPIVolumeOnboardings(),
			"ibm_pi_volume_remote_copy_relationship":        power.DataSourceIBMPIVolumeRemoteCopyRelationship(),
			"ibm_pi_volume_snapshot":                        power.DataSourceIBMPIVolumeSnapshot(),
			"ibm_pi_volume_snapshots":                       power.DataSourceIBMPIVolumeSnapshots(),
			"ibm_pi_volume":                                 power.DataSourceIBMPIVolume(),
			"ibm_pi_workspace":                              power.DatasourceIBMPIWorkspace(),
			"ibm_pi_workspaces":                             power.DatasourceIBMPIWorkspaces(),

			// Added for private dns zones

			"ibm_dns_zones":                            dnsservices.DataSourceIBMPrivateDNSZones(),
			"ibm_dns_permitted_networks":               dnsservices.DataSourceIBMPrivateDNSPermittedNetworks(),
			"ibm_dns_resource_records":                 dnsservices.DataSourceIBMPrivateDNSResourceRecords(),
			"ibm_dns_glb_monitors":                     dnsservices.DataSourceIBMPrivateDNSGLBMonitors(),
			"ibm_dns_glb_pools":                        dnsservices.DataSourceIBMPrivateDNSGLBPools(),
			"ibm_dns_glbs":                             dnsservices.DataSourceIBMPrivateDNSGLBs(),
			"ibm_dns_custom_resolvers":                 dnsservices.DataSourceIBMPrivateDNSCustomResolver(),
			"ibm_dns_custom_resolver_forwarding_rules": dnsservices.DataSourceIBMPrivateDNSForwardingRules(),
			"ibm_dns_custom_resolver_secondary_zones":  dnsservices.DataSourceIBMPrivateDNSSecondaryZones(),

			// Added for Direct Link

			"ibm_dl_gateways":             directlink.DataSourceIBMDLGateways(),
			"ibm_dl_offering_speeds":      directlink.DataSourceIBMDLOfferingSpeeds(),
			"ibm_dl_port":                 directlink.DataSourceIBMDirectLinkPort(),
			"ibm_dl_ports":                directlink.DataSourceIBMDirectLinkPorts(),
			"ibm_dl_gateway":              directlink.DataSourceIBMDLGateway(),
			"ibm_dl_locations":            directlink.DataSourceIBMDLLocations(),
			"ibm_dl_routers":              directlink.DataSourceIBMDLRouters(),
			"ibm_dl_provider_ports":       directlink.DataSourceIBMDirectLinkProviderPorts(),
			"ibm_dl_provider_gateways":    directlink.DataSourceIBMDirectLinkProviderGateways(),
			"ibm_dl_route_reports":        directlink.DataSourceIBMDLRouteReports(),
			"ibm_dl_route_report":         directlink.DataSourceIBMDLRouteReport(),
			"ibm_dl_export_route_filters": directlink.DataSourceIBMDLExportRouteFilters(),
			"ibm_dl_export_route_filter":  directlink.DataSourceIBMDLExportRouteFilter(),
			"ibm_dl_import_route_filters": directlink.DataSourceIBMDLImportRouteFilters(),
			"ibm_dl_import_route_filter":  directlink.DataSourceIBMDLImportRouteFilter(),

			// Added for Transit Gateway
			"ibm_tg_gateway":                   transitgateway.DataSourceIBMTransitGateway(),
			"ibm_tg_gateways":                  transitgateway.DataSourceIBMTransitGateways(),
			"ibm_tg_connection_prefix_filter":  transitgateway.DataSourceIBMTransitGatewayConnectionPrefixFilter(),
			"ibm_tg_connection_prefix_filters": transitgateway.DataSourceIBMTransitGatewayConnectionPrefixFilters(),
			"ibm_tg_locations":                 transitgateway.DataSourceIBMTransitGatewaysLocations(),
			"ibm_tg_location":                  transitgateway.DataSourceIBMTransitGatewaysLocation(),
			"ibm_tg_route_report":              transitgateway.DataSourceIBMTransitGatewayRouteReport(),
			"ibm_tg_route_reports":             transitgateway.DataSourceIBMTransitGatewayRouteReports(),

			// Added for BSS Enterprise
			"ibm_enterprises":               enterprise.DataSourceIBMEnterprises(),
			"ibm_enterprise_account_groups": enterprise.DataSourceIBMEnterpriseAccountGroups(),
			"ibm_enterprise_accounts":       enterprise.DataSourceIBMEnterpriseAccounts(),

			// //Added for Usage Reports
			"ibm_billing_snapshot_list": usagereports.DataSourceIBMBillingSnapshotList(),

			// Added for Secrets Manager
			"ibm_sm_secret_group":  secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmSecretGroup()),
			"ibm_sm_secret_groups": secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmSecretGroups()),
			"ibm_sm_private_certificate_configuration_intermediate_ca":           secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPrivateCertificateConfigurationIntermediateCA()),
			"ibm_sm_private_certificate_configuration_root_ca":                   secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPrivateCertificateConfigurationRootCA()),
			"ibm_sm_private_certificate_configuration_template":                  secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPrivateCertificateConfigurationTemplate()),
			"ibm_sm_public_certificate_configuration_ca_lets_encrypt":            secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPublicCertificateConfigurationCALetsEncrypt()),
			"ibm_sm_public_certificate_configuration_dns_cis":                    secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmConfigurationPublicCertificateDNSCis()),
			"ibm_sm_public_certificate_configuration_dns_classic_infrastructure": secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructure()),
			"ibm_sm_iam_credentials_configuration":                               secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmIamCredentialsConfiguration()),
			"ibm_sm_configurations":                                              secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmConfigurations()),
			"ibm_sm_secrets":                                                     secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmSecrets()),
			"ibm_sm_arbitrary_secret_metadata":                                   secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmArbitrarySecretMetadata()),
			"ibm_sm_imported_certificate_metadata":                               secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmImportedCertificateMetadata()),
			"ibm_sm_public_certificate_metadata":                                 secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPublicCertificateMetadata()),
			"ibm_sm_private_certificate_metadata":                                secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPrivateCertificateMetadata()),
			"ibm_sm_iam_credentials_secret_metadata":                             secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmIamCredentialsSecretMetadata()),
			"ibm_sm_service_credentials_secret_metadata":                         secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmServiceCredentialsSecretMetadata()),
			"ibm_sm_kv_secret_metadata":                                          secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmKvSecretMetadata()),
			"ibm_sm_username_password_secret_metadata":                           secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmUsernamePasswordSecretMetadata()),
			"ibm_sm_arbitrary_secret":                                            secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmArbitrarySecret()),
			"ibm_sm_imported_certificate":                                        secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmImportedCertificate()),
			"ibm_sm_public_certificate":                                          secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPublicCertificate()),
			"ibm_sm_private_certificate":                                         secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmPrivateCertificate()),
			"ibm_sm_iam_credentials_secret":                                      secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmIamCredentialsSecret()),
			"ibm_sm_kv_secret":                                                   secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmKvSecret()),
			"ibm_sm_username_password_secret":                                    secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmUsernamePasswordSecret()),
			"ibm_sm_service_credentials_secret":                                  secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmServiceCredentialsSecret()),
			"ibm_sm_en_registration":                                             secretsmanager.AddInstanceFields(secretsmanager.DataSourceIbmSmEnRegistration()),

			// Added for Satellite
			"ibm_satellite_location":                            satellite.DataSourceIBMSatelliteLocation(),
			"ibm_satellite_location_nlb_dns":                    satellite.DataSourceIBMSatelliteLocationNLBDNS(),
			"ibm_satellite_attach_host_script":                  satellite.DataSourceIBMSatelliteAttachHostScript(),
			"ibm_satellite_cluster":                             satellite.DataSourceIBMSatelliteCluster(),
			"ibm_satellite_cluster_worker_pool":                 satellite.DataSourceIBMSatelliteClusterWorkerPool(),
			"ibm_satellite_link":                                satellite.DataSourceIBMSatelliteLink(),
			"ibm_satellite_endpoint":                            satellite.DataSourceIBMSatelliteEndpoint(),
			"ibm_satellite_cluster_worker_pool_zone_attachment": satellite.DataSourceIBMSatelliteClusterWorkerPoolAttachment(),
			"ibm_satellite_storage_configuration":               satellite.DataSourceIBMSatelliteStorageConfiguration(),
			"ibm_satellite_storage_assignment":                  satellite.DataSourceIBMSatelliteStorageAssignment(),

			// Catalog related resources
			"ibm_cm_catalog":           catalogmanagement.DataSourceIBMCmCatalog(),
			"ibm_cm_offering":          catalogmanagement.DataSourceIBMCmOffering(),
			"ibm_cm_version":           catalogmanagement.DataSourceIBMCmVersion(),
			"ibm_cm_offering_instance": catalogmanagement.DataSourceIBMCmOfferingInstance(),
			"ibm_cm_preset":            catalogmanagement.DataSourceIBMCmPreset(),
			"ibm_cm_object":            catalogmanagement.DataSourceIBMCmObject(),
			"ibm_cm_account":           catalogmanagement.DataSourceIBMCmAccount(),

			// Added for Resource Tag
			"ibm_resource_tag":   globaltagging.DataSourceIBMResourceTag(),
			"ibm_iam_access_tag": globaltagging.DataSourceIBMIamAccessTag(),

			// Atracker
			"ibm_atracker_targets": atracker.DataSourceIBMAtrackerTargets(),
			"ibm_atracker_routes":  atracker.DataSourceIBMAtrackerRoutes(),

			// Metrics Router
			"ibm_metrics_router_targets": metricsrouter.DataSourceIBMMetricsRouterTargets(),
			"ibm_metrics_router_routes":  metricsrouter.DataSourceIBMMetricsRouterRoutes(),

			// MQ on Cloud
			"ibm_mqcloud_queue_manager_options":             mqcloud.DataSourceIbmMqcloudQueueManagerOptions(),
			"ibm_mqcloud_queue_manager":                     mqcloud.DataSourceIbmMqcloudQueueManager(),
			"ibm_mqcloud_queue_manager_status":              mqcloud.DataSourceIbmMqcloudQueueManagerStatus(),
			"ibm_mqcloud_application":                       mqcloud.DataSourceIbmMqcloudApplication(),
			"ibm_mqcloud_user":                              mqcloud.DataSourceIbmMqcloudUser(),
			"ibm_mqcloud_truststore_certificate":            mqcloud.DataSourceIbmMqcloudTruststoreCertificate(),
			"ibm_mqcloud_keystore_certificate":              mqcloud.DataSourceIbmMqcloudKeystoreCertificate(),
			"ibm_mqcloud_virtual_private_endpoint_gateways": mqcloud.DataSourceIbmMqcloudVirtualPrivateEndpointGateways(),
			"ibm_mqcloud_virtual_private_endpoint_gateway":  mqcloud.DataSourceIbmMqcloudVirtualPrivateEndpointGateway(),

			// Security and Complaince Center(soon to be deprecated)
			"ibm_scc_account_location":              scc.DataSourceIBMSccAccountLocation(),
			"ibm_scc_account_locations":             scc.DataSourceIBMSccAccountLocations(),
			"ibm_scc_account_location_settings":     scc.DataSourceIBMSccAccountLocationSettings(),
			"ibm_scc_account_notification_settings": scc.DataSourceIBMSccNotificationSettings(),

			// Security and Compliance Center
			"ibm_scc_instance_settings":        scc.DataSourceIbmSccInstanceSettings(),
			"ibm_scc_control_library":          scc.DataSourceIbmSccControlLibrary(),
			"ibm_scc_control_libraries":        scc.DataSourceIbmSccControlLibraries(),
			"ibm_scc_profile":                  scc.DataSourceIbmSccProfile(),
			"ibm_scc_profiles":                 scc.DataSourceIbmSccProfiles(),
			"ibm_scc_profile_attachment":       scc.DataSourceIbmSccProfileAttachment(),
			"ibm_scc_provider_type":            scc.DataSourceIbmSccProviderType(),
			"ibm_scc_provider_types":           scc.DataSourceIbmSccProviderTypes(),
			"ibm_scc_provider_type_collection": scc.DataSourceIbmSccProviderTypeCollection(),
			"ibm_scc_provider_type_instance":   scc.DataSourceIbmSccProviderTypeInstance(),
			"ibm_scc_latest_reports":           scc.DataSourceIbmSccLatestReports(),
			"ibm_scc_scope":                    scc.DataSourceIbmSccScope(),
			"ibm_scc_scope_collection":         scc.DataSourceIbmSccScopeCollection(),
			"ibm_scc_report":                   scc.DataSourceIbmSccReport(),
			"ibm_scc_report_controls":          scc.DataSourceIbmSccReportControls(),
			"ibm_scc_report_evaluations":       scc.DataSourceIbmSccReportEvaluations(),
			"ibm_scc_report_resources":         scc.DataSourceIbmSccReportResources(),
			"ibm_scc_report_rule":              scc.DataSourceIbmSccReportRule(),
			"ibm_scc_report_summary":           scc.DataSourceIbmSccReportSummary(),
			"ibm_scc_report_tags":              scc.DataSourceIbmSccReportTags(),
			"ibm_scc_report_violation_drift":   scc.DataSourceIbmSccReportViolationDrift(),
			"ibm_scc_rule":                     scc.DataSourceIbmSccRule(),

			// Security Services
			"ibm_pag_instance": pag.DataSourceIBMPag(),

			// Added for Context Based Restrictions
			"ibm_cbr_zone":           contextbasedrestrictions.DataSourceIBMCbrZone(),
			"ibm_cbr_zone_addresses": contextbasedrestrictions.DataSourceIBMCbrZoneAddresses(),
			"ibm_cbr_rule":           contextbasedrestrictions.DataSourceIBMCbrRule(),

			// Added for Event Notifications
			"ibm_en_source":                     eventnotification.DataSourceIBMEnSource(),
			"ibm_en_destinations":               eventnotification.DataSourceIBMEnDestinations(),
			"ibm_en_topic":                      eventnotification.DataSourceIBMEnTopic(),
			"ibm_en_topics":                     eventnotification.DataSourceIBMEnTopics(),
			"ibm_en_subscriptions":              eventnotification.DataSourceIBMEnSubscriptions(),
			"ibm_en_destination_webhook":        eventnotification.DataSourceIBMEnWebhookDestination(),
			"ibm_en_destination_android":        eventnotification.DataSourceIBMEnFCMDestination(),
			"ibm_en_destination_ios":            eventnotification.DataSourceIBMEnAPNSDestination(),
			"ibm_en_destination_chrome":         eventnotification.DataSourceIBMEnChromeDestination(),
			"ibm_en_destination_firefox":        eventnotification.DataSourceIBMEnFirefoxDestination(),
			"ibm_en_destination_slack":          eventnotification.DataSourceIBMEnSlackDestination(),
			"ibm_en_subscription_sms":           eventnotification.DataSourceIBMEnSMSSubscription(),
			"ibm_en_subscription_email":         eventnotification.DataSourceIBMEnEmailSubscription(),
			"ibm_en_subscription_webhook":       eventnotification.DataSourceIBMEnWebhookSubscription(),
			"ibm_en_subscription_android":       eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_subscription_ios":           eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_subscription_chrome":        eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_subscription_firefox":       eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_subscription_slack":         eventnotification.DataSourceIBMEnSlackSubscription(),
			"ibm_en_subscription_safari":        eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_destination_safari":         eventnotification.DataSourceIBMEnSafariDestination(),
			"ibm_en_destination_msteams":        eventnotification.DataSourceIBMEnMSTeamsDestination(),
			"ibm_en_subscription_msteams":       eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_destination_pagerduty":      eventnotification.DataSourceIBMEnPagerDutyDestination(),
			"ibm_en_subscription_pagerduty":     eventnotification.DataSourceIBMEnPagerDutySubscription(),
			"ibm_en_integration":                eventnotification.DataSourceIBMEnIntegration(),
			"ibm_en_integrations":               eventnotification.DataSourceIBMEnIntegrations(),
			"ibm_en_destination_sn":             eventnotification.DataSourceIBMEnServiceNowDestination(),
			"ibm_en_subscription_sn":            eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_destination_ce":             eventnotification.DataSourceIBMEnCodeEngineDestination(),
			"ibm_en_subscription_ce":            eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_destination_cos":            eventnotification.DataSourceIBMEnCOSDestination(),
			"ibm_en_subscription_cos":           eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_destination_huawei":         eventnotification.DataSourceIBMEnHuaweiDestination(),
			"ibm_en_subscription_huawei":        eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_sources":                    eventnotification.DataSourceIBMEnSources(),
			"ibm_en_destination_custom_email":   eventnotification.DataSourceIBMEnCustomEmailDestination(),
			"ibm_en_subscription_custom_email":  eventnotification.DataSourceIBMEnCustomEmailSubscription(),
			"ibm_en_email_template":             eventnotification.DataSourceIBMEnEmailTemplate(),
			"ibm_en_email_templates":            eventnotification.DataSourceIBMEnTemplates(),
			"ibm_en_destination_custom_sms":     eventnotification.DataSourceIBMEnCustomSMSDestination(),
			"ibm_en_subscription_custom_sms":    eventnotification.DataSourceIBMEnCustomSMSSubscription(),
			"ibm_en_integration_cos":            eventnotification.DataSourceIBMEnCOSIntegration(),
			"ibm_en_smtp_configuration":         eventnotification.DataSourceIBMEnSMTPConfiguration(),
			"ibm_en_smtp_configurations":        eventnotification.DataSourceIBMEnSMTPCOnfigurations(),
			"ibm_en_smtp_user":                  eventnotification.DataSourceIBMEnSMTPUser(),
			"ibm_en_smtp_users":                 eventnotification.DataSourceIBMEnSMTPUsers(),
			"ibm_en_slack_template":             eventnotification.DataSourceIBMEnSlackTemplate(),
			"ibm_en_metrics":                    eventnotification.DataSourceIBMEnMetrics(),
			"ibm_en_smtp_allowed_ips":           eventnotification.DataSourceIBMEnSMTPAllowedIps(),
			"ibm_en_webhook_template":           eventnotification.DataSourceIBMEnWebhookTemplate(),
			"ibm_en_subscription_scheduler":     eventnotification.DataSourceIBMEnFCMSubscription(),
			"ibm_en_pagerduty_template":         eventnotification.DataSourceIBMEnPagerDutyTemplate(),
			"ibm_en_destination_event_streams":  eventnotification.DataSourceIBMEnEventStreamsDestination(),
			"ibm_en_subscription_event_streams": eventnotification.DataSourceIBMEnEventStreamsSubscription(),
			"ibm_en_event_streams_template":     eventnotification.DataSourceIBMEnEventStreamsTemplate(),

			// Added for Toolchain
			"ibm_cd_toolchain":                         cdtoolchain.DataSourceIBMCdToolchain(),
			"ibm_cd_toolchains":                        cdtoolchain.DataSourceIBMCdToolchains(),
			"ibm_cd_toolchain_tool_keyprotect":         cdtoolchain.DataSourceIBMCdToolchainToolKeyprotect(),
			"ibm_cd_toolchain_tool_secretsmanager":     cdtoolchain.DataSourceIBMCdToolchainToolSecretsmanager(),
			"ibm_cd_toolchain_tool_bitbucketgit":       cdtoolchain.DataSourceIBMCdToolchainToolBitbucketgit(),
			"ibm_cd_toolchain_tool_githubconsolidated": cdtoolchain.DataSourceIBMCdToolchainToolGithubconsolidated(),
			"ibm_cd_toolchain_tool_gitlab":             cdtoolchain.DataSourceIBMCdToolchainToolGitlab(),
			"ibm_cd_toolchain_tool_hostedgit":          cdtoolchain.DataSourceIBMCdToolchainToolHostedgit(),
			"ibm_cd_toolchain_tool_artifactory":        cdtoolchain.DataSourceIBMCdToolchainToolArtifactory(),
			"ibm_cd_toolchain_tool_custom":             cdtoolchain.DataSourceIBMCdToolchainToolCustom(),
			"ibm_cd_toolchain_tool_pipeline":           cdtoolchain.DataSourceIBMCdToolchainToolPipeline(),
			"ibm_cd_toolchain_tool_devopsinsights":     cdtoolchain.DataSourceIBMCdToolchainToolDevopsinsights(),
			"ibm_cd_toolchain_tool_slack":              cdtoolchain.DataSourceIBMCdToolchainToolSlack(),
			"ibm_cd_toolchain_tool_sonarqube":          cdtoolchain.DataSourceIBMCdToolchainToolSonarqube(),
			"ibm_cd_toolchain_tool_hashicorpvault":     cdtoolchain.DataSourceIBMCdToolchainToolHashicorpvault(),
			"ibm_cd_toolchain_tool_securitycompliance": cdtoolchain.DataSourceIBMCdToolchainToolSecuritycompliance(),
			"ibm_cd_toolchain_tool_privateworker":      cdtoolchain.DataSourceIBMCdToolchainToolPrivateworker(),
			"ibm_cd_toolchain_tool_appconfig":          cdtoolchain.DataSourceIBMCdToolchainToolAppconfig(),
			"ibm_cd_toolchain_tool_jenkins":            cdtoolchain.DataSourceIBMCdToolchainToolJenkins(),
			"ibm_cd_toolchain_tool_nexus":              cdtoolchain.DataSourceIBMCdToolchainToolNexus(),
			"ibm_cd_toolchain_tool_pagerduty":          cdtoolchain.DataSourceIBMCdToolchainToolPagerduty(),
			"ibm_cd_toolchain_tool_saucelabs":          cdtoolchain.DataSourceIBMCdToolchainToolSaucelabs(),
			"ibm_cd_toolchain_tool_jira":               cdtoolchain.DataSourceIBMCdToolchainToolJira(),
			"ibm_cd_toolchain_tool_eventnotifications": cdtoolchain.DataSourceIBMCdToolchainToolEventnotifications(),

			// Added for Tekton Pipeline
			"ibm_cd_tekton_pipeline_definition":       cdtektonpipeline.DataSourceIBMCdTektonPipelineDefinition(),
			"ibm_cd_tekton_pipeline_trigger_property": cdtektonpipeline.DataSourceIBMCdTektonPipelineTriggerProperty(),
			"ibm_cd_tekton_pipeline_property":         cdtektonpipeline.DataSourceIBMCdTektonPipelineProperty(),
			"ibm_cd_tekton_pipeline_trigger":          cdtektonpipeline.DataSourceIBMCdTektonPipelineTrigger(),
			"ibm_cd_tekton_pipeline":                  cdtektonpipeline.DataSourceIBMCdTektonPipeline(),

			// Added for Code Engine
			"ibm_code_engine_allowed_outbound_destination": codeengine.DataSourceIbmCodeEngineAllowedOutboundDestination(),
			"ibm_code_engine_app":                          codeengine.DataSourceIbmCodeEngineApp(),
			"ibm_code_engine_binding":                      codeengine.DataSourceIbmCodeEngineBinding(),
			"ibm_code_engine_build":                        codeengine.DataSourceIbmCodeEngineBuild(),
			"ibm_code_engine_config_map":                   codeengine.DataSourceIbmCodeEngineConfigMap(),
			"ibm_code_engine_domain_mapping":               codeengine.DataSourceIbmCodeEngineDomainMapping(),
			"ibm_code_engine_function":                     codeengine.DataSourceIbmCodeEngineFunction(),
			"ibm_code_engine_job":                          codeengine.DataSourceIbmCodeEngineJob(),
			"ibm_code_engine_project":                      codeengine.DataSourceIbmCodeEngineProject(),
			"ibm_code_engine_secret":                       codeengine.DataSourceIbmCodeEngineSecret(),

			// Added for Project
			"ibm_project":             project.DataSourceIbmProject(),
			"ibm_project_config":      project.DataSourceIbmProjectConfig(),
			"ibm_project_environment": project.DataSourceIbmProjectEnvironment(),

			// Added for VMware as a Service
			"ibm_vmaas_vdc": vmware.DataSourceIbmVmaasVdc(),
			// Logs Service
			"ibm_logs_alert":              logs.AddLogsInstanceFields(logs.DataSourceIbmLogsAlert()),
			"ibm_logs_alerts":             logs.AddLogsInstanceFields(logs.DataSourceIbmLogsAlerts()),
			"ibm_logs_rule_group":         logs.AddLogsInstanceFields(logs.DataSourceIbmLogsRuleGroup()),
			"ibm_logs_rule_groups":        logs.AddLogsInstanceFields(logs.DataSourceIbmLogsRuleGroups()),
			"ibm_logs_policy":             logs.AddLogsInstanceFields(logs.DataSourceIbmLogsPolicy()),
			"ibm_logs_policies":           logs.AddLogsInstanceFields(logs.DataSourceIbmLogsPolicies()),
			"ibm_logs_dashboard":          logs.AddLogsInstanceFields(logs.DataSourceIbmLogsDashboard()),
			"ibm_logs_e2m":                logs.AddLogsInstanceFields(logs.DataSourceIbmLogsE2m()),
			"ibm_logs_e2ms":               logs.AddLogsInstanceFields(logs.DataSourceIbmLogsE2ms()),
			"ibm_logs_outgoing_webhook":   logs.AddLogsInstanceFields(logs.DataSourceIbmLogsOutgoingWebhook()),
			"ibm_logs_outgoing_webhooks":  logs.AddLogsInstanceFields(logs.DataSourceIbmLogsOutgoingWebhooks()),
			"ibm_logs_view_folder":        logs.AddLogsInstanceFields(logs.DataSourceIbmLogsViewFolder()),
			"ibm_logs_view_folders":       logs.AddLogsInstanceFields(logs.DataSourceIbmLogsViewFolders()),
			"ibm_logs_view":               logs.AddLogsInstanceFields(logs.DataSourceIbmLogsView()),
			"ibm_logs_views":              logs.AddLogsInstanceFields(logs.DataSourceIbmLogsViews()),
			"ibm_logs_dashboard_folders":  logs.AddLogsInstanceFields(logs.DataSourceIbmLogsDashboardFolders()),
			"ibm_logs_data_usage_metrics": logs.AddLogsInstanceFields(logs.DataSourceIbmLogsDataUsageMetrics()),
			"ibm_logs_enrichments":        logs.AddLogsInstanceFields(logs.DataSourceIbmLogsEnrichments()),
			"ibm_logs_data_access_rules":  logs.AddLogsInstanceFields(logs.DataSourceIbmLogsDataAccessRules()),
			"ibm_logs_stream":             logs.AddLogsInstanceFields(logs.DataSourceIbmLogsStream()),
			"ibm_logs_streams":            logs.AddLogsInstanceFields(logs.DataSourceIbmLogsStreams()),

			// Logs Router Service
			"ibm_logs_router_tenants": logsrouting.DataSourceIBMLogsRouterTenants(),
			"ibm_logs_router_targets": logsrouting.DataSourceIBMLogsRouterTargets(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ibm_backup_recovery_agent_upgrade_task":                             backuprecovery.ResourceIbmBackupRecoveryAgentUpgradeTask(),
			"ibm_backup_recovery_protection_group_run_request":                   backuprecovery.ResourceIbmBackupRecoveryProtectionGroupRunRequest(),
			"ibm_backup_recovery_data_source_connection":                         backuprecovery.ResourceIbmBackupRecoveryDataSourceConnection(),
			"ibm_backup_recovery_data_source_connector_patch":                    backuprecovery.ResourceIbmBackupRecoveryDataSourceConnectorPatch(),
			"ibm_backup_recovery_download_files_folders":                         backuprecovery.ResourceIbmBackupRecoveryDownloadFilesFolders(),
			"ibm_backup_recovery_restore_points":                                 backuprecovery.ResourceIbmBackupRecoveryRestorePoints(),
			"ibm_backup_recovery_perform_action_on_protection_group_run_request": backuprecovery.ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequest(),
			"ibm_backup_recovery_protection_group":                               backuprecovery.ResourceIbmBackupRecoveryProtectionGroup(),
			"ibm_backup_recovery_protection_policy":                              backuprecovery.ResourceIbmBackupRecoveryProtectionPolicy(),
			"ibm_backup_recovery":                                                backuprecovery.ResourceIbmBackupRecovery(),
			"ibm_backup_recovery_source_registration":                            backuprecovery.ResourceIbmBackupRecoverySourceRegistration(),
			"ibm_backup_recovery_update_protection_group_run_request":            backuprecovery.ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequest(),
			"ibm_backup_recovery_connection_registration_token":                  backuprecovery.ResourceIbmBackupRecoveryConnectionRegistrationToken(),

			"ibm_api_gateway_endpoint":              apigateway.ResourceIBMApiGatewayEndPoint(),
			"ibm_api_gateway_endpoint_subscription": apigateway.ResourceIBMApiGatewayEndpointSubscription(),
			"ibm_app":                               cloudfoundry.ResourceIBMApp(),
			"ibm_app_domain_private":                cloudfoundry.ResourceIBMAppDomainPrivate(),
			"ibm_app_domain_shared":                 cloudfoundry.ResourceIBMAppDomainShared(),
			"ibm_app_route":                         cloudfoundry.ResourceIBMAppRoute(),
			"ibm_config_aggregator_settings":        configurationaggregator.AddConfigurationAggregatorInstanceFields(configurationaggregator.ResourceIbmConfigAggregatorSettings()),

			// AppID
			"ibm_appid_action_url":               appid.ResourceIBMAppIDActionURL(),
			"ibm_appid_apm":                      appid.ResourceIBMAppIDAPM(),
			"ibm_appid_application":              appid.ResourceIBMAppIDApplication(),
			"ibm_appid_application_scopes":       appid.ResourceIBMAppIDApplicationScopes(),
			"ibm_appid_application_roles":        appid.ResourceIBMAppIDApplicationRoles(),
			"ibm_appid_audit_status":             appid.ResourceIBMAppIDAuditStatus(),
			"ibm_appid_cloud_directory_template": appid.ResourceIBMAppIDCloudDirectoryTemplate(),
			"ibm_appid_cloud_directory_user":     appid.ResourceIBMAppIDCloudDirectoryUser(),
			"ibm_appid_idp_cloud_directory":      appid.ResourceIBMAppIDIDPCloudDirectory(),
			"ibm_appid_idp_custom":               appid.ResourceIBMAppIDIDPCustom(),
			"ibm_appid_idp_facebook":             appid.ResourceIBMAppIDIDPFacebook(),
			"ibm_appid_idp_google":               appid.ResourceIBMAppIDIDPGoogle(),
			"ibm_appid_idp_saml":                 appid.ResourceIBMAppIDIDPSAML(),
			"ibm_appid_languages":                appid.ResourceIBMAppIDLanguages(),
			"ibm_appid_mfa":                      appid.ResourceIBMAppIDMFA(),
			"ibm_appid_mfa_channel":              appid.ResourceIBMAppIDMFAChannel(),
			"ibm_appid_password_regex":           appid.ResourceIBMAppIDPasswordRegex(),
			"ibm_appid_token_config":             appid.ResourceIBMAppIDTokenConfig(),
			"ibm_appid_redirect_urls":            appid.ResourceIBMAppIDRedirectURLs(),
			"ibm_appid_role":                     appid.ResourceIBMAppIDRole(),
			"ibm_appid_theme_color":              appid.ResourceIBMAppIDThemeColor(),
			"ibm_appid_theme_text":               appid.ResourceIBMAppIDThemeText(),
			"ibm_appid_user_roles":               appid.ResourceIBMAppIDUserRoles(),

			"ibm_function_action":    functions.ResourceIBMFunctionAction(),
			"ibm_function_package":   functions.ResourceIBMFunctionPackage(),
			"ibm_function_rule":      functions.ResourceIBMFunctionRule(),
			"ibm_function_trigger":   functions.ResourceIBMFunctionTrigger(),
			"ibm_function_namespace": functions.ResourceIBMFunctionNamespace(),

			"ibm_cis":                                 cis.ResourceIBMCISInstance(),
			"ibm_database":                            database.ResourceIBMDatabaseInstance(),
			"ibm_db2":                                 db2.ResourceIBMDb2Instance(),
			"ibm_cis_domain":                          cis.ResourceIBMCISDomain(),
			"ibm_cis_domain_settings":                 cis.ResourceIBMCISSettings(),
			"ibm_cis_firewall":                        cis.ResourceIBMCISFirewallRecord(),
			"ibm_cis_range_app":                       cis.ResourceIBMCISRangeApp(),
			"ibm_cis_healthcheck":                     cis.ResourceIBMCISHealthCheck(),
			"ibm_cis_origin_pool":                     cis.ResourceIBMCISPool(),
			"ibm_cis_global_load_balancer":            cis.ResourceIBMCISGlb(),
			"ibm_cis_certificate_upload":              cis.ResourceIBMCISCertificateUpload(),
			"ibm_cis_dns_record":                      cis.ResourceIBMCISDnsRecord(),
			"ibm_cis_dns_records_import":              cis.ResourceIBMCISDNSRecordsImport(),
			"ibm_cis_rate_limit":                      cis.ResourceIBMCISRateLimit(),
			"ibm_cis_page_rule":                       cis.ResourceIBMCISPageRule(),
			"ibm_cis_edge_functions_action":           cis.ResourceIBMCISEdgeFunctionsAction(),
			"ibm_cis_edge_functions_trigger":          cis.ResourceIBMCISEdgeFunctionsTrigger(),
			"ibm_cis_tls_settings":                    cis.ResourceIBMCISTLSSettings(),
			"ibm_cis_waf_package":                     cis.ResourceIBMCISWAFPackage(),
			"ibm_cis_webhook":                         cis.ResourceIBMCISWebhooks(),
			"ibm_cis_origin_auth":                     cis.ResourceIBMCISOriginAuthPull(),
			"ibm_cis_mtls":                            cis.ResourceIBMCISMtls(),
			"ibm_cis_mtls_app":                        cis.ResourceIBMCISMtlsApp(),
			"ibm_cis_bot_management":                  cis.ResourceIBMCISBotManagement(),
			"ibm_cis_logpush_job":                     cis.ResourceIBMCISLogPushJob(),
			"ibm_cis_alert":                           cis.ResourceIBMCISAlert(),
			"ibm_cis_routing":                         cis.ResourceIBMCISRouting(),
			"ibm_cis_waf_group":                       cis.ResourceIBMCISWAFGroup(),
			"ibm_cis_cache_settings":                  cis.ResourceIBMCISCacheSettings(),
			"ibm_cis_custom_page":                     cis.ResourceIBMCISCustomPage(),
			"ibm_cis_waf_rule":                        cis.ResourceIBMCISWAFRule(),
			"ibm_cis_certificate_order":               cis.ResourceIBMCISCertificateOrder(),
			"ibm_cis_filter":                          cis.ResourceIBMCISFilter(),
			"ibm_cis_firewall_rule":                   cis.ResourceIBMCISFirewallrules(),
			"ibm_cis_ruleset":                         cis.ResourceIBMCISRuleset(),
			"ibm_cis_ruleset_version_detach":          cis.ResourceIBMCISRulesetVersionDetach(),
			"ibm_cis_ruleset_rule":                    cis.ResourceIBMCISRulesetRule(),
			"ibm_cis_ruleset_entrypoint_version":      cis.ResourceIBMCISRulesetEntryPointVersion(),
			"ibm_cis_advanced_certificate_pack_order": cis.ResourceIBMCISAdvancedCertificatePackOrder(),
			"ibm_cis_origin_certificate_order":        cis.ResourceIBMCISOriginCertificateOrder(),

			"ibm_cloudant":                                  cloudant.ResourceIBMCloudant(),
			"ibm_cloudant_database":                         cloudant.ResourceIBMCloudantDatabase(),
			"ibm_cloud_shell_account_settings":              cloudshell.ResourceIBMCloudShellAccountSettings(),
			"ibm_compute_autoscale_group":                   classicinfrastructure.ResourceIBMComputeAutoScaleGroup(),
			"ibm_compute_autoscale_policy":                  classicinfrastructure.ResourceIBMComputeAutoScalePolicy(),
			"ibm_compute_bare_metal":                        classicinfrastructure.ResourceIBMComputeBareMetal(),
			"ibm_compute_dedicated_host":                    classicinfrastructure.ResourceIBMComputeDedicatedHost(),
			"ibm_compute_monitor":                           classicinfrastructure.ResourceIBMComputeMonitor(),
			"ibm_compute_placement_group":                   classicinfrastructure.ResourceIBMComputePlacementGroup(),
			"ibm_compute_reserved_capacity":                 classicinfrastructure.ResourceIBMComputeReservedCapacity(),
			"ibm_compute_provisioning_hook":                 classicinfrastructure.ResourceIBMComputeProvisioningHook(),
			"ibm_compute_ssh_key":                           classicinfrastructure.ResourceIBMComputeSSHKey(),
			"ibm_compute_ssl_certificate":                   classicinfrastructure.ResourceIBMComputeSSLCertificate(),
			"ibm_compute_user":                              classicinfrastructure.ResourceIBMComputeUser(),
			"ibm_compute_vm_instance":                       classicinfrastructure.ResourceIBMComputeVmInstance(),
			"ibm_container_addons":                          kubernetes.ResourceIBMContainerAddOns(),
			"ibm_container_alb":                             kubernetes.ResourceIBMContainerALB(),
			"ibm_container_alb_create":                      kubernetes.ResourceIBMContainerAlbCreate(),
			"ibm_container_api_key_reset":                   kubernetes.ResourceIBMContainerAPIKeyReset(),
			"ibm_container_vpc_alb":                         kubernetes.ResourceIBMContainerVpcALB(),
			"ibm_container_vpc_alb_create":                  kubernetes.ResourceIBMContainerVpcAlbCreateNew(),
			"ibm_container_vpc_worker_pool":                 kubernetes.ResourceIBMContainerVpcWorkerPool(),
			"ibm_container_vpc_worker":                      kubernetes.ResourceIBMContainerVpcWorker(),
			"ibm_container_vpc_cluster":                     kubernetes.ResourceIBMContainerVpcCluster(),
			"ibm_container_alb_cert":                        kubernetes.ResourceIBMContainerALBCert(),
			"ibm_container_ingress_instance":                kubernetes.ResourceIBMContainerIngressInstance(),
			"ibm_container_ingress_secret_tls":              kubernetes.ResourceIBMContainerIngressSecretTLS(),
			"ibm_container_ingress_secret_opaque":           kubernetes.ResourceIBMContainerIngressSecretOpaque(),
			"ibm_container_cluster":                         kubernetes.ResourceIBMContainerCluster(),
			"ibm_container_cluster_feature":                 kubernetes.ResourceIBMContainerClusterFeature(),
			"ibm_container_bind_service":                    kubernetes.ResourceIBMContainerBindService(),
			"ibm_container_worker_pool":                     kubernetes.ResourceIBMContainerWorkerPool(),
			"ibm_container_worker_pool_zone_attachment":     kubernetes.ResourceIBMContainerWorkerPoolZoneAttachment(),
			"ibm_container_storage_attachment":              kubernetes.ResourceIBMContainerVpcWorkerVolumeAttachment(),
			"ibm_container_nlb_dns":                         kubernetes.ResourceIBMContainerNlbDns(),
			"ibm_container_dedicated_host_pool":             kubernetes.ResourceIBMContainerDedicatedHostPool(),
			"ibm_container_dedicated_host":                  kubernetes.ResourceIBMContainerDedicatedHost(),
			"ibm_cr_namespace":                              registry.ResourceIBMCrNamespace(),
			"ibm_cr_retention_policy":                       registry.ResourceIBMCrRetentionPolicy(),
			"ibm_ob_logging":                                kubernetes.ResourceIBMObLogging(),
			"ibm_ob_monitoring":                             kubernetes.ResourceIBMObMonitoring(),
			"ibm_cos_bucket":                                cos.ResourceIBMCOSBucket(),
			"ibm_cos_bucket_replication_rule":               cos.ResourceIBMCOSBucketReplicationConfiguration(),
			"ibm_cos_bucket_object":                         cos.ResourceIBMCOSBucketObject(),
			"ibm_cos_bucket_object_lock_configuration":      cos.ResourceIBMCOSBucketObjectlock(),
			"ibm_cos_bucket_website_configuration":          cos.ResourceIBMCOSBucketWebsiteConfiguration(),
			"ibm_cos_bucket_lifecycle_configuration":        cos.ResourceIBMCOSBucketLifecycleConfiguration(),
			"ibm_cos_backup_vault":                          cos.ResourceIBMCOSBackupVault(),
			"ibm_cos_backup_policy":                         cos.ResourceIBMCOSBackupPolicy(),
			"ibm_dns_domain":                                classicinfrastructure.ResourceIBMDNSDomain(),
			"ibm_dns_domain_registration_nameservers":       classicinfrastructure.ResourceIBMDNSDomainRegistrationNameservers(),
			"ibm_dns_secondary":                             classicinfrastructure.ResourceIBMDNSSecondary(),
			"ibm_dns_record":                                classicinfrastructure.ResourceIBMDNSRecord(),
			"ibm_event_streams_topic":                       eventstreams.ResourceIBMEventStreamsTopic(),
			"ibm_event_streams_schema":                      eventstreams.ResourceIBMEventStreamsSchema(),
			"ibm_event_streams_schema_global_rule":          eventstreams.ResourceIBMEventStreamsSchemaGlobalCompatibilityRule(),
			"ibm_event_streams_quota":                       eventstreams.ResourceIBMEventStreamsQuota(),
			"ibm_event_streams_mirroring_config":            eventstreams.ResourceIBMEventStreamsMirroringConfig(),
			"ibm_firewall":                                  classicinfrastructure.ResourceIBMFirewall(),
			"ibm_firewall_policy":                           classicinfrastructure.ResourceIBMFirewallPolicy(),
			"ibm_hpcs":                                      hpcs.ResourceIBMHPCS(),
			"ibm_hpcs_managed_key":                          hpcs.ResourceIbmManagedKey(),
			"ibm_hpcs_key_template":                         hpcs.ResourceIbmKeyTemplate(),
			"ibm_hpcs_keystore":                             hpcs.ResourceIbmKeystore(),
			"ibm_hpcs_vault":                                hpcs.ResourceIbmVault(),
			"ibm_iam_access_group":                          iamaccessgroup.ResourceIBMIAMAccessGroup(),
			"ibm_iam_access_group_account_settings":         iamaccessgroup.ResourceIBMIAMAccessGroupAccountSettings(),
			"ibm_iam_account_settings":                      iamidentity.ResourceIBMIAMAccountSettings(),
			"ibm_iam_access_group_template":                 iamaccessgroup.ResourceIBMIAMAccessGroupTemplate(),
			"ibm_iam_access_group_template_version":         iamaccessgroup.ResourceIBMIAMAccessGroupTemplateVersion(),
			"ibm_iam_access_group_template_assignment":      iamaccessgroup.ResourceIBMIAMAccessGroupTemplateAssignment(),
			"ibm_iam_custom_role":                           iampolicy.ResourceIBMIAMCustomRole(),
			"ibm_iam_access_group_dynamic_rule":             iamaccessgroup.ResourceIBMIAMDynamicRule(),
			"ibm_iam_access_group_members":                  iamaccessgroup.ResourceIBMIAMAccessGroupMembers(),
			"ibm_iam_access_group_policy":                   iampolicy.ResourceIBMIAMAccessGroupPolicy(),
			"ibm_iam_authorization_policy":                  iampolicy.ResourceIBMIAMAuthorizationPolicy(),
			"ibm_iam_authorization_policy_detach":           iampolicy.ResourceIBMIAMAuthorizationPolicyDetach(),
			"ibm_iam_user_policy":                           iampolicy.ResourceIBMIAMUserPolicy(),
			"ibm_iam_user_settings":                         iamidentity.ResourceIBMIAMUserSettings(),
			"ibm_iam_service_id":                            iamidentity.ResourceIBMIAMServiceID(),
			"ibm_iam_service_api_key":                       iamidentity.ResourceIBMIAMServiceAPIKey(),
			"ibm_iam_service_policy":                        iampolicy.ResourceIBMIAMServicePolicy(),
			"ibm_iam_user_invite":                           iampolicy.ResourceIBMIAMUserInvite(),
			"ibm_iam_api_key":                               iamidentity.ResourceIBMIAMApiKey(),
			"ibm_iam_trusted_profile":                       iamidentity.ResourceIBMIAMTrustedProfile(),
			"ibm_iam_trusted_profile_identity":              iamidentity.ResourceIBMIamTrustedProfileIdentity(),
			"ibm_iam_trusted_profile_claim_rule":            iamidentity.ResourceIBMIAMTrustedProfileClaimRule(),
			"ibm_iam_trusted_profile_link":                  iamidentity.ResourceIBMIAMTrustedProfileLink(),
			"ibm_iam_trusted_profile_policy":                iampolicy.ResourceIBMIAMTrustedProfilePolicy(),
			"ibm_iam_account_settings_template":             iamidentity.ResourceIBMAccountSettingsTemplate(),
			"ibm_iam_trusted_profile_template":              iamidentity.ResourceIBMTrustedProfileTemplate(),
			"ibm_iam_account_settings_template_assignment":  iamidentity.ResourceIBMAccountSettingsTemplateAssignment(),
			"ibm_iam_trusted_profile_template_assignment":   iamidentity.ResourceIBMTrustedProfileTemplateAssignment(),
			"ibm_ipsec_vpn":                                 classicinfrastructure.ResourceIBMIPSecVPN(),
			"ibm_iam_policy_template":                       iampolicy.ResourceIBMIAMPolicyTemplate(),
			"ibm_iam_policy_template_version":               iampolicy.ResourceIBMIAMPolicyTemplateVersion(),
			"ibm_iam_policy_assignment":                     iampolicy.ResourceIBMIAMPolicyAssignment(),
			"ibm_iam_account_settings_external_interaction": iampolicy.ResourceIBMIAMAccountSettingsExternalInteraction(),

			"ibm_is_backup_policy":      vpc.ResourceIBMIsBackupPolicy(),
			"ibm_is_backup_policy_plan": vpc.ResourceIBMIsBackupPolicyPlan(),

			// cluster
			"ibm_is_cluster_network_interface":           vpc.ResourceIBMIsClusterNetworkInterface(),
			"ibm_is_cluster_network_subnet_reserved_ip":  vpc.ResourceIBMIsClusterNetworkSubnetReservedIP(),
			"ibm_is_cluster_network_subnet":              vpc.ResourceIBMIsClusterNetworkSubnet(),
			"ibm_is_cluster_network":                     vpc.ResourceIBMIsClusterNetwork(),
			"ibm_is_instance_cluster_network_attachment": vpc.ResourceIBMIsInstanceClusterNetworkAttachment(),

			// bare_metal_server
			"ibm_is_bare_metal_server_action":                        vpc.ResourceIBMIsBareMetalServerAction(),
			"ibm_is_bare_metal_server_disk":                          vpc.ResourceIBMIsBareMetalServerDisk(),
			"ibm_is_bare_metal_server_initialization":                vpc.ResourceIBMIsBareMetalServerInitialization(),
			"ibm_is_bare_metal_server_network_attachment":            vpc.ResourceIBMIsBareMetalServerNetworkAttachment(),
			"ibm_is_bare_metal_server_network_interface_allow_float": vpc.ResourceIBMIsBareMetalServerNetworkInterfaceAllowFloat(),
			"ibm_is_bare_metal_server_network_interface_floating_ip": vpc.ResourceIBMIsBareMetalServerNetworkInterfaceFloatingIp(),
			"ibm_is_bare_metal_server_network_interface":             vpc.ResourceIBMIsBareMetalServerNetworkInterface(),
			"ibm_is_bare_metal_server":                               vpc.ResourceIBMIsBareMetalServer(),

			"ibm_is_dedicated_host":                              vpc.ResourceIbmIsDedicatedHost(),
			"ibm_is_dedicated_host_group":                        vpc.ResourceIbmIsDedicatedHostGroup(),
			"ibm_is_dedicated_host_disk_management":              vpc.ResourceIBMISDedicatedHostDiskManagement(),
			"ibm_is_placement_group":                             vpc.ResourceIbmIsPlacementGroup(),
			"ibm_is_floating_ip":                                 vpc.ResourceIBMISFloatingIP(),
			"ibm_is_flow_log":                                    vpc.ResourceIBMISFlowLog(),
			"ibm_is_instance":                                    vpc.ResourceIBMISInstance(),
			"ibm_is_instance_action":                             vpc.ResourceIBMISInstanceAction(),
			"ibm_is_instance_network_attachment":                 vpc.ResourceIBMIsInstanceNetworkAttachment(),
			"ibm_is_instance_network_interface":                  vpc.ResourceIBMIsInstanceNetworkInterface(),
			"ibm_is_instance_network_interface_floating_ip":      vpc.ResourceIBMIsInstanceNetworkInterfaceFloatingIp(),
			"ibm_is_instance_disk_management":                    vpc.ResourceIBMISInstanceDiskManagement(),
			"ibm_is_instance_group":                              vpc.ResourceIBMISInstanceGroup(),
			"ibm_is_instance_group_membership":                   vpc.ResourceIBMISInstanceGroupMembership(),
			"ibm_is_instance_group_manager":                      vpc.ResourceIBMISInstanceGroupManager(),
			"ibm_is_instance_group_manager_policy":               vpc.ResourceIBMISInstanceGroupManagerPolicy(),
			"ibm_is_instance_group_manager_action":               vpc.ResourceIBMISInstanceGroupManagerAction(),
			"ibm_is_instance_volume_attachment":                  vpc.ResourceIBMISInstanceVolumeAttachment(),
			"ibm_is_virtual_endpoint_gateway":                    vpc.ResourceIBMISEndpointGateway(),
			"ibm_is_virtual_endpoint_gateway_ip":                 vpc.ResourceIBMISEndpointGatewayIP(),
			"ibm_is_instance_template":                           vpc.ResourceIBMISInstanceTemplate(),
			"ibm_is_ike_policy":                                  vpc.ResourceIBMISIKEPolicy(),
			"ibm_is_ipsec_policy":                                vpc.ResourceIBMISIPSecPolicy(),
			"ibm_is_lb":                                          vpc.ResourceIBMISLB(),
			"ibm_is_lb_listener":                                 vpc.ResourceIBMISLBListener(),
			"ibm_is_lb_listener_policy":                          vpc.ResourceIBMISLBListenerPolicy(),
			"ibm_is_lb_listener_policy_rule":                     vpc.ResourceIBMISLBListenerPolicyRule(),
			"ibm_is_lb_pool":                                     vpc.ResourceIBMISLBPool(),
			"ibm_is_lb_pool_member":                              vpc.ResourceIBMISLBPoolMember(),
			"ibm_is_network_acl":                                 vpc.ResourceIBMISNetworkACL(),
			"ibm_is_network_acl_rule":                            vpc.ResourceIBMISNetworkACLRule(),
			"ibm_is_public_gateway":                              vpc.ResourceIBMISPublicGateway(),
			"ibm_is_private_path_service_gateway_account_policy": vpc.ResourceIBMIsPrivatePathServiceGatewayAccountPolicy(),
			"ibm_is_private_path_service_gateway":                vpc.ResourceIBMIsPrivatePathServiceGateway(),
			"ibm_is_private_path_service_gateway_revoke_account": vpc.ResourceIBMIsPrivatePathServiceGatewayRevokeAccount(),
			"ibm_is_private_path_service_gateway_endpoint_gateway_binding_operations": vpc.ResourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperations(),
			"ibm_is_private_path_service_gateway_operations":                          vpc.ResourceIBMIsPrivatePathServiceGatewayOperations(),
			"ibm_is_security_group":                        vpc.ResourceIBMISSecurityGroup(),
			"ibm_is_security_group_rule":                   vpc.ResourceIBMISSecurityGroupRule(),
			"ibm_is_security_group_target":                 vpc.ResourceIBMISSecurityGroupTarget(),
			"ibm_is_share":                                 vpc.ResourceIbmIsShare(),
			"ibm_is_share_replica_operations":              vpc.ResourceIbmIsShareReplicaOperations(),
			"ibm_is_share_mount_target":                    vpc.ResourceIBMIsShareMountTarget(),
			"ibm_is_share_delete_accessor_binding":         vpc.ResourceIbmIsShareDeleteAccessorBinding(),
			"ibm_is_share_snapshot":                        vpc.ResourceIBMIsShareSnapshot(),
			"ibm_is_subnet":                                vpc.ResourceIBMISSubnet(),
			"ibm_is_reservation":                           vpc.ResourceIBMISReservation(),
			"ibm_is_reservation_activate":                  vpc.ResourceIBMISReservationActivate(),
			"ibm_is_subnet_reserved_ip":                    vpc.ResourceIBMISReservedIP(),
			"ibm_is_subnet_reserved_ip_patch":              vpc.ResourceIBMISReservedIPPatch(),
			"ibm_is_subnet_network_acl_attachment":         vpc.ResourceIBMISSubnetNetworkACLAttachment(),
			"ibm_is_subnet_public_gateway_attachment":      vpc.ResourceIBMISSubnetPublicGatewayAttachment(),
			"ibm_is_subnet_routing_table_attachment":       vpc.ResourceIBMISSubnetRoutingTableAttachment(),
			"ibm_is_ssh_key":                               vpc.ResourceIBMISSSHKey(),
			"ibm_is_snapshot":                              vpc.ResourceIBMSnapshot(),
			"ibm_is_virtual_network_interface":             vpc.ResourceIBMIsVirtualNetworkInterface(),
			"ibm_is_virtual_network_interface_floating_ip": vpc.ResourceIBMIsVirtualNetworkInterfaceFloatingIP(),
			"ibm_is_virtual_network_interface_ip":          vpc.ResourceIBMIsVirtualNetworkInterfaceIP(),
			"ibm_is_snapshot_consistency_group":            vpc.ResourceIBMIsSnapshotConsistencyGroup(),
			"ibm_is_volume":                                vpc.ResourceIBMISVolume(),
			"ibm_is_vpn_gateway":                           vpc.ResourceIBMISVPNGateway(),
			"ibm_is_vpn_gateway_connection":                vpc.ResourceIBMISVPNGatewayConnection(),
			"ibm_is_vpc":                                   vpc.ResourceIBMISVPC(),
			"ibm_is_vpc_address_prefix":                    vpc.ResourceIBMISVpcAddressPrefix(),
			"ibm_is_vpc_dns_resolution_binding":            vpc.ResourceIBMIsVPCDnsResolutionBinding(),
			"ibm_is_vpc_routing_table":                     vpc.ResourceIBMISVPCRoutingTable(),
			"ibm_is_vpc_routing_table_route":               vpc.ResourceIBMISVPCRoutingTableRoute(),
			"ibm_is_vpn_server":                            vpc.ResourceIBMIsVPNServer(),
			"ibm_is_vpn_server_client":                     vpc.ResourceIBMIsVPNServerClient(),
			"ibm_is_vpn_server_route":                      vpc.ResourceIBMIsVPNServerRoute(),
			"ibm_is_image":                                 vpc.ResourceIBMISImage(),
			"ibm_is_image_deprecate":                       vpc.ResourceIBMISImageDeprecate(),
			"ibm_is_image_export_job":                      vpc.ResourceIBMIsImageExportJob(),
			"ibm_is_image_obsolete":                        vpc.ResourceIBMISImageObsolete(),
			"ibm_lb":                                       classicinfrastructure.ResourceIBMLb(),
			"ibm_lbaas":                                    classicinfrastructure.ResourceIBMLbaas(),
			"ibm_lbaas_health_monitor":                     classicinfrastructure.ResourceIBMLbaasHealthMonitor(),
			"ibm_lbaas_server_instance_attachment":         classicinfrastructure.ResourceIBMLbaasServerInstanceAttachment(),
			"ibm_lb_service":                               classicinfrastructure.ResourceIBMLbService(),
			"ibm_lb_service_group":                         classicinfrastructure.ResourceIBMLbServiceGroup(),
			"ibm_lb_vpx":                                   classicinfrastructure.ResourceIBMLbVpx(),
			"ibm_lb_vpx_ha":                                classicinfrastructure.ResourceIBMLbVpxHa(),
			"ibm_lb_vpx_service":                           classicinfrastructure.ResourceIBMLbVpxService(),
			"ibm_lb_vpx_vip":                               classicinfrastructure.ResourceIBMLbVpxVip(),
			"ibm_multi_vlan_firewall":                      classicinfrastructure.ResourceIBMMultiVlanFirewall(),
			"ibm_network_gateway":                          classicinfrastructure.ResourceIBMNetworkGateway(),
			"ibm_network_gateway_vlan_association":         classicinfrastructure.ResourceIBMNetworkGatewayVlanAttachment(),
			"ibm_network_interface_sg_attachment":          classicinfrastructure.ResourceIBMNetworkInterfaceSGAttachment(),
			"ibm_network_public_ip":                        classicinfrastructure.ResourceIBMNetworkPublicIp(),
			"ibm_network_vlan":                             classicinfrastructure.ResourceIBMNetworkVlan(),
			"ibm_network_vlan_spanning":                    classicinfrastructure.ResourceIBMNetworkVlanSpan(),
			"ibm_object_storage_account":                   classicinfrastructure.ResourceIBMObjectStorageAccount(),
			"ibm_org":                                      cloudfoundry.ResourceIBMOrg(),
			"ibm_pn_application_chrome":                    pushnotification.ResourceIBMPNApplicationChrome(),
			"ibm_app_config_environment":                   appconfiguration.ResourceIBMAppConfigEnvironment(),
			"ibm_app_config_collection":                    appconfiguration.ResourceIBMAppConfigCollection(),
			"ibm_app_config_feature":                       appconfiguration.ResourceIBMIbmAppConfigFeature(),
			"ibm_app_config_property":                      appconfiguration.ResourceIBMIbmAppConfigProperty(),
			"ibm_app_config_segment":                       appconfiguration.ResourceIBMIbmAppConfigSegment(),
			"ibm_app_config_snapshot":                      appconfiguration.ResourceIBMIbmAppConfigSnapshot(),
			"ibm_kms_key":                                  kms.ResourceIBMKmskey(),
			"ibm_kms_key_with_policy_overrides":            kms.ResourceIBMKmsKeyWithPolicyOverrides(),
			"ibm_kms_key_alias":                            kms.ResourceIBMKmskeyAlias(),
			"ibm_kms_key_rings":                            kms.ResourceIBMKmskeyRings(),
			"ibm_kms_key_policies":                         kms.ResourceIBMKmskeyPolicies(),
			"ibm_kp_key":                                   kms.ResourceIBMkey(),
			"ibm_kms_instance_policies":                    kms.ResourceIBMKmsInstancePolicy(),
			"ibm_kms_kmip_adapter":                         kms.ResourceIBMKmsKMIPAdapter(),
			"ibm_kms_kmip_client_cert":                     kms.ResourceIBMKmsKMIPClientCertificate(),
			"ibm_resource_group":                           resourcemanager.ResourceIBMResourceGroup(),
			"ibm_resource_instance":                        resourcecontroller.ResourceIBMResourceInstance(),
			"ibm_resource_key":                             resourcecontroller.ResourceIBMResourceKey(),
			"ibm_security_group":                           classicinfrastructure.ResourceIBMSecurityGroup(),
			"ibm_security_group_rule":                      classicinfrastructure.ResourceIBMSecurityGroupRule(),
			"ibm_service_instance":                         cloudfoundry.ResourceIBMServiceInstance(),
			"ibm_service_key":                              cloudfoundry.ResourceIBMServiceKey(),
			"ibm_space":                                    cloudfoundry.ResourceIBMSpace(),
			"ibm_storage_evault":                           classicinfrastructure.ResourceIBMStorageEvault(),
			"ibm_storage_block":                            classicinfrastructure.ResourceIBMStorageBlock(),
			"ibm_storage_file":                             classicinfrastructure.ResourceIBMStorageFile(),
			"ibm_subnet":                                   classicinfrastructure.ResourceIBMSubnet(),
			"ibm_dns_reverse_record":                       classicinfrastructure.ResourceIBMDNSReverseRecord(),
			"ibm_ssl_certificate":                          classicinfrastructure.ResourceIBMSSLCertificate(),
			"ibm_cdn":                                      classicinfrastructure.ResourceIBMCDN(),
			"ibm_hardware_firewall_shared":                 classicinfrastructure.ResourceIBMFirewallShared(),

			// Software Defined Storage as a Service
			"ibm_sds_volume": sdsaas.ResourceIBMSdsVolume(),
			"ibm_sds_host":   sdsaas.ResourceIBMSdsHost(),

			// Partner Center Sell
			"ibm_onboarding_registration":       partnercentersell.ResourceIbmOnboardingRegistration(),
			"ibm_onboarding_product":            partnercentersell.ResourceIbmOnboardingProduct(),
			"ibm_onboarding_iam_registration":   partnercentersell.ResourceIbmOnboardingIamRegistration(),
			"ibm_onboarding_catalog_product":    partnercentersell.ResourceIbmOnboardingCatalogProduct(),
			"ibm_onboarding_catalog_plan":       partnercentersell.ResourceIbmOnboardingCatalogPlan(),
			"ibm_onboarding_catalog_deployment": partnercentersell.ResourceIbmOnboardingCatalogDeployment(),
			"ibm_onboarding_resource_broker":    partnercentersell.ResourceIbmOnboardingResourceBroker(),

			// Added for Power Colo
			"ibm_pi_capture":                         power.ResourceIBMPICapture(),
			"ibm_pi_cloud_connection_network_attach": power.ResourceIBMPICloudConnectionNetworkAttach(),
			"ibm_pi_cloud_connection":                power.ResourceIBMPICloudConnection(),
			"ibm_pi_console_language":                power.ResourceIBMPIInstanceConsoleLanguage(),
			"ibm_pi_dhcp":                            power.ResourceIBMPIDhcp(),
			"ibm_pi_host_group":                      power.ResourceIBMPIHostGroup(),
			"ibm_pi_host":                            power.ResourceIBMPIHost(),
			"ibm_pi_ike_policy":                      power.ResourceIBMPIIKEPolicy(),
			"ibm_pi_image_export":                    power.ResourceIBMPIImageExport(),
			"ibm_pi_image":                           power.ResourceIBMPIImage(),
			"ibm_pi_instance_action":                 power.ResourceIBMPIInstanceAction(),
			"ibm_pi_instance":                        power.ResourceIBMPIInstance(),
			"ibm_pi_instance_snapshot":               power.ResourceIBMPIInstanceSnapshot(),
			"ibm_pi_ipsec_policy":                    power.ResourceIBMPIIPSecPolicy(),
			"ibm_pi_key":                             power.ResourceIBMPIKey(),
			"ibm_pi_network_address_group_member":    power.ResourceIBMPINetworkAddressGroupMember(),
			"ibm_pi_network_address_group":           power.ResourceIBMPINetworkAddressGroup(),
			"ibm_pi_network_interface":               power.ResourceIBMPINetworkInterface(),
			"ibm_pi_network_port_attach":             power.ResourceIBMPINetworkPortAttach(),
			"ibm_pi_network_security_group_action":   power.ResourceIBMPINetworkSecurityGroupAction(),
			"ibm_pi_network_security_group_member":   power.ResourceIBMPINetworkSecurityGroupMember(),
			"ibm_pi_network_security_group_rule":     power.ResourceIBMPINetworkSecurityGroupRule(),
			"ibm_pi_network_security_group":          power.ResourceIBMPINetworkSecurityGroup(),
			"ibm_pi_network":                         power.ResourceIBMPINetwork(),
			"ibm_pi_placement_group":                 power.ResourceIBMPIPlacementGroup(),
			"ibm_pi_shared_processor_pool":           power.ResourceIBMPISharedProcessorPool(),
			"ibm_pi_snapshot":                        power.ResourceIBMPISnapshot(),
			"ibm_pi_spp_placement_group":             power.ResourceIBMPISPPPlacementGroup(),
			"ibm_pi_virtual_serial_number":           power.ResourceIBMPIVirtualSerialNumber(),
			"ibm_pi_volume_attach":                   power.ResourceIBMPIVolumeAttach(),
			"ibm_pi_volume_clone":                    power.ResourceIBMPIVolumeClone(),
			"ibm_pi_volume_group_action":             power.ResourceIBMPIVolumeGroupAction(),
			"ibm_pi_volume_group":                    power.ResourceIBMPIVolumeGroup(),
			"ibm_pi_volume_onboarding":               power.ResourceIBMPIVolumeOnboarding(),
			"ibm_pi_volume":                          power.ResourceIBMPIVolume(),
			"ibm_pi_vpn_connection":                  power.ResourceIBMPIVPNConnection(),
			"ibm_pi_workspace":                       power.ResourceIBMPIWorkspace(),

			// Private DNS related resources
			"ibm_dns_zone":              dnsservices.ResourceIBMPrivateDNSZone(),
			"ibm_dns_permitted_network": dnsservices.ResourceIBMPrivateDNSPermittedNetwork(),
			"ibm_dns_resource_record":   dnsservices.ResourceIBMPrivateDNSResourceRecord(),
			"ibm_dns_glb_monitor":       dnsservices.ResourceIBMPrivateDNSGLBMonitor(),
			"ibm_dns_glb_pool":          dnsservices.ResourceIBMPrivateDNSGLBPool(),
			"ibm_dns_glb":               dnsservices.ResourceIBMPrivateDNSGLB(),

			// Added for Custom Resolver
			"ibm_dns_custom_resolver":                 dnsservices.ResourceIBMPrivateDNSCustomResolver(),
			"ibm_dns_custom_resolver_forwarding_rule": dnsservices.ResourceIBMPrivateDNSForwardingRule(),
			"ibm_dns_custom_resolver_secondary_zone":  dnsservices.ResourceIBMPrivateDNSSecondaryZone(),
			"ibm_dns_linked_zone":                     dnsservices.ResourceIBMDNSLinkedZone(),

			// Direct Link related resources
			"ibm_dl_gateway":            directlink.ResourceIBMDLGateway(),
			"ibm_dl_virtual_connection": directlink.ResourceIBMDLGatewayVC(),
			"ibm_dl_provider_gateway":   directlink.ResourceIBMDLProviderGateway(),
			"ibm_dl_route_report":       directlink.ResourceIBMDLGatewayRouteReport(),
			"ibm_dl_gateway_action":     directlink.ResourceIBMDLGatewayAction(),

			// Added for Transit Gateway
			"ibm_tg_gateway":                  transitgateway.ResourceIBMTransitGateway(),
			"ibm_tg_connection":               transitgateway.ResourceIBMTransitGatewayConnection(),
			"ibm_tg_connection_action":        transitgateway.ResourceIBMTransitGatewayConnectionAction(),
			"ibm_tg_connection_prefix_filter": transitgateway.ResourceIBMTransitGatewayConnectionPrefixFilter(),
			"ibm_tg_route_report":             transitgateway.ResourceIBMTransitGatewayRouteReport(),
			"ibm_tg_connection_rgre_tunnel":   transitgateway.ResourceIBMTransitGatewayConnectionRgreTunnel(),

			// Catalog related resources
			"ibm_cm_offering_instance": catalogmanagement.ResourceIBMCmOfferingInstance(),
			"ibm_cm_catalog":           catalogmanagement.ResourceIBMCmCatalog(),
			"ibm_cm_offering":          catalogmanagement.ResourceIBMCmOffering(),
			"ibm_cm_version":           catalogmanagement.ResourceIBMCmVersion(),
			"ibm_cm_validation":        catalogmanagement.ResourceIBMCmValidation(),
			"ibm_cm_object":            catalogmanagement.ResourceIBMCmObject(),
			"ibm_cm_account":           catalogmanagement.ResourceIBMCmAccount(),

			// Added for enterprise
			"ibm_enterprise":               enterprise.ResourceIBMEnterprise(),
			"ibm_enterprise_account_group": enterprise.ResourceIBMEnterpriseAccountGroup(),
			"ibm_enterprise_account":       enterprise.ResourceIBMEnterpriseAccount(),

			// //Added for Usage Reports
			"ibm_billing_report_snapshot": usagereports.ResourceIBMBillingReportSnapshot(),

			// Added for Schematics
			"ibm_schematics_workspace":      schematics.ResourceIBMSchematicsWorkspace(),
			"ibm_schematics_action":         schematics.ResourceIBMSchematicsAction(),
			"ibm_schematics_job":            schematics.ResourceIBMSchematicsJob(),
			"ibm_schematics_inventory":      schematics.ResourceIBMSchematicsInventory(),
			"ibm_schematics_resource_query": schematics.ResourceIBMSchematicsResourceQuery(),
			"ibm_schematics_policy":         schematics.ResourceIbmSchematicsPolicy(),
			"ibm_schematics_agent":          schematics.ResourceIbmSchematicsAgent(),
			"ibm_schematics_agent_prs":      schematics.ResourceIbmSchematicsAgentPrs(),
			"ibm_schematics_agent_deploy":   schematics.ResourceIbmSchematicsAgentDeploy(),
			"ibm_schematics_agent_health":   schematics.ResourceIbmSchematicsAgentHealth(),

			// Added for Secrets Manager
			"ibm_sm_secret_group":                                                secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmSecretGroup()),
			"ibm_sm_arbitrary_secret":                                            secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmArbitrarySecret()),
			"ibm_sm_imported_certificate":                                        secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmImportedCertificate()),
			"ibm_sm_public_certificate":                                          secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPublicCertificate()),
			"ibm_sm_private_certificate":                                         secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificate()),
			"ibm_sm_iam_credentials_secret":                                      secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmIamCredentialsSecret()),
			"ibm_sm_service_credentials_secret":                                  secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmServiceCredentialsSecret()),
			"ibm_sm_username_password_secret":                                    secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmUsernamePasswordSecret()),
			"ibm_sm_kv_secret":                                                   secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmKvSecret()),
			"ibm_sm_public_certificate_configuration_ca_lets_encrypt":            secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPublicCertificateConfigurationCALetsEncrypt()),
			"ibm_sm_public_certificate_configuration_dns_cis":                    secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmConfigurationPublicCertificateDNSCis()),
			"ibm_sm_public_certificate_configuration_dns_classic_infrastructure": secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructure()),
			"ibm_sm_private_certificate_configuration_root_ca":                   secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificateConfigurationRootCA()),
			"ibm_sm_private_certificate_configuration_intermediate_ca":           secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificateConfigurationIntermediateCA()),
			"ibm_sm_private_certificate_configuration_template":                  secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificateConfigurationTemplate()),
			"ibm_sm_iam_credentials_configuration":                               secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmIamCredentialsConfiguration()),
			"ibm_sm_public_certificate_action_validate_manual_dns":               secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPublicCertificateActionValidateManualDns()),
			"ibm_sm_en_registration":                                             secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmEnRegistration()),
			"ibm_sm_private_certificate_configuration_action_sign_csr":           secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificateConfigurationActionSignCsr()),
			"ibm_sm_private_certificate_configuration_action_set_signed":         secretsmanager.AddInstanceFields(secretsmanager.ResourceIbmSmPrivateCertificateConfigurationActionSetSigned()),

			// satellite  resources
			"ibm_satellite_location":                            satellite.ResourceIBMSatelliteLocation(),
			"ibm_satellite_host":                                satellite.ResourceIBMSatelliteHost(),
			"ibm_satellite_cluster":                             satellite.ResourceIBMSatelliteCluster(),
			"ibm_satellite_cluster_worker_pool":                 satellite.ResourceIBMSatelliteClusterWorkerPool(),
			"ibm_satellite_link":                                satellite.ResourceIBMSatelliteLink(),
			"ibm_satellite_storage_configuration":               satellite.ResourceIBMSatelliteStorageConfiguration(),
			"ibm_satellite_storage_assignment":                  satellite.ResourceIBMSatelliteStorageAssignment(),
			"ibm_satellite_endpoint":                            satellite.ResourceIBMSatelliteEndpoint(),
			"ibm_satellite_location_nlb_dns":                    satellite.ResourceIBMSatelliteLocationNlbDns(),
			"ibm_satellite_cluster_worker_pool_zone_attachment": satellite.ResourceIbmSatelliteClusterWorkerPoolZoneAttachment(),

			// Added for Resource Tag
			"ibm_resource_tag": globaltagging.ResourceIBMResourceTag(),

			// Added for Iam Access Tag
			"ibm_iam_access_tag": globaltagging.ResourceIBMIamAccessTag(),

			// Atracker
			"ibm_atracker_target":   atracker.ResourceIBMAtrackerTarget(),
			"ibm_atracker_route":    atracker.ResourceIBMAtrackerRoute(),
			"ibm_atracker_settings": atracker.ResourceIBMAtrackerSettings(),

			// Metrics Router
			"ibm_metrics_router_target":   metricsrouter.ResourceIBMMetricsRouterTarget(),
			"ibm_metrics_router_route":    metricsrouter.ResourceIBMMetricsRouterRoute(),
			"ibm_metrics_router_settings": metricsrouter.ResourceIBMMetricsRouterSettings(),

			// MQ on Cloud
			"ibm_mqcloud_queue_manager":                    mqcloud.ResourceIbmMqcloudQueueManager(),
			"ibm_mqcloud_application":                      mqcloud.ResourceIbmMqcloudApplication(),
			"ibm_mqcloud_user":                             mqcloud.ResourceIbmMqcloudUser(),
			"ibm_mqcloud_keystore_certificate":             mqcloud.ResourceIbmMqcloudKeystoreCertificate(),
			"ibm_mqcloud_truststore_certificate":           mqcloud.ResourceIbmMqcloudTruststoreCertificate(),
			"ibm_mqcloud_virtual_private_endpoint_gateway": mqcloud.ResourceIbmMqcloudVirtualPrivateEndpointGateway(),

			// Security and Compliance Center(soon to be deprecated)
			"ibm_scc_account_settings":    scc.ResourceIBMSccAccountSettings(),
			"ibm_scc_rule_attachment":     scc.ResourceIBMSccRuleAttachment(),
			"ibm_scc_template":            scc.ResourceIBMSccTemplate(),
			"ibm_scc_template_attachment": scc.ResourceIBMSccTemplateAttachment(),

			// Security and Compliance Center
			"ibm_scc_instance_settings":      scc.ResourceIbmSccInstanceSettings(),
			"ibm_scc_rule":                   scc.ResourceIbmSccRule(),
			"ibm_scc_control_library":        scc.ResourceIbmSccControlLibrary(),
			"ibm_scc_profile":                scc.ResourceIbmSccProfile(),
			"ibm_scc_profile_attachment":     scc.ResourceIbmSccProfileAttachment(),
			"ibm_scc_provider_type_instance": scc.ResourceIbmSccProviderTypeInstance(),
			"ibm_scc_scope":                  scc.ResourceIbmSccScope(),

			// Security Services
			"ibm_pag_instance": pag.ResourceIBMPag(),

			// Added for Context Based Restrictions
			"ibm_cbr_zone":           contextbasedrestrictions.ResourceIBMCbrZone(),
			"ibm_cbr_zone_addresses": contextbasedrestrictions.ResourceIBMCbrZoneAddresses(),
			"ibm_cbr_rule":           contextbasedrestrictions.ResourceIBMCbrRule(),

			// Added for Event Notifications
			"ibm_en_source":                     eventnotification.ResourceIBMEnSource(),
			"ibm_en_topic":                      eventnotification.ResourceIBMEnTopic(),
			"ibm_en_destination_webhook":        eventnotification.ResourceIBMEnWebhookDestination(),
			"ibm_en_destination_android":        eventnotification.ResourceIBMEnFCMDestination(),
			"ibm_en_destination_chrome":         eventnotification.ResourceIBMEnChromeDestination(),
			"ibm_en_destination_firefox":        eventnotification.ResourceIBMEnFirefoxDestination(),
			"ibm_en_destination_ios":            eventnotification.ResourceIBMEnAPNSDestination(),
			"ibm_en_destination_slack":          eventnotification.ResourceIBMEnSlackDestination(),
			"ibm_en_subscription_sms":           eventnotification.ResourceIBMEnSMSSubscription(),
			"ibm_en_subscription_email":         eventnotification.ResourceIBMEnEmailSubscription(),
			"ibm_en_subscription_webhook":       eventnotification.ResourceIBMEnWebhookSubscription(),
			"ibm_en_subscription_android":       eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_subscription_ios":           eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_subscription_chrome":        eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_subscription_firefox":       eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_subscription_slack":         eventnotification.ResourceIBMEnSlackSubscription(),
			"ibm_en_subscription_safari":        eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_destination_safari":         eventnotification.ResourceIBMEnSafariDestination(),
			"ibm_en_destination_msteams":        eventnotification.ResourceIBMEnMSTeamsDestination(),
			"ibm_en_subscription_msteams":       eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_destination_pagerduty":      eventnotification.ResourceIBMEnPagerDutyDestination(),
			"ibm_en_subscription_pagerduty":     eventnotification.ResourceIBMEnPagerDutySubscription(),
			"ibm_en_integration":                eventnotification.ResourceIBMEnIntegration(),
			"ibm_en_destination_sn":             eventnotification.ResourceIBMEnServiceNowDestination(),
			"ibm_en_subscription_sn":            eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_destination_ce":             eventnotification.ResourceIBMEnCodeEngineDestination(),
			"ibm_en_subscription_ce":            eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_destination_cos":            eventnotification.ResourceIBMEnCOSDestination(),
			"ibm_en_subscription_cos":           eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_destination_huawei":         eventnotification.ResourceIBMEnHuaweiDestination(),
			"ibm_en_subscription_huawei":        eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_ibmsource":                  eventnotification.ResourceIBMEnIBMSource(),
			"ibm_en_destination_custom_email":   eventnotification.ResourceIBMEnCustomEmailDestination(),
			"ibm_en_subscription_custom_email":  eventnotification.ResourceIBMEnCustomEmailSubscription(),
			"ibm_en_email_template":             eventnotification.ResourceIBMEnEmailTemplate(),
			"ibm_en_integration_cos":            eventnotification.ResourceIBMEnCOSIntegration(),
			"ibm_en_destination_custom_sms":     eventnotification.ResourceIBMEnCustomSMSDestination(),
			"ibm_en_subscription_custom_sms":    eventnotification.ResourceIBMEnCustomSMSSubscription(),
			"ibm_en_smtp_configuration":         eventnotification.ResourceIBMEnSMTPConfiguration(),
			"ibm_en_smtp_user":                  eventnotification.ResourceIBMEnSMTPUser(),
			"ibm_en_slack_template":             eventnotification.ResourceIBMEnSlackTemplate(),
			"ibm_en_smtp_setting":               eventnotification.ResourceIBMEnSMTPSetting(),
			"ibm_en_webhook_template":           eventnotification.ResourceIBMEnWebhookTemplate(),
			"ibm_en_subscription_scheduler":     eventnotification.ResourceIBMEnFCMSubscription(),
			"ibm_en_pagerduty_template":         eventnotification.ResourceIBMEnPagerDutyTemplate(),
			"ibm_en_destination_event_streams":  eventnotification.ResourceIBMEnEventStreamsDestination(),
			"ibm_en_subscription_event_streams": eventnotification.ResourceIBMEnEventStreamsSubscription(),
			"ibm_en_event_streams_template":     eventnotification.ResourceIBMEnEventStreamsTemplate(),

			// Added for Toolchain
			"ibm_cd_toolchain":                         cdtoolchain.ResourceIBMCdToolchain(),
			"ibm_cd_toolchain_tool_keyprotect":         cdtoolchain.ResourceIBMCdToolchainToolKeyprotect(),
			"ibm_cd_toolchain_tool_secretsmanager":     cdtoolchain.ResourceIBMCdToolchainToolSecretsmanager(),
			"ibm_cd_toolchain_tool_bitbucketgit":       cdtoolchain.ResourceIBMCdToolchainToolBitbucketgit(),
			"ibm_cd_toolchain_tool_githubconsolidated": cdtoolchain.ResourceIBMCdToolchainToolGithubconsolidated(),
			"ibm_cd_toolchain_tool_gitlab":             cdtoolchain.ResourceIBMCdToolchainToolGitlab(),
			"ibm_cd_toolchain_tool_hostedgit":          cdtoolchain.ResourceIBMCdToolchainToolHostedgit(),
			"ibm_cd_toolchain_tool_artifactory":        cdtoolchain.ResourceIBMCdToolchainToolArtifactory(),
			"ibm_cd_toolchain_tool_custom":             cdtoolchain.ResourceIBMCdToolchainToolCustom(),
			"ibm_cd_toolchain_tool_pipeline":           cdtoolchain.ResourceIBMCdToolchainToolPipeline(),
			"ibm_cd_toolchain_tool_devopsinsights":     cdtoolchain.ResourceIBMCdToolchainToolDevopsinsights(),
			"ibm_cd_toolchain_tool_slack":              cdtoolchain.ResourceIBMCdToolchainToolSlack(),
			"ibm_cd_toolchain_tool_sonarqube":          cdtoolchain.ResourceIBMCdToolchainToolSonarqube(),
			"ibm_cd_toolchain_tool_hashicorpvault":     cdtoolchain.ResourceIBMCdToolchainToolHashicorpvault(),
			"ibm_cd_toolchain_tool_securitycompliance": cdtoolchain.ResourceIBMCdToolchainToolSecuritycompliance(),
			"ibm_cd_toolchain_tool_privateworker":      cdtoolchain.ResourceIBMCdToolchainToolPrivateworker(),
			"ibm_cd_toolchain_tool_appconfig":          cdtoolchain.ResourceIBMCdToolchainToolAppconfig(),
			"ibm_cd_toolchain_tool_jenkins":            cdtoolchain.ResourceIBMCdToolchainToolJenkins(),
			"ibm_cd_toolchain_tool_nexus":              cdtoolchain.ResourceIBMCdToolchainToolNexus(),
			"ibm_cd_toolchain_tool_pagerduty":          cdtoolchain.ResourceIBMCdToolchainToolPagerduty(),
			"ibm_cd_toolchain_tool_saucelabs":          cdtoolchain.ResourceIBMCdToolchainToolSaucelabs(),
			"ibm_cd_toolchain_tool_jira":               cdtoolchain.ResourceIBMCdToolchainToolJira(),
			"ibm_cd_toolchain_tool_eventnotifications": cdtoolchain.ResourceIBMCdToolchainToolEventnotifications(),

			// Added for Tekton Pipeline
			"ibm_cd_tekton_pipeline_definition":       cdtektonpipeline.ResourceIBMCdTektonPipelineDefinition(),
			"ibm_cd_tekton_pipeline_trigger_property": cdtektonpipeline.ResourceIBMCdTektonPipelineTriggerProperty(),
			"ibm_cd_tekton_pipeline_property":         cdtektonpipeline.ResourceIBMCdTektonPipelineProperty(),
			"ibm_cd_tekton_pipeline_trigger":          cdtektonpipeline.ResourceIBMCdTektonPipelineTrigger(),
			"ibm_cd_tekton_pipeline":                  cdtektonpipeline.ResourceIBMCdTektonPipeline(),

			// Added for Code Engine
			"ibm_code_engine_allowed_outbound_destination": codeengine.ResourceIbmCodeEngineAllowedOutboundDestination(),
			"ibm_code_engine_app":                          codeengine.ResourceIbmCodeEngineApp(),
			"ibm_code_engine_binding":                      codeengine.ResourceIbmCodeEngineBinding(),
			"ibm_code_engine_build":                        codeengine.ResourceIbmCodeEngineBuild(),
			"ibm_code_engine_config_map":                   codeengine.ResourceIbmCodeEngineConfigMap(),
			"ibm_code_engine_domain_mapping":               codeengine.ResourceIbmCodeEngineDomainMapping(),
			"ibm_code_engine_function":                     codeengine.ResourceIbmCodeEngineFunction(),
			"ibm_code_engine_job":                          codeengine.ResourceIbmCodeEngineJob(),
			"ibm_code_engine_project":                      codeengine.ResourceIbmCodeEngineProject(),
			"ibm_code_engine_secret":                       codeengine.ResourceIbmCodeEngineSecret(),

			// Added for Project
			"ibm_project":             project.ResourceIbmProject(),
			"ibm_project_config":      project.ResourceIbmProjectConfig(),
			"ibm_project_environment": project.ResourceIbmProjectEnvironment(),

			// Added for VMware as a Service
			"ibm_vmaas_vdc": vmware.ResourceIbmVmaasVdc(),
			// Logs Service
			"ibm_logs_alert":              logs.AddLogsInstanceFields(logs.ResourceIbmLogsAlert()),
			"ibm_logs_rule_group":         logs.AddLogsInstanceFields(logs.ResourceIbmLogsRuleGroup()),
			"ibm_logs_policy":             logs.AddLogsInstanceFields(logs.ResourceIbmLogsPolicy()),
			"ibm_logs_dashboard":          logs.AddLogsInstanceFields(logs.ResourceIbmLogsDashboard()),
			"ibm_logs_e2m":                logs.AddLogsInstanceFields(logs.ResourceIbmLogsE2m()),
			"ibm_logs_outgoing_webhook":   logs.AddLogsInstanceFields(logs.ResourceIbmLogsOutgoingWebhook()),
			"ibm_logs_view_folder":        logs.AddLogsInstanceFields(logs.ResourceIbmLogsViewFolder()),
			"ibm_logs_view":               logs.AddLogsInstanceFields(logs.ResourceIbmLogsView()),
			"ibm_logs_dashboard_folder":   logs.AddLogsInstanceFields(logs.ResourceIbmLogsDashboardFolder()),
			"ibm_logs_data_usage_metrics": logs.AddLogsInstanceFields(logs.ResourceIbmLogsDataUsageMetrics()),
			"ibm_logs_enrichment":         logs.AddLogsInstanceFields(logs.ResourceIbmLogsEnrichment()),
			"ibm_logs_data_access_rule":   logs.AddLogsInstanceFields(logs.ResourceIbmLogsDataAccessRule()),
			"ibm_logs_stream":             logs.AddLogsInstanceFields(logs.ResourceIbmLogsStream()),

			// Logs Router Service
			"ibm_logs_router_tenant": logsrouting.ResourceIBMLogsRouterTenant(),
		},

		ConfigureFunc: providerConfigure,
	}

	wrappedProvider := wrapProvider(provider)
	return &wrappedProvider
}

func wrapProvider(provider schema.Provider) schema.Provider {
	wrappedResourcesMap := map[string]*schema.Resource{}
	wrappedDataSourcesMap := map[string]*schema.Resource{}

	for key, value := range provider.ResourcesMap {
		wrappedResourcesMap[key] = wrapResource(key, value)
	}

	for key, value := range provider.DataSourcesMap {
		wrappedDataSourcesMap[key] = wrapDataSource(key, value)
	}

	return schema.Provider{
		Schema:         provider.Schema,
		DataSourcesMap: wrappedDataSourcesMap,
		ResourcesMap:   wrappedResourcesMap,
		ConfigureFunc:  provider.ConfigureFunc,
	}
}

func wrapResource(name string, resource *schema.Resource) *schema.Resource {
	return &schema.Resource{
		Schema:               resource.Schema,
		SchemaVersion:        resource.SchemaVersion,
		MigrateState:         resource.MigrateState,
		StateUpgraders:       resource.StateUpgraders,
		Exists:               resource.Exists,
		CreateContext:        wrapFunction(name, "create", resource.CreateContext, resource.Create, false),
		ReadContext:          wrapFunction(name, "read", resource.ReadContext, resource.Read, false),
		UpdateContext:        wrapFunction(name, "update", resource.UpdateContext, resource.Update, false),
		DeleteContext:        wrapFunction(name, "delete", resource.DeleteContext, resource.Delete, false),
		CreateWithoutTimeout: wrapFunction(name, "create", resource.CreateWithoutTimeout, nil, false),
		ReadWithoutTimeout:   wrapFunction(name, "read", resource.ReadWithoutTimeout, nil, false),
		UpdateWithoutTimeout: wrapFunction(name, "update", resource.UpdateWithoutTimeout, nil, false),
		DeleteWithoutTimeout: wrapFunction(name, "delete", resource.DeleteWithoutTimeout, nil, false),
		CustomizeDiff:        wrapCustomizeDiff(name, resource.CustomizeDiff),
		Importer:             resource.Importer,
		DeprecationMessage:   resource.DeprecationMessage,
		Timeouts:             resource.Timeouts,
		Description:          resource.Description,
		UseJSONNumber:        resource.UseJSONNumber,
	}
}

func wrapDataSource(name string, resource *schema.Resource) *schema.Resource {
	return &schema.Resource{
		Schema:             resource.Schema,
		SchemaVersion:      resource.SchemaVersion,
		MigrateState:       resource.MigrateState,
		StateUpgraders:     resource.StateUpgraders,
		Exists:             resource.Exists,
		ReadContext:        wrapFunction(name, "read", resource.ReadContext, resource.Read, true),
		ReadWithoutTimeout: wrapFunction(name, "read", resource.ReadWithoutTimeout, nil, true),
		Importer:           resource.Importer,
		DeprecationMessage: resource.DeprecationMessage,
		Timeouts:           resource.Timeouts,
		Description:        resource.Description,
		UseJSONNumber:      resource.UseJSONNumber,
	}
}

func wrapFunction(
	resourceName, operationName string,
	function func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics,
	fallback func(*schema.ResourceData, interface{}) error,
	isDataSource bool,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	if function != nil {
		return func(context context.Context, schema *schema.ResourceData, meta interface{}) diag.Diagnostics {

			// only allow deletion if the resource is not marked as protected
			if operationName == "delete" && schema.Get("deletion_protection") != nil {
				// we check the value in state, not current config. Current config will always be null for a delete

				if schema.Get("deletion_protection") == true {
					log.Printf("[DEBUG] Resource has deletion protection turned on %s", resourceName)
					var diags diag.Diagnostics
					summary := fmt.Sprintf("Deletion protection is enabled for resource %s to prevent accidential deletion", schema.Get("name"))
					return append(
						diags,
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  summary,
							Detail:   "Set deletion_protection to false, apply and then destroy if deletion should proceed",
						},
					)
				}
			}

			return function(context, schema, meta)
		}
	} else if fallback != nil {
		return func(context context.Context, schema *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return wrapError(fallback(schema, meta), resourceName, operationName, isDataSource)
		}
	}

	return nil
}

func wrapError(err error, resourceName, operationName string, isDataSource bool) diag.Diagnostics {
	if err == nil {
		return nil
	}

	var diags diag.Diagnostics

	// Distinguish data sources from resources. Data sources technically are resources but
	// they may have the same names and we need to tell them apart.
	if isDataSource {
		resourceName = fmt.Sprintf("(Data) %s", resourceName)
	}

	var tfError *flex.TerraformProblem
	if errors.As(err, &tfError) {
		tfError.Resource = resourceName
		tfError.Operation = operationName
	} else {
		tfError = flex.TerraformErrorf(err, "", resourceName, operationName)
	}

	log.Printf("[DEBUG] %s", tfError.GetDebugMessage())
	return append(
		diags,
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  tfError.Error(),
			Detail:   tfError.GetConsoleMessage(),
		},
	)
}

func wrapCustomizeDiff(resourceName string, function schema.CustomizeDiffFunc) schema.CustomizeDiffFunc {
	if function == nil {
		return nil
	}

	return func(c context.Context, rd *schema.ResourceDiff, i interface{}) error {
		return wrapDiffErrors(function(c, rd, i), resourceName)
	}
}

func wrapDiffErrors(err error, resourceName string) error {
	if err != nil {
		// CustomizeDiff fields often use the customizediff.All() method, which concatenates the errors
		// returned from multiple functions using errors.Join(). Individual errors are still embedded in the
		// error and will be extracted when the error is unwrapped by the Go core.
		tfError := flex.TerraformErrorf(err, err.Error(), resourceName, "CustomizeDiff")

		// By the time this error gets printed by the Terraform code, we've lost control of it and the
		// message that gets printed comes from the Error() method (and we only see the Summary).
		// Although it would be ideal to return the full TerraformError object, it is sufficient
		// to package the console message into a new error so that the user gets the information.
		log.Printf("[DEBUG] %s", tfError.GetDebugMessage())
		return errors.New(tfError.GetConsoleMessage())
	}

	// Return the nil error.
	return err
}

var (
	globalValidatorDict validate.ValidatorDict
	initOnce            sync.Once
)

func init() {
	validate.SetValidatorDict(Validator())
}

// Validator return validator
func Validator() validate.ValidatorDict {
	initOnce.Do(func() {
		globalValidatorDict = validate.ValidatorDict{
			ResourceValidatorDictionary: map[string]*validate.ResourceValidator{
				"ibm_iam_trusted_profile_template_assignment":  iamidentity.ResourceIBMTrustedProfileTemplateAssignmentValidator(),
				"ibm_iam_account_settings_template_assignment": iamidentity.ResourceIBMAccountSettingsTemplateAssignmentValidator(),
				"ibm_iam_account_settings":                     iamidentity.ResourceIBMIAMAccountSettingsValidator(),
				"ibm_iam_custom_role":                          iampolicy.ResourceIBMIAMCustomRoleValidator(),
				"ibm_cis_healthcheck":                          cis.ResourceIBMCISHealthCheckValidator(),
				"ibm_cis_rate_limit":                           cis.ResourceIBMCISRateLimitValidator(),
				"ibm_cis":                                      cis.ResourceIBMCISValidator(),
				"ibm_cis_domain_settings":                      cis.ResourceIBMCISDomainSettingValidator(),
				"ibm_cis_domain":                               cis.ResourceIBMCISDomainValidator(),
				"ibm_cis_tls_settings":                         cis.ResourceIBMCISTLSSettingsValidator(),
				"ibm_cis_routing":                              cis.ResourceIBMCISRoutingValidator(),
				"ibm_cis_page_rule":                            cis.ResourceIBMCISPageRuleValidator(),
				"ibm_cis_waf_package":                          cis.ResourceIBMCISWAFPackageValidator(),
				"ibm_cis_waf_group":                            cis.ResourceIBMCISWAFGroupValidator(),
				"ibm_cis_certificate_upload":                   cis.ResourceIBMCISCertificateUploadValidator(),
				"ibm_cis_cache_settings":                       cis.ResourceIBMCISCacheSettingsValidator(),
				"ibm_cis_custom_page":                          cis.ResourceIBMCISCustomPageValidator(),
				"ibm_cis_firewall":                             cis.ResourceIBMCISFirewallValidator(),
				"ibm_cis_range_app":                            cis.ResourceIBMCISRangeAppValidator(),
				"ibm_cis_waf_rule":                             cis.ResourceIBMCISWAFRuleValidator(),
				"ibm_cis_certificate_order":                    cis.ResourceIBMCISCertificateOrderValidator(),
				"ibm_cis_filter":                               cis.ResourceIBMCISFilterValidator(),
				"ibm_cis_firewall_rules":                       cis.ResourceIBMCISFirewallrulesValidator(),
				"ibm_cis_webhook":                              cis.ResourceIBMCISWebhooksValidator(),
				"ibm_cis_alert":                                cis.ResourceIBMCISAlertValidator(),
				"ibm_cis_dns_record":                           cis.ResourceIBMCISDnsRecordValidator(),
				"ibm_cis_dns_records_import":                   cis.ResourceIBMCISDnsRecordsImportValidator(),
				"ibm_cis_edge_functions_action":                cis.ResourceIBMCISEdgeFunctionsActionValidator(),
				"ibm_cis_edge_functions_trigger":               cis.ResourceIBMCISEdgeFunctionsTriggerValidator(),
				"ibm_cis_global_load_balancer":                 cis.ResourceIBMCISGlbValidator(),
				"ibm_cis_logpush_job":                          cis.ResourceIBMCISLogPushJobValidator(),
				"ibm_cis_mtls_app":                             cis.ResourceIBMCISMtlsAppValidator(),
				"ibm_cis_mtls":                                 cis.ResourceIBMCISMtlsValidator(),
				"ibm_cis_bot_management":                       cis.ResourceIBMCISBotManagementValidator(),
				"ibm_cis_origin_auth":                          cis.ResourceIBMCISOriginAuthPullValidator(),
				"ibm_cis_origin_pool":                          cis.ResourceIBMCISPoolValidator(),
				"ibm_cis_ruleset":                              cis.ResourceIBMCISRulesetValidator(),
				"ibm_cis_ruleset_entrypoint_version":           cis.ResourceIBMCISRulesetEntryPointVersionValidator(),
				"ibm_cis_ruleset_rule":                         cis.ResourceIBMCISRulesetRuleValidator(),
				"ibm_cis_ruleset_version_detach":               cis.ResourceIBMCISRulesetVersionDetachValidator(),
				"ibm_cis_advanced_certificate_pack_order":      cis.ResourceIBMCISAdvancedCertificatePackOrderValidator(),
				"ibm_cis_origin_certificate_order":             cis.ResourceIBMCISOriginCertificateOrderValidator(),
				"ibm_container_cluster":                        kubernetes.ResourceIBMContainerClusterValidator(),
				"ibm_container_worker_pool":                    kubernetes.ResourceIBMContainerWorkerPoolValidator(),
				"ibm_container_vpc_worker_pool":                kubernetes.ResourceIBMContainerVPCWorkerPoolValidator(),
				"ibm_container_vpc_worker":                     kubernetes.ResourceIBMContainerVPCWorkerValidator(),
				"ibm_container_vpc_cluster":                    kubernetes.ResourceIBMContainerVpcClusterValidator(),
				"ibm_cos_bucket":                               cos.ResourceIBMCOSBucketValidator(),
				"ibm_cr_namespace":                             registry.ResourceIBMCrNamespaceValidator(),
				"ibm_tg_gateway":                               transitgateway.ResourceIBMTGValidator(),
				"ibm_app_config_feature":                       appconfiguration.ResourceIBMAppConfigFeatureValidator(),
				"ibm_tg_connection":                            transitgateway.ResourceIBMTransitGatewayConnectionValidator(),
				"ibm_tg_connection_action":                     transitgateway.ResourceIBMTransitGatewayConnectionActionValidator(),
				"ibm_tg_connection_prefix_filter":              transitgateway.ResourceIBMTransitGatewayConnectionPrefixFilterValidator(),
				"ibm_tg_connection_rgre_tunnel":                transitgateway.ResourceIBMTransitGatewayConnectionRgreTunnelValidator(),
				"ibm_dl_virtual_connection":                    directlink.ResourceIBMDLGatewayVCValidator(),
				"ibm_dl_gateway":                               directlink.ResourceIBMDLGatewayValidator(),
				"ibm_dl_provider_gateway":                      directlink.ResourceIBMDLProviderGatewayValidator(),
				"ibm_dl_gateway_action":                        directlink.ResourceIBMDLGatewayActionValidator(),
				"ibm_database":                                 database.ResourceIBMICDValidator(),
				"ibm_function_package":                         functions.ResourceIBMFuncPackageValidator(),
				"ibm_function_action":                          functions.ResourceIBMFuncActionValidator(),
				"ibm_function_rule":                            functions.ResourceIBMFuncRuleValidator(),
				"ibm_function_trigger":                         functions.ResourceIBMFuncTriggerValidator(),
				"ibm_function_namespace":                       functions.ResourceIBMFuncNamespaceValidator(),
				"ibm_hpcs":                                     hpcs.ResourceIBMHPCSValidator(),
				"ibm_hpcs_managed_key":                         hpcs.ResourceIbmManagedKeyValidator(),
				"ibm_hpcs_keystore":                            hpcs.ResourceIbmKeystoreValidator(),
				"ibm_hpcs_key_template":                        hpcs.ResourceIbmKeyTemplateValidator(),
				"ibm_hpcs_vault":                               hpcs.ResourceIbmVaultValidator(),
				"ibm_config_aggregator_settings":               configurationaggregator.ResourceIbmConfigAggregatorSettingsValidator(),

				// Cloudshell
				"ibm_cloud_shell_account_settings": cloudshell.ResourceIBMCloudShellAccountSettingsValidator(),

				// MQ on Cloud
				"ibm_mqcloud_queue_manager":                    mqcloud.ResourceIbmMqcloudQueueManagerValidator(),
				"ibm_mqcloud_application":                      mqcloud.ResourceIbmMqcloudApplicationValidator(),
				"ibm_mqcloud_user":                             mqcloud.ResourceIbmMqcloudUserValidator(),
				"ibm_mqcloud_keystore_certificate":             mqcloud.ResourceIbmMqcloudKeystoreCertificateValidator(),
				"ibm_mqcloud_truststore_certificate":           mqcloud.ResourceIbmMqcloudTruststoreCertificateValidator(),
				"ibm_mqcloud_virtual_private_endpoint_gateway": mqcloud.ResourceIbmMqcloudVirtualPrivateEndpointGatewayValidator(),

				"ibm_is_backup_policy":      vpc.ResourceIBMIsBackupPolicyValidator(),
				"ibm_is_backup_policy_plan": vpc.ResourceIBMIsBackupPolicyPlanValidator(),

				// bare_metal_server
				"ibm_is_bare_metal_server_disk":               vpc.ResourceIBMIsBareMetalServerDiskValidator(),
				"ibm_is_bare_metal_server_network_attachment": vpc.ResourceIBMIsBareMetalServerNetworkAttachmentValidator(),
				"ibm_is_bare_metal_server_network_interface":  vpc.ResourceIBMIsBareMetalServerNetworkInterfaceValidator(),
				"ibm_is_bare_metal_server":                    vpc.ResourceIBMIsBareMetalServerValidator(),

				// cluster

				"ibm_is_cluster_network_interface":           vpc.ResourceIBMIsClusterNetworkInterfaceValidator(),
				"ibm_is_cluster_network_subnet":              vpc.ResourceIBMIsClusterNetworkSubnetValidator(),
				"ibm_is_cluster_network_subnet_reserved_ip":  vpc.ResourceIBMIsClusterNetworkSubnetReservedIPValidator(),
				"ibm_is_cluster_network":                     vpc.ResourceIBMIsClusterNetworkValidator(),
				"ibm_is_instance_cluster_network_attachment": vpc.ResourceIBMIsInstanceClusterNetworkAttachmentValidator(),

				"ibm_is_dedicated_host_group":                        vpc.ResourceIbmIsDedicatedHostGroupValidator(),
				"ibm_is_dedicated_host":                              vpc.ResourceIbmIsDedicatedHostValidator(),
				"ibm_is_dedicated_host_disk_management":              vpc.ResourceIBMISDedicatedHostDiskManagementValidator(),
				"ibm_is_flow_log":                                    vpc.ResourceIBMISFlowLogValidator(),
				"ibm_is_instance_group":                              vpc.ResourceIBMISInstanceGroupValidator(),
				"ibm_is_instance_group_membership":                   vpc.ResourceIBMISInstanceGroupMembershipValidator(),
				"ibm_is_instance_group_manager":                      vpc.ResourceIBMISInstanceGroupManagerValidator(),
				"ibm_is_instance_group_manager_policy":               vpc.ResourceIBMISInstanceGroupManagerPolicyValidator(),
				"ibm_is_instance_group_manager_action":               vpc.ResourceIBMISInstanceGroupManagerActionValidator(),
				"ibm_is_floating_ip":                                 vpc.ResourceIBMISFloatingIPValidator(),
				"ibm_is_ike_policy":                                  vpc.ResourceIBMISIKEValidator(),
				"ibm_is_image":                                       vpc.ResourceIBMISImageValidator(),
				"ibm_is_image_export_job":                            vpc.ResourceIBMIsImageExportValidator(),
				"ibm_is_instance_template":                           vpc.ResourceIBMISInstanceTemplateValidator(),
				"ibm_is_instance":                                    vpc.ResourceIBMISInstanceValidator(),
				"ibm_is_instance_action":                             vpc.ResourceIBMISInstanceActionValidator(),
				"ibm_is_instance_network_attachment":                 vpc.ResourceIBMIsInstanceNetworkAttachmentValidator(),
				"ibm_is_instance_network_interface":                  vpc.ResourceIBMIsInstanceNetworkInterfaceValidator(),
				"ibm_is_instance_disk_management":                    vpc.ResourceIBMISInstanceDiskManagementValidator(),
				"ibm_is_instance_volume_attachment":                  vpc.ResourceIBMISInstanceVolumeAttachmentValidator(),
				"ibm_is_ipsec_policy":                                vpc.ResourceIBMISIPSECValidator(),
				"ibm_is_lb_listener_policy_rule":                     vpc.ResourceIBMISLBListenerPolicyRuleValidator(),
				"ibm_is_lb_listener_policy":                          vpc.ResourceIBMISLBListenerPolicyValidator(),
				"ibm_is_lb_listener":                                 vpc.ResourceIBMISLBListenerValidator(),
				"ibm_is_lb_pool_member":                              vpc.ResourceIBMISLBPoolMemberValidator(),
				"ibm_is_lb_pool":                                     vpc.ResourceIBMISLBPoolValidator(),
				"ibm_is_lb":                                          vpc.ResourceIBMISLBValidator(),
				"ibm_is_network_acl":                                 vpc.ResourceIBMISNetworkACLValidator(),
				"ibm_is_network_acl_rule":                            vpc.ResourceIBMISNetworkACLRuleValidator(),
				"ibm_is_public_gateway":                              vpc.ResourceIBMISPublicGatewayValidator(),
				"ibm_is_private_path_service_gateway":                vpc.ResourceIBMIsPrivatePathServiceGatewayValidator(),
				"ibm_is_private_path_service_gateway_account_policy": vpc.ResourceIBMIsPrivatePathServiceGatewayAccountPolicyValidator(),
				"ibm_is_placement_group":                             vpc.ResourceIbmIsPlacementGroupValidator(),
				"ibm_is_security_group_target":                       vpc.ResourceIBMISSecurityGroupTargetValidator(),
				"ibm_is_security_group_rule":                         vpc.ResourceIBMISSecurityGroupRuleValidator(),
				"ibm_is_security_group":                              vpc.ResourceIBMISSecurityGroupValidator(),
				"ibm_is_share":                                       vpc.ResourceIbmIsShareValidator(),
				"ibm_is_share_replica_operations":                    vpc.ResourceIbmIsShareReplicaOperationsValidator(),
				"ibm_is_share_mount_target":                          vpc.ResourceIBMIsShareMountTargetValidator(),
				"ibm_is_share_snapshot":                              vpc.ResourceIBMIsShareSnapshotValidator(),
				"ibm_is_snapshot":                                    vpc.ResourceIBMISSnapshotValidator(),
				"ibm_is_snapshot_consistency_group":                  vpc.ResourceIBMIsSnapshotConsistencyGroupValidator(),
				"ibm_is_ssh_key":                                     vpc.ResourceIBMISSHKeyValidator(),
				"ibm_is_subnet":                                      vpc.ResourceIBMISSubnetValidator(),
				"ibm_is_subnet_reserved_ip":                          vpc.ResourceIBMISSubnetReservedIPValidator(),
				"ibm_is_volume":                                      vpc.ResourceIBMISVolumeValidator(),
				"ibm_is_virtual_network_interface":                   vpc.ResourceIBMIsVirtualNetworkInterfaceValidator(),
				"ibm_is_address_prefix":                              vpc.ResourceIBMISAddressPrefixValidator(),
				"ibm_is_vpc":                                         vpc.ResourceIBMISVPCValidator(),
				"ibm_is_vpc_routing_table":                           vpc.ResourceIBMISVPCRoutingTableValidator(),
				"ibm_is_vpc_routing_table_route":                     vpc.ResourceIBMISVPCRoutingTableRouteValidator(),
				"ibm_is_vpn_gateway_connection":                      vpc.ResourceIBMISVPNGatewayConnectionValidator(),
				"ibm_is_vpn_gateway":                                 vpc.ResourceIBMISVPNGatewayValidator(),
				"ibm_is_vpn_server":                                  vpc.ResourceIBMIsVPNServerValidator(),
				"ibm_is_vpn_server_route":                            vpc.ResourceIBMIsVPNServerRouteValidator(),
				"ibm_is_reservation":                                 vpc.ResourceIBMISReservationValidator(),
				"ibm_kms_key_rings":                                  kms.ResourceIBMKeyRingValidator(),
				"ibm_dns_glb_monitor":                                dnsservices.ResourceIBMPrivateDNSGLBMonitorValidator(),
				"ibm_dns_custom_resolver":                            dnsservices.ResourceIBMPrivateDNSCustomResolverValidator(),
				"ibm_dns_custom_resolver_forwarding_rule":            dnsservices.ResourceIBMPrivateDNSForwardingRuleValidator(),
				"ibm_schematics_action":                              schematics.ResourceIBMSchematicsActionValidator(),
				"ibm_schematics_job":                                 schematics.ResourceIBMSchematicsJobValidator(),
				"ibm_schematics_workspace":                           schematics.ResourceIBMSchematicsWorkspaceValidator(),
				"ibm_schematics_inventory":                           schematics.ResourceIBMSchematicsInventoryValidator(),
				"ibm_schematics_resource_query":                      schematics.ResourceIBMSchematicsResourceQueryValidator(),
				"ibm_schematics_policy":                              schematics.ResourceIbmSchematicsPolicyValidator(),
				"ibm_resource_instance":                              resourcecontroller.ResourceIBMResourceInstanceValidator(),
				"ibm_resource_key":                                   resourcecontroller.ResourceIBMResourceKeyValidator(),
				"ibm_is_virtual_endpoint_gateway":                    vpc.ResourceIBMISEndpointGatewayValidator(),
				"ibm_resource_tag":                                   globaltagging.ResourceIBMResourceTagValidator(),
				"ibm_iam_access_tag":                                 globaltagging.ResourceIBMIamAccessTagValidator(),
				"ibm_satellite_location":                             satellite.ResourceIBMSatelliteLocationValidator(),
				"ibm_satellite_cluster":                              satellite.ResourceIBMSatelliteClusterValidator(),
				"ibm_pi_volume":                                      power.ResourceIBMPIVolumeValidator(),
				"ibm_atracker_target":                                atracker.ResourceIBMAtrackerTargetValidator(),
				"ibm_atracker_route":                                 atracker.ResourceIBMAtrackerRouteValidator(),
				"ibm_atracker_settings":                              atracker.ResourceIBMAtrackerSettingsValidator(),
				"ibm_metrics_router_target":                          metricsrouter.ResourceIBMMetricsRouterTargetValidator(),
				"ibm_metrics_router_route":                           metricsrouter.ResourceIBMMetricsRouterRouteValidator(),
				"ibm_metrics_router_settings":                        metricsrouter.ResourceIBMMetricsRouterSettingsValidator(),
				"ibm_satellite_endpoint":                             satellite.ResourceIBMSatelliteEndpointValidator(),
				"ibm_satellite_host":                                 satellite.ResourceIBMSatelliteHostValidator(),

				// Partner Center Sell
				"ibm_onboarding_registration":       partnercentersell.ResourceIbmOnboardingRegistrationValidator(),
				"ibm_onboarding_product":            partnercentersell.ResourceIbmOnboardingProductValidator(),
				"ibm_onboarding_iam_registration":   partnercentersell.ResourceIbmOnboardingIamRegistrationValidator(),
				"ibm_onboarding_catalog_product":    partnercentersell.ResourceIbmOnboardingCatalogProductValidator(),
				"ibm_onboarding_catalog_plan":       partnercentersell.ResourceIbmOnboardingCatalogPlanValidator(),
				"ibm_onboarding_catalog_deployment": partnercentersell.ResourceIbmOnboardingCatalogDeploymentValidator(),
				"ibm_onboarding_resource_broker":    partnercentersell.ResourceIbmOnboardingResourceBrokerValidator(),

				// Added for Context Based Restrictions
				"ibm_cbr_zone":           contextbasedrestrictions.ResourceIBMCbrZoneValidator(),
				"ibm_cbr_zone_addresses": contextbasedrestrictions.ResourceIBMCbrZoneAddressesValidator(),
				"ibm_cbr_rule":           contextbasedrestrictions.ResourceIBMCbrRuleValidator(),

				// Added for SCC
				"ibm_scc_instance_settings":      scc.ResourceIbmSccInstanceSettingsValidator(),
				"ibm_scc_rule":                   scc.ResourceIbmSccRuleValidator(),
				"ibm_scc_control_library":        scc.ResourceIbmSccControlLibraryValidator(),
				"ibm_scc_profile":                scc.ResourceIbmSccProfileValidator(),
				"ibm_scc_profile_attachment":     scc.ResourceIbmSccProfileAttachmentValidator(),
				"ibm_scc_provider_type_instance": scc.ResourceIbmSccProviderTypeInstanceValidator(),
				"ibm_scc_scope":                  scc.ResourceIbmSccScopeValidator(),

				// Added for Toolchains
				"ibm_cd_toolchain":                         cdtoolchain.ResourceIBMCdToolchainValidator(),
				"ibm_cd_toolchain_tool_keyprotect":         cdtoolchain.ResourceIBMCdToolchainToolKeyprotectValidator(),
				"ibm_cd_toolchain_tool_secretsmanager":     cdtoolchain.ResourceIBMCdToolchainToolSecretsmanagerValidator(),
				"ibm_cd_toolchain_tool_bitbucketgit":       cdtoolchain.ResourceIBMCdToolchainToolBitbucketgitValidator(),
				"ibm_cd_toolchain_tool_githubconsolidated": cdtoolchain.ResourceIBMCdToolchainToolGithubconsolidatedValidator(),
				"ibm_cd_toolchain_tool_gitlab":             cdtoolchain.ResourceIBMCdToolchainToolGitlabValidator(),
				"ibm_cd_toolchain_tool_hostedgit":          cdtoolchain.ResourceIBMCdToolchainToolHostedgitValidator(),
				"ibm_cd_toolchain_tool_artifactory":        cdtoolchain.ResourceIBMCdToolchainToolArtifactoryValidator(),
				"ibm_cd_toolchain_tool_custom":             cdtoolchain.ResourceIBMCdToolchainToolCustomValidator(),
				"ibm_cd_toolchain_tool_pipeline":           cdtoolchain.ResourceIBMCdToolchainToolPipelineValidator(),
				"ibm_cd_toolchain_tool_slack":              cdtoolchain.ResourceIBMCdToolchainToolSlackValidator(),
				"ibm_cd_toolchain_tool_devopsinsights":     cdtoolchain.ResourceIBMCdToolchainToolDevopsinsightsValidator(),
				"ibm_cd_toolchain_tool_sonarqube":          cdtoolchain.ResourceIBMCdToolchainToolSonarqubeValidator(),
				"ibm_cd_toolchain_tool_hashicorpvault":     cdtoolchain.ResourceIBMCdToolchainToolHashicorpvaultValidator(),
				"ibm_cd_toolchain_tool_securitycompliance": cdtoolchain.ResourceIBMCdToolchainToolSecuritycomplianceValidator(),
				"ibm_cd_toolchain_tool_privateworker":      cdtoolchain.ResourceIBMCdToolchainToolPrivateworkerValidator(),
				"ibm_cd_toolchain_tool_appconfig":          cdtoolchain.ResourceIBMCdToolchainToolAppconfigValidator(),
				"ibm_cd_toolchain_tool_jenkins":            cdtoolchain.ResourceIBMCdToolchainToolJenkinsValidator(),
				"ibm_cd_toolchain_tool_nexus":              cdtoolchain.ResourceIBMCdToolchainToolNexusValidator(),
				"ibm_cd_toolchain_tool_pagerduty":          cdtoolchain.ResourceIBMCdToolchainToolPagerdutyValidator(),
				"ibm_cd_toolchain_tool_saucelabs":          cdtoolchain.ResourceIBMCdToolchainToolSaucelabsValidator(),
				"ibm_cd_toolchain_tool_jira":               cdtoolchain.ResourceIBMCdToolchainToolJiraValidator(),
				"ibm_cd_toolchain_tool_eventnotifications": cdtoolchain.ResourceIBMCdToolchainToolEventnotificationsValidator(),

				// // Added for Tekton Pipeline
				"ibm_cd_tekton_pipeline_definition":       cdtektonpipeline.ResourceIBMCdTektonPipelineDefinitionValidator(),
				"ibm_cd_tekton_pipeline_trigger_property": cdtektonpipeline.ResourceIBMCdTektonPipelineTriggerPropertyValidator(),
				"ibm_cd_tekton_pipeline_property":         cdtektonpipeline.ResourceIBMCdTektonPipelinePropertyValidator(),
				"ibm_cd_tekton_pipeline_trigger":          cdtektonpipeline.ResourceIBMCdTektonPipelineTriggerValidator(),

				"ibm_container_addons":                      kubernetes.ResourceIBMContainerAddOnsValidator(),
				"ibm_container_alb_create":                  kubernetes.ResourceIBMContainerAlbCreateValidator(),
				"ibm_container_nlb_dns":                     kubernetes.ResourceIBMContainerNlbDnsValidator(),
				"ibm_container_vpc_alb_create":              kubernetes.ResourceIBMContainerVpcAlbCreateNewValidator(),
				"ibm_container_storage_attachment":          kubernetes.ResourceIBMContainerVpcWorkerVolumeAttachmentValidator(),
				"ibm_container_worker_pool_zone_attachment": kubernetes.ResourceIBMContainerWorkerPoolZoneAttachmentValidator(),
				"ibm_container_bind_service":                kubernetes.ResourceIBMContainerBindServiceValidator(),
				"ibm_container_alb_cert":                    kubernetes.ResourceIBMContainerALBCertValidator(),
				"ibm_container_ingress_instance":            kubernetes.ResourceIBMContainerIngressInstanceValidator(),
				"ibm_container_ingress_secret_tls":          kubernetes.ResourceIBMContainerIngressSecretTLSValidator(),
				"ibm_container_ingress_secret_opaque":       kubernetes.ResourceIBMContainerIngressSecretOpaqueValidator(),
				"ibm_container_cluster_feature":             kubernetes.ResourceIBMContainerClusterFeatureValidator(),

				"ibm_iam_access_group_dynamic_rule":        iamaccessgroup.ResourceIBMIAMDynamicRuleValidator(),
				"ibm_iam_access_group_members":             iamaccessgroup.ResourceIBMIAMAccessGroupMembersValidator(),
				"ibm_iam_access_group_template":            iamaccessgroup.ResourceIBMIAMAccessGroupTemplateValidator(),
				"ibm_iam_access_group_template_version":    iamaccessgroup.ResourceIBMIAMAccessGroupTemplateVersionValidator(),
				"ibm_iam_access_group_template_assignment": iamaccessgroup.ResourceIBMIAMAccessGroupTemplateAssignmentValidator(),
				"ibm_iam_trusted_profile_claim_rule":       iamidentity.ResourceIBMIAMTrustedProfileClaimRuleValidator(),
				"ibm_iam_trusted_profile_link":             iamidentity.ResourceIBMIAMTrustedProfileLinkValidator(),
				"ibm_iam_service_api_key":                  iamidentity.ResourceIBMIAMServiceAPIKeyValidator(),
				"ibm_iam_trusted_profile_identity":         iamidentity.ResourceIBMIamTrustedProfileIdentityValidator(),

				"ibm_iam_trusted_profile_policy":  iampolicy.ResourceIBMIAMTrustedProfilePolicyValidator(),
				"ibm_iam_access_group_policy":     iampolicy.ResourceIBMIAMAccessGroupPolicyValidator(),
				"ibm_iam_service_policy":          iampolicy.ResourceIBMIAMServicePolicyValidator(),
				"ibm_iam_authorization_policy":    iampolicy.ResourceIBMIAMAuthorizationPolicyValidator(),
				"ibm_iam_policy_template":         iampolicy.ResourceIBMIAMPolicyTemplateValidator(),
				"ibm_iam_policy_template_version": iampolicy.ResourceIBMIAMPolicyTemplateVersionValidator(),

				// // Added for Usage Reports
				"ibm_billing_report_snapshot": usagereports.ResourceIBMBillingReportSnapshotValidator(),

				// // Added for Secrets Manager
				"ibm_sm_secret_group":                                                secretsmanager.ResourceIbmSmSecretGroupValidator(),
				"ibm_sm_en_registration":                                             secretsmanager.ResourceIbmSmEnRegistrationValidator(),
				"ibm_sm_public_certificate_configuration_dns_cis":                    secretsmanager.ResourceIbmSmConfigurationPublicCertificateDNSCisValidator(),
				"ibm_sm_public_certificate_configuration_dns_classic_infrastructure": secretsmanager.ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureValidator(),

				// // Added for Code Engine
				"ibm_code_engine_allowed_outbound_destination": codeengine.ResourceIbmCodeEngineAllowedOutboundDestinationValidator(),
				"ibm_code_engine_app":                          codeengine.ResourceIbmCodeEngineAppValidator(),
				"ibm_code_engine_binding":                      codeengine.ResourceIbmCodeEngineBindingValidator(),
				"ibm_code_engine_build":                        codeengine.ResourceIbmCodeEngineBuildValidator(),
				"ibm_code_engine_config_map":                   codeengine.ResourceIbmCodeEngineConfigMapValidator(),
				"ibm_code_engine_domain_mapping":               codeengine.ResourceIbmCodeEngineDomainMappingValidator(),
				"ibm_code_engine_function":                     codeengine.ResourceIbmCodeEngineFunctionValidator(),
				"ibm_code_engine_job":                          codeengine.ResourceIbmCodeEngineJobValidator(),
				"ibm_code_engine_project":                      codeengine.ResourceIbmCodeEngineProjectValidator(),
				"ibm_code_engine_secret":                       codeengine.ResourceIbmCodeEngineSecretValidator(),

				// Added for Project
				"ibm_project":             project.ResourceIbmProjectValidator(),
				"ibm_project_config":      project.ResourceIbmProjectConfigValidator(),
				"ibm_project_environment": project.ResourceIbmProjectEnvironmentValidator(),

				// Added for Event Notifications

				"ibm_en_smtp_configuration":       eventnotification.ResourceIBMEnSMTPConfigurationValidator(),
				"ibm_en_smtp_user":                eventnotification.ResourceIBMEnSMTPUserValidator(),
				"ibm_en_destination_custom_email": eventnotification.ResourceIBMEnEmailDestinationValidator(),

				// Added for VMware as a Service
				"ibm_vmaas_vdc":             vmware.ResourceIbmVmaasVdcValidator(),
				"ibm_logs_alert":            logs.ResourceIbmLogsAlertValidator(),
				"ibm_logs_rule_group":       logs.ResourceIbmLogsRuleGroupValidator(),
				"ibm_logs_outgoing_webhook": logs.ResourceIbmLogsOutgoingWebhookValidator(),
				"ibm_logs_policy":           logs.ResourceIbmLogsPolicyValidator(),
				"ibm_logs_dashboard":        logs.ResourceIbmLogsDashboardValidator(),
				"ibm_logs_e2m":              logs.ResourceIbmLogsE2mValidator(),
				"ibm_logs_view":             logs.ResourceIbmLogsViewValidator(),
				"ibm_logs_view_folder":      logs.ResourceIbmLogsViewFolderValidator(),
				"ibm_logs_dashboard_folder": logs.ResourceIbmLogsDashboardFolderValidator(),
				"ibm_logs_enrichment":       logs.ResourceIbmLogsEnrichmentValidator(),
				"ibm_logs_data_access_rule": logs.ResourceIbmLogsDataAccessRuleValidator(),
				"ibm_logs_stream":           logs.ResourceIbmLogsStreamValidator(),

				// Added for Logs Router Service
				"ibm_logs_router_tenant": logsrouting.ResourceIBMLogsRouterTenantValidator(),

				// Added for Software Defined Storage as a Service
				"ibm_sds_volume": sdsaas.ResourceIBMSdsVolumeValidator(),
				"ibm_sds_host":   sdsaas.ResourceIBMSdsHostValidator(),
			},
			DataSourceValidatorDictionary: map[string]*validate.ResourceValidator{
				"ibm_is_subnet":                     vpc.DataSourceIBMISSubnetValidator(),
				"ibm_is_snapshot_consistency_group": vpc.DataSourceIBMISSnapshotConsistencyGroupValidator(),
				"ibm_is_snapshot":                   vpc.DataSourceIBMISSnapshotValidator(),
				"ibm_is_images":                     vpc.DataSourceIBMISImagesValidator(),
				"ibm_dl_offering_speeds":            directlink.DataSourceIBMDLOfferingSpeedsValidator(),
				"ibm_dl_routers":                    directlink.DataSourceIBMDLRoutersValidator(),
				"ibm_resource_instance":             resourcecontroller.DataSourceIBMResourceInstanceValidator(),
				"ibm_resource_key":                  resourcecontroller.DataSourceIBMResourceKeyValidator(),
				"ibm_resource_group":                resourcemanager.DataSourceIBMResourceGroupValidator(),

				// bare_metal_server
				"ibm_is_bare_metal_server": vpc.DataSourceIBMIsBareMetalServerValidator(),

				"ibm_is_vpc":                          vpc.DataSourceIBMISVpcValidator(),
				"ibm_is_volume":                       vpc.DataSourceIBMISVolumeValidator(),
				"ibm_cis_webhooks":                    cis.DataSourceIBMCISAlertWebhooksValidator(),
				"ibm_cis_alerts":                      cis.DataSourceIBMCISAlertsValidator(),
				"ibm_cis_bot_managements":             cis.DataSourceIBMCISBotManagementValidator(),
				"ibm_cis_bot_analytics":               cis.DataSourceIBMCISBotAnalyticsValidator(),
				"ibm_cis_cache_settings":              cis.DataSourceIBMCISCacheSettingsValidator(),
				"ibm_cis_custom_certificates":         cis.DataSourceIBMCISCustomCertificatesValidator(),
				"ibm_cis_custom_pages":                cis.DataSourceIBMCISCustomPagesValidator(),
				"ibm_cis_dns_records":                 cis.DataSourceIBMCISDNSRecordsValidator(),
				"ibm_cis_domain":                      cis.DataSourceIBMCISDomainValidator(),
				"ibm_cis_certificates":                cis.DataSourceIBMCISCertificatesValidator(),
				"ibm_cis_edge_functions_actions":      cis.DataSourceIBMCISEdgeFunctionsActionsValidator(),
				"ibm_cis_edge_functions_triggers":     cis.DataSourceIBMCISEdgeFunctionsTriggersValidator(),
				"ibm_cis_filters":                     cis.DataSourceIBMCISFiltersValidator(),
				"ibm_cis_firewall_rules":              cis.DataSourceIBMCISFirewallRulesValidator(),
				"ibm_cis_firewall":                    cis.DataSourceIBMCISFirewallsRecordValidator(),
				"ibm_cis_global_load_balancers":       cis.DataSourceIBMCISGlbsValidator(),
				"ibm_cis_healthchecks":                cis.DataSourceIBMCISHealthChecksValidator(),
				"ibm_cis_mtls_apps":                   cis.DataSourceIBMCISMtlsAppValidator(),
				"ibm_cis_mtlss":                       cis.DataSourceIBMCISMtlsValidator(),
				"ibm_cis_origin_auths":                cis.DataSourceIBMCISOriginAuthPullValidator(),
				"ibm_cis_origin_pools":                cis.DataSourceIBMCISOriginPoolsValidator(),
				"ibm_cis_page_rules":                  cis.DataSourceIBMCISPageRulesValidator(),
				"ibm_cis_range_apps":                  cis.DataSourceIBMCISRangeAppsValidator(),
				"ibm_cis_rate_limit":                  cis.DataSourceIBMCISRateLimitValidator(),
				"ibm_cis_rulesets":                    cis.DataSourceIBMCISRulesetsValidator(),
				"ibm_cis_ruleset_versions":            cis.DataSourceIBMCISRulesetVersionsValidator(),
				"ibm_cis_ruleset_rules_by_tag":        cis.DataSourceIBMCISRulesetRulesByTagValidator(),
				"ibm_cis_ruleset_entrypoint_versions": cis.DataSourceIBMCISRulesetEntrypointVersionsValidator(),
				"ibm_cis_waf_groups":                  cis.DataSourceIBMCISWAFGroupsValidator(),
				"ibm_cis_waf_packages":                cis.DataSourceIBMCISWAFPackagesValidator(),
				"ibm_cis_waf_rules":                   cis.DataSourceIBMCISWAFRulesValidator(),
				"ibm_cis_logpush_jobs":                cis.DataSourceIBMCISLogPushJobsValidator(),
				"ibm_cis_origin_certificates":         cis.DataIBMCISOriginCertificateOrderValidator(),

				"ibm_config_aggregator_configurations": configurationaggregator.DataSourceIbmConfigAggregatorValidator(),
				"ibm_cos_bucket":                       cos.DataSourceIBMCosBucketValidator(),

				"ibm_database_backups":                database.DataSourceIBMDatabaseBackupsValidator(),
				"ibm_database_connection":             database.DataSourceIBMDatabaseConnectionValidator(),
				"ibm_database_point_in_time_recovery": database.DataSourceIBMDatabasePointInTimeRecoveryValidator(),
				"ibm_database_remotes":                database.DataSourceIBMDatabaseRemotesValidator(),
				"ibm_database_tasks":                  database.DataSourceIBMDatabaseTasksValidator(),
				"ibm_database":                        database.DataSourceIBMDatabaseInstanceValidator(),

				"ibm_container_addons":                  kubernetes.DataSourceIBMContainerAddOnsValidator(),
				"ibm_container_nlb_dns":                 kubernetes.DataSourceIBMContainerNLBDNSValidator(),
				"ibm_container_storage_attachment":      kubernetes.DataSourceIBMContainerVpcWorkerVolumeAttachmentValidator(),
				"ibm_container_vpc_cluster_worker_pool": kubernetes.DataSourceIBMContainerVpcClusterWorkerPoolValidator(),
				"ibm_container_worker_pool":             kubernetes.DataSourceIBMContainerWorkerPoolValidator(),
				"ibm_container_bind_service":            kubernetes.DataSourceIBMContainerBindServiceValidator(),
				"ibm_container_cluster_config":          kubernetes.DataSourceIBMContainerClusterConfigValidator(),
				"ibm_container_cluster":                 kubernetes.DataSourceIBMContainerClusterValidator(),
				"ibm_container_vpc_cluster_worker":      kubernetes.DataSourceIBMContainerVPCClusterWorkerValidator(),
				"ibm_container_vpc_cluster":             kubernetes.DataSourceIBMContainerVPCClusterValidator(),
				"ibm_container_alb_cert":                kubernetes.DataSourceIBMContainerALBCertValidator(),
				"ibm_container_ingress_instance":        kubernetes.DataSourceIBMContainerIngressInstanceValidator(),
				"ibm_container_ingress_secret_tls":      kubernetes.DataSourceIBMContainerIngressSecretTLSValidator(),
				"ibm_container_ingress_secret_opaque":   kubernetes.DataSourceIBMContainerIngressSecretOpaqueValidator(),

				"ibm_iam_access_group": iamaccessgroup.DataSourceIBMIAMAccessGroupValidator(),

				"ibm_iam_service_id":                  iamidentity.DataSourceIBMIAMServiceIDValidator(),
				"ibm_iam_trusted_profile_claim_rule":  iamidentity.DataSourceIBMIamTrustedProfileClaimRuleValidator(),
				"ibm_iam_trusted_profile_link":        iamidentity.DataSourceIBMIamTrustedProfileLinkValidator(),
				"ibm_iam_trusted_profile_links":       iamidentity.DataSourceIBMIamTrustedProfileLinksValidator(),
				"ibm_iam_trusted_profile":             iamidentity.DataSourceIBMIamTrustedProfileValidator(),
				"ibm_iam_trusted_profile_claim_rules": iamidentity.DataSourceIBMIamTrustedProfileClaimRulesValidator(),

				"ibm_iam_access_group_policy":    iampolicy.DataSourceIBMIAMAccessGroupPolicyValidator(),
				"ibm_iam_service_policy":         iampolicy.DataSourceIBMIAMServicePolicyValidator(),
				"ibm_iam_trusted_profile_policy": iampolicy.DataSourceIBMIAMTrustedProfilePolicyValidator(),
			},
		}
	})
	return globalValidatorDict
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var bluemixAPIKey string
	var bluemixTimeout int
	var iamToken, iamRefreshToken, iamTrustedProfileId string
	if key, ok := d.GetOk("bluemix_api_key"); ok {
		bluemixAPIKey = key.(string)
	}
	if key, ok := d.GetOk("ibmcloud_api_key"); ok {
		bluemixAPIKey = key.(string)
	}
	if itoken, ok := d.GetOk("iam_token"); ok {
		iamToken = itoken.(string)
	}
	if rtoken, ok := d.GetOk("iam_refresh_token"); ok {
		iamRefreshToken = rtoken.(string)
	}
	if ttoken, ok := d.GetOk("iam_profile_id"); ok {
		iamTrustedProfileId = ttoken.(string)
	}
	var softlayerUsername, softlayerAPIKey, softlayerEndpointUrl string
	var softlayerTimeout int
	if username, ok := d.GetOk("softlayer_username"); ok {
		softlayerUsername = username.(string)
	}
	if username, ok := d.GetOk("iaas_classic_username"); ok {
		softlayerUsername = username.(string)
	}
	if apikey, ok := d.GetOk("softlayer_api_key"); ok {
		softlayerAPIKey = apikey.(string)
	}
	if apikey, ok := d.GetOk("iaas_classic_api_key"); ok {
		softlayerAPIKey = apikey.(string)
	}
	if endpoint, ok := d.GetOk("softlayer_endpoint_url"); ok {
		softlayerEndpointUrl = endpoint.(string)
	}
	if endpoint, ok := d.GetOk("iaas_classic_endpoint_url"); ok {
		softlayerEndpointUrl = endpoint.(string)
	}
	if tm, ok := d.GetOk("softlayer_timeout"); ok {
		softlayerTimeout = tm.(int)
	}
	if tm, ok := d.GetOk("iaas_classic_timeout"); ok {
		softlayerTimeout = tm.(int)
	}

	if tm, ok := d.GetOk("bluemix_timeout"); ok {
		bluemixTimeout = tm.(int)
	}
	if tm, ok := d.GetOk("ibmcloud_timeout"); ok {
		bluemixTimeout = tm.(int)
	}
	var visibility string
	if v, ok := d.GetOk("visibility"); ok {
		visibility = v.(string)
	}
	var file string
	if f, ok := d.GetOk("endpoints_file_path"); ok {
		file = f.(string)
	}

	resourceGrp := d.Get("resource_group").(string)
	region := d.Get("region").(string)
	zone := d.Get("zone").(string)
	retryCount := d.Get("max_retries").(int)
	wskNameSpace := d.Get("function_namespace").(string)
	riaasEndPoint := d.Get("riaas_endpoint").(string)

	wskEnvVal, err := schema.EnvDefaultFunc("FUNCTION_NAMESPACE", "")()
	if err != nil {
		return nil, err
	}
	// Set environment variable to be used in DiffSupressFunction
	if wskEnvVal.(string) == "" {
		os.Setenv("FUNCTION_NAMESPACE", wskNameSpace)
	}

	config := conns.Config{
		BluemixAPIKey:        bluemixAPIKey,
		Region:               region,
		ResourceGroup:        resourceGrp,
		BluemixTimeout:       time.Duration(bluemixTimeout) * time.Second,
		SoftLayerTimeout:     time.Duration(softlayerTimeout) * time.Second,
		SoftLayerUserName:    softlayerUsername,
		SoftLayerAPIKey:      softlayerAPIKey,
		RetryCount:           retryCount,
		SoftLayerEndpointURL: softlayerEndpointUrl,
		RetryDelay:           conns.RetryAPIDelay,
		FunctionNameSpace:    wskNameSpace,
		RiaasEndPoint:        riaasEndPoint,
		IAMToken:             iamToken,
		IAMRefreshToken:      iamRefreshToken,
		Zone:                 zone,
		Visibility:           visibility,
		EndpointsFile:        file,
		IAMTrustedProfileID:  iamTrustedProfileId,
	}

	return config.ClientSession()
}
