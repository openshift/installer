package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func networkFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Network name.",
			Required:    true,
		},
		"network_source": {
			Type:        schema.TypeList,
			Description: "NetworkSource represents the network type and the source interface that should be connected to the virtual machine.",
			MaxItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"pod": {
						Type:        schema.TypeList,
						Description: "Pod network.",
						MaxItems:    1,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"vm_network_cidr": {
									Type:        schema.TypeString,
									Description: "CIDR for vm network.",
									Required:    true,
								},
							},
						},
					},
					"multus": {
						Type:        schema.TypeList,
						Description: "Multus network.",
						MaxItems:    1,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"network_name": {
									Type:        schema.TypeString,
									Description: "References to a NetworkAttachmentDefinition CRD object. Format: <networkName>, <namespace>/<networkName>. If namespace is not specified, VMI namespace is assumed.",
									Required:    true,
								},
								"default": {
									Type:        schema.TypeBool,
									Description: "Select the default network and add it to the multus-cni.io/default-network annotation.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func networksSchema() *schema.Schema {
	fields := networkFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("List of networks that can be attached to a vm's virtual interface."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandNetworks(networks []interface{}) []kubevirtapiv1.Network {
	result := make([]kubevirtapiv1.Network, len(networks))

	if len(networks) == 0 || networks[0] == nil {
		return result
	}

	for i, network := range networks {
		in := network.(map[string]interface{})

		if v, ok := in["name"].(string); ok {
			result[i].Name = v
		}
		if v, ok := in["network_source"].([]interface{}); ok {
			result[i].NetworkSource = expandNetworkSource(v)
		}
	}

	return result
}

func expandNetworkSource(networkSource []interface{}) kubevirtapiv1.NetworkSource {
	result := kubevirtapiv1.NetworkSource{}

	if len(networkSource) == 0 || networkSource[0] == nil {
		return result
	}

	in := networkSource[0].(map[string]interface{})

	if v, ok := in["pod"].([]interface{}); ok {
		result.Pod = expandPodNetwork(v)
	}
	if v, ok := in["multus"].([]interface{}); ok {
		result.Multus = expandMultusNetwork(v)
	}

	return result
}

func expandPodNetwork(pod []interface{}) *kubevirtapiv1.PodNetwork {
	if len(pod) == 0 || pod[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.PodNetwork{}
	in := pod[0].(map[string]interface{})

	if v, ok := in["vm_network_cidr"].(string); ok {
		result.VMNetworkCIDR = v
	}

	return result
}

func expandMultusNetwork(multus []interface{}) *kubevirtapiv1.MultusNetwork {
	if len(multus) == 0 || multus[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.MultusNetwork{}
	in := multus[0].(map[string]interface{})

	if v, ok := in["network_name"].(string); ok {
		result.NetworkName = v
	}
	if v, ok := in["default"].(bool); ok {
		result.Default = v
	}

	return result
}

func flattenNetworks(in []kubevirtapiv1.Network) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})

		c["name"] = v.Name
		c["network_source"] = flattenNetworkSource(v.NetworkSource)

		att[i] = c
	}

	return att
}

func flattenNetworkSource(in kubevirtapiv1.NetworkSource) []interface{} {
	att := make(map[string]interface{})

	if in.Pod != nil {
		att["pod"] = flattenPodNetwork(*in.Pod)
	}
	if in.Multus != nil {
		att["multus"] = flattenMultusNetwork(*in.Multus)
	}

	return []interface{}{att}
}

func flattenPodNetwork(in kubevirtapiv1.PodNetwork) []interface{} {
	att := make(map[string]interface{})

	att["vm_network_cidr"] = in.VMNetworkCIDR

	return []interface{}{att}
}

func flattenMultusNetwork(in kubevirtapiv1.MultusNetwork) []interface{} {
	att := make(map[string]interface{})

	att["network_name"] = in.NetworkName
	att["default"] = in.Default

	return []interface{}{att}
}
