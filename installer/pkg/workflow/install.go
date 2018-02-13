package workflow

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
)

const (
	configFileName             = "config.yaml"
	terraformVariablesFileName = "terraform.tfvars"
)

// NewInstallWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallWorkflow(configFile string) Workflow {
	config, err := config.ParseFile(configFile)
	if err != nil {
		log.Fatalf("%s is not a valid config file: %s", configFile, err)
	}
	cluster := config.Clusters[0]

	return simpleWorkflow{
		metadata: metadata{
			Cluster:    cluster,
			configFile: configFile,
		},
		steps: []Step{
			tectonicPrepareStep,
			tectonicGenerateClusterConfig,
			tectonicGenerateTerraformVariables,
			terraformInitStep,
			terraformApplyStep,
		},
	}
}

func tectonicGenerateClusterConfig(m *metadata) error {
	return tectonic.GenerateClusterConfig(m.Cluster, m.statePath)
}

func tectonicGenerateTerraformVariables(m *metadata) error {
	configFilePath := filepath.Join(m.statePath, terraformVariablesFileName)

	return tectonic.GenerateTerraformVars(m.Cluster, configFilePath)
}

func tectonicPrepareStep(m *metadata) error {
	if m.statePath == "" {
		m.statePath = tectonic.NewBuildLocation(m.Cluster.Name)
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

func terraformApplyStep(m *metadata) error {
	log.Printf("Installation is running...")

	return terraformExec(m, "apply")
}

func terraformInitStep(m *metadata) error {
	log.Printf("Initializing cluster ...")

	return terraformExec(m, "init")
}
