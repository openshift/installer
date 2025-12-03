package openstack

import (
	configv1 "github.com/openshift/api/config/v1"
)

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	// Deprecated: this value is not used by the installer.
	// +optional
	DeprecatedRegion string `json:"region,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on OpenStack for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Cloud is the name of OpenStack cloud to use from clouds.yaml.
	Cloud string `json:"cloud"`

	// ExternalNetwork is name of the external network in your OpenStack cluster.
	// +optional
	ExternalNetwork string `json:"externalNetwork,omitempty"`

	// DeprecatedFlavorName is the name of the flavor to use for instances in this cluster.
	// Deprecated: use FlavorName in DefaultMachinePlatform to define default flavor.
	// +optional
	DeprecatedFlavorName string `json:"computeFlavor,omitempty"`

	// LbFloatingIP is the IP address of an available floating IP in your OpenStack cluster
	// to associate with the OpenShift load balancer.
	// Deprecated: this value has been renamed to apiFloatingIP.
	// +optional
	DeprecatedLbFloatingIP string `json:"lbFloatingIP,omitempty"`

	// APIFloatingIP is the IP address of an available floating IP in your OpenStack cluster
	// to associate with the OpenShift API load balancer.
	// +optional
	APIFloatingIP string `json:"apiFloatingIP,omitempty"`

	// IngressFloatingIP is the ID of an available floating IP in your OpenStack cluster
	// that will be associated with the OpenShift ingress port
	// +optional
	IngressFloatingIP string `json:"ingressFloatingIP,omitempty"`

	// ExternalDNS holds the IP addresses of dns servers that will
	// be added to the dns resolution of all instances in the cluster.
	// +optional
	ExternalDNS []string `json:"externalDNS"`

	// TrunkSupport holds a `0` or `1` value that indicates whether or not to use trunk ports
	// in your OpenShift cluster.
	// Deprecated: the machine manifest should be used to specify that trunk should be used.
	// +optional
	DeprecatedTrunkSupport string `json:"trunkSupport,omitempty"`

	// OctaviaSupport holds a `0` or `1` value that indicates whether your OpenStack
	// cluster supports Octavia Loadbalancing.
	// Deprecated: this value is set by the installer automatically.
	// +optional
	DeprecatedOctaviaSupport string `json:"octaviaSupport,omitempty"`

	// ClusterOSImage is either a URL with `http(s)` or `file` scheme to override
	// the default OS image for cluster nodes, or an existing Glance image name.
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// ClusterOSImageProperties is a list of properties to be added to the metadata of the uploaded Glance ClusterOSImage.
	// Default: the default is to not set any properties.
	// +optional
	ClusterOSImageProperties map[string]string `json:"clusterOSImageProperties,omitempty"`

	// DeprecatedAPIVIP is the static IP on the nodes subnet that the api port for openshift will be assigned
	// Default: will be set to the 5 on the first entry in the machineNetwork CIDR
	// Deprecated: Use APIVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedAPIVIP string `json:"apiVIP,omitempty"`

	// APIVIPs contains the VIP(s) on the nodes subnet that the api port for
	// openshift will be assigned. In dual stack clusters it contains an IPv4
	// and IPv6 address, otherwise only one VIP
	// Default: will be set to the 5 on the first entry in the machineNetwork
	// CIDR
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// DeprecatedIngressVIP is the static IP on the nodes subnet that the apps port for openshift will be assigned
	// Default: will be set to the 7 on the first entry in the machineNetwork CIDR
	// Deprecated: Use IngressVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedIngressVIP string `json:"ingressVIP,omitempty"`

	// IngressVIPs contains the VIP(s) on the nodes subnet that the apps port
	// for openshift will be assigned. In dual stack clusters it contains an
	// IPv4 and IPv6 address, otherwise only one VIP
	// Default: will be set to the 7 on the first entry in the machineNetwork
	// CIDR
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

	// DeprecatedMachinesSubnet is a string of the UUIDv4 of an openstack subnet. This subnet will be used by all nodes created by the installer.
	// By setting this, the installer will no longer create a network and subnet.
	// The subnet and network specified in MachinesSubnet will not be deleted or modified by the installer.
	// Deprecated: Use ControlPlanePort
	// +optional
	DeprecatedMachinesSubnet string `json:"machinesSubnet,omitempty"`

	// ControlPlanePort contains details of the network attached to the control plane port, with the network either containing one openstack
	// subnet for IPv4 or two openstack subnets for dualstack clusters. Providing this configuration will prevent OpenShift from managing
	// or updating this network and its subnets. The network and its subnets will be used during creation of all nodes.
	// +optional
	ControlPlanePort *PortTarget `json:"controlPlanePort,omitempty"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// +optional
	LoadBalancer *configv1.OpenStackPlatformLoadBalancer `json:"loadBalancer,omitempty"`

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
}
