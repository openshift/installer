// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"os"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/mutexkv"
)

// This is a global MutexKV for use within this plugin.
var ibmMutexKV = mutexkv.NewMutexKV()

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
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
				//DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_GENERATION", "IBMCLOUD_GENERATION"}, nil),
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
				ValidateFunc: validateAllowedStringValue([]string{"public", "private", "public-and-private"}),
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
			"ibm_api_gateway":        dataSourceIBMApiGateway(),
			"ibm_account":            dataSourceIBMAccount(),
			"ibm_app":                dataSourceIBMApp(),
			"ibm_app_domain_private": dataSourceIBMAppDomainPrivate(),
			"ibm_app_domain_shared":  dataSourceIBMAppDomainShared(),
			"ibm_app_route":          dataSourceIBMAppRoute(),

			// AppID
			"ibm_appid_action_url":               dataSourceIBMAppIDActionURL(),
			"ibm_appid_apm":                      dataSourceIBMAppIDAPM(),
			"ibm_appid_application":              dataSourceIBMAppIDApplication(),
			"ibm_appid_application_scopes":       dataSourceIBMAppIDApplicationScopes(),
			"ibm_appid_application_roles":        dataSourceIBMAppIDApplicationRoles(),
			"ibm_appid_applications":             dataSourceIBMAppIDApplications(),
			"ibm_appid_audit_status":             dataSourceIBMAppIDAuditStatus(),
			"ibm_appid_cloud_directory_template": dataSourceIBMAppIDCloudDirectoryTemplate(),
			"ibm_appid_cloud_directory_user":     dataSourceIBMAppIDCloudDirectoryUser(),
			"ibm_appid_idp_cloud_directory":      dataSourceIBMAppIDIDPCloudDirectory(),
			"ibm_appid_idp_custom":               dataSourceIBMAppIDIDPCustom(),
			"ibm_appid_idp_facebook":             dataSourceIBMAppIDIDPFacebook(),
			"ibm_appid_idp_google":               dataSourceIBMAppIDIDPGoogle(),
			"ibm_appid_idp_saml":                 dataSourceIBMAppIDIDPSAML(),
			"ibm_appid_idp_saml_metadata":        dataSourceIBMAppIDIDPSAMLMetadata(),
			"ibm_appid_languages":                dataSourceIBMAppIDLanguages(),
			"ibm_appid_mfa":                      dataSourceIBMAppIDMFA(),
			"ibm_appid_mfa_channel":              dataSourceIBMAppIDMFAChannel(),
			"ibm_appid_password_regex":           dataSourceIBMAppIDPasswordRegex(),
			"ibm_appid_token_config":             dataSourceIBMAppIDTokenConfig(),
			"ibm_appid_redirect_urls":            dataSourceIBMAppIDRedirectURLs(),
			"ibm_appid_role":                     dataSourceIBMAppIDRole(),
			"ibm_appid_roles":                    dataSourceIBMAppIDRoles(),
			"ibm_appid_theme_color":              dataSourceIBMAppIDThemeColor(),
			"ibm_appid_theme_text":               dataSourceIBMAppIDThemeText(),
			"ibm_appid_user_roles":               dataSourceIBMAppIDUserRoles(),

			"ibm_function_action":                    dataSourceIBMFunctionAction(),
			"ibm_function_package":                   dataSourceIBMFunctionPackage(),
			"ibm_function_rule":                      dataSourceIBMFunctionRule(),
			"ibm_function_trigger":                   dataSourceIBMFunctionTrigger(),
			"ibm_function_namespace":                 dataSourceIBMFunctionNamespace(),
			"ibm_certificate_manager_certificates":   dataIBMCertificateManagerCertificates(),
			"ibm_certificate_manager_certificate":    dataIBMCertificateManagerCertificate(),
			"ibm_cis":                                dataSourceIBMCISInstance(),
			"ibm_cis_dns_records":                    dataSourceIBMCISDNSRecords(),
			"ibm_cis_certificates":                   dataIBMCISCertificates(),
			"ibm_cis_global_load_balancers":          dataSourceIBMCISGlbs(),
			"ibm_cis_origin_pools":                   dataSourceIBMCISOriginPools(),
			"ibm_cis_healthchecks":                   dataSourceIBMCISHealthChecks(),
			"ibm_cis_domain":                         dataSourceIBMCISDomain(),
			"ibm_cis_firewall":                       dataIBMCISFirewallsRecord(),
			"ibm_cis_cache_settings":                 dataSourceIBMCISCacheSetting(),
			"ibm_cis_waf_packages":                   dataSourceIBMCISWAFPackages(),
			"ibm_cis_range_apps":                     dataSourceIBMCISRangeApps(),
			"ibm_cis_custom_certificates":            dataSourceIBMCISCustomCertificates(),
			"ibm_cis_rate_limit":                     dataSourceIBMCISRateLimit(),
			"ibm_cis_ip_addresses":                   dataSourceIBMCISIP(),
			"ibm_cis_waf_groups":                     dataSourceIBMCISWAFGroups(),
			"ibm_cis_edge_functions_actions":         dataSourceIBMCISEdgeFunctionsActions(),
			"ibm_cis_edge_functions_triggers":        dataSourceIBMCISEdgeFunctionsTriggers(),
			"ibm_cis_custom_pages":                   dataSourceIBMCISCustomPages(),
			"ibm_cis_page_rules":                     dataSourceIBMCISPageRules(),
			"ibm_cis_waf_rules":                      dataSourceIBMCISWAFRules(),
			"ibm_cis_filters":                        dataSourceIBMCISFilters(),
			"ibm_cis_firewall_rules":                 dataSourceIBMCISFirewallRules(),
			"ibm_cloudant":                           dataSourceIBMCloudant(),
			"ibm_database":                           dataSourceIBMDatabaseInstance(),
			"ibm_compute_bare_metal":                 dataSourceIBMComputeBareMetal(),
			"ibm_compute_image_template":             dataSourceIBMComputeImageTemplate(),
			"ibm_compute_placement_group":            dataSourceIBMComputePlacementGroup(),
			"ibm_compute_reserved_capacity":          dataSourceIBMComputeReservedCapacity(),
			"ibm_compute_ssh_key":                    dataSourceIBMComputeSSHKey(),
			"ibm_compute_vm_instance":                dataSourceIBMComputeVmInstance(),
			"ibm_container_addons":                   datasourceIBMContainerAddOns(),
			"ibm_container_alb":                      dataSourceIBMContainerALB(),
			"ibm_container_alb_cert":                 dataSourceIBMContainerALBCert(),
			"ibm_container_bind_service":             dataSourceIBMContainerBindService(),
			"ibm_container_cluster":                  dataSourceIBMContainerCluster(),
			"ibm_container_cluster_config":           dataSourceIBMContainerClusterConfig(),
			"ibm_container_cluster_versions":         dataSourceIBMContainerClusterVersions(),
			"ibm_container_cluster_worker":           dataSourceIBMContainerClusterWorker(),
			"ibm_container_nlb_dns":                  dataSourceIBMContainerNLBDNS(),
			"ibm_container_vpc_cluster_alb":          dataSourceIBMContainerVPCClusterALB(),
			"ibm_container_vpc_alb":                  dataSourceIBMContainerVPCClusterALB(),
			"ibm_container_vpc_cluster":              dataSourceIBMContainerVPCCluster(),
			"ibm_container_vpc_cluster_worker":       dataSourceIBMContainerVPCClusterWorker(),
			"ibm_container_vpc_cluster_worker_pool":  dataSourceIBMContainerVpcClusterWorkerPool(),
			"ibm_container_vpc_worker_pool":          dataSourceIBMContainerVpcClusterWorkerPool(),
			"ibm_container_worker_pool":              dataSourceIBMContainerWorkerPool(),
			"ibm_container_storage_attachment":       dataSourceIBMContainerVpcWorkerVolumeAttachment(),
			"ibm_cr_namespaces":                      dataIBMContainerRegistryNamespaces(),
			"ibm_cloud_shell_account_settings":       dataSourceIBMCloudShellAccountSettings(),
			"ibm_cos_bucket":                         dataSourceIBMCosBucket(),
			"ibm_cos_bucket_object":                  dataSourceIBMCosBucketObject(),
			"ibm_dns_domain_registration":            dataSourceIBMDNSDomainRegistration(),
			"ibm_dns_domain":                         dataSourceIBMDNSDomain(),
			"ibm_dns_secondary":                      dataSourceIBMDNSSecondary(),
			"ibm_event_streams_topic":                dataSourceIBMEventStreamsTopic(),
			"ibm_event_streams_schema":               dataSourceIBMEventStreamsSchema(),
			"ibm_hpcs":                               dataSourceIBMHPCS(),
			"ibm_iam_access_group":                   dataSourceIBMIAMAccessGroup(),
			"ibm_iam_access_group_policy":            dataSourceIBMIAMAccessGroupPolicy(),
			"ibm_iam_account_settings":               dataSourceIBMIAMAccountSettings(),
			"ibm_iam_auth_token":                     dataSourceIBMIAMAuthToken(),
			"ibm_iam_role_actions":                   datasourceIBMIAMRoleAction(),
			"ibm_iam_users":                          dataSourceIBMIAMUsers(),
			"ibm_iam_roles":                          datasourceIBMIAMRole(),
			"ibm_iam_user_policy":                    dataSourceIBMIAMUserPolicy(),
			"ibm_iam_authorization_policies":         dataSourceIBMIAMAuthorizationPolicies(),
			"ibm_iam_user_profile":                   dataSourceIBMIAMUserProfile(),
			"ibm_iam_service_id":                     dataSourceIBMIAMServiceID(),
			"ibm_iam_service_policy":                 dataSourceIBMIAMServicePolicy(),
			"ibm_iam_api_key":                        dataSourceIbmIamApiKey(),
			"ibm_iam_trusted_profile":                dataSourceIBMIamTrustedProfile(),
			"ibm_iam_trusted_profile_claim_rule":     dataSourceIBMIamTrustedProfileClaimRule(),
			"ibm_iam_trusted_profile_link":           dataSourceIBMIamTrustedProfileLink(),
			"ibm_iam_trusted_profile_claim_rules":    dataSourceIBMIamTrustedProfileClaimRules(),
			"ibm_iam_trusted_profile_links":          dataSourceIBMIamTrustedProfileLinks(),
			"ibm_iam_trusted_profiles":               dataSourceIBMIamTrustedProfiles(),
			"ibm_iam_trusted_profile_policy":         dataSourceIBMIAMTrustedProfilePolicy(),
			"ibm_is_dedicated_host":                  dataSourceIbmIsDedicatedHost(),
			"ibm_is_dedicated_hosts":                 dataSourceIbmIsDedicatedHosts(),
			"ibm_is_dedicated_host_profile":          dataSourceIbmIsDedicatedHostProfile(),
			"ibm_is_dedicated_host_profiles":         dataSourceIbmIsDedicatedHostProfiles(),
			"ibm_is_dedicated_host_group":            dataSourceIbmIsDedicatedHostGroup(),
			"ibm_is_dedicated_host_groups":           dataSourceIbmIsDedicatedHostGroups(),
			"ibm_is_dedicated_host_disk":             dataSourceIbmIsDedicatedHostDisk(),
			"ibm_is_dedicated_host_disks":            dataSourceIbmIsDedicatedHostDisks(),
			"ibm_is_placement_group":                 dataSourceIbmIsPlacementGroup(),
			"ibm_is_placement_groups":                dataSourceIbmIsPlacementGroups(),
			"ibm_is_floating_ip":                     dataSourceIBMISFloatingIP(),
			"ibm_is_floating_ips":                    dataSourceIBMIsFloatingIps(),
			"ibm_is_flow_log":                        dataSourceIBMIsFlowLog(),
			"ibm_is_flow_logs":                       dataSourceIBMISFlowLogs(),
			"ibm_is_image":                           dataSourceIBMISImage(),
			"ibm_is_images":                          dataSourceIBMISImages(),
			"ibm_is_endpoint_gateway_targets":        dataSourceIBMISEndpointGatewayTargets(),
			"ibm_is_instance_group":                  dataSourceIBMISInstanceGroup(),
			"ibm_is_instance_group_memberships":      dataSourceIBMISInstanceGroupMemberships(),
			"ibm_is_instance_group_membership":       dataSourceIBMISInstanceGroupMembership(),
			"ibm_is_instance_group_manager":          dataSourceIBMISInstanceGroupManager(),
			"ibm_is_instance_group_managers":         dataSourceIBMISInstanceGroupManagers(),
			"ibm_is_instance_group_manager_policies": dataSourceIBMISInstanceGroupManagerPolicies(),
			"ibm_is_instance_group_manager_policy":   dataSourceIBMISInstanceGroupManagerPolicy(),
			"ibm_is_instance_group_manager_action":   dataSourceIBMISInstanceGroupManagerAction(),
			"ibm_is_instance_group_manager_actions":  dataSourceIBMISInstanceGroupManagerActions(),
			"ibm_is_virtual_endpoint_gateways":       dataSourceIBMISEndpointGateways(),
			"ibm_is_virtual_endpoint_gateway_ips":    dataSourceIBMISEndpointGatewayIPs(),
			"ibm_is_virtual_endpoint_gateway":        dataSourceIBMISEndpointGateway(),
			"ibm_is_instance_template":               dataSourceIBMISInstanceTemplate(),
			"ibm_is_instance_templates":              dataSourceIBMISInstanceTemplates(),
			"ibm_is_instance_profile":                dataSourceIBMISInstanceProfile(),
			"ibm_is_instance_profiles":               dataSourceIBMISInstanceProfiles(),
			"ibm_is_instance":                        dataSourceIBMISInstance(),
			"ibm_is_instances":                       dataSourceIBMISInstances(),
			"ibm_is_instance_network_interface":      dataSourceIBMIsInstanceNetworkInterface(),
			"ibm_is_instance_network_interfaces":     dataSourceIBMIsInstanceNetworkInterfaces(),
			"ibm_is_instance_disk":                   dataSourceIbmIsInstanceDisk(),
			"ibm_is_instance_disks":                  dataSourceIbmIsInstanceDisks(),
			"ibm_is_instance_volume_attachment":      dataSourceIBMISInstanceVolumeAttachment(),
			"ibm_is_instance_volume_attachments":     dataSourceIBMISInstanceVolumeAttachments(),
			"ibm_is_lb":                              dataSourceIBMISLB(),
			"ibm_is_lb_profiles":                     dataSourceIBMISLbProfiles(),
			"ibm_is_lbs":                             dataSourceIBMISLBS(),
			"ibm_is_public_gateway":                  dataSourceIBMISPublicGateway(),
			"ibm_is_public_gateways":                 dataSourceIBMISPublicGateways(),
			"ibm_is_region":                          dataSourceIBMISRegion(),
			"ibm_is_regions":                         dataSourceIBMISRegions(),
			"ibm_is_ssh_key":                         dataSourceIBMISSSHKey(),
			"ibm_is_subnet":                          dataSourceIBMISSubnet(),
			"ibm_is_subnets":                         dataSourceIBMISSubnets(),
			"ibm_is_subnet_reserved_ip":              dataSourceIBMISReservedIP(),
			"ibm_is_subnet_reserved_ips":             dataSourceIBMISReservedIPs(),
			"ibm_is_security_group":                  dataSourceIBMISSecurityGroup(),
			"ibm_is_security_group_target":           dataSourceIBMISSecurityGroupTarget(),
			"ibm_is_security_group_targets":          dataSourceIBMISSecurityGroupTargets(),
			"ibm_is_snapshot":                        dataSourceSnapshot(),
			"ibm_is_snapshots":                       dataSourceSnapshots(),
			"ibm_is_volume":                          dataSourceIBMISVolume(),
			"ibm_is_volume_profile":                  dataSourceIBMISVolumeProfile(),
			"ibm_is_volume_profiles":                 dataSourceIBMISVolumeProfiles(),
			"ibm_is_vpc":                             dataSourceIBMISVPC(),
			"ibm_is_vpcs":                            dataSourceIBMISVPCs(),
			"ibm_is_vpn_gateways":                    dataSourceIBMISVPNGateways(),
			"ibm_is_vpc_address_prefixes":            dataSourceIbmIsVpcAddressPrefixes(),
			"ibm_is_vpn_gateway_connections":         dataSourceIBMISVPNGatewayConnections(),
			"ibm_is_vpc_default_routing_table":       dataSourceIBMISVPCDefaultRoutingTable(),
			"ibm_is_vpc_routing_tables":              dataSourceIBMISVPCRoutingTables(),
			"ibm_is_vpc_routing_table_routes":        dataSourceIBMISVPCRoutingTableRoutes(),
			"ibm_is_zone":                            dataSourceIBMISZone(),
			"ibm_is_zones":                           dataSourceIBMISZones(),
			"ibm_is_operating_system":                dataSourceIBMISOperatingSystem(),
			"ibm_is_operating_systems":               dataSourceIBMISOperatingSystems(),
			"ibm_is_network_acls":                    dataSourceIBMIsNetworkAcls(),
			"ibm_is_network_acl":                     dataSourceIBMIsNetworkACL(),
			"ibm_is_network_acl_rule":                dataSourceIBMISNetworkACLRule(),
			"ibm_is_network_acl_rules":               dataSourceIBMISNetworkACLRules(),
			"ibm_lbaas":                              dataSourceIBMLbaas(),
			"ibm_network_vlan":                       dataSourceIBMNetworkVlan(),
			"ibm_org":                                dataSourceIBMOrg(),
			"ibm_org_quota":                          dataSourceIBMOrgQuota(),
			"ibm_kp_key":                             dataSourceIBMkey(),
			"ibm_kms_key_rings":                      dataSourceIBMKMSkeyRings(),
			"ibm_kms_key_policies":                   dataSourceIBMKMSkeyPolicies(),
			"ibm_kms_keys":                           dataSourceIBMKMSkeys(),
			"ibm_pn_application_chrome":              dataSourceIBMPNApplicationChrome(),
			"ibm_app_config_environment":             dataSourceIbmAppConfigEnvironment(),
			"ibm_app_config_environments":            dataSourceIbmAppConfigEnvironments(),
			"ibm_app_config_feature":                 dataSourceIbmAppConfigFeature(),
			"ibm_app_config_features":                dataSourceIbmAppConfigFeatures(),
			"ibm_kms_key":                            dataSourceIBMKMSkey(),
			"ibm_resource_quota":                     dataSourceIBMResourceQuota(),
			"ibm_resource_group":                     dataSourceIBMResourceGroup(),
			"ibm_resource_instance":                  dataSourceIBMResourceInstance(),
			"ibm_resource_key":                       dataSourceIBMResourceKey(),
			"ibm_security_group":                     dataSourceIBMSecurityGroup(),
			"ibm_service_instance":                   dataSourceIBMServiceInstance(),
			"ibm_service_key":                        dataSourceIBMServiceKey(),
			"ibm_service_plan":                       dataSourceIBMServicePlan(),
			"ibm_space":                              dataSourceIBMSpace(),

			// Added for Schematics
			"ibm_schematics_workspace":      dataSourceIBMSchematicsWorkspace(),
			"ibm_schematics_output":         dataSourceIBMSchematicsOutput(),
			"ibm_schematics_state":          dataSourceIBMSchematicsState(),
			"ibm_schematics_action":         dataSourceIBMSchematicsAction(),
			"ibm_schematics_job":            dataSourceIBMSchematicsJob(),
			"ibm_schematics_inventory":      dataSourceIBMSchematicsInventory(),
			"ibm_schematics_resource_query": dataSourceIBMSchematicsResourceQuery(),

			// Added for Power Resources

			"ibm_pi_key":                dataSourceIBMPIKey(),
			"ibm_pi_keys":               dataSourceIBMPIKeys(),
			"ibm_pi_image":              dataSourceIBMPIImage(),
			"ibm_pi_instance":           dataSourceIBMPIInstance(),
			"ibm_pi_tenant":             dataSourceIBMPITenant(),
			"ibm_pi_network":            dataSourceIBMPINetwork(),
			"ibm_pi_volume":             dataSourceIBMPIVolume(),
			"ibm_pi_instance_volumes":   dataSourceIBMPIInstanceVolumes(),
			"ibm_pi_public_network":     dataSourceIBMPIPublicNetwork(),
			"ibm_pi_images":             dataSourceIBMPIImages(),
			"ibm_pi_instance_ip":        dataSourceIBMPIInstanceIP(),
			"ibm_pi_instance_snapshots": dataSourceIBMPISnapshots(),
			"ibm_pi_pvm_snapshots":      dataSourceIBMPISnapshot(),
			"ibm_pi_network_port":       dataSourceIBMPINetworkPort(),
			"ibm_pi_cloud_instance":     dataSourceIBMPICloudInstance(),
			"ibm_pi_catalog_images":     dataSourceIBMPICatalogImages(),
			"ibm_pi_dhcp":               dataSourceIBMPIDhcp(),
			"ibm_pi_dhcps":              dataSourceIBMPIDhcps(),
			"ibm_pi_cloud_connection":   dataSourceIBMPICloudConnection(),
			"ibm_pi_cloud_connections":  dataSourceIBMPICloudConnections(),
			"ibm_pi_sap_profiles":       dataSourceIBMPISAPProfiles(),
			"ibm_pi_sap_profile":        dataSourceIBMPISAPProfile(),
			"ibm_pi_console_languages":  dataSourceIBMPIInstanceConsoleLanguages(),
			"ibm_pi_placement_group":    dataSourceIBMPIPlacementGroup(),
			"ibm_pi_placement_groups":   dataSourceIBMPIPlacementGroups(),

			// Added for private dns zones

			"ibm_dns_zones":                            dataSourceIBMPrivateDNSZones(),
			"ibm_dns_permitted_networks":               dataSourceIBMPrivateDNSPermittedNetworks(),
			"ibm_dns_resource_records":                 dataSourceIBMPrivateDNSResourceRecords(),
			"ibm_dns_glb_monitors":                     dataSourceIBMPrivateDNSGLBMonitors(),
			"ibm_dns_glb_pools":                        dataSourceIBMPrivateDNSGLBPools(),
			"ibm_dns_glbs":                             dataSourceIBMPrivateDNSGLBs(),
			"ibm_dns_custom_resolvers":                 dataSourceIBMDNSCustomResolver(),
			"ibm_dns_custom_resolver_forwarding_rules": dataSourceIBMPrivateDNSForwardingRules(),

			// Added for Direct Link

			"ibm_dl_gateways":          dataSourceIBMDLGateways(),
			"ibm_dl_offering_speeds":   dataSourceIBMDLOfferingSpeeds(),
			"ibm_dl_port":              dataSourceIBMDirectLinkPort(),
			"ibm_dl_ports":             dataSourceIBMDirectLinkPorts(),
			"ibm_dl_gateway":           dataSourceIBMDLGateway(),
			"ibm_dl_locations":         dataSourceIBMDLLocations(),
			"ibm_dl_routers":           dataSourceIBMDLRouters(),
			"ibm_dl_provider_ports":    dataSourceIBMDirectLinkProviderPorts(),
			"ibm_dl_provider_gateways": dataSourceIBMDirectLinkProviderGateways(),

			//Added for Transit Gateway
			"ibm_tg_gateway":   dataSourceIBMTransitGateway(),
			"ibm_tg_gateways":  dataSourceIBMTransitGateways(),
			"ibm_tg_locations": dataSourceIBMTransitGatewaysLocations(),
			"ibm_tg_location":  dataSourceIBMTransitGatewaysLocation(),

			//Added for BSS Enterprise
			"ibm_enterprises":               dataSourceIbmEnterprises(),
			"ibm_enterprise_account_groups": dataSourceIbmEnterpriseAccountGroups(),
			"ibm_enterprise_accounts":       dataSourceIbmEnterpriseAccounts(),

			//Added for Secrets Manager
			"ibm_secrets_manager_secrets": dataSourceIBMSecretsManagerSecrets(),
			"ibm_secrets_manager_secret":  dataSourceIBMSecretsManagerSecret(),

			//Added for Satellite
			"ibm_satellite_location":                            dataSourceIBMSatelliteLocation(),
			"ibm_satellite_location_nlb_dns":                    dataSourceIBMSatelliteLocationNLBDNS(),
			"ibm_satellite_attach_host_script":                  dataSourceIBMSatelliteAttachHostScript(),
			"ibm_satellite_cluster":                             dataSourceIBMSatelliteCluster(),
			"ibm_satellite_cluster_worker_pool":                 dataSourceIBMSatelliteClusterWorkerPool(),
			"ibm_satellite_link":                                dataSourceIbmSatelliteLink(),
			"ibm_satellite_endpoint":                            dataSourceIbmSatelliteEndpoint(),
			"ibm_satellite_cluster_worker_pool_zone_attachment": dataSourceIBMSatelliteClusterWorkerPoolAttachment(),

			// Catalog related resources
			"ibm_cm_catalog":           dataSourceIBMCmCatalog(),
			"ibm_cm_offering":          dataSourceIBMCmOffering(),
			"ibm_cm_version":           dataSourceIBMCmVersion(),
			"ibm_cm_offering_instance": dataSourceIBMCmOfferingInstance(),

			//Added for Resource Tag
			"ibm_resource_tag": dataSourceIBMResourceTag(),

			// Atracker
			"ibm_atracker_targets":   dataSourceIBMAtrackerTargets(),
			"ibm_atracker_routes":    dataSourceIBMAtrackerRoutes(),
			"ibm_atracker_endpoints": dataSourceIBMAtrackerEndpoints(),

			//Security and Compliance Center
			"ibm_scc_si_providers":      dataSourceIBMSccSiProviders(),
			"ibm_scc_si_note":           dataSourceIBMSccSiNote(),
			"ibm_scc_si_notes":          dataSourceIBMSccSiNotes(),
			"ibm_scc_account_location":  dataSourceIbmSccAccountLocation(),
			"ibm_scc_account_locations": dataSourceIbmSccAccountLocations(),
			"ibm_scc_account_settings":  dataSourceIbmSccAccountLocationSettings(),
			"ibm_scc_si_occurrence":     dataSourceIBMSccSiOccurrence(),
			"ibm_scc_si_occurrences":    dataSourceIBMSccSiOccurrences(),

			// Compliance Posture Management
			"ibm_scc_posture_scopes":            dataSourceIBMSccPostureScopes(),
			"ibm_scc_posture_latest_scans":      dataSourceIBMSccPostureLatestScans(),
			"ibm_scc_posture_profiles":          dataSourceIBMSccPostureProfiles(),
			"ibm_scc_posture_scan_summary":      dataSourceIBMSccPostureScansSummary(),
			"ibm_scc_posture_scan_summaries":    dataSourceIBMSccPostureScanSummaries(),
			"ibm_scc_posture_profile":           dataSourceIBMSccPostureProfileDetails(),
			"ibm_scc_posture_group_profile":     dataSourceIBMSccPostureGroupProfileDetails(),
			"ibm_scc_posture_scope_correlation": dataSourceIBMSccPostureScopeCorrelation(),

			// Added for Context Based Restrictions
			"ibm_cbr_zone": dataSourceIBMCbrZone(),
			"ibm_cbr_rule": dataSourceIBMCbrRule(),

			// Added for Event Notifications
			"ibm_en_destination":   dataSourceIBMEnDestination(),
			"ibm_en_destinations":  dataSourceIBMEnDestinations(),
			"ibm_en_topic":         dataSourceIBMEnTopic(),
			"ibm_en_topics":        dataSourceIBMEnTopics(),
			"ibm_en_subscription":  dataSourceIBMEnSubscription(),
			"ibm_en_subscriptions": dataSourceIBMEnSubscriptions(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ibm_api_gateway_endpoint":              resourceIBMApiGatewayEndPoint(),
			"ibm_api_gateway_endpoint_subscription": resourceIBMApiGatewayEndpointSubscription(),
			"ibm_app":                               resourceIBMApp(),
			"ibm_app_domain_private":                resourceIBMAppDomainPrivate(),
			"ibm_app_domain_shared":                 resourceIBMAppDomainShared(),
			"ibm_app_route":                         resourceIBMAppRoute(),

			// AppID
			"ibm_appid_action_url":               resourceIBMAppIDActionURL(),
			"ibm_appid_apm":                      resourceIBMAppIDAPM(),
			"ibm_appid_application":              resourceIBMAppIDApplication(),
			"ibm_appid_application_scopes":       resourceIBMAppIDApplicationScopes(),
			"ibm_appid_application_roles":        resourceIBMAppIDApplicationRoles(),
			"ibm_appid_audit_status":             resourceIBMAppIDAuditStatus(),
			"ibm_appid_cloud_directory_template": resourceIBMAppIDCloudDirectoryTemplate(),
			"ibm_appid_cloud_directory_user":     resourceIBMAppIDCloudDirectoryUser(),
			"ibm_appid_idp_cloud_directory":      resourceIBMAppIDIDPCloudDirectory(),
			"ibm_appid_idp_custom":               resourceIBMAppIDIDPCustom(),
			"ibm_appid_idp_facebook":             resourceIBMAppIDIDPFacebook(),
			"ibm_appid_idp_google":               resourceIBMAppIDIDPGoogle(),
			"ibm_appid_idp_saml":                 resourceIBMAppIDIDPSAML(),
			"ibm_appid_languages":                resourceIBMAppIDLanguages(),
			"ibm_appid_mfa":                      resourceIBMAppIDMFA(),
			"ibm_appid_mfa_channel":              resourceIBMAppIDMFAChannel(),
			"ibm_appid_password_regex":           resourceIBMAppIDPasswordRegex(),
			"ibm_appid_token_config":             resourceIBMAppIDTokenConfig(),
			"ibm_appid_redirect_urls":            resourceIBMAppIDRedirectURLs(),
			"ibm_appid_role":                     resourceIBMAppIDRole(),
			"ibm_appid_theme_color":              resourceIBMAppIDThemeColor(),
			"ibm_appid_theme_text":               resourceIBMAppIDThemeText(),
			"ibm_appid_user_roles":               resourceIBMAppIDUserRoles(),

			"ibm_function_action":                                resourceIBMFunctionAction(),
			"ibm_function_package":                               resourceIBMFunctionPackage(),
			"ibm_function_rule":                                  resourceIBMFunctionRule(),
			"ibm_function_trigger":                               resourceIBMFunctionTrigger(),
			"ibm_function_namespace":                             resourceIBMFunctionNamespace(),
			"ibm_cis":                                            resourceIBMCISInstance(),
			"ibm_database":                                       resourceIBMDatabaseInstance(),
			"ibm_certificate_manager_import":                     resourceIBMCertificateManagerImport(),
			"ibm_certificate_manager_order":                      resourceIBMCertificateManagerOrder(),
			"ibm_cis_domain":                                     resourceIBMCISDomain(),
			"ibm_cis_domain_settings":                            resourceIBMCISSettings(),
			"ibm_cis_firewall":                                   resourceIBMCISFirewallRecord(),
			"ibm_cis_range_app":                                  resourceIBMCISRangeApp(),
			"ibm_cis_healthcheck":                                resourceIBMCISHealthCheck(),
			"ibm_cis_origin_pool":                                resourceIBMCISPool(),
			"ibm_cis_global_load_balancer":                       resourceIBMCISGlb(),
			"ibm_cis_certificate_upload":                         resourceIBMCISCertificateUpload(),
			"ibm_cis_dns_record":                                 resourceIBMCISDnsRecord(),
			"ibm_cis_dns_records_import":                         resourceIBMCISDNSRecordsImport(),
			"ibm_cis_rate_limit":                                 resourceIBMCISRateLimit(),
			"ibm_cis_page_rule":                                  resourceIBMCISPageRule(),
			"ibm_cis_edge_functions_action":                      resourceIBMCISEdgeFunctionsAction(),
			"ibm_cis_edge_functions_trigger":                     resourceIBMCISEdgeFunctionsTrigger(),
			"ibm_cis_tls_settings":                               resourceIBMCISTLSSettings(),
			"ibm_cis_waf_package":                                resourceIBMCISWAFPackage(),
			"ibm_cis_routing":                                    resourceIBMCISRouting(),
			"ibm_cis_waf_group":                                  resourceIBMCISWAFGroup(),
			"ibm_cis_cache_settings":                             resourceIBMCISCacheSettings(),
			"ibm_cis_custom_page":                                resourceIBMCISCustomPage(),
			"ibm_cis_waf_rule":                                   resourceIBMCISWAFRule(),
			"ibm_cis_certificate_order":                          resourceIBMCISCertificateOrder(),
			"ibm_cis_filter":                                     resourceIBMCISFilter(),
			"ibm_cis_firewall_rule":                              resourceIBMCISFirewallrules(),
			"ibm_cloudant":                                       resourceIBMCloudant(),
			"ibm_cloud_shell_account_settings":                   resourceIBMCloudShellAccountSettings(),
			"ibm_compute_autoscale_group":                        resourceIBMComputeAutoScaleGroup(),
			"ibm_compute_autoscale_policy":                       resourceIBMComputeAutoScalePolicy(),
			"ibm_compute_bare_metal":                             resourceIBMComputeBareMetal(),
			"ibm_compute_dedicated_host":                         resourceIBMComputeDedicatedHost(),
			"ibm_compute_monitor":                                resourceIBMComputeMonitor(),
			"ibm_compute_placement_group":                        resourceIBMComputePlacementGroup(),
			"ibm_compute_reserved_capacity":                      resourceIBMComputeReservedCapacity(),
			"ibm_compute_provisioning_hook":                      resourceIBMComputeProvisioningHook(),
			"ibm_compute_ssh_key":                                resourceIBMComputeSSHKey(),
			"ibm_compute_ssl_certificate":                        resourceIBMComputeSSLCertificate(),
			"ibm_compute_user":                                   resourceIBMComputeUser(),
			"ibm_compute_vm_instance":                            resourceIBMComputeVmInstance(),
			"ibm_container_addons":                               resourceIBMContainerAddOns(),
			"ibm_container_alb":                                  resourceIBMContainerALB(),
			"ibm_container_alb_create":                           resourceIBMContainerAlbCreate(),
			"ibm_container_api_key_reset":                        resourceIBMContainerAPIKeyReset(),
			"ibm_container_vpc_alb":                              resourceIBMContainerVpcALB(),
			"ibm_container_vpc_alb_create":                       resourceIBMContainerVpcAlbCreateNew(),
			"ibm_container_vpc_worker_pool":                      resourceIBMContainerVpcWorkerPool(),
			"ibm_container_vpc_cluster":                          resourceIBMContainerVpcCluster(),
			"ibm_container_alb_cert":                             resourceIBMContainerALBCert(),
			"ibm_container_cluster":                              resourceIBMContainerCluster(),
			"ibm_container_cluster_feature":                      resourceIBMContainerClusterFeature(),
			"ibm_container_bind_service":                         resourceIBMContainerBindService(),
			"ibm_container_worker_pool":                          resourceIBMContainerWorkerPool(),
			"ibm_container_worker_pool_zone_attachment":          resourceIBMContainerWorkerPoolZoneAttachment(),
			"ibm_container_storage_attachment":                   resourceIBMContainerVpcWorkerVolumeAttachment(),
			"ibm_container_nlb_dns":                              resourceIbmContainerNlbDns(),
			"ibm_cr_namespace":                                   resourceIBMCrNamespace(),
			"ibm_cr_retention_policy":                            resourceIBMCrRetentionPolicy(),
			"ibm_ob_logging":                                     resourceIBMObLogging(),
			"ibm_ob_monitoring":                                  resourceIBMObMonitoring(),
			"ibm_cos_bucket":                                     resourceIBMCOSBucket(),
			"ibm_cos_bucket_object":                              resourceIBMCOSBucketObject(),
			"ibm_dns_domain":                                     resourceIBMDNSDomain(),
			"ibm_dns_domain_registration_nameservers":            resourceIBMDNSDomainRegistrationNameservers(),
			"ibm_dns_secondary":                                  resourceIBMDNSSecondary(),
			"ibm_dns_record":                                     resourceIBMDNSRecord(),
			"ibm_event_streams_topic":                            resourceIBMEventStreamsTopic(),
			"ibm_event_streams_schema":                           resourceIBMEventStreamsSchema(),
			"ibm_firewall":                                       resourceIBMFirewall(),
			"ibm_firewall_policy":                                resourceIBMFirewallPolicy(),
			"ibm_hpcs":                                           resourceIBMHPCS(),
			"ibm_iam_access_group":                               resourceIBMIAMAccessGroup(),
			"ibm_iam_account_settings":                           resourceIbmIamAccountSettings(),
			"ibm_iam_custom_role":                                resourceIBMIAMCustomRole(),
			"ibm_iam_access_group_dynamic_rule":                  resourceIBMIAMDynamicRule(),
			"ibm_iam_access_group_members":                       resourceIBMIAMAccessGroupMembers(),
			"ibm_iam_access_group_policy":                        resourceIBMIAMAccessGroupPolicy(),
			"ibm_iam_authorization_policy":                       resourceIBMIAMAuthorizationPolicy(),
			"ibm_iam_authorization_policy_detach":                resourceIBMIAMAuthorizationPolicyDetach(),
			"ibm_iam_user_policy":                                resourceIBMIAMUserPolicy(),
			"ibm_iam_user_settings":                              resourceIBMUserSettings(),
			"ibm_iam_service_id":                                 resourceIBMIAMServiceID(),
			"ibm_iam_service_api_key":                            resourceIBMIAMServiceAPIKey(),
			"ibm_iam_service_policy":                             resourceIBMIAMServicePolicy(),
			"ibm_iam_user_invite":                                resourceIBMUserInvite(),
			"ibm_iam_api_key":                                    resourceIbmIamApiKey(),
			"ibm_iam_trusted_profile":                            resourceIBMIamTrustedProfile(),
			"ibm_iam_trusted_profile_claim_rule":                 resourceIBMIamTrustedProfileClaimRule(),
			"ibm_iam_trusted_profile_link":                       resourceIBMIamTrustedProfileLink(),
			"ibm_iam_trusted_profile_policy":                     resourceIBMIAMTrustedProfilePolicy(),
			"ibm_ipsec_vpn":                                      resourceIBMIPSecVPN(),
			"ibm_is_dedicated_host":                              resourceIbmIsDedicatedHost(),
			"ibm_is_dedicated_host_group":                        resourceIbmIsDedicatedHostGroup(),
			"ibm_is_dedicated_host_disk_management":              resourceIBMISDedicatedHostDiskManagement(),
			"ibm_is_placement_group":                             resourceIbmIsPlacementGroup(),
			"ibm_is_floating_ip":                                 resourceIBMISFloatingIP(),
			"ibm_is_flow_log":                                    resourceIBMISFlowLog(),
			"ibm_is_instance":                                    resourceIBMISInstance(),
			"ibm_is_instance_action":                             resourceIBMISInstanceAction(),
			"ibm_is_instance_network_interface":                  resourceIBMIsInstanceNetworkInterface(),
			"ibm_is_instance_disk_management":                    resourceIBMISInstanceDiskManagement(),
			"ibm_is_instance_group":                              resourceIBMISInstanceGroup(),
			"ibm_is_instance_group_membership":                   resourceIBMISInstanceGroupMembership(),
			"ibm_is_instance_group_manager":                      resourceIBMISInstanceGroupManager(),
			"ibm_is_instance_group_manager_policy":               resourceIBMISInstanceGroupManagerPolicy(),
			"ibm_is_instance_group_manager_action":               resourceIBMISInstanceGroupManagerAction(),
			"ibm_is_instance_volume_attachment":                  resourceIBMISInstanceVolumeAttachment(),
			"ibm_is_virtual_endpoint_gateway":                    resourceIBMISEndpointGateway(),
			"ibm_is_virtual_endpoint_gateway_ip":                 resourceIBMISEndpointGatewayIP(),
			"ibm_is_instance_template":                           resourceIBMISInstanceTemplate(),
			"ibm_is_ike_policy":                                  resourceIBMISIKEPolicy(),
			"ibm_is_ipsec_policy":                                resourceIBMISIPSecPolicy(),
			"ibm_is_lb":                                          resourceIBMISLB(),
			"ibm_is_lb_listener":                                 resourceIBMISLBListener(),
			"ibm_is_lb_listener_policy":                          resourceIBMISLBListenerPolicy(),
			"ibm_is_lb_listener_policy_rule":                     resourceIBMISLBListenerPolicyRule(),
			"ibm_is_lb_pool":                                     resourceIBMISLBPool(),
			"ibm_is_lb_pool_member":                              resourceIBMISLBPoolMember(),
			"ibm_is_network_acl":                                 resourceIBMISNetworkACL(),
			"ibm_is_network_acl_rule":                            resourceIBMISNetworkACLRule(),
			"ibm_is_public_gateway":                              resourceIBMISPublicGateway(),
			"ibm_is_security_group":                              resourceIBMISSecurityGroup(),
			"ibm_is_security_group_rule":                         resourceIBMISSecurityGroupRule(),
			"ibm_is_security_group_target":                       resourceIBMISSecurityGroupTarget(),
			"ibm_is_security_group_network_interface_attachment": resourceIBMISSecurityGroupNetworkInterfaceAttachment(),
			"ibm_is_subnet":                                      resourceIBMISSubnet(),
			"ibm_is_subnet_reserved_ip":                          resourceIBMISReservedIP(),
			"ibm_is_subnet_network_acl_attachment":               resourceIBMISSubnetNetworkACLAttachment(),
			"ibm_is_ssh_key":                                     resourceIBMISSSHKey(),
			"ibm_is_snapshot":                                    resourceIBMSnapshot(),
			"ibm_is_volume":                                      resourceIBMISVolume(),
			"ibm_is_vpn_gateway":                                 resourceIBMISVPNGateway(),
			"ibm_is_vpn_gateway_connection":                      resourceIBMISVPNGatewayConnection(),
			"ibm_is_vpc":                                         resourceIBMISVPC(),
			"ibm_is_vpc_address_prefix":                          resourceIBMISVpcAddressPrefix(),
			"ibm_is_vpc_route":                                   resourceIBMISVpcRoute(),
			"ibm_is_vpc_routing_table":                           resourceIBMISVPCRoutingTable(),
			"ibm_is_vpc_routing_table_route":                     resourceIBMISVPCRoutingTableRoute(),
			"ibm_is_image":                                       resourceIBMISImage(),
			"ibm_lb":                                             resourceIBMLb(),
			"ibm_lbaas":                                          resourceIBMLbaas(),
			"ibm_lbaas_health_monitor":                           resourceIBMLbaasHealthMonitor(),
			"ibm_lbaas_server_instance_attachment":               resourceIBMLbaasServerInstanceAttachment(),
			"ibm_lb_service":                                     resourceIBMLbService(),
			"ibm_lb_service_group":                               resourceIBMLbServiceGroup(),
			"ibm_lb_vpx":                                         resourceIBMLbVpx(),
			"ibm_lb_vpx_ha":                                      resourceIBMLbVpxHa(),
			"ibm_lb_vpx_service":                                 resourceIBMLbVpxService(),
			"ibm_lb_vpx_vip":                                     resourceIBMLbVpxVip(),
			"ibm_multi_vlan_firewall":                            resourceIBMMultiVlanFirewall(),
			"ibm_network_gateway":                                resourceIBMNetworkGateway(),
			"ibm_network_gateway_vlan_association":               resourceIBMNetworkGatewayVlanAttachment(),
			"ibm_network_interface_sg_attachment":                resourceIBMNetworkInterfaceSGAttachment(),
			"ibm_network_public_ip":                              resourceIBMNetworkPublicIp(),
			"ibm_network_vlan":                                   resourceIBMNetworkVlan(),
			"ibm_network_vlan_spanning":                          resourceIBMNetworkVlanSpan(),
			"ibm_object_storage_account":                         resourceIBMObjectStorageAccount(),
			"ibm_org":                                            resourceIBMOrg(),
			"ibm_pn_application_chrome":                          resourceIBMPNApplicationChrome(),
			"ibm_app_config_environment":                         resourceIbmAppConfigEnvironment(),
			"ibm_app_config_feature":                             resourceIbmIbmAppConfigFeature(),
			"ibm_kms_key":                                        resourceIBMKmskey(),
			"ibm_kms_key_alias":                                  resourceIBMKmskeyAlias(),
			"ibm_kms_key_rings":                                  resourceIBMKmskeyRings(),
			"ibm_kms_key_policies":                               resourceIBMKmskeyPolicies(),
			"ibm_kp_key":                                         resourceIBMkey(),
			"ibm_resource_group":                                 resourceIBMResourceGroup(),
			"ibm_resource_instance":                              resourceIBMResourceInstance(),
			"ibm_resource_key":                                   resourceIBMResourceKey(),
			"ibm_security_group":                                 resourceIBMSecurityGroup(),
			"ibm_security_group_rule":                            resourceIBMSecurityGroupRule(),
			"ibm_service_instance":                               resourceIBMServiceInstance(),
			"ibm_service_key":                                    resourceIBMServiceKey(),
			"ibm_space":                                          resourceIBMSpace(),
			"ibm_storage_evault":                                 resourceIBMStorageEvault(),
			"ibm_storage_block":                                  resourceIBMStorageBlock(),
			"ibm_storage_file":                                   resourceIBMStorageFile(),
			"ibm_subnet":                                         resourceIBMSubnet(),
			"ibm_dns_reverse_record":                             resourceIBMDNSReverseRecord(),
			"ibm_ssl_certificate":                                resourceIBMSSLCertificate(),
			"ibm_cdn":                                            resourceIBMCDN(),
			"ibm_hardware_firewall_shared":                       resourceIBMFirewallShared(),

			//Added for Power Colo

			"ibm_pi_key":                 resourceIBMPIKey(),
			"ibm_pi_volume":              resourceIBMPIVolume(),
			"ibm_pi_network":             resourceIBMPINetwork(),
			"ibm_pi_instance":            resourceIBMPIInstance(),
			"ibm_pi_operations":          resourceIBMPIIOperations(),
			"ibm_pi_volume_attach":       resourceIBMPIVolumeAttach(),
			"ibm_pi_capture":             resourceIBMPICapture(),
			"ibm_pi_image":               resourceIBMPIImage(),
			"ibm_pi_image_export":        resourceIBMPIImageExport(),
			"ibm_pi_network_port":        resourceIBMPINetworkPort(),
			"ibm_pi_snapshot":            resourceIBMPISnapshot(),
			"ibm_pi_network_port_attach": resourceIBMPINetworkPortAttach(),
			"ibm_pi_dhcp":                resourceIBMPIDhcp(),
			"ibm_pi_cloud_connection":    resourceIBMPICloudConnection(),
			"ibm_pi_ike_policy":          resourceIBMPIIKEPolicy(),
			"ibm_pi_ipsec_policy":        resourceIBMPIIPSecPolicy(),
			"ibm_pi_vpn_connection":      resourceIBMPIVPNConnection(),
			"ibm_pi_console_language":    resourceIBMPIInstanceConsoleLanguage(),
			"ibm_pi_placement_group":     resourceIBMPIPlacementGroup(),

			//Private DNS related resources
			"ibm_dns_zone":              resourceIBMPrivateDNSZone(),
			"ibm_dns_permitted_network": resourceIBMPrivateDNSPermittedNetwork(),
			"ibm_dns_resource_record":   resourceIBMPrivateDNSResourceRecord(),
			"ibm_dns_glb_monitor":       resourceIBMPrivateDNSGLBMonitor(),
			"ibm_dns_glb_pool":          resourceIBMPrivateDNSGLBPool(),
			"ibm_dns_glb":               resourceIBMPrivateDNSGLB(),

			//Added for Custom Resolver
			"ibm_dns_custom_resolver":                 resouceIBMPrivateDNSCustomResolver(),
			"ibm_dns_custom_resolver_location":        resourceIBMPrivateDNSCRLocation(),
			"ibm_dns_custom_resolver_forwarding_rule": resourceIBMPrivateDNSForwardingRule(),

			//Direct Link related resources
			"ibm_dl_gateway":            resourceIBMDLGateway(),
			"ibm_dl_virtual_connection": resourceIBMDLGatewayVC(),
			"ibm_dl_provider_gateway":   resourceIBMDLProviderGateway(),
			//Added for Transit Gateway
			"ibm_tg_gateway":    resourceIBMTransitGateway(),
			"ibm_tg_connection": resourceIBMTransitGatewayConnection(),

			//Catalog related resources
			"ibm_cm_offering_instance": resourceIBMCmOfferingInstance(),
			"ibm_cm_catalog":           resourceIBMCmCatalog(),
			"ibm_cm_offering":          resourceIBMCmOffering(),
			"ibm_cm_version":           resourceIBMCmVersion(),

			//Added for enterprise
			"ibm_enterprise":               resourceIbmEnterprise(),
			"ibm_enterprise_account_group": resourceIbmEnterpriseAccountGroup(),
			"ibm_enterprise_account":       resourceIbmEnterpriseAccount(),

			//Added for Schematics
			"ibm_schematics_workspace":      resourceIBMSchematicsWorkspace(),
			"ibm_schematics_action":         resourceIBMSchematicsAction(),
			"ibm_schematics_job":            resourceIBMSchematicsJob(),
			"ibm_schematics_inventory":      resourceIBMSchematicsInventory(),
			"ibm_schematics_resource_query": resourceIBMSchematicsResourceQuery(),

			//satellite  resources
			"ibm_satellite_location":                            resourceIBMSatelliteLocation(),
			"ibm_satellite_host":                                resourceIBMSatelliteHost(),
			"ibm_satellite_cluster":                             resourceIBMSatelliteCluster(),
			"ibm_satellite_cluster_worker_pool":                 resourceIBMSatelliteClusterWorkerPool(),
			"ibm_satellite_link":                                resourceIbmSatelliteLink(),
			"ibm_satellite_endpoint":                            resourceIbmSatelliteEndpoint(),
			"ibm_satellite_location_nlb_dns":                    resourceIbmSatelliteLocationNlbDns(),
			"ibm_satellite_cluster_worker_pool_zone_attachment": resourceIbmSatelliteClusterWorkerPoolZoneAttachment(),

			//Added for Resource Tag
			"ibm_resource_tag": resourceIBMResourceTag(),

			// Atracker
			"ibm_atracker_target": resourceIBMAtrackerTarget(),
			"ibm_atracker_route":  resourceIBMAtrackerRoute(),

			//Security and Compliance Center
			"ibm_scc_si_note":          resourceIBMSccSiNote(),
			"ibm_scc_account_settings": resourceIBMSccAccountSettings(),
			"ibm_scc_si_occurrence":    resourceIBMSccSiOccurrence(),

			//Security and Compliance Center - PostureManagement
			"ibm_scc_posture_collector":  resourceIBMSccPostureCollectors(),
			"ibm_scc_posture_scope":      resourceIBMSccPostureScopes(),
			"ibm_scc_posture_credential": resourceIBMSccPostureCredentials(),

			// Added for Context Based Restrictions
			"ibm_cbr_zone": resourceIBMCbrZone(),
			"ibm_cbr_rule": resourceIBMCbrRule(),

			// Added for Event Notifications
			"ibm_en_destination":  resourceIBMEnDestination(),
			"ibm_en_topic":        resourceIBMEnTopic(),
			"ibm_en_subscription": resourceIBMEnSubscription(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var globalValidatorDict ValidatorDict
var initOnce sync.Once

// Validator return validator
func Validator() ValidatorDict {
	initOnce.Do(func() {
		globalValidatorDict = ValidatorDict{
			ResourceValidatorDictionary: map[string]*ResourceValidator{
				"ibm_iam_account_settings":                resourceIBMIAMAccountSettingsValidator(),
				"ibm_iam_custom_role":                     resourceIBMIAMCustomRoleValidator(),
				"ibm_cis_healthcheck":                     resourceIBMCISHealthCheckValidator(),
				"ibm_cis_rate_limit":                      resourceIBMCISRateLimitValidator(),
				"ibm_cis":                                 resourceIBMCISValidator(),
				"ibm_cis_domain_settings":                 resourceIBMCISDomainSettingValidator(),
				"ibm_cis_tls_settings":                    resourceIBMCISTLSSettingsValidator(),
				"ibm_cis_routing":                         resourceIBMCISRoutingValidator(),
				"ibm_cis_page_rule":                       resourceCISPageRuleValidator(),
				"ibm_cis_waf_package":                     resourceIBMCISWAFPackageValidator(),
				"ibm_cis_waf_group":                       resourceIBMCISWAFGroupValidator(),
				"ibm_cis_certificate_upload":              resourceCISCertificateUploadValidator(),
				"ibm_cis_cache_settings":                  resourceIBMCISCacheSettingsValidator(),
				"ibm_cis_custom_page":                     resourceIBMCISCustomPageValidator(),
				"ibm_cis_firewall":                        resourceIBMCISFirewallValidator(),
				"ibm_cis_range_app":                       resourceIBMCISRangeAppValidator(),
				"ibm_cis_waf_rule":                        resourceIBMCISWAFRuleValidator(),
				"ibm_cis_certificate_order":               resourceIBMCISCertificateOrderValidator(),
				"ibm_cis_filter":                          resourceIBMCISFilterValidator(),
				"ibm_cis_firewall_rules":                  resourceIBMCISFirewallrulesValidator(),
				"ibm_container_cluster":                   resourceIBMContainerClusterValidator(),
				"ibm_container_worker_pool":               resourceContainerWorkerPoolValidator(),
				"ibm_container_vpc_worker_pool":           resourceContainerVPCWorkerPoolValidator(),
				"ibm_container_vpc_cluster":               resourceIBMContainerVpcClusterValidator(),
				"ibm_cr_namespace":                        resourceIBMCrNamespaceValidator(),
				"ibm_tg_gateway":                          resourceIBMTGValidator(),
				"ibm_app_config_feature":                  resourceIbmAppConfigFeatureValidator(),
				"ibm_tg_connection":                       resourceIBMTransitGatewayConnectionValidator(),
				"ibm_dl_virtual_connection":               resourceIBMdlGatewayVCValidator(),
				"ibm_dl_gateway":                          resourceIBMDLGatewayValidator(),
				"ibm_dl_provider_gateway":                 resourceIBMDLProviderGatewayValidator(),
				"ibm_database":                            resourceIBMICDValidator(),
				"ibm_function_package":                    resourceIBMFuncPackageValidator(),
				"ibm_function_action":                     resourceIBMFuncActionValidator(),
				"ibm_function_rule":                       resourceIBMFuncRuleValidator(),
				"ibm_function_trigger":                    resourceIBMFuncTriggerValidator(),
				"ibm_function_namespace":                  resourceIBMFuncNamespaceValidator(),
				"ibm_hpcs":                                resourceIBMHPCSValidator(),
				"ibm_is_dedicated_host_group":             resourceIbmIsDedicatedHostGroupValidator(),
				"ibm_is_dedicated_host":                   resourceIbmIsDedicatedHostValidator(),
				"ibm_is_dedicated_host_disk_management":   resourceIBMISDedicatedHostDiskManagementValidator(),
				"ibm_is_flow_log":                         resourceIBMISFlowLogValidator(),
				"ibm_is_instance_group":                   resourceIBMISInstanceGroupValidator(),
				"ibm_is_instance_group_membership":        resourceIBMISInstanceGroupMembershipValidator(),
				"ibm_is_instance_group_manager":           resourceIBMISInstanceGroupManagerValidator(),
				"ibm_is_instance_group_manager_policy":    resourceIBMISInstanceGroupManagerPolicyValidator(),
				"ibm_is_instance_group_manager_action":    resourceIBMISInstanceGroupManagerActionValidator(),
				"ibm_is_floating_ip":                      resourceIBMISFloatingIPValidator(),
				"ibm_is_ike_policy":                       resourceIBMISIKEValidator(),
				"ibm_is_image":                            resourceIBMISImageValidator(),
				"ibm_is_instance_template":                resourceIBMISInstanceTemplateValidator(),
				"ibm_is_instance":                         resourceIBMISInstanceValidator(),
				"ibm_is_instance_action":                  resourceIBMISInstanceActionValidator(),
				"ibm_is_instance_network_interface":       resourceIBMIsInstanceNetworkInterfaceValidator(),
				"ibm_is_instance_disk_management":         resourceIBMISInstanceDiskManagementValidator(),
				"ibm_is_instance_volume_attachment":       resourceIBMISInstanceVolumeAttachmentValidator(),
				"ibm_is_ipsec_policy":                     resourceIBMISIPSECValidator(),
				"ibm_is_lb_listener_policy_rule":          resourceIBMISLBListenerPolicyRuleValidator(),
				"ibm_is_lb_listener_policy":               resourceIBMISLBListenerPolicyValidator(),
				"ibm_is_lb_listener":                      resourceIBMISLBListenerValidator(),
				"ibm_is_lb_pool_member":                   resourceIBMISLBPoolMemberValidator(),
				"ibm_is_lb_pool":                          resourceIBMISLBPoolValidator(),
				"ibm_is_lb":                               resourceIBMISLBValidator(),
				"ibm_is_network_acl":                      resourceIBMISNetworkACLValidator(),
				"ibm_is_network_acl_rule":                 resourceIBMISNetworkACLRuleValidator(),
				"ibm_is_public_gateway":                   resourceIBMISPublicGatewayValidator(),
				"ibm_is_placement_group":                  resourceIbmIsPlacementGroupValidator(),
				"ibm_is_security_group_target":            resourceIBMISSecurityGroupTargetValidator(),
				"ibm_is_security_group_rule":              resourceIBMISSecurityGroupRuleValidator(),
				"ibm_is_security_group":                   resourceIBMISSecurityGroupValidator(),
				"ibm_is_snapshot":                         resourceIBMISSnapshotValidator(),
				"ibm_is_ssh_key":                          resourceIBMISSHKeyValidator(),
				"ibm_is_subnet":                           resourceIBMISSubnetValidator(),
				"ibm_is_subnet_reserved_ip":               resourceIBMISSubnetReservedIPValidator(),
				"ibm_is_volume":                           resourceIBMISVolumeValidator(),
				"ibm_is_address_prefix":                   resourceIBMISAddressPrefixValidator(),
				"ibm_is_route":                            resourceIBMISRouteValidator(),
				"ibm_is_vpc":                              resourceIBMISVPCValidator(),
				"ibm_is_vpc_routing_table":                resourceIBMISVPCRoutingTableValidator(),
				"ibm_is_vpc_routing_table_route":          resourceIBMISVPCRoutingTableRouteValidator(),
				"ibm_is_vpn_gateway_connection":           resourceIBMISVPNGatewayConnectionValidator(),
				"ibm_is_vpn_gateway":                      resourceIBMISVPNGatewayValidator(),
				"ibm_kms_key_rings":                       resourceIBMKeyRingValidator(),
				"ibm_dns_glb_monitor":                     resourceIBMPrivateDNSGLBMonitorValidator(),
				"ibm_dns_glb_pool":                        resourceIBMPrivateDNSGLBPoolValidator(),
				"ibm_dns_custom_resolver_forwarding_rule": resourceIBMPrivateDNSForwardingRuleValidator(),
				"ibm_schematics_action":                   resourceIBMSchematicsActionValidator(),
				"ibm_schematics_job":                      resourceIBMSchematicsJobValidator(),
				"ibm_schematics_workspace":                resourceIBMSchematicsWorkspaceValidator(),
				"ibm_schematics_inventory":                resourceIBMSchematicsInventoryValidator(),
				"ibm_schematics_resource_query":           resourceIBMSchematicsResourceQueryValidator(),
				"ibm_resource_instance":                   resourceIBMResourceInstanceValidator(),
				"ibm_is_virtual_endpoint_gateway":         resourceIBMISEndpointGatewayValidator(),
				"ibm_resource_tag":                        resourceIBMResourceTagValidator(),
				"ibm_satellite_location":                  resourceIBMSatelliteLocationValidator(),
				"ibm_satellite_cluster":                   resourceIBMSatelliteClusterValidator(),
				"ibm_pi_volume":                           resourceIBMPIVolumeValidator(),
				"ibm_atracker_target":                     resourceIBMAtrackerTargetValidator(),
				"ibm_atracker_route":                      resourceIBMAtrackerRouteValidator(),
				"ibm_satellite_endpoint":                  resourceIbmSatelliteEndpointValidator(),
				"ibm_scc_si_note":                         resourceIBMSccSiNoteValidator(),
				"ibm_scc_account_settings":                resourceIBMSccAccountSettingsValidator(),
				"ibm_scc_si_occurrence":                   resourceIBMSccSiOccurrenceValidator(),
				"ibm_cbr_zone":                            resourceIBMCbrZoneValidator(),
				"ibm_cbr_rule":                            resourceIBMCbrRuleValidator(),
				"ibm_scc_posture_collector":               resourceIBMSccPostureCollectorsValidator(),
				"ibm_scc_posture_scope":                   resourceIBMSccPostureScopesValidator(),
				"ibm_scc_posture_credential":              resourceIBMSccPostureCredentialsValidator(),

				// Added for Event Notifications
				"ibm_en_destination": resourceIBMEnDestinationValidator(),
			},
			DataSourceValidatorDictionary: map[string]*ResourceValidator{
				"ibm_is_subnet":               dataSourceIBMISSubnetValidator(),
				"ibm_is_snapshot":             dataSourceIBMISSnapshotValidator(),
				"ibm_dl_offering_speeds":      datasourceIBMDLOfferingSpeedsValidator(),
				"ibm_dl_routers":              datasourceIBMDLRoutersValidator(),
				"ibm_is_vpc":                  dataSourceIBMISVpcValidator(),
				"ibm_is_volume":               dataSourceIBMISVolumeValidator(),
				"ibm_scc_si_notes":            dataSourceIBMSccSiNotesValidator(),
				"ibm_scc_si_occurrences":      dataSourceIBMSccSiOccurrencesValidator(),
				"ibm_secrets_manager_secret":  datasourceIBMSecretsManagerSecretValidator(),
				"ibm_secrets_manager_secrets": datasourceIBMSecretsManagerSecretsValidator(),
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
	//Set environment variable to be used in DiffSupressFunction
	if wskEnvVal.(string) == "" {
		os.Setenv("FUNCTION_NAMESPACE", wskNameSpace)
	}

	config := Config{
		BluemixAPIKey:        bluemixAPIKey,
		Region:               region,
		ResourceGroup:        resourceGrp,
		BluemixTimeout:       time.Duration(bluemixTimeout) * time.Second,
		SoftLayerTimeout:     time.Duration(softlayerTimeout) * time.Second,
		SoftLayerUserName:    softlayerUsername,
		SoftLayerAPIKey:      softlayerAPIKey,
		RetryCount:           retryCount,
		SoftLayerEndpointURL: softlayerEndpointUrl,
		RetryDelay:           RetryAPIDelay,
		FunctionNameSpace:    wskNameSpace,
		RiaasEndPoint:        riaasEndPoint,
		IAMToken:             iamToken,
		IAMRefreshToken:      iamRefreshToken,
		Zone:                 zone,
		Visibility:           visibility,
		EndpointsFile:        file,
		IAMTrustedProfileID:  iamTrustedProfileId,
		//PowerServiceInstance: powerServiceInstance,
	}

	return config.ClientSession()
}
