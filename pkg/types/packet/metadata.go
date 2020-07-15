package packet

// Metadata contains packet metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// FacilityCode represents the Packet region and datacenter where your devices will be provisioned (https://www.packet.com/developers/docs/getting-started/facilities/)
	FacilityCode string `json:"facility_code,omitempty"`

	// ProjectID represents the Packet project used for logical grouping and invoicing (https://www.packet.com/developers/docs/API/getting-started/)
	ProjectID string `json:"project_id,omitempty"`
}
