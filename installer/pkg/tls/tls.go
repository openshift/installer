package tls

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"time"
)

const (
	keySize = 2048
)

// CertCfg contains all needed fields to configure a new certificate
type CertCfg struct {
	DNSNames     []string
	ExtKeyUsages []x509.ExtKeyUsage
	IPAddresses  []net.IP
	KeyUsages    x509.KeyUsage
	Subject      pkix.Name
	Validity     time.Duration
	IsCA         bool
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// PrivateKey generates an RSA Private key and returns the value
func PrivateKey() (*rsa.PrivateKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("error generating RSA private key: %v", err)
	}

	return rsaKey, nil
}

// SelfSignedCACert Creates a self signed CA certificate
func SelfSignedCACert(cfg *CertCfg, key *rsa.PrivateKey) (*x509.Certificate, error) {
	var err error

	cert := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:         cfg.IsCA,
		KeyUsage:     cfg.KeyUsages,
		NotAfter:     time.Now().Add(cfg.Validity),
		NotBefore:    time.Now(),
		SerialNumber: new(big.Int).SetInt64(0),
		Subject:      cfg.Subject,
	}
	// verifies that the CN and/or OU for the cert is set
	if len(cfg.Subject.CommonName) == 0 || len(cfg.Subject.OrganizationalUnit) == 0 {
		return nil, fmt.Errorf("certification's subject is not set, or invalid")
	}
	pub := key.Public()
	cert.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to set subject key identifier: %v", err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, key.Public(), key)
	if err != nil {
		return nil, fmt.Errorf("error creating certificate: %v", err)
	}
	return x509.ParseCertificate(certBytes)
}

// SignedCertificate creates a new X.509 certificate based on a template.
func SignedCertificate(
	cfg *CertCfg,
	csr *x509.CertificateRequest,
	key *rsa.PrivateKey,
	caCert *x509.Certificate,
	caKey *rsa.PrivateKey,
) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	certTmpl := x509.Certificate{
		DNSNames:              csr.DNSNames,
		ExtKeyUsage:           cfg.ExtKeyUsages,
		IPAddresses:           csr.IPAddresses,
		KeyUsage:              cfg.KeyUsages,
		NotAfter:              time.Now().Add(cfg.Validity),
		NotBefore:             caCert.NotBefore,
		SerialNumber:          serial,
		Subject:               csr.Subject,
		IsCA:                  cfg.IsCA,
		Version:               3,
		BasicConstraintsValid: true,
	}
	pub := caCert.PublicKey.(*rsa.PublicKey)
	certTmpl.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to set subject key identifier: %v", err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, fmt.Errorf("error creating signed certificate: %v", err)
	}
	return x509.ParseCertificate(certBytes)
}

// generateSubjectKeyID generates a SHA-1 hash of the subject public key.
func generateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var publicKeyBytes []byte
	var err error

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(rsaPublicKey{N: pub.N, E: pub.E})
		if err != nil {
			return nil, err
		}
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	default:
		return nil, errors.New("only RSA and ECDSA public keys supported")
	}

	hash := sha1.Sum(publicKeyBytes)
	return hash[:], nil
}
