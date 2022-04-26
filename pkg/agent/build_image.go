package agent

import (
	"github.com/sirupsen/logrus"
)

// BuildImage builds the image required by the agent installer.
func BuildImage() error {

	// baseImage, err := isosource.EnsureIso()
	// if err != nil {
	// 	return err
	// }

	// err = imagebuilder.BuildImage(baseImage)
	// if err != nil {
	// 	return err
	// }

	logrus.Info("BuildImage command")

	return nil
}
