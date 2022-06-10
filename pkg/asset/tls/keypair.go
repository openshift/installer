package tls

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

// KeyPairInterface contains a private key and a public key.
type KeyPairInterface interface {
	// Private returns the private key.
	Private() []byte
	// Public returns the public key.
	Public() []byte
}

// KeyPair contains a private key and a public key.
type KeyPair struct {
	asset.DefaultFileListWriter

	Pvt []byte
	Pub []byte
}

// Generate generates the rsa private / public key pair.
func (k *KeyPair) Generate(filenameBase string) error {
	key, err := PrivateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate private key")
	}

	pubkeyData, err := PublicKeyToPem(&key.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to get public key data from private key")
	}

	k.Pvt = PrivateKeyToPem(key)
	k.Pub = pubkeyData

	k.FileList = []*asset.File{
		{
			Filename: assetFilePath(filenameBase + ".key"),
			Data:     k.Pvt,
		},
		{
			Filename: assetFilePath(filenameBase + ".pub"),
			Data:     k.Pub,
		},
	}

	return nil
}

// Public returns the public key.
func (k *KeyPair) Public() []byte {
	return k.Pub
}

// Private returns the private key.
func (k *KeyPair) Private() []byte {
	return k.Pvt
}