// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/crn"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

type config struct {
	APIKey                 string `json:"powervs_api_key"`
	SSHKey                 string `json:"powervs_ssh_key"`
	PowerVSRegion          string `json:"powervs_region"`
	PowerVSZone            string `json:"powervs_zone"`
	VPCRegion              string `json:"powervs_vpc_region"`
	VPCZone                string `json:"powervs_vpc_zone"`
	COSRegion              string `json:"powervs_cos_region"`
	PowerVSResourceGroup   string `json:"powervs_resource_group"`
	CISInstanceCRN         string `json:"powervs_cis_crn"`
	DNSInstanceGUID        string `json:"powervs_dns_guid"`
	ImageBucketName        string `json:"powervs_image_bucket_name"`
	ImageBucketFileName    string `json:"powervs_image_bucket_file_name"`
	VPCName                string `json:"powervs_vpc_name"`
	VPCSubnetName          string `json:"powervs_vpc_subnet_name"`
	VPCPermitted           bool   `json:"powervs_vpc_permitted"`
	VPCGatewayName         string `json:"powervs_vpc_gateway_name"`
	VPCGatewayAttached     bool   `json:"powervs_vpc_gateway_attached"`
	BootstrapMemory        int32  `json:"powervs_bootstrap_memory"`
	BootstrapProcessors    string `json:"powervs_bootstrap_processors"`
	MasterMemory           int32  `json:"powervs_master_memory"`
	MasterProcessors       string `json:"powervs_master_processors"`
	ProcType               string `json:"powervs_proc_type"`
	SysType                string `json:"powervs_sys_type"`
	PublishStrategy        string `json:"powervs_publish_strategy"`
	EnableSNAT             bool   `json:"powervs_enable_snat"`
	AttachedTransitGateway string `json:"powervs_attached_transit_gateway"`
	TGConnectionVPCID      string `json:"powervs_tg_connection_vpc_id"`
	ServiceInstanceName    string `json:"powervs_service_instance_name"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfigs          []*machinev1.PowerVSMachineProviderConfig
	APIKey                 string
	SSHKey                 string
	Region                 string
	Zone                   string
	ImageBucketName        string
	ImageBucketFileName    string
	PowerVSResourceGroup   string
	CISInstanceCRN         string
	DNSInstanceCRN         string
	VPCRegion              string
	VPCZone                string
	VPCName                string
	VPCSubnetName          string
	VPCPermitted           bool
	VPCGatewayName         string
	VPCGatewayAttached     bool
	PublishStrategy        types.PublishingStrategy
	EnableSNAT             bool
	AttachedTransitGateway string
	TGConnectionVPCID      string
	ServiceInstanceName    string
}

// TFVars generates Power VS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	cosRegion, err := powervstypes.COSRegionForVPCRegion(sources.VPCRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to find COS region for VPC region")
	}

	var processor, dnsGUID string

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
		APIKey:                 sources.APIKey,
		SSHKey:                 sources.SSHKey,
		PowerVSRegion:          sources.Region,
		PowerVSZone:            sources.Zone,
		VPCRegion:              sources.VPCRegion,
		VPCZone:                sources.VPCZone,
		COSRegion:              cosRegion,
		PowerVSResourceGroup:   sources.PowerVSResourceGroup,
		CISInstanceCRN:         sources.CISInstanceCRN,
		DNSInstanceGUID:        dnsGUID,
		ImageBucketName:        sources.ImageBucketName,
		ImageBucketFileName:    sources.ImageBucketFileName,
		VPCName:                sources.VPCName,
		VPCSubnetName:          sources.VPCSubnetName,
		VPCPermitted:           sources.VPCPermitted,
		VPCGatewayName:         sources.VPCGatewayName,
		VPCGatewayAttached:     sources.VPCGatewayAttached,
		BootstrapMemory:        masterConfig.MemoryGiB,
		BootstrapProcessors:    processor,
		MasterMemory:           masterConfig.MemoryGiB,
		MasterProcessors:       processor,
		ProcType:               strings.ToLower(string(masterConfig.ProcessorType)),
		SysType:                masterConfig.SystemType,
		PublishStrategy:        string(sources.PublishStrategy),
		EnableSNAT:             sources.EnableSNAT,
		AttachedTransitGateway: sources.AttachedTransitGateway,
		TGConnectionVPCID:      sources.TGConnectionVPCID,
		ServiceInstanceName:    sources.ServiceInstanceName,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
