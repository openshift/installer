package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func resourceNetworkingSubnetRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSubnetRouteV2Create,
		ReadContext:   resourceNetworkingSubnetRouteV2Read,
		DeleteContext: resourceNetworkingSubnetRouteV2Delete,
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

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingSubnetRouteV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	subnetID := d.Get("subnet_id").(string)
	config.MutexKV.Lock(subnetID)
	defer config.MutexKV.Unlock(subnetID)

	subnet, err := subnets.Get(networkingClient, subnetID).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving openstack_networking_subnet_v2: %s", err)
	}

	destCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	for _, r := range subnet.HostRoutes {
		if r.DestinationCIDR == destCIDR && r.NextHop == nextHop {
			return diag.Errorf(
				"openstack_networking_subnet_v2 %s already has a route to %s via %s",
				subnetID,
				r.DestinationCIDR,
				r.NextHop,
			)
		}
	}

	// Add a new route.
	subnet.HostRoutes = append(subnet.HostRoutes, subnets.HostRoute{
		DestinationCIDR: destCIDR,
		NextHop:         nextHop,
	})

	log.Printf(
		"[DEBUG] Adding openstack_networking_subnet_v2 %s route to %s via %s",
		subnetID,
		destCIDR,
		nextHop,
	)
	updateOpts := subnets.UpdateOpts{
		HostRoutes: &subnet.HostRoutes,
	}
	log.Printf("[DEBUG] Updating openstack_networking_subnet_v2 %s with options: %+v", subnetID, updateOpts)
	_, err = subnets.Update(networkingClient, subnetID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_subnet_v2: %s", err)
	}

	d.SetId(resourceNetworkingSubnetRouteV2BuildID(subnetID, destCIDR, nextHop))

	return resourceNetworkingSubnetRouteV2Read(ctx, d, meta)
}

func resourceNetworkingSubnetRouteV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	subnetID, destCIDR, nextHop, err := resourceNetworkingSubnetRouteV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_subnet_route_v2 ID %s: %s", d.Id(), err)
	}

	subnet, err := subnets.Get(networkingClient, subnetID).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving openstack_networking_subnet_v2: %s", err)
	}

	exists := false
	for _, r := range subnet.HostRoutes {
		if r.DestinationCIDR == destCIDR && r.NextHop == nextHop {
			exists = true
		}
	}
	if !exists {
		return diag.Errorf(
			"openstack_networking_subnet_v2 %s doesn't have a route to %s via %s",
			subnetID,
			destCIDR,
			nextHop,
		)
	}

	d.Set("subnet_id", subnetID)
	d.Set("next_hop", nextHop)
	d.Set("destination_cidr", destCIDR)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingSubnetRouteV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	subnetID := d.Get("subnet_id").(string)
	config.MutexKV.Lock(subnetID)
	defer config.MutexKV.Unlock(subnetID)

	subnet, err := subnets.Get(networkingClient, subnetID).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return nil
		}

		return diag.Errorf("Error retrieving openstack_networking_subnet_v2: %s", err)
	}

	var destCIDR = d.Get("destination_cidr").(string)
	var nextHop = d.Get("next_hop").(string)

	oldRoutes := subnet.HostRoutes
	newRoutes := make([]subnets.HostRoute, 0, 1)

	for _, r := range oldRoutes {
		if r.DestinationCIDR != destCIDR || r.NextHop != nextHop {
			newRoutes = append(newRoutes, r)
		}
	}

	if len(oldRoutes) == len(newRoutes) {
		return diag.Errorf(
			"openstack_networking_subnet_v2 %s already doesn't have a route to %s via %s",
			subnetID,
			destCIDR,
			nextHop,
		)
	}

	log.Printf(
		"[DEBUG] Deleting openstack_networking_subnet_v2 %s route to %s via %s",
		subnetID,
		destCIDR,
		nextHop,
	)
	updateOpts := subnets.UpdateOpts{
		HostRoutes: &newRoutes,
	}
	log.Printf("[DEBUG] Updating openstack_networking_subnet_v2 %s with options: %#v", subnetID, updateOpts)
	_, err = subnets.Update(networkingClient, subnetID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_subnet_v2: %s", err)
	}

	return nil
}
