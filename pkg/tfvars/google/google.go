package google

// GCP converts GCP related config.
type GCP struct {
	ImageNameOverride string            `json:"google_image_name_override,omitempty"`
	ExtraLabels       map[string]string `json:"google_extra_labels,omitempty"`
	Master            `json:",inline"`
	Region            string `json:"google_region,omitempty"`
	Worker            `json:",inline"`
}

// Master converts master related config.
type Master struct {
	InstanceType     string `json:"google_master_instance_type,omitempty"`
	MasterRootVolume `json:",inline"`
}

// MasterRootVolume converts master rool volume related config.
type MasterRootVolume struct {
	Size int    `json:"google_master_root_volume_size,omitempty"`
	Type string `json:"google_master_root_volume_type,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
}
