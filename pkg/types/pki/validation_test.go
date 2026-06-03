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
		expectError bool
		errorCount  int
		errorMsg    string
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
			name:        "empty PKI config - algorithm required",
			pkiConfig:   &types.PKIConfig{},
			expectError: true,
			errorCount:  1,
			errorMsg:    "algorithm must be specified",
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
			errorMsg:    "must be a multiple of 1024 from 2048 to 8192",
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
			errorMsg:    "supported values",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidatePKIConfig(tc.pkiConfig, fldPath)
			if tc.expectError {
				if assert.Len(t, errs, tc.errorCount) {
					assert.Regexp(t, tc.errorMsg, errs.ToAggregate())
				}
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
		expectError bool
		errorCount  int
		errorMsg    string
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
			name: "valid RSA 3072",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 3072},
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
			name: "missing algorithm",
			config: types.KeyConfig{
				RSA: &types.RSAKeyConfig{KeySize: 4096},
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "algorithm must be specified",
		},
		// Invalid: unsupported algorithm
		{
			name: "unsupported algorithm",
			config: types.KeyConfig{
				Algorithm: "Ed25519",
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "supported values",
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
			errorMsg:    "must be a multiple of 1024 from 2048 to 8192",
		},
		{
			name: "RSA key size too large",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 9216},
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "must be a multiple of 1024 from 2048 to 8192",
		},
		{
			name: "RSA key size not multiple of 1024",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
				RSA:       &types.RSAKeyConfig{KeySize: 5000},
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "must be a multiple of 1024 from 2048 to 8192",
		},
		{
			name: "RSA missing rsa config",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmRSA,
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "rsa is required when algorithm is RSA",
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
			errorMsg:    "supported values",
		},
		{
			name: "ECDSA missing ecdsa config",
			config: types.KeyConfig{
				Algorithm: types.KeyAlgorithmECDSA,
			},
			expectError: true,
			errorCount:  1,
			errorMsg:    "ecdsa is required when algorithm is ECDSA",
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
			errorMsg:    "ecdsa must not be set when algorithm is RSA",
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
			errorMsg:    "rsa must not be set when algorithm is ECDSA",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateKeyConfig(tc.config, fldPath)
			if tc.expectError {
				if assert.Len(t, errs, tc.errorCount, "expected %d errors, got: %v", tc.errorCount, errs) {
					assert.Regexp(t, tc.errorMsg, errs.ToAggregate())
				}
			} else {
				assert.Empty(t, errs, "expected no errors, got: %v", errs)
			}
		})
	}
}
