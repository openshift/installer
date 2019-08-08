package openstack

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on OpenStack for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Cloud
	// Name of OpenStack cloud to use from clouds.yaml
	Cloud string `json:"cloud"`

	// ExternalNetwork
	// The OpenStack external network name to be used for installation.
	ExternalNetwork string `json:"externalNetwork"`

	// FlavorName
	// The OpenStack compute flavor to use for servers.
	FlavorName string `json:"computeFlavor"`

	// LbFloatingIP
	// Existing Floating IP to associate with the OpenStack load balancer.
	LbFloatingIP string `json:"lbFloatingIP"`

	// APIVIP
	// The internal virtual IP address (VIP) put in front of the
	// Kubernetes API server for use by components inside the cluster.
	// The DNS static pods running on the nodes resolve the api-int
	// record to APIVIP.
	//
	// The value is set by the installer from the MachineCIDR range.
	APIVIP string `json:"apiVIP"`

	// DNSVIP
	// The internal virtual IP address (VIP) put in front of the DNS
	// static pods running on the nodes. Unlike the DNS operator these
	// services provide name resolution for the nodes themselves.
	//
	// The value is set by the installer from the MachineCIDR range.
	DNSVIP string `json:"dnsVIP"`

	// IngressVIP
	// The internal virtual IP address (VIP) put in front of the
	// OpenShift router pods. This provides the internal accessibility
	// to the internal pods running on the worker nodes, e.g.
	// `console`. The DNS static pods running on the nodes resolve the
	// wildcard apps record to IngressVIP.
	//
	// The value is set by the installer from the MachineCIDR range.
	IngressVIP string `json:"ingressVIP"`

	// TrunkSupport
	// Whether OpenStack ports can be trunked
	TrunkSupport string `json:"trunkSupport"`

	// OctaviaSupport
	// Whether OpenStack has Octavia support
	OctaviaSupport string `json:"octaviaSupport"`
}
