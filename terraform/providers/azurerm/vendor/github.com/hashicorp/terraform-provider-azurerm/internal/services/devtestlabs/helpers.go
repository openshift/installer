package devtestlabs

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func schemaDevTestVirtualMachineInboundNatRule() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		// since these aren't returned from the API
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"protocol": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(virtualmachines.TransportProtocolTcp),
						string(virtualmachines.TransportProtocolUdp),
					}, false),
				},

				"backend_port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.PortNumber,
				},

				"frontend_port": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func expandDevTestLabVirtualMachineNatRules(input *pluginsdk.Set) []virtualmachines.InboundNatRule {
	rules := make([]virtualmachines.InboundNatRule, 0)
	if input == nil {
		return rules
	}

	for _, val := range input.List() {
		v := val.(map[string]interface{})
		backendPort := v["backend_port"].(int)
		protocol := virtualmachines.TransportProtocol(v["protocol"].(string))

		rule := virtualmachines.InboundNatRule{
			TransportProtocol: &protocol,
			BackendPort:       utils.Int64(int64(backendPort)),
		}

		rules = append(rules, rule)
	}

	return rules
}

func expandDevTestLabVirtualMachineGalleryImageReference(input []interface{}, osType string) *virtualmachines.GalleryImageReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	offer := v["offer"].(string)
	publisher := v["publisher"].(string)
	sku := v["sku"].(string)
	version := v["version"].(string)

	return &virtualmachines.GalleryImageReference{
		Offer:     utils.String(offer),
		OsType:    utils.String(osType),
		Publisher: utils.String(publisher),
		Sku:       utils.String(sku),
		Version:   utils.String(version),
	}
}

func schemaDevTestVirtualMachineGalleryImageReference() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"offer": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"publisher": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"sku": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
				"version": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func flattenDevTestVirtualMachineGalleryImage(input *virtualmachines.GalleryImageReference) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		output := make(map[string]interface{})

		if input.Offer != nil {
			output["offer"] = *input.Offer
		}

		if input.Publisher != nil {
			output["publisher"] = *input.Publisher
		}

		if input.Sku != nil {
			output["sku"] = *input.Sku
		}

		if input.Version != nil {
			output["version"] = *input.Version
		}

		results = append(results, output)
	}

	return results
}
