package datavolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils/patch"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func DataVolumeFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": k8s.NamespacedMetadataSchema("DataVolume", false),
		"spec":     dataVolumeSpecSchema(),
		"status":   dataVolumeStatusSchema(),
	}
}

func DataVolumeTemplatesSchema() *schema.Schema {
	fields := DataVolumeFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("dataVolumeTemplates is a list of dataVolumes that the VirtualMachineInstance template can reference."),
		Required:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandDataVolumeTemplates(dataVolumes []interface{}) ([]cdiv1.DataVolume, error) {
	result := make([]cdiv1.DataVolume, len(dataVolumes))

	if len(dataVolumes) == 0 || dataVolumes[0] == nil {
		return result, nil
	}

	for i, dataVolume := range dataVolumes {
		in := dataVolume.(map[string]interface{})

		if v, ok := in["metadata"].([]interface{}); ok {
			result[i].ObjectMeta = k8s.ExpandMetadata(v)
		}
		if v, ok := in["spec"].([]interface{}); ok {
			spec, err := expandDataVolumeSpec(v)
			if err != nil {
				return result, err
			}
			result[i].Spec = spec
		}
		if v, ok := in["status"].([]interface{}); ok {
			result[i].Status = expandDataVolumeStatus(v)
		}
	}

	return result, nil
}

func FlattenDataVolumeTemplates(in []cdiv1.DataVolume) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})
		c["metadata"] = k8s.FlattenMetadata(v.ObjectMeta)
		c["spec"] = flattenDataVolumeSpec(v.Spec)
		c["status"] = flattenDataVolumeStatus(v.Status)
		att[i] = c
	}

	return att
}

func FromResourceData(resourceData *schema.ResourceData) (*cdiv1.DataVolume, error) {
	result := &cdiv1.DataVolume{}

	result.ObjectMeta = k8s.ExpandMetadata(resourceData.Get("metadata").([]interface{}))
	spec, err := expandDataVolumeSpec(resourceData.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = spec
	result.Status = expandDataVolumeStatus(resourceData.Get("status").([]interface{}))

	return result, nil
}

func ToResourceData(dv cdiv1.DataVolume, resourceData *schema.ResourceData) error {
	if err := resourceData.Set("metadata", k8s.FlattenMetadata(dv.ObjectMeta)); err != nil {
		return err
	}
	if err := resourceData.Set("spec", flattenDataVolumeSpec(dv.Spec)); err != nil {
		return err
	}
	if err := resourceData.Set("status", flattenDataVolumeStatus(dv.Status)); err != nil {
		return err
	}

	return nil
}

func AppendPatchOps(keyPrefix, pathPrefix string, resourceData *schema.ResourceData, ops []patch.PatchOperation) patch.PatchOperations {
	return k8s.AppendPatchOps(keyPrefix+"metadata.0.", pathPrefix+"/metadata/", resourceData, ops)
}
