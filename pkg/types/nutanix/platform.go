package nutanix

import (
	configv1 "github.com/openshift/api/config/v1"
)

// Platform stores any global configuration used for Nutanix platforms.
type Platform struct {
	// PrismCentral is the endpoint (address and port) and credentials to
	// connect to the Prism Central.
	PrismCentral PrismCentral `json:"prismCentral"`

	// PrismElements holds a list of Prism Elements (clusters). A Prism Element encompasses all Nutanix resources (VMs, subnets, etc.)
	// used to host the OpenShift cluster. Currently only a single Prism Element may be defined.
	PrismElements []PrismElement `json:"prismElements"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS Image
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

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
	// +kubebuilder:validation:UniqueItems=true
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
	// +kubebuilder:validation:UniqueItems=true
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
	Endpoint PrismEndpoint `json:"endpoint"`

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
