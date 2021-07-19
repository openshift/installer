package ibmcloud

// Metadata contains IBM Cloud metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	AccountID         string   `json:"accountID"`
	BaseDomain        string   `json:"baseDomain"`
	CISInstanceCRN    string   `json:"cisInstanceCRN"`
	Region            string   `json:"region,omitempty"`
	ResourceGroupName string   `json:"resourceGroupName,omitempty"`
	VPC               string   `json:"vpc,omitempty"`
	Subnets           []string `json:"subnets,omitempty"`
}
