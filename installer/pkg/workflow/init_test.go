package workflow

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	jose "gopkg.in/square/go-jose.v2"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func generatePullSecretAndLicense(name string, expiration time.Time) (*os.File, *os.File, error) {
	pullBytes, err := json.Marshal(&struct{}{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal pull secret: %v", err)
	}
	p, err := ioutil.TempFile("", fmt.Sprintf("%s_pull_secret", name))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pull secret file: %v", err)
	}
	if _, err := p.Write(pullBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to write pull secret file: %v", err)
	}
	p.Close()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create RSA key pair: %v", err)
	}
	s, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create license signer: %v", err)
	}
	buf, err := json.Marshal(struct {
		ExpirationDate time.Time `json:"expirationDate"`
	}{expiration})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal license: %v", err)
	}
	jws, err := s.Sign(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign license: %v", err)
	}
	license, err := jws.CompactSerialize()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to serialize license: %v", err)
	}
	l, err := ioutil.TempFile("", fmt.Sprintf("%s_license", name))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create license file: %v", err)
	}
	if _, err := l.WriteString(license); err != nil {
		return nil, nil, fmt.Errorf("failed to write license file: %v", err)
	}
	l.Close()

	return p, l, nil
}

func initTestCluster(cfg, pullSecret, license string) (*config.Cluster, error) {
	testConfig, err := config.ParseConfigFile(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test config: %v", err)
	}
	testConfig.PullSecretPath = pullSecret
	testConfig.LicensePath = license
	if len(testConfig.Validate()) != 0 {
		return nil, errors.New("failed to validate test conifg")
	}
	return testConfig, nil
}

func TestGenerateTerraformVariablesStep(t *testing.T) {
	expectedTfVarsFilePath := "./fixtures/terraform.tfvars"
	clusterDir := "."
	gotTfVarsFilePath := filepath.Join(clusterDir, terraformVariablesFileName)

	// clean up
	defer func() {
		if err := os.Remove(gotTfVarsFilePath); err != nil {
			t.Errorf("failed to clean up generated tf vars file: %v", err)
		}
	}()

	ps, lic, err := generatePullSecretAndLicense("init_workflow", time.Now().AddDate(1, 0, 0))
	if err != nil {
		t.Fatalf("failed to generate pull secret and license: %v", err)
	}
	defer os.Remove(ps.Name())
	defer os.Remove(lic.Name())

	cluster, err := initTestCluster("./fixtures/aws.basic.yaml", ps.Name(), lic.Name())
	if err != nil {
		t.Fatalf("failed to init cluster: %v", err)
	}

	// Remove auto-generated license and pull secret for comparison.
	cluster.LicensePath = ""
	cluster.PullSecretPath = ""

	m := &metadata{
		cluster:    *cluster,
		clusterDir: clusterDir,
	}

	generateTerraformVariablesStep(m)
	gotData, err := ioutil.ReadFile(gotTfVarsFilePath)
	if err != nil {
		t.Errorf("failed to load generated tf vars file: %v", err)
	}
	got := string(gotData)

	expectedData, err := ioutil.ReadFile(expectedTfVarsFilePath)
	if err != nil {
		t.Errorf("failed to load expected tf vars file: %v", err)
	}
	expected := string(expectedData)

	if got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}

func TestBuildInternalConfig(t *testing.T) {
	testClusterDir := "."
	internalFilePath := filepath.Join(testClusterDir, internalFileName)

	// clean up
	defer func() {
		if err := os.Remove(internalFilePath); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	errorTestCases := []struct {
		test     string
		got      string
		expected string
	}{
		{
			test:     "no clusterDir exists",
			got:      buildInternalConfig("").Error(),
			expected: "no cluster dir given for building internal config",
		},
	}

	for _, tc := range errorTestCases {
		if tc.got != tc.expected {
			t.Errorf("test case %s: expected: %s, got: %s", tc.test, tc.expected, tc.got)
		}
	}

	if err := buildInternalConfig(testClusterDir); err != nil {
		t.Errorf("failed to run buildInternalStep, %v", err)
	}

	if _, err := os.Stat(internalFilePath); err != nil {
		t.Errorf("failed to create internal file, %v", err)
	}

	testInternal, err := config.ParseInternalFile(internalFilePath)
	if err != nil {
		t.Errorf("failed to parse internal file, %v", err)
	}
	testCases := []struct {
		test     string
		got      string
		expected string
	}{
		{
			test:     "clusterId",
			got:      testInternal.ClusterID,
			expected: "^[a-zA-Z0-9_-]*$",
		},
	}

	for _, tc := range testCases {
		match, _ := regexp.MatchString(tc.expected, tc.got)
		if !match {
			t.Errorf("test case %s: expected: %s, got: %s", tc.test, tc.expected, tc.got)
		}
	}
}
