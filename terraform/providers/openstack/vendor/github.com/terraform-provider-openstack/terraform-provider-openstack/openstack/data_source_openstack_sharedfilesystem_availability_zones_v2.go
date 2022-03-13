package openstack

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/availabilityzones"
	"github.com/gophercloud/utils/terraform/hashcode"
)

func dataSourceSharedFilesystemAvailabilityZonesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemAvailabilityZonesV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSharedFilesystemAvailabilityZonesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	allPages, err := availabilityzones.List(client).AllPages()
	if err != nil {
		return diag.Errorf("Error retrieving openstack_sharedfilesystem_availability_zones_v2: %s", err)
	}
	zoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return diag.Errorf("Error extracting openstack_sharedfilesystem_availability_zones_v2 from response: %s", err)
	}

	zones := make([]string, 0, len(zoneInfo))
	for _, z := range zoneInfo {
		zones = append(zones, z.Name)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(zones)

	d.SetId(hashcode.Strings(zones))
	d.Set("names", zones)
	d.Set("region", GetRegion(d, config))

	return nil
}
