package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
)

func resourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupV3Create,
		ReadContext:   resourceIdentityGroupV3Read,
		UpdateContext: resourceIdentityGroupV3Update,
		DeleteContext: resourceIdentityGroupV3Delete,
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

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIdentityGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := groups.CreateOpts{
		DomainID:    d.Get("domain_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] openstack_identity_group_v3 create options: %#v", createOpts)
	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_group_v3: %s", err)
	}

	d.SetId(group.ID)

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_group_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_group_v3: %#v", group)

	d.Set("domain_id", group.DomainID)
	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if hasChange {
		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_group_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_group_v3"))
	}

	return nil
}
