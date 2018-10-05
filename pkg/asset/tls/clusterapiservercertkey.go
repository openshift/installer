package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// ClusterAPIServerCertKey is the asset that generates the cluster API server key/cert pair.
type ClusterAPIServerCertKey struct {
	CertKey
}

var _ asset.Asset = (*ClusterAPIServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *ClusterAPIServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AggregatorCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *ClusterAPIServerCertKey) Generate(dependencies asset.Parents) error {
	aggregatorCA := &AggregatorCA{}
	dependencies.Get(aggregatorCA)

	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "clusterapi.openshift-cluster-api.svc", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return a.CertKey.Generate(cfg, aggregatorCA, "cluster-apiserver-ca", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *ClusterAPIServerCertKey) Name() string {
	return "Certificate (clusterapi.openshift-cluster-api.svc)"
}
