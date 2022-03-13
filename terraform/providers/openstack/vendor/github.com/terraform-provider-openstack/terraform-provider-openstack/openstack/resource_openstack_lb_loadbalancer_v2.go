package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	neutronloadbalancers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
)

func resourceLoadBalancerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV2Create,
		ReadContext:   resourceLoadBalancerV2Read,
		UpdateContext: resourceLoadBalancerV2Update,
		DeleteContext: resourceLoadBalancerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vip_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_port_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceLoadBalancerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		lbID       string
		vipPortID  string
		lbProvider string
	)

	if v, ok := d.GetOk("loadbalancer_provider"); ok {
		lbProvider = v.(string)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	if lbClient.Type == octaviaLBClientType {
		createOpts := octavialoadbalancers.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			VipNetworkID: d.Get("vip_network_id").(string),
			VipSubnetID:  d.Get("vip_subnet_id").(string),
			VipPortID:    d.Get("vip_port_id").(string),
			ProjectID:    d.Get("tenant_id").(string),
			VipAddress:   d.Get("vip_address").(string),
			AdminStateUp: &adminStateUp,
			FlavorID:     d.Get("flavor_id").(string),
			Provider:     lbProvider,
		}

		// availability_zone requires octavia minor version 2.14. Only set when specified.
		if v, ok := d.GetOk("availability_zone"); ok {
			aZ := v.(string)
			createOpts.AvailabilityZone = aZ
		}

		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			createOpts.Tags = expandToStringSlice(tags)
		}

		log.Printf("[DEBUG][Octavia] openstack_lb_loadbalancer_v2 create options: %#v", createOpts)
		lb, err := octavialoadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating openstack_lb_loadbalancer_v2: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	} else {
		createOpts := neutronloadbalancers.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			VipSubnetID:  d.Get("vip_subnet_id").(string),
			TenantID:     d.Get("tenant_id").(string),
			VipAddress:   d.Get("vip_address").(string),
			AdminStateUp: &adminStateUp,
			FlavorID:     d.Get("flavor_id").(string),
			Provider:     lbProvider,
		}

		log.Printf("[DEBUG][Neutron] openstack_lb_loadbalancer_v2 create options: %#v", createOpts)
		lb, err := neutronloadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating openstack_lb_loadbalancer_v2: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	}

	// Wait for load-balancer to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Once the load-balancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}
	if err := resourceLoadBalancerV2SetSecurityGroups(networkingClient, vipPortID, d); err != nil {
		return diag.Errorf("Error setting openstack_lb_loadbalancer_v2 security groups: %s", err)
	}

	d.SetId(lbID)

	return resourceLoadBalancerV2Read(ctx, d, meta)
}

func resourceLoadBalancerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var vipPortID string

	if lbClient.Type == octaviaLBClientType {
		lb, err := octavialoadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve openstack_lb_loadbalancer_v2"))
		}

		log.Printf("[DEBUG][Octavia] Retrieved openstack_lb_loadbalancer_v2 %s: %#v", d.Id(), lb)

		d.Set("name", lb.Name)
		d.Set("description", lb.Description)
		d.Set("vip_subnet_id", lb.VipSubnetID)
		d.Set("vip_network_id", lb.VipNetworkID)
		d.Set("tenant_id", lb.ProjectID)
		d.Set("vip_address", lb.VipAddress)
		d.Set("vip_port_id", lb.VipPortID)
		d.Set("admin_state_up", lb.AdminStateUp)
		d.Set("flavor_id", lb.FlavorID)
		d.Set("loadbalancer_provider", lb.Provider)
		d.Set("availability_zone", lb.AvailabilityZone)
		d.Set("region", GetRegion(d, config))
		d.Set("tags", lb.Tags)
		vipPortID = lb.VipPortID
	} else {
		lb, err := neutronloadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Unable to retrieve openstack_lb_loadbalancer_v2"))
		}

		log.Printf("[DEBUG][Neutron] Retrieved openstack_lb_loadbalancer_v2 %s: %#v", d.Id(), lb)

		d.Set("name", lb.Name)
		d.Set("description", lb.Description)
		d.Set("vip_subnet_id", lb.VipSubnetID)
		d.Set("tenant_id", lb.TenantID)
		d.Set("vip_address", lb.VipAddress)
		d.Set("vip_port_id", lb.VipPortID)
		d.Set("admin_state_up", lb.AdminStateUp)
		d.Set("flavor_id", lb.FlavorID)
		d.Set("loadbalancer_provider", lb.Provider)
		d.Set("region", GetRegion(d, config))
		vipPortID = lb.VipPortID
	}

	// Get any security groups on the VIP Port.
	if vipPortID != "" {
		networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating OpenStack networking client: %s", err)
		}
		if err := resourceLoadBalancerV2GetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error getting port security groups for openstack_lb_loadbalancer_v2: %s", err)
		}
	}

	return nil
}

func resourceLoadBalancerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	updateOpts, err := chooseLBV2LoadbalancerUpdateOpts(d, config)
	if err != nil {
		return diag.Errorf("Error building openstack_lb_loadbalancer_v2 update options: %s", err)
	}

	if updateOpts != nil {
		// Wait for load-balancer to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBV2LoadBalancer(ctx, lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating openstack_lb_loadbalancer_v2 %s with options: %#v", d.Id(), updateOpts)
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err = neutronloadbalancers.Update(lbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.Errorf("Error updating openstack_lb_loadbalancer_v2 %s: %s", d.Id(), err)
		}

		// Wait for load-balancer to become active before continuing.
		err = waitForLBV2LoadBalancer(ctx, lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Security Groups get updated separately.
	if d.HasChange("security_group_ids") {
		networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating OpenStack networking client: %s", err)
		}
		vipPortID := d.Get("vip_port_id").(string)
		if err := resourceLoadBalancerV2SetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error setting openstack_lb_loadbalancer_v2 security groups: %s", err)
		}
	}

	return resourceLoadBalancerV2Read(ctx, d, meta)
}

func resourceLoadBalancerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_lb_loadbalancer_v2 %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronloadbalancers.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_loadbalancer_v2"))
	}

	// Wait for load-balancer to become deleted.
	err = waitForLBV2LoadBalancer(ctx, lbClient, d.Id(), "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
