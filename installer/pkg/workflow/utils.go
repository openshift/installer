package workflow

import (
	"os"
	"os/exec"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
)

func terraformExec(m *metadata, command ...string) error {
	command = append(command, tectonic.FindTemplatesForType(m.Cluster.Platform))

	tf := exec.Command("terraform", command...)
	tf.Dir = m.statePath
	tf.Stdin = os.Stdin
	tf.Stdout = os.Stdout
	tf.Stderr = os.Stderr
	if err := tf.Run(); err != nil {
		return err
	}

	return nil
}
