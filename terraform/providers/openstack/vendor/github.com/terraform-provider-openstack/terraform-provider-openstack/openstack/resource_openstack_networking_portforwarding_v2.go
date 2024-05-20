package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/portforwarding"
)

func resourceNetworkingPortForwardingV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPortForwardingV2Create,
		ReadContext:   resourceNetworkPortForwardingV2Read,
		UpdateContext: resourceNetworkPortForwardingV2Update,
		DeleteContext: resourceNetworkPortForwardingV2Delete,

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

			"floatingip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"internal_port_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"external_port": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"protocol": {
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

func resourceNetworkPortForwardingV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)
	createOpts := portforwarding.CreateOpts{
		InternalIPAddress: d.Get("internal_ip_address").(string),
		ExternalPort:      d.Get("external_port").(int),
		InternalPort:      d.Get("internal_port").(int),
		InternalPortID:    d.Get("internal_port_id").(string),
		Protocol:          d.Get("protocol").(string),
	}

	// TODO: add description.

	log.Printf("[DEBUG] openstack_networking_portforwarding_v2 create options: %#v", createOpts)

	pf, err := portforwarding.Create(networkingClient, fipID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_portforwarding_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_portforwarding_v2 %s to become available.", pf.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingPortForwardingV2StateRefreshFunc(networkingClient, fipID, pf.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_portforwarding_v2 %s to become available: %s", pf.ID, err)
	}

	d.SetId(pf.ID)

	log.Printf("[DEBUG] Created openstack_networking_portforwarding_v2 %s: %#v", pf.ID, pf)
	return resourceNetworkPortForwardingV2Read(ctx, d, meta)
}

func resourceNetworkPortForwardingV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)

	pf, err := portforwarding.Get(networkingClient, fipID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_portforwarding_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_portforwarding_v2 %s: %#v", d.Id(), pf)

	d.Set("id", pf.ID)
	d.Set("internal_port_id", pf.InternalPortID)
	d.Set("internal_ip_address", pf.InternalIPAddress)
	d.Set("internal_port", pf.InternalPort)
	d.Set("external_port", pf.ExternalPort)
	d.Set("protocol", pf.Protocol)
	d.Set("region", GetRegion(d, config))

	// TODO: add description.

	return nil
}

func resourceNetworkPortForwardingV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var hasChange bool
	var updateOpts portforwarding.UpdateOpts

	fipID := d.Get("floating_IP_ID").(string)

	if d.HasChange("internal_port_id") {
		hasChange = true
		internalPortID := d.Get("internal_port_id").(string)
		updateOpts.InternalPortID = internalPortID
	}

	if d.HasChange("external_port") {
		hasChange = true
		externalPort := d.Get("external_port").(int)
		updateOpts.ExternalPort = externalPort
	}

	if d.HasChange("internal_port") {
		hasChange = true
		internalPort := d.Get("internal_port").(int)
		updateOpts.InternalPort = internalPort
	}
	if d.HasChange("protocol") {
		hasChange = true
		protocol := d.Get("protocol").(string)
		updateOpts.Protocol = protocol
	}

	// TODO: add description.

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_portforwarding_v2 %s update options: %#v", d.Id(), updateOpts)
		_, err = portforwarding.Update(networkingClient, fipID, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_portforwarding_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkPortForwardingV2Read(ctx, d, meta)
}

func resourceNetworkPortForwardingV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)

	if err := portforwarding.Delete(networkingClient, fipID, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_portforwarding_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingPortForwardingV2StateRefreshFunc(networkingClient, fipID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_portforwarding_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
