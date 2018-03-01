package workflow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func initCluster(t *testing.T, file string) config.Cluster {
	testConfig, err := config.ParseFile(file)
	if err != nil {
		t.Errorf("Test case TestGenerateTerraformVariablesStep: failed to parse test config, %s", err)
	}
	return testConfig.Clusters[0]
}

func TestGenerateTerraformVariablesStep(t *testing.T) {

	cluster := initCluster(t, "./fixtures/aws.basic.yaml")
	clusterDir := "./"
	m := &metadata{
		cluster:    cluster,
		clusterDir: clusterDir,
	}
	generateTerraformVariablesStep(m)
	terraformVariablesFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	gotData, err := ioutil.ReadFile(terraformVariablesFilePath)
	if err != nil {
		t.Errorf("Test case TestGenerateTerraformVariablesStep: failed to ReadFile(): %s", err)
	}
	got := string(gotData)

	expectedData, err := ioutil.ReadFile("./fixtures/terraform.tfvars")
	if err != nil {
		t.Errorf("Test case TestGenerateTerraformVariablesStep: failed to ReadFile(): %s", err)
	}
	expected := string(expectedData)

	if got != expected {
		t.Errorf("Test case TestGenerateTerraformVariablesStep: expected: %s, got: %s", expected, got)
	}

	// clean up
	if err := os.Remove(terraformVariablesFilePath); err != nil {
		t.Errorf("TestGenerateTerraformVariablesStep: failed to clean up temp file: %s", err)
	}
}
