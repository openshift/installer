package manifests

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/stream-metadata-go/arch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
)

// InfraEnv generates the infraenv.yaml file.
type InfraEnv struct {
	InfraEnvFile
}

var _ asset.WritableAsset = (*InfraEnv)(nil)

// Name returns a human friendly name for the asset.
func (*InfraEnv) Name() string {
	return "InfraEnv Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*InfraEnv) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
		&agent.OptionalInstallConfig{},
		&agentconfig.AgentConfig{},
	}
}

// Generate generates the InfraEnv manifest.
func (i *InfraEnv) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	clusterInfo := &joiner.ClusterInfo{}
	installConfig := &agent.OptionalInstallConfig{}
	agentConfig := &agentconfig.AgentConfig{}
	dependencies.Get(installConfig, agentConfig, agentWorkflow, clusterInfo)

	rendezvousIP := ""
	if agentConfig.Config != nil {
		rendezvousIP = agentConfig.Config.RendezvousIP
	}
	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		if installConfig.Config != nil {
			err := i.generateManifest(installConfig.ClusterName(), installConfig.ClusterNamespace(), installConfig.Config.SSHKey, installConfig.Config.AdditionalTrustBundle, installConfig.Config.Proxy, string(installConfig.Config.ControlPlane.Architecture), &installConfig.Config.Networking.MachineNetwork, rendezvousIP)
			if err != nil {
				return err
			}

			if agentConfig.Config != nil {
				i.Config.Spec.AdditionalNTPSources = agentConfig.Config.AdditionalNTPSources
			}

			if installConfig.Config.BareMetal != nil {
				if i.Config.Spec.AdditionalNTPSources == nil && installConfig.Config.BareMetal.AdditionalNTPServers != nil {
					i.Config.Spec.AdditionalNTPSources = installConfig.Config.BareMetal.AdditionalNTPServers
				}
			}

		}

	case workflow.AgentWorkflowTypeAddNodes:
		err := i.generateManifest(clusterInfo.ClusterName, clusterInfo.Namespace, clusterInfo.SSHKey, clusterInfo.UserCaBundle, clusterInfo.Proxy, clusterInfo.Architecture, nil, rendezvousIP)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	return i.finish()
}

func (i *InfraEnv) generateManifest(clusterName, clusterNamespace, sshKey, additionalTrustBundle string, proxy *types.Proxy, architecture string, machineNetwork *[]types.MachineNetworkEntry, rendezvousIP string) error {
	infraEnv := &aiv1beta1.InfraEnv{
		TypeMeta: metav1.TypeMeta{
			Kind:       "InfraEnv",
			APIVersion: aiv1beta1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterName,
			Namespace: clusterNamespace,
		},
		Spec: aiv1beta1.InfraEnvSpec{
			ClusterRef: &aiv1beta1.ClusterReference{
				Name:      clusterName,
				Namespace: clusterNamespace,
			},
			SSHAuthorizedKey: strings.Trim(sshKey, "|\n\t"),
			PullSecretRef: &corev1.LocalObjectReference{
				Name: getPullSecretName(clusterName),
			},
			NMStateConfigLabelSelector: metav1.LabelSelector{
				MatchLabels: getNMStateConfigLabels(clusterName),
			},
		},
	}

	// Input values (amd64, arm64) must be converted to rpmArch because infraEnv.Spec.CpuArchitecture expects x86_64 or aarch64.
	if architecture != "" {
		infraEnv.Spec.CpuArchitecture = arch.RpmArch(architecture)
	}
	if additionalTrustBundle != "" {
		infraEnv.Spec.AdditionalTrustBundle = additionalTrustBundle
	}
	if proxy != nil {
		infraEnv.Spec.Proxy = getProxy(proxy, machineNetwork, rendezvousIP)
	}

	i.Config = infraEnv

	return nil
}
