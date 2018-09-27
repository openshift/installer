package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// MCSCertKey is the asset that generates the MCS key/cert pair.
type MCSCertKey struct {
	CertKey
}

var _ asset.Asset = (*MCSCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *MCSCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *MCSCertKey) Generate(dependencies asset.Parents) error {
	rootCA := &RootCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(rootCA, installConfig)

	hostname := apiAddress(installConfig.Config)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: hostname},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears,
		DNSNames:     []string{hostname},
	}

	return a.CertKey.Generate(cfg, rootCA, "machine-config-server", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *MCSCertKey) Name() string {
	return "Certificate (mcs)"
}
