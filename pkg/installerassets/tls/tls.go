// Package tls installs installerassets.Rebuilders for TLS assets.
package tls

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/assets"
	"github.com/pkg/errors"
)

const (
	keySize               = 2048
	validityTenYears      = time.Hour * 24 * 365 * 10
	validityThirtyMinutes = time.Minute * 30
)

type templateAdjuster func(ctx context.Context, asset *assets.Asset, getByName assets.GetByString, template *x509.Certificate) (err error)

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

func privateKey(ctx context.Context) (data []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	keyInBytes := x509.MarshalPKCS1PrivateKey(key)
	keyinPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyInBytes,
		},
	)
	return keyinPem, nil
}

// PEMToPrivateKey converts PEM data to a rsa.PrivateKey.
func PEMToPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// PublicKeyToPEM converts an rsa.PublicKey object to PEM.
func PublicKeyToPEM(key *rsa.PublicKey) ([]byte, error) {
	keyInBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
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

func pemToCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.Errorf("could not find a PEM block in the certificate")
	}
	return x509.ParseCertificate(block.Bytes)
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

func certificateRebuilder(name string, key string, caCert string, caKey string, template *x509.Certificate, templateAdjust templateAdjuster) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          name,
			RebuildHelper: certificateRebuilder(name, key, caKey, caCert, template, templateAdjust),
		}

		parents, err := asset.GetParents(ctx, getByName, key, caCert, caKey)
		if err != nil {
			return nil, err
		}

		keys := make(map[string]*rsa.PrivateKey)
		for _, keyName := range []string{key, caKey} {
			keyPEM, err := PEMToPrivateKey(parents[keyName].Data)
			if err != nil {
				return nil, err
			}

			keys[keyName] = keyPEM
		}

		template.SubjectKeyId, err = generateSubjectKeyID(keys[key].Public())
		if err != nil {
			return nil, errors.Wrap(err, "failed to set subject key identifier")
		}

		cert, err := pemToCertificate(parents[caCert].Data)
		if err != nil {
			return nil, err
		}

		if templateAdjust != nil {
			err = templateAdjust(ctx, asset, getByName, template)
			if err != nil {
				return nil, err
			}
		}

		der, err := x509.CreateCertificate(rand.Reader, template, cert, keys[key].Public(), keys[caKey])
		if err != nil {
			return nil, errors.Wrap(err, "failed to create certificate")
		}

		asset.Data = pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: der,
		})
		return asset, nil
	}
}

func certificateChainRebuilder(name string, parents ...string) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          name,
			RebuildHelper: certificateChainRebuilder(name, parents...),
		}

		parentAssets, err := asset.GetParents(ctx, getByName, parents...)
		if err != nil {
			return nil, err
		}

		data := make([][]byte, 0, len(parents))
		for _, parentName := range parents {
			data = append(data, parentAssets[parentName].Data)
		}

		asset.Data = bytes.Join(data, []byte("\n"))

		return asset, nil
	}
}
