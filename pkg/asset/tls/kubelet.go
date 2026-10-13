package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// KubeletCSRSignerCertKey is a key/cert pair that signs the kubelet client certs.
type KubeletCSRSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeletCSRSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeletCSRSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{&installconfig.InstallConfig{}}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeletCSRSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	parents.Get(installConfig)
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kubelet-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneDay(installConfig),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kubelet-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeletCSRSignerCertKey) Name() string {
	return "Certificate (kubelet-signer)"
}

// KubeletClientCABundle is the asset the generates the kubelet-client-ca-bundle,
// which contains all the individual client CAs.
type KubeletClientCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeletClientCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeletClientCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeletClientCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kubelet-client-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletClientCABundle) Name() string {
	return "Certificate (kubelet-client-ca-bundle)"
}

// KubeletServingCABundle is the asset the generates the kubelet-serving-ca-bundle,
// which contains all the individual client CAs.
type KubeletServingCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeletServingCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeletServingCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeletServingCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kubelet-serving-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletServingCABundle) Name() string {
	return "Certificate (kubelet-serving-ca-bundle)"
}

// KubeletBootstrapCertSigner is a key/cert pair that signs the kubelet bootstrap kubeconfig client certs that the kubelet
// uses to create CSRs for it's real certificates
type KubeletBootstrapCertSigner struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeletBootstrapCertSigner)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeletBootstrapCertSigner) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeletBootstrapCertSigner) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kubelet-bootstrap-kubeconfig-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kubelet-bootstrap-kubeconfig-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeletBootstrapCertSigner) Name() string {
	return "Certificate (kubelet-bootstrap-kubeconfig-signer)"
}

// KubeletBootstrapCABundle is the asset the generates the admin-kubeconfig-ca-bundle,
// which contains all the individual client CAs.
type KubeletBootstrapCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeletBootstrapCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeletBootstrapCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletBootstrapCertSigner{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeletBootstrapCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kubelet-bootstrap-kubeconfig-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletBootstrapCABundle) Name() string {
	return "Certificate (kubelet-bootstrap-kubeconfig-ca-bundle)"
}

// KubeletClientCertKey is the asset that generates the key/cert pair for kubelet client to apiserver.
// This credential can be revoked by deleting the configmap containing its signer.
type KubeletClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeletClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *KubeletClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeletBootstrapCertSigner{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeletClientCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeletBootstrapCertSigner{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper", Organization: []string{"system:serviceaccounts:openshift-machine-config-operator", "system:serviceaccounts"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears(),
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kubelet-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletClientCertKey) Name() string {
	return "Certificate (kubelet-client)"
}
