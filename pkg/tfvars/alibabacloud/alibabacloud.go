package alibabacloud

import (
	"encoding/json"
	"fmt"

	alibabacloudprovider "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
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
	Region                string            `json:"ali_region_id"`
	ZoneIDs               []string          `json:"ali_zone_ids"`
	NatGatewayZoneID      string            `json:"ali_nat_gateway_zone_id"`
	ResourceGroupID       string            `json:"ali_resource_group_id"`
	BootstrapInstanceType string            `json:"ali_bootstrap_instance_type"`
	MasterInstanceType    string            `json:"ali_master_instance_type"`
	ImageID               string            `json:"ali_image_id"`
	SystemDiskSize        int               `json:"ali_system_disk_size"`
	SystemDiskCategory    string            `json:"ali_system_disk_category"`
	ExtraTags             map[string]string `json:"ali_extra_tags,omitempty"`
	IgnitionBucket        string            `json:"ali_ignition_bucket"`
	BootstrapIgnitionStub string            `json:"ali_bootstrap_stub_ignition"`
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
	Publish               types.PublishingStrategy
	AdditionalTrustBundle string
	Architecture          types.Architecture
}

// TFVars generates AlibabaCloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	zoneIDs := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		zoneIDs[i] = c.ZoneID
	}

	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Key] = tag.Value
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
		ExtraTags:             tags,
		IgnitionBucket:        sources.IgnitionBucket,
	}

	stubIgn, err := bootstrap.GenerateIgnitionShimWithCertBundle(sources.IgnitionPresignedURL, sources.AdditionalTrustBundle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stub Ignition config for bootstrap")
	}

	// Check the size of the raw ignition stub is less than 16KB for user-data
	if len(stubIgn) > 16000 {
		return nil, fmt.Errorf("rendered bootstrap ignition shim exceeds the 16KB limit for Alibaba Cloud user data -- try reducing the size of your CA cert bundle")
	}
	cfg.BootstrapIgnitionStub = string(stubIgn)

	return json.MarshalIndent(cfg, "", "  ")
}
