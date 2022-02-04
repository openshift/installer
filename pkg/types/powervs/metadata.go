package powervs

// Metadata contains Power VS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	CISInstanceCRN string `json:"cisInstanceCRN"`
	Region         string `json:"region"`
	Zone           string `json:"zone"`
}
