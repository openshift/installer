// Package ovirt contains ovirt-specific Terraform-variable logic.
package ovirt

import (
	"encoding/json"
	"fmt"

	"github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types/ovirt"
)

// Auth is the collection of credentials that will be used by terraform.
type Auth struct {
	URL      string `json:"ovirt_url"`
	Username string `json:"ovirt_username"`
	Password string `json:"ovirt_password"`
	Cafile   string `json:"ovirt_cafile"`
	Cabundle string `json:"ovirt_ca_bundle"`
	Insecure bool   `json:"ovirt_insecure"`
}

type config struct {
	Auth                    `json:",inline"`
	ClusterID               string                   `json:"ovirt_cluster_id"`
	StorageDomainID         string                   `json:"ovirt_storage_domain_id"`
	NetworkName             string                   `json:"ovirt_network_name,omitempty"`
	VNICProfileID           string                   `json:"ovirt_vnic_profile_id,omitempty"`
	AffinityGroups          []map[string]interface{} `json:"ovirt_affinity_groups,omitempty"`
	BaseImageName           string                   `json:"openstack_base_image_name,omitempty"`
	BaseImageLocalFilePath  string                   `json:"openstack_base_image_local_file_path,omitempty"`
	MasterInstanceTypeID    string                   `json:"ovirt_master_instance_type_id"`
	MasterVMType            string                   `json:"ovirt_master_vm_type,omitempty"`
	MasterMemory            int32                    `json:"ovirt_master_memory"`
	MasterCores             int32                    `json:"ovirt_master_cores"`
	MasterSockets           int32                    `json:"ovirt_master_sockets"`
	MasterOsDiskGB          int64                    `json:"ovirt_master_os_disk_gb"`
	MasterAffinityGroups    []string                 `json:"ovirt_master_affinity_groups"`
	MasterAutoPinningPolicy string                   `json:"ovirt_master_auto_pinning_policy,omitempty"`
	MasterHugePages         int32                    `json:"ovirt_master_hugepages"`
}

// TFVars generates ovirt-specific Terraform variables.
func TFVars(
	auth Auth,
	clusterID string,
	storageDomainID string,
	networkName string,
	vnicProfileID string,
	baseImage string,
	infraID string,
	masterSpec *v1beta1.OvirtMachineProviderSpec,
	affinityGroups []ovirt.AffinityGroup,
) ([]byte, error) {
	cfg := config{
		Auth:                    auth,
		ClusterID:               clusterID,
		StorageDomainID:         storageDomainID,
		NetworkName:             networkName,
		VNICProfileID:           vnicProfileID,
		BaseImageName:           baseImage,
		MasterInstanceTypeID:    masterSpec.InstanceTypeId,
		MasterVMType:            masterSpec.VMType,
		MasterOsDiskGB:          masterSpec.OSDisk.SizeGB,
		MasterMemory:            masterSpec.MemoryMB,
		MasterAffinityGroups:    masterSpec.AffinityGroupsNames,
		MasterAutoPinningPolicy: masterSpec.AutoPinningPolicy,
		MasterHugePages:         masterSpec.Hugepages,
	}
	if masterSpec.CPU != nil {
		cfg.MasterCores = masterSpec.CPU.Cores
		cfg.MasterSockets = masterSpec.CPU.Sockets
	}

	imageName, isURL := rhcos.GenerateOpenStackImageName(baseImage, infraID)
	cfg.BaseImageName = imageName
	if isURL {
		imageFilePath, err := cache.DownloadImageFile(baseImage)
		if err != nil {
			return nil, err
		}
		cfg.BaseImageLocalFilePath = imageFilePath
	}
	if len(affinityGroups) > 0 {
		cfg.AffinityGroups = handleAffinityGroups(affinityGroups, infraID)
	}
	return json.MarshalIndent(cfg, "", "  ")
}

func handleAffinityGroups(ags []ovirt.AffinityGroup, infraID string) []map[string]interface{} {
	tfAffinityGroups := make([]map[string]interface{}, len(ags))
	for i, ag := range ags {
		tfAffinityGroups[i] = map[string]interface{}{
			"name":        fmt.Sprintf("%s-%s", infraID, ag.Name),
			"priority":    ag.Priority,
			"description": ag.Description,
			"enforcing":   ag.Enforcing,
		}
	}
	return tfAffinityGroups
}
