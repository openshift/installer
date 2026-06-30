package tls

import (
	"context"

	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
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

// Load reads the install-config.yaml directly from disk and extracts the PKI
// configuration. Returns (true, nil) when the config is found and parsed, or
// (false, nil) when the file is missing or unparseable — never errors and
// never triggers interactive prompts.
//
// Returning (false, nil) when the file is missing allows the asset store to
// fall back to the state file between multi-step invocations (e.g. create
// manifests followed by create cluster, where install-config.yaml is consumed
// after the first step). When no state file exists either (e.g. agent flow),
// the asset store calls Generate() which leaves PKIConfig nil (RSA-2048).
//
// We parse the install-config directly rather than depending on the
// InstallConfig asset because store.load() recursively loads all dependencies
// from disk, and InstallConfig.Load() runs full platform validation via
// finish() -> ValidateInstallConfig(). In the agent flow, configs are validated
// by OptionalInstallConfig with agent-specific rules that are more lenient
// (e.g. vSphere without credentials is valid for agent installs). Adding
// InstallConfig as a dependency here would pull standard validation into the
// agent flow, rejecting configs that OptionalInstallConfig accepts.
//
// Since we only need the PKI and FeatureSet fields, we can safely unmarshal
// without validation. Any actual validation errors will be caught by the
// authoritative InstallConfig or OptionalInstallConfig asset elsewhere in the
// pipeline.
func (s *SignerKeyParams) Load(f asset.FileFetcher) (bool, error) {
	// Matches installConfigFilename in pkg/asset/installconfig/installconfig.go.
	file, err := f.FetchByName("install-config.yaml")
	if err != nil {
		return false, nil
	}
	config := &types.InstallConfig{}
	if err := yaml.Unmarshal(file.Data, config); err != nil {
		return false, nil
	}
	s.PKIConfig = pkidefaults.EffectiveSignerPKIConfig(config)
	return true, nil
}
