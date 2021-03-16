package virtualmachine

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils/patch"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func VirtualMachineFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": k8s.NamespacedMetadataSchema("VirtualMachine", false),
		"spec":     virtualMachineSpecSchema(),
		"status":   virtualMachineStatusSchema(),
	}
}

func ExpandVirtualMachine(virtualMachine []interface{}) (*kubevirtapiv1.VirtualMachine, error) {
	result := &kubevirtapiv1.VirtualMachine{}

	if len(virtualMachine) == 0 || virtualMachine[0] == nil {
		return result, nil
	}

	in := virtualMachine[0].(map[string]interface{})

	if v, ok := in["metadata"].([]interface{}); ok {
		result.ObjectMeta = k8s.ExpandMetadata(v)
	}
	if v, ok := in["spec"].([]interface{}); ok {
		spec, err := expandVirtualMachineSpec(v)
		if err != nil {
			return result, err
		}
		result.Spec = spec
	}
	if v, ok := in["status"].([]interface{}); ok {
		status, err := expandVirtualMachineStatus(v)
		if err != nil {
			return result, err
		}
		result.Status = status
	}

	return result, nil
}

func FlattenVirtualMachine(in kubevirtapiv1.VirtualMachine) []interface{} {
	att := make(map[string]interface{})

	att["metadata"] = k8s.FlattenMetadata(in.ObjectMeta)
	att["spec"] = flattenVirtualMachineSpec(in.Spec)
	att["status"] = flattenVirtualMachineStatus(in.Status)

	return []interface{}{att}
}

func FromResourceData(resourceData *schema.ResourceData) (*kubevirtapiv1.VirtualMachine, error) {
	result := &kubevirtapiv1.VirtualMachine{}

	result.ObjectMeta = k8s.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandVirtualMachineSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	status, err := expandVirtualMachineStatus(resourceData.Get("status").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Status = status

	return result, nil
}

func ToResourceData(vm kubevirtapiv1.VirtualMachine, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", k8s.FlattenMetadata(vm.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenVirtualMachineSpec(vm.Spec)); err != nil {
		return err
	}
	if err := resourceData.Set("status", flattenVirtualMachineStatus(vm.Status)); err != nil {
		return err
	}

	return nil
}

func AppendPatchOps(keyPrefix, pathPrefix string, resourceData *schema.ResourceData, ops []patch.PatchOperation) patch.PatchOperations {
	return k8s.AppendPatchOps(keyPrefix+"metadata.0.", pathPrefix+"/metadata/", resourceData, ops)
}
