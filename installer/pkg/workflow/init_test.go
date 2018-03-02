package workflow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func initTestCluster(t *testing.T, file string) config.Cluster {
	testConfig, err := config.ParseFile(file)
	if err != nil {
		t.Errorf("Test case TestGenerateTerraformVariablesStep: failed to parse test config, %s", err)
	}
	return testConfig.Clusters[0]
}

func TestGenerateTerraformVariablesStep(t *testing.T) {
	cluster := initTestCluster(t, "./fixtures/aws.basic.yaml")
	expectedTfVarsFilePath := "./fixtures/terraform.tfvars"
	clusterDir := "."
	m := &metadata{
		cluster:    cluster,
		clusterDir: clusterDir,
	}

	generateTerraformVariablesStep(m)
	gotTfVarsFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	gotData, err := ioutil.ReadFile(gotTfVarsFilePath)
	if err != nil {
		t.Errorf("Failed to load generated tf vars file: %s", err)
	}
	got := string(gotData)

	expectedData, err := ioutil.ReadFile(expectedTfVarsFilePath)
	if err != nil {
		t.Errorf("Failed to load expected tf vars file: %s", err)
	}
	expected := string(expectedData)

	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}

	// clean up
	if err := os.Remove(gotTfVarsFilePath); err != nil {
		t.Errorf("Failed to clean up generated tf vars file: %s", err)
	}
}
