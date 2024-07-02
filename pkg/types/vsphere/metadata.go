package vsphere

// Metadata contains vSphere metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter,omitempty"`
	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username,omitempty"`
	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password,omitempty"`
	// TerraformPlatform is the type...
	TerraformPlatform string `json:"terraform_platform"`
	// VCenters collection of vcenters when multi vcenter support is enabled
	VCenters []VCenters
}

// VCenters contains information on individual vcenter.
type VCenters struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter"`
	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username"`
	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password"`
}
