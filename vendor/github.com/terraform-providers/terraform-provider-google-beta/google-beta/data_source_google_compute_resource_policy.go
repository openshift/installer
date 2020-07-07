package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleComputeResourcePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeResourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeResourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resourcePolicy, err := config.clientCompute.ResourcePolicies.Get(project, region, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ResourcePolicy Not Found : %s", name))
	}
	d.Set("self_link", resourcePolicy.SelfLink)
	d.Set("description", resourcePolicy.Description)
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/resourcePolicies/%s", project, region, name))
	return nil
}
