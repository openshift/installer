package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// KubeletCertKey is the asset that generates the kubelet key/cert pair.
type KubeletCertKey struct {
	CertKey
}

var _ asset.Asset = (*KubeletCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *KubeletCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeletCertKey) Generate(dependencies asset.Parents) error {
	kubeCA := &KubeCA{}
	dependencies.Get(kubeCA)

	cfg := &CertCfg{
		// system:masters is a hack to get the kubelet up without kube-core
		// TODO(node): make kubelet bootstrapping secure with minimal permissions eventually switching to system:node:* CommonName
		Subject:      pkix.Name{CommonName: "system:serviceaccount:kube-system:default", Organization: []string{"system:serviceaccounts:kube-system", "system:masters"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityThirtyMinutes,
	}

	return a.CertKey.Generate(cfg, kubeCA, "kubelet", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeletCertKey) Name() string {
	return "Certificate (system:serviceaccount:kube-system:default)"
}
