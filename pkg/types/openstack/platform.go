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
	// IP in the machineCIDR to use for api-int.
	APIVIP string `json:"apiVIP"`

	// DNSVIP
	// IP in the machineCIDR to use for simulated route53.
	DNSVIP string `json:"dnsVIP"`

	// TrunkSupport
	// Whether OpenStack ports can be trunked
	TrunkSupport string `json:"trunkSupport"`

	// OctaviaSupport
	// Whether OpenStack has Octavia support
	OctaviaSupport string `json:"octaviaSupport"`
}
