// Package powervs contains Power Virtual Servers-specific Terraform-variable logic.
package powervs

import (
	"encoding/json"

	"github.com/openshift/cluster-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
)

// powervsRegionToVPCRegion based on:
// https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server&locale=en#creating-service
// https://github.com/ocp-power-automation/ocp4-upi-powervs/blob/master/docs/var.tfvars-doc.md
var powervsRegionToIBMRegion = map[string]string{
	"dal":     "us-south",
	"us-east": "us-east",
	"sao":     "br-sao",
	"tor":     "ca-tor",
	"mon":     "ca-mon",
	"eu-de-1": "eu-de",
	"eu-de-2": "eu-de",
	"lon":     "eu-gb",
	"syd":     "au-syd",
	"tok":     "jp-tok",
	"osa":     "jp-osa",
}

type config struct {
	APIKey               string `json:"powervs_api_key"`
	PowerVSRegion        string `json:"powervs_region"`
	VPCRegion            string `json:"powervs_vpc_region"`
	PowerVSResourceGroup string `json:"powervs_resource_group"`
	SSHKey               string `json:"powervs_ssh_key"`
	ImageID              string `json:"powervs_image_name"`
	NetworkIDs           string `json:"powervs_network_name"`
	BootstrapMemory      string `json:"powervs_bootstrap_memory"`
	BootstrapProcessors  string `json:"powervs_bootstrap_processors"`
	MasterMemory         string `json:"powervs_master_memory"`
	MasterProcessors     string `json:"powervs_master_processors"`
	ProcType             string `json:"powervs_proc_type"`
	SysType              string `json:"powervs_sys_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfigs []*v1alpha1.PowerVSMachineProviderConfig
	APIKey        string
	SSHKey        string
	PowerVSRegion string
}

// TFVars generates Power VS-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]

	//@TODO: Add resource group to platform
	//  -- change ImageID to ImageURL here?
	//  --
	cfg := &config{
		APIKey:               sources.APIKey,
		PowerVSRegion:        sources.PowerVSRegion,
		VPCRegion:            powervsRegionToIBMRegion[sources.PowerVSRegion],
		PowerVSResourceGroup: "powervs-ipi-resource-group",
		SSHKey:               sources.SSHKey,
		ImageID:              masterConfig.ImageID,
		NetworkIDs:           masterConfig.NetworkIDs[0],
		BootstrapMemory:      masterConfig.Memory,
		BootstrapProcessors:  masterConfig.Processors,
		MasterMemory:         masterConfig.Memory,
		MasterProcessors:     masterConfig.Processors,
		ProcType:             masterConfig.ProcType,
		SysType:              masterConfig.SysType,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
