package powervs

import configv1 "github.com/openshift/api/config/v1"

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	BaseDomain           string                            `json:"BaseDomain"`
	CISInstanceCRN       string                            `json:"cisInstanceCRN"`
	DNSInstanceCRN       string                            `json:"dnsInstanceCRN"`
	PowerVSResourceGroup string                            `json:"powerVSResourceGroup"`
	Region               string                            `json:"region"`
	VPCRegion            string                            `json:"vpcRegion"`
	Zone                 string                            `json:"zone"`
	ServiceInstanceGUID  string                            `json:"serviceInstanceGUID"`
	ServiceEndpoints     []configv1.PowerVSServiceEndpoint `json:"serviceEndpoints,omitempty"`
	TransitGateway       string                            `json:"transitGatewayName"`
	VPC                  string                            `json:"vpcName"`
}
