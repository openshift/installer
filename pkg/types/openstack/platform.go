package openstack

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	// Deprecated: this value is not used by the installer.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on OpenStack for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Cloud is the name of OpenStack cloud to use from clouds.yaml.
	Cloud string `json:"cloud"`

	// ExternalNetwork is name of the external network in your OpenStack cluster.
	ExternalNetwork string `json:"externalNetwork"`

	// FlavorName is the name of the compute flavor to use for instances in this cluster.
	FlavorName string `json:"computeFlavor"`

	// LbFloatingIP is the IP address of an available floating IP in your OpenStack cluster
	// to associate with the OpenShift load balancer.
	LbFloatingIP string `json:"lbFloatingIP"`

	// ExternalDNS holds the IP addresses of dns servers that will
	// be added to the dns resolution of all instances in the cluster.
	// +optional
	ExternalDNS []string `json:"externalDNS"`

	// TrunkSupport holds a `0` or `1` value that indicates whether or not to use trunk ports
	// in your OpenShift cluster.
	TrunkSupport string `json:"trunkSupport"`

	// OctaviaSupport holds a `0` or `1` value that indicates whether your OpenStack
	// cluster supports Octavia Loadbalancing.
	OctaviaSupport string `json:"octaviaSupport"`

	// ClusterOSImage is either a URL with `http(s)` or `file` scheme to override
	// the default OS image for cluster nodes, or an existing Glance image name.
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// APIVIP is the static IP on the nodes subnet that the api port for openshift will be assigned
	// Default: will be set to the 5 on the first entry in the machineNetwork CIDR
	// +optional
	APIVIP string `json:"apiVIP,omitempty"`

	// IngressVIP is the static IP on the nodes subnet that the apps port for openshift will be assigned
	// Default: will be set to the 7 on the first entry in the machineNewtwork CIDR
	// +optional
	IngressVIP string `json:"ingressVIP,omitempty"`

	// MachinesSubnet is the UUIDv4 of an openstack subnet. This subnet will be used by all nodes created by the installer.
	// By setting this, the installer will no longer create a network and subnet.
	// The subnet and network specified in MachinesSubnet will not be deleted or modified by the installer.
	// +optional
	MachinesSubnet string `json:"machinesSubnet,omitempty"`
}
