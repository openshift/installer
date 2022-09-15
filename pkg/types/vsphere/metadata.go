package vsphere

// Metadata contains vSphere metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter"`
	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username"`
	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password"`
	// TerraformPlatform is the type...
	TerraformPlatform string `json:"terraform_platform"`
}
