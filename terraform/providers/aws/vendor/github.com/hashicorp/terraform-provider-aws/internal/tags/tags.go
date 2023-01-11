package tags

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TagsSchema returns the schema to use for tags.
func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func TagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func TagsSchemaComputedDeprecated(message string) *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeMap,
		Computed:   true,
		Elem:       &schema.Schema{Type: schema.TypeString},
		Deprecated: message,
	}
}

func TagsSchemaForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}
