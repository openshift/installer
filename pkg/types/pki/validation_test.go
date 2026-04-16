package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

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
				SignerCertificates: types.CertificateConfig{
					Key: types.KeyConfig{
						Algorithm: types.KeyAlgorithmRSA,
						RSA:       &types.RSAKeyConfig{KeySize: 4096},
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid ECDSA signer config",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: types.CertificateConfig{
					Key: types.KeyConfig{
						Algorithm: types.KeyAlgorithmECDSA,
						ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
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
				SignerCertificates: types.CertificateConfig{
					Key: types.KeyConfig{
						Algorithm: types.KeyAlgorithmRSA,
						RSA:       &types.RSAKeyConfig{KeySize: 1024},
					},
				},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "invalid ECDSA curve",
			pkiConfig: &types.PKIConfig{
				SignerCertificates: types.CertificateConfig{
					Key: types.KeyConfig{
						Algorithm: types.KeyAlgorithmECDSA,
						ECDSA:     &types.ECDSAKeyConfig{Curve: "P224"},
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
		config      types.KeyConfig
		fips        bool
		expectError bool
		errorCount  int
	}{
		// Valid RSA configs
		{
			name: "valid RSA 2048",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 2048},
			},
		},
		{
			name: "valid RSA 4096",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 4096},
			},
		},
		{
			name: "valid RSA 8192",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 8192},
			},
		},
		// Valid ECDSA configs
		{
			name: "valid ECDSA P256",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
				ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP256},
			},
		},
		{
			name: "valid ECDSA P384",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
				ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
			},
		},
		{
			name: "valid ECDSA P521",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
				ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP521},
			},
		},
		// Invalid: missing algorithm
		{
			name:        "missing algorithm",
			config:      types.KeyConfig{},
			expectError: true,
			errorCount:  1,
		},
		// Invalid: unsupported algorithm
		{
			name: "unsupported algorithm",
			config: types.KeyConfig{
				Algorithm: "Ed25519",
			},
			expectError: true,
			errorCount:  1,
		},
		// Invalid RSA key sizes
		{
			name: "RSA key size too small",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 1024},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA key size too large",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 9216},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA key size not multiple of 1024",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 5000},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "RSA missing rsa field",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
			},
			expectError: true,
			errorCount:  1,
		},
		// Invalid ECDSA curves
		{
			name: "ECDSA invalid curve",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
				ECDSA:     &types.ECDSAKeyConfig{Curve: "P224"},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "ECDSA missing ecdsa field",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
			},
			expectError: true,
			errorCount:  1,
		},
		// Mismatched algorithm/params
		{
			name: "RSA with ECDSA config",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 4096},
				ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP256},
			},
			expectError: true,
			errorCount:  1,
		},
		{
			name: "ECDSA with RSA config",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
				ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
				RSA:       &types.RSAKeyConfig{KeySize: 4096},
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
