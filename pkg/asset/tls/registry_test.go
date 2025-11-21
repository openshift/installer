package tls

import (
	"context"
	"crypto/x509"
	"net"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInternalReleaseRegistrySignerCertKeyGenerate(t *testing.T) {
	ca := &InternalReleaseRegistrySignerCertKey{}
	parents := asset.Parents{}

	if err := ca.Generate(context.TODO(), parents); err != nil {
		t.Fatalf("failed to generate internal release registry signer cert key: %v", err)
	}

	if len(ca.Cert()) == 0 {
		t.Error("expected certificate data")
	}
	if len(ca.Key()) == 0 {
		t.Error("expected key data")
	}

	cert, err := PemToCertificate(ca.Cert())
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	if !cert.IsCA {
		t.Error("expected certificate to be a CA")
	}
	if cert.Subject.CommonName != "internal-release-registry-signer" {
		t.Errorf("expected CommonName to be 'internal-release-registry-signer', got %s", cert.Subject.CommonName)
	}
}

func TestInternalReleaseRegistryCertKeyGenerate(t *testing.T) {
	ca := &InternalReleaseRegistrySignerCertKey{}
	if err := ca.Generate(context.TODO(), asset.Parents{}); err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	installConfig := &installconfig.InstallConfig{}
	installConfig.Config = &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		BaseDomain: "test.example.com",
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{{
				CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
			}},
			ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
		},
	}

	parents := asset.Parents{}
	parents.Add(ca, installConfig)

	serverCert := &InternalReleaseRegistryCertKey{}
	if err := serverCert.Generate(context.TODO(), parents); err != nil {
		t.Fatalf("failed to generate server certificate: %v", err)
	}

	if len(serverCert.Cert()) == 0 {
		t.Error("expected certificate data")
	}
	if len(serverCert.Key()) == 0 {
		t.Error("expected key data")
	}

	cert, err := PemToCertificate(serverCert.Cert())
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	// Check SANs include localhost and api-int
	expectedDNSNames := map[string]bool{
		"localhost":                         true,
		"api-int.test-cluster.test.example.com": true,
	}

	for _, dnsName := range cert.DNSNames {
		if !expectedDNSNames[dnsName] {
			t.Errorf("unexpected DNS name in certificate: %s", dnsName)
		}
		delete(expectedDNSNames, dnsName)
	}

	if len(expectedDNSNames) > 0 {
		t.Errorf("missing expected DNS names: %v", expectedDNSNames)
	}

	// Check IP addresses
	expectedIPs := map[string]bool{
		"127.0.0.1": true,
		"::1":       true,
	}

	for _, ipAddr := range cert.IPAddresses {
		ipStr := ipAddr.String()
		if !expectedIPs[ipStr] {
			t.Errorf("unexpected IP address in certificate: %s", ipStr)
		}
		delete(expectedIPs, ipStr)
	}

	if len(expectedIPs) > 0 {
		t.Errorf("missing expected IP addresses: %v", expectedIPs)
	}

	// Verify ExtKeyUsage includes ServerAuth
	hasServerAuth := false
	for _, usage := range cert.ExtKeyUsage {
		if usage == x509.ExtKeyUsageServerAuth {
			hasServerAuth = true
			break
		}
	}
	if !hasServerAuth {
		t.Error("expected certificate to have ExtKeyUsageServerAuth")
	}
}

func TestInternalReleaseRegistryLocalhostCertKeyGenerate(t *testing.T) {
	ca := &InternalReleaseRegistrySignerCertKey{}
	if err := ca.Generate(context.TODO(), asset.Parents{}); err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	parents := asset.Parents{}
	parents.Add(ca)

	localhostCert := &InternalReleaseRegistryLocalhostCertKey{}
	if err := localhostCert.Generate(context.TODO(), parents); err != nil {
		t.Fatalf("failed to generate localhost certificate: %v", err)
	}

	if len(localhostCert.Cert()) == 0 {
		t.Error("expected certificate data")
	}
	if len(localhostCert.Key()) == 0 {
		t.Error("expected key data")
	}

	cert, err := PemToCertificate(localhostCert.Cert())
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	// Check SANs include only localhost
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "localhost" {
		t.Errorf("expected DNSNames to be [localhost], got %v", cert.DNSNames)
	}

	// Check IP addresses
	expectedIPs := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	if len(cert.IPAddresses) != len(expectedIPs) {
		t.Errorf("expected %d IP addresses, got %d", len(expectedIPs), len(cert.IPAddresses))
	}

	// Verify ExtKeyUsage includes ServerAuth
	hasServerAuth := false
	for _, usage := range cert.ExtKeyUsage {
		if usage == x509.ExtKeyUsageServerAuth {
			hasServerAuth = true
			break
		}
	}
	if !hasServerAuth {
		t.Error("expected certificate to have ExtKeyUsageServerAuth")
	}
}

func TestInternalReleaseRegistrySignerCertKeyLoadFromDisk(t *testing.T) {
	ca := &InternalReleaseRegistrySignerCertKey{}

	// Before loading, LoadedFromDisk should be false
	if ca.LoadedFromDisk {
		t.Error("expected LoadedFromDisk to be false before loading")
	}

	// Generate a certificate
	if err := ca.Generate(context.TODO(), asset.Parents{}); err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	// LoadedFromDisk should still be false after generation
	if ca.LoadedFromDisk {
		t.Error("expected LoadedFromDisk to be false after generation")
	}

	// Note: Testing actual file loading would require a FileFetcher mock
	// and is typically done in integration tests
}
