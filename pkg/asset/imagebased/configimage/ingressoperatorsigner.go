package configimage

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1" //nolint: gosec
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
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

	cfg := &tls.CertCfg{
		Subject:   pkix.Name{CommonName: signerName},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  tls.ValidityOneYear(installConfig) * 2,
		IsCA:      true,
	}

	key, crt, err := generateSelfSignedCertificate(cfg)
	if err != nil {
		return err
	}

	a.KeyRaw = tls.PrivateKeyToPem(key)
	a.CertRaw = tls.CertToPem(crt)

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

// selfSignedCertificate creates a self signed certificate.
//
// TODO: this is a modified tls.SelfSignedCertificate function to allow the
// certificate's Subject.OU to be empty. We should revisit this by sharing as
// much code as possible with the tls package.
//
// https://github.com/openshift/installer/blob/6723dfd18056a6d002f792afce5547fc24874908/pkg/asset/tls/tls.go
func selfSignedCertificate(cfg *tls.CertCfg, key *rsa.PrivateKey) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	cert := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  cfg.IsCA,
		KeyUsage:              cfg.KeyUsages,
		NotAfter:              time.Now().Add(cfg.Validity),
		NotBefore:             time.Now(),
		SerialNumber:          serial,
		Subject:               cfg.Subject,
	}
	// verifies that the CN for the cert is set
	if len(cfg.Subject.CommonName) == 0 {
		return nil, errors.New("certification's subject is not set, or invalid")
	}
	pub := key.Public()
	cert.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to set subject key identifier: %w", err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, key.Public(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}
	return x509.ParseCertificate(certBytes)
}

// generateSelfSignedCertificate generates a key/cert pair defined by CertCfg.
func generateSelfSignedCertificate(cfg *tls.CertCfg) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.PrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	crt, err := selfSignedCertificate(cfg, key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create self-signed certificate: %w", err)
	}
	return key, crt, nil
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// generateSubjectKeyID generates a SHA-1 hash of the subject public key.
func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var publicKeyBytes []byte
	var err error

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(rsaPublicKey{N: pub.N, E: pub.E})
		if err != nil {
			return nil, fmt.Errorf("failed to Marshal ans1 public key: %w", err)
		}
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y) //nolint: staticcheck
	default:
		return nil, errors.New("only RSA and ECDSA public keys supported")
	}

	hash := sha1.Sum(publicKeyBytes) //nolint: gosec
	return hash[:], nil
}
