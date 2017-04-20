package matchbox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns the provider schema to Terraform.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// Provider configuration
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_cert": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ca": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"matchbox_profile": resourceProfile(),
			"matchbox_group":   resourceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	ca := d.Get("ca").(string)
	clientCert := d.Get("client_cert").(string)
	clientKey := d.Get("client_key").(string)

	config := &Config{
		Endpoint:   d.Get("endpoint").(string),
		ClientCert: []byte(clientCert),
		ClientKey:  []byte(clientKey),
		CA:         []byte(ca),
	}

	return NewMatchboxClient(config)
}
