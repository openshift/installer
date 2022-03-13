package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/accept"
)

// TransferAcceptCreateOpts represents the attributes used when creating a new transfer accept.
type TransferAcceptCreateOpts struct {
	accept.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToTransferAcceptCreateMap casts a CreateOpts struct to a map.
// It overrides accept.ToTransferAcceptCreateMap to add the ValueSpecs field.
func (opts TransferAcceptCreateOpts) ToTransferAcceptCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

func dnsTransferAcceptV2RefreshFunc(dnsClient *gophercloud.ServiceClient, transferAcceptID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		transferAccept, err := accept.Get(dnsClient, transferAcceptID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return transferAccept, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_transfer_accept_v2 %s current status: %s", transferAccept.ID, transferAccept.Status)
		return transferAccept, transferAccept.Status, nil
	}
}
