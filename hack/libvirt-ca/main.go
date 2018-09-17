package main

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/installer/pkg/asset/tls"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	orgUnit          = "openshift"
	caCommonName     = "libvirt"
	serverCommonName = "libvirt"
	clientCommonName = "openshift-cluster-api"
	defaultNetwork   = "192.168.124.0/24"
	defaultOutDir    = "."

	// Libvirt is picky about the naming here:
	// https://libvirt.org/remote.html
	caCertFile     = "cacert.pem"
	caKeyFile      = "cakey.pem"
	serverCertFile = "servercert.pem"
	serverKeyFile  = "serverkey.pem"
	clientCertFile = "clientcert.pem"
	clientKeyFile  = "clientkey.pem"
)

var (
	network = kingpin.Flag("network", "Cluster network CIDR.").
		Short('n').Default(defaultNetwork).String()

	outDir = kingpin.Flag("out", "Output directory.").
		Short('o').Default(defaultOutDir).ExistingDir()
)

type certificate struct {
	key  *rsa.PrivateKey
	cert *x509.Certificate
}

func (c *certificate) WritePEMs(certPath, keyPath string) error {
	if err := ioutil.WriteFile(keyPath, []byte(tls.PrivateKeyToPem(c.key)), 0600); err != nil {
		return err
	}

	if err := ioutil.WriteFile(certPath, []byte(tls.CertToPem(c.cert)), 0644); err != nil {
		return err
	}

	return nil
}

func getGateway(network string) (net.IP, error) {
	_, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}

	gateway, err := cidr.Host(ipNet, 1)
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

func generateCA() (*certificate, error) {
	var ca certificate
	var err error

	cfg := &tls.CertCfg{
		Subject: pkix.Name{
			CommonName:         caCommonName,
			OrganizationalUnit: []string{orgUnit},
		},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  tls.ValidityTenYears,
		IsCA:      true,
	}

	ca.key, ca.cert, err = tls.GenerateRootCertKey(cfg)
	if err != nil {
		return nil, err
	}

	return &ca, err
}

func generateServerCert(dnsNames []string, ips []net.IP, ca *certificate) (*certificate, error) {
	var server certificate
	var err error

	cfg := &tls.CertCfg{
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Subject: pkix.Name{
			CommonName:         serverCommonName,
			OrganizationalUnit: []string{orgUnit},
		},
		DNSNames:    dnsNames,
		Validity:    tls.ValidityTenYears,
		IPAddresses: ips,
		IsCA:        false,
	}

	server.key, server.cert, err = tls.GenerateCert(ca.key, ca.cert, cfg)
	if err != nil {
		return nil, err
	}

	return &server, err
}

func generateClientCert(ca *certificate) (*certificate, error) {
	var client certificate
	var err error

	cfg := &tls.CertCfg{
		Subject: pkix.Name{
			CommonName:         clientCommonName,
			OrganizationalUnit: []string{orgUnit},
		},
		KeyUsages:    x509.KeyUsageKeyEncipherment,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     tls.ValidityTenYears,
	}

	client.key, client.cert, err = tls.GenerateCert(ca.key, ca.cert, cfg)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func writeCertificates(dir string, ca, server, client *certificate) error {
	// Certificate Authority
	caCertPath := filepath.Join(dir, caCertFile)
	caKeyPath := filepath.Join(dir, caKeyFile)
	if err := ca.WritePEMs(caCertPath, caKeyPath); err != nil {
		return err
	}

	// Server Certificate
	serverCertPath := filepath.Join(dir, serverCertFile)
	serverKeyPath := filepath.Join(dir, serverKeyFile)
	if err := server.WritePEMs(serverCertPath, serverKeyPath); err != nil {
		return err
	}

	// Client Certificate
	clientCertPath := filepath.Join(dir, clientCertFile)
	clientKeyPath := filepath.Join(dir, clientKeyFile)
	if err := server.WritePEMs(clientCertPath, clientKeyPath); err != nil {
		return err
	}

	return nil
}

func main() {
	kingpin.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to determine hostname: %v\n", err)
		os.Exit(1)
	}

	gateway, err := getGateway(*network)
	if err != nil {
		fmt.Printf("Failed to read network: %v\n", err)
		os.Exit(1)
	}

	ca, err := generateCA()
	if err != nil {
		fmt.Printf("Failed to generate CA: %v\n", err)
		os.Exit(1)
	}

	server, err := generateServerCert([]string{hostname}, []net.IP{gateway}, ca)
	if err != nil {
		fmt.Printf("Failed to generate server certificate: %v\n", err)
		os.Exit(1)
	}

	client, err := generateClientCert(ca)
	if err != nil {
		fmt.Printf("Failed to generate client certificate: %v\n", err)
		os.Exit(1)
	}

	if err := writeCertificates(*outDir, ca, server, client); err != nil {
		fmt.Printf("Failed to write certificate files: %v\n", err)
		os.Exit(1)
	}
}
