// Package aws contains AWS-specific Terraform-variable logic.
package aws

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws/defaults"
)

type config struct {
	AMI                     string            `json:"aws_ami"`
	ExtraTags               map[string]string `json:"aws_extra_tags,omitempty"`
	BootstrapInstanceType   string            `json:"aws_bootstrap_instance_type,omitempty"`
	MasterInstanceType      string            `json:"aws_master_instance_type,omitempty"`
	MasterAvailabilityZones []string          `json:"aws_master_availability_zones"`
	WorkerAvailabilityZones []string          `json:"aws_worker_availability_zones"`
	IOPS                    int64             `json:"aws_master_root_volume_iops"`
	Size                    int64             `json:"aws_master_root_volume_size,omitempty"`
	Type                    string            `json:"aws_master_root_volume_type,omitempty"`
	Encrypted               bool              `json:"aws_master_root_volume_encrypted"`
	KMSKeyID                string            `json:"aws_master_root_volume_kms_key_id,omitempty"`
	Region                  string            `json:"aws_region,omitempty"`
	VPC                     string            `json:"aws_vpc,omitempty"`
	PrivateSubnets          []string          `json:"aws_private_subnets,omitempty"`
	PublicSubnets           *[]string         `json:"aws_public_subnets,omitempty"`
	PublishStrategy         string            `json:"aws_publish_strategy,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	VPC                           string
	PrivateSubnets, PublicSubnets []string

	Publish types.PublishingStrategy

	MasterConfigs, WorkerConfigs []*v1beta1.AWSMachineProviderConfig
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Name] = tag.Value
	}

	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Placement.AvailabilityZone
	}

	exists := struct{}{}
	availabilityZoneMap := map[string]struct{}{}
	for _, c := range sources.WorkerConfigs {
		availabilityZoneMap[c.Placement.AvailabilityZone] = exists
	}
	workerAvailabilityZones := make([]string, 0, len(availabilityZoneMap))
	for zone := range availabilityZoneMap {
		workerAvailabilityZones = append(workerAvailabilityZones, zone)
	}

	if len(masterConfig.BlockDevices) == 0 {
		return nil, errors.New("block device slice cannot be empty")
	}

	rootVolume := masterConfig.BlockDevices[0]
	if rootVolume.EBS == nil {
		return nil, errors.New("EBS information must be configured for the root volume")
	}

	if rootVolume.EBS.VolumeType == nil {
		return nil, errors.New("EBS volume type must be configured for the root volume")
	}

	if rootVolume.EBS.VolumeSize == nil {
		return nil, errors.New("EBS volume size must be configured for the root volume")
	}

	if *rootVolume.EBS.VolumeType == "io1" && rootVolume.EBS.Iops == nil {
		return nil, errors.New("EBS IOPS must be configured for the io1 root volume")
	}

	instanceClass := defaults.InstanceClass(masterConfig.Placement.Region)

	cfg := &config{
		Region:                  masterConfig.Placement.Region,
		ExtraTags:               tags,
		AMI:                     *masterConfig.AMI.ID,
		MasterAvailabilityZones: masterAvailabilityZones,
		WorkerAvailabilityZones: workerAvailabilityZones,
		BootstrapInstanceType:   fmt.Sprintf("%s.large", instanceClass),
		MasterInstanceType:      masterConfig.InstanceType,
		Size:                    *rootVolume.EBS.VolumeSize,
		Type:                    *rootVolume.EBS.VolumeType,
		VPC:                     sources.VPC,
		PrivateSubnets:          sources.PrivateSubnets,
		PublishStrategy:         string(sources.Publish),
	}

	if len(sources.PublicSubnets) == 0 {
		if cfg.VPC != "" {
			cfg.PublicSubnets = &[]string{}
		}
	} else {
		cfg.PublicSubnets = &sources.PublicSubnets
	}

	if rootVolume.EBS.Iops != nil {
		cfg.IOPS = *rootVolume.EBS.Iops
	}

	cfg.Encrypted = true
	if rootVolume.EBS.Encrypted != nil {
		cfg.Encrypted = *rootVolume.EBS.Encrypted
	}
	if rootVolume.EBS.KMSKey.ID != nil && *rootVolume.EBS.KMSKey.ID != "" {
		cfg.KMSKeyID = *rootVolume.EBS.KMSKey.ID
	} else if rootVolume.EBS.KMSKey.ARN != nil && *rootVolume.EBS.KMSKey.ARN != "" {
		cfg.KMSKeyID = *rootVolume.EBS.KMSKey.ARN
	}

	return json.MarshalIndent(cfg, "", "  ")
}
