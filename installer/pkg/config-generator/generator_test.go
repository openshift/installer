package configgenerator

import (
	"io/ioutil"
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func initConfig(t *testing.T, file string) ConfigGenerator {
	cluster, err := config.ParseConfigFile("./fixtures/" + file)
	if err != nil {
		t.Errorf("Test case TestUrlFunctions: failed to parse test config, %s", err)
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
		if tc.got != tc.expected {
			t.Errorf("Test case %s: expected: %s, got: %s", tc.test, tc.expected, tc.got)
		}
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
		if got != tc.expected {
			t.Errorf("Test case %s: expected: %s, got: %s", tc.test, tc.expected, got)
		}
	}
}

func TestKubeSystem(t *testing.T) {
	config := initConfig(t, "test-aws.yaml")
	got, err := config.KubeSystem()
	if err != nil {
		t.Errorf("Test case TestKubeSystem: failed to get KubeSystem(): %s", err)
	}
	expected, err := ioutil.ReadFile("./fixtures/kube-system.yaml")
	if err != nil {
		t.Errorf("Test case TestKubeSystem: failed to ReadFile(): %s", err)
	}

	if got != string(expected) {
		t.Errorf("Test case TestKubeSystem: expected: %s, got: %s", expected, got)
	}
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
		if got != tc.expected {
			t.Errorf("Test case %s: expected: %s, got: %s", tc.test, tc.expected, got)
		}
	}
}
