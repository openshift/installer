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

// FRR's "timers <keepalive> <hold>" takes bare seconds; the install-config
// carries the typed install-config fields, so the generated peer data must
// carry the boolean-string and whole-second decimal-string shapes of the
// merged runtimecfg/MCO/CNO contract (installer#10718 review 4919561631).
func TestBGPVIPConfigMapContractMapping(t *testing.T) {
	ic := bgpInstallConfig()
	ic.Config.Platform.BareMetal.BGPVIPConfig.Peers[0].HoldTimeSeconds = 90
	ic.Config.Platform.BareMetal.BGPVIPConfig.Peers[0].KeepaliveTimeSeconds = 30
	ic.Config.Platform.BareMetal.BGPVIPConfig.Peers[0].FailureDetection = baremetal.BGPFailureDetectionBFD
	ic.Config.Platform.BareMetal.BGPVIPConfig.Peers[0].PeerReachability = baremetal.BGPPeerReachabilityMultiHop
	ic.Config.Platform.BareMetal.Hosts[0].BGPPeers[0].HoldTimeSeconds = 90
	ic.Config.Platform.BareMetal.Hosts[0].BGPPeers[0].KeepaliveTimeSeconds = 30
	ic.Config.Platform.BareMetal.Hosts[0].BGPPeers[0].FailureDetection = baremetal.BGPFailureDetectionHoldTimer

	cm := &BGPVIPConfigMap{}
	parents := asset.Parents{}
	parents.Add(ic)
	assert.NoError(t, cm.Generate(context.Background(), parents))

	var data bgpVIPConfigJSON
	assert.NoError(t, json.Unmarshal([]byte(cm.ConfigMap.Data["config.json"]), &data))
	assert.Equal(t, "90", data.DefaultPeers[0].HoldTime)
	assert.Equal(t, "30", data.DefaultPeers[0].KeepaliveTime)
	assert.Equal(t, "true", data.DefaultPeers[0].BFDEnabled)
	assert.Equal(t, "true", data.DefaultPeers[0].EBGPMultiHop)
	assert.Equal(t, "90", data.HostOverrides["master-0"][0].HoldTime)
	assert.Equal(t, "30", data.HostOverrides["master-0"][0].KeepaliveTime)
	// HoldTimer / DirectlyConnected / omitted map to omitted contract keys.
	assert.Empty(t, data.HostOverrides["master-0"][0].BFDEnabled)
	assert.Empty(t, data.HostOverrides["master-0"][0].EBGPMultiHop)

	// The raw payload must carry the stable contract keys, not the
	// install-config field names.
	raw := cm.ConfigMap.Data["config.json"]
	assert.Contains(t, raw, `"bfdEnabled"`)
	assert.Contains(t, raw, `"holdTime"`)
	assert.NotContains(t, raw, `"failureDetection"`)
	assert.NotContains(t, raw, `"holdTimeSeconds"`)
}
