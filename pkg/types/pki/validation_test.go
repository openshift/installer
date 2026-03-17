package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

func TestValidatePKIConfig(t *testing.T) {
	fldPath := field.NewPath("pki")

	cases := []struct {
		name        string
		pkiConfig   *types.PKIConfig
		fips        bool
		expectError bool
		errorCount  int
	}{
		{
			name:        "nil config is valid",
			pkiConfig:   nil,
			expectError: false,
		},
		{
			name: "valid RSA signer config",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: configv1alpha1.CertificateConfig{
					Key: configv1alpha1.KeyConfig{
						Algorithm: configv1alpha1.KeyAlgorithmRSA,
						RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid ECDSA signer config",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: configv1alpha1.CertificateConfig{
					Key: configv1alpha1.KeyConfig{
						Algorithm: configv1alpha1.KeyAlgorithmECDSA,
						ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP384},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "empty PKI config - signerCertificates required",
			pkiConfig:   &types.PKIConfig{},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "invalid RSA key size",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: configv1alpha1.CertificateConfig{
					Key: configv1alpha1.KeyConfig{
						Algorithm: configv1alpha1.KeyAlgorithmRSA,
						RSA:       configv1alpha1.RSAKeyConfig{KeySize: 1024},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "invalid ECDSA curve",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: configv1alpha1.CertificateConfig{
					Key: configv1alpha1.KeyConfig{
						Algorithm: configv1alpha1.KeyAlgorithmECDSA,
						ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: "P224"},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidatePKIConfig(tc.pkiConfig, fldPath, tc.fips)
			if tc.expectError {
				assert.Len(t, errs, tc.errorCount)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidateKeyConfig(t *testing.T) {
	fldPath := field.NewPath("pki", "signerCertificates", "key")

	cases := []struct {
		name        string
		config      configv1alpha1.KeyConfig
		fips        bool
		expectError bool
		errorCount  int
	}{
		// Valid RSA configs
		{
			name: "valid RSA 2048",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 2048},
			},
		},
		{
			name: "valid RSA 3072",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 3072},
			},
		},
		{
			name: "valid RSA 4096",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
			},
		},
		{
			name: "valid RSA 8192",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 8192},
			},
		},
		// Valid ECDSA configs
		{
			name: "valid ECDSA P256",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP256},
			},
		},
		{
			name: "valid ECDSA P384",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP384},
			},
		},
		{
			name: "valid ECDSA P521",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP521},
			},
		},
		// Invalid: missing algorithm
		{
			name: "missing algorithm",
			config: configv1alpha1.KeyConfig{
				RSA: configv1alpha1.RSAKeyConfig{KeySize: 4096},
			},
			expectError: true,
			errorCount:  1,
		},
		// Invalid: unsupported algorithm
		{
			name: "unsupported algorithm",
			config: configv1alpha1.KeyConfig{
				Algorithm: "Ed25519",
			},
			expectError: true,
			errorCount:  1,
		},
		// Invalid RSA key sizes
		{
			name: "RSA key size too small",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 1024},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA key size too large",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 9216},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA key size not multiple of 1024",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 5000},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA missing key size",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
			},
			expectError: true,
			errorCount:  1,
		},
		// Invalid ECDSA curves
		{
			name: "ECDSA invalid curve",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: "P224"},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "ECDSA missing curve",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
			},
			expectError: true,
			errorCount:  1,
		},
		// Mismatched algorithm/params
		{
			name: "RSA with ECDSA config",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP256},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "ECDSA with RSA config",
			config: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP384},
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateKeyConfig(tc.config, fldPath, tc.fips)
			if tc.expectError {
				assert.Len(t, errs, tc.errorCount, "expected %d errors, got: %v", tc.errorCount, errs)
			} else {
				assert.Empty(t, errs, "expected no errors, got: %v", errs)
			}
		})
	}
}
