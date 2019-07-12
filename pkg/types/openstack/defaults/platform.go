package defaults

import (
	"fmt"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Defaults for the openstack platform.
const (
	APIVIP = ""
	DNSVIP = ""
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *openstack.Platform, c *types.InstallConfig) {
	if p.APIVIP == "" {
		vip, err := cidr.Host(&c.Networking.MachineCIDR.IPNet, 5)
		if err != nil {
			p.APIVIP = fmt.Sprintf("failed to get APIVIP from MachineCIDR: %s", err.Error())
		} else {
			p.APIVIP = vip.String()
		}
	}

	if p.DNSVIP == "" {
		vip, err := cidr.Host(&c.Networking.MachineCIDR.IPNet, 6)
		if err != nil {
			p.DNSVIP = fmt.Sprintf("failed to get DNSVIP from MachineCIDR: %s", err.Error())
		} else {
			p.DNSVIP = vip.String()
		}
	}
}
