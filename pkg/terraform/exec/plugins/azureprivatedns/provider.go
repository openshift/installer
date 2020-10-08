package azureprivatedns

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"azureprivatedns_zone":                      resourceArmPrivateDNSZone(),
			"azureprivatedns_a_record":                  resourceArmPrivateDNSARecord(),
			"azureprivatedns_aaaa_record":               resourceArmPrivateDNSAAAARecord(),
			"azureprivatedns_srv_record":                resourceArmPrivateDNSSrvRecord(),
			"azureprivatedns_zone_virtual_network_link": resourceArmPrivateDNSZoneVirtualNetworkLink(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		builder := &authentication.Builder{
			SubscriptionID: d.Get("subscription_id").(string),
			ClientID:       d.Get("client_id").(string),
			ClientSecret:   d.Get("client_secret").(string),
			TenantID:       d.Get("tenant_id").(string),
			Environment:    d.Get("environment").(string),

			// Feature Toggles
			SupportsClientSecretAuth: true,
		}

		config, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("error building AzureRM Client: %s", err)
		}

		client, err := getArmClient(config)

		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error {
			client.StopContext = p.StopContext()
			return nil
		}

		return client, nil
	}
}
