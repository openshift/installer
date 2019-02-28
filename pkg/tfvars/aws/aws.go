// Package aws contains AWS-specific Terraform-variable logic.
package aws

import (
	"encoding/json"
	"fmt"

	"github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1beta1"
)

type config struct {
	AMI                   string            `json:"aws_ami"`
	ExtraTags             map[string]string `json:"aws_extra_tags,omitempty"`
	BootstrapInstanceType string            `json:"aws_bootstrap_instance_type,omitempty"`
	MasterInstanceType    string            `json:"aws_master_instance_type,omitempty"`
	AvailabilityZones     []string          `json:"aws_master_availability_zones"`
	IOPS                  int64             `json:"aws_master_root_volume_iops"`
	Size                  int64             `json:"aws_master_root_volume_size,omitempty"`
	Type                  string            `json:"aws_master_root_volume_type,omitempty"`
	Region                string            `json:"aws_region,omitempty"`
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(masterConfigs []*v1beta1.AWSMachineProviderConfig) ([]byte, error) {
	masterConfig := masterConfigs[0]

	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Name] = tag.Value
	}

	availabilityZones := make([]string, len(masterConfigs))
	for i, c := range masterConfigs {
		availabilityZones[i] = c.Placement.AvailabilityZone
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
		Region:                masterConfig.Placement.Region,
		ExtraTags:             tags,
		AMI:                   *masterConfig.AMI.ID,
		AvailabilityZones:     availabilityZones,
		BootstrapInstanceType: fmt.Sprintf("%s.large", instanceClass),
		MasterInstanceType:    masterConfig.InstanceType,
		Size:                  *rootVolume.EBS.VolumeSize,
		Type:                  *rootVolume.EBS.VolumeType,
	}

	if rootVolume.EBS.Iops != nil {
		cfg.IOPS = *rootVolume.EBS.Iops
	}

	return json.MarshalIndent(cfg, "", "  ")
}
