package alibabacloud

import (
	"encoding/json"
	"fmt"

	machinev1 "github.com/openshift/api/machine/v1"
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
	Auth                      `json:",inline"`
	Region                    string            `json:"ali_region_id"`
	VpcID                     string            `json:"ali_vpc_id"`
	VSwitchIDs                []string          `json:"ali_vswitch_ids"`
	MasterAvailabilityZoneIDs []string          `json:"ali_master_availability_zone_ids"`
	WorkerAvailabilityZoneIDs []string          `json:"ali_worker_availability_zone_ids"`
	PrivateZoneID             string            `json:"ali_private_zone_id"`
	NatGatewayZoneID          string            `json:"ali_nat_gateway_zone_id"`
	ResourceGroupID           string            `json:"ali_resource_group_id"`
	BootstrapInstanceType     string            `json:"ali_bootstrap_instance_type"`
	MasterInstanceType        string            `json:"ali_master_instance_type"`
	ImageID                   string            `json:"ali_image_id"`
	SystemDiskSize            int               `json:"ali_system_disk_size"`
	SystemDiskCategory        string            `json:"ali_system_disk_category"`
	ExtraTags                 map[string]string `json:"ali_extra_tags"`
	IgnitionBucket            string            `json:"ali_ignition_bucket"`
	BootstrapIgnitionStub     string            `json:"ali_bootstrap_stub_ignition"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                  Auth
	VpcID                 string
	VSwitchIDs            []string
	PrivateZoneID         string
	ResourceGroupID       string
	BaseDomain            string
	NatGatewayZoneID      string
	MasterConfigs         []*machinev1.AlibabaCloudMachineProviderConfig
	WorkerConfigs         []*machinev1.AlibabaCloudMachineProviderConfig
	IgnitionBucket        string
	IgnitionPresignedURL  string
	Publish               types.PublishingStrategy
	AdditionalTrustBundle string
	Architecture          types.Architecture
}

// TFVars generates AlibabaCloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	masterAvailabilityZoneIDs := make([]string, len(sources.MasterConfigs))
	for i, masterConfig := range sources.MasterConfigs {
		masterAvailabilityZoneIDs[i] = masterConfig.ZoneID
	}

	workerAvailabilityZones := map[string]bool{}
	for _, workerConfig := range sources.WorkerConfigs {
		if !workerAvailabilityZones[workerConfig.ZoneID] {
			workerAvailabilityZones[workerConfig.ZoneID] = true
		}
	}
	workerAvailabilityZoneIDs := make([]string, 0, len(workerAvailabilityZones))
	for zoneID := range workerAvailabilityZones {
		workerAvailabilityZoneIDs = append(workerAvailabilityZoneIDs, zoneID)
	}

	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Key] = tag.Value
	}

	cfg := &config{
		Auth:                      sources.Auth,
		Region:                    masterConfig.RegionID,
		MasterAvailabilityZoneIDs: masterAvailabilityZoneIDs,
		WorkerAvailabilityZoneIDs: workerAvailabilityZoneIDs,
		VpcID:                     sources.VpcID,
		VSwitchIDs:                sources.VSwitchIDs,
		PrivateZoneID:             sources.PrivateZoneID,
		NatGatewayZoneID:          sources.NatGatewayZoneID,
		ResourceGroupID:           sources.ResourceGroupID,
		BootstrapInstanceType:     masterConfig.InstanceType,
		MasterInstanceType:        masterConfig.InstanceType,
		ImageID:                   masterConfig.ImageID,
		SystemDiskSize:            int(masterConfig.SystemDisk.Size),
		SystemDiskCategory:        masterConfig.SystemDisk.Category,
		ExtraTags:                 tags,
		IgnitionBucket:            sources.IgnitionBucket,
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
