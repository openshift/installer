package ibmcloud

import (
	"encoding/json"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	APIKey string `json:"ibmcloud_api_key,omitempty"`
}

type config struct {
	Auth                    `json:",inline"`
	Region                  string   `json:"ibmcloud_region,omitempty"`
	CISInstanceCRN          string   `json:"ibmcloud_cis_crn,omitempty"`
	MasterAvailabilityZones []string `json:"ibmcloud_master_availability_zones"`
	ResourceGroupName       string   `json:"ibmcloud_resource_group_name,omitempty"`
	VSIImage                string   `json:"ibmcloud_vsi_image,omitempty"`

	// TODO: IBM[#98]: Set different instance types
	// BootstrapInstanceType   string   `json:"ibmcloud_bootstrap_instance_type,omitempty"`
	// MasterInstanceType      string   `json:"ibmcloud_master_instance_type,omitempty"`

	// TODO: IBM[#100]: Support publish strategy modes
	// PublishStrategy         string   `json:"ibmcloud_publish_strategy,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                    Auth
	Region                  string
	CISInstanceCRN          string
	ResourceGroupName       string
	VSIImage                string
	MasterAvailabilityZones []string

	// TODO: IBM: Future support
	// MasterConfigs      []*ibmcloudprovider.ibmcloudMachineProviderSpec
	// WorkerConfigs      []*ibmcloudprovider.ibmcloudMachineProviderSpec
	// PublishStrategy types.PublishingStrategy
}

// TFVars generates ibmcloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	// TODO: IBM: Future support
	// masterConfig := sources.MasterConfigs[0]
	// workerConfig := sources.WorkerConfigs[0]
	// masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	// for i, c := range sources.MasterConfigs {
	// 	masterAvailabilityZones[i] = c.Zone
	// }

	// TODO: IBM: Fetch most of the config from masterConfig instead

	cfg := &config{
		Auth:                    sources.Auth,
		Region:                  sources.Region,
		CISInstanceCRN:          sources.CISInstanceCRN,
		MasterAvailabilityZones: sources.MasterAvailabilityZones,
		ResourceGroupName:       sources.ResourceGroupName,
		VSIImage:                sources.VSIImage,

		// TODO: IBM: Future support
		// Region                   masterConfig.Region,
		// BootstrapInstanceType:   masterConfig.MachineType,
		// MasterInstanceType:      masterConfig.MachineType,
		// MasterAvailabilityZones: masterAvailabilityZones,
		// PublishStrategy:         string(sources.PublishStrategy),
	}

	return json.MarshalIndent(cfg, "", "  ")
}
