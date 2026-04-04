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

var _ asset.Asset = (*SignerPKIConfig)(nil)

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

// Name returns the human-friendly name of the asset.
func (*SignerPKIConfig) Name() string {
	return "PKI Signer Configuration"
}
