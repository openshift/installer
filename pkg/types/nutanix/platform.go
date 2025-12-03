package nutanix

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
)

// CredentialsSecretName is the default nutanix credentials secret name.
//
//nolint:gosec
const CredentialsSecretName = "nutanix-credentials"

// Platform stores any global configuration used for Nutanix platforms.
type Platform struct {
	// PrismCentral is the endpoint (address and port) and credentials to
	// connect to the Prism Central.
	// This serves as the default Prism-Central.
	PrismCentral PrismCentral `json:"prismCentral"`

	// PrismElements holds a list of Prism Elements (clusters). A Prism Element encompasses all Nutanix resources (VMs, subnets, etc.)
	// used to host the OpenShift cluster. Currently only a single Prism Element may be defined.
	// This serves as the default Prism-Element.
	PrismElements []PrismElement `json:"prismElements"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS Image.
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// PreloadedOSImageName uses the named preloaded RHCOS image from PC/PE,
	// instead of create and upload a new image for each cluster.
	//
	// +optional
	PreloadedOSImageName string `json:"preloadedOSImageName,omitempty"`

	// DeprecatedAPIVIP is the virtual IP address for the api endpoint
	// Deprecated: use APIVIPs
	//
	// +kubebuilder:validation:format=ip
	// +optional
	DeprecatedAPIVIP string `json:"apiVIP,omitempty"`

	// APIVIPs contains the VIP(s) for the api endpoint. In dual stack clusters
	// it contains an IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// DeprecatedIngressVIP is the virtual IP address for ingress
	// Deprecated: use IngressVIPs
	//
	// +kubebuilder:validation:format=ip
	// +optional
	DeprecatedIngressVIP string `json:"ingressVIP,omitempty"`

	// IngressVIPs contains the VIP(s) for ingress. In dual stack clusters
	// it contains an IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Nutanix for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// SubnetUUIDs identifies the network subnets to be used by the cluster.
	// Currently we only support one subnet for an OpenShift cluster.
	SubnetUUIDs []string `json:"subnetUUIDs"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// LoadBalancer is available in TechPreview.
	// +optional
	LoadBalancer *configv1.NutanixPlatformLoadBalancer `json:"loadBalancer,omitempty"`

	// dnsRecordsType determines whether records for api, api-int, and ingress
	// are provided by the internal DNS service or externally.
	// Allowed values are `Internal`, `External`, and omitted.
	// When set to `Internal`, records are provided by the internal infrastructure and
	// no additional user configuration is required for the cluster to function.
	// When set to `External`, records are not provided by the internal infrastructure
	// and must be configured by the user on a DNS server outside the cluster.
	// Cluster nodes must use this external server for their upstream DNS requests.
	// This value may only be set when loadBalancer.type is set to UserManaged.
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default is `Internal`.
	// +openshift:enable:FeatureGate=OnPremDNSRecords
	// +optional
	DNSRecordsType configv1.DNSRecordsType `json:"dnsRecordsType,omitempty"`

	// FailureDomains configures failure domains for the Nutanix platform.
	// +optional
	FailureDomains []FailureDomain `json:"failureDomains,omitempty"`

	// PrismAPICallTimeout sets the timeout (in minutes) for the prism-api calls.
	// If not configured, the default value of 10 minutes will be used as the prism-api call timeout.
	// +optional
	PrismAPICallTimeout *int `json:"prismAPICallTimeout,omitempty"`
}

// PrismCentral holds the endpoint and credentials data used to connect to the Prism Central
type PrismCentral struct {
	// Endpoint holds the address and port of the Prism Central
	Endpoint PrismEndpoint `json:"endpoint"`

	// Username is the name of the user to connect to the Prism Central
	Username string `json:"username"`

	// Password is the password for the user to connect to the Prism Central
	Password string `json:"password"`
}

// PrismElement holds the uuid, endpoint of the Prism Element (cluster)
type PrismElement struct {
	// UUID is the UUID of the Prism Element (cluster)
	UUID string `json:"uuid"`

	// Endpoint holds the address and port of the Prism Element
	// +optional
	Endpoint PrismEndpoint `json:"endpoint,omitempty"`

	// Name is prism endpoint Name
	Name string `json:"name,omitempty"`
}

// PrismEndpoint holds the endpoint address and port to access the Nutanix Prism Central or Element (cluster)
type PrismEndpoint struct {
	// address is the endpoint address (DNS name or IP address) of the Nutanix Prism Central or Element (cluster)
	Address string `json:"address"`

	// port is the port number to access the Nutanix Prism Central or Element (cluster)
	Port int32 `json:"port"`
}

// FailureDomain configures failure domain information for the Nutanix platform.
type FailureDomain struct {
	// Name defines the unique name of a failure domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z_.-@/]+$`
	Name string `json:"name"`

	// prismElement holds the identification (name, uuid) and the optional endpoint address and port of the Nutanix Prism Element.
	// When a cluster-wide proxy is installed, by default, this endpoint will be accessed via the proxy.
	// Should you wish for communication with this endpoint not to be proxied, please add the endpoint to the
	// proxy spec.noProxy list.
	// +kubebuilder:validation:Required
	PrismElement PrismElement `json:"prismElement"`

	// SubnetUUIDs identifies the network subnets of the Prism Element.
	// Currently we only support one subnet for a failure domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	SubnetUUIDs []string `json:"subnetUUIDs"`

	// StorageContainers identifies the storage containers in the Prism Element.
	// +optional
	StorageContainers []StorageResourceReference `json:"storageContainers,omitempty"`

	// DataSourceImages identifies the datasource images in the Prism Element.
	// +optional
	DataSourceImages []StorageResourceReference `json:"dataSourceImages,omitempty"`
}

// GetFailureDomainByName returns the NutanixFailureDomain pointer with the input name.
// Returns nil if not found.
func (p Platform) GetFailureDomainByName(fdName string) (*FailureDomain, error) {
	for _, fd := range p.FailureDomains {
		if fd.Name == fdName {
			return &fd, nil
		}
	}

	return nil, fmt.Errorf("not found the defined failure domain with name %q", fdName)
}

// GetStorageContainerFromFailureDomain returns the storage container configuration with the provided reference and failuer domain names.
// Returns nil and error if not found.
func (p Platform) GetStorageContainerFromFailureDomain(fdName, storageContainerRefName string) (*StorageResourceReference, error) {
	for _, fd := range p.FailureDomains {
		if fd.Name == fdName {
			for _, sc := range fd.StorageContainers {
				if sc.ReferenceName == storageContainerRefName {
					return &sc, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found the storage container with reference name %q in failureDomain %q", storageContainerRefName, fdName)
}

// GetDataSourceImageFromFailureDomain returns the datasource image configuration with the provided reference and failuer domain names.
// Returns nil and error if not found.
func (p Platform) GetDataSourceImageFromFailureDomain(fdName, dataSourceRefName string) (*StorageResourceReference, error) {
	for _, fd := range p.FailureDomains {
		if fd.Name == fdName {
			for _, dsi := range fd.DataSourceImages {
				if dsi.ReferenceName == dataSourceRefName {
					return &dsi, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found the datasource image with reference name %q in failureDomain %q", dataSourceRefName, fdName)
}
