package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	libpki "github.com/openshift/library-go/pkg/pki"
)

// JournalCertKey is the asset that generates the key/cert pair for
// journal-gatewayd on the bootstrap node. This is a peer certificate:
// journal-gatewayd uses it to serve HTTPS (ServerAuth), and the
// installer uses the same cert/key to authenticate as a client when
// streaming logs via curl (ClientAuth).
type JournalCertKey struct {
	SignedCertKey
}

var _ asset.WritableAsset = (*JournalCertKey)(nil)

// Dependencies returns the dependencies.
func (a *JournalCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
		&SignerKeyParams{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *JournalCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &RootCA{}
	pkiCfg := &SignerKeyParams{}
	dependencies.Get(ca, pkiCfg)

	if !pkiCfg.ConfigurablePKIEnabled {
		cfg := &CertCfg{
			Subject:      pkix.Name{CommonName: "journal-gatewayd", Organization: []string{"OpenShift Bootstrap"}},
			KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			Validity:     ValidityTenYears(),
		}
		return a.SignedCertKey.Generate(ctx, cfg, ca, "journal-gatewayd", DoNotAppendParent, nil)
	}

	keyGen, err := resolveKeyGen(pkiCfg, libpki.CertificateTypePeer, "installer.journal-gateway")
	if err != nil {
		return err
	}
	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "journal-gatewayd", Organization: []string{"OpenShift Bootstrap"}},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears(),
		DNSNames:     []string{"journal-gatewayd"},
		CertType:     libpki.CertificateTypePeer,
	}
	return a.SignedCertKey.Generate(ctx, cfg, ca, "journal-gatewayd", DoNotAppendParent, keyGen)
}

// Name returns the human-friendly name of the asset.
func (a *JournalCertKey) Name() string {
	return "Certificate (journal-gatewayd)"
}
