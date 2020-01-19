// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceSearchSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"criteria": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"max": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"case_sensitive": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

// Helper method to support the customization of search field
func dataSourceSearchSchemaWith(vs ...string) *schema.Schema {
	schemaMap := make(map[string]*schema.Schema)
	for _, v := range vs {
		switch v {
		case "criteria":
			schemaMap[v] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
		case "max":
			schemaMap[v] = &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			}
		case "case_sensitive":
			schemaMap[v] = &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			}
		}
	}
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: schemaMap,
		},
	}
}
