package nutanix

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func categoriesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Set: func(v interface{}) int {
			category := v.(map[string]interface{})
			return hashcode.String(category["name"].(string) + category["value"].(string))
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func expandCategories(categoriesSet interface{}) map[string]string {
	categories := categoriesSet.(*schema.Set).List()
	output := make(map[string]string)

	for _, v := range categories {
		category := v.(map[string]interface{})
		output[category["name"].(string)] = category["value"].(string)
	}

	return output
}

func flattenCategories(categories map[string]string) []interface{} {
	c := make([]interface{}, 0)

	for name, value := range categories {
		c = append(c, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}

	return c
}

func categoriesSchemaOptional() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Set: func(v interface{}) int {
			category := v.(map[string]interface{})
			return hashcode.String(category["name"].(string) + category["value"].(string))
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}
