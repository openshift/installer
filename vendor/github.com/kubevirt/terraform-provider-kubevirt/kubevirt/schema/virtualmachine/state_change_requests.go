package virtualmachine

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	k8stypes "k8s.io/apimachinery/pkg/types"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func virtualMachineStateChangeRequestFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:        schema.TypeString,
			Description: "Indicates the type of action that is requested. e.g. Start or Stop.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Start",
				"Stop",
			}, false),
		},
		"data": {
			Type:        schema.TypeMap,
			Description: "Provides additional data in order to perform the Action.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"uid": {
			Type:        schema.TypeString,
			Description: "Indicates the UUID of an existing Virtual Machine Instance that this change request applies to -- if applicable.",
			Optional:    true,
		},
	}
}

func virtualMachineStateChangeRequestsSchema() *schema.Schema {
	fields := virtualMachineStateChangeRequestFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("StateChangeRequests indicates a list of actions that should be taken on a VMI."),
		Required:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandVirtualMachineStateChangeRequests(virtualMachineStateChangeRequests []interface{}) []kubevirtapiv1.VirtualMachineStateChangeRequest {
	result := make([]kubevirtapiv1.VirtualMachineStateChangeRequest, len(virtualMachineStateChangeRequests))

	if len(virtualMachineStateChangeRequests) == 0 || virtualMachineStateChangeRequests[0] == nil {
		return result
	}

	for i, virtualMachineStateChangeRequest := range virtualMachineStateChangeRequests {
		in := virtualMachineStateChangeRequest.(map[string]interface{})

		if v, ok := in["action"].(string); ok {
			result[i].Action = kubevirtapiv1.StateChangeRequestAction(v)
		}
		if v, ok := in["data"].(map[string]interface{}); ok && len(v) > 0 {
			result[i].Data = utils.ExpandStringMap(v)
		}
		if v, ok := in["uid"].(string); ok {
			uid := k8stypes.UID(v)
			result[i].UID = &uid
		}
	}

	return result
}

func flattenVirtualMachineStateChangeRequests(in []kubevirtapiv1.VirtualMachineStateChangeRequest) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})

		c["action"] = string(v.Action)
		c["data"] = v.Data
		if v.UID != nil {
			c["uid"] = string(*v.UID)
		}

		att[i] = c
	}

	return att
}
