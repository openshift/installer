package main

import (
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/agent"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/password"
)

func newAgentCmd() *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Commands for supporting cluster installation using agent installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	agentCmd.AddCommand(newAgentCreateCmd())
	agentCmd.AddCommand(agent.NewWaitForCmd())
	return agentCmd
}

var (
	agentConfigTarget = target{
		// TODO: remove template wording when interactive survey has been implemented
		name: "Agent Config Template",
		command: &cobra.Command{
			Use:   "agent-config-template",
			Short: "Generates a template of the agent config manifest used by the agent installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&agentconfig.AgentConfig{},
		},
	}

	agentManifestsTarget = target{
		name: "Cluster Manifests",
		command: &cobra.Command{
			Use:   "cluster-manifests",
			Short: "Generates the cluster definition manifests used by the agent installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&manifests.AgentManifests{},
			&mirror.RegistriesConf{},
			&mirror.CaBundle{},
		},
	}

	agentImageTarget = target{
		name: "Agent ISO Image",
		command: &cobra.Command{
			Use:   "image",
			Short: "Generates a bootable image containing all the information needed to deploy a cluster",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&image.AgentImage{},
			&kubeconfig.AgentAdminClient{},
			&password.KubeadminPassword{},
		},
	}

	//nolint:varcheck,deadcode
	agentPXEFilesTarget = target{
		name: "Agent PXE Files",
		command: &cobra.Command{
			Use:   "pxe-files",
			Short: "Generates PXE bootable image files containing all the information needed to deploy a cluster",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&image.AgentPXEFiles{},
			&kubeconfig.AgentAdminClient{},
			&password.KubeadminPassword{},
		},
	}

	agentTargets = []target{agentConfigTarget, agentManifestsTarget, agentImageTarget}
)

func newAgentCreateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating agent installation artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range agentTargets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}
