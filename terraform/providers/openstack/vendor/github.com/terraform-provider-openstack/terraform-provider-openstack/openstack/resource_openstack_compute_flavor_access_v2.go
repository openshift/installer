package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

func resourceComputeFlavorAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFlavorAccessV2Create,
		ReadContext:   resourceComputeFlavorAccessV2Read,
		DeleteContext: resourceComputeFlavorAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeFlavorAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	flavorID := d.Get("flavor_id").(string)
	tenantID := d.Get("tenant_id").(string)

	accessOpts := flavors.AddAccessOpts{
		Tenant: tenantID,
	}
	log.Printf("[DEBUG] Flavor Access Options: %#v", accessOpts)

	if _, err := flavors.AddAccess(computeClient, flavorID, accessOpts).Extract(); err != nil {
		return diag.Errorf("Error adding access to tenant %s for flavor %s: %s", tenantID, flavorID, err)
	}

	id := fmt.Sprintf("%s/%s", flavorID, tenantID)
	d.SetId(id)

	return resourceComputeFlavorAccessV2Read(ctx, d, meta)
}

func resourceComputeFlavorAccessV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	flavorAccess, err := getFlavorAccess(computeClient, d)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting flavor access"))
	}

	d.Set("region", GetRegion(d, config))
	d.Set("flavor_id", flavorAccess.FlavorID)
	d.Set("tenant_id", flavorAccess.TenantID)

	return nil
}

func resourceComputeFlavorAccessV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	flavorAccess, err := getFlavorAccess(computeClient, d)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting flavor access"))
	}

	removeAccessOpts := flavors.RemoveAccessOpts{Tenant: flavorAccess.TenantID}
	log.Printf("[DEBUG] RemoveAccess Options: %#v", removeAccessOpts)

	if _, err := flavors.RemoveAccess(computeClient, flavorAccess.FlavorID, removeAccessOpts).Extract(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, fmt.Sprintf("Error removing tenant %s access from flavor %s", flavorAccess.TenantID, flavorAccess.FlavorID)))
	}

	return nil
}

func parseComputeFlavorAccessID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine flavor access ID")
	}

	flavorID := idParts[0]
	tenantID := idParts[1]

	return flavorID, tenantID, nil
}

func getFlavorAccess(computeClient *gophercloud.ServiceClient, d *schema.ResourceData) (flavors.FlavorAccess, error) {
	var access flavors.FlavorAccess
	flavorID, tenantID, err := parseComputeFlavorAccessID(d.Id())
	if err != nil {
		return access, err
	}

	found := false
	pager := flavors.ListAccesses(computeClient, flavorID)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		accessList, err := flavors.ExtractAccesses(page)
		if err != nil {
			return false, err
		}

		for _, a := range accessList {
			if a.TenantID == tenantID && a.FlavorID == flavorID {
				access = a
				found = true
				return false, nil
			}
		}

		return true, nil
	})

	if !found {
		return access, gophercloud.ErrDefault404{}
	}

	return access, err
}
