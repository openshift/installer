package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/utils/terraform/auth"
)

// This is a global MutexKV for use within this plugin.
var osMutexKV = mutexkv.NewMutexKV()

// Use openstackbase.Config as the base/foundation of this provider's
// Config struct.
type Config struct {
	auth.Config
}

// Provider returns a schema.Provider for OpenStack.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: descriptions["auth_url"],
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"application_credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_ID", ""),
				Description: descriptions["application_credential_id"],
			},

			"application_credential_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_NAME", ""),
				Description: descriptions["application_credential_name"],
			},

			"application_credential_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_APPLICATION_CREDENTIAL_SECRET", ""),
				Description: descriptions["application_credential_secret"],
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: descriptions["token"],
			},

			"user_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", ""),
				Description: descriptions["user_domain_name"],
			},

			"user_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_DOMAIN_ID", ""),
				Description: descriptions["user_domain_id"],
			},

			"project_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_NAME", ""),
				Description: descriptions["project_domain_name"],
			},

			"project_domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PROJECT_DOMAIN_ID", ""),
				Description: descriptions["project_domain_id"],
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_ID", ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DOMAIN_NAME", ""),
				Description: descriptions["domain_name"],
			},

			"default_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DEFAULT_DOMAIN", "default"),
				Description: descriptions["default_domain"],
			},

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", nil),
				Description: descriptions["insecure"],
			},

			"endpoint_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", false),
				Description: descriptions["swauth"],
			},

			"use_octavia": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USE_OCTAVIA", false),
				Description: descriptions["use_octavia"],
			},

			"delayed_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_DELAYED_AUTH", false),
				Description: descriptions["delayed_auth"],
			},

			"allow_reauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ALLOW_REAUTH", false),
				Description: descriptions["allow_reauth"],
			},

			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: descriptions["cloud"],
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: descriptions["max_retries"],
			},

			"endpoint_overrides": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: descriptions["endpoint_overrides"],
			},

			"disable_no_cache_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["disable_no_cache_header"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"openstack_blockstorage_availability_zones_v3":       dataSourceBlockStorageAvailabilityZonesV3(),
			"openstack_blockstorage_snapshot_v2":                 dataSourceBlockStorageSnapshotV2(),
			"openstack_blockstorage_snapshot_v3":                 dataSourceBlockStorageSnapshotV3(),
			"openstack_blockstorage_volume_v2":                   dataSourceBlockStorageVolumeV2(),
			"openstack_blockstorage_volume_v3":                   dataSourceBlockStorageVolumeV3(),
			"openstack_compute_availability_zones_v2":            dataSourceComputeAvailabilityZonesV2(),
			"openstack_compute_flavor_v2":                        dataSourceComputeFlavorV2(),
			"openstack_compute_keypair_v2":                       dataSourceComputeKeypairV2(),
			"openstack_containerinfra_clustertemplate_v1":        dataSourceContainerInfraClusterTemplateV1(),
			"openstack_containerinfra_cluster_v1":                dataSourceContainerInfraCluster(),
			"openstack_dns_zone_v2":                              dataSourceDNSZoneV2(),
			"openstack_fw_policy_v1":                             dataSourceFWPolicyV1(),
			"openstack_identity_role_v3":                         dataSourceIdentityRoleV3(),
			"openstack_identity_project_v3":                      dataSourceIdentityProjectV3(),
			"openstack_identity_user_v3":                         dataSourceIdentityUserV3(),
			"openstack_identity_auth_scope_v3":                   dataSourceIdentityAuthScopeV3(),
			"openstack_identity_endpoint_v3":                     dataSourceIdentityEndpointV3(),
			"openstack_identity_service_v3":                      dataSourceIdentityServiceV3(),
			"openstack_identity_group_v3":                        dataSourceIdentityGroupV3(),
			"openstack_images_image_v2":                          dataSourceImagesImageV2(),
			"openstack_networking_addressscope_v2":               dataSourceNetworkingAddressScopeV2(),
			"openstack_networking_network_v2":                    dataSourceNetworkingNetworkV2(),
			"openstack_networking_qos_bandwidth_limit_rule_v2":   dataSourceNetworkingQoSBandwidthLimitRuleV2(),
			"openstack_networking_qos_dscp_marking_rule_v2":      dataSourceNetworkingQoSDSCPMarkingRuleV2(),
			"openstack_networking_qos_minimum_bandwidth_rule_v2": dataSourceNetworkingQoSMinimumBandwidthRuleV2(),
			"openstack_networking_qos_policy_v2":                 dataSourceNetworkingQoSPolicyV2(),
			"openstack_networking_subnet_v2":                     dataSourceNetworkingSubnetV2(),
			"openstack_networking_secgroup_v2":                   dataSourceNetworkingSecGroupV2(),
			"openstack_networking_subnetpool_v2":                 dataSourceNetworkingSubnetPoolV2(),
			"openstack_networking_floatingip_v2":                 dataSourceNetworkingFloatingIPV2(),
			"openstack_networking_router_v2":                     dataSourceNetworkingRouterV2(),
			"openstack_networking_port_v2":                       dataSourceNetworkingPortV2(),
			"openstack_networking_port_ids_v2":                   dataSourceNetworkingPortIDsV2(),
			"openstack_networking_trunk_v2":                      dataSourceNetworkingTrunkV2(),
			"openstack_sharedfilesystem_availability_zones_v2":   dataSourceSharedFilesystemAvailabilityZonesV2(),
			"openstack_sharedfilesystem_sharenetwork_v2":         dataSourceSharedFilesystemShareNetworkV2(),
			"openstack_sharedfilesystem_share_v2":                dataSourceSharedFilesystemShareV2(),
			"openstack_sharedfilesystem_snapshot_v2":             dataSourceSharedFilesystemSnapshotV2(),
			"openstack_keymanager_secret_v1":                     dataSourceKeyManagerSecretV1(),
			"openstack_keymanager_container_v1":                  dataSourceKeyManagerContainerV1(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"openstack_blockstorage_quotaset_v2":                 resourceBlockStorageQuotasetV2(),
			"openstack_blockstorage_quotaset_v3":                 resourceBlockStorageQuotasetV3(),
			"openstack_blockstorage_volume_v1":                   resourceBlockStorageVolumeV1(),
			"openstack_blockstorage_volume_v2":                   resourceBlockStorageVolumeV2(),
			"openstack_blockstorage_volume_v3":                   resourceBlockStorageVolumeV3(),
			"openstack_blockstorage_volume_attach_v2":            resourceBlockStorageVolumeAttachV2(),
			"openstack_blockstorage_volume_attach_v3":            resourceBlockStorageVolumeAttachV3(),
			"openstack_compute_flavor_v2":                        resourceComputeFlavorV2(),
			"openstack_compute_flavor_access_v2":                 resourceComputeFlavorAccessV2(),
			"openstack_compute_instance_v2":                      resourceComputeInstanceV2(),
			"openstack_compute_interface_attach_v2":              resourceComputeInterfaceAttachV2(),
			"openstack_compute_keypair_v2":                       resourceComputeKeypairV2(),
			"openstack_compute_secgroup_v2":                      resourceComputeSecGroupV2(),
			"openstack_compute_servergroup_v2":                   resourceComputeServerGroupV2(),
			"openstack_compute_quotaset_v2":                      resourceComputeQuotasetV2(),
			"openstack_compute_floatingip_v2":                    resourceComputeFloatingIPV2(),
			"openstack_compute_floatingip_associate_v2":          resourceComputeFloatingIPAssociateV2(),
			"openstack_compute_volume_attach_v2":                 resourceComputeVolumeAttachV2(),
			"openstack_containerinfra_clustertemplate_v1":        resourceContainerInfraClusterTemplateV1(),
			"openstack_containerinfra_cluster_v1":                resourceContainerInfraClusterV1(),
			"openstack_db_instance_v1":                           resourceDatabaseInstanceV1(),
			"openstack_db_user_v1":                               resourceDatabaseUserV1(),
			"openstack_db_configuration_v1":                      resourceDatabaseConfigurationV1(),
			"openstack_db_database_v1":                           resourceDatabaseDatabaseV1(),
			"openstack_dns_recordset_v2":                         resourceDNSRecordSetV2(),
			"openstack_dns_zone_v2":                              resourceDNSZoneV2(),
			"openstack_fw_firewall_v1":                           resourceFWFirewallV1(),
			"openstack_fw_policy_v1":                             resourceFWPolicyV1(),
			"openstack_fw_rule_v1":                               resourceFWRuleV1(),
			"openstack_identity_endpoint_v3":                     resourceIdentityEndpointV3(),
			"openstack_identity_project_v3":                      resourceIdentityProjectV3(),
			"openstack_identity_role_v3":                         resourceIdentityRoleV3(),
			"openstack_identity_role_assignment_v3":              resourceIdentityRoleAssignmentV3(),
			"openstack_identity_service_v3":                      resourceIdentityServiceV3(),
			"openstack_identity_user_v3":                         resourceIdentityUserV3(),
			"openstack_identity_application_credential_v3":       resourceIdentityApplicationCredentialV3(),
			"openstack_images_image_v2":                          resourceImagesImageV2(),
			"openstack_images_image_access_v2":                   resourceImagesImageAccessV2(),
			"openstack_images_image_access_accept_v2":            resourceImagesImageAccessAcceptV2(),
			"openstack_lb_member_v1":                             resourceLBMemberV1(),
			"openstack_lb_monitor_v1":                            resourceLBMonitorV1(),
			"openstack_lb_pool_v1":                               resourceLBPoolV1(),
			"openstack_lb_vip_v1":                                resourceLBVipV1(),
			"openstack_lb_loadbalancer_v2":                       resourceLoadBalancerV2(),
			"openstack_lb_listener_v2":                           resourceListenerV2(),
			"openstack_lb_pool_v2":                               resourcePoolV2(),
			"openstack_lb_member_v2":                             resourceMemberV2(),
			"openstack_lb_monitor_v2":                            resourceMonitorV2(),
			"openstack_lb_l7policy_v2":                           resourceL7PolicyV2(),
			"openstack_lb_l7rule_v2":                             resourceL7RuleV2(),
			"openstack_networking_floatingip_v2":                 resourceNetworkingFloatingIPV2(),
			"openstack_networking_floatingip_associate_v2":       resourceNetworkingFloatingIPAssociateV2(),
			"openstack_networking_network_v2":                    resourceNetworkingNetworkV2(),
			"openstack_networking_port_v2":                       resourceNetworkingPortV2(),
			"openstack_networking_rbac_policy_v2":                resourceNetworkingRBACPolicyV2(),
			"openstack_networking_port_secgroup_associate_v2":    resourceNetworkingPortSecGroupAssociateV2(),
			"openstack_networking_qos_bandwidth_limit_rule_v2":   resourceNetworkingQoSBandwidthLimitRuleV2(),
			"openstack_networking_qos_dscp_marking_rule_v2":      resourceNetworkingQoSDSCPMarkingRuleV2(),
			"openstack_networking_qos_minimum_bandwidth_rule_v2": resourceNetworkingQoSMinimumBandwidthRuleV2(),
			"openstack_networking_qos_policy_v2":                 resourceNetworkingQoSPolicyV2(),
			"openstack_networking_quota_v2":                      resourceNetworkingQuotaV2(),
			"openstack_networking_router_v2":                     resourceNetworkingRouterV2(),
			"openstack_networking_router_interface_v2":           resourceNetworkingRouterInterfaceV2(),
			"openstack_networking_router_route_v2":               resourceNetworkingRouterRouteV2(),
			"openstack_networking_secgroup_v2":                   resourceNetworkingSecGroupV2(),
			"openstack_networking_secgroup_rule_v2":              resourceNetworkingSecGroupRuleV2(),
			"openstack_networking_subnet_v2":                     resourceNetworkingSubnetV2(),
			"openstack_networking_subnet_route_v2":               resourceNetworkingSubnetRouteV2(),
			"openstack_networking_subnetpool_v2":                 resourceNetworkingSubnetPoolV2(),
			"openstack_networking_addressscope_v2":               resourceNetworkingAddressScopeV2(),
			"openstack_networking_trunk_v2":                      resourceNetworkingTrunkV2(),
			"openstack_objectstorage_container_v1":               resourceObjectStorageContainerV1(),
			"openstack_objectstorage_object_v1":                  resourceObjectStorageObjectV1(),
			"openstack_objectstorage_tempurl_v1":                 resourceObjectstorageTempurlV1(),
			"openstack_orchestration_stack_v1":                   resourceOrchestrationStackV1(),
			"openstack_vpnaas_ipsec_policy_v2":                   resourceIPSecPolicyV2(),
			"openstack_vpnaas_service_v2":                        resourceServiceV2(),
			"openstack_vpnaas_ike_policy_v2":                     resourceIKEPolicyV2(),
			"openstack_vpnaas_endpoint_group_v2":                 resourceEndpointGroupV2(),
			"openstack_vpnaas_site_connection_v2":                resourceSiteConnectionV2(),
			"openstack_sharedfilesystem_securityservice_v2":      resourceSharedFilesystemSecurityServiceV2(),
			"openstack_sharedfilesystem_sharenetwork_v2":         resourceSharedFilesystemShareNetworkV2(),
			"openstack_sharedfilesystem_share_v2":                resourceSharedFilesystemShareV2(),
			"openstack_sharedfilesystem_share_access_v2":         resourceSharedFilesystemShareAccessV2(),
			"openstack_keymanager_secret_v1":                     resourceKeyManagerSecretV1(),
			"openstack_keymanager_container_v1":                  resourceKeyManagerContainerV1(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return configureProvider(d, terraformVersion)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"auth_url": "The Identity authentication URL.",

		"region": "The OpenStack region to connect to.",

		"user_name": "Username to login with.",

		"user_id": "User ID to login with.",

		"application_credential_id": "Application Credential ID to login with.",

		"application_credential_name": "Application Credential name to login with.",

		"application_credential_secret": "Application Credential secret to login with.",

		"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
			"to login with.",

		"password": "Password to login with.",

		"token": "Authentication token to use as an alternative to username/password.",

		"user_domain_name": "The name of the domain where the user resides (Identity v3).",

		"user_domain_id": "The ID of the domain where the user resides (Identity v3).",

		"project_domain_name": "The name of the domain where the project resides (Identity v3).",

		"project_domain_id": "The ID of the domain where the proejct resides (Identity v3).",

		"domain_id": "The ID of the Domain to scope to (Identity v3).",

		"domain_name": "The name of the Domain to scope to (Identity v3).",

		"default_domain": "The name of the Domain ID to scope to if no other domain is specified. Defaults to `default` (Identity v3).",

		"insecure": "Trust self-signed certificates.",

		"cacert_file": "A Custom CA certificate.",

		"endpoint_type": "The catalog endpoint type to use.",

		"cert": "A client certificate to authenticate with.",

		"key": "A client private key to authenticate with.",

		"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
			"interaction with Swift.",

		"use_octavia": "If set to `true`, API requests will go the Load Balancer\n" +
			"service (Octavia) instead of the Networking service (Neutron).",

		"delayed_auth": "If set to `true`, OpenStack authorization will be perfomed,\n" +
			"when the service provider client is called.",

		"allow_reauth": "If set to `true`, OpenStack authorization will be perfomed\n" +
			"automatically, if the initial auth token get expired. This is useful,\n" +
			"when the token TTL is low or the overall Terraform provider execution\n" +
			"time expected to be greater than the initial token TTL.",

		"cloud": "An entry in a `clouds.yaml` file to use.",

		"max_retries": "How many times HTTP connection should be retried until giving up.",

		"endpoint_overrides": "A map of services with an endpoint to override what was\n" +
			"from the Keystone catalog",

		"disable_no_cache_header": "If set to `true`, the HTTP `Cache-Control: no-cache` header will not be added by default to all API requests.",
	}
}

func configureProvider(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		auth.Config{
			CACertFile:                  d.Get("cacert_file").(string),
			ClientCertFile:              d.Get("cert").(string),
			ClientKeyFile:               d.Get("key").(string),
			Cloud:                       d.Get("cloud").(string),
			DefaultDomain:               d.Get("default_domain").(string),
			DomainID:                    d.Get("domain_id").(string),
			DomainName:                  d.Get("domain_name").(string),
			EndpointOverrides:           d.Get("endpoint_overrides").(map[string]interface{}),
			EndpointType:                d.Get("endpoint_type").(string),
			IdentityEndpoint:            d.Get("auth_url").(string),
			Password:                    d.Get("password").(string),
			ProjectDomainID:             d.Get("project_domain_id").(string),
			ProjectDomainName:           d.Get("project_domain_name").(string),
			Region:                      d.Get("region").(string),
			Swauth:                      d.Get("swauth").(bool),
			Token:                       d.Get("token").(string),
			TenantID:                    d.Get("tenant_id").(string),
			TenantName:                  d.Get("tenant_name").(string),
			UserDomainID:                d.Get("user_domain_id").(string),
			UserDomainName:              d.Get("user_domain_name").(string),
			Username:                    d.Get("user_name").(string),
			UserID:                      d.Get("user_id").(string),
			ApplicationCredentialID:     d.Get("application_credential_id").(string),
			ApplicationCredentialName:   d.Get("application_credential_name").(string),
			ApplicationCredentialSecret: d.Get("application_credential_secret").(string),
			UseOctavia:                  d.Get("use_octavia").(bool),
			DelayedAuth:                 d.Get("delayed_auth").(bool),
			AllowReauth:                 d.Get("allow_reauth").(bool),
			MaxRetries:                  d.Get("max_retries").(int),
			DisableNoCacheHeader:        d.Get("disable_no_cache_header").(bool),
			TerraformVersion:            terraformVersion,
			SDKVersion:                  meta.SDKVersionString(),
		},
	}

	v, ok := d.GetOkExists("insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
