package tls //nolint:revive // pre-existing package name

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openshift/installer/pkg/types"
)

func TestPrivateKeyToPemRoundtrip(t *testing.T) {
	cases := []struct {
		name       string
		genFunc    func() (interface{}, error)
		expectType interface{}
	}{
		{
			name: "RSA key",
			genFunc: func() (interface{}, error) {
				return GenerateRSAPrivateKey(2048)
			},
			expectType: &rsa.PrivateKey{},
		},
		{
			name: "ECDSA P256 key",
			genFunc: func() (interface{}, error) {
				return GenerateECDSAPrivateKey(types.ECDSACurveP256)
			},
			expectType: &ecdsa.PrivateKey{},
		},
		{
			name: "ECDSA P384 key",
			genFunc: func() (interface{}, error) {
				return GenerateECDSAPrivateKey(types.ECDSACurveP384)
			},
			expectType: &ecdsa.PrivateKey{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key, err := tc.genFunc()
			require.NoError(t, err)

			pemBytes, err := PrivateKeyToPem(key)
			require.NoError(t, err)
			require.NotEmpty(t, pemBytes)

			decoded, err := PemToPrivateKey(pemBytes)
			require.NoError(t, err)
			assert.IsType(t, tc.expectType, decoded)

			reEncoded, reErr := PrivateKeyToPem(decoded)
			require.NoError(t, reErr)
			assert.Equal(t, pemBytes, reEncoded, "round-tripped key material must match original")
		})
	}
}

func TestPemToPrivateKeyFormats(t *testing.T) {
	t.Run("invalid PEM", func(t *testing.T) {
		_, err := PemToPrivateKey([]byte("not a PEM"))
		assert.Error(t, err)
	})

	t.Run("empty data", func(t *testing.T) {
		_, err := PemToPrivateKey([]byte{})
		assert.Error(t, err)
	})

	t.Run("RSA PEM block", func(t *testing.T) {
		key, err := GenerateRSAPrivateKey(2048)
		require.NoError(t, err)
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})

		decoded, err := PemToPrivateKey(pemBytes)
		require.NoError(t, err)
		_, ok := decoded.(*rsa.PrivateKey)
		assert.True(t, ok, "expected *rsa.PrivateKey")
	})

	t.Run("EC PEM block", func(t *testing.T) {
		key, err := GenerateECDSAPrivateKey(types.ECDSACurveP256)
		require.NoError(t, err)
		derBytes, err := x509.MarshalECPrivateKey(key)
		require.NoError(t, err)
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: derBytes,
		})

		decoded, err := PemToPrivateKey(pemBytes)
		require.NoError(t, err)
		_, ok := decoded.(*ecdsa.PrivateKey)
		assert.True(t, ok, "expected *ecdsa.PrivateKey")
	})

	t.Run("PKCS#8 PEM block", func(t *testing.T) {
		key, err := GenerateRSAPrivateKey(2048)
		require.NoError(t, err)
		pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
		require.NoError(t, err)
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		})

		decoded, err := PemToPrivateKey(pemBytes)
		require.NoError(t, err)
		_, ok := decoded.(*rsa.PrivateKey)
		assert.True(t, ok, "expected *rsa.PrivateKey")
	})

	t.Run("EC key with EC PARAMETERS prefix", func(t *testing.T) {
		ecKey, err := GenerateECDSAPrivateKey(types.ECDSACurveP256)
		assert.NoError(t, err)
		ecDER, err := x509.MarshalECPrivateKey(ecKey)
		assert.NoError(t, err)

		paramBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PARAMETERS",
			Bytes: []byte{0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07},
		})
		keyBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: ecDER,
		})
		combined := make([]byte, 0, len(paramBlock)+len(keyBlock))
		combined = append(combined, paramBlock...)
		combined = append(combined, keyBlock...)

		decoded, err := PemToPrivateKey(combined)
		assert.NoError(t, err)
		_, ok := decoded.(*ecdsa.PrivateKey)
		assert.True(t, ok, "expected *ecdsa.PrivateKey")
	})

	t.Run("unsupported block type only", func(t *testing.T) {
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte("not a key"),
		})
		_, err := PemToPrivateKey(pemBytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find a PEM block with a supported private key type")
	})

	t.Run("malformed key block followed by valid key", func(t *testing.T) {
		badBlock := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: []byte("malformed key data"),
		})
		key, err := GenerateECDSAPrivateKey(types.ECDSACurveP256)
		assert.NoError(t, err)
		goodPem, err := PrivateKeyToPem(key)
		assert.NoError(t, err)
		combined := make([]byte, 0, len(badBlock)+len(goodPem))
		combined = append(combined, badBlock...)
		combined = append(combined, goodPem...)

		decoded, err := PemToPrivateKey(combined)
		assert.NoError(t, err)
		_, ok := decoded.(*ecdsa.PrivateKey)
		assert.True(t, ok, "expected *ecdsa.PrivateKey")
	})
}
