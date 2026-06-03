package tls

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"

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
			assert.NoError(t, err)

			pemBytes, err := PrivateKeyToPem(key)
			assert.NoError(t, err)
			assert.NotEmpty(t, pemBytes)

			decoded, err := PemToPrivateKey(pemBytes)
			assert.NoError(t, err)
			assert.IsType(t, tc.expectType, decoded)
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
		assert.NoError(t, err)
		pemBytes, pemErr := PrivateKeyToPem(key)
		assert.NoError(t, pemErr)

		decoded, err := PemToPrivateKey(pemBytes)
		assert.NoError(t, err)
		_, ok := decoded.(*rsa.PrivateKey)
		assert.True(t, ok, "expected *rsa.PrivateKey")
	})

	t.Run("EC PEM block", func(t *testing.T) {
		key, err := GenerateECDSAPrivateKey(types.ECDSACurveP256)
		assert.NoError(t, err)
		pemBytes, pemErr := PrivateKeyToPem(key)
		assert.NoError(t, pemErr)

		decoded, err := PemToPrivateKey(pemBytes)
		assert.NoError(t, err)
		_, ok := decoded.(*ecdsa.PrivateKey)
		assert.True(t, ok, "expected *ecdsa.PrivateKey")
	})

	t.Run("PKCS#8 PEM block", func(t *testing.T) {
		key, err := GenerateRSAPrivateKey(2048)
		assert.NoError(t, err)
		pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
		assert.NoError(t, err)
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: pkcs8Bytes,
		})

		decoded, err := PemToPrivateKey(pemBytes)
		assert.NoError(t, err)
		_, ok := decoded.(*rsa.PrivateKey)
		assert.True(t, ok, "expected *rsa.PrivateKey")
	})
}
