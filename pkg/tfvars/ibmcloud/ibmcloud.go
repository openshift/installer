package ibmcloud

import (
	"encoding/json"

	ibmcloudprovider "github.com/openshift/cluster-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1beta1"

	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	APIKey string `json:"ibmcloud_api_key,omitempty"`
}

// DedicatedHost is the format used by terraform.
type DedicatedHost struct {
	ID      string `json:"id,omitempty"`
	Profile string `json:"profile,omitempty"`
}

type config struct {
	Auth                    `json:",inline"`
	Region                  string          `json:"ibmcloud_region,omitempty"`
	BootstrapInstanceType   string          `json:"ibmcloud_bootstrap_instance_type,omitempty"`
	CISInstanceCRN          string          `json:"ibmcloud_cis_crn,omitempty"`
	ExtraTags               []string        `json:"ibmcloud_extra_tags,omitempty"`
	MasterAvailabilityZones []string        `json:"ibmcloud_master_availability_zones"`
	WorkerAvailabilityZones []string        `json:"ibmcloud_worker_availability_zones"`
	MasterInstanceType      string          `json:"ibmcloud_master_instance_type,omitempty"`
	MasterDedicatedHosts    []DedicatedHost `json:"ibmcloud_master_dedicated_hosts,omitempty"`
	WorkerDedicatedHosts    []DedicatedHost `json:"ibmcloud_worker_dedicated_hosts,omitempty"`
	PublishStrategy         string          `json:"ibmcloud_publish_strategy,omitempty"`
	ResourceGroupName       string          `json:"ibmcloud_resource_group_name,omitempty"`
	ImageFilePath           string          `json:"ibmcloud_image_filepath,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                 Auth
	CISInstanceCRN       string
	ImageURL             string
	MasterConfigs        []*ibmcloudprovider.IBMCloudMachineProviderSpec
	MasterDedicatedHosts []DedicatedHost
	PublishStrategy      types.PublishingStrategy
	ResourceGroupName    string
	WorkerConfigs        []*ibmcloudprovider.IBMCloudMachineProviderSpec
	WorkerDedicatedHosts []DedicatedHost
}

// TFVars generates ibmcloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached ibmcloud image")
	}

	masterConfig := sources.MasterConfigs[0]
	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Zone
	}
	workerAvailabilityZones := make([]string, len(sources.WorkerConfigs))
	for i, c := range sources.WorkerConfigs {
		workerAvailabilityZones[i] = c.Zone
	}

	cfg := &config{
		Auth:                    sources.Auth,
		BootstrapInstanceType:   masterConfig.Profile,
		CISInstanceCRN:          sources.CISInstanceCRN,
		ImageFilePath:           cachedImage,
		MasterAvailabilityZones: masterAvailabilityZones,
		MasterDedicatedHosts:    sources.MasterDedicatedHosts,
		MasterInstanceType:      masterConfig.Profile,
		PublishStrategy:         string(sources.PublishStrategy),
		Region:                  masterConfig.Region,
		ResourceGroupName:       sources.ResourceGroupName,
		WorkerAvailabilityZones: workerAvailabilityZones,
		WorkerDedicatedHosts:    sources.WorkerDedicatedHosts,

		// TODO: IBM: Future support
		// ExtraTags:               masterConfig.Tags,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
