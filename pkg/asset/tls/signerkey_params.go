package tls //nolint:revive // pre-existing package name

import (
	"context"

	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	pkidefaults "github.com/openshift/installer/pkg/types/pki"
)

// SignerKeyParams resolves the effective PKI configuration for certificate
// generation. It has no asset dependencies and reads install-config.yaml
// directly from disk so that signer certs can be generated without triggering
// standard InstallConfig validation. When no install-config is present (e.g.
// agent create certificates, node-joiner add-nodes), Generate() leaves
// ConfigurablePKIEnabled false, which causes all cert assets to take the
// legacy RSA-2048 code path.
type SignerKeyParams struct {
	// Profile is the resolved PKI profile. Only meaningful when
	// ConfigurablePKIEnabled is true.
	Profile configv1alpha1.PKIProfile

	// ConfigurablePKIEnabled indicates whether the ConfigurablePKI feature
	// gate is active. When false, cert assets take the legacy code path.
	ConfigurablePKIEnabled bool
}

var _ asset.WritableAsset = (*SignerKeyParams)(nil)

// Name returns a human-friendly name for the asset.
func (*SignerKeyParams) Name() string {
	return "Signer Key Parameters"
}

// Dependencies returns no dependencies. See Load() for why this asset does
// not depend on the standard InstallConfig asset.
func (*SignerKeyParams) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate is the fallback when Load() finds no install-config on disk
// (e.g. agent flow). It leaves ConfigurablePKIEnabled false so all cert
// assets take the legacy RSA-2048 path.
func (s *SignerKeyParams) Generate(_ context.Context, _ asset.Parents) error {
	return nil
}

// Files returns nil — this asset has no on-disk representation.
func (*SignerKeyParams) Files() []*asset.File {
	return nil
}

// Load reads install-config.yaml through the standard LoadFromFile pipeline
// (strict YAML, deprecated field conversion, defaults) and extracts the
// effective PKI profile. Returns (false, nil) when the file is missing,
// allowing the asset store to fall back to the state file between
// multi-step invocations (e.g. create manifests followed by create cluster).
//
// This asset does not depend on the InstallConfig asset because that would
// pull platform validation and cloud API calls into codepaths that generate
// signer certs without credentials (agent create certificates, node-joiner).
func (s *SignerKeyParams) Load(f asset.FileFetcher) (bool, error) {
	base := &installconfig.AssetBase{}
	found, err := base.LoadFromFile(f)
	if !found || err != nil {
		return found, err
	}
	s.Profile, s.ConfigurablePKIEnabled = pkidefaults.EffectiveProfile(base.Config)
	return true, nil
}

// LegacyPKIConfig returns a *types.PKIConfig for the legacy code path.
// Returns nil when ConfigurablePKI is disabled (RSA-2048 default).
func (s *SignerKeyParams) LegacyPKIConfig() *types.PKIConfig {
	if !s.ConfigurablePKIEnabled {
		return nil
	}
	keyConfig := types.KeyConfig{
		Algorithm: types.KeyAlgorithm(s.Profile.SignerCertificates.Key.Algorithm),
	}
	switch keyConfig.Algorithm {
	case types.KeyAlgorithmRSA:
		keyConfig.RSA = &types.RSAKeyConfig{KeySize: s.Profile.SignerCertificates.Key.RSA.KeySize}
	case types.KeyAlgorithmECDSA:
		keyConfig.ECDSA = &types.ECDSAKeyConfig{Curve: types.ECDSACurve(s.Profile.SignerCertificates.Key.ECDSA.Curve)}
	}
	return &types.PKIConfig{
		SignerCertificates: types.CertificateConfig{Key: keyConfig},
	}
}
