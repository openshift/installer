package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func virtualMachineInstanceTemplateSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": k8s.NamespacedMetadataSchema("VirtualMachineInstanceTemplateSpec", false),
		"spec":     virtualMachineInstanceSpecSchema(),
	}
}

func VirtualMachineInstanceTemplateSpecSchema() *schema.Schema {
	fields := virtualMachineInstanceTemplateSpecFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Template is the direct specification of VirtualMachineInstance."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandVirtualMachineInstanceTemplateSpec(virtualMachine []interface{}) (*kubevirtapiv1.VirtualMachineInstanceTemplateSpec, error) {
	if len(virtualMachine) == 0 || virtualMachine[0] == nil {
		return nil, nil
	}

	result := &kubevirtapiv1.VirtualMachineInstanceTemplateSpec{}

	in := virtualMachine[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = k8s.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandVirtualMachineInstanceSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}

	return result, nil
}

func FlattenVirtualMachineInstanceTemplateSpec(in kubevirtapiv1.VirtualMachineInstanceTemplateSpec) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = k8s.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenVirtualMachineInstanceSpec(in.Spec)

	return []interface{}{att}
}
