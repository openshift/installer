package configgenerator

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/tls"
)

const (
	rootCACertPath           = "generated/newTLS/root-ca.crt"
	rootCAKeyPath            = "generated/newTLS/root-ca.key"
	kubeCACertPath           = "generated/newTLS/kube-ca.key"
	kubeCAKeyPath            = "generated/newTLS/kube-ca.crt"
	aggregatorCAKeyPath      = "generated/newTLS/aggregator-ca.key"
	aggregatorCACertPath     = "generated/newTLS/aggregator-ca.crt"
	serviceServiceCAKeyPath  = "generated/newTLS/service-serving-ca.key"
	serviceServiceCACertPath = "generated/newTLS/service-serving-ca.crt"
	etcdClientKeyPath        = "generated/newTLS/etcd-client-ca.key"
	etcdClientCertPath       = "generated/newTLS/etcd-client-ca.crt"
)

// GenerateTLSConfig fetches and validates the TLS cert files
// If no file paths were provided, the certs will be auto-generated
func (c *ConfigGenerator) GenerateTLSConfig(clusterDir string) error {
	if c.CA.RootCAKeyPath == "" && c.CA.RootCACertPath == "" {
		// generate key and certificate
		key, err := generatePrivateKey(clusterDir, rootCAKeyPath)
		if err != nil {
			return fmt.Errorf("failed to generate private key: %v", err)
		}
		if _, err := generateRootCA(clusterDir, key); err != nil {
			return fmt.Errorf("failed to create a certificate: %v", err)
		}
	} else {
		// copy key and certificates
		keyDst := filepath.Join(clusterDir, rootCAKeyPath)
		if err := copyFile(c.CA.RootCAKeyPath, keyDst); err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}

		certDst := filepath.Join(clusterDir, rootCACertPath)
		if err := copyFile(c.CA.RootCACertPath, certDst); err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
	}
	return nil
}

func generateRootCA(path string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	fileTargetPath := filepath.Join(path, rootCACertPath)
	cfg := &tls.CertCfg{
		Subject: pkix.Name{
			CommonName:         "root-ca",
			OrganizationalUnit: []string{"openshift"},
		},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	cert, err := tls.SelfSignedCACert(cfg, key)
	if err != nil {
		return nil, fmt.Errorf("error generating self signed certificate: %v", err)
	}
	if err := writeFile(fileTargetPath, certToPem(cert)); err != nil {
		return nil, err
	}
	return cert, nil
}
