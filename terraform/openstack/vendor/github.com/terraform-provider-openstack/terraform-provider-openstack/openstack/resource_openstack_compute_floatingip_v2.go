package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
)

func resourceComputeFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFloatingIPV2Create,
		ReadContext:   resourceComputeFloatingIPV2Read,
		DeleteContext: resourceComputeFloatingIPV2Delete,
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

			"pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_POOL_NAME", nil),
			},

			// computed-only
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeFloatingIPV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	createOpts := &floatingips.CreateOpts{
		Pool: d.Get("pool").(string),
	}

	log.Printf("[DEBUG] openstack_compute_floatingip_v2 Create Options: %#v", createOpts)

	newFip, err := floatingips.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_compute_floatingip_v2: %s", err)
	}

	d.SetId(newFip.ID)

	return resourceComputeFloatingIPV2Read(ctx, d, meta)
}

func resourceComputeFloatingIPV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	fip, err := floatingips.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_floatingip_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_floatingip_v2 %s: %#v", d.Id(), fip)

	d.Set("pool", fip.Pool)
	d.Set("instance_id", fip.InstanceID)
	d.Set("address", fip.IP)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeFloatingIPV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	if err := floatingips.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_compute_floatingip_v2"))
	}

	return nil
}
