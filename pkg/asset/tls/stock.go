package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Stock is the stock of TLS assets that can be generated.
type Stock interface {
	// RootCA is the asset that generates the root-ca key/cert pair.
	RootCA() asset.Asset
	// KubeCA is the asset that generates the kube-ca key/cert pair.
	KubeCA() asset.Asset
	// EtcdCA is the asset that generates the etcd-ca key/cert pair.
	EtcdCA() asset.Asset
	// AggregatorCA is the asset that generates the aggregator-ca key/cert pair.
	AggregatorCA() asset.Asset
	// ServiceServingCA is the asset that generates the service-serving-ca key/cert pair.
	ServiceServingCA() asset.Asset
	// EtcdClientCertKey is the asset that generates the etcd client key/cert pair.
	EtcdClientCertKey() asset.Asset
	// AdminCertKey is the asset that generates the admin key/cert pair.
	AdminCertKey() asset.Asset
	// IngressCertKey is the asset that generates the ingress key/cert pair.
	IngressCertKey() asset.Asset
	// APIServerCertKey is the asset that generates the API server key/cert pair.
	APIServerCertKey() asset.Asset
	// OpenshiftAPIServerCertKey is the asset that generates the Openshift API server key/cert pair.
	OpenshiftAPIServerCertKey() asset.Asset
	// APIServerProxyCertKey is the asset that generates the API server proxy key/cert pair.
	APIServerProxyCertKey() asset.Asset
	// KubeletCertKey is the asset that generates the kubelet key/cert pair.
	KubeletCertKey() asset.Asset
	// TNCCertKey is the asset that generates the TNC key/cert pair.
	TNCCertKey() asset.Asset
	// ClusterAPIServerCertKey is the asset that generates the cluster API server key/cert pair.
	ClusterAPIServerCertKey() asset.Asset
	// ServiceAccountKeyPair is the asset that generates the service-account public/private key pair.
	ServiceAccountKeyPair() asset.Asset
}

// StockImpl implements the Stock interface for tls assets.
type StockImpl struct {
	rootCA                    asset.Asset
	kubeCA                    asset.Asset
	etcdCA                    asset.Asset
	aggregatorCA              asset.Asset
	serviceServingCA          asset.Asset
	etcdClientCertKey         asset.Asset
	adminCertKey              asset.Asset
	ingressCertKey            asset.Asset
	apiServerCertKey          asset.Asset
	openshiftAPIServerCertKey asset.Asset
	apiServerProxyCertKey     asset.Asset
	kubeletCertKey            asset.Asset
	tncCertKey                asset.Asset
	clusterAPIServerCertKey   asset.Asset
	serviceAccountKeyPair     asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(rootDir string, stock installconfig.Stock) {
	s.rootCA = &RootCA{rootDir: rootDir}
	s.kubeCA = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "kube-ca", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:      ValidityTenYears,
		KeyFileName:   "kube-ca.key",
		CertFileName:  "kube-ca.crt",

		IsCA:     true,
		ParentCA: s.rootCA,
	}

	s.etcdCA = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:      ValidityTenYears,
		KeyFileName:   "etcd-ca.key",
		CertFileName:  "etcd-ca.crt",

		IsCA:     true,
		ParentCA: s.rootCA,
	}

	s.aggregatorCA = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "aggregator", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:      ValidityTenYears,
		KeyFileName:   "aggregator-ca.key",
		CertFileName:  "aggregator-ca.crt",

		IsCA:     true,
		ParentCA: s.rootCA,
	}

	s.serviceServingCA = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "service-serving", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:      ValidityTenYears,
		KeyFileName:   "service-serving-ca.key",
		CertFileName:  "service-serving-ca.crt",

		IsCA:     true,
		ParentCA: s.rootCA,
	}

	s.etcdClientCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:      ValidityTenYears,
		KeyFileName:   "etcd-client.key",
		CertFileName:  "etcd-client.crt",

		ParentCA: s.etcdCA,
	}

	s.adminCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "system:admin", Organization: []string{"system:masters"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:      ValidityTenYears,
		KeyFileName:   "admin.key",
		CertFileName:  "admin.crt",

		ParentCA: s.kubeCA,
	}

	s.ingressCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:      ValidityTenYears,
		KeyFileName:   "ingress.key",
		CertFileName:  "ingress.crt",

		ParentCA:     s.kubeCA,
		AppendParent: true,
		GenSubject:   genSubjectForIngressCertKey,
		GenDNSNames:  genDNSNamesForIngressCertKey,
	}

	s.apiServerCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Subject:       pkix.Name{CommonName: "kube-apiserver", Organization: []string{"kube-master"}},
		Validity:      ValidityTenYears,
		KeyFileName:   "apiserver.key",
		CertFileName:  "apiserver.crt",

		ParentCA:       s.kubeCA,
		AppendParent:   true,
		GenDNSNames:    genDNSNamesForAPIServerCertKey,
		GenIPAddresses: genIPAddressesForAPIServerCertKey,
	}

	s.openshiftAPIServerCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Subject:       pkix.Name{CommonName: "openshift-apiserver", Organization: []string{"kube-master"}},
		Validity:      ValidityTenYears,
		KeyFileName:   "openshift-apiserver.key",
		CertFileName:  "openshift-apiserver.crt",

		ParentCA:       s.aggregatorCA,
		AppendParent:   true,
		GenDNSNames:    genDNSNamesForOpenshiftAPIServerCertKey,
		GenIPAddresses: genIPAddressesForOpenshiftAPIServerCertKey,
	}

	s.apiServerProxyCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Subject:       pkix.Name{CommonName: "kube-apiserver-proxy", Organization: []string{"kube-master"}},
		Validity:      ValidityTenYears,
		KeyFileName:   "apiserver-proxy.key",
		CertFileName:  "apiserver-proxy.crt",

		ParentCA: s.aggregatorCA,
	}

	s.kubeletCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Subject:       pkix.Name{CommonName: "system:serviceaccount:kube-system:default", Organization: []string{"system:serviceaccounts:kube-system"}},
		Validity:      ValidityThirtyMinutes,
		KeyFileName:   "kubelet.key",
		CertFileName:  "kubelet.crt",

		ParentCA: s.kubeCA,
	}

	s.tncCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:      ValidityTenYears,
		KeyFileName:   "tnc.key",
		CertFileName:  "tnc.crt",

		ParentCA:    s.rootCA,
		GenDNSNames: genDNSNamesForTNCCertKey,
		GenSubject:  genSubjectForTNCCertKey,
	}

	s.clusterAPIServerCertKey = &CertKey{
		rootDir:       rootDir,
		installConfig: stock.InstallConfig(),
		Subject:       pkix.Name{CommonName: "cluster-apiserver", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:      ValidityTenYears,
		KeyFileName:   "cluster-apiserver-ca.key",
		CertFileName:  "cluster-apiserver-ca.crt",
		IsCA:          true,

		ParentCA:     s.aggregatorCA,
		AppendParent: true,
	}

	s.serviceAccountKeyPair = &KeyPair{
		rootDir:         rootDir,
		PrivKeyFileName: "service-account.key",
		PubKeyFileName:  "service-account.pub",
	}
}

// RootCA is the asset that generates the root-ca key/cert pair.
func (s *StockImpl) RootCA() asset.Asset { return s.rootCA }

// KubeCA is the asset that generates the kube-ca key/cert pair.
func (s *StockImpl) KubeCA() asset.Asset { return s.kubeCA }

// EtcdCA is the asset that generates the etcd-ca key/cert pair.
func (s *StockImpl) EtcdCA() asset.Asset { return s.etcdCA }

// AggregatorCA is the asset that generates the aggregator-ca key/cert pair.
func (s *StockImpl) AggregatorCA() asset.Asset { return s.aggregatorCA }

// ServiceServingCA is the asset that generates the service-serving-ca key/cert pair.
func (s *StockImpl) ServiceServingCA() asset.Asset { return s.serviceServingCA }

// EtcdClientCertKey is the asset that generates the etcd client key/cert pair.
func (s *StockImpl) EtcdClientCertKey() asset.Asset { return s.etcdClientCertKey }

// AdminCertKey is the asset that generates the admin key/cert pair.
func (s *StockImpl) AdminCertKey() asset.Asset { return s.adminCertKey }

// IngressCertKey is the asset that generates the ingress key/cert pair.
func (s *StockImpl) IngressCertKey() asset.Asset { return s.ingressCertKey }

// APIServerCertKey is the asset that generates the API server key/cert pair.
func (s *StockImpl) APIServerCertKey() asset.Asset { return s.apiServerCertKey }

// OpenshiftAPIServerCertKey is the asset that generates the Openshift API server key/cert pair.
func (s *StockImpl) OpenshiftAPIServerCertKey() asset.Asset { return s.openshiftAPIServerCertKey }

// APIServerProxyCertKey is the asset that generates the API server proxy key/cert pair.
func (s *StockImpl) APIServerProxyCertKey() asset.Asset { return s.apiServerProxyCertKey }

// KubeletCertKey is the asset that generates the kubelet key/cert pair.
func (s *StockImpl) KubeletCertKey() asset.Asset { return s.kubeletCertKey }

// TNCCertKey is the asset that generates the TNC key/cert pair.
func (s *StockImpl) TNCCertKey() asset.Asset { return s.tncCertKey }

// ClusterAPIServerCertKey is the asset that generates the cluster API server key/cert pair.
func (s *StockImpl) ClusterAPIServerCertKey() asset.Asset { return s.clusterAPIServerCertKey }

// ServiceAccountKeyPair is the asset that generates the service-account public/private key pair.
func (s *StockImpl) ServiceAccountKeyPair() asset.Asset { return s.serviceAccountKeyPair }
