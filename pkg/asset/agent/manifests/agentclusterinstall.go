package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

var (
	agentClusterInstallFilename = filepath.Join(clusterManifestDir, "agent-cluster-install.yaml")
)

// AgentClusterInstall generates the agent-cluster-install.yaml file.
type AgentClusterInstall struct {
	asset.DefaultFileWriter

	Config *hiveext.AgentClusterInstall
}

var _ asset.WritableAsset = (*AgentClusterInstall)(nil)

// Name returns a human friendly name for the asset.
func (*AgentClusterInstall) Name() string {
	return "AgentClusterInstall Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*AgentClusterInstall) Dependencies() []asset.Asset {
	return []asset.Asset{
		// &installconfig.InstallConfig{},
	}
}

// Generate generates the AgentClusterInstall manifest.
func (a *AgentClusterInstall) Generate(dependencies asset.Parents) error {
	// installConfig := &installconfig.InstallConfig{}
	// dependencies.Get(installConfig)

	// agentClusterInstall := &hiveext.AgentClusterInstall{
	// 	ObjectMeta: v1.ObjectMeta{
	// 		Name:      "agent-cluster-install",
	// 		Namespace: installConfig.Config.Namespace,
	// 	},
	// 	Spec: hiveext.AgentClusterInstallSpec{
	// 		ClusterDeploymentRef: corev1.LocalObjectReference{
	// 			Name: installConfig.Config.ObjectMeta.Name,
	// 		},
	// 		SSHPublicKey: installConfig.Config.SSHKey,
	// 		ProvisionRequirements: hiveext.ProvisionRequirements{
	// 			ControlPlaneAgents: int(*installConfig.Config.ControlPlane.Replicas),
	// 		},
	// 	},
	// }

	// var numberOfWorkers int = 0
	// for _, compute := range installConfig.Config.Compute {
	// 	numberOfWorkers = numberOfWorkers + int(*compute.Replicas)
	// }
	// agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents = numberOfWorkers

	// agentClusterInstallData, err := yaml.Marshal(agentClusterInstall)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to marshal agent installer AgentClusterInstall")
	// }

	// a.File = &asset.File{
	// 	Filename: agentClusterInstallFilename,
	// 	Data:     agentClusterInstallData,
	// }

	return nil
}

// Load returns agentclusterinstall asset from the disk.
func (a *AgentClusterInstall) Load(f asset.FileFetcher) (bool, error) {

	agentClusterInstallFile, err := f.FetchByName(agentClusterInstallFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", agentClusterInstallFilename))
	}

	a.File = agentClusterInstallFile

	agentClusterInstall := &hiveext.AgentClusterInstall{}
	if err := yaml.UnmarshalStrict(agentClusterInstallFile.Data, agentClusterInstall); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal %s", agentClusterInstallFilename)
		return false, err
	}
	a.Config = agentClusterInstall

	return true, nil
}
