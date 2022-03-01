package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
)

func dataSourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityGroupV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityGroupV3Read performs the group lookup.
func dataSourceIdentityGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	listOpts := groups.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] openstack_identity_group_v3 list options: %#v", listOpts)

	var group groups.Group
	allPages, err := groups.List(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_group_v3: %s", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_identity_group_v3: %s", err)
	}

	if len(allGroups) < 1 {
		return diag.Errorf("Your openstack_identity_group_v3 query returned no results. " +
			"Please change your search criteria and try again")
	}

	if len(allGroups) > 1 {
		return diag.Errorf("Your openstack_identity_group_v3 query returned more than one result")
	}

	group = allGroups[0]

	dataSourceIdentityGroupV3Attributes(d, config, &group)

	return nil
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityGroupV3Attributes(d *schema.ResourceData, config *Config, group *groups.Group) {
	log.Printf("[DEBUG] openstack_identity_group_v3 details: %#v", group)

	d.SetId(group.ID)
	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("domain_id", group.DomainID)
	d.Set("region", GetRegion(d, config))
}
