package ibmcloud

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
)

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

// GetRegionAndEndpointsFlag will return the IBM Cloud region and any service endpoint overrides formatted as the IBM Cloud CAPI command line argument.
func (m *Metadata) GetRegionAndEndpointsFlag() string {
	// If there are no endpoints, return an empty string (rather than just the region).
	if len(m.ServiceEndpoints) == 0 {
		return ""
	}

	flag := m.Region
	for index, endpoint := range m.ServiceEndpoints {
		// IBM Cloud CAPI has pre-defined endpoint service names that do not follow naming scheme, use those instead.
		// TODO(cjschaef): See about opening a CAPI GH issue to link here for this restriction.
		var serviceName string
		switch endpoint.Name {
		case configv1.IBMCloudServiceCOS:
			serviceName = "cos"
		case configv1.IBMCloudServiceGlobalTagging:
			serviceName = "globaltagging"
		case configv1.IBMCloudServiceResourceController:
			serviceName = "rc"
		case configv1.IBMCloudServiceResourceManager:
			serviceName = "rm"
		case configv1.IBMCloudServiceVPC:
			serviceName = "vpc"
		default:
			// Any additional Service Endpoint overrides should be ignored, as they are not supported by CAPI.
			// NOTE(cjschaef): IAM is not supported as an option for CAPI's endpoint flag (argument), it must be passed in as an environment variable instead.
			continue
		}

		// Format for first (and perhaps only) endpoint is unique, remaining are similar
		if index == 0 {
			flag = fmt.Sprintf("%s:%s=%s", flag, serviceName, endpoint.URL)
		} else {
			flag = fmt.Sprintf("%s,%s=%s", flag, serviceName, endpoint.URL)
		}
	}
	return flag
}
