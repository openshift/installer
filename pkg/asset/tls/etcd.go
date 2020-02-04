package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// EtcdSignerCertKey is a key/cert pair that signs the etcd client and peer certs.
type EtcdSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*EtcdSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *EtcdSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *EtcdSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "etcd-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "etcd-signer")
}

// Name returns the human-friendly name of the asset.
func (c *EtcdSignerCertKey) Name() string {
	return "Certificate (etcd-signer)"
}

// EtcdCABundle is the asset the generates the etcd-ca-bundle,
// which contains all the individual client CAs.
type EtcdCABundle struct {
	CertBundle
}

var _ asset.Asset = (*EtcdCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *EtcdCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *EtcdCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("etcd-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdCABundle) Name() string {
	return "Certificate (etcd-ca-bundle)"
}

// EtcdSignerClientCertKey is the asset that generates the etcd client key/cert pair.
type EtcdSignerClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*EtcdSignerClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdSignerClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdSignerClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &EtcdSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "etcd-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdSignerClientCertKey) Name() string {
	return "Certificate (etcd-client)"
}



// EtcdKubeAPIServerClientCertSigner is a key/cert pair that signs the a client cert for the etcd static pods to use
// to access the kube-apiserver to get the current list of endpoints in a happy path.
type EtcdKubeAPIServerClientCertSigner struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*EtcdKubeAPIServerClientCertSigner)(nil)

// Dependencies returns the dependency of the EtcdKubeAPIServerClientCertSigner, which is empty.
func (c *EtcdKubeAPIServerClientCertSigner) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *EtcdKubeAPIServerClientCertSigner) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "etcd-kube-apiserver-client-cert-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "etcd-kube-apiserver-client-cert-signer")
}

// Name returns the human-friendly name of the asset.
func (c *EtcdKubeAPIServerClientCertSigner) Name() string {
	return "Certificate (ketcd-kube-apiserver-client-cert-signer)"
}

// EtcdKubeAPIServerClientCertCABundle is the CA bundle for EtcdKubeAPIServerClientCertSigner
type EtcdKubeAPIServerClientCertCABundle struct {
	CertBundle
}

var _ asset.Asset = (*EtcdKubeAPIServerClientCertCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *EtcdKubeAPIServerClientCertCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdKubeAPIServerClientCertSigner{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *EtcdKubeAPIServerClientCertCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("etcd-kube-apiserver-client-cert-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdKubeAPIServerClientCertCABundle) Name() string {
	return "Certificate (etcd-kube-apiserver-client-cert-ca-bundle)"
}

// EtcdKubeAPIServerClientCert is the asset that generates the key/cert pair for etcd to communicate to the kube-apiserver.
type EtcdKubeAPIServerClientCert struct {
	SignedCertKey
}

var _ asset.Asset = (*EtcdKubeAPIServerClientCert)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdKubeAPIServerClientCert) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdKubeAPIServerClientCertSigner{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdKubeAPIServerClientCert) Generate(dependencies asset.Parents) error {
	ca := &EtcdKubeAPIServerClientCertSigner{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:serviceaccount:openshift-etcd:default", Organization: []string{"system:serviceaccounts:openshift-etcd"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "etcd-kube-apiserver-client-cert", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdKubeAPIServerClientCert) Name() string {
	return "Certificate (etcd-kube-apiserver-client-cert)"
}
