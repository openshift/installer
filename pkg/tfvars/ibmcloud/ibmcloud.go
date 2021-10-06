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

type config struct {
	Auth                    `json:",inline"`
	Region                  string   `json:"ibmcloud_region,omitempty"`
	BootstrapInstanceType   string   `json:"ibmcloud_bootstrap_instance_type,omitempty"`
	CISInstanceCRN          string   `json:"ibmcloud_cis_crn,omitempty"`
	ExtraTags               []string `json:"ibmcloud_extra_tags,omitempty"`
	MasterAvailabilityZones []string `json:"ibmcloud_master_availability_zones"`
	MasterInstanceType      string   `json:"ibmcloud_master_instance_type,omitempty"`
	PublishStrategy         string   `json:"ibmcloud_publish_strategy,omitempty"`
	ResourceGroupName       string   `json:"ibmcloud_resource_group_name,omitempty"`
	ImageFilePath           string   `json:"ibmcloud_image_filepath,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth              Auth
	CISInstanceCRN    string
	ImageURL          string
	MasterConfigs     []*ibmcloudprovider.IBMCloudMachineProviderSpec
	PublishStrategy   types.PublishingStrategy
	ResourceGroupName string
	WorkerConfigs     []*ibmcloudprovider.IBMCloudMachineProviderSpec
}

// TFVars generates ibmcloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached ibmcloud image")
	}

	masterConfig := sources.MasterConfigs[0]
	// workerConfig := sources.WorkerConfigs[0]
	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Zone
	}

	cfg := &config{
		Auth:                    sources.Auth,
		BootstrapInstanceType:   masterConfig.Profile,
		CISInstanceCRN:          sources.CISInstanceCRN,
		ImageFilePath:           cachedImage,
		MasterAvailabilityZones: masterAvailabilityZones,
		MasterInstanceType:      masterConfig.Profile,
		PublishStrategy:         string(sources.PublishStrategy),
		Region:                  masterConfig.Region,
		ResourceGroupName:       sources.ResourceGroupName,

		// TODO: IBM: Future support
		// ExtraTags:               masterConfig.Tags,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
