package portstrustedvif

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

// PortCreateOptsExt adds port trusted VIF options to the base ports.CreateOpts.
type PortCreateOptsExt struct {
	ports.CreateOptsBuilder

	// PortTrustedVIF toggles the port's trusted VIF status.
	PortTrustedVIF *bool `json:"trusted,omitempty"`
}

// To PortCreateMap casts a CreateOpts struct to a map
func (opts PortCreateOptsExt) ToPortCreateMap() (map[string]any, error) {
	base, err := opts.CreateOptsBuilder.ToPortCreateMap()
	if err != nil {
		return nil, err
	}

	port := base["port"].(map[string]any)

	if opts.PortTrustedVIF != nil {
		port["trusted"] = *opts.PortTrustedVIF
	}

	return base, nil
}

// PortUpdateOptsExt adds port trusted VIF options to the base ports.UpdateOpts.
type PortUpdateOptsExt struct {
	ports.UpdateOptsBuilder

	// PortTrustedVIF updates the port's trusted VIF status.
	PortTrustedVIF *bool `json:"trusted,omitempty"`
}

// ToPortUpdateMap casts a UpdateOpts struct to a map.
func (opts PortUpdateOptsExt) ToPortUpdateMap() (map[string]any, error) {
	base, err := opts.UpdateOptsBuilder.ToPortUpdateMap()
	if err != nil {
		return nil, err
	}

	port := base["port"].(map[string]any)

	if opts.PortTrustedVIF != nil {
		port["trusted"] = *opts.PortTrustedVIF
	}

	return base, nil
}
