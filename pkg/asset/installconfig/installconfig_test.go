package installconfig

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

// TestClusterDNSIP tests the ClusterDNSIP function.
func TestClusterDNSIP(t *testing.T) {
	_, cidr, err := net.ParseCIDR("10.0.1.0/24")
	assert.NoError(t, err, "unexpected error parsing CIDR")
	installConfig := &types.InstallConfig{
		Networking: types.Networking{
			ServiceCIDR: ipnet.IPNet{
				IPNet: *cidr,
			},
		},
	}
	expected := "10.0.1.10"
	actual, err := ClusterDNSIP(installConfig)
	assert.NoError(t, err, "unexpected error get cluster DNS IP")
	assert.Equal(t, expected, actual, "unexpected DNS IP")
}
