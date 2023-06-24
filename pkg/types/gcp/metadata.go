package gcp

// Metadata contains GCP metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region            string `json:"region"`
	ProjectID         string `json:"projectID"`
	NetworkProjectID  string `json:"networkProjectID,omitempty"`
	PrivateZoneDomain string `json:"privateZoneDomain,omitempty"`
}
