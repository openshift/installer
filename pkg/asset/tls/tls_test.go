package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func TestSelfSignedCertificate(t *testing.T) {
	key, err := PrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate Private Key: %v", err)
	}
	cases := []struct {
		cfg *CertCfg
		err bool
	}{
		{
			cfg: &CertCfg{
				Validity:  time.Hour * 5,
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					CommonName:         "root_ca",
					OrganizationalUnit: []string{"openshift"},
				},
				IsCA: true,
			},
			err: false,
		},
		{
			cfg: &CertCfg{
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					CommonName: "root_ca",
				},
				IsCA: false,
			},
			err: true,
		},
		{
			cfg: &CertCfg{
				KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
				Subject: pkix.Name{
					OrganizationalUnit: []string{"openshift"},
				},
			},
			err: true,
		},
	}
	for i, c := range cases {
		if _, err := SelfSignedCertificate(c.cfg, key); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestSignedCertificate(t *testing.T) {
	key, err := PrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	cases := []struct {
		Subject            pkix.Name
		SignatureAlgorithm x509.SignatureAlgorithm
		err                bool
	}{
		{
			Subject: pkix.Name{
				CommonName:         "csr",
				OrganizationalUnit: []string{"openshift"},
			},
			err: false,
		},
		{
			Subject: pkix.Name{},
			err:     false,
		},
		{
			Subject: pkix.Name{
				CommonName:         "csr-wrong-alg",
				OrganizationalUnit: []string{"openshift"},
			},
			SignatureAlgorithm: 123,
			err:                true,
		},
	}
	for i, c := range cases {
		csrTmpl := x509.CertificateRequest{
			Subject:            c.Subject,
			SignatureAlgorithm: c.SignatureAlgorithm,
		}
		if _, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestGenerateRSAPrivateKey(t *testing.T) {
	cases := []struct {
		name    string
		keySize int32
	}{
		{"RSA 2048", 2048},
		{"RSA 4096", 4096},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := GenerateRSAPrivateKey(tc.keySize)
			assert.NoError(t, err)
			assert.IsType(t, &rsa.PrivateKey{}, key)
			assert.Equal(t, int(tc.keySize), key.N.BitLen())
		})
	}
}

func TestGenerateECDSAPrivateKey(t *testing.T) {
	cases := []struct {
		name      string
		curve     types.ECDSACurve
		expected  elliptic.Curve
		expectErr bool
	}{
		{"P256", types.ECDSACurveP256, elliptic.P256(), false},
		{"P384", types.ECDSACurveP384, elliptic.P384(), false},
		{"P521", types.ECDSACurveP521, elliptic.P521(), false},
		{"invalid", "P224", nil, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := GenerateECDSAPrivateKey(tc.curve)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.IsType(t, &ecdsa.PrivateKey{}, key)
			assert.Equal(t, tc.expected, key.Curve)
		})
	}
}

func TestGenerateSelfSignedCertificateWithParams(t *testing.T) {
	cases := []struct {
		name            string
		params          PrivateKeyParams
		expectKeyType   interface{}
		expectPubKeyAlg x509.PublicKeyAlgorithm
	}{
		{
			name: "RSA 4096 CA",
			params: PrivateKeyParams{
				Algorithm:  types.KeyAlgorithmRSA,
				RSAKeySize: 4096,
			},
			expectKeyType:   &rsa.PrivateKey{},
			expectPubKeyAlg: x509.RSA,
		},
		{
			name: "ECDSA P384 CA",
			params: PrivateKeyParams{
				Algorithm:  types.KeyAlgorithmECDSA,
				ECDSACurve: types.ECDSACurveP384,
			},
			expectKeyType:   &ecdsa.PrivateKey{},
			expectPubKeyAlg: x509.ECDSA,
		},
		{
			name: "RSA 2048 CA (default)",
			params: PrivateKeyParams{
				Algorithm:  types.KeyAlgorithmRSA,
				RSAKeySize: 2048,
			},
			expectKeyType:   &rsa.PrivateKey{},
			expectPubKeyAlg: x509.RSA,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &CertCfg{
				Subject:  pkix.Name{CommonName: "test-ca", OrganizationalUnit: []string{"openshift"}},
				Validity: time.Hour,
				IsCA:     true,
			}
			key, cert, err := GenerateSelfSignedCertificate(cfg, tc.params)
			assert.NoError(t, err)
			assert.IsType(t, tc.expectKeyType, key)
			assert.Equal(t, tc.expectPubKeyAlg, cert.PublicKeyAlgorithm)
			assert.True(t, cert.IsCA)
		})
	}
}

func TestKeyUsageForAlgorithm(t *testing.T) {
	cases := []struct {
		name      string
		params    PrivateKeyParams
		isCA      bool
		wantUsage x509.KeyUsage
		notUsage  x509.KeyUsage
	}{
		{
			name:      "RSA signer",
			params:    PrivateKeyParams{Algorithm: types.KeyAlgorithmRSA, RSAKeySize: 2048},
			isCA:      true,
			wantUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		},
		{
			name:      "ECDSA signer",
			params:    PrivateKeyParams{Algorithm: types.KeyAlgorithmECDSA, ECDSACurve: types.ECDSACurveP256},
			isCA:      true,
			wantUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			notUsage:  x509.KeyUsageKeyEncipherment,
		},
		{
			name:      "RSA non-CA",
			params:    PrivateKeyParams{Algorithm: types.KeyAlgorithmRSA, RSAKeySize: 2048},
			isCA:      false,
			wantUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			notUsage:  x509.KeyUsageCertSign,
		},
		{
			name:      "ECDSA non-CA",
			params:    PrivateKeyParams{Algorithm: types.KeyAlgorithmECDSA, ECDSACurve: types.ECDSACurveP384},
			isCA:      false,
			wantUsage: x509.KeyUsageDigitalSignature,
			notUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &CertCfg{
				Subject:  pkix.Name{CommonName: "test", OrganizationalUnit: []string{"openshift"}},
				Validity: time.Hour,
				IsCA:     tc.isCA,
			}
			_, cert, err := GenerateSelfSignedCertificate(cfg, tc.params)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantUsage, cert.KeyUsage, "KeyUsage mismatch")
			if tc.notUsage != 0 {
				assert.Zero(t, cert.KeyUsage&tc.notUsage, "unexpected KeyUsage bits set")
			}
		})
	}
}

func TestSignatureAlgorithmAutoDetection(t *testing.T) {
	cases := []struct {
		name     string
		params   PrivateKeyParams
		expected x509.SignatureAlgorithm
	}{
		{
			name:     "RSA",
			params:   PrivateKeyParams{Algorithm: types.KeyAlgorithmRSA, RSAKeySize: 2048},
			expected: x509.SHA256WithRSA,
		},
		{
			name:     "ECDSA P256",
			params:   PrivateKeyParams{Algorithm: types.KeyAlgorithmECDSA, ECDSACurve: types.ECDSACurveP256},
			expected: x509.ECDSAWithSHA256,
		},
		{
			name:     "ECDSA P384",
			params:   PrivateKeyParams{Algorithm: types.KeyAlgorithmECDSA, ECDSACurve: types.ECDSACurveP384},
			expected: x509.ECDSAWithSHA384,
		},
		{
			name:     "ECDSA P521",
			params:   PrivateKeyParams{Algorithm: types.KeyAlgorithmECDSA, ECDSACurve: types.ECDSACurveP521},
			expected: x509.ECDSAWithSHA512,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &CertCfg{
				Subject:  pkix.Name{CommonName: "test-sig", OrganizationalUnit: []string{"openshift"}},
				Validity: time.Hour,
				IsCA:     true,
			}
			_, cert, err := GenerateSelfSignedCertificate(cfg, tc.params)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, cert.SignatureAlgorithm)
		})
	}
}
