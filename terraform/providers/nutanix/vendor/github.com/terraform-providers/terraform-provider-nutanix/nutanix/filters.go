package nutanix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

func DataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
func BuildFiltersDataSource(set *schema.Set) []*client.AdditionalFilter {
	filters := []*client.AdditionalFilter{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, &client.AdditionalFilter{
			Name:   m["name"].(string),
			Values: filterValues,
		})
	}
	return filters
}

func ReplaceFilterPrefixes(filters []*client.AdditionalFilter, mappings map[string]string) []*client.AdditionalFilter {
	if mappings == nil {
		return filters
	}

	for _, filter := range filters {
		filterPath := strings.Split(filter.Name, ".")
		fmt.Println(filterPath)
		if len(filterPath) > 0 {
			replacedBase, ok := mappings[filterPath[0]]
			fmt.Println(replacedBase)
			if ok {
				filterPath[0] = replacedBase
			}
		}
		filter.Name = strings.Join(filterPath, ".")
	}

	return filters
}
