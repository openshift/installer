package ibmcloud

import configv1 "github.com/openshift/api/config/v1"

// Metadata contains IBM Cloud metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	AccountID         string                             `json:"accountID"`
	BaseDomain        string                             `json:"baseDomain"`
	CISInstanceCRN    string                             `json:"cisInstanceCRN,omitempty"`
	DNSInstanceID     string                             `json:"dnsInstanceID,omitempty"`
	Region            string                             `json:"region,omitempty"`
	ResourceGroupName string                             `json:"resourceGroupName,omitempty"`
	ServiceEndpoints  []configv1.IBMCloudServiceEndpoint `json:"serviceEndpoints,omitempty"`
	Subnets           []string                           `json:"subnets,omitempty"`
	VPC               string                             `json:"vpc,omitempty"`
}
