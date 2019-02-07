// Package tfvars generates Terraform variables for launching the cluster.
package tfvars

import (
	"encoding/json"
	"net"
)

type config struct {
	ClusterID                   string `json:"cluster_id,omitempty"`
	Name                        string `json:"cluster_name,omitempty"`
	BaseDomain                  string `json:"base_domain,omitempty"`
	MachineCIDR                 string `json:"machine_cidr"`
	ControlPlaneCount           int    `json:"control_plane_count,omitempty"`
	ControlPlaneMachinePoolName string `json:"control_plane_machine_pool_name"`

	IgnitionBootstrap    string `json:"ignition_bootstrap,omitempty"`
	IgnitionControlPlane string `json:"ignition_control_plane,omitempty"`
}

// TFVars generates terraform.tfvar JSON for launching the cluster.
func TFVars(clusterID string, clusterName string, baseDomain string, machineCIDR *net.IPNet, bootstrapIgn string, controlPlaneIgn string, controlPlaneCount int, controlPlaneMachinePoolName string) ([]byte, error) {
	config := &config{
		ClusterID:                   clusterID,
		Name:                        clusterName,
		BaseDomain:                  baseDomain,
		MachineCIDR:                 machineCIDR.String(),
		ControlPlaneCount:           controlPlaneCount,
		ControlPlaneMachinePoolName: controlPlaneMachinePoolName,
		IgnitionBootstrap:           bootstrapIgn,
		IgnitionControlPlane:        controlPlaneIgn,
	}

	return json.MarshalIndent(config, "", "  ")
}
