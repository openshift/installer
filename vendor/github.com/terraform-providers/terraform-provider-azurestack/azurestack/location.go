package azurestack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func locationSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		StateFunc:        azureStackNormalizeLocation,
		DiffSuppressFunc: azureStackSuppressLocationDiff,
	}
}

func locationForDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

// azureStackNormalizeLocation is a function which normalises human-readable region/location
// names (e.g. "West US") to the values used and returned by the Azure API (e.g. "westus").
// In state we track the API internal version as it is easier to go from the human form
// to the canonical form than the other way around.
func azureStackNormalizeLocation(location interface{}) string {
	input := location.(string)
	return strings.Replace(strings.ToLower(input), " ", "", -1)
}

func azureStackSuppressLocationDiff(k, old, new string, d *schema.ResourceData) bool {
	return azureStackNormalizeLocation(old) == azureStackNormalizeLocation(new)
}
