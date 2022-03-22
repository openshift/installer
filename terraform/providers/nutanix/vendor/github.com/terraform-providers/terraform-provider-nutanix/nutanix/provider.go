package nutanix

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() terraform.ResourceProvider {
	// defines descriptions for ResourceProvider schema definitions
	descriptions := map[string]string{
		"username": "User name for Nutanix Prism. Could be\n" +
			"local cluster auth (e.g. 'admin') or directory auth.",

		"password": "Password for provided user name.",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"session_auth": "Use session authentification instead of basic auth for each request",

		"port": "Port for Nutanix Prism.",

		"wait_timeout": "Set if you know that the creation o update of a resource may take long time (minutes)",

		"endpoint": "URL for Nutanix Prism (e.g IP or FQDN for cluster VIP\n" +
			"note, this is never the data services VIP, and should not be an\n" +
			"individual CVM address, as this would cause calls to fail during\n" +
			"cluster lifecycle management operations, such as AOS upgrades.",
	}

	// Nutanix provider schema
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_USERNAME", nil),
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PASSWORD", nil),
				Description: descriptions["password"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_INSECURE", false),
				Description: descriptions["insecure"],
			},
			"session_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_SESSION_AUTH", false),
				Description: descriptions["session_auth"],
			},
			"port": {
				Type:        schema.TypeString,
				Default:     "9440",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PORT", false),
				Description: descriptions["port"],
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_ENDPOINT", nil),
				Description: descriptions["endpoint"],
			},
			"wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_WAIT_TIMEOUT", nil),
				Description: descriptions["wait_timeout"],
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PROXY_URL", nil),
				Description: descriptions["proxy_url"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nutanix_image":                     dataSourceNutanixImage(),
			"nutanix_subnet":                    dataSourceNutanixSubnet(),
			"nutanix_subnets":                   dataSourceNutanixSubnets(),
			"nutanix_cluster":                   dataSourceNutanixCluster(),
			"nutanix_clusters":                  dataSourceNutanixClusters(),
			"nutanix_virtual_machine":           dataSourceNutanixVirtualMachine(),
			"nutanix_category_key":              dataSourceNutanixCategoryKey(),
			"nutanix_network_security_rule":     dataSourceNutanixNetworkSecurityRule(),
			"nutanix_host":                      dataSourceNutanixHost(),
			"nutanix_hosts":                     dataSourceNutanixHosts(),
			"nutanix_access_control_policy":     dataSourceNutanixAccessControlPolicy(),
			"nutanix_access_control_policies":   dataSourceNutanixAccessControlPolicies(),
			"nutanix_project":                   dataSourceNutanixProject(),
			"nutanix_projects":                  dataSourceNutanixProjects(),
			"nutanix_role":                      dataSourceNutanixRole(),
			"nutanix_roles":                     dataSourceNutanixRoles(),
			"nutanix_user":                      dataSourceNutanixUser(),
			"nutanix_users":                     dataSourceNutanixUsers(),
			"nutanix_user_group":                dataSourceNutanixUserGroup(),
			"nutanix_user_groups":               dataSourceNutanixUserGroups(),
			"nutanix_permission":                dataSourceNutanixPermission(),
			"nutanix_permissions":               dataSourceNutanixPermissions(),
			"nutanix_karbon_cluster_kubeconfig": dataSourceNutanixKarbonClusterKubeconfig(),
			"nutanix_karbon_cluster":            dataSourceNutanixKarbonCluster(),
			"nutanix_karbon_clusters":           dataSourceNutanixKarbonClusters(),
			"nutanix_karbon_cluster_ssh":        dataSourceNutanixKarbonClusterSSH(),
			"nutanix_karbon_private_registry":   dataSourceNutanixKarbonPrivateRegistry(),
			"nutanix_karbon_private_registries": dataSourceNutanixKarbonPrivateRegistries(),
			"nutanix_protection_rule":           dataSourceNutanixProtectionRule(),
			"nutanix_protection_rules":          dataSourceNutanixProtectionRules(),
			"nutanix_recovery_plan":             dataSourceNutanixRecoveryPlan(),
			"nutanix_recovery_plans":            dataSourceNutanixRecoveryPlans(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nutanix_virtual_machine":         resourceNutanixVirtualMachine(),
			"nutanix_image":                   resourceNutanixImage(),
			"nutanix_subnet":                  resourceNutanixSubnet(),
			"nutanix_category_key":            resourceNutanixCategoryKey(),
			"nutanix_category_value":          resourceNutanixCategoryValue(),
			"nutanix_network_security_rule":   resourceNutanixNetworkSecurityRule(),
			"nutanix_access_control_policy":   resourceNutanixAccessControlPolicy(),
			"nutanix_project":                 resourceNutanixProject(),
			"nutanix_role":                    resourceNutanixRole(),
			"nutanix_user":                    resourceNutanixUser(),
			"nutanix_karbon_cluster":          resourceNutanixKarbonCluster(),
			"nutanix_karbon_private_registry": resourceNutanixKarbonPrivateRegistry(),
			"nutanix_protection_rule":         resourceNutanixProtectionRule(),
			"nutanix_recovery_plan":           resourceNutanixRecoveryPlan(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// This function used to fetch the configuration params given to our provider which
// we will use to initialize a dummy client that interacts with API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[DEBUG] config wait_timeout %d", d.Get("wait_timeout").(int))

	config := Config{
		Endpoint:    d.Get("endpoint").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		Insecure:    d.Get("insecure").(bool),
		SessionAuth: d.Get("session_auth").(bool),
		Port:        d.Get("port").(string),
		WaitTimeout: int64(d.Get("wait_timeout").(int)),
		ProxyURL:    d.Get("proxy_url").(string),
	}

	return config.Client()
}
