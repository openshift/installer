package packet

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a schema.Provider for managing Packet infrastructure.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PACKET_AUTH_TOKEN", nil),
				Description: "The API auth key for API operations.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"packet_precreated_ip_block": dataSourcePacketPreCreatedIPBlock(),
			"packet_operating_system":    dataSourceOperatingSystem(),
			"packet_spot_market_price":   dataSourceSpotMarketPrice(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"packet_device":               resourcePacketDevice(),
			"packet_ssh_key":              resourcePacketSSHKey(),
			"packet_project_ssh_key":      resourcePacketProjectSSHKey(),
			"packet_project":              resourcePacketProject(),
			"packet_organization":         resourcePacketOrganization(),
			"packet_volume":               resourcePacketVolume(),
			"packet_volume_attachment":    resourcePacketVolumeAttachment(),
			"packet_reserved_ip_block":    resourcePacketReservedIPBlock(),
			"packet_ip_attachment":        resourcePacketIPAttachment(),
			"packet_spot_market_request":  resourcePacketSpotMarketRequest(),
			"packet_vlan":                 resourcePacketVlan(),
			"packet_bgp_session":          resourcePacketBGPSession(),
			"packet_port_vlan_attachment": resourcePacketPortVlanAttachment(),
			"packet_connect":              resourcePacketConnect(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AuthToken: d.Get("auth_token").(string),
	}
	return config.Client(), nil
}

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
