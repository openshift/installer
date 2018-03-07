package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

func initTestCluster(file string) (*config.Cluster, error) {
	testConfig, err := config.ParseFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse test config: %v", err)
	}
	if len((&testConfig.Clusters[0]).Validate()) != 0 {
		return nil, errors.New("failed to validate test conifg")
	}
	return &testConfig.Clusters[0], nil
}

func TestGenerateTerraformVariablesStep(t *testing.T) {
	cluster, err := initTestCluster("./fixtures/aws.basic.yaml")
	if err != nil {
		t.Errorf("failed to init cluster: %v", err)
	}
	expectedTfVarsFilePath := "./fixtures/terraform.tfvars"
	clusterDir := "."
	m := &metadata{
		cluster:    *cluster,
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
