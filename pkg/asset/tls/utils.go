package tls

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PrivateKeyToPem converts a private key (RSA or ECDSA) to PEM format.
func PrivateKeyToPem(key crypto.PrivateKey) ([]byte, error) {
	var block *pem.Block

	switch k := key.(type) {
	case *rsa.PrivateKey:
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ECDSA private key: %w", err)
		}
		block = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: bytes,
		}
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", key)
	}

	return pem.EncodeToMemory(block), nil
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
		logrus.Debugf("Failed to marshal PKIX public key: %s", err)
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

// PemToPrivateKey converts a PEM data block to a private key (RSA or ECDSA).
func PemToPrivateKey(data []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the private key")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		switch key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("unsupported PKCS#8 key type: %T", key)
		}
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
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
