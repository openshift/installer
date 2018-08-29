package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

type Stock interface {
	RootCA() asset.Asset
	KubeCA() asset.Asset
	EtcdCA() asset.Asset
	AggregatorCA() asset.Asset
	ServiceServingCA() asset.Asset
	EtcdClientCertKey() asset.Asset
	AdminCertKey() asset.Asset
	IngressCertKey() asset.Asset
	APIServerCertKey() asset.Asset
	OpenshiftAPIServerCertKey() asset.Asset
	APIServerProxyCertKey() asset.Asset
	KubeletCertKey() asset.Asset
	TNCCertKey() asset.Asset
	ClusterAPIServerCertKey() asset.Asset
	ServiceAccountKeyPair() asset.Asset
}

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

func (s *StockImpl) RootCA() asset.Asset                    { return s.rootCA }
func (s *StockImpl) KubeCA() asset.Asset                    { return s.kubeCA }
func (s *StockImpl) EtcdCA() asset.Asset                    { return s.etcdCA }
func (s *StockImpl) AggregatorCA() asset.Asset              { return s.aggregatorCA }
func (s *StockImpl) ServiceServingCA() asset.Asset          { return s.serviceServingCA }
func (s *StockImpl) EtcdClientCertKey() asset.Asset         { return s.etcdClientCertKey }
func (s *StockImpl) AdminCertKey() asset.Asset              { return s.adminCertKey }
func (s *StockImpl) IngressCertKey() asset.Asset            { return s.ingressCertKey }
func (s *StockImpl) APIServerCertKey() asset.Asset          { return s.apiServerCertKey }
func (s *StockImpl) OpenshiftAPIServerCertKey() asset.Asset { return s.openshiftAPIServerCertKey }
func (s *StockImpl) APIServerProxyCertKey() asset.Asset     { return s.apiServerProxyCertKey }
func (s *StockImpl) KubeletCertKey() asset.Asset            { return s.kubeletCertKey }
func (s *StockImpl) TNCCertKey() asset.Asset                { return s.tncCertKey }
func (s *StockImpl) ClusterAPIServerCertKey() asset.Asset   { return s.clusterAPIServerCertKey }
func (s *StockImpl) ServiceAccountKeyPair() asset.Asset     { return s.serviceAccountKeyPair }
