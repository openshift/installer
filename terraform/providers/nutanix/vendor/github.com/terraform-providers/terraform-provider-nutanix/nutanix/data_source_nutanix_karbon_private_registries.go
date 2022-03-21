package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixKarbonPrivateRegistries() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixKarbonPrivateRegistriesRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"private_registries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: KarbonPrivateRegistryElementDataSourceMap(),
				},
			},
		},
	}
}

func dataSourceNutanixKarbonPrivateRegistriesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	resp, err := conn.PrivateRegistry.ListKarbonPrivateRegistries()
	if err != nil {
		d.SetId("")
		return nil
	}

	privateRegistries := make([]map[string]interface{}, len(*resp))

	for k, v := range *resp {
		privateRegistry := make(map[string]interface{})
		if err != nil {
			return fmt.Errorf("error searching for private registry via legacy API: %s", err)
		}

		privateRegistry["name"] = utils.StringValue(v.Name)

		privateRegistry["endpoint"] = utils.StringValue(v.Endpoint)
		privateRegistry["uuid"] = utils.StringValue(v.UUID)
		privateRegistries[k] = privateRegistry
	}

	if err := d.Set("private_registries", privateRegistries); err != nil {
		return fmt.Errorf("failed to set private_registries output: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}
