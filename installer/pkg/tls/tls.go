package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
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
	cert := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:         true,
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
		DNSNames:     csr.DNSNames,
		ExtKeyUsage:  cfg.ExtKeyUsages,
		IPAddresses:  csr.IPAddresses,
		KeyUsage:     cfg.KeyUsages,
		NotAfter:     time.Now().Add(cfg.Validity),
		NotBefore:    caCert.NotBefore,
		SerialNumber: serial,
		Subject:      csr.Subject,
		IsCA:         true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, fmt.Errorf("error creating signed certificate: %v", err)
	}
	return x509.ParseCertificate(certBytes)
}
