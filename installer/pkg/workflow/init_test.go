package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/openshift/installer/installer/pkg/config"
)

func initTestCluster(cfg string) (*config.Cluster, error) {
	testConfig, err := config.ParseConfigFile(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test config: %v", err)
	}
	testConfig.PullSecret = "{\"auths\": {}}"
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

	cluster, err := initTestCluster("./fixtures/aws.basic.yaml")
	if err != nil {
		t.Fatalf("failed to init cluster: %v", err)
	}

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

	if got+"\n" != expected {
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
