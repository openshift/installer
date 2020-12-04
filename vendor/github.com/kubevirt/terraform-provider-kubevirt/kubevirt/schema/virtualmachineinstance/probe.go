package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func probeFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// TODO nargaman
	}
}

func probeSchema() *schema.Schema {
	fields := probeFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Specification of the desired behavior of the VirtualMachineInstance on the host."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandProbe(probe []interface{}) *kubevirtapiv1.Probe {
	if len(probe) == 0 || probe[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.Probe{}

	_ = probe[0].(map[string]interface{})

	// TODO nargaman

	return result
}

func flattenProbe(in kubevirtapiv1.Probe) []interface{} {
	att := make(map[string]interface{})

	// att["spec"] = flattenVirtualMachineInstanceSpecSpec(in.Spec)
	// att["status"] = flattenVirtualMachineInstanceSpecStatus(in.Status)
	// TODO nargaman

	return []interface{}{att}
}
