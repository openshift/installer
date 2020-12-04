package virtualmachine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	k8sv1 "k8s.io/api/core/v1"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func virtualMachineConditionsFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "VirtualMachineConditionType represent the type of the VM as concluded from its VMi status.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Failure",
				"Ready",
				"Paused",
				"RenameOperation",
			}, false),
		},
		"status": {
			Type:        schema.TypeString,
			Description: "ConditionStatus represents the status of this VM condition, if the VM currently in the condition.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"True",
				"False",
				"Unknown",
			}, false),
		},
		// TODO nargaman -  Add following values
		// "last_probe_time": {
		// 	Type:        schema.TypeString,
		// 	Description: "Last probe time.",
		// 	Optional:    true,
		// },
		// "last_transition_time": {
		// 	Type:        schema.TypeString,
		// 	Description: "Last transition time.",
		// 	Optional:    true,
		// },
		"reason": {
			Type:        schema.TypeString,
			Description: "Condition reason.",
			Optional:    true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "Condition message.",
			Optional:    true,
		},
	}
}

func virtualMachineConditionsSchema() *schema.Schema {
	fields := virtualMachineConditionsFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Hold the state information of the VirtualMachine and its VirtualMachineInstance."),
		Required:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandVirtualMachineConditions(conditions []interface{}) ([]kubevirtapiv1.VirtualMachineCondition, error) {
	result := make([]kubevirtapiv1.VirtualMachineCondition, len(conditions))

	if len(conditions) == 0 || conditions[0] == nil {
		return result, nil
	}

	for i, condition := range conditions {
		in := condition.(map[string]interface{})

		if v, ok := in["type"].(string); ok {
			result[i].Type = kubevirtapiv1.VirtualMachineConditionType(v)
		}
		if v, ok := in["status"].(string); ok {
			result[i].Status = k8sv1.ConditionStatus(v)
		}
		if v, ok := in["reason"].(string); ok {
			result[i].Reason = v
		}
		if v, ok := in["message"].(string); ok {
			result[i].Message = v
		}
	}

	return result, nil
}

func flattenVirtualMachineConditions(in []kubevirtapiv1.VirtualMachineCondition) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})
		c["type"] = string(v.Type)
		c["status"] = string(v.Status)
		c["reason"] = v.Reason
		c["message"] = v.Message

		att[i] = c
	}

	return att
}
