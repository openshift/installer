package tectonic

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// ClusterNameFromConfig determines the name of the cluster form a
// configuration object specified by the user.
// TEMPORARY: This is a stub, until we wire in the cluster config object.
func ClusterNameFromConfig(varfile string) (string, error) {
	// TODO @spangenberg: implement this based on parsed config object.
	return "", errors.New("not found")
}

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
	return filepath.Join(pwd, "platforms", buildType)
}
