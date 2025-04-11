package ibmcloud

import (
	"fmt"
	"strings"

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

	capiEndpoints := make([]string, 0)
	for _, endpoint := range m.ServiceEndpoints {
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
			// https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/blob/91d63f492c4b9b16a67b0312be26325056953111/pkg/endpoints/endpoints.go#L48
			// NOTE(cjschaef): IAM is not supported as an option for CAPI's endpoint flag (argument), it must be passed in as an environment variable instead.
			continue
		}

		capiEndpoints = append(capiEndpoints, fmt.Sprintf("%s=%s", serviceName, endpoint.URL))
	}

	// If no IBM Cloud CAPI endpoints exist, nothing should be returned for the flag.
	if len(capiEndpoints) == 0 {
		return ""
	}

	// IBM Cloud CAPI expects endpoint flag formatted as:
	// "region":"service-name1"="url","service-name2"="url",...
	return fmt.Sprintf("%s:%s", m.Region, strings.Join(capiEndpoints, ","))
}
