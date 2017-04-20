package local

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns the provider definition.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"local_file": resourceLocalFile(),
		},
	}
}
