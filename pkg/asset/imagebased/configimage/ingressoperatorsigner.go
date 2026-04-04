package configimage

import (
	"context"
	"crypto/x509/pkix"
	"fmt"
	"time"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	libcrypto "github.com/openshift/library-go/pkg/crypto"
)

// Name returns the human-friendly name of the asset.
func (a *IngressOperatorSignerCertKey) Name() string {
	return "Certificate (ingress-operator-signer)"
}

// IngressOperatorSignerCertKey is the asset that generates the ingress operator
// key/cert pair.
type IngressOperatorSignerCertKey struct {
	tls.SelfSignedCertKey
}

var _ asset.Asset = (*IngressOperatorSignerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair.
func (a *IngressOperatorSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{&installconfig.InstallConfig{}}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *IngressOperatorSignerCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	signerName := fmt.Sprintf("%s@%d", "ingress-operator", time.Now().Unix())

	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	keyGen := libcrypto.RSAKeyPairGenerator{Bits: 2048}
	tlsCfg, err := libcrypto.NewSigningCertificate(
		signerName,
		keyGen,
		libcrypto.WithLifetime(tls.ValidityOneYear(installConfig)*2),
		libcrypto.WithSubject(pkix.Name{CommonName: signerName}),
	)
	if err != nil {
		return fmt.Errorf("failed to generate ingress operator signer certificate: %w", err)
	}

	certPEM, keyPEM, err := tlsCfg.GetPEMBytes()
	if err != nil {
		return fmt.Errorf("failed to encode ingress operator signer to PEM: %w", err)
	}

	a.CertRaw = certPEM
	a.KeyRaw = keyPEM

	return nil
}

// IngressOperatorCABundle is the asset the generates the ingress-operator-signer-ca-bundle,
// which contains all the ingrees operator signer CA.
type IngressOperatorCABundle struct {
	tls.CertBundle
}

var _ asset.Asset = (*IngressOperatorCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *IngressOperatorCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&IngressOperatorSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *IngressOperatorCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	certs := []tls.CertInterface{}
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(tls.CertInterface))
	}
	return a.CertBundle.Generate(ctx, "ingress-operator-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *IngressOperatorCABundle) Name() string {
	return "Certificate (ingress-operator-ca-bundle)"
}
