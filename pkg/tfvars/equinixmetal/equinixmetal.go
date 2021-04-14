package equinixmetal

import (
	"encoding/json"

	equinixprovider "github.com/openshift/cluster-api-provider-equinix-metal/pkg/apis/equinixmetal/v1beta1"
)

type Auth struct {
	APIURL string `json:"metal_api_url"`
	APIKey string `json:"metal_auth_token"`
}
type config struct {
	Auth          `json:",inline"`
	Roles         []string `json:"metal_roles,omitempty"`
	Facility      string   `json:"metal_facility,omitempty"`
	Metro         string   `json:"metal_metro,omitempty"`
	OS            string   `json:"metal_os"`
	ProjectID     string   `json:"metal_project_id"`
	BillingCycle  string   `json:"metal_billing_cycle"`
	MachineType   string   `json:"metal_machine_type"`
	SshKeys       []string `json:"metal_ssh_keys,omitempty"`
	IPXEScriptURL string   `json:"metal_ipxe_script_url,omitempty"`
	CustomData    string   `json:"metal_custom_data"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*equinixprovider.EquinixMetalMachineProviderConfig
	Auth                Auth
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
		Auth: sources.Auth,
		// Roles:        roles,
		Facility: plane0.Facility,
		// Metro:         plane0.Metro,
		OS:            plane0.OS,
		ProjectID:     plane0.ProjectID,
		BillingCycle:  plane0.BillingCycle,
		MachineType:   plane0.MachineType,
		IPXEScriptURL: plane0.IPXEScriptURL,

		// SshKeys:      plane0.SshKeys,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
