package equinixmetal

import (
	"encoding/json"

	equinixprovider "github.com/openshift/cluster-api-provider-equinix-metal/pkg/apis/equinixmetal/v1beta1"
)

type config struct {
	Roles        []string `json:"metal_roles,omitempty"`
	Facility     []string `json:"metal_facility,omitempty"`
	OS           string   `json:"metal_os"`
	ProjectID    string   `json:"metal_project_id"`
	BillingCycle string   `json:"metal_billing_cycle"`
	MachineType  string   `json:"metal_machine_type"`
	SshKeys      []string `json:"metal_ssh_keys,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*equinixprovider.EquinixMetalMachineProviderConfig
	APIURL              string
	APIKey              string
}

//TFVars generate EquinixMetal-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	plane0 := sources.ControlPlaneConfigs[0]

	/*
		 roles := make([]string, len(plane0.Roles))
		for _, r := range plane0.Roles {
			roles = append(roles, string(r))
		}
	*/
	// TODO(displague) fill in the tf vars
	cfg := &config{
		// Roles:        roles,
		Facility:     []string{plane0.Facility},
		OS:           plane0.OS,
		ProjectID:    plane0.ProjectID,
		BillingCycle: plane0.BillingCycle,
		MachineType:  plane0.MachineType,
		// SshKeys:      plane0.SshKeys,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
