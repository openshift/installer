package alibabacloud

import (
	"encoding/json"

	// TODO Alibaba: In the future, will use this repo: github.com/openshift/cluster-api-provider-alibaba
	alibabacloudprovider "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	AccessKey string `json:"ali_access_key"`
	SecretKey string `json:"ali_secret_key"`
}

type config struct {
	Auth                  `json:",inline"`
	Region                string                     `json:"ali_region_id"`
	ZoneIDs               []string                   `json:"ali_zone_ids"`
	NatGatewayZoneID      string                     `json:"ali_nat_gateway_zone_id"`
	ResourceGroupID       string                     `json:"ali_resource_group_id"`
	BootstrapInstanceType string                     `json:"ali_bootstrap_instance_type"`
	MasterInstanceType    string                     `json:"ali_master_instance_type"`
	ImageID               string                     `json:"ali_image_id"`
	SystemDiskSize        int                        `json:"ali_system_disk_size"`
	SystemDiskCategory    string                     `json:"ali_system_disk_category"`
	Tags                  []alibabacloudprovider.Tag `json:"tag,omitempty"`
	IgnitionBucket        string                     `json:"ali_ignition_bucket"`
	BootstrapIgnitionStub string                     `json:"ali_bootstrap_stub_ignition"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                  Auth
	ResourceGroupID       string
	BaseDomain            string
	NatGatewayZoneID      string
	MasterConfigs         []*alibabacloudprovider.AlibabaCloudMachineProviderConfig
	WorkerConfigs         []*alibabacloudprovider.AlibabaCloudMachineProviderConfig
	IgnitionBucket        string
	IgnitionPresignedURL  string
	IgnitionFile          string
	ImageID               string
	SSHKey                string
	Publish               types.PublishingStrategy
	AdditionalTrustBundle string
	Architecture          types.Architecture
}

// TFVars generates AlibabaCloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	// workerConfig := sources.WorkerConfigs[0]

	zoneIDs := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		zoneIDs[i] = c.ZoneID
	}

	cfg := &config{
		Auth:                  sources.Auth,
		Region:                masterConfig.RegionID,
		ZoneIDs:               zoneIDs,
		NatGatewayZoneID:      sources.NatGatewayZoneID,
		ResourceGroupID:       sources.ResourceGroupID,
		BootstrapInstanceType: masterConfig.InstanceType,
		MasterInstanceType:    masterConfig.InstanceType,
		ImageID:               masterConfig.ImageID,
		SystemDiskSize:        masterConfig.SystemDiskSize,
		SystemDiskCategory:    masterConfig.SystemDiskCategory,
		Tags:                  masterConfig.Tags,
		IgnitionBucket:        sources.IgnitionBucket,
	}

	stubIgn, err := generateIgnitionShim(sources.IgnitionPresignedURL, sources.AdditionalTrustBundle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stub Ignition config for bootstrap")
	}
	cfg.BootstrapIgnitionStub = stubIgn

	return json.MarshalIndent(cfg, "", "  ")
}
