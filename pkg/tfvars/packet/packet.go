package packet

import "encoding/json"

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
}

//TFVars generate Packet-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	// TODO(displague) fill in the tf vars
	return json.MarshalIndent(struct{}{}, "", "  ")
}
