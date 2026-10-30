package manifests

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	mcfgv1alpha1 "github.com/openshift/api/machineconfiguration/v1alpha1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

var (
	bgpVIPConfigCRFileName = path.Join(manifestDir, "cluster-bgp-vip-02-config.yml")

	bgpNameSanitizeRE = regexp.MustCompile(`[:.]+`)
	bgpNameSqueezeRE  = regexp.MustCompile(`-+`)
)

const (
	bgpVIPConfigCRName       = "cluster"
	bgpPeerSecretNamespace   = "openshift-config" //nolint:gosec // namespace the peer password Secrets are created in, not a credential
	bgpPeerSecretPasswordKey = "password"         //nolint:gosec // key name inside the basic-auth Secret, not a credential
)

// bgpSanitizeNameComponent turns an address or hostname into an RFC 1123
// label fragment: lowercase, ':' and '.' replaced by '-', repeated '-'
// squeezed, leading/trailing '-' trimmed.
func bgpSanitizeNameComponent(s string) string {
	s = strings.ToLower(s)
	s = bgpNameSanitizeRE.ReplaceAllString(s, "-")
	s = bgpNameSqueezeRE.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// bgpPeerSecretName returns the deterministic Secret name for a peer's BGP
// password. hostname is empty for default peers.
func bgpPeerSecretName(hostname, address string) string {
	if hostname == "" {
		return fmt.Sprintf("bgp-vip-peer-%s", bgpSanitizeNameComponent(address))
	}
	return fmt.Sprintf("bgp-vip-peer-%s-%s", bgpSanitizeNameComponent(hostname), bgpSanitizeNameComponent(address))
}

// bgpPeerPasswordSecret builds the kubernetes.io/basic-auth Secret carrying
// the peer's BGP password in the openshift-config namespace, where the
// BGPVIPConfig consumers resolve passwordSecret references.
func bgpPeerPasswordSecret(name, password string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: bgpPeerSecretNamespace,
		},
		Type: corev1.SecretTypeBasicAuth,
		Data: map[string][]byte{
			bgpPeerSecretPasswordKey: []byte(password),
		},
	}
}

// bgpFailureDetectionMode maps the install-config enum onto the CRD enum. An
// unset value stays unset: the API defaults by omission, not by value.
func bgpFailureDetectionMode(v baremetal.BGPFailureDetection) mcfgv1alpha1.FailureDetectionMode {
	switch v {
	case baremetal.BGPFailureDetectionBFD:
		return mcfgv1alpha1.FailureDetectionBFD
	case baremetal.BGPFailureDetectionHoldTimer:
		return mcfgv1alpha1.FailureDetectionHoldTimer
	}
	return ""
}

// bgpPeerReachabilityMode maps the install-config enum onto the CRD enum,
// keeping unset values unset.
func bgpPeerReachabilityMode(v baremetal.BGPPeerReachability) mcfgv1alpha1.PeerReachabilityMode {
	switch v {
	case baremetal.BGPPeerReachabilityDirectlyConnected:
		return mcfgv1alpha1.PeerReachabilityDirectlyConnected
	case baremetal.BGPPeerReachabilityMultiHop:
		return mcfgv1alpha1.PeerReachabilityMultiHop
	}
	return ""
}

// bgpPeerToCR converts one install-config peer into a CRD peer. An inline
// install-config password is moved into a Secret appended to secrets; the
// CRD peer carries only the reference.
func bgpPeerToCR(p baremetal.BGPPeerConfig, hostname string, secrets *[]*corev1.Secret) mcfgv1alpha1.BGPVIPPeer {
	peer := mcfgv1alpha1.BGPVIPPeer{
		PeerAddress:          p.PeerAddress,
		PeerASN:              p.PeerASN,
		Port:                 p.Port,
		FailureDetection:     bgpFailureDetectionMode(p.FailureDetection),
		PeerReachability:     bgpPeerReachabilityMode(p.PeerReachability),
		HoldTimeSeconds:      p.HoldTimeSeconds,
		KeepaliveTimeSeconds: p.KeepaliveTimeSeconds,
	}
	if p.Password != "" {
		name := bgpPeerSecretName(hostname, p.PeerAddress)
		peer.PasswordSecret = mcfgv1alpha1.BGPVIPPasswordSecret{Name: name}
		*secrets = append(*secrets, bgpPeerPasswordSecret(name, p.Password))
	}
	return peer
}

// bgpVIPConfigFromInstallConfig converts the baremetal platform's
// bgpVIPConfig install-config section into the machineconfiguration
// BGPVIPConfig singleton named "cluster", plus the peer password Secrets it
// references. The VIPs themselves are not part of the CR (they live in the
// Infrastructure status); install-config validation already enforces the
// CRD's limits (1-16 peers, canonical addresses, community ranges), so the
// conversion is a plain mapping.
func bgpVIPConfigFromInstallConfig(ic *types.InstallConfig) (*mcfgv1alpha1.BGPVIPConfig, []*corev1.Secret) {
	bm := ic.Platform.BareMetal
	bgpConfig := bm.BGPVIPConfig
	var secrets []*corev1.Secret

	cr := &mcfgv1alpha1.BGPVIPConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1alpha1.GroupVersion.String(),
			Kind:       "BGPVIPConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: bgpVIPConfigCRName,
		},
		Spec: mcfgv1alpha1.BGPVIPConfigSpec{
			LocalASN: bgpConfig.LocalASN,
		},
	}
	for _, c := range bgpConfig.Communities {
		cr.Spec.Communities = append(cr.Spec.Communities, mcfgv1alpha1.BGPCommunity(c))
	}
	for _, p := range bgpConfig.Peers {
		cr.Spec.DefaultPeers = append(cr.Spec.DefaultPeers, bgpPeerToCR(p, "", &secrets))
	}
	for _, host := range bm.Hosts {
		if host == nil || len(host.BGPPeers) == 0 {
			continue
		}
		override := mcfgv1alpha1.BGPVIPHostPeers{Hostname: host.Name}
		for _, p := range host.BGPPeers {
			override.Peers = append(override.Peers, bgpPeerToCR(p, host.Name, &secrets))
		}
		cr.Spec.HostOverrides = append(cr.Spec.HostOverrides, override)
	}

	return cr, secrets
}

// BGPVIPConfigCR generates the BGPVIPConfig custom resource and the peer
// password Secrets for BGP-based VIP management. It is emitted alongside the
// bgp-vip-config ConfigMap (BGPVIPConfigMap) while the operators move from
// the ConfigMap to the CR; the ConfigMap goes away once they have.
type BGPVIPConfigCR struct {
	Config   *mcfgv1alpha1.BGPVIPConfig
	FileList []*asset.File
}

var _ asset.WritableAsset = (*BGPVIPConfigCR)(nil)

// Name returns a human friendly name for the asset.
func (*BGPVIPConfigCR) Name() string {
	return "BGP VIP Config CR"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*BGPVIPConfigCR) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the BGPVIPConfig CR manifest and the peer password
// Secret manifests.
func (b *BGPVIPConfigCR) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	if installConfig.Config.Platform.Name() != baremetal.Name ||
		installConfig.Config.Platform.BareMetal == nil ||
		installConfig.Config.Platform.BareMetal.BGPVIPConfig == nil {
		return nil
	}

	cr, secrets := bgpVIPConfigFromInstallConfig(installConfig.Config)

	crData, err := yaml.Marshal(cr)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifest", b.Name())
	}
	b.Config = cr
	b.FileList = []*asset.File{{
		Filename: bgpVIPConfigCRFileName,
		Data:     crData,
	}}
	for _, secret := range secrets {
		secretData, err := yaml.Marshal(secret)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s peer password Secret manifest", b.Name())
		}
		b.FileList = append(b.FileList, &asset.File{
			Filename: path.Join(manifestDir, secret.Name+"-secret.yaml"),
			Data:     secretData,
		})
	}
	asset.SortFiles(b.FileList)
	return nil
}

// Files returns the files generated by the asset.
func (b *BGPVIPConfigCR) Files() []*asset.File {
	return b.FileList
}

// Load returns false: the manifests are regenerated from the install-config,
// and the Manifests asset loads user-edited copies from disk.
func (b *BGPVIPConfigCR) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
