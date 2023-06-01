// Package aws contains AWS-specific Terraform-variable logic.
package aws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/types"
	typesaws "github.com/openshift/installer/pkg/types/aws"
)

type config struct {
	AMI                             string            `json:"aws_ami"`
	AMIRegion                       string            `json:"aws_ami_region"`
	CustomEndpoints                 map[string]string `json:"custom_endpoints,omitempty"`
	ExtraTags                       map[string]string `json:"aws_extra_tags,omitempty"`
	BootstrapInstanceType           string            `json:"aws_bootstrap_instance_type,omitempty"`
	MasterInstanceType              string            `json:"aws_master_instance_type,omitempty"`
	MasterAvailabilityZones         []string          `json:"aws_master_availability_zones"`
	WorkerAvailabilityZones         []string          `json:"aws_worker_availability_zones"`
	IOPS                            int64             `json:"aws_master_root_volume_iops"`
	Size                            int64             `json:"aws_master_root_volume_size,omitempty"`
	Type                            string            `json:"aws_master_root_volume_type,omitempty"`
	Encrypted                       bool              `json:"aws_master_root_volume_encrypted"`
	KMSKeyID                        string            `json:"aws_master_root_volume_kms_key_id,omitempty"`
	Region                          string            `json:"aws_region,omitempty"`
	VPC                             string            `json:"aws_vpc,omitempty"`
	PrivateSubnets                  []string          `json:"aws_private_subnets,omitempty"`
	PublicSubnets                   *[]string         `json:"aws_public_subnets,omitempty"`
	InternalZone                    string            `json:"aws_internal_zone,omitempty"`
	InternalZoneRole                string            `json:"aws_internal_zone_role,omitempty"`
	PublishStrategy                 string            `json:"aws_publish_strategy,omitempty"`
	IgnitionBucket                  string            `json:"aws_ignition_bucket"`
	BootstrapIgnitionStub           string            `json:"aws_bootstrap_stub_ignition"`
	MasterIAMRoleName               string            `json:"aws_master_iam_role_name,omitempty"`
	WorkerIAMRoleName               string            `json:"aws_worker_iam_role_name,omitempty"`
	MasterMetadataAuthentication    string            `json:"aws_master_instance_metadata_authentication,omitempty"`
	BootstrapMetadataAuthentication string            `json:"aws_bootstrap_instance_metadata_authentication,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	VPC                            string
	PrivateSubnets, PublicSubnets  []string
	InternalZone, InternalZoneRole string
	Services                       []typesaws.ServiceEndpoint

	Publish types.PublishingStrategy

	AMIID, AMIRegion string

	MasterConfigs, WorkerConfigs []*machinev1beta1.AWSMachineProviderConfig

	IgnitionBucket, IgnitionPresignedURL string

	AdditionalTrustBundle string

	MasterIAMRoleName, WorkerIAMRoleName string

	MasterMetadataAuthentication string

	Architecture types.Architecture

	Proxy *types.Proxy
}

// TFVars generates AWS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	endpoints := make(map[string]string)
	for _, service := range sources.Services {
		service := service
		endpoints[service.Name] = service.URL
	}

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

	cfg := &config{
		CustomEndpoints:         endpoints,
		Region:                  masterConfig.Placement.Region,
		ExtraTags:               tags,
		MasterAvailabilityZones: masterAvailabilityZones,
		WorkerAvailabilityZones: workerAvailabilityZones,
		BootstrapInstanceType:   masterConfig.InstanceType,
		MasterInstanceType:      masterConfig.InstanceType,
		Size:                    *rootVolume.EBS.VolumeSize,
		Type:                    *rootVolume.EBS.VolumeType,
		VPC:                     sources.VPC,
		PrivateSubnets:          sources.PrivateSubnets,
		InternalZone:            sources.InternalZone,
		InternalZoneRole:        sources.InternalZoneRole,
		PublishStrategy:         string(sources.Publish),
		IgnitionBucket:          sources.IgnitionBucket,
		MasterIAMRoleName:       sources.MasterIAMRoleName,
		WorkerIAMRoleName:       sources.WorkerIAMRoleName,
	}

	stubIgn, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(sources.IgnitionPresignedURL, sources.AdditionalTrustBundle, sources.Proxy)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stub Ignition config for bootstrap")
	}

	// Check the size of the raw ignition stub is less than 16KB for aws user-data
	// see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html
	if len(stubIgn) > 16000 {
		return nil, fmt.Errorf("rendered bootstrap ignition shim exceeds the 16KB limit for AWS user data -- try reducing the size of your CA cert bundle")
	}
	cfg.BootstrapIgnitionStub = string(stubIgn)

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

	if masterConfig.AMI.ID != nil && *masterConfig.AMI.ID != "" {
		cfg.AMI = *masterConfig.AMI.ID
		cfg.AMIRegion = masterConfig.Placement.Region
	} else {
		cfg.AMI = sources.AMIID
		cfg.AMIRegion = sources.AMIRegion
	}

	if masterConfig.MetadataServiceOptions.Authentication != "" {
		cfg.MasterMetadataAuthentication = strings.ToLower(string(masterConfig.MetadataServiceOptions.Authentication))
		cfg.BootstrapMetadataAuthentication = cfg.MasterMetadataAuthentication
	}

	return json.MarshalIndent(cfg, "", "  ")
}
