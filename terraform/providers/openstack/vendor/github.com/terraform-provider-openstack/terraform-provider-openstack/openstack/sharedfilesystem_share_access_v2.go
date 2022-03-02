package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func sharedFilesystemShareAccessV2StateRefreshFunc(client *gophercloud.ServiceClient, shareID string, accessID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		access, err := shares.ListAccessRights(client, shareID).Extract()
		if err != nil {
			return nil, "", err
		}
		for _, v := range access {
			if v.ID == accessID {
				return v, v.State, nil
			}
		}
		return nil, "", gophercloud.ErrDefault404{}
	}
}
