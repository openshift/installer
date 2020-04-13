// Package tfvars generates Terraform variables for launching the cluster.
package tfvars

import (
	"encoding/json"
	"net"
	"strings"
)

type config struct {
	ClusterID     string `json:"cluster_id,omitempty"`
	ClusterDomain string `json:"cluster_domain,omitempty"`
	BaseDomain    string `json:"base_domain,omitempty"`
	// DeprecatedMachineCIDR has been replaced with machine_v4_cidrs, use the first
	// entry from there instead.
	DeprecatedMachineCIDR string   `json:"machine_cidr"`
	Masters               int      `json:"master_count,omitempty"`
	MachineV4CIDRs        []string `json:"machine_v4_cidrs"`
	MachineV6CIDRs        []string `json:"machine_v6_cidrs"`

	UseIPv4 bool `json:"use_ipv4"`
	UseIPv6 bool `json:"use_ipv6"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`
}

// TFVars generates terraform.tfvar JSON for launching the cluster.
func TFVars(clusterID string, clusterDomain string, baseDomain string, deprecatedMachineCIDR *net.IPNet, machineV4CIDRs []string, machineV6CIDRs []string, useIPv4, useIPv6 bool, bootstrapIgn string, masterIgn string, masterCount int) ([]byte, error) {
	config := &config{
		ClusterID:             clusterID,
		ClusterDomain:         strings.TrimSuffix(clusterDomain, "."),
		BaseDomain:            strings.TrimSuffix(baseDomain, "."),
		DeprecatedMachineCIDR: deprecatedMachineCIDR.String(),
		MachineV4CIDRs:        machineV4CIDRs,
		MachineV6CIDRs:        machineV6CIDRs,
		UseIPv4:               useIPv4,
		UseIPv6:               useIPv6,
		Masters:               masterCount,
		IgnitionBootstrap:     bootstrapIgn,
		IgnitionMaster:        masterIgn,
	}

	return json.MarshalIndent(config, "", "  ")
}
