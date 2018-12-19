package openstack

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/hashicorp/terraform/helper/schema"
)

func expandNetworkingPortDHCPOptsV2Create(dhcpOpts *schema.Set) []extradhcpopts.CreateExtraDHCPOpt {
	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.CreateExtraDHCPOpt, len(rawDHCPOpts))
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})

		ipVersion := rawMap["ip_version"].(int)
		optName := rawMap["name"].(string)
		optValue := rawMap["value"].(string)

		extraDHCPOpts[i] = extradhcpopts.CreateExtraDHCPOpt{
			OptName:   optName,
			OptValue:  optValue,
			IPVersion: gophercloud.IPVersion(ipVersion),
		}
	}

	return extraDHCPOpts
}

func expandNetworkingPortDHCPOptsV2Update(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.UpdateExtraDHCPOpt, len(rawDHCPOpts))
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})

		ipVersion := rawMap["ip_version"].(int)
		optName := rawMap["name"].(string)
		optValue := rawMap["value"].(string)

		extraDHCPOpts[i] = extradhcpopts.UpdateExtraDHCPOpt{
			OptName:   optName,
			OptValue:  &optValue,
			IPVersion: gophercloud.IPVersion(ipVersion),
		}
	}

	return extraDHCPOpts
}

func expandNetworkingPortDHCPOptsV2Delete(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	if dhcpOpts == nil {
		return []extradhcpopts.UpdateExtraDHCPOpt{}
	}

	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.UpdateExtraDHCPOpt, len(rawDHCPOpts))
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})
		extraDHCPOpts[i] = extradhcpopts.UpdateExtraDHCPOpt{
			OptName:  rawMap["name"].(string),
			OptValue: nil,
		}
	}

	return extraDHCPOpts
}

func flattenNetworkingPortDHCPOptsV2(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) []map[string]interface{} {
	dhcpOptsSet := make([]map[string]interface{}, len(dhcpOpts.ExtraDHCPOpts))

	for i, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet[i] = map[string]interface{}{
			"ip_version": dhcpOpt.IPVersion,
			"name":       dhcpOpt.OptName,
			"value":      dhcpOpt.OptValue,
		}
	}

	return dhcpOptsSet
}
