package tls

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
)

func rootCARebuilder(ctx context.Context, getByName assets.GetByString) (asset *assets.Asset, err error) {
	asset = &assets.Asset{
		Name:          "tls/root-ca.crt",
		RebuildHelper: rootCARebuilder,
	}

	cert := &x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:         true,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		NotAfter:     time.Now().Add(validityTenYears),
		NotBefore:    time.Now(),
		SerialNumber: new(big.Int).SetInt64(0),
		Subject:      pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
	}

	parents, err := asset.GetParents(ctx, getByName, "tls/root-ca.key")
	if err != nil {
		return nil, err
	}

	key, err := PEMToPrivateKey(parents["tls/root-ca.key"].Data)
	if err != nil {
		return nil, err
	}
	pub := key.Public()

	cert.SubjectKeyId, err = generateSubjectKeyID(pub)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set subject key identifier")
	}

	der, err := x509.CreateCertificate(rand.Reader, cert, cert, key.Public(), key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create certificate")
	}

	asset.Data = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})
	return asset, nil
}

func init() {
	installerassets.Defaults["tls/root-ca.key"] = privateKey
	installerassets.Rebuilders["tls/root-ca.crt"] = rootCARebuilder
}
