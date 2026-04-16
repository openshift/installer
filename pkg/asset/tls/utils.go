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

	libcrypto "github.com/openshift/library-go/pkg/crypto"
)

// PrivateKeyToPem converts a private key (RSA or ECDSA) to PEM format.
//
// Deprecated: Use libcrypto.EncodeKey directly for new code.
func PrivateKeyToPem(key crypto.PrivateKey) ([]byte, error) {
	return libcrypto.EncodeKey(key)
}

// CertToPem converts an x509.Certificate object to a pem string.
//
// Deprecated: Use libcrypto.EncodeCertificates directly for new code.
func CertToPem(cert *x509.Certificate) []byte {
	b, err := libcrypto.EncodeCertificates(cert)
	if err != nil {
		return nil
	}
	return b
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
//
// Deprecated: Use libcrypto.GetCAFromBytes or libcrypto.GetTLSCertificateConfigFromBytes for new code.
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
//
// Deprecated: Use libcrypto.GetCAFromBytes or libcrypto.GetTLSCertificateConfigFromBytes for new code.
func PemToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}
