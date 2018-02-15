package workflow

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
)

const (
	configFileName             = "config.yaml"
	terraformVariablesFileName = "terraform.tfvars"
	kubeConfig                 = "/generated/auth/kubeconfig"
)

// NewInstallWorkflow creates new instances of the 'install' workflow,
// responsible for running the actions necessary to install a new cluster.
func NewInstallWorkflow(configFile string) Workflow {

	// TODO: move to tectonicGenerateClusterConfig/tectonicGenerateTerraformVariables and get this dynamically
	clusterName := "cluster-aws"
	platform := "aws"

	if platform == "aws" {
		return simpleWorkflow{
			metadata: metadata{
				clusterName: clusterName,
				configFile:  configFile,
			},
			steps: []Step{
				terraformPrepareStep,
				assetsStep,
				bootstrapStep,
				joiningStep,
			},
		}
	}
	return simpleWorkflow{
		metadata: metadata{
			clusterName: clusterName,
			configFile:  configFile,
			platform:    platform,
		},
		steps: []Step{
			terraformPrepareStep,
			terraformInitStep,
			terraformApplyStep,
		},
	}
}

//func tectonicGenerateClusterConfig(m *metadata) error {
//	return tectonic.GenerateClusterConfig(m.Cluster, m.statePath)
//}
//
//func tectonicGenerateTerraformVariables(m *metadata) error {
//	configFilePath := filepath.Join(m.statePath, terraformVariablesFileName)
//
//	return tectonic.GenerateTerraformVars(m.Cluster, configFilePath)
//}

// NewAssetsWorkflow creates new instances of the 'assets' workflow,
// responsible for running the actions necessary to generate cluster assets.
func NewAssetsWorkflow(configFile string) Workflow {
	// TODO: move to tectonicGenerateClusterConfig/tectonicGenerateTerraformVariables and get this dynamically
	clusterName := "cluster-aws"
	return simpleWorkflow{
		metadata: metadata{
			clusterName: clusterName,
			configFile:  configFile,
		},
		steps: []Step{
			terraformPrepareStep,
			assetsStep,
		},
	}
}

// NewBootstrapWorkflow creates new instances of the 'bootstrap' workflow,
// responsible for running the actions necessary to generate a single bootstrap machine cluster.
func NewBootstrapWorkflow(configFile string) Workflow {
	// TODO: move to tectonicGenerateClusterConfig/tectonicGenerateTerraformVariables and get this dynamically
	clusterName := "cluster-aws"
	return simpleWorkflow{
		metadata: metadata{
			clusterName: clusterName,
			configFile:  configFile,
		},
		steps: []Step{
			terraformPrepareStep,
			bootstrapStep,
		},
	}
}

// NewJoinWorkflow creates new instances of the 'join' workflow,
// responsible for running the actions necessary to scale the machines of the cluster.
func NewJoinWorkflow(configFile string) Workflow {
	// TODO: move to tectonicGenerateClusterConfig/tectonicGenerateTerraformVariables and get this dynamically
	clusterName := "cluster-aws"
	return simpleWorkflow{
		metadata: metadata{
			clusterName: clusterName,
			configFile:  configFile,
		},
		steps: []Step{
			terraformPrepareStep,
			joiningStep,
		},
	}
}

func terraformPrepareStep(m *metadata) error {
	if m.statePath == "" {
		m.statePath = tectonic.NewBuildLocation(m.clusterName)
	}
	varfile := filepath.Join(m.statePath, m.configFile)
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
	return tfInit(m.statePath, tectonic.FindTemplatesForType(m.platform))
}

func terraformApplyStep(m *metadata) error {
	log.Printf("Installation is running...")
	return tfApply(m.statePath, "state", tectonic.FindTemplatesForType(m.platform))
}

func assetsStep(m *metadata) error {
	log.Printf("Installation is running...")
	return runStep(m.statePath, "assets")
}

func bootstrapStep(m *metadata) error {
	log.Printf("Installation is running...")
	err := runStep(m.statePath, "bootstrap")
	if err != nil {
		return err
	}
	err = waitForNcg(m)
	if err != nil {
		return err
	}
	err = destroyCname(m)
	if err != nil {
		return err
	}
	return nil
}

func joiningStep(m *metadata) error {
	// TODO: import will fail after a first run, error is ignored for now
	importAutoScalingGroup(m)
	log.Printf("Installation is running...")
	return runStep(m.statePath, "joining")
}

func runStep(buildPath string, step string) error {
	codePath := tectonic.FindTemplatesForStep(step)
	err := tfInit(buildPath, codePath)
	if err != nil {
		return err
	}

	err = tfApply(buildPath, step, codePath)
	if err != nil {
		return err
	}
	return nil
}
