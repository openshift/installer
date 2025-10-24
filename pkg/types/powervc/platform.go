package powervc

import (
	configv1 "github.com/openshift/api/config/v1"
)

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
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

	// ClusterOSImage is either a URL with `http(s)` or `file` scheme to override
	// the default OS image for cluster nodes, or an existing Glance image name.
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// ClusterOSImageProperties is a list of properties to be added to the metadata of the uploaded Glance ClusterOSImage.
	// Default: the default is to not set any properties.
	// +optional
	ClusterOSImageProperties map[string]string `json:"clusterOSImageProperties,omitempty"`

	// APIVIPs contains the VIP(s) on the nodes subnet that the api port for
	// openshift will be assigned. In dual stack clusters it contains an IPv4
	// and IPv6 address, otherwise only one VIP
	// Default: will be set to the 5 on the first entry in the machineNetwork
	// CIDR
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// IngressVIPs contains the VIP(s) on the nodes subnet that the apps port
	// for openshift will be assigned. In dual stack clusters it contains an
	// IPv4 and IPv6 address, otherwise only one VIP
	// Default: will be set to the 7 on the first entry in the machineNetwork
	// CIDR
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

	// ControlPlanePort contains details of the network attached to the control plane port, with the network either containing one openstack
	// subnet for IPv4 or two openstack subnets for dualstack clusters. Providing this configuration will prevent OpenShift from managing
	// or updating this network and its subnets. The network and its subnets will be used during creation of all nodes.
	// +optional
	ControlPlanePort *PortTarget `json:"controlPlanePort,omitempty"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// +optional
	LoadBalancer *configv1.OpenStackPlatformLoadBalancer `json:"loadBalancer,omitempty"`

	// +optional
	ImageName string `json:"imageName,omitempty"`
}
