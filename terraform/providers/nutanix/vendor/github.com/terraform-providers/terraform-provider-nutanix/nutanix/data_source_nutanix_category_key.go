package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixCategoryKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixCategoryKeyRead,

		Schema: map[string]*schema.Schema{
			"system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNutanixCategoryKeyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CategoryKey: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryKey(d.Get("name").(string))

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}

		return err
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Name))
	d.Set("description", utils.StringValue(resp.Description))
	d.Set("system_defined", utils.BoolValue(resp.SystemDefined))

	d.SetId(utils.StringValue(resp.Name))

	list, err := conn.V3.ListAllCategoryValues(d.Get("name").(string), "")

	if err != nil {
		return err
	}

	values := make([]string, len(list.Entities))

	for k, v := range list.Entities {
		values[k] = utils.StringValue(v.Value)
	}

	return d.Set("values", values)
}
