package equinixmetal

// Metadata contains equinixmetal metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// FacilityCode represents the Equinix Metal region and datacenter where your devices will be provisioned (https://metal.equinix.com/developers/docs/getting-started/facilities/)
	FacilityCode string `json:"facility_code,omitempty"`

	// ProjectID represents the Equinix Metal project used for logical grouping and invoicing (https://metal.equinix.com/developers/docs/API/getting-started/)
	ProjectID string `json:"project_id,omitempty"`
}
