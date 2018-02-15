package tectonic

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
	"github.com/coreos/tectonic-installer/installer/pkg/config-generator"
	"github.com/coreos/tectonic-installer/installer/pkg/terraform-generator"
)

const (
	kubeSystemFileName     = "kube-system.yml"
	tectonicSystemFileName = "tectonic-system.yml"
)

// NewBuildLocation creates a new directory on disk that will become
// the root location for all statefull artefacts of the current cluster build.
func NewBuildLocation(clusterName string) string {
	var err error
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory because: %v", err)
	}
	buildPath := filepath.Join(pwd, clusterName)
	err = os.MkdirAll(buildPath, os.ModeDir|0755)
	if err != nil {
		log.Fatalf("Failed to create build folder at %s", buildPath)
	}
	return buildPath
}

// FindTemplatesForType determines the location of top-level
// Terraform templates for a given type (platform) of build.
// TEMPORARY: implement actual detection of templates from released artefacts.
func FindTemplatesForType(buildType string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "platforms", strings.ToLower(buildType))
}

// GenerateClusterConfig writes, if successful, the cluster configuration.
func GenerateClusterConfig(cluster config.Cluster, configPath string) error {
	configGenerator := configgenerator.New(cluster)

	kubeSystem, err := configGenerator.KubeSystem()
	if err != nil {
		return err
	}

	kubeSystemConfigFilePath := filepath.Join(configPath, kubeSystemFileName)
	if err := writeFile(kubeSystemConfigFilePath, kubeSystem); err != nil {
		return err
	}

	tectonicSystem, err := configGenerator.TectonicSystem()
	if err != nil {
		return err
	}

	tectonicSystemConfigFilePath := filepath.Join(configPath, tectonicSystemFileName)
	return writeFile(tectonicSystemConfigFilePath, tectonicSystem)
}

// GenerateTerraformVars writes, if successful, the terraform variables.
func GenerateTerraformVars(cluster config.Cluster, configFilePath string) error {
	terraformGenerator := terraformgenerator.New(cluster)

	vars, err := terraformGenerator.TFVars()
	if err != nil {
		return err
	}

	return writeFile(configFilePath, vars)
}

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := fmt.Fprintln(w, content); err != nil {
		return err
	}
	w.Flush()

	return nil
}

// FindTemplatesForStep determines the location of top-level
// Terraform templates for a given step of build.
func FindTemplatesForStep(step ...string) string {
	pwd, _ := os.Getwd()
	step = append([]string{pwd, "steps"}, step...)
	return filepath.Join(step...)
}
