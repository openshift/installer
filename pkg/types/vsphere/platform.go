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
}
