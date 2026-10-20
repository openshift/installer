package manifests

import (
	"context"
	"encoding/json"
	"path"
	"strconv"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/baremetal"
)

var (
	bgpVIPConfigMapFileName = path.Join(manifestDir, "bgp-vip-config.yaml")
)

const (
	bgpVIPConfigMapName      = "bgp-vip-config"
	bgpVIPConfigMapNamespace = "openshift-network-operator"
	bgpVIPConfigMapDataKey   = "config.json"
)

// bgpVIPConfigJSON is the JSON structure stored in the ConfigMap. Its field
// names match baremetal-runtimecfg's FRRPeerMapping so the same JSON can be
// written verbatim to /etc/kubernetes/static-pod-resources/frr-k8s/frr-peers.json.
type bgpVIPConfigJSON struct {
	LocalASN      int64                    `json:"localASN"`
	DefaultPeers  []frrPeerJSON            `json:"defaultPeers"`
	Communities   []string                 `json:"communities,omitempty"`
	APIVIPs       []string                 `json:"apiVIPs"`
	IngressVIPs   []string                 `json:"ingressVIPs"`
	HostOverrides map[string][]frrPeerJSON `json:"hostOverrides,omitempty"`
}

// frrPeerJSON is one peer in the ConfigMap payload. It is a dedicated type
// (not baremetal.BGPPeerConfig) because the install-config surface and this
// contract evolve independently: the contract keeps the boolean-string and
// bare-seconds shapes the merged runtimecfg/MCO/CNO consumers parse.
type frrPeerJSON struct {
	PeerAddress   string `json:"peerAddress"`
	PeerASN       int64  `json:"peerASN"`
	Password      string `json:"password,omitempty"` //nolint:gosec // BGP TCP MD5 session password carried to the renderers, not a hardcoded credential
	Port          int32  `json:"port,omitempty"`
	BFDEnabled    string `json:"bfdEnabled,omitempty"`
	EBGPMultiHop  string `json:"ebgpMultiHop,omitempty"`
	HoldTime      string `json:"holdTime,omitempty"`
	KeepaliveTime string `json:"keepaliveTime,omitempty"`
}

// toFRRPeers maps the install-config peers onto the ConfigMap contract:
// the failureDetection/peerReachability enums become the boolean strings
// the consumers parse, and the integer-second timers become the bare-second
// decimal strings FRR's "timers <keepalive> <hold>" takes.
func toFRRPeers(peers []baremetal.BGPPeerConfig) []frrPeerJSON {
	out := make([]frrPeerJSON, 0, len(peers))
	for _, p := range peers {
		j := frrPeerJSON{
			PeerAddress: p.PeerAddress,
			PeerASN:     p.PeerASN,
			Password:    p.Password,
			Port:        p.Port,
		}
		if p.FailureDetection == baremetal.BGPFailureDetectionBFD {
			j.BFDEnabled = "true"
		}
		if p.PeerReachability == baremetal.BGPPeerReachabilityMultiHop {
			j.EBGPMultiHop = "true"
		}
		if p.HoldTimeSeconds != 0 {
			j.HoldTime = strconv.FormatInt(int64(p.HoldTimeSeconds), 10)
		}
		if p.KeepaliveTimeSeconds != 0 {
			j.KeepaliveTime = strconv.FormatInt(int64(p.KeepaliveTimeSeconds), 10)
		}
		out = append(out, j)
	}
	return out
}

// BGPVIPConfigMap generates the bgp-vip-config ConfigMap for CNO.
type BGPVIPConfigMap struct {
	ConfigMap *corev1.ConfigMap
	File      *asset.File
}

var _ asset.WritableAsset = (*BGPVIPConfigMap)(nil)

// Name returns a human friendly name for the asset.
func (*BGPVIPConfigMap) Name() string {
	return "BGP VIP Config ConfigMap"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*BGPVIPConfigMap) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the BGP VIP Config ConfigMap.
func (b *BGPVIPConfigMap) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	// Only generate for baremetal platform with BGPVIPConfig set.
	if installConfig.Config.Platform.Name() != baremetal.Name ||
		installConfig.Config.Platform.BareMetal == nil ||
		installConfig.Config.Platform.BareMetal.BGPVIPConfig == nil {
		return nil
	}

	bm := installConfig.Config.Platform.BareMetal
	bgpConfig := bm.BGPVIPConfig

	configData := bgpVIPConfigJSON{
		LocalASN:     bgpConfig.LocalASN,
		DefaultPeers: toFRRPeers(bgpConfig.Peers),
		Communities:  bgpConfig.Communities,
		APIVIPs:      bm.APIVIPs,
		IngressVIPs:  bm.IngressVIPs,
	}

	// Collect per-host BGP peer overrides.
	hostOverrides := make(map[string][]frrPeerJSON)
	for _, host := range bm.Hosts {
		if host != nil && len(host.BGPPeers) > 0 {
			hostOverrides[host.Name] = toFRRPeers(host.BGPPeers)
		}
	}
	if len(hostOverrides) > 0 {
		configData.HostOverrides = hostOverrides
	}

	jsonBytes, err := json.Marshal(configData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal BGP VIP config JSON")
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: bgpVIPConfigMapNamespace,
			Name:      bgpVIPConfigMapName,
		},
		Data: map[string]string{
			bgpVIPConfigMapDataKey: string(jsonBytes),
		},
	}

	cmData, err := yaml.Marshal(cm)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifest", b.Name())
	}
	b.ConfigMap = cm
	b.File = &asset.File{
		Filename: bgpVIPConfigMapFileName,
		Data:     cmData,
	}
	return nil
}

// Files returns the files generated by the asset.
func (b *BGPVIPConfigMap) Files() []*asset.File {
	if b.File != nil {
		return []*asset.File{b.File}
	}
	return []*asset.File{}
}

// Load loads the already-rendered files back from disk.
func (b *BGPVIPConfigMap) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
