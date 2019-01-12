package libvirt

// Network is the configuration of the libvirt network.
type Network struct {
	// +optional
	// Default is tt0.
	IfName string `json:"if,omitempty"`
}
