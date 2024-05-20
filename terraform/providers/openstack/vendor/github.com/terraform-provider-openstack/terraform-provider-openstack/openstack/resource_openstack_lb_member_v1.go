package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas/members"
)

func resourceLBMemberV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBMemberV1Create,
		ReadContext:   resourceLBMemberV1Read,
		UpdateContext: resourceLBMemberV1Update,
		DeleteContext: resourceLBMemberV1Delete,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceLBMemberV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := members.CreateOpts{
		TenantID:     d.Get("tenant_id").(string),
		PoolID:       d.Get("pool_id").(string),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("port").(int),
	}

	log.Printf("[DEBUG] OpenStack LB Member Create Options: %#v", createOpts)
	m, err := members.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating OpenStack LB member: %s", err)
	}
	log.Printf("[INFO] LB member ID: %s", m.ID)

	log.Printf("[DEBUG] Waiting for OpenStack LB member (%s) to become available.", m.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE", "INACTIVE", "CREATED", "DOWN"},
		Refresh:    waitForLBMemberActive(networkingClient, m.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(m.ID)

	// Due to the way Gophercloud is currently set up, AdminStateUp must be set post-create
	asu := d.Get("admin_state_up").(bool)
	updateOpts := members.UpdateOpts{
		AdminStateUp: &asu,
	}

	log.Printf("[DEBUG] OpenStack LB Member Update Options: %#v", createOpts)
	_, err = members.Update(networkingClient, m.ID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating OpenStack LB member: %s", err)
	}

	return resourceLBMemberV1Read(ctx, d, meta)
}

func resourceLBMemberV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	m, err := members.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "LB member"))
	}

	log.Printf("[DEBUG] Retrieved OpenStack LB member %s: %+v", d.Id(), m)

	d.Set("address", m.Address)
	d.Set("pool_id", m.PoolID)
	d.Set("port", m.ProtocolPort)
	d.Set("weight", m.Weight)
	d.Set("admin_state_up", m.AdminStateUp)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceLBMemberV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts members.UpdateOpts
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	log.Printf("[DEBUG] Updating LB member %s with options: %+v", d.Id(), updateOpts)

	_, err = members.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating OpenStack LB member: %s", err)
	}

	return resourceLBMemberV1Read(ctx, d, meta)
}

func resourceLBMemberV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = members.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		if err := CheckDeleted(d, err, "LB member"); err != nil {
			log.Printf("%s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "PENDING_DELETE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForLBMemberDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting OpenStack LB member: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForLBMemberActive(networkingClient *gophercloud.ServiceClient, memberID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		m, err := members.Get(networkingClient, memberID).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenStack LB member: %+v", m)
		if m.Status == "ACTIVE" {
			return m, "ACTIVE", nil
		}

		return m, m.Status, nil
	}
}

func waitForLBMemberDelete(networkingClient *gophercloud.ServiceClient, memberID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenStack LB member %s", memberID)

		m, err := members.Get(networkingClient, memberID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenStack LB member %s", memberID)
				return m, "DELETED", nil
			}
			return m, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenStack LB member %s still active.", memberID)
		return m, "ACTIVE", nil
	}
}
