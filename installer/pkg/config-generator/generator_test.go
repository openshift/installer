package configgenerator

import (
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func initConfig(t *testing.T, file string) ConfigGenerator {
	testConfig, err := config.ParseFile("./fixtures/" + file)
	if err != nil {
		t.Errorf("Test case TestUrlFunctions: failed to parse test config, %s", err)
	}
	cluster := testConfig.Clusters[0]

	return ConfigGenerator{
		cluster,
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
			test:     "getApiServerUrl",
			got:      config.getApiServerUrl(),
			expected: "https://test-api.cluster.com:443",
		},
		{
			test:     "getBaseAddress",
			got:      config.getBaseAddress(),
			expected: "test.cluster.com",
		},
		{
			test:     "getOicdIssuerUrl",
			got:      config.getOicdIssuerUrl(),
			expected: "test.cluster.com/identity",
		},
	}
	for _, tc := range testCases {
		if tc.got != tc.expected {
			t.Errorf("Test case %s: expected: %s, got: %s", tc.test, tc.expected, tc.got)
		}
	}
}

func TestGetEtcdServersUrls(t *testing.T) {
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
		got := config.getEtcdServersUrls()
		if got != tc.expected {
			t.Errorf("Test case %s: expected: %s, got: %s", tc.test, tc.expected, got)
		}
	}

}
