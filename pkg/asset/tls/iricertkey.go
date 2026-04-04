package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	features "github.com/openshift/api/features"
	libpki "github.com/openshift/library-go/pkg/pki"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/manifests"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// IRICertKey is the asset that generates the InternalReleaseImage registry key/cert pair.
type IRICertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*IRICertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *IRICertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
		&installconfig.InstallConfig{},
		&manifests.InternalReleaseImage{},
		&SignerPKIConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *IRICertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &RootCA{}
	installConfig := &installconfig.InstallConfig{}
	iri := &manifests.InternalReleaseImage{}
	pkiCfg := &SignerPKIConfig{}
	dependencies.Get(ca, installConfig, iri, pkiCfg)

	if !installConfig.Config.Enabled(features.FeatureGateNoRegistryClusterInstall) {
		return nil
	}

	// Skip if InternalReleaseImage manifest wasn't found.
	if len(iri.FileList) == 0 {
		return nil
	}

	apiInt := internalAPIAddress(installConfig.Config)

	keyGen, err := resolveKeyGen(pkiCfg, libpki.CertificateTypeServing, "installer.ingress-router-initial")
	if err != nil {
		return err
	}

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:internal-release-image"},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears(),
		CertType:     libpki.CertificateTypeServing,
	}

	var vips []string
	switch installConfig.Config.Platform.Name() {
	case baremetaltypes.Name:
		vips = installConfig.Config.BareMetal.APIVIPs
	case nutanixtypes.Name:
		vips = installConfig.Config.Nutanix.APIVIPs
	case vspheretypes.Name:
		vips = installConfig.Config.VSphere.APIVIPs
	}

	cfg.IPAddresses = []net.IP{}
	cfg.DNSNames = []string{
		"localhost",
		apiInt,
	}
	localIPs := []string{
		"127.0.0.1",
		"::1",
	}
	for _, vip := range vips {
		cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(vip))
		cfg.DNSNames = append(cfg.DNSNames, vip)
	}
	for _, i := range localIPs {
		if ip := net.ParseIP(i); ip != nil {
			cfg.IPAddresses = append(cfg.IPAddresses, ip)
		}
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "internal-release-image", DoNotAppendParent, keyGen)
}

// Name returns the human-friendly name of the asset.
func (a *IRICertKey) Name() string {
	return "Certificate (InternalReleaseImage)"
}
