package portstrustedvif

type PortTrustedVIFExt struct {
	// PortTrustedVIF stores information about whether a SR-IOV port should be trusted
	PortTrustedVIF *bool `json:"trusted"`
}
