package alibabacloud

import (
	"encoding/json"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}

type config struct {
	Auth        `json:",inline"`
	Region      string `json:"region,omitempty"`
	ZoneId      string `json:"zone_id,omitempty"`
	VpcName     string `json:"vpc_name,omitempty"`
	VSwitchName string `json:"vswitch_name,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth Auth
}

// TFVars generates AlibabaCloud-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	cfg := &config{
		Auth: sources.Auth,
		// Test
		Region:      "cn-hangzhou",
		ZoneId:      "cn-hangzhou-i",
		VpcName:     "sh-test-vpc",
		VSwitchName: "sh-test-vs",
	}

	return json.MarshalIndent(cfg, "", "  ")
}
