package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// KubeAPIServerToKubeletSignerCertKey is a key/cert pair that signs the kube-apiserver to kubelet client certs.
type KubeAPIServerToKubeletSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeAPIServerToKubeletSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeAPIServerToKubeletSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{&installconfig.InstallConfig{}}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeAPIServerToKubeletSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	parents.Get(installConfig)
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-apiserver-to-kubelet-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneYear(installConfig),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kube-apiserver-to-kubelet-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeAPIServerToKubeletSignerCertKey) Name() string {
	return "Certificate (kube-apiserver-to-kubelet-signer)"
}

// KubeAPIServerToKubeletCABundle is the asset the generates the kube-apiserver-to-kubelet-ca-bundle,
// which contains all the individual client CAs.
type KubeAPIServerToKubeletCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerToKubeletCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerToKubeletCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerToKubeletSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerToKubeletCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-to-kubelet-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerToKubeletCABundle) Name() string {
	return "Certificate (kube-apiserver-to-kubelet-ca-bundle)"
}

// KubeAPIServerToKubeletClientCertKey is the asset that generates the kube-apiserver to kubelet client key/cert pair.
type KubeAPIServerToKubeletClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeAPIServerToKubeletClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeAPIServerToKubeletClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerToKubeletSignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeAPIServerToKubeletClientCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeAPIServerToKubeletSignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneYear(installConfig),
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-apiserver-to-kubelet-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerToKubeletClientCertKey) Name() string {
	return "Certificate (kube-apiserver-to-kubelet-client)"
}

// KubeAPIServerLocalhostSignerCertKey is a key/cert pair that signs the kube-apiserver server cert for SNI localhost.
type KubeAPIServerLocalhostSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeAPIServerLocalhostSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeAPIServerLocalhostSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeAPIServerLocalhostSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-apiserver-localhost-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kube-apiserver-localhost-signer")
}

// Load reads the asset files from disk.
func (c *KubeAPIServerLocalhostSignerCertKey) Load(f asset.FileFetcher) (bool, error) {
	return c.loadCertKey(f, "kube-apiserver-localhost-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeAPIServerLocalhostSignerCertKey) Name() string {
	return "Certificate (kube-apiserver-localhost-signer)"
}

// KubeAPIServerLocalhostCABundle is the asset the generates the kube-apiserver-localhost-ca-bundle,
// which contains all the individual client CAs.
type KubeAPIServerLocalhostCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerLocalhostCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerLocalhostCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLocalhostSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerLocalhostCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-localhost-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerLocalhostCABundle) Name() string {
	return "Certificate (kube-apiserver-localhost-ca-bundle)"
}

// KubeAPIServerLocalhostServerCertKey is the asset that generates the kube-apiserver serving key/cert pair for SNI localhost.
type KubeAPIServerLocalhostServerCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeAPIServerLocalhostServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeAPIServerLocalhostServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLocalhostSignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeAPIServerLocalhostServerCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeAPIServerLocalhostSignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityOneDay(installConfig),
		DNSNames: []string{
			"localhost",
		},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-apiserver-localhost-server", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerLocalhostServerCertKey) Name() string {
	return "Certificate (kube-apiserver-localhost-server)"
}

// KubeAPIServerServiceNetworkSignerCertKey is a key/cert pair that signs the kube-apiserver server cert for SNI service network.
type KubeAPIServerServiceNetworkSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeAPIServerServiceNetworkSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeAPIServerServiceNetworkSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeAPIServerServiceNetworkSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-apiserver-service-network-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kube-apiserver-service-network-signer")
}

// Load reads the asset files from disk.
func (c *KubeAPIServerServiceNetworkSignerCertKey) Load(f asset.FileFetcher) (bool, error) {
	return c.loadCertKey(f, "kube-apiserver-service-network-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeAPIServerServiceNetworkSignerCertKey) Name() string {
	return "Certificate (kube-apiserver-service-network-signer)"
}

// KubeAPIServerServiceNetworkCABundle is the asset the generates the kube-apiserver-service-network-ca-bundle,
// which contains all the individual client CAs.
type KubeAPIServerServiceNetworkCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerServiceNetworkCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerServiceNetworkCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerServiceNetworkSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerServiceNetworkCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-service-network-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerServiceNetworkCABundle) Name() string {
	return "Certificate (kube-apiserver-service-network-ca-bundle)"
}

// KubeAPIServerServiceNetworkServerCertKey is the asset that generates the kube-apiserver serving key/cert pair for SNI service network.
type KubeAPIServerServiceNetworkServerCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeAPIServerServiceNetworkServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeAPIServerServiceNetworkServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerServiceNetworkSignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeAPIServerServiceNetworkServerCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeAPIServerServiceNetworkSignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)
	/*
		serviceAddress, err := cidrhost(installConfig.Config.Networking.ServiceNetwork[0].IPNet, 1)
		if err != nil {
			return errors.Wrap(err, "failed to get service address for kube-apiserver from InstallConfig")
		}
	*/

	serviceAddresses := make([]net.IP, len(installConfig.Config.Networking.ServiceNetwork))
	for i, svcNet := range installConfig.Config.Networking.ServiceNetwork {
		serviceAddress, err := cidrhost(svcNet.IPNet, 1)
		if err != nil {
			return fmt.Errorf("failed to get service address for kube-apiserver from InstallConfig: %w", err)
		}
		serviceAddresses[i] = net.ParseIP(serviceAddress)
	}

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityOneDay(installConfig),
		DNSNames: []string{
			"kubernetes", "kubernetes.default",
			"kubernetes.default.svc",
			"kubernetes.default.svc.cluster.local",
			"openshift", "openshift.default",
			"openshift.default.svc",
			"openshift.default.svc.cluster.local",
		},
		IPAddresses: serviceAddresses,
		//IPAddresses: []net.IP{net.ParseIP(serviceAddress)},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-apiserver-service-network-server", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerServiceNetworkServerCertKey) Name() string {
	return "Certificate (kube-apiserver-service-network-server)"
}

// KubeAPIServerLBSignerCertKey is a key/cert pair that signs the kube-apiserver server cert for SNI load balancer.
type KubeAPIServerLBSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*KubeAPIServerLBSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *KubeAPIServerLBSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *KubeAPIServerLBSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "kube-apiserver-lb-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "kube-apiserver-lb-signer")
}

// Load reads the asset files from disk.
func (c *KubeAPIServerLBSignerCertKey) Load(f asset.FileFetcher) (bool, error) {
	return c.loadCertKey(f, "kube-apiserver-lb-signer")
}

// Name returns the human-friendly name of the asset.
func (c *KubeAPIServerLBSignerCertKey) Name() string {
	return "Certificate (kube-apiserver-lb-signer)"
}

// KubeAPIServerLBCABundle is the asset the generates the kube-apiserver-lb-ca-bundle,
// which contains all the individual client CAs.
type KubeAPIServerLBCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerLBCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerLBCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLBSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerLBCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-lb-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerLBCABundle) Name() string {
	return "Certificate (kube-apiserver-lb-ca-bundle)"
}

// KubeAPIServerExternalLBServerCertKey is the asset that generates the kube-apiserver serving key/cert pair for SNI external load balancer.
type KubeAPIServerExternalLBServerCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeAPIServerExternalLBServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeAPIServerExternalLBServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLBSignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeAPIServerExternalLBServerCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeAPIServerLBSignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityOneDay(installConfig),
		DNSNames: []string{
			apiAddress(installConfig.Config),
		},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-apiserver-lb-server", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerExternalLBServerCertKey) Name() string {
	return "Certificate (kube-apiserver-external-lb-server)"
}

// KubeAPIServerInternalLBServerCertKey is the asset that generates the kube-apiserver serving key/cert pair for SNI internal load balancer.
type KubeAPIServerInternalLBServerCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*KubeAPIServerInternalLBServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *KubeAPIServerInternalLBServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLBSignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *KubeAPIServerInternalLBServerCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &KubeAPIServerLBSignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityOneDay(installConfig),
		DNSNames: []string{
			internalAPIAddress(installConfig.Config),
		},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "kube-apiserver-internal-lb-server", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerInternalLBServerCertKey) Name() string {
	return "Certificate (kube-apiserver-internal-lb-server)"
}

// KubeAPIServerCompleteCABundle is the asset the generates the kube-apiserver-complete-server-ca-bundle,
// which contains all the certs that are valid to confirm the kube-apiserver identity.
type KubeAPIServerCompleteCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerCompleteCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerCompleteCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeAPIServerLocalhostCABundle{},
		&KubeAPIServerServiceNetworkCABundle{},
		&KubeAPIServerLBCABundle{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerCompleteCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-complete-server-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerCompleteCABundle) Name() string {
	return "Certificate (kube-apiserver-complete-server-ca-bundle)"
}

// KubeAPIServerCompleteClientCABundle is the asset the generates the kube-apiserver-complete-client-ca-bundle,
// which contains all the certs that are valid for the kube-apiserver to trust for clients.
type KubeAPIServerCompleteClientCABundle struct {
	CertBundle
}

var _ asset.Asset = (*KubeAPIServerCompleteClientCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *KubeAPIServerCompleteClientCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AdminKubeConfigCABundle{},        // admin.kubeconfig
		&KubeletClientCABundle{},          // signed kubelet certs
		&KubeControlPlaneCABundle{},       // controller-manager, scheduler
		&KubeAPIServerToKubeletCABundle{}, // kube-apiserver to kubelet (kubelet piggy-backs on KAS client-ca)
		&KubeletBootstrapCABundle{},       // used to create the kubelet kubeconfig files that are used to create CSRs
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *KubeAPIServerCompleteClientCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "kube-apiserver-complete-client-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *KubeAPIServerCompleteClientCABundle) Name() string {
	return "Certificate (kube-apiserver-complete-client-ca-bundle)"
}
