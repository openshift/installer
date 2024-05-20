package openstack

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
)

func dataSourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkingSecGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := groups.ListOpts{
		ID:          d.Get("secgroup_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	allSecGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve security groups: %s", err)
	}

	if len(allSecGroups) < 1 {
		return diag.Errorf("No Security Group found with name: %s", d.Get("name"))
	}

	if len(allSecGroups) > 1 {
		return diag.Errorf("More than one Security Group found with name: %s", d.Get("name"))
	}

	secGroup := allSecGroups[0]

	log.Printf("[DEBUG] Retrieved Security Group %s: %+v", secGroup.ID, secGroup)
	d.SetId(secGroup.ID)

	d.Set("name", secGroup.Name)
	d.Set("description", secGroup.Description)
	d.Set("tenant_id", secGroup.TenantID)
	d.Set("all_tags", secGroup.Tags)
	d.Set("region", GetRegion(d, config))

	return nil
}
