package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// MCSCertKey is the asset that generates the MCS key/cert pair.
type MCSCertKey struct {
	SignedCertKey
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
	ca := &RootCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	hostname := internalAPIAddress(installConfig.Config)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: hostname},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears,
	}

	switch installConfig.Config.Platform.Name() {
	case baremetaltypes.Name:
		cfg.IPAddresses = []net.IP{net.ParseIP(installConfig.Config.BareMetal.APIVIP)}
		cfg.DNSNames = []string{hostname, installConfig.Config.BareMetal.APIVIP}
	case openstacktypes.Name:
		apiVIP, err := openstackdefaults.APIVIP(installConfig.Config.Platform.OpenStack)
		if err != nil {
			return err
		}
		cfg.IPAddresses = []net.IP{apiVIP}
		cfg.DNSNames = []string{hostname, apiVIP.String()}
	default:
		cfg.DNSNames = []string{hostname}
	}

	return a.SignedCertKey.Generate(cfg, ca, "machine-config-server", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *MCSCertKey) Name() string {
	return "Certificate (mcs)"
}
