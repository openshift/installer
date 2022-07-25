package openstack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func resourceNetworkingQoSRuleV2BuildID(qosPolicyID, qosRuleID string) string {
	return fmt.Sprintf("%s/%s", qosPolicyID, qosRuleID)
}

func resourceNetworkingQoSRuleV2ParseID(qosRuleID string) (string, string, error) {
	qosRuleIDParts := strings.Split(qosRuleID, "/")
	if len(qosRuleIDParts) != 2 {
		return "", "", fmt.Errorf("invalid ID format: %s", qosRuleID)
	}

	return qosRuleIDParts[0], qosRuleIDParts[1], nil
}

func networkingQoSBandwidthLimitRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetBandwidthLimitRule(client, policyID, ruleID).ExtractBandwidthLimitRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSDSCPMarkingRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetDSCPMarkingRule(client, policyID, ruleID).ExtractDSCPMarkingRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSMinimumBandwidthRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetMinimumBandwidthRule(client, policyID, ruleID).ExtractMinimumBandwidthRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}
