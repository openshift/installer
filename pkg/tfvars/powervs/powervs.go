// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types"
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
	DNSInstanceGUID      string `json:"powervs_dns_guid"`
	ImageBucketName      string `json:"powervs_image_bucket_name"`
	ImageBucketFileName  string `json:"powervs_image_bucket_file_name"`
	NetworkName          string `json:"powervs_network_name"`
	VPCName              string `json:"powervs_vpc_name"`
	VPCSubnetName        string `json:"powervs_vpc_subnet_name"`
	VPCPermitted         bool   `json:"powervs_vpc_permitted"`
	VPCGatewayName       string `json:"powervs_vpc_gateway_name"`
	VPCGatewayAttached   bool   `json:"powervs_vpc_gateway_attached"`
	CloudConnectionName  string `json:"powervs_ccon_name"`
	BootstrapMemory      int32  `json:"powervs_bootstrap_memory"`
	BootstrapProcessors  string `json:"powervs_bootstrap_processors"`
	MasterMemory         int32  `json:"powervs_master_memory"`
	MasterProcessors     string `json:"powervs_master_processors"`
	ProcType             string `json:"powervs_proc_type"`
	SysType              string `json:"powervs_sys_type"`
	PublishStrategy      string `json:"powervs_publish_strategy"`
	EnableSNAT           bool   `json:"powervs_enable_snat"`
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
	CloudConnectionName  string
	CISInstanceCRN       string
	DNSInstanceCRN       string
	VPCName              string
	VPCSubnetName        string
	VPCPermitted         bool
	VPCGatewayName       string
	VPCGatewayAttached   bool
	PublishStrategy      types.PublishingStrategy
	EnableSNAT           bool
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

	var serviceInstanceID, processor, dnsGUID string
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

	// Parse GUID from DNS CRN
	if sources.DNSInstanceCRN != "" {
		dnsCRN, err := crn.Parse(sources.DNSInstanceCRN)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DNSInstanceCRN")
		}
		dnsGUID = dnsCRN.ServiceInstance
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
		DNSInstanceGUID:      dnsGUID,
		ImageBucketName:      sources.ImageBucketName,
		ImageBucketFileName:  sources.ImageBucketFileName,
		VPCName:              sources.VPCName,
		VPCSubnetName:        sources.VPCSubnetName,
		VPCPermitted:         sources.VPCPermitted,
		VPCGatewayName:       sources.VPCGatewayName,
		VPCGatewayAttached:   sources.VPCGatewayAttached,
		CloudConnectionName:  sources.CloudConnectionName,
		BootstrapMemory:      masterConfig.MemoryGiB,
		BootstrapProcessors:  processor,
		MasterMemory:         masterConfig.MemoryGiB,
		MasterProcessors:     processor,
		ProcType:             strings.ToLower(string(masterConfig.ProcessorType)),
		SysType:              masterConfig.SystemType,
		PublishStrategy:      string(sources.PublishStrategy),
		EnableSNAT:           sources.EnableSNAT,
	}
	if masterConfig.Network.Name != nil {
		cfg.NetworkName = *masterConfig.Network.Name
	}

	return json.MarshalIndent(cfg, "", "  ")
}
