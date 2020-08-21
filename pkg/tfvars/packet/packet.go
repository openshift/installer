package packet

import (
	"encoding/json"

	packetprovider "github.com/packethost/cluster-api-provider-packet/pkg/apis/packetprovider/v1alpha1"
)

type config struct {
	Roles        []string `json:"packet_roles,omitempty"`
	Facility     []string `json:"packet_facility,omitempty"`
	OS           string   `json:"packet_os"`
	ProjectID    string   `json:"packet_project_id"`
	BillingCycle string   `json:"packet_billing_cycle"`
	MachineType  string   `json:"packet_machine_type"`
	SshKeys      []string `json:"packet_ssh_keys,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*packetprovider.PacketMachineProviderSpec
	APIURL              string
	APIKey              string
}

//TFVars generate Packet-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	plane0 := sources.ControlPlaneConfigs[0]

	roles := make([]string, len(plane0.Roles))
	for _, r := range plane0.Roles {
		roles = append(roles, string(r))
	}
	// TODO(displague) fill in the tf vars
	cfg := &config{
		Roles:        roles,
		Facility:     plane0.Facility,
		OS:           plane0.OS,
		ProjectID:    plane0.ProjectID,
		BillingCycle: plane0.BillingCycle,
		MachineType:  plane0.MachineType,
		SshKeys:      plane0.SshKeys,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
