package gcp

import (
	"encoding/json"

	gcpprovider "github.com/openshift/cluster-api-provider-gcp/pkg/apis/gcpprovider/v1beta1"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	ProjectID      string `json:"gcp_project_id,omitempty"`
	ServiceAccount string `json:"gcp_service_account,ompitempty"`
}

type config struct {
	Auth                    `json:",inline"`
	Region                  string   `json:"gcp_region,omitempty"`
	BootstrapInstanceType   string   `json:"gcp_bootstrap_instance_type,omitempty"`
	MasterInstanceType      string   `json:"gcp_master_instance_type,omitempty"`
	MasterAvailabilityZones []string `json:"gcp_master_availability_zones"`
	ImageURI                string   `json:"gcp_image_uri,omitempty"`
	VolumeType              string   `json:"gcp_master_root_volume_type"`
	VolumeSize              int64    `json:"gcp_master_root_volume_size"`
	PublicZoneName          string   `json:"gcp_public_dns_zone_name,omitempty"`
}

// TFVars generates gcp-specific Terraform variables launching the cluster.
func TFVars(auth Auth, masterConfigs []*gcpprovider.GCPMachineProviderSpec, imageURI, publicZoneName string) ([]byte, error) {
	masterConfig := masterConfigs[0]
	masterAvailabilityZones := make([]string, len(masterConfigs))
	for i, c := range masterConfigs {
		masterAvailabilityZones[i] = c.Zone
	}
	cfg := &config{
		Auth:                    auth,
		Region:                  masterConfig.Region,
		BootstrapInstanceType:   masterConfig.MachineType,
		MasterInstanceType:      masterConfig.MachineType,
		MasterAvailabilityZones: masterAvailabilityZones,
		VolumeType:              masterConfig.Disks[0].Type,
		VolumeSize:              masterConfig.Disks[0].SizeGb,
		ImageURI:                imageURI,
		PublicZoneName:          publicZoneName,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
