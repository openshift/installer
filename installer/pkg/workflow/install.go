package workflow

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
)

const configFileName string = "config.yaml"

// NewInstallWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallWorkflow(configFile string) Workflow {
	clusterName, err := tectonic.ClusterNameFromConfig(configFile) // TODO @spangenberg: re-implement with config object
	if err != nil {
		log.Fatalf("%s is not a valid config file", configFile)
	}
	return simpleWorkflow{
		metadata: metadata{
			clusterName: clusterName,
			configFile:  configFile,
		},
		steps: []Step{
			terraformPrepareStep,
			terraformInitStep,
			terraformApplyStep,
		},
	}
}

func terraformPrepareStep(m *metadata) error {
	if m.statePath == "" {
		m.statePath = tectonic.NewBuildLocation(m.clusterName)
	}
	varfile := filepath.Join(m.statePath, configFileName)
	if _, err := os.Stat(varfile); os.IsNotExist(err) {
		from, err := os.Open(m.configFile)
		if err != nil {
			return err
		}
		defer from.Close()
		to, err := os.OpenFile(varfile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer to.Close()
		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
	}
	return nil
}

func terraformInitStep(m *metadata) error {
	log.Printf("Initializing cluster ...")
	tfInit := exec.Command("terraform", "init", tectonic.FindTemplatesForType("aws")) // TODO: get from cluster config
	tfInit.Dir = m.statePath
	tfInit.Stdin = os.Stdin
	tfInit.Stdout = os.Stdout
	tfInit.Stderr = os.Stderr
	err := tfInit.Run()
	if err != nil {
		return err
	}
	return nil
}

func terraformApplyStep(m *metadata) error {
	log.Printf("Installation is running...")
	tfInit := exec.Command("terraform", "apply", tectonic.FindTemplatesForType("aws")) // TODO: get from cluster config
	tfInit.Dir = m.statePath
	tfInit.Stdin = os.Stdin
	tfInit.Stdout = os.Stdout
	tfInit.Stderr = os.Stderr
	err := tfInit.Run()
	if err != nil {
		return err
	}
	return nil
}
