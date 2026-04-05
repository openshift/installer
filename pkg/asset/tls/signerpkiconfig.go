package tls

import (
	"context"

	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	features "github.com/openshift/api/features"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	pkidefaults "github.com/openshift/installer/pkg/types/pki"
)

// SignerPKIConfig resolves the effective PKI configuration for certificate generation.
// Profile is always non-nil — when the feature gate is off, it defaults to RSA 2048.
type SignerPKIConfig struct {
	// Profile is the resolved PKI profile. Never nil.
	// Contains both signer and leaf cert defaults.
	Profile configv1alpha1.PKIProfile

	// ConfigurablePKIEnabled indicates whether the ConfigurablePKI feature gate is active.
	// When true, the PKI CR manifest should be generated.
	// When false, the installer is using legacy defaults (RSA 2048)
	// and no PKI CR is emitted.
	ConfigurablePKIEnabled bool
}

var _ asset.WritableAsset = (*SignerPKIConfig)(nil)

// Dependencies returns the dependency of the SignerPKIConfig.
func (*SignerPKIConfig) Dependencies() []asset.Asset {
	return []asset.Asset{&installconfig.InstallConfig{}}
}

// Generate resolves the effective PKI profile from the install config.
func (p *SignerPKIConfig) Generate(_ context.Context, parents asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	parents.Get(ic)

	p.ConfigurablePKIEnabled = ic.Config.Enabled(features.FeatureGateConfigurablePKI)
	p.Profile = pkidefaults.EffectiveSignerPKIProfile(ic.Config)

	return nil
}

// Files returns no files — this asset is not written to disk.
func (p *SignerPKIConfig) Files() []*asset.File {
	return nil
}

// Load returns the default RSA 2048 profile when no install-config is available.
// This allows commands like "agent create certificates" to generate signer certs
// without requiring an install-config.yaml, preserving backward compatibility.
// When an install-config is present, Load returns false and Generate is used instead.
func (p *SignerPKIConfig) Load(f asset.FileFetcher) (bool, error) {
	// Check if install-config exists on disk. If it does, return false
	// so the normal Generate path runs with the full dependency chain.
	_, err := f.FetchByName("install-config.yaml")
	if err == nil {
		return false, nil
	}

	// No install-config available — use legacy RSA 2048 defaults.
	// This is the same behavior as before SignerPKIConfig was introduced.
	rsa2048 := configv1alpha1.KeyConfig{
		Algorithm: configv1alpha1.KeyAlgorithmRSA,
		RSA:       configv1alpha1.RSAKeyConfig{KeySize: 2048},
	}
	p.Profile = configv1alpha1.PKIProfile{
		Defaults: configv1alpha1.DefaultCertificateConfig{
			Key: rsa2048,
		},
		SignerCertificates: configv1alpha1.CertificateConfig{
			Key: rsa2048,
		},
	}
	p.ConfigurablePKIEnabled = false
	return true, nil
}

// Name returns the human-friendly name of the asset.
func (*SignerPKIConfig) Name() string {
	return "PKI Signer Configuration"
}
