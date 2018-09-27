package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// IngressCertKey is the asset that generates the ingress key/cert pair.
type IngressCertKey struct {
	CertKey
}

var _ asset.Asset = (*IngressCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *IngressCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeCA{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *IngressCertKey) Generate(dependencies asset.Parents) error {
	kubeCA := &KubeCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(kubeCA, installConfig)

	baseAddress := fmt.Sprintf("%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: baseAddress, Organization: []string{"ingress"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
		DNSNames: []string{
			baseAddress,
			fmt.Sprintf("*.%s", baseAddress),
		},
	}

	return a.CertKey.Generate(cfg, kubeCA, "ingress", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *IngressCertKey) Name() string {
	return "Certificate (ingress)"
}
