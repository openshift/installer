package openstack

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
)

func networkingV2ReadAttributesTags(d *schema.ResourceData, tags []string) {
	expandObjectReadTags(d, tags)
}

func networkingV2UpdateAttributesTags(d *schema.ResourceData) []string {
	return expandObjectUpdateTags(d)
}

func networkingV2AttributesTags(d *schema.ResourceData) []string {
	return expandObjectTags(d)
}

type neutronErrorWrap struct {
	NeutronError neutronError
}

type neutronError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Detail  string `json:"detail"`
}

func retryOn409(err error) bool {
	switch err := err.(type) {
	case gophercloud.ErrDefault409:
		neutronError, e := decodeNeutronError(err.ErrUnexpectedResponseCode.Body)
		if e != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", e)
			return true
		}
		if neutronError.Type == "IpAddressGenerationFailure" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case gophercloud.ErrDefault400:
		neutronError, e := decodeNeutronError(err.ErrUnexpectedResponseCode.Body)
		if e != nil {
			// retry, when error type cannot be detected
			log.Printf("[DEBUG] failed to decode a neutron error: %s", e)
			return true
		}
		if neutronError.Type == "ExternalIpAddressExhausted" {
			return true
		}

		// don't retry on quota or other errors
		return false
	case gophercloud.ErrDefault404: // this case is handled mostly for functional tests
		return true
	}

	return false
}

func decodeNeutronError(body []byte) (*neutronError, error) {
	e := &neutronErrorWrap{}
	if err := json.Unmarshal(body, e); err != nil {
		return nil, err
	}

	return &e.NeutronError, nil
}
