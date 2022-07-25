package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func resourceComputeFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFloatingIPAssociateV2Create,
		ReadContext:   resourceComputeFloatingIPAssociateV2Read,
		DeleteContext: resourceComputeFloatingIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floating_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"wait_until_associated": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeFloatingIPAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	fixedIP := d.Get("fixed_ip").(string)
	instanceID := d.Get("instance_id").(string)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP,
		FixedIP:    fixedIP,
	}
	log.Printf("[DEBUG] openstack_compute_floatingip_associate_v2 create options: %#v", associateOpts)

	err = floatingips.AssociateInstance(computeClient, instanceID, associateOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_compute_floatingip_associate_v2: %s", err)
	}

	// This API call should be synchronous, but we've had reports where it isn't.
	// If the user opted in to wait for association, then poll here.
	var waitUntilAssociated bool
	if v, ok := d.GetOkExists("wait_until_associated"); ok {
		if wua, ok := v.(bool); ok {
			waitUntilAssociated = wua
		}
	}

	if waitUntilAssociated {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"NOT_ASSOCIATED"},
			Target:     []string{"ASSOCIATED"},
			Refresh:    computeFloatingIPAssociateV2CheckAssociation(computeClient, instanceID, floatingIP),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// There's an API call to get this information, but it has been
	// deprecated. The Neutron API could be used, but I'm trying not
	// to mix service APIs. Therefore, a faux ID will be used.
	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceID, fixedIP)
	d.SetId(id)

	return resourceComputeFloatingIPAssociateV2Read(ctx, d, meta)
}

func resourceComputeFloatingIPAssociateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	floatingIP, instanceID, fixedIP, err := parseComputeFloatingIPAssociateID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Now check and see whether the floating IP still exists.
	// First try to do this by querying the Network API.
	networkEnabled := true
	networkClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		networkEnabled = false
	}

	var exists bool
	if networkEnabled {
		log.Printf("[DEBUG] Checking for openstack_compute_floatingip_associate_v2 %s existence via Network API", d.Id())
		exists, err = computeFloatingIPAssociateV2NetworkExists(networkClient, floatingIP)
	} else {
		log.Printf("[DEBUG] Checking for openstack_compute_floatingip_associate_v2 %s existence via Compute API", d.Id())
		exists, err = computeFloatingIPAssociateV2ComputeExists(computeClient, floatingIP)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		d.SetId("")
	}

	// Next, see if the instance still exists
	instance, err := servers.Get(computeClient, instanceID).Extract()
	if err != nil {
		if CheckDeleted(d, err, "instance") == nil {
			return nil
		}
	}

	// Finally, check and see if the floating ip is still associated with the instance.
	var associated bool
	for _, networkAddresses := range instance.Addresses {
		for _, element := range networkAddresses.([]interface{}) {
			address := element.(map[string]interface{})
			if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == floatingIP {
				associated = true
			}
		}
	}

	if !associated {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	d.Set("floating_ip", floatingIP)
	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", fixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeFloatingIPAssociateV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	instanceID := d.Get("instance_id").(string)

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: floatingIP,
	}
	log.Printf("[DEBUG] openstack_compute_floatingip_associate_v2 %s delete options: %#v", d.Id(), disassociateOpts)

	err = floatingips.DisassociateInstance(computeClient, instanceID, disassociateOpts).ExtractErr()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault409); ok {
			// 409 is returned when floating ip address is not associated with an instance.
			log.Printf("[DEBUG] openstack_compute_floatingip_associate_v2 %s is not associated with instance %s", d.Id(), instanceID)
		} else {
			return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_compute_floatingip_associate_v2"))
		}
	}

	return nil
}
