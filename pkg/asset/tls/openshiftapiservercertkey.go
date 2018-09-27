package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// OpenshiftAPIServerCertKey is the asset that generates the Openshift API server key/cert pair.
type OpenshiftAPIServerCertKey struct {
	CertKey
}

var _ asset.Asset = (*OpenshiftAPIServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *OpenshiftAPIServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AggregatorCA{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *OpenshiftAPIServerCertKey) Generate(dependencies asset.Parents) error {
	aggregatorCA := &AggregatorCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(aggregatorCA, installConfig)

	apiServerAddress, err := cidrhost(installConfig.Config.Networking.ServiceCIDR.IPNet, 1)
	if err != nil {
		return errors.Wrap(err, "failed to get API Server address from InstallConfig")
	}

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "openshift-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
		DNSNames: []string{
			apiAddress(installConfig.Config),
			"openshift-apiserver",
			"openshift-apiserver.kube-system",
			"openshift-apiserver.kube-system.svc",
			"openshift-apiserver.kube-system.svc.cluster.local",
			"localhost", "127.0.0.1",
		},
		IPAddresses: []net.IP{net.ParseIP(apiServerAddress)},
	}

	return a.CertKey.Generate(cfg, aggregatorCA, "openshift-apiserver", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *OpenshiftAPIServerCertKey) Name() string {
	return "Certificate (openshift-apiserver)"
}
