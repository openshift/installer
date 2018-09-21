package tls

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset"
)

// KeyPair implements the Asset interface and
// generates an RSA public/private key pair.
type KeyPair struct {
	PrivKeyFileName string
	PubKeyFileName  string
}

var _ asset.Asset = (*KeyPair)(nil)

// Dependencies returns the dependency of an rsa private / public key pair.
func (k *KeyPair) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the rsa private / public key pair.
func (k *KeyPair) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	key, err := PrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	pubkeyData, err := PublicKeyToPem(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key data: %v", err)
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: assetFilePath(k.PrivKeyFileName),
				Data: []byte(PrivateKeyToPem(key)),
			},
			{
				Name: assetFilePath(k.PubKeyFileName),
				Data: []byte(pubkeyData),
			},
		},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (k *KeyPair) Name() string {
	return fmt.Sprintf("Key Pair (%s)", k.PubKeyFileName)
}
