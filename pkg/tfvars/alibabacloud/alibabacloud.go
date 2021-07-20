package alibabacloud

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/types"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type config struct {
	Auth               `json:",inline"`
	Region             string         `json:"region_id"`
	ZoneIDs            []string       `json:"zone_ids"`
	ResourceGroupID    string         `json:"resource_group_id"`
	InstanceType       string         `json:"instance_type"`
	ImageID            string         `json:"image_id"`
	SystemDiskSize     string         `json:"system_disk_size"`
	SystemDiskCategory string         `json:"system_disk_category"`
	KeyName            string         `json:"key_name"`
	Tags               []*InstanceTag `json:"resource_tags"`
	IgnitionBucket     string         `json:"ignition_bucket,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                  Auth
	ResourceGroupID       string
	BaseDomain            string
	MasterConfigs         []*MachineProviderSpec
	WorkerConfigs         []*MachineProviderSpec
	IgnitionBucket        string
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
	workerConfig := sources.WorkerConfigs[0]

	zoneIDs := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		zoneIDs[i] = c.ZoneID
	}

	cfg := &config{
		Auth:               sources.Auth,
		Region:             masterConfig.RegionID,
		ZoneIDs:            zoneIDs,
		ResourceGroupID:    sources.ResourceGroupID,
		InstanceType:       masterConfig.InstanceType,
		ImageID:            masterConfig.ImageID,
		SystemDiskSize:     masterConfig.SystemDiskSize,
		SystemDiskCategory: masterConfig.SystemDiskCategory,
		KeyName:            workerConfig.KeyPairName,
		Tags:               masterConfig.Tag,
		IgnitionBucket:     sources.IgnitionBucket,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
