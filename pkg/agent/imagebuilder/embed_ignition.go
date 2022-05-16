package imagebuilder

import (
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
)

const (
	outputImage = "output/agent.iso"
)

// BuildImage builds an ISO with ignition content from a base image, and writes
// the result to disk.
func BuildImage(baseImage string) error {
	configBuilder, err := New()
	if err != nil {
		return err
	}

	ignition, err := configBuilder.Ignition()
	if err != nil {
		return err
	}
	ignitionContent := &isoeditor.IgnitionContent{Config: ignition}

	custom, err := isoeditor.NewRHCOSStreamReader(baseImage, ignitionContent, nil)
	if err != nil {
		return err
	}
	defer custom.Close()

	output, err := os.Create(outputImage)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, custom)
	return err
}
