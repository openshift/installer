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
	Region             string            `json:"region_id"`
	ZoneIds            []string          `json:"zone_ids"`
	ResourceGroupId    string            `json:"resource_group_id"`
	VpcCidrBlock       string            `json:"vpc_cidr_block"`
	VSwitchCidrBlocks  []string          `json:"vswitch_cidr_blocks"`
	InstanceType       string            `json:"instance_type"`
	ImageId            string            `json:"image_id"`
	SystemDiskSize     string            `json:"system_disk_size"`
	SystemDiskCategory string            `json:"system_disk_category"`
	DataDiskSize       string            `json:"data_disk_size"`
	DataDiskCategory   string            `json:"data_disk_category"`
	KeyName            string            `json:"key_name"`
	Tags               map[string]string `json:"tags"`
	IgnitionBucket     string            `json:"ignition_bucket,omitempty"`
	IgnitionStub       string            `json:"ignition_stub,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth              Auth
	EnvironmentName   string
	ResourceGroupId   string
	VpcCidrBlock      string
	VSwitchCidrBlocks []string
	Publish           types.PublishingStrategy
	BaseDomain        string
	MasterConfigs     []*AlibabacloudMachineProviderSpec
	WorkerConfigs     []*AlibabacloudMachineProviderSpec
	IgnitionBucket    string
	IgnitionFile      string
	IgnitionStub      string
	ImageId           string
}

// TFVars generates AlibabaCloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	workerConfig := sources.WorkerConfigs[0]

	zone_ids := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		zone_ids[i] = c.ZoneId
	}

	cfg := &config{
		Auth:               sources.Auth,
		Region:             masterConfig.RegionId,
		ZoneIds:            zone_ids,
		ResourceGroupId:    masterConfig.ResourceGroupId,
		VpcCidrBlock:       sources.VpcCidrBlock,
		VSwitchCidrBlocks:  sources.VSwitchCidrBlocks,
		InstanceType:       masterConfig.InstanceType,
		ImageId:            masterConfig.ImageId,
		SystemDiskSize:     masterConfig.SystemDiskSize,
		SystemDiskCategory: masterConfig.SystemDiskCategory,
		DataDiskSize:       masterConfig.DataDisk[0].Size,
		DataDiskCategory:   masterConfig.DataDisk[0].Category,
		KeyName:            workerConfig.KeyPairName,
		IgnitionBucket:     sources.IgnitionBucket,
		IgnitionStub:       sources.IgnitionStub,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
