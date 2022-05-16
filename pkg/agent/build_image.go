package agent

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/agent/imagebuilder"
	"github.com/openshift/installer/pkg/agent/isosource"
)

// BuildImage builds the image required by the agent installer.
func BuildImage(assetsDir string) error {

	baseImage, err := isosource.EnsureIso()
	if err != nil {
		return err
	}

	err = imagebuilder.BuildImage(assetsDir, baseImage)
	if err != nil {
		return err
	}

	logrus.Info("BuildImage command")

	return nil
}
