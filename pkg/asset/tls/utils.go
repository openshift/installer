package tls

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKeyToPem converts an rsa.PrivateKey object to pem string
func PrivateKeyToPem(key *rsa.PrivateKey) string {
	keyInBytes := x509.MarshalPKCS1PrivateKey(key)
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyInBytes,
		},
	)
	return string(keyinPem)
}

// CertToPem converts an x509.Certificate object to a pem string
func CertToPem(cert *x509.Certificate) string {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return string(certInPem)
}

// CSRToPem converts an x509.CertificateRequest to a pem string
func CSRToPem(cert *x509.CertificateRequest) string {
	certInPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: cert.Raw,
		},
	)
	return string(certInPem)
}

// PublicKeyToPem converts an rsa.PublicKey object to pem string
func PublicKeyToPem(key *rsa.PublicKey) (string, error) {
	keyInBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyInBytes,
		},
	)
	return string(keyinPem), nil
}

// PemToPrivateKey converts a data block to rsa.PrivateKey.
func PemToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// PemToCertificate converts a data block to x509.Certificate.
func PemToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	return x509.ParseCertificate(block.Bytes)
}
