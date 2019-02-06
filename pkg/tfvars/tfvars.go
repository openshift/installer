// Package tfvars generates Terraform variables for launching the cluster.
package tfvars

import (
	"encoding/json"
	"net"
)

type config struct {
	ClusterID             string `json:"cluster_id,omitempty"`
	Name                  string `json:"cluster_name,omitempty"`
	BaseDomain            string `json:"base_domain,omitempty"`
	MachineCIDR           string `json:"machine_cidr"`
	Masters               int    `json:"master_count,omitempty"`
	MasterMachinePoolName string `json:"master_machine_pool_name"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`
}

// TFVars generates terraform.tfvar JSON for launching the cluster.
func TFVars(clusterID string, clusterName string, baseDomain string, machineCIDR *net.IPNet, bootstrapIgn string, masterIgn string, masterCount int, masterMachinePoolName string) ([]byte, error) {
	config := &config{
		ClusterID:             clusterID,
		Name:                  clusterName,
		BaseDomain:            baseDomain,
		MachineCIDR:           machineCIDR.String(),
		Masters:               masterCount,
		MasterMachinePoolName: masterMachinePoolName,
		IgnitionBootstrap:     bootstrapIgn,
		IgnitionMaster:        masterIgn,
	}

	return json.MarshalIndent(config, "", "  ")
}
