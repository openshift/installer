package azure

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/types/azure/defaults"
)

type config struct {
	ExtraTags             map[string]string `json:"azure_extra_tags,omitempty"`
	BootstrapInstanceType string            `json:"azure_bootstrap_vm_type,omitempty"`
	MasterInstanceType    string            `json:"azure_master_vm_type,omitempty"`
	Size                  int64             `json:"azure_master_root_volume_size,omitempty"`
	Region                string            `json:"azure_region,omitempty"`
	ClusterName           string            `json:"cluster_name,omitempty"`
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(clustername string) ([]byte, error) {
	var region = "eastus"
	var volumeSize int64
	volumeSize = 30
	cfg := &config{
		Region:                region,
		ExtraTags:             map[string]string{},
		BootstrapInstanceType: defaults.InstanceClass(region),
		MasterInstanceType:    defaults.InstanceClass(region),
		Size:                  volumeSize,
		ClusterName:           clustername,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
