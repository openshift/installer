package ibmcloud

import (
	configv1 "github.com/openshift/api/config/v1"
)

const (
	// IBM Cloud Service Endpoint variables are supplied via documentation:
	// https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/guides/custom-service-endpoints#supported-endpoint-customizations

	// IBMCloudServiceCISVar is the variable name used by the IBM Cloud Terraform Provider to override the CIS endpoint.
	IBMCloudServiceCISVar string = "IBMCLOUD_CIS_API_ENDPOINT"

	// IBMCloudServiceCOSVar is the variable name used by the IBM Cloud Terraform Provider to override the COS endpoint.
	IBMCloudServiceCOSVar string = "IBMCLOUD_COS_ENDPOINT"

	// IBMCloudServiceCOSConfigVar is the variable name used by IBM Cloud Terraform Provider to override the COS Config endpoint.
	IBMCloudServiceCOSConfigVar string = "IBMCLOUD_COS_CONFIG_ENDPOINT"

	// IBMCloudServiceDNSServicesVar is the variable name used by the IBM Cloud Terraform Provider to override the DNS Services endpoint.
	IBMCloudServiceDNSServicesVar string = "IBMCLOUD_PRIVATE_DNS_API_ENDPOINT"

	// IBMCloudServiceGlobalCatalogVar is the variable name used by the IBM Cloud Terraform Provider to override the Global Catalog endpoint.
	IBMCloudServiceGlobalCatalogVar string = "IBMCLOUD_RESOURCE_CATALOG_API_ENDPOINT"

	// IBMCloudServiceGlobalSearchVar is the variable name used by the IBM Cloud Terraform Provider to override the Global Search endpoint.
	IBMCloudServiceGlobalSearchVar string = "IBMCLOUD_GS_API_ENDPOINT" //nolint:gosec // not hardcoded creds

	// IBMCloudServiceGlobalTaggingVar is the variable name used by the IBM Cloud Terraform Provider to override the Global Tagging endpoint.
	IBMCloudServiceGlobalTaggingVar string = "IBMCLOUD_GT_API_ENDPOINT" //nolint:gosec // not hardcoded creds

	// IBMCloudServiceHyperProtectVar is the variable name used by the IBM Cloud Terraform Provider to override the Hyper Protect endpoint.
	IBMCloudServiceHyperProtectVar string = "IBMCLOUD_HPCS_API_ENDPOINT"

	// IBMCloudServiceIAMVar is the variable name used by the IBM Cloud Terraform Provider to override the IAM endpoint.
	IBMCloudServiceIAMVar string = "IBMCLOUD_IAM_API_ENDPOINT"

	// IBMCloudServiceKeyProtectVar is the variable name used by the IBM Cloud Terraform Provider to override the Key Protect endpoint.
	IBMCloudServiceKeyProtectVar string = "IBMCLOUD_KP_API_ENDPOINT" //nolint:gosec // not hardcoded creds

	// IBMCloudServiceResourceControllerVar is the variable name used by the IBM Cloud Terraform Provider to override the Resource Controller endpoint.
	IBMCloudServiceResourceControllerVar string = "IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT"

	// IBMCloudServiceResourceManagerVar is the variable name used by the IBM Cloud Terraform Provider to override the Resource Manager endpoint.
	IBMCloudServiceResourceManagerVar string = "IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT"

	// IBMCloudServiceVPCVar is the variable name used by the IBM Cloud Terraform Provider to override the VPC endpoint.
	IBMCloudServiceVPCVar string = "IBMCLOUD_IS_NG_API_ENDPOINT"

	// IBMCloudInfrastructureSecurityGroupTargetReference is the name of a generic IBM Cloud Infrastructure Security Group target.
	IBMCloudInfrastructureSecurityGroupTargetReference string = "*vpcv1.SecurityGroupTargetReference"
)

var (
	// IBMCloudServiceOverrides is a set of IBM Cloud services allowed to have their endpoints overridden mapped to their override variable.
	IBMCloudServiceOverrides = map[configv1.IBMCloudServiceName]string{
		configv1.IBMCloudServiceCIS:                IBMCloudServiceCISVar,
		configv1.IBMCloudServiceCOS:                IBMCloudServiceCOSVar,
		configv1.IBMCloudServiceCOSConfig:          IBMCloudServiceCOSConfigVar,
		configv1.IBMCloudServiceDNSServices:        IBMCloudServiceDNSServicesVar,
		configv1.IBMCloudServiceGlobalCatalog:      IBMCloudServiceGlobalCatalogVar,
		configv1.IBMCloudServiceGlobalSearch:       IBMCloudServiceGlobalSearchVar,
		configv1.IBMCloudServiceGlobalTagging:      IBMCloudServiceGlobalTaggingVar,
		configv1.IBMCloudServiceHyperProtect:       IBMCloudServiceHyperProtectVar,
		configv1.IBMCloudServiceIAM:                IBMCloudServiceIAMVar,
		configv1.IBMCloudServiceKeyProtect:         IBMCloudServiceKeyProtectVar,
		configv1.IBMCloudServiceResourceController: IBMCloudServiceResourceControllerVar,
		configv1.IBMCloudServiceResourceManager:    IBMCloudServiceResourceManagerVar,
		configv1.IBMCloudServiceVPC:                IBMCloudServiceVPCVar,
	}
)

// EndpointsJSON represents the JSON format to override IBM Cloud Terraform provider utilized service endpoints.
// https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/guides/custom-service-endpoints#file-structure-for-endpoints-file
type EndpointsJSON struct {
	// IBMCloudEndpointCIS contains endpoint mapping for IBM Cloud CIS.
	IBMCloudEndpointCIS *EndpointsVisibility `json:"IBMCLOUD_CIS_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointCOS contains endpoint mapping for IBM Cloud COS.
	IBMCloudEndpointCOS *EndpointsVisibility `json:"IBMCLOUD_COS_ENDPOINT,omitempty"`

	// IBMCloudEndpointCOSConfig contains endpoint mapping for IBM Cloud COS Config.
	IBMCloudEndpointCOSConfig *EndpointsVisibility `json:"IBMCLOUD_COS_CONFIG_ENDPOINT,omitempty"`

	// IBMCloudEndpointDNSServices contains endpoint mapping for IBM Cloud DNS Services.
	IBMCloudEndpointDNSServices *EndpointsVisibility `json:"IBMCLOUD_PRIVATE_DNS_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointGlobalCatalog contains endpoint mapping for IBM Cloud Global Catalog.
	IBMCloudEndpointGlobalCatalog *EndpointsVisibility `json:"IBMCLOUD_RESOURCE_CATALOG_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointGlobalSearch contains endpoint mapping for IBM Cloud Global Search.
	IBMCloudEndpointGlobalSearch *EndpointsVisibility `json:"IBMCLOUD_GS_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointGlobalTagging contains endpoint mapping for IBM Cloud Global Tagging.
	IBMCloudEndpointGlobalTagging *EndpointsVisibility `json:"IBMCLOUD_GT_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointHyperProtect contains endpoint mapping for IBM Cloud Hyper Protect.
	IBMCloudEndpointHyperProtect *EndpointsVisibility `json:"IBMCLOUD_HPCS_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointIAM contains endpoint mapping for IBM Cloud IAM.
	IBMCloudEndpointIAM *EndpointsVisibility `json:"IBMCLOUD_IAM_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointKeyProtect contains endpoint mapping for IBM Cloud Key Protect.
	IBMCloudEndpointKeyProtect *EndpointsVisibility `json:"IBMCLOUD_KP_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointResourceController contains endpoint mapping for IBM Cloud Resource Controller.
	IBMCloudEndpointResourceController *EndpointsVisibility `json:"IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointResourceManager contains endpoint mapping for IBM Cloud Resource Manager.
	IBMCloudEndpointResourceManager *EndpointsVisibility `json:"IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT,omitempty"`

	// IBMCloudEndpointVPC contains endpoint mapping for IBM Cloud VPC.
	IBMCloudEndpointVPC *EndpointsVisibility `json:"IBMCLOUD_IS_NG_API_ENDPOINT,omitempty"`
}

// EndpointsVisibility contains region mapped endpoint for a service.
type EndpointsVisibility struct {
	// Private is a string-string map of a region name to endpoint URL
	// To prevent maintaining a list of supported regions here, we simply use a map instead of a struct
	Private map[string]string `json:"private"`

	// Public is a string-string map of a region name to endpoint URL
	// To prevent maintaining a list of supported regions here, we simply use a map instead of a struct
	Public map[string]string `json:"public"`
}

// CheckServiceEndpointOverride checks whether a service has an override endpoint.
func CheckServiceEndpointOverride(service configv1.IBMCloudServiceName, serviceEndpoints []configv1.IBMCloudServiceEndpoint) string {
	if len(serviceEndpoints) > 0 {
		for _, endpoint := range serviceEndpoints {
			if endpoint.Name == service {
				return endpoint.URL
			}
		}
	}
	return ""
}

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the IBM Cloud region where the cluster will be
	// created.
	Region string `json:"region"`

	// ResourceGroupName is the name of an already existing resource group where the
	// cluster should be installed. If empty, a new resource group will be created
	// for the cluster.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`

	// NetworkResourceGroupName is the name of an already existing resource group
	// where an existing VPC and set of Subnets exist, to be used during cluster
	// creation.
	// +optional
	NetworkResourceGroupName string `json:"networkResourceGroupName,omitempty"`

	// VPCName is the name of an already existing VPC to be used during cluster
	// creation.
	// +optional
	VPCName string `json:"vpcName,omitempty"`

	// ControlPlaneSubnets are the names of already existing subnets where the
	// cluster control plane nodes should be created.
	// +optional
	ControlPlaneSubnets []string `json:"controlPlaneSubnets,omitempty"`

	// ComputeSubnets are the names of already existing subnets where the cluster
	// compute nodes should be created.
	// +optional
	ComputeSubnets []string `json:"computeSubnets,omitempty"`

	// DefaultMachinePlatform is the default configuration used when installing
	// on IBM Cloud for machine pools which do not define their own platform
	// configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// ServiceEndpoints is a list which contains custom endpoints to override default
	// service endpoints of IBM Cloud Services.
	// There must only be one ServiceEndpoint for a service (no duplicates).
	// +optional
	ServiceEndpoints []configv1.IBMCloudServiceEndpoint `json:"serviceEndpoints,omitempty"`
}

// ClusterResourceGroupName returns the name of the resource group for the cluster.
func (p *Platform) ClusterResourceGroupName(infraID string) string {
	if len(p.ResourceGroupName) > 0 {
		return p.ResourceGroupName
	}
	return infraID
}

// GetVPCName returns the user provided name of the VPC for the cluster.
func (p *Platform) GetVPCName() string {
	if len(p.VPCName) > 0 {
		return p.VPCName
	}
	return ""
}
