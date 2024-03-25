package ibi

import (
	"os"

	"github.com/openshift/installer/pkg/asset"
)

// ImageBasedInstallArtifacts is an asset that generates all the artifacts
// that could be used for a subsequent generation of an ISO image, starting
// from the content of the RHCOS image enriched with IBI specific files.
type ImageBasedInstallArtifacts struct {
	TmpPath string
}

// Dependencies returns the assets on which the AgentArtifacts asset depends.
func (i *ImageBasedInstallArtifacts) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the configurations for the IBI ISO image assets.
func (i *ImageBasedInstallArtifacts) Generate(dependencies asset.Parents) error {
	// Create a tmp folder to store all the pieces required to generate the agent artifacts.
	tmpPath, err := os.MkdirTemp("", "ibi")
	if err != nil {
		return err
	}
	i.TmpPath = tmpPath

	return nil
}

// Name returns the human-friendly name of the asset.
func (i *ImageBasedInstallArtifacts) Name() string {
	return "Image-based Installer Artifacts"
}
