package tls

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
)

// ServiceAccountKeyPair is the asset that generates the service-account public/private key pair.
type ServiceAccountKeyPair struct {
	KeyPair
}

var _ asset.WritableAsset = (*ServiceAccountKeyPair)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *ServiceAccountKeyPair) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *ServiceAccountKeyPair) Generate(ctx context.Context, dependencies asset.Parents) error {
	return a.KeyPair.Generate(ctx, "service-account")
}

// Name returns the human-friendly name of the asset.
func (a *ServiceAccountKeyPair) Name() string {
	return "Key Pair (service-account.pub)"
}

// Load is a no-op because the service account keypair is not written to disk.
func (a *ServiceAccountKeyPair) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
