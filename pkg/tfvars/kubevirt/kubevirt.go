// Package kubevirt contains kubevirt-specific Terraform-variable logic.
package kubevirt

import (
	"encoding/json"
	v1 "github.com/openshift/cluster-api-provider-kubevirt/pkg/apis/kubevirtprovider/v1alpha1"
)

type config struct {
	Namespace                  string            `json:"kubevirt_namespace"`
	ImageURL                   string            `json:"kubevirt_image_url"`
	SourcePVCName              string            `json:"kubevirt_source_pvc_name"`
	Memory                     string            `json:"kubevirt_master_memory"`
	CPU                        uint32            `json:"kubevirt_master_cpu"`
	Storage                    string            `json:"kubevirt_master_storage"`
	StorageClass               string            `json:"kubevirt_storage_class"`
	NetworkName                string            `json:"kubevirt_network_name"`
	InterfaceBindingMethod     string            `json:"kubevirt_interface_binding_method"`
	PersistentVolumeAccessMode string            `json:"kubevirt_pv_access_mode"`
	ResourcesLabels            map[string]string `json:"kubevirt_labels"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterSpecs     []*v1.KubevirtMachineProviderSpec
	ImageURL        string
	Namespace       string
	ResourcesLabels map[string]string
}

// TFVars generates kubevirt-specific Terraform variables.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterSpec := sources.MasterSpecs[0]

	// For optional parameters, set only if not nil
	cfg := config{
		Namespace:                  sources.Namespace,
		ImageURL:                   sources.ImageURL,
		SourcePVCName:              masterSpec.SourcePvcName,
		Memory:                     masterSpec.RequestedMemory,
		CPU:                        masterSpec.RequestedCPU,
		Storage:                    masterSpec.RequestedStorage,
		StorageClass:               masterSpec.StorageClassName,
		NetworkName:                masterSpec.NetworkName,
		InterfaceBindingMethod:     masterSpec.InterfaceBindingMethod,
		PersistentVolumeAccessMode: masterSpec.PersistentVolumeAccessMode,
		ResourcesLabels:            sources.ResourcesLabels,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
