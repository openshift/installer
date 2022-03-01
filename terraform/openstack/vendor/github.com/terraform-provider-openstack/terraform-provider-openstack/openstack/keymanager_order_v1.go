package openstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/orders"
)

func keyManagerOrderV1WaitForOrderDeletion(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := orders.Delete(kmClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func keyManagerOrderV1OrderType(v string) orders.OrderType {
	var otype orders.OrderType
	switch v {
	case "asymmetric":
		otype = orders.AsymmetricOrder
	case "key":
		otype = orders.KeyOrder
	}

	return otype
}

func keyManagerOrderV1WaitForOrderCreation(kmClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		order, err := orders.Get(kmClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "NOT_CREATED", nil
			}

			return "", "NOT_CREATED", err
		}

		if order.Status == "ERROR" {
			return "", order.Status, fmt.Errorf("Error creating order")
		}

		return order, order.Status, nil
	}
}

func keyManagerOrderV1GetUUIDfromOrderRef(ref string) string {
	// order ref has form https://{barbican_host}/v1/orders/{order_uuid}
	// so we are only interested in the last part
	refSplit := strings.Split(ref, "/")
	uuid := refSplit[len(refSplit)-1]
	return uuid
}

func expandKeyManagerOrderV1Meta(s []interface{}) orders.MetaOpts {
	var meta orders.MetaOpts
	m := s[0].(map[string]interface{})

	if v, ok := m["algorithm"]; ok {
		meta.Algorithm = v.(string)
	}

	if v, ok := m["bit_length"]; ok {
		meta.BitLength = v.(int)
	}

	if v, ok := m["expiration"]; ok {
		if t, _ := time.Parse(time.RFC3339, v.(string)); t != (time.Time{}) {
			meta.Expiration = &t
		}
	}

	if v, ok := m["mode"]; ok {
		meta.Mode = v.(string)
	}

	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}

	if v, ok := m["payload_content_type"]; ok {
		meta.PayloadContentType = v.(string)
	}

	return meta
}

func flattenKeyManagerOrderV1Meta(m orders.Meta) []map[string]interface{} {
	var meta []map[string]interface{}
	s := make(map[string]interface{})

	if m.Algorithm != "" {
		s["algorithm"] = m.Algorithm
	}

	if m.BitLength != 0 {
		s["bit_length"] = m.BitLength
	}

	if !m.Expiration.IsZero() {
		s["expiration"] = m.Expiration.UTC().Format(time.RFC3339)
	}

	if m.Mode != "" {
		s["mode"] = m.Mode
	}

	if m.Name != "" {
		s["name"] = m.Name
	}

	if m.PayloadContentType != "" {
		s["payload_content_type"] = m.PayloadContentType
	}

	return append(meta, s)
}
