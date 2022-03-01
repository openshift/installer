package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
)

func resourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleV3Create,
		ReadContext:   resourceIdentityRoleV3Read,
		UpdateContext: resourceIdentityRoleV3Update,
		DeleteContext: resourceIdentityRoleV3Delete,
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
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIdentityRoleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := roles.CreateOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] openstack_identity_role_v3 create options: %#v", createOpts)
	role, err := roles.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_role_v3: %s", err)
	}

	d.SetId(role.ID)

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	role, err := roles.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_role_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_role_v3: %#v", role)

	d.Set("domain_id", role.DomainID)
	d.Set("name", role.Name)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityRoleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts roles.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		_, err := roles.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_role_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = roles.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_role_v3"))
	}

	return nil
}
