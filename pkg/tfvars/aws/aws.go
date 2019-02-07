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
	EC2Type        string            `json:"aws_control_plane_ec2_type,omitempty"`
	IOPS           int64             `json:"aws_control_plane_root_volume_iops"`
	Size           int64             `json:"aws_control_plane_root_volume_size,omitempty"`
	Type           string            `json:"aws_control_plane_root_volume_type,omitempty"`
	Region         string            `json:"aws_region,omitempty"`
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(controlPlaneConfig *v1beta1.AWSMachineProviderConfig) ([]byte, error) {
	tags := make(map[string]string, len(controlPlaneConfig.Tags))
	for _, tag := range controlPlaneConfig.Tags {
		tags[tag.Name] = tag.Value
	}

	if len(controlPlaneConfig.BlockDevices) == 0 {
		return nil, errors.New("block device slice cannot be empty")
	}

	rootVolume := controlPlaneConfig.BlockDevices[0]
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
		Region:         controlPlaneConfig.Placement.Region,
		ExtraTags:      tags,
		EC2AMIOverride: *controlPlaneConfig.AMI.ID,
		EC2Type:        controlPlaneConfig.InstanceType,
		Size:           *rootVolume.EBS.VolumeSize,
		Type:           *rootVolume.EBS.VolumeType,
	}

	if rootVolume.EBS.Iops != nil {
		cfg.IOPS = *rootVolume.EBS.Iops
	}

	return json.MarshalIndent(cfg, "", "  ")
}
