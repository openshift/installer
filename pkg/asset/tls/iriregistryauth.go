package tls //nolint:revive // pre-existing package name

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/manifests"
)

const (
	// IRIRegistryUsername is the fixed username for IRI registry authentication.
	IRIRegistryUsername = "openshift"
	// PasswordBytes is the number of random bytes to generate for the password (256-bit entropy).
	PasswordBytes = 32
)

// IRIRegistryCredentials is the asset for the IRI registry authentication credentials.
// This is an in-memory-only asset: credentials are consumed by other assets
// (operators.go, bootstrap/common.go) but not written to disk.
//
// This must NOT write files to the auth/ directory. In agent-based installs,
// assisted-service moves kubeadmin-password and kubeconfig out of auth/ and
// then calls os.Remove("auth") to delete the directory. That call fails if
// any extra files remain, which would break the deployment. See:
// https://github.com/openshift/assisted-service/blob/89897ade7135/internal/ignition/installmanifests.go#L356
type IRIRegistryCredentials struct {
	Username        string
	Password        string //nolint:gosec // this is a credential holder, not a hardcoded secret
	HtpasswdContent string
}

var _ asset.Asset = (*IRIRegistryCredentials)(nil)

// Dependencies returns the dependencies for generating IRI registry auth.
func (a *IRIRegistryCredentials) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&manifests.InternalReleaseImage{},
	}
}

// Generate generates the IRI registry authentication credentials.
func (a *IRIRegistryCredentials) Generate(ctx context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	iri := &manifests.InternalReleaseImage{}
	dependencies.Get(installConfig, iri)

	// Only generate if NoRegistryClusterInstall feature is enabled
	if !installConfig.Config.EnabledFeatureGates().Enabled(features.FeatureGateNoRegistryClusterInstall) {
		return nil
	}

	// Skip if InternalReleaseImage manifest wasn't found
	if len(iri.FileList) == 0 {
		return nil
	}

	// Generate random password (32 bytes = 256-bit entropy)
	passwordBytes := make([]byte, PasswordBytes)
	if _, err := rand.Read(passwordBytes); err != nil {
		return fmt.Errorf("failed to generate random password: %w", err)
	}
	a.Password = base64.StdEncoding.EncodeToString(passwordBytes)
	a.Username = IRIRegistryUsername

	// Create bcrypt hash
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create htpasswd format: username:bcrypt-hash
	a.HtpasswdContent = fmt.Sprintf("%s:%s\n", a.Username, string(hash))

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *IRIRegistryCredentials) Name() string {
	return "IRI Registry Authentication"
}
