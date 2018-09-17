package installconfig

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
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

func TestInstallConfigGenerate(t *testing.T) {
	cases := []struct {
		name                 string
		platformContents     []string
		expectedPlatformYaml string
	}{
		{
			name: "aws",
			platformContents: []string{
				"aws",
				"test-region",
			},
			expectedPlatformYaml: `  aws:
    region: test-region`,
		},
		{
			name: "libvirt",
			platformContents: []string{
				"libvirt",
				"test-uri",
			},
			expectedPlatformYaml: `  libvirt:
    URI: test-uri
    network: {}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stock := &StockImpl{
				clusterID:    &testAsset{},
				emailAddress: &testAsset{},
				password:     &testAsset{},
				sshKey:       &testAsset{},
				baseDomain:   &testAsset{},
				clusterName:  &testAsset{},
				pullSecret:   &testAsset{},
				platform:     &testAsset{},
			}

			dir, err := ioutil.TempDir("", "TestInstallConfigGenerate")
			if err != nil {
				t.Skipf("could not create temporary directory: %v", err)
			}
			defer os.RemoveAll(dir)

			installConfig := &installConfig{
				assetStock: stock,
				directory:  dir,
			}

			states := map[asset.Asset]*asset.State{
				stock.clusterID: {
					Contents: []asset.Content{{Data: []byte("test-cluster-id")}},
				},
				stock.emailAddress: {
					Contents: []asset.Content{{Data: []byte("test-email")}},
				},
				stock.password: {
					Contents: []asset.Content{{Data: []byte("test-password")}},
				},
				stock.sshKey: {
					Contents: []asset.Content{{Data: []byte("test-sshkey")}},
				},
				stock.baseDomain: {
					Contents: []asset.Content{{Data: []byte("test-domain")}},
				},
				stock.clusterName: {
					Contents: []asset.Content{{Data: []byte("test-cluster-name")}},
				},
				stock.pullSecret: {
					Contents: []asset.Content{{Data: []byte("test-pull-secret")}},
				},
				stock.platform: {
					Contents: make([]asset.Content, len(tc.platformContents)),
				},
			}
			for i, c := range tc.platformContents {
				states[stock.platform].Contents[i].Data = []byte(c)
			}

			state, err := installConfig.Generate(states)
			assert.NoError(t, err, "unexpected error generating asset")
			assert.NotNil(t, state, "unexpected nil for asset state")

			filename := filepath.Join(dir, "install-config.yml")
			assert.Equal(t, 1, len(state.Contents), "unexpected number of contents in asset state")
			assert.Equal(t, filename, state.Contents[0].Name, "unexpected filename in asset state")

			exp := fmt.Sprintf(`admin:
  email: test-email
  password: test-password
  sshKey: test-sshkey
baseDomain: test-domain
clusterID: test-cluster-id
machines:
- name: master
  platform: {}
  replicas: 3
- name: worker
  platform: {}
  replicas: 3
metadata:
  creationTimestamp: null
  name: test-cluster-name
networking:
  podCIDR: 10.2.0.0/16
  serviceCIDR: 10.3.0.0/16
platform:
%s
pullSecret: test-pull-secret
`, tc.expectedPlatformYaml)

			assert.Equal(t, exp, string(state.Contents[0].Data), "unexpected data in install-config.yml")
		})
	}
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
