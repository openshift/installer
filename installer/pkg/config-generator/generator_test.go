package configgenerator

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types/config"
	"github.com/stretchr/testify/assert"
)

func initConfig(t *testing.T, file string) ConfigGenerator {
	cluster, err := config.ParseConfigFile("./fixtures/" + file)
	if err != nil {
		t.Fatalf("Test case TestUrlFunctions: failed to parse test config, %s", err)
	}

	return ConfigGenerator{
		*cluster,
	}
}
func TestUrlFunctions(t *testing.T) {
	config := initConfig(t, "test.yaml")

	testCases := []struct {
		test     string
		got      string
		expected string
	}{
		{
			test:     "getAPIServerURL",
			got:      config.getAPIServerURL(),
			expected: "https://test-api.cluster.com:6443",
		},
		{
			test:     "getBaseAddress",
			got:      config.getBaseAddress(),
			expected: "test.cluster.com",
		},
		{
			test:     "getOicdIssuerURL",
			got:      config.getOicdIssuerURL(),
			expected: "https://test.cluster.com/identity",
		},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.got)
	}
}

func TestGetEtcdServersURLs(t *testing.T) {
	testCases := []struct {
		test       string
		configFile string
		expected   string
	}{
		{
			test:       "No ExternalServers",
			configFile: "test.yaml",
			expected:   "https://test-etcd-0.cluster.com:2379,https://test-etcd-1.cluster.com:2379,https://test-etcd-2.cluster.com:2379",
		},
	}
	for _, tc := range testCases {

		config := initConfig(t, tc.configFile)
		got := config.getEtcdServersURLs()
		assert.Equal(t, tc.expected, got)
	}
}

func TestKubeSystem(t *testing.T) {
	config := initConfig(t, "test-aws.yaml")
	got, err := config.KubeSystem("./fixtures")
	if err != nil {
		t.Errorf("Test case TestKubeSystem: failed to get KubeSystem(): %s", err)
	}
	expected, err := ioutil.ReadFile("./fixtures/kube-system.yaml")
	if err != nil {
		t.Errorf("Test case TestKubeSystem: failed to ReadFile(): %s", err)
	}

	assert.Equal(t, string(expected), got)
}

func TestCIDRHost(t *testing.T) {
	testCases := []struct {
		test     string
		iprange  string
		hostNum  int
		expected string
	}{
		{
			test:     "10.0.0.0/8",
			iprange:  "10.0.0.0/8",
			hostNum:  8,
			expected: "10.0.0.8",
		},
		{
			test:     "10.3.0.0/16",
			iprange:  "10.3.0.0/16",
			hostNum:  10,
			expected: "10.3.0.10",
		},
	}
	for _, tc := range testCases {
		got, err := cidrhost(tc.iprange, tc.hostNum)
		if err != nil {
			t.Errorf("Test case %s: failed to run cidrhost(): %s", tc.test, err)
		}
		assert.Equal(t, tc.expected, got)
	}
}

func TestGenerateCert(t *testing.T) {
	caKey, err := tls.PrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate Private Key: %v", err)
	}
	caCfg := &tls.CertCfg{
		Subject: pkix.Name{
			CommonName:         "test-self-signed-ca",
			OrganizationalUnit: []string{"openshift"},
		},
		Validity: tls.ValidityTenYears,
	}
	caCert, err := tls.SelfSignedCACert(caCfg, caKey)
	if err != nil {
		t.Fatalf("failed to generate self signed certificate: %v", err)
	}
	keyPath := "./test.key"
	certPath := "./test.crt"

	cases := []struct {
		cfg        *tls.CertCfg
		clusterDir string
		err        bool
	}{
		{
			cfg: &tls.CertCfg{
				Subject:      pkix.Name{CommonName: "test-cert", OrganizationalUnit: []string{"test"}},
				KeyUsages:    x509.KeyUsageKeyEncipherment,
				DNSNames:     []string{"test-api.kubernetes.default"},
				ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
				Validity:     tls.ValidityTenYears,
				IsCA:         false,
			},
			clusterDir: "./",
			err:        false,
		},
	}
	for i, c := range cases {
		_, _, err := generateCert(c.clusterDir, caKey, caCert, keyPath, certPath, c.cfg, false)
		if err != nil {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}

		if err := os.Remove(keyPath); err != nil {
			t.Errorf("test case %d: failed to cleanup test key: %s, got %v", i, keyPath, err)
		}
		if err := os.Remove(certPath); err != nil {
			t.Errorf("test case %d: failed to cleanup test certificate: %s, got %v", i, certPath, err)
		}
	}
}

func TestLibvirtURI(t *testing.T) {
	escapedPKIPath := url.QueryEscape(libvirtPKIPath)

	cases := []struct {
		label    string
		uri      string
		network  string
		expected string
		err      error
	}{
		{
			label:    "defaults",
			uri:      "qemu:///system",
			network:  "192.168.124.0/24",
			expected: "qemu+tls://192.168.124.1/system?pkipath=" + escapedPKIPath,
			err:      nil,
		},
		{
			label:    "qemu-localhost",
			uri:      "qemu://127.0.0.1/system",
			network:  "192.168.124.0/24",
			expected: "qemu+tls://192.168.124.1/system?pkipath=" + escapedPKIPath,
			err:      nil,
		},
		{
			label:    "custom-network",
			uri:      "qemu:///system",
			network:  "172.16.128.0/17",
			expected: "qemu+tls://172.16.128.1/system?pkipath=" + escapedPKIPath,
			err:      nil,
		},
		{
			label:    "preserve-query",
			uri:      "qemu:///system?foo=bar",
			network:  "192.168.124.0/24",
			expected: "qemu+tls://192.168.124.1/system?foo=bar&pkipath=" + escapedPKIPath,
			err:      nil,
		},
		{
			label:    "ssh-scheme",
			uri:      "qemu+ssh://127.0.0.1/system",
			network:  "192.168.124.0/24",
			expected: "qemu+tls://192.168.124.1/system?pkipath=" + escapedPKIPath,
			err:      nil,
		},
		{
			label:    "bad-scheme",
			uri:      "qemu+foo+bar:///system",
			network:  "192.168.124.0/24",
			expected: "",
			err:      errBadLibvirtScheme,
		},
	}

	for _, tt := range cases {
		t.Run(tt.label, func(t *testing.T) {
			uri, err := libvirtURI(tt.uri, tt.network)
			if err != tt.err {
				t.Errorf("unexpected error: %v", err)
			}

			if uri == nil || err != nil {
				return
			}

			if uri.String() != tt.expected {
				t.Errorf("got %q, want %q", uri.String(), tt.expected)
			}
		})
	}
}
