package google

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// Region specifies the GCP region where the cluster will be created.
	Region string `json:"region"`

	// UserTags specifies additional tags for GCP resources created for the cluster.
	UserTags map[string]string `json:"userTags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on GCP for machine pools which do not define their own
	// platform configuration.
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
