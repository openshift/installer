package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

// KubeCA is the asset that generates the kube-ca key/cert pair.
type KubeCA struct {
	CertKey
}

var _ asset.Asset = (*KubeCA)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *KubeCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeCA) Generate(dependencies asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-ca", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	key, crt, err := GenerateRootCertKey(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to generate Kube CA")
	}

	a.KeyRaw = PrivateKeyToPem(key)
	a.CertRaw = CertToPem(crt)

	a.generateFiles("kube-ca")

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *KubeCA) Name() string {
	return "Certificate (kube-ca)"
}
