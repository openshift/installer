package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// KubeControlPlaneSignerCertKey is a key/cert pair that signs the kube control-plane client certs.
type KubeControlPlaneSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeControlPlaneSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeControlPlaneSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeControlPlaneSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-control-plane-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneYear,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kube-control-plane-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeControlPlaneSignerCertKey) Name() string {
	return "Certificate (kube-control-plane-signer)"
}

// KubeControlPlaneCABundle is the asset the generates the kube-control-plane-ca-bundle,
// which contains all the individual client CAs.
type KubeControlPlaneCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeControlPlaneCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeControlPlaneCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeControlPlaneSignerCertKey{},
		&KubeAPIServerLBSignerCertKey{},
		&KubeAPIServerLocalhostSignerCertKey{},
		&KubeAPIServerServiceNetworkSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeControlPlaneCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-control-plane-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeControlPlaneCABundle) Name() string {
	return "Certificate (kube-control-plane-ca-bundle)"
}

// KubeControlPlaneKubeControllerManagerClientCertKey is the asset that generates the kube-controller-manger client key/cert pair.
type KubeControlPlaneKubeControllerManagerClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeControlPlaneKubeControllerManagerClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeControlPlaneKubeControllerManagerClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeControlPlaneSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeControlPlaneKubeControllerManagerClientCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeControlPlaneSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:admin", Organization: []string{"system:masters"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneYear,
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-control-plane-kube-controller-manager-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeControlPlaneKubeControllerManagerClientCertKey) Name() string {
	return "Certificate (kube-control-plane-kube-controller-manager-client)"
}

// KubeControlPlaneKubeSchedulerClientCertKey is the asset that generates the kube-scheduler client key/cert pair.
type KubeControlPlaneKubeSchedulerClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeControlPlaneKubeSchedulerClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeControlPlaneKubeSchedulerClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeControlPlaneSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeControlPlaneKubeSchedulerClientCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeControlPlaneSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:admin", Organization: []string{"system:masters"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneYear,
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-control-plane-kube-scheduler-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeControlPlaneKubeSchedulerClientCertKey) Name() string {
	return "Certificate (kube-control-plane-kube-scheduler-client)"
}
