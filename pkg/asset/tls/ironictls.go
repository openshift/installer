package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// IronicTLSCert is the asset that generates the key/cert pair that is used for
// enabling TLS for virtual media in ironic.
type IronicTLSCert struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*IronicTLSCert)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *IronicTLSCert) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *IronicTLSCert) Generate(ctx context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)
	if installConfig.Config.Platform.BareMetal == nil {
		return nil
	}

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "ironic", OrganizationalUnit: []string{"openshift"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityOneDay(installConfig),
	}

	vips := installConfig.Config.BareMetal.APIVIPs
	if installConfig.Config.BareMetal.BootstrapProvisioningIP != "" {
		vips = append(vips, installConfig.Config.BareMetal.BootstrapProvisioningIP)
	}
	cfg.IPAddresses = []net.IP{}
	for _, vip := range vips {
		cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(vip))
	}
	hostname := internalAPIAddress(installConfig.Config)
	cfg.DNSNames = []string{hostname}

	logrus.Debugf("Generating TLS certificate for ironic (virtual media)")
	return a.SelfSignedCertKey.Generate(ctx, cfg, "ironic/tls")
}

// Name returns the human-friendly name of the asset.
func (a *IronicTLSCert) Name() string {
	return "Certificate (ironic virtual media TLS)"
}
