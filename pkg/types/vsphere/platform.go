package vsphere

// Platform stores any global configuration used for vsphere platforms.
type Platform struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter"`

	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username"`

	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password"`

	// Datacenter is the name of the datacenter to use in the vCenter.
	Datacenter string `json:"datacenter"`

	// DefaultDatastore is the default datastore to use for provisioning volumes.
	DefaultDatastore string `json:"defaultDatastore"`

	// Folder is the name of the folder that will be used and/or created for
	// virtual machines.
	Folder string `json:"folder,omitempty"`

	// Cluster is the name of the cluster virtual machines will be cloned into.
	Cluster string `json:"cluster,omitempty"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS OVA
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// APIVIP is the virtual IP address for the api endpoint
	APIVIP string `json:"apiVIP,omitempty"`

	// IngressVIP is the virtual IP address for ingress
	IngressVIP string `json:"ingressVIP,omitempty"`

	// DNSVIP is the virtual IP address for DNS
	DNSVIP string `json:"dnsVIP,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on VSphere for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network specifies the name of the network to be used by the cluster.
	Network string `json:"network,omitempty"`
}
