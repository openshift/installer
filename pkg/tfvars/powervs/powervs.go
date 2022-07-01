// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types/powervs"
)

type config struct {
	ServiceInstanceID    string `json:"powervs_cloud_instance_id"`
	APIKey               string `json:"powervs_api_key"`
	SSHKey               string `json:"powervs_ssh_key"`
	PowerVSRegion        string `json:"powervs_region"`
	PowerVSZone          string `json:"powervs_zone"`
	VPCRegion            string `json:"powervs_vpc_region"`
	VPCZone              string `json:"powervs_vpc_zone"`
	PowerVSResourceGroup string `json:"powervs_resource_group"`
	CISInstanceCRN       string `json:"powervs_cis_crn"`
	ImageBucketName      string `json:"powervs_image_bucket_name"`
	ImageBucketFileName  string `json:"powervs_image_bucket_file_name"`
	NetworkName          string `json:"powervs_network_name"`
	VPCName              string `json:"powervs_vpc_name"`
	VPCSubnetName        string `json:"powervs_vpc_subnet_name"`
	BootstrapMemory      int32  `json:"powervs_bootstrap_memory"`
	BootstrapProcessors  string `json:"powervs_bootstrap_processors"`
	MasterMemory         int32  `json:"powervs_master_memory"`
	MasterProcessors     string `json:"powervs_master_processors"`
	ProcType             string `json:"powervs_proc_type"`
	SysType              string `json:"powervs_sys_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfigs        []*machinev1.PowerVSMachineProviderConfig
	APIKey               string
	SSHKey               string
	Region               string
	Zone                 string
	ImageBucketName      string
	ImageBucketFileName  string
	NetworkName          string
	PowerVSResourceGroup string
	CISInstanceCRN       string
	VPCName              string
	VPCSubnetName        string
}

// TFVars generates Power VS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	// TODO(mjturek): Allow user to specify vpcRegion in install config like we're doing for vpcZone
	vpcRegion := powervs.Regions[sources.Region].VPCRegion

	// Randomly select a zone in the VPC region.
	// @TODO: Align this with a region later.
	rand.Seed(time.Now().UnixNano())
	// All supported Regions are MZRs and have Zones named "region-[1-3]"
	vpcZone := fmt.Sprintf("%s-%d", vpcRegion, rand.Intn(2)+1)

	var serviceInstanceID, processor string
	if masterConfig.ServiceInstance.ID != nil {
		serviceInstanceID = *masterConfig.ServiceInstance.ID
	} else {
		return nil, fmt.Errorf("serviceInstanceID is nil")
	}

	if masterConfig.Processors.StrVal != "" {
		processor = masterConfig.Processors.StrVal
	} else {
		processor = fmt.Sprintf("%d", masterConfig.Processors.IntVal)
	}

	cfg := &config{
		ServiceInstanceID:    serviceInstanceID,
		APIKey:               sources.APIKey,
		SSHKey:               sources.SSHKey,
		PowerVSRegion:        sources.Region,
		PowerVSZone:          sources.Zone,
		VPCRegion:            vpcRegion,
		VPCZone:              vpcZone,
		PowerVSResourceGroup: sources.PowerVSResourceGroup,
		CISInstanceCRN:       sources.CISInstanceCRN,
		ImageBucketName:      sources.ImageBucketName,
		ImageBucketFileName:  sources.ImageBucketFileName,
		VPCName:              sources.VPCName,
		VPCSubnetName:        sources.VPCSubnetName,
		BootstrapMemory:      masterConfig.MemoryGiB,
		BootstrapProcessors:  processor,
		MasterMemory:         masterConfig.MemoryGiB,
		MasterProcessors:     processor,
		ProcType:             strings.ToLower(string(masterConfig.ProcessorType)),
		SysType:              masterConfig.SystemType,
	}
	if masterConfig.Network.Name != nil {
		cfg.NetworkName = *masterConfig.Network.Name
	}

	return json.MarshalIndent(cfg, "", "  ")
}
