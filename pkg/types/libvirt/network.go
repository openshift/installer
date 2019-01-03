package libvirt

// Network is the configuration of the libvirt network.
type Network struct {
	// IfName is the name of the network interface.
	IfName string `json:"if"`
}
