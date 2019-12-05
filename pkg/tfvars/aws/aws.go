// Package aws contains AWS-specific Terraform-variable logic.
package aws

import (
	"encoding/json"
	"fmt"
	"os"

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
	Region                  string            `json:"aws_region,omitempty"`
	VPC                     string            `json:"aws_vpc,omitempty"`
	PrivateSubnets          []string          `json:"aws_private_subnets,omitempty"`
	PublicSubnets           *[]string         `json:"aws_public_subnets,omitempty"`
	PublishStrategy         string            `json:"aws_publish_strategy,omitempty"`
	UseIPv6                 bool              `json:"aws_use_ipv6,omitempty"`
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(vpc string, privateSubnets []string, publicSubnets []string, publish types.PublishingStrategy, masterConfigs []*v1beta1.AWSMachineProviderConfig, workerConfigs []*v1beta1.AWSMachineProviderConfig) ([]byte, error) {
	masterConfig := masterConfigs[0]

	tags := make(map[string]string, len(masterConfig.Tags))
	for _, tag := range masterConfig.Tags {
		tags[tag.Name] = tag.Value
	}

	masterAvailabilityZones := make([]string, len(masterConfigs))
	for i, c := range masterConfigs {
		masterAvailabilityZones[i] = c.Placement.AvailabilityZone
	}

	exists := struct{}{}
	availabilityZoneMap := map[string]struct{}{}
	for _, c := range workerConfigs {
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

	useIPv6 := false
	if os.Getenv("OPENSHIFT_AWS_USE_IPV6") == "true" {
		useIPv6 = true
	}

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
		VPC:                     vpc,
		PrivateSubnets:          privateSubnets,
		PublishStrategy:         string(publish),
		UseIPv6:                 useIPv6,
	}

	if len(publicSubnets) == 0 {
		if cfg.VPC != "" {
			cfg.PublicSubnets = &[]string{}
		}
	} else {
		cfg.PublicSubnets = &publicSubnets
	}

	if rootVolume.EBS.Iops != nil {
		cfg.IOPS = *rootVolume.EBS.Iops
	}

	return json.MarshalIndent(cfg, "", "  ")
}
