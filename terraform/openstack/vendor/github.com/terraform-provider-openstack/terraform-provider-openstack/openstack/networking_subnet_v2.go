package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

// networkingSubnetV2StateRefreshFunc returns a standard resource.StateRefreshFunc to wait for subnet status.
func networkingSubnetV2StateRefreshFunc(client *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := subnets.Get(client, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return subnet, "DELETED", nil
			}

			return nil, "", err
		}

		return subnet, "ACTIVE", nil
	}
}

// networkingSubnetV2StateRefreshFuncDelete returns a special case resource.StateRefreshFunc to try to delete a subnet.
func networkingSubnetV2StateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete openstack_networking_subnet_v2 %s", subnetID)

		s, err := subnets.Get(networkingClient, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_subnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}

			return s, "ACTIVE", err
		}

		err = subnets.Delete(networkingClient, subnetID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_subnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			// Subnet is still in use - we can retry.
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}

			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_subnet_v2 %s is still active", subnetID)

		return s, "ACTIVE", nil
	}
}

// networkingSubnetV2GetRawAllocationPoolsValueToExpand selects the resource argument to populate
// the allocations pool value.
func networkingSubnetV2GetRawAllocationPoolsValueToExpand(d *schema.ResourceData) []interface{} {
	// First check allocation_pool since that is the new argument.
	result := d.Get("allocation_pool").(*schema.Set).List()

	if len(result) == 0 {
		// If no allocation_pool was specified, check allocation_pools
		// which is the older legacy argument.
		result = d.Get("allocation_pools").([]interface{})
	}

	return result
}

// expandNetworkingSubnetV2AllocationPools returns a slice of subnets.AllocationPool structs.
func expandNetworkingSubnetV2AllocationPools(allocationPools []interface{}) []subnets.AllocationPool {
	result := make([]subnets.AllocationPool, len(allocationPools))
	for i, raw := range allocationPools {
		rawMap := raw.(map[string]interface{})

		result[i] = subnets.AllocationPool{
			Start: rawMap["start"].(string),
			End:   rawMap["end"].(string),
		}
	}

	return result
}

// flattenNetworkingSubnetV2AllocationPools allows to flatten slice of subnets.AllocationPool structs into
// a slice of maps.
func flattenNetworkingSubnetV2AllocationPools(allocationPools []subnets.AllocationPool) []map[string]interface{} {
	result := make([]map[string]interface{}, len(allocationPools))
	for i, allocationPool := range allocationPools {
		pool := make(map[string]interface{})
		pool["start"] = allocationPool.Start
		pool["end"] = allocationPool.End

		result[i] = pool
	}

	return result
}

// expandNetworkingSubnetV2HostRoutes returns a slice of HostRoute structures.
func expandNetworkingSubnetV2HostRoutes(rawHostRoutes []interface{}) []subnets.HostRoute {
	result := make([]subnets.HostRoute, len(rawHostRoutes))
	for i, raw := range rawHostRoutes {
		rawMap := raw.(map[string]interface{})

		result[i] = subnets.HostRoute{
			DestinationCIDR: rawMap["destination_cidr"].(string),
			NextHop:         rawMap["next_hop"].(string),
		}
	}

	return result
}

func networkingSubnetV2AllocationPoolsCustomizeDiff(diff *schema.ResourceDiff) error {
	if diff.Id() != "" && diff.HasChange("allocation_pools") {
		o, n := diff.GetChange("allocation_pools")
		oldPools := o.([]interface{})
		newPools := n.([]interface{})

		samePools := networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)

		if samePools {
			log.Printf("[DEBUG] allocation_pools have not changed. clearing diff")
			return diff.Clear("allocation_pools")
		}
	}

	return nil
}

func networkingSubnetV2AllocationPoolsMatch(oldPools, newPools []interface{}) bool {
	if len(oldPools) != len(newPools) {
		return false
	}

	for _, newPool := range newPools {
		var found bool

		newPoolPool := newPool.(map[string]interface{})
		newStart := newPoolPool["start"].(string)
		newEnd := newPoolPool["end"].(string)

		for _, oldPool := range oldPools {
			oldPoolPool := oldPool.(map[string]interface{})
			oldStart := oldPoolPool["start"].(string)
			oldEnd := oldPoolPool["end"].(string)

			if oldStart == newStart && oldEnd == newEnd {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func networkingSubnetV2DNSNameserverAreUnique(raw []interface{}) error {
	set := make(map[string]struct{})
	for _, rawNS := range raw {
		nameserver, ok := rawNS.(string)
		if ok {
			if _, exists := set[nameserver]; exists {
				return fmt.Errorf("got duplicate nameserver %s", nameserver)
			}
			set[nameserver] = struct{}{}
		}
	}

	return nil
}
