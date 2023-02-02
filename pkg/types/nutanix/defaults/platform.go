package defaults

import (
	"log"
	"os"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/nutanix"
)

var (
	defaultMachineCIDR = ipnet.MustParseCIDR("10.0.0.0/16")
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *nutanix.Platform) {}

func GetMachineCIDR() *ipnet.IPNet {
	cidrOverride, ok := os.LookupEnv("NUTANIX_MACHINE_CIDR_OVERRIDE")
	if ok && cidrOverride != "" {
		log.Println("NUTANIX_MACHINE_CIDR_OVERRIDE is set, using it instead of default machine CIDR", "NUTANIX_MACHINE_CIDR_OVERRIDE", cidrOverride)
		cidr, err := ipnet.ParseCIDR(cidrOverride)
		if err != nil {
			return defaultMachineCIDR
		}
		return cidr
	}

	return defaultMachineCIDR
}
