package tftime

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"time_offset":   resourceTimeOffset(),
			"time_rotating": resourceTimeRotating(),
			"time_sleep":    resourceTimeSleep(),
			"time_static":   resourceTimeStatic(),
		},
	}
}
