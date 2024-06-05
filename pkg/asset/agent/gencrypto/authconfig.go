package gencrypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/golang-jwt/jwt/v4"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
)

// AuthConfig is an asset that generates ECDSA public/private keys, JWT token.
type AuthConfig struct {
	PublicKey, Token string
}

var _ asset.Asset = (*AuthConfig)(nil)

// LocalJWTKeyType suggests the key type to be used for the token.
type LocalJWTKeyType string

const (
	// InfraEnvKey is used to generate token using infra env id.
	InfraEnvKey LocalJWTKeyType = "infra_env_id"
)

var _ asset.Asset = (*AuthConfig)(nil)

// Dependencies returns the assets on which the AuthConfig asset depends.
func (a *AuthConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&common.InfraEnvID{},
	}
}

// Generate generates the auth config for agent installer APIs.
func (a *AuthConfig) Generate(dependencies asset.Parents) error {
	infraEnvID := &common.InfraEnvID{}
	dependencies.Get(infraEnvID)

	PublicKey, PrivateKey, err := keyPairPEM()
	if err != nil {
		return err
	}
	// Encode to Base64 (Standard encoding)
	encodedPubKeyPEM := base64.StdEncoding.EncodeToString([]byte(PublicKey))

	a.PublicKey = encodedPubKeyPEM

	token, err := localJWTForKey(infraEnvID.ID, PrivateKey)
	if err != nil {
		return err
	}
	a.Token = token

	return nil
}

// Name returns the human-friendly name of the asset.
func (*AuthConfig) Name() string {
	return "Agent Installer API Auth Config"
}

func keyPairPEM() (string, string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// encode private key to PEM string
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}

	var privKeyPEM bytes.Buffer
	err = pem.Encode(&privKeyPEM, block)
	if err != nil {
		return "", "", err
	}

	// encode public key to PEM string
	pubBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return "", "", err
	}

	block = &pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: pubBytes,
	}

	var pubKeyPEM bytes.Buffer
	err = pem.Encode(&pubKeyPEM, block)
	if err != nil {
		return "", "", err
	}

	return pubKeyPEM.String(), privKeyPEM.String(), nil
}

func localJWTForKey(id string, privateKkeyPem string) (string, error) {
	priv, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKkeyPem))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		string(InfraEnvKey): id,
	})

	tokenString, err := token.SignedString(priv)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
