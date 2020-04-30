package tls

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

// PrivateKeyToPem converts an rsa.PrivateKey object to pem string
func PrivateKeyToPem(key *rsa.PrivateKey) []byte {
	keyInBytes := x509.MarshalPKCS1PrivateKey(key)
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyInBytes,
		},
	)
	return keyinPem
}

// CertToPem converts an x509.Certificate object to a pem string
func CertToPem(cert *x509.Certificate) []byte {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return certInPem
}

// CSRToPem converts an x509.CertificateRequest to a pem string
func CSRToPem(cert *x509.CertificateRequest) []byte {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: cert.Raw,
		},
	)
	return certInPem
}

// PublicKeyToPem converts an rsa.PublicKey object to pem string
func PublicKeyToPem(key *rsa.PublicKey) ([]byte, error) {
	keyInBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to MarshalPKIXPublicKey")
	}
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyInBytes,
		},
	)
	return keyinPem, nil
}

// PemToPrivateKey converts a data block to rsa.PrivateKey.
func PemToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// PemToPublicKey converts a data block to rsa.PublicKey.
func PemToPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the public key")
	}
	obji, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := obji.(*rsa.PublicKey)
	if !ok {
		return nil, errors.Errorf("invalid public key format, expected RSA")
	}
	return publicKey, nil
}

// PemToCertificate converts a data block to x509.Certificate.
func PemToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}
