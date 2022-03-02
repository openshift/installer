package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/request"
)

// TransferRequestCreateOpts represents the attributes used when creating a new transfer request.
type TransferRequestCreateOpts struct {
	request.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToTransferRequestCreateMap casts a CreateOpts struct to a map.
// It overrides request.ToTransferRequestCreateMap to add the ValueSpecs field.
func (opts TransferRequestCreateOpts) ToTransferRequestCreateMap() (map[string]interface{}, error) {
	b, err := BuildRequest(opts, "")
	if err != nil {
		return nil, err
	}

	if m, ok := b[""].(map[string]interface{}); ok {
		return m, nil
	}

	return nil, fmt.Errorf("Expected map but got %T", b[""])
}

func dnsTransferRequestV2RefreshFunc(dnsClient *gophercloud.ServiceClient, transferRequestID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		transferRequest, err := request.Get(dnsClient, transferRequestID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return transferRequest, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_dns_transfer_request_v2 %s current status: %s", transferRequest.ID, transferRequest.Status)
		return transferRequest, transferRequest.Status, nil
	}
}
