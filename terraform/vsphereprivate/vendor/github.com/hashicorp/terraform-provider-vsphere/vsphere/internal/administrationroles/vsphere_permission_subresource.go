package administrationroles

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func VspherePermissionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_or_group": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User or group receiving access.",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if strings.ToLower(old) == strings.ToLower(new) {
					return true
				}
				return false
			},
		},
		"propagate": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Whether or not this permission propagates down the hierarchy to sub-entities.",
		},
		"is_group": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Whether user_or_group field refers to a user or a group. True for a group and false for a user.",
		},
		"role_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Reference to the role providing the access.",
		},
	}
}
