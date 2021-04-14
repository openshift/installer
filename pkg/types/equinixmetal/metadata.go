package equinixmetal

// Metadata contains equinixmetal metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// Metro represents the Equinix Metal metro code where your devices will be
	// provisioned
	// (https://metal.equinix.com/developers/docs/getting-started/facilities/)
	Metro string `json:"metro,omitempty"`

	// Facility represents the Equinix Metal facility code where your devices
	// will be provisioned
	// (https://metal.equinix.com/developers/docs/getting-started/facilities/)
	Facility string `json:"facility,omitempty"`

	// ProjectID represents the Equinix Metal project used for logical grouping
	// and invoicing
	// (https://metal.equinix.com/developers/docs/API/getting-started/)
	ProjectID string `json:"project_id,omitempty"`
}
