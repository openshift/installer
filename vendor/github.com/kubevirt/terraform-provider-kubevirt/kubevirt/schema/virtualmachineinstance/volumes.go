package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func volumesFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Volume's name.",
			Required:    true,
		},
		"volume_source": {
			Type:        schema.TypeList,
			Description: "VolumeSource represents the location and type of the mounted volume. Defaults to Disk, if no type is specified.",
			MaxItems:    1,
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"data_volume": {
						Type:        schema.TypeList,
						Description: "DataVolume represents the dynamic creation a PVC for this volume as well as the process of populating that PVC with a disk image.",
						MaxItems:    1,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Description: "Name represents the name of the DataVolume in the same namespace.",
									Required:    true,
								},
							},
						},
					},
					"cloud_init_config_drive": {
						Type:        schema.TypeList,
						Description: "CloudInitConfigDrive represents a cloud-init Config Drive user-data source.",
						MaxItems:    1,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"user_data_secret_ref": k8s.LocalObjectReferenceSchema("UserDataSecretRef references a k8s secret that contains config drive userdata."),
								"user_data_base64": {
									Type:        schema.TypeString,
									Description: "UserDataBase64 contains config drive cloud-init userdata as a base64 encoded string.",
									Optional:    true,
								},
								"user_data": {
									Type:        schema.TypeString,
									Description: "UserData contains config drive inline cloud-init userdata.",
									Optional:    true,
								},
								"network_data_secret_ref": k8s.LocalObjectReferenceSchema("NetworkDataSecretRef references a k8s secret that contains config drive networkdata."),
								"network_data_base64": {
									Type:        schema.TypeString,
									Description: "NetworkDataBase64 contains config drive cloud-init networkdata as a base64 encoded string.",
									Optional:    true,
								},
								"network_data": {
									Type:        schema.TypeString,
									Description: "NetworkData contains config drive inline cloud-init networkdata.",
									Optional:    true,
								},
							},
						},
					},
					"service_account": {
						Type:        schema.TypeList,
						Description: "ServiceAccountVolumeSource represents a reference to a service account.",
						MaxItems:    1,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"service_account_name": {
									Type:        schema.TypeString,
									Description: "Name of the service account in the pod's namespace to use.",
									Required:    true,
								},
							},
						},
					},
					// TODO nargaman - Add other data volume source types
				},
			},
		},
	}
}

func volumesSchema() *schema.Schema {
	fields := volumesFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Specification of the desired behavior of the VirtualMachineInstance on the host."),
		Optional:    true,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandVolumes(volumes []interface{}) []kubevirtapiv1.Volume {
	result := make([]kubevirtapiv1.Volume, len(volumes))

	if len(volumes) == 0 || volumes[0] == nil {
		return result
	}

	for i, condition := range volumes {
		in := condition.(map[string]interface{})

		if v, ok := in["name"].(string); ok {
			result[i].Name = v
		}
		if v, ok := in["volume_source"].([]interface{}); ok {
			result[i].VolumeSource = expandVolumeSource(v)
		}
	}

	return result
}

func expandVolumeSource(volumeSource []interface{}) kubevirtapiv1.VolumeSource {
	result := kubevirtapiv1.VolumeSource{}

	if len(volumeSource) == 0 || volumeSource[0] == nil {
		return result
	}

	in := volumeSource[0].(map[string]interface{})

	if v, ok := in["data_volume"].([]interface{}); ok {
		result.DataVolume = expandDataVolume(v)
	}
	if v, ok := in["cloud_init_config_drive"].([]interface{}); ok {
		result.CloudInitConfigDrive = expandCloudInitConfigDrive(v)
	}
	if v, ok := in["service_account"].([]interface{}); ok {
		result.ServiceAccount = expandServiceAccount(v)
	}

	return result
}

func expandDataVolume(dataVolumeSource []interface{}) *kubevirtapiv1.DataVolumeSource {
	if len(dataVolumeSource) == 0 || dataVolumeSource[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.DataVolumeSource{}
	in := dataVolumeSource[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok {
		result.Name = v
	}

	return result
}

func expandCloudInitConfigDrive(cloudInitConfigDriveSource []interface{}) *kubevirtapiv1.CloudInitConfigDriveSource {
	if len(cloudInitConfigDriveSource) == 0 || cloudInitConfigDriveSource[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.CloudInitConfigDriveSource{}
	in := cloudInitConfigDriveSource[0].(map[string]interface{})

	if v, ok := in["user_data_secret_ref"].([]interface{}); ok {
		result.UserDataSecretRef = k8s.ExpandLocalObjectReferences(v)
	}
	if v, ok := in["user_data_base64"].(string); ok {
		result.UserDataBase64 = v
	}
	if v, ok := in["user_data"].(string); ok {
		result.UserData = v
	}
	if v, ok := in["network_data_secret_ref"].([]interface{}); ok {
		result.NetworkDataSecretRef = k8s.ExpandLocalObjectReferences(v)
	}
	if v, ok := in["network_data_base64"].(string); ok {
		result.NetworkDataBase64 = v
	}
	if v, ok := in["network_data"].(string); ok {
		result.NetworkData = v
	}

	return result
}

func expandServiceAccount(serviceAccountSource []interface{}) *kubevirtapiv1.ServiceAccountVolumeSource {
	if len(serviceAccountSource) == 0 || serviceAccountSource[0] == nil {
		return nil
	}

	result := &kubevirtapiv1.ServiceAccountVolumeSource{}
	in := serviceAccountSource[0].(map[string]interface{})

	if v, ok := in["service_account_name"].(string); ok {
		result.ServiceAccountName = v
	}

	return result
}

func flattenVolumes(in []kubevirtapiv1.Volume) []interface{} {
	att := make([]interface{}, len(in))

	for i, v := range in {
		c := make(map[string]interface{})

		c["name"] = v.Name
		c["volume_source"] = flattenVolumeSource(v.VolumeSource)

		att[i] = c
	}

	return att
}

func flattenVolumeSource(in kubevirtapiv1.VolumeSource) []interface{} {
	att := make(map[string]interface{})

	if in.DataVolume != nil {
		att["data_volume"] = flattenDataVolume(*in.DataVolume)
	}
	if in.CloudInitConfigDrive != nil {
		att["cloud_init_config_drive"] = flattenCloudInitConfigDrive(*in.CloudInitConfigDrive)
	}
	if in.ServiceAccount != nil {
		att["service_account"] = flattenServiceAccount(*in.ServiceAccount)
	}

	return []interface{}{att}
}

func flattenDataVolume(in kubevirtapiv1.DataVolumeSource) []interface{} {
	att := make(map[string]interface{})

	att["name"] = in.Name

	return []interface{}{att}
}

func flattenCloudInitConfigDrive(in kubevirtapiv1.CloudInitConfigDriveSource) []interface{} {
	att := make(map[string]interface{})

	if in.UserDataSecretRef != nil {
		att["user_data_secret_ref"] = k8s.FlattenLocalObjectReferences(*in.UserDataSecretRef)
	}
	att["user_data_base64"] = in.UserDataBase64
	att["user_data"] = in.UserData
	if in.NetworkDataSecretRef != nil {
		att["network_data_secret_ref"] = k8s.FlattenLocalObjectReferences(*in.NetworkDataSecretRef)
	}
	att["network_data_base64"] = in.NetworkDataBase64
	att["network_data"] = in.NetworkData

	return []interface{}{att}
}

func flattenServiceAccount(in kubevirtapiv1.ServiceAccountVolumeSource) []interface{} {
	att := make(map[string]interface{})

	att["service_account_name"] = in.ServiceAccountName

	return []interface{}{att}
}
