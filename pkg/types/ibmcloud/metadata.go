package ibmcloud

// Metadata contains IBM Cloud metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region            string `json:"region"`
	ResourceGroupName string `json:"resourceGroupName"`
}
