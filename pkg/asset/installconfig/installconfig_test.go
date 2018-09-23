package installconfig

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

type testAsset struct {
	name string
}

func (a *testAsset) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

func (a *testAsset) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	return nil, nil
}

func (a *testAsset) Name() string {
	return "Test Asset"
}

func TestInstallConfigDependencies(t *testing.T) {
	stock := &StockImpl{
		clusterID:    &testAsset{name: "test-cluster-id"},
		emailAddress: &testAsset{name: "test-email"},
		password:     &testAsset{name: "test-password"},
		sshKey:       &testAsset{name: "test-sshkey"},
		baseDomain:   &testAsset{name: "test-domain"},
		clusterName:  &testAsset{name: "test-cluster"},
		pullSecret:   &testAsset{name: "test-pull-secret"},
		platform:     &testAsset{name: "test-platform"},
	}
	installConfig := &installConfig{
		assetStock: stock,
	}
	exp := []string{
		"test-cluster-id",
		"test-email",
		"test-password",
		"test-sshkey",
		"test-domain",
		"test-cluster",
		"test-pull-secret",
		"test-platform",
	}
	deps := installConfig.Dependencies()
	act := make([]string, len(deps))
	for i, d := range deps {
		a, ok := d.(*testAsset)
		assert.True(t, ok, "expected dependency to be a *testAsset")
		act[i] = a.name
	}
	assert.Equal(t, exp, act, "unexpected dependency")
}

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
