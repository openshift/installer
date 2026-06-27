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
	"fmt"
	"math"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

const (
	// DefaultRSAKeySize is the default RSA key size used when PKI config is not specified.
	DefaultRSAKeySize int32 = 2048
	// DefaultKeyAlgorithm is the default key algorithm used when PKI config is not specified.
	DefaultKeyAlgorithm = types.KeyAlgorithmRSA
)

// PrivateKeyParams specifies parameters for private key generation.
type PrivateKeyParams struct {
	Algorithm  types.KeyAlgorithm
	RSAKeySize int32
	ECDSACurve types.ECDSACurve
}

// PKIConfigToKeyParams converts PKI config to key generation parameters.
// If pkiConfig is nil, returns default RSA 2048 parameters.
func PKIConfigToKeyParams(pkiConfig *types.PKIConfig) PrivateKeyParams {
	if pkiConfig == nil {
		return PrivateKeyParams{
			Algorithm:  DefaultKeyAlgorithm,
			RSAKeySize: DefaultRSAKeySize,
		}
	}

	keyConfig := pkiConfig.SignerCertificates.Key
	params := PrivateKeyParams{
		Algorithm: keyConfig.Algorithm,
	}

	switch keyConfig.Algorithm {
	case types.KeyAlgorithmRSA:
		params.RSAKeySize = keyConfig.RSA.KeySize
	case types.KeyAlgorithmECDSA:
		params.ECDSACurve = keyConfig.ECDSA.Curve
	}

	return params
}

// GeneratePrivateKeyWithParams generates a private key with the specified parameters.
func GeneratePrivateKeyWithParams(params PrivateKeyParams) (crypto.PrivateKey, error) {
	switch params.Algorithm {
	case types.KeyAlgorithmRSA:
		return GenerateRSAPrivateKey(params.RSAKeySize)
	case types.KeyAlgorithmECDSA:
		return GenerateECDSAPrivateKey(params.ECDSACurve)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", params.Algorithm)
	}
}

// GenerateRSAPrivateKey generates an RSA private key with the specified size.
func GenerateRSAPrivateKey(keySize int32) (*rsa.PrivateKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, int(keySize))
	if err != nil {
		return nil, errors.Wrap(err, "error generating RSA private key")
	}
	return rsaKey, nil
}

// GenerateECDSAPrivateKey generates an ECDSA private key with the specified curve.
func GenerateECDSAPrivateKey(curve types.ECDSACurve) (*ecdsa.PrivateKey, error) {
	var c elliptic.Curve

	switch curve {
	case types.ECDSACurveP256:
		c = elliptic.P256()
	case types.ECDSACurveP384:
		c = elliptic.P384()
	case types.ECDSACurveP521:
		c = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported ECDSA curve: %s", curve)
	}

	ecdsaKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "error generating ECDSA private key")
	}
	return ecdsaKey, nil
}

// keyUsageForAlgorithm returns appropriate x509.KeyUsage flags for the given algorithm.
// ECDSA keys can only perform digital signatures — they cannot perform key encipherment.
// RSA keys support both digital signatures and key encipherment.
func keyUsageForAlgorithm(algorithm types.KeyAlgorithm) x509.KeyUsage {
	switch algorithm {
	case types.KeyAlgorithmECDSA:
		return x509.KeyUsageDigitalSignature
	default: // RSA
		return x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	}
}

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

// PrivateKey generates an RSA private key with default parameters (for leaf certs).
func PrivateKey() (*rsa.PrivateKey, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, int(DefaultRSAKeySize))
	if err != nil {
		return nil, errors.Wrap(err, "error generating RSA private key")
	}
	return rsaKey, nil
}

// SelfSignedCertificate creates a self signed certificate
func SelfSignedCertificate(cfg *CertCfg, key crypto.Signer) (*x509.Certificate, error) {
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
	// verifies that the CN and/or OU for the cert is set
	if len(cfg.Subject.CommonName) == 0 || len(cfg.Subject.OrganizationalUnit) == 0 {
		return nil, errors.Errorf("certification's subject is not set, or invalid")
	}
	pub := key.Public()
	cert.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set subject key identifier")
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, key.Public(), key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create certificate")
	}
	return x509.ParseCertificate(certBytes)
}

// SignedCertificate creates a new X.509 certificate based on a template.
func SignedCertificate(
	cfg *CertCfg,
	csr *x509.CertificateRequest,
	key *rsa.PrivateKey,
	caCert *x509.Certificate,
	caKey crypto.PrivateKey,
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
	pub := key.Public()
	certTmpl.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set subject key identifier")
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create x509 certificate")
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
			return nil, errors.Wrap(err, "failed to Marshal ans1 public key")
		}
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	default:
		return nil, errors.New("only RSA and ECDSA public keys supported")
	}

	hash := sha1.Sum(publicKeyBytes)
	return hash[:], nil
}

// GenerateSignedCertificate generate a key and cert defined by CertCfg and signed by CA.
func GenerateSignedCertificate(caKey crypto.PrivateKey, caCert *x509.Certificate,
	cfg *CertCfg) (*rsa.PrivateKey, *x509.Certificate, error) {

	// create a private key
	key, err := PrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate private key")
	}

	// create a CSR
	csrTmpl := x509.CertificateRequest{Subject: cfg.Subject, DNSNames: cfg.DNSNames, IPAddresses: cfg.IPAddresses}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create certificate request")
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		logrus.Debugf("Failed to parse x509 certificate request: %s", err)
		return nil, nil, errors.Wrap(err, "error parsing x509 certificate request")
	}

	// create a cert
	cert, err := SignedCertificate(cfg, csr, key, caCert, caKey)
	if err != nil {
		logrus.Debugf("Failed to create a signed certificate: %s", err)
		return nil, nil, errors.Wrap(err, "failed to create a signed certificate")
	}
	return key, cert, nil
}

// GenerateSelfSignedCertificate generates a key/cert pair defined by CertCfg
// using the specified key parameters. KeyUsage is automatically adjusted
// based on the algorithm (ECDSA cannot have KeyEncipherment).
func GenerateSelfSignedCertificate(cfg *CertCfg, params PrivateKeyParams) (crypto.PrivateKey, *x509.Certificate, error) {
	key, err := GeneratePrivateKeyWithParams(params)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate private key")
	}

	// Set KeyUsage based on algorithm — ECDSA keys cannot perform key encipherment
	adjustedCfg := *cfg
	baseUsage := keyUsageForAlgorithm(params.Algorithm)
	if cfg.IsCA {
		adjustedCfg.KeyUsages = baseUsage | x509.KeyUsageCertSign
	} else {
		adjustedCfg.KeyUsages = baseUsage
	}

	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, nil, fmt.Errorf("generated key does not implement crypto.Signer")
	}

	crt, err := SelfSignedCertificate(&adjustedCfg, signer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create self-signed certificate")
	}
	return key, crt, nil
}

// hasShortCertRotationEnabled returns true if ShortCertRotation featuregate is enabled.
func hasShortCertRotationEnabled(installConfig *installconfig.InstallConfig) bool {
	fgs := installConfig.Config.EnabledFeatureGates()
	return fgs.Enabled(features.FeatureShortCertRotation)
}

// ValidityOneDay sets the validity of a cert to 24 hours - or 1 hour when ShortRotationEnabled featuregate is enabled.
func ValidityOneDay(installConfig *installconfig.InstallConfig) time.Duration {
	if hasShortCertRotationEnabled(installConfig) {
		return time.Hour * 2
	}
	return time.Hour * 24
}

// ValidityOneYear sets the validity of a cert to 1 year - or two hours when ShortRotationEnabled featuregate is enabled.
func ValidityOneYear(installConfig *installconfig.InstallConfig) time.Duration {
	if hasShortCertRotationEnabled(installConfig) {
		return time.Hour * 4
	}
	return time.Hour * 24 * 365
}

// ValidityTenYears sets the validity of a cert to 10 years.
func ValidityTenYears() time.Duration {
	return time.Hour * 24 * 365 * 10
}
