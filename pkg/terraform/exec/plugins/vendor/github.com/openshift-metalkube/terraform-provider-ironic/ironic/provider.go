package ironic

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	noauthintrospection "github.com/gophercloud/gophercloud/openstack/baremetalintrospection/noauth"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

// Clients stores the client connection information for Ironic and Inspector
type Clients struct {
	Ironic    *gophercloud.ServiceClient
	Inspector *gophercloud.ServiceClient
}

// Provider Ironic
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_ENDPOINT", ""),
				Description: descriptions["url"],
			},
			"inspector": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_INSPECTOR_ENDPOINT", ""),
				Description: descriptions["inspector"],
			},
			"microversion": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("IRONIC_MICROVERSION", "1.52"),
				Description: descriptions["microversion"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ironic_node_v1":       resourceNodeV1(),
			"ironic_port_v1":       resourcePortV1(),
			"ironic_allocation_v1": resourceAllocationV1(),
			"ironic_deployment":    resourceDeployment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ironic_introspection": dataSourceIronicIntrospection(),
		},
		ConfigureFunc: configureProvider,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"url":          "The authentication endpoint for Ironic",
		"inspector":    "The endpoint for Ironic inspector",
		"microversion": "The microversion to use for Ironic",
	}
}

// Creates a noauth Ironic client
func configureProvider(schema *schema.ResourceData) (interface{}, error) {
	var clients Clients

	url := schema.Get("url").(string)
	if url == "" {
		return nil, fmt.Errorf("url is required for ironic provider")
	}
	log.Printf("[DEBUG] Ironic endpoint is %s", url)

	ironic, err := noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: url,
	})
	if err != nil {
		return nil, err
	}
	ironic.Microversion = schema.Get("microversion").(string)
	clients.Ironic = ironic

	inspectorURL := schema.Get("inspector").(string)
	if inspectorURL != "" {
		log.Printf("[DEBUG] Inspector endpoint is %s", inspectorURL)
		inspector, err := noauthintrospection.NewBareMetalIntrospectionNoAuth(noauthintrospection.EndpointOpts{
			IronicInspectorEndpoint: inspectorURL,
		})
		if err != nil {
			return nil, fmt.Errorf("could not configure inspector endpoint: %s", err.Error())
		}
		clients.Inspector = inspector
	}

	return clients, err
}
