package virtualmachine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func virtualMachineStatusFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created": &schema.Schema{
			Type:        schema.TypeBool,
			Description: "Created indicates if the virtual machine is created in the cluster.",
			Optional:    true,
		},
		"ready": &schema.Schema{
			Type:        schema.TypeBool,
			Description: "Ready indicates if the virtual machine is running and ready.",
			Optional:    true,
		},
		"conditions":            virtualMachineConditionsSchema(),
		"state_change_requests": virtualMachineStateChangeRequestsSchema(),
	}
}

func virtualMachineStatusSchema() *schema.Schema {
	fields := virtualMachineStatusFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("VirtualMachineStatus represents the status returned by the controller to describe how the VirtualMachine is doing."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandVirtualMachineStatus(virtualMachineStatus []interface{}) (kubevirtapiv1.VirtualMachineStatus, error) {
	result := kubevirtapiv1.VirtualMachineStatus{}

	if len(virtualMachineStatus) == 0 || virtualMachineStatus[0] == nil {
		return result, nil
	}

	in := virtualMachineStatus[0].(map[string]interface{})

	if v, ok := in["created"].(bool); ok {
		result.Created = v
	}
	if v, ok := in["ready"].(bool); ok {
		result.Ready = v
	}
	if v, ok := in["conditions"].([]interface{}); ok {
		conditions, err := expandVirtualMachineConditions(v)
		if err != nil {
			return result, err
		}
		result.Conditions = conditions
	}
	if v, ok := in["state_change_requests"].([]interface{}); ok {
		result.StateChangeRequests = expandVirtualMachineStateChangeRequests(v)
	}

	return result, nil
}

func flattenVirtualMachineStatus(in kubevirtapiv1.VirtualMachineStatus) []interface{} {
	att := make(map[string]interface{})

	att["created"] = in.Created
	att["ready"] = in.Ready
	att["conditions"] = flattenVirtualMachineConditions(in.Conditions)
	att["state_change_requests"] = flattenVirtualMachineStateChangeRequests(in.StateChangeRequests)

	return []interface{}{att}
}
