package agent

import (
	"github.com/sirupsen/logrus"
)

// WaitFor wait for the installation complete triggered by the agent installer.
func WaitFor() error {
	logrus.Info("WaitFor command")

	return nil
}
