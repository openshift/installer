package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	mcfgv1alpha1 "github.com/openshift/api/machineconfiguration/v1alpha1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestBGPVIPConfigFromInstallConfig(t *testing.T) {
	ic := &types.InstallConfig{Platform: types.Platform{BareMetal: &baremetal.Platform{
		APIVIPs:     []string{"192.168.111.5"},
		IngressVIPs: []string{"192.168.111.4"},
		BGPVIPConfig: &baremetal.BGPVIPConfig{
			LocalASN:    64512,
			Communities: []string{"64512:100", "64512:200"},
			Peers: []baremetal.BGPPeerConfig{{
				PeerAddress:          "192.168.111.1",
				PeerASN:              64513,
				Password:             "s3cret",
				Port:                 180,
				FailureDetection:     baremetal.BGPFailureDetectionBFD,
				PeerReachability:     baremetal.BGPPeerReachabilityMultiHop,
				HoldTimeSeconds:      90,
				KeepaliveTimeSeconds: 30,
			}},
		},
		Hosts: []*baremetal.Host{
			{Name: "master-0"},
			{Name: "worker-0", BGPPeers: []baremetal.BGPPeerConfig{{
				PeerAddress:      "fd2e::1",
				PeerASN:          64514,
				FailureDetection: baremetal.BGPFailureDetectionHoldTimer,
				PeerReachability: baremetal.BGPPeerReachabilityDirectlyConnected,
			}}},
		},
	}}}

	got, secrets := bgpVIPConfigFromInstallConfig(ic)

	assert.Equal(t, "cluster", got.Name)
	assert.Equal(t, "machineconfiguration.openshift.io/v1alpha1", got.APIVersion)
	assert.Equal(t, "BGPVIPConfig", got.Kind)
	assert.Equal(t, int64(64512), got.Spec.LocalASN)
	assert.Equal(t, []mcfgv1alpha1.BGPCommunity{"64512:100", "64512:200"}, got.Spec.Communities)

	if !assert.Len(t, got.Spec.DefaultPeers, 1) {
		return
	}
	p := got.Spec.DefaultPeers[0]
	assert.Equal(t, "192.168.111.1", p.PeerAddress)
	assert.Equal(t, int64(64513), p.PeerASN)
	assert.Equal(t, int32(180), p.Port)
	assert.Equal(t, mcfgv1alpha1.FailureDetectionBFD, p.FailureDetection)
	assert.Equal(t, mcfgv1alpha1.PeerReachabilityMultiHop, p.PeerReachability)
	assert.Equal(t, int32(90), p.HoldTimeSeconds)
	assert.Equal(t, int32(30), p.KeepaliveTimeSeconds)
	assert.Equal(t, "bgp-vip-peer-192-168-111-1", p.PasswordSecret.Name)

	if assert.Len(t, secrets, 1) {
		s := secrets[0]
		assert.Equal(t, "bgp-vip-peer-192-168-111-1", s.Name)
		assert.Equal(t, "openshift-config", s.Namespace)
		assert.Equal(t, corev1.SecretTypeBasicAuth, s.Type)
		assert.Equal(t, []byte("s3cret"), s.Data["password"])
	}

	// hosts without bgpPeers produce no override; the override peer keeps its
	// own enums and has no Secret.
	if !assert.Len(t, got.Spec.HostOverrides, 1) {
		return
	}
	override := got.Spec.HostOverrides[0]
	assert.Equal(t, "worker-0", override.Hostname)
	if assert.Len(t, override.Peers, 1) {
		assert.Equal(t, "fd2e::1", override.Peers[0].PeerAddress)
		assert.Equal(t, mcfgv1alpha1.FailureDetectionHoldTimer, override.Peers[0].FailureDetection)
		assert.Equal(t, mcfgv1alpha1.PeerReachabilityDirectlyConnected, override.Peers[0].PeerReachability)
		assert.Empty(t, override.Peers[0].PasswordSecret.Name)
	}
}

// Unset install-config fields stay unset in the CR: the API applies its
// defaults by omission (the CRD rejects "" enum values), and the Secret
// reference is only produced for peers that carry a password.
func TestBGPVIPConfigFromInstallConfigOmittedStaysOmitted(t *testing.T) {
	ic := &types.InstallConfig{Platform: types.Platform{BareMetal: &baremetal.Platform{
		BGPVIPConfig: &baremetal.BGPVIPConfig{
			LocalASN: 64512,
			Peers:    []baremetal.BGPPeerConfig{{PeerAddress: "192.168.111.1", PeerASN: 64513}},
		},
	}}}

	got, secrets := bgpVIPConfigFromInstallConfig(ic)

	p := got.Spec.DefaultPeers[0]
	assert.Equal(t, mcfgv1alpha1.FailureDetectionMode(""), p.FailureDetection)
	assert.Equal(t, mcfgv1alpha1.PeerReachabilityMode(""), p.PeerReachability)
	assert.Equal(t, int32(0), p.HoldTimeSeconds)
	assert.Equal(t, int32(0), p.KeepaliveTimeSeconds)
	assert.Equal(t, int32(0), p.Port)
	assert.Equal(t, mcfgv1alpha1.BGPVIPPasswordSecret{}, p.PasswordSecret)
	assert.Nil(t, got.Spec.Communities)
	assert.Empty(t, got.Spec.HostOverrides)
	assert.Empty(t, secrets)
}

func TestBGPSanitizeNameComponent(t *testing.T) {
	cases := []struct{ in, want string }{
		{"192.168.111.1", "192-168-111-1"},
		{"fd2e::1", "fd2e-1"},
		{"fd2e:6f44:5dd8::1", "fd2e-6f44-5dd8-1"},
		{"::1", "1"},
		{"2001:db8::", "2001-db8"},
		{"Master-0.Example.COM", "master-0-example-com"},
		{"worker--0", "worker-0"},
		{"-host-", "host"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, bgpSanitizeNameComponent(c.in), "input %q", c.in)
	}
}

func TestBGPPeerSecretName(t *testing.T) {
	assert.Equal(t, "bgp-vip-peer-192-168-111-1", bgpPeerSecretName("", "192.168.111.1"))
	assert.Equal(t, "bgp-vip-peer-master-0-fd2e-1", bgpPeerSecretName("master-0", "fd2e::1"))
}

// Peers with passwords in hostOverrides get hostname-scoped Secret names, so
// the same peer address under two hosts cannot collide.
func TestBGPVIPConfigOverridePeerSecretNaming(t *testing.T) {
	ic := &types.InstallConfig{Platform: types.Platform{BareMetal: &baremetal.Platform{
		BGPVIPConfig: &baremetal.BGPVIPConfig{
			LocalASN: 64512,
			Peers:    []baremetal.BGPPeerConfig{{PeerAddress: "192.168.111.1", PeerASN: 64513}},
		},
		Hosts: []*baremetal.Host{
			{Name: "Worker-0", BGPPeers: []baremetal.BGPPeerConfig{{PeerAddress: "fd2e::1", PeerASN: 64514, Password: "hunter2"}}},
			{Name: "worker-1", BGPPeers: []baremetal.BGPPeerConfig{{PeerAddress: "fd2e::1", PeerASN: 64514, Password: "hunter3"}}},
		},
	}}}

	got, secrets := bgpVIPConfigFromInstallConfig(ic)

	if assert.Len(t, secrets, 2) {
		assert.Equal(t, "bgp-vip-peer-worker-0-fd2e-1", secrets[0].Name)
		assert.Equal(t, "bgp-vip-peer-worker-1-fd2e-1", secrets[1].Name)
	}
	assert.Equal(t, "bgp-vip-peer-worker-0-fd2e-1", got.Spec.HostOverrides[0].Peers[0].PasswordSecret.Name)
	assert.Equal(t, "bgp-vip-peer-worker-1-fd2e-1", got.Spec.HostOverrides[1].Peers[0].PasswordSecret.Name)
}

// Without passwords the asset renders exactly one file: the CR, with no
// passwordSecret key (omitzero) and no VIPs.
func TestBGPVIPConfigCRAssetRendersCR(t *testing.T) {
	parents := asset.Parents{}
	parents.Add(bgpInstallConfig())

	a := &BGPVIPConfigCR{}
	if !assert.NoError(t, a.Generate(context.Background(), parents)) {
		return
	}
	files := a.Files()
	if !assert.Len(t, files, 1) {
		return
	}
	assert.Equal(t, "manifests/cluster-bgp-vip-02-config.yml", files[0].Filename)
	data := string(files[0].Data)
	assert.Contains(t, data, "kind: BGPVIPConfig")
	assert.Contains(t, data, "apiVersion: machineconfiguration.openshift.io/v1alpha1")
	assert.Contains(t, data, "name: cluster")
	assert.Contains(t, data, "localASN: 64512")
	assert.Contains(t, data, "hostname: master-0")
	assert.NotContains(t, data, "kind: ConfigMap")
	assert.NotContains(t, data, "apiVIPs")
	assert.NotContains(t, data, "ingressVIPs")
	assert.NotContains(t, data, "passwordSecret")
	assert.NotContains(t, data, "password")
}

// Peers with install-config passwords produce basic-auth Secret manifests
// referenced by name from the CR; the plaintext never appears in the CR.
func TestBGPVIPConfigCRAssetRendersPeerSecrets(t *testing.T) {
	ic := installconfig.MakeAsset(&types.InstallConfig{
		Platform: types.Platform{
			BareMetal: &baremetal.Platform{
				BGPVIPConfig: &baremetal.BGPVIPConfig{
					LocalASN: 64512,
					Peers: []baremetal.BGPPeerConfig{
						{PeerAddress: "192.168.111.1", PeerASN: 64513, Password: "s3cret"},
					},
				},
				Hosts: []*baremetal.Host{{
					Name: "master-0",
					BGPPeers: []baremetal.BGPPeerConfig{
						{PeerAddress: "fd2e::1", PeerASN: 64514, Password: "hunter2"},
					},
				}},
			},
		},
	})
	parents := asset.Parents{}
	parents.Add(ic)

	a := &BGPVIPConfigCR{}
	if !assert.NoError(t, a.Generate(context.Background(), parents)) {
		return
	}
	files := a.Files()
	if !assert.Len(t, files, 3) {
		return
	}
	byName := map[string]string{}
	for _, f := range files {
		byName[f.Filename] = string(f.Data)
	}

	cr := byName["manifests/cluster-bgp-vip-02-config.yml"]
	if !assert.NotEmpty(t, cr) {
		return
	}
	assert.Contains(t, cr, "passwordSecret:")
	assert.Contains(t, cr, "name: bgp-vip-peer-192-168-111-1")
	assert.Contains(t, cr, "name: bgp-vip-peer-master-0-fd2e-1")
	assert.NotContains(t, cr, "s3cret")
	assert.NotContains(t, cr, "hunter2")

	for name, want := range map[string]string{
		"bgp-vip-peer-192-168-111-1":   "s3cret",
		"bgp-vip-peer-master-0-fd2e-1": "hunter2",
	} {
		data := byName["manifests/"+name+"-secret.yaml"]
		if !assert.NotEmpty(t, data, "secret manifest for %s", name) {
			continue
		}
		var secret corev1.Secret
		if !assert.NoError(t, yaml.Unmarshal([]byte(data), &secret)) {
			continue
		}
		assert.Equal(t, name, secret.Name)
		assert.Equal(t, "openshift-config", secret.Namespace)
		assert.Equal(t, corev1.SecretTypeBasicAuth, secret.Type)
		assert.Equal(t, []byte(want), secret.Data["password"])
	}
}

// Non-baremetal platforms and baremetal without bgpVIPConfig render nothing.
func TestBGPVIPConfigCRAssetSkipsWhenNotConfigured(t *testing.T) {
	for name, ic := range map[string]*installconfig.InstallConfig{
		"baremetal without bgpVIPConfig": installconfig.MakeAsset(&types.InstallConfig{
			Platform: types.Platform{BareMetal: &baremetal.Platform{APIVIPs: []string{"192.168.111.5"}}},
		}),
		"no platform": installconfig.MakeAsset(&types.InstallConfig{}),
	} {
		parents := asset.Parents{}
		parents.Add(ic)
		a := &BGPVIPConfigCR{}
		assert.NoError(t, a.Generate(context.Background(), parents), name)
		assert.Empty(t, a.Files(), name)
	}
}
