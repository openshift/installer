package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/addressscopes"
)

func resourceNetworkingAddressScopeV2StateRefreshFunc(client *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		a, err := addressscopes.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return a, "DELETED", nil
			}

			return nil, "", err
		}

		return a, "ACTIVE", nil
	}
}
