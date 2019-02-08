// Package aws contains AWS-specific Terraform-variable logic.
package aws

import (
	"encoding/json"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
)

type config struct {
	EC2AMIOverride string            `json:"aws_ec2_ami_override,omitempty"`
	ExtraTags      map[string]string `json:"aws_extra_tags,omitempty"`
	EC2Type        string            `json:"aws_master_ec2_type,omitempty"`
	IOPS           int64             `json:"aws_master_root_volume_iops,omitempty"`
	Size           int64             `json:"aws_master_root_volume_size,omitempty"`
	Type           string            `json:"aws_master_root_volume_type,omitempty"`
	Region         string            `json:"aws_region,omitempty"`
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(masterConfig *v1beta1.AWSMachineProviderConfig) ([]byte, error) {
	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Name] = tag.Value
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

	cfg := &config{
		Region:         masterConfig.Placement.Region,
		ExtraTags:      tags,
		EC2AMIOverride: *masterConfig.AMI.ID,
		EC2Type:        masterConfig.InstanceType,
		Size:           *rootVolume.EBS.VolumeSize,
		Type:           *rootVolume.EBS.VolumeType,
	}

	if rootVolume.EBS.Iops != nil {
		cfg.IOPS = *rootVolume.EBS.Iops
	}

	return json.MarshalIndent(cfg, "", "  ")
}
