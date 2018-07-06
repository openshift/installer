package configgenerator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/openshift/installer/installer/pkg/tls"
)

const (
	aggregatorCACertPath     = "generated/newTLS/aggregator-ca.crt"
	aggregatorCAKeyPath      = "generated/newTLS/aggregator-ca.key"
	etcdClientCertPath       = "generated/newTLS/etcd-client-ca.crt"
	etcdClientKeyPath        = "generated/newTLS/etcd-client-ca.key"
	kubeCACertPath           = "generated/newTLS/kube-ca.crt"
	kubeCAKeyPath            = "generated/newTLS/kube-ca.key"
	rootCACertPath           = "generated/newTLS/root-ca.crt"
	rootCAKeyPath            = "generated/newTLS/root-ca.key"
	serviceServiceCACertPath = "generated/newTLS/service-serving-ca.crt"
	serviceServiceCAKeyPath  = "generated/newTLS/service-serving-ca.key"
	ingressCACertPath        = "generated/newTLS/ingress-ca.crt"
	ingressCertPath          = "generated/newTLS/ingress.crt"
	ingressKeyPath           = "generated/newTLS/ingress.key"
)

// GenerateTLSConfig fetches and validates the TLS cert files
// If no file paths were provided, the certs will be auto-generated
func (c *ConfigGenerator) GenerateTLSConfig(clusterDir string) error {
	var caKey *rsa.PrivateKey
	var caCert *x509.Certificate
	var kubeCAKey *rsa.PrivateKey
	var kubeCACert *x509.Certificate
	var err error
	var baseAddress string

	baseAddress = c.getBaseAddress()

	if c.CA.RootCAKeyPath == "" && c.CA.RootCACertPath == "" {
		caCert, caKey, err = generateRootCert(clusterDir)
		if err != nil {
			return fmt.Errorf("failed to generate root CA certificate and key pair: %v", err)
		}
	} else {
		// copy key and certificates
		caCert, caKey, err = getCertFiles(clusterDir, c.CA.RootCACertPath, c.CA.RootCAKeyPath)
		if err != nil {
			return fmt.Errorf("failed to process CA certificate and key pair: %v", err)
		}
	}

	// generate kube CA
	cfg := &tls.CertCfg{
		Subject:   pkix.Name{CommonName: "kube-ca", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:      true,
	}
	kubeCAKey, kubeCACert, err = generateCert(clusterDir, caKey, caCert, kubeCAKeyPath, kubeCACertPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate kube CAs: %v", err)
	}
	// generate aggregator CA
	cfg = &tls.CertCfg{
		Subject:   pkix.Name{CommonName: "aggregator", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:      true,
	}
	if _, _, err := generateCert(clusterDir, caKey, caCert, aggregatorCAKeyPath, aggregatorCACertPath, cfg); err != nil {
		return fmt.Errorf("failed to generate aggregator CAs: %v", err)
	}

	// generate service-serving CA
	cfg = &tls.CertCfg{
		Subject:   pkix.Name{CommonName: "service-serving", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:      true,
	}
	if _, _, err := generateCert(clusterDir, caKey, caCert, serviceServiceCAKeyPath, serviceServiceCACertPath, cfg); err != nil {
		return fmt.Errorf("failed to generate service-serving CAs: %v", err)
	}

	// Ingress certs
	if copyFile(kubeCACertPath, ingressCACertPath); err != nil {
		return fmt.Errorf("failed to import kube CA cert into ingress-ca.crt: %v", err)
	}

	cfg = &tls.CertCfg{
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:     []string{baseAddress, fmt.Sprintf("%s.%s", "*", baseAddress)},
		Subject:      pkix.Name{CommonName: baseAddress, Organization: []string{"ingress"}},
		IsCA:         false}

	if _, _, err := generateCert(clusterDir, kubeCAKey, kubeCACert, ingressKeyPath, ingressCertPath, cfg); err != nil {
		return fmt.Errorf("failed to generate ingress CAs: %v", err)
	}
	return nil
}

// generatePrivateKey generates and writes the private key to disk
func generatePrivateKey(clusterDir string, path string) (*rsa.PrivateKey, error) {
	fileTargetPath := filepath.Join(clusterDir, path)
	key, err := tls.PrivateKey()
	if err != nil {
		return nil, fmt.Errorf("error writing private key: %v", err)
	}
	if err := writeFile(fileTargetPath, tls.PrivateKeyToPem(key)); err != nil {
		return nil, err
	}
	return key, nil
}

// generateRootCert creates the rootCAKey and rootCACert
func generateRootCert(clusterDir string) (cert *x509.Certificate, key *rsa.PrivateKey, err error) {
	// generate key and certificate
	caKey, err := generatePrivateKey(clusterDir, rootCAKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	caCert, err := generateRootCA(clusterDir, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a certificate: %v", err)
	}
	return caCert, caKey, nil
}

// getCertFiles copies the given cert/key files into the generated folder and returns their contents
func getCertFiles(clusterDir string, certPath string, keyPath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	keyDst := filepath.Join(clusterDir, rootCAKeyPath)
	if err := copyFile(keyPath, keyDst); err != nil {
		return nil, nil, fmt.Errorf("failed to write file: %v", err)
	}

	certDst := filepath.Join(clusterDir, rootCACertPath)
	if err := copyFile(certPath, certDst); err != nil {
		return nil, nil, fmt.Errorf("failed to write file: %v", err)
	}
	// content validation occurs in pkg/config/validate.go
	// if it fails here, something went wrong
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		panic(err)
	}
	certPem, _ := pem.Decode([]byte(string(certData)))
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	keyPem, _ := pem.Decode([]byte(string(keyData)))
	key, err := x509.ParsePKCS1PrivateKey(keyPem.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process private key: %v", err)
	}
	certs, err := x509.ParseCertificates(certPem.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process certificate: %v", err)
	}

	return certs[0], key, nil
}

// generateCert creates a key, csr & a signed cert
func generateCert(clusterDir string,
	caKey *rsa.PrivateKey,
	caCert *x509.Certificate,
	keyPath string,
	certPath string,
	cfg *tls.CertCfg) (*rsa.PrivateKey, *x509.Certificate, error) {

	// create a private key
	key, err := generatePrivateKey(clusterDir, keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// create a CSR
	csrTmpl := x509.CertificateRequest{Subject: cfg.Subject, DNSNames: cfg.DNSNames}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating certificate request: %v", err)
	}
	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing certificate request: %v", err)
	}

	// create a cert
	cert, err := generateSignedCert(cfg, csr, key, caKey, caCert, clusterDir, certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a certificate: %v", err)
	}
	return key, cert, nil
}

// generateRootCA creates and returns the root CA
func generateRootCA(path string, key *rsa.PrivateKey) (*x509.Certificate, error) {
	fileTargetPath := filepath.Join(path, rootCACertPath)
	cfg := &tls.CertCfg{
		Subject: pkix.Name{
			CommonName:         "root-ca",
			OrganizationalUnit: []string{"openshift"},
		},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:      true,
	}
	cert, err := tls.SelfSignedCACert(cfg, key)
	if err != nil {
		return nil, fmt.Errorf("error generating self signed certificate: %v", err)
	}
	if err := writeFile(fileTargetPath, tls.CertToPem(cert)); err != nil {
		return nil, err
	}
	return cert, nil
}

func generateSignedCert(cfg *tls.CertCfg,
	csr *x509.CertificateRequest,
	key *rsa.PrivateKey,
	caKey *rsa.PrivateKey,
	caCert *x509.Certificate,
	clusterDir string,
	path string) (*x509.Certificate, error) {
	cert, err := tls.SignedCertificate(cfg, csr, key, caCert, caKey)
	if err != nil {
		return nil, fmt.Errorf("error signing certificate: %v", err)
	}
	fileTargetPath := filepath.Join(clusterDir, path)
	if err := writeFile(fileTargetPath, tls.CertToPem(cert)); err != nil {
		return nil, err
	}
	return cert, nil
}
