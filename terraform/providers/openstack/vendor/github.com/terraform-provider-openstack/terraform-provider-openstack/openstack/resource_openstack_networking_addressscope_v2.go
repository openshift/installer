package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/addressscopes"
)

func resourceNetworkingAddressScopeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingAddressScopeV2Create,
		ReadContext:   resourceNetworkingAddressScopeV2Read,
		UpdateContext: resourceNetworkingAddressScopeV2Update,
		DeleteContext: resourceNetworkingAddressScopeV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
				ForceNew: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingAddressScopeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := addressscopes.CreateOpts{
		Name:      d.Get("name").(string),
		ProjectID: d.Get("project_id").(string),
		IPVersion: d.Get("ip_version").(int),
		Shared:    d.Get("shared").(bool),
	}

	log.Printf("[DEBUG] openstack_networking_addressscope_v2 create options: %#v", createOpts)
	a, err := addressscopes.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_addressscope_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_addressscope_v2 %s to become available", a.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    resourceNetworkingAddressScopeV2StateRefreshFunc(networkingClient, a.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_addressscope_v2 %s to become available: %s", a.ID, err)
	}

	d.SetId(a.ID)

	log.Printf("[DEBUG] Created openstack_networking_addressscope_v2 %s: %#v", a.ID, a)
	return resourceNetworkingAddressScopeV2Read(ctx, d, meta)
}

func resourceNetworkingAddressScopeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	a, err := addressscopes.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_addressscope_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_addressscope_v2 %s: %#v", d.Id(), a)

	d.Set("region", GetRegion(d, config))
	d.Set("name", a.Name)
	d.Set("project_id", a.ProjectID)
	d.Set("ip_version", a.IPVersion)
	d.Set("shared", a.Shared)

	return nil
}

func resourceNetworkingAddressScopeV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts addressscopes.UpdateOpts
	)

	if d.HasChange("name") {
		hasChange = true
		v := d.Get("name").(string)
		updateOpts.Name = &v
	}

	if d.HasChange("shared") {
		hasChange = true
		v := d.Get("shared").(bool)
		updateOpts.Shared = &v
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_addressscope_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = addressscopes.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_addressscope_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkingAddressScopeV2Read(ctx, d, meta)
}

func resourceNetworkingAddressScopeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := addressscopes.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_addressscope_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingAddressScopeV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_addressscope_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
