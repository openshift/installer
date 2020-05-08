package ovirt

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// The target cluster under which all VMs will run
	ClusterID string `json:"ovirt_cluster_id"`

	// The target storage domain under which all VM disk would be created.
	StorageDomainID string `json:"ovirt_storage_domain_id"`

	// NetworkName is the target network of all the network interfaces of the nodes.
	// When no ovirt_network_name is provided it defaults to `ovirtmgmt` network, which is a default network for every ovirt cluster.
	// +optional
	NetworkName string `json:"ovirt_network_name,omitempty"`

	// VNICProfileID defines the VNIC profile ID to use the the VM network interfaces.
	// When no vnicProfileID is provided it will be set to the profile of the network. If there are multiple
	// profiles for the network, the installer requires you to explicitly set the vnicProfileID.
	// +optional
	VNICProfileID string `json:"vnicProfileID,omitempty"`

	// APIVIP is an IP which will be served by bootstrap and then pivoted masters, using keepalived
	APIVIP string `json:"api_vip"`

	// DNSVIP is the IP of the internal DNS which will be operated by the cluster
	DNSVIP string `json:"dns_vip"`

	// IngressIP is an external IP which routes to the default ingress controller.
	// The IP is a suitable target of a wildcard DNS record used to resolve default route host names.
	IngressVIP string `json:"ingress_vip"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on ovirt for machine pools which do not define their
	// own platform configuration.
	// Default will set the image field to the latest RHCOS image.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
