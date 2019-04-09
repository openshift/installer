package vsphere

// Platform stores any global configuration used for vsphere platforms.
type Platform struct {
	// VirtualCenters are the configurations for the vCenters.
	VirtualCenters []VirtualCenter `json:"virtualCenters"`
	// Workspace is the configuration for the workspace.
	Workspace Workspace `json:"workspace"`
	// SCSIControllerType is the SCSI controller type in use.
	// +optional
	// Default is pvscsi.
	SCSIControllerType string `json:"scsiControllerType"`
	// PublicNetwork is the name of the VM network to use.
	PublicNetwork string `json:"publicNetwork"`
}

// VirtualCenter is the configuration of a vCenter.
type VirtualCenter struct {
	// Name of the vCenter. This is the domain name or the IP address of the vCenter.
	Name string `json:"name"`
	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username"`
	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password"`
	// Datacenters are the names of the datacenters to use in the vCenter.
	Datacenters []string `json:"datacenters"`
	// Insecure should be true if the vCenter uses a self-signed cert.
	// +optional
	// The default is false.
	Insecure bool `json:"insecure,omitempty"`
}

// Workspace is the configuration of the vSphere workspace.
type Workspace struct {
	// Server is the server to use for provisioning.
	// +optional
	// Default is the name of the vCenter, if there is only a single vCenter.
	Server string `json:"server"`
	// Datacenter is the datacenter to use for provisioning.
	// +optional
	// Default is the datacenter in Server, if that vCenter has only a single
	// datacenter.
	Datacenter string `json:"datacenter"`
	// DefaultDatastore is the default datastore to use for provisioning volumes.
	DefaultDatastore string `json:"defaultDatastore"`
	// ResourcePoolPath is the resource pool to use in the datacenter.
	// +optional
	// Default is the name of the cluster.
	ResourcePoolPath string `json:"resourcePoolPath"`
	// Folder is the vCenter VM folder path in the datacenter.
	// +optional
	// Default is the name of the cluster.
	Folder string `json:"folder"`
}
