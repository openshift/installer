package manifests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func bgpInstallConfig() *installconfig.InstallConfig {
	return installconfig.MakeAsset(&types.InstallConfig{
		Platform: types.Platform{
			BareMetal: &baremetal.Platform{
				APIVIPs:     []string{"192.168.111.5"},
				IngressVIPs: []string{"192.168.111.4"},
				BGPVIPConfig: &baremetal.BGPVIPConfig{
					LocalASN: 64512,
					Peers: []baremetal.BGPPeerConfig{
						{PeerAddress: "192.168.111.1", PeerASN: 64513},
					},
				},
				Hosts: []*baremetal.Host{
					{
						Name: "master-0",
						BGPPeers: []baremetal.BGPPeerConfig{
							{PeerAddress: "192.168.1.1", PeerASN: 64513},
						},
					},
				},
			},
		},
	})
}

// The config.json schema must match baremetal-runtimecfg's FRRPeerMapping:
// top-level "defaultPeers" (not "peers") and hostOverrides as a flat
// map[hostname][]peer (not map[hostname]{"peers":[...]}).
func TestBGPVIPConfigMapSchemaMatchesRuntimecfg(t *testing.T) {
	parents := asset.Parents{}
	parents.Add(bgpInstallConfig())

	cmAsset := &BGPVIPConfigMap{}
	if !assert.NoError(t, cmAsset.Generate(context.Background(), parents)) {
		return
	}
	if !assert.NotNil(t, cmAsset.ConfigMap) {
		return
	}

	var got map[string]json.RawMessage
	if !assert.NoError(t, json.Unmarshal([]byte(cmAsset.ConfigMap.Data["config.json"]), &got)) {
		return
	}

	assert.Contains(t, got, "defaultPeers")
	assert.NotContains(t, got, "peers")

	var overrides map[string][]baremetal.BGPPeerConfig
	if !assert.NoError(t, json.Unmarshal(got["hostOverrides"], &overrides)) {
		return
	}
	if !assert.Len(t, overrides["master-0"], 1) {
		return
	}
	assert.Equal(t, "192.168.1.1", overrides["master-0"][0].PeerAddress)
}
