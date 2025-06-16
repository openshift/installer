package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervctypes "github.com/openshift/installer/pkg/types/powervc"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
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
func (a *MCSCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &RootCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	hostname := internalAPIAddress(installConfig.Config)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:machine-config-server"},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears(),
	}

	var vips []string
	switch installConfig.Config.Platform.Name() {
	case baremetaltypes.Name:
		vips = installConfig.Config.BareMetal.APIVIPs
	case nutanixtypes.Name:
		vips = installConfig.Config.Nutanix.APIVIPs
	case openstacktypes.Name, powervctypes.Name:
		vips = installConfig.Config.OpenStack.APIVIPs
	case ovirttypes.Name:
		vips = installConfig.Config.Ovirt.APIVIPs
	case vspheretypes.Name:
		vips = installConfig.Config.VSphere.APIVIPs
	}

	cfg.IPAddresses = []net.IP{}
	cfg.DNSNames = []string{hostname}
	for _, vip := range vips {
		cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(vip))
		cfg.DNSNames = append(cfg.DNSNames, vip)
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "machine-config-server", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *MCSCertKey) Name() string {
	return "Certificate (mcs)"
}

// RegenerateMCSCertKey generates the cert/key pair based on input values.
func RegenerateMCSCertKey(ic *installconfig.InstallConfig, ca *RootCA, privateLBs []string) ([]byte, []byte, error) {
	hostname := internalAPIAddress(ic.Config)
	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:machine-config-server"},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears(),
	}
	cfg.IPAddresses = []net.IP{}
	cfg.DNSNames = []string{hostname}
	for _, ip := range privateLBs {
		cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(ip))
		cfg.DNSNames = append(cfg.DNSNames, ip)
	}
	return RegenerateSignedCertKey(cfg, ca, DoNotAppendParent)
}
