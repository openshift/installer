package tlsutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math"
	"math/big"
	"net"
	"time"
)

// Certificate and key constants.
const (
	RSAKeySize   = 2048
	Duration365d = time.Hour * 24 * 365
)

// CertConfig is the TLS distinguished name configuration.
type CertConfig struct {
	CommonName   string
	Organization []string
	AltNames     AltNames
}

// AltNames represent TLS Subject Alternative Names.
type AltNames struct {
	DNSNames []string
	IPs      []net.IP
}

// NewPrivateKey returns a new private key.
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSAKeySize)
}

// ParsePEMEncodedCert parses raw certificate bytes and returns an x509
// Certificate.
func ParsePEMEncodedCert(pemdata []byte) (*x509.Certificate, error) {
	decoded, _ := pem.Decode(pemdata)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParseCertificate(decoded.Bytes)
}

// ParsePEMEncodedPrivateKey parses raw private keys and returns a private
// key.
func ParsePEMEncodedPrivateKey(pemdata []byte) (*rsa.PrivateKey, error) {
	decoded, _ := pem.Decode(pemdata)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParsePKCS1PrivateKey(decoded.Bytes)
}

// EncodeCertificatePEM returns encoded bytes for the given Certificate.
func EncodeCertificatePEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// EncodePrivateKeyPEM returns encoded bytes of the given private key.
func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

// NewServerCertificate returns a new x509 server certificate, signed by the
// CA with the given certificate and key.
func NewServerCertificate(cfg CertConfig, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, customDuration *time.Duration) (*x509.Certificate, error) {
	return newCertificate(cfg, key, caCert, caKey, customDuration, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})
}

// NewClientCertificate returns a new x509 client certificate, signed by the
// CA with the given certificate and key.
func NewClientCertificate(cfg CertConfig, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, customDuration *time.Duration) (*x509.Certificate, error) {
	return newCertificate(cfg, key, caCert, caKey, customDuration, []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth})
}

func newCertificate(cfg CertConfig, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey, customDuration *time.Duration, usages []x509.ExtKeyUsage) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}

	dur := Duration365d
	if customDuration != nil {
		dur = *customDuration
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: caCert.Subject.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     time.Now().Add(dur),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  usages,
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}
