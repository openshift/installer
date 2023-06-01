package defaults

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/powervs"
)

var (
	// DefaultMachineCIDR is the PowerVS default IP address space from
	// which to assign machine IPs.
	DefaultMachineCIDR = ipnet.MustParseCIDR("192.168.0.0/24")
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *powervs.Platform) {
	n, err := rand.Int(rand.Reader, big.NewInt(253))
	if err != nil {
		panic(err)
	}
	subnet := n.Int64()
	// MustParseCIDR parses a CIDR from its string representation. If the parse fails, the function will panic.
	DefaultMachineCIDR = ipnet.MustParseCIDR(fmt.Sprintf("192.168.%d.0/24", subnet))
}
