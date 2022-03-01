package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func resourceNetworkingRouterV2StateRefreshFunc(client *gophercloud.ServiceClient, routerID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := routers.Get(client, routerID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}

func expandNetworkingRouterExternalFixedIPsV2(externalFixedIPs []interface{}) []routers.ExternalFixedIP {
	fixedIPs := make([]routers.ExternalFixedIP, len(externalFixedIPs))

	for i, raw := range externalFixedIPs {
		rawMap := raw.(map[string]interface{})

		fixedIPs[i] = routers.ExternalFixedIP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		}
	}

	return fixedIPs
}

func expandNetworkingRouterExternalSubnetIDsV2(externalSubnetIDs []interface{}) []routers.ExternalFixedIP {
	subnetIDs := make([]routers.ExternalFixedIP, len(externalSubnetIDs))

	for i, raw := range externalSubnetIDs {
		subnetIDs[i] = routers.ExternalFixedIP{
			SubnetID: raw.(string),
		}
	}

	return subnetIDs
}

func flattenNetworkingRouterExternalFixedIPsV2(externalFixedIPs []routers.ExternalFixedIP) []map[string]string {
	fixedIPs := make([]map[string]string, len(externalFixedIPs))

	for i, fixedIP := range externalFixedIPs {
		fixedIPs[i] = map[string]string{
			"subnet_id":  fixedIP.SubnetID,
			"ip_address": fixedIP.IPAddress,
		}
	}

	return fixedIPs
}
