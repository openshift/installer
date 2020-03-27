package azureprivatedns

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
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

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = mutexkv.NewMutexKV()

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
