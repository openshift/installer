package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// JournalCertKey is the asset that generates the key/cert pair that is used to
// authenticate with journal-gatewayd on the bootstrap node.
type JournalCertKey struct {
	SignedCertKey
}

var _ asset.WritableAsset = (*JournalCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *JournalCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *JournalCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &RootCA{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "journal-gatewayd", Organization: []string{"OpenShift Bootstrap"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "journal-gatewayd", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *JournalCertKey) Name() string {
	return "Certificate (journal-gatewayd)"
}
