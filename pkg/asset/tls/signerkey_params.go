package tls //nolint:revive // pre-existing package name

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	pkidefaults "github.com/openshift/installer/pkg/types/pki"
)

// SignerKeyParams resolves the effective PKI configuration for signer
// certificates. It has no asset dependencies and reads install-config.yaml
// directly from disk so that signer certs can be generated without triggering
// standard InstallConfig validation. When no install-config is present (e.g.
// agent create certificates, node-joiner add-nodes), it defaults to nil
// PKIConfig which maps to RSA-2048.
type SignerKeyParams struct {
	PKIConfig *types.PKIConfig
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

// Generate is a no-op that leaves PKIConfig as nil (RSA-2048). This is the
// fallback when Load() finds no install-config on disk (e.g. agent flow).
func (s *SignerKeyParams) Generate(_ context.Context, _ asset.Parents) error {
	return nil
}

// Files returns nil — this asset has no on-disk representation.
func (*SignerKeyParams) Files() []*asset.File {
	return nil
}

// Load reads install-config.yaml through the standard LoadFromFile pipeline
// (strict YAML, deprecated field conversion, defaults) and extracts the
// effective PKI config. Returns (false, nil) when the file is missing,
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
	s.PKIConfig = pkidefaults.EffectiveSignerPKIConfig(base.Config)
	return true, nil
}
