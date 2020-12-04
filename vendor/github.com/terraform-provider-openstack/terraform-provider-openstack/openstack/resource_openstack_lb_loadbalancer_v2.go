package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	neutronloadbalancers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
)

func resourceLoadBalancerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerV2Create,
		Read:   resourceLoadBalancerV2Read,
		Update: resourceLoadBalancerV2Update,
		Delete: resourceLoadBalancerV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				ForceNew: true,
			},

			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceLoadBalancerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		lbID      string
		vipPortID string
	)

	// Choose either the Octavia or Neutron create options.
	createOpts := chooseLBV2LoadBalancerCreateOpts(d, config)

	if lbClient.Type == octaviaLBClientType {
		log.Printf("[DEBUG][Octavia] openstack_lb_loadbalancer_v2 create options: %#v", createOpts)
		lb, err := octavialoadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error creating openstack_lb_loadbalancer_v2: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	} else {
		log.Printf("[DEBUG][Neutron] openstack_lb_loadbalancer_v2 create options: %#v", createOpts)
		lb, err := neutronloadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error creating openstack_lb_loadbalancer_v2: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	}

	// Wait for load-balancer to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return err
	}

	// Once the load-balancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	if err := resourceLoadBalancerV2SetSecurityGroups(networkingClient, vipPortID, d); err != nil {
		return fmt.Errorf("Error setting openstack_lb_loadbalancer_v2 security groups: %s", err)
	}

	d.SetId(lbID)

	return resourceLoadBalancerV2Read(d, meta)
}

func resourceLoadBalancerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var vipPortID string

	if lbClient.Type == octaviaLBClientType {
		lb, err := octavialoadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "Unable to retrieve openstack_lb_loadbalancer_v2")
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
		d.Set("region", GetRegion(d, config))
		vipPortID = lb.VipPortID
	} else {
		lb, err := neutronloadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "Unable to retrieve openstack_lb_loadbalancer_v2")
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
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}
		if err := resourceLoadBalancerV2GetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return fmt.Errorf("Error getting port security groups for openstack_lb_loadbalancer_v2: %s", err)
		}
	}

	return nil
}

func resourceLoadBalancerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts neutronloadbalancers.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if updateOpts != (neutronloadbalancers.UpdateOpts{}) {
		// Wait for load-balancer to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBV2LoadBalancer(lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return err
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
			return fmt.Errorf("Error updating openstack_lb_loadbalancer_v2 %s: %s", d.Id(), err)
		}

		// Wait for load-balancer to become active before continuing.
		err = waitForLBV2LoadBalancer(lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return err
		}
	}

	// Security Groups get updated separately.
	if d.HasChange("security_group_ids") {
		networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}
		vipPortID := d.Get("vip_port_id").(string)
		if err := resourceLoadBalancerV2SetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return fmt.Errorf("Error setting openstack_lb_loadbalancer_v2 security groups: %s", err)
		}
	}

	return resourceLoadBalancerV2Read(d, meta)
}

func resourceLoadBalancerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := chooseLBV2Client(d, config)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
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
		return CheckDeleted(d, err, "Error deleting openstack_lb_loadbalancer_v2")
	}

	// Wait for load-balancer to become deleted.
	err = waitForLBV2LoadBalancer(lbClient, d.Id(), "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return err
	}

	return nil
}
