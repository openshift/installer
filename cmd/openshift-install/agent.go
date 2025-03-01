package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/agent"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/configimage"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/password"
)

func newAgentCmd(ctx context.Context) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Commands for supporting cluster installation using agent installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	agentCmd.AddCommand(newAgentCreateCmd(ctx))
	agentCmd.AddCommand(agent.NewWaitForCmd())
	agentCmd.AddCommand(newAgentGraphCmd())
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

	agentConfigImageTarget = target{
		name: "Agent Config Image",
		command: &cobra.Command{
			Use:   "config-image",
			Short: "Generates an ISO containing configuration files only",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&configimage.ConfigImage{},
			&kubeconfig.AgentAdminClient{},
			&password.KubeadminPassword{},
		},
	}

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

	agentUnconfiguredIgnitionTarget = target{
		name: "Agent unconfigured ignition",
		command: &cobra.Command{
			Use:    "unconfigured-ignition",
			Short:  "Generates an agent ignition that excludes cluster configuration",
			Args:   cobra.ExactArgs(0),
			Hidden: true,
		},
		assets: []asset.WritableAsset{
			&image.UnconfiguredIgnition{},
		},
	}

	agentInteractiveDisconnectedIgnitionTarget = target{
		name: "Agent interactive disconnected install ignition",
		command: &cobra.Command{
			Use:    "interactive-disconnected-ignition",
			Short:  "Generates the agent ignition required to support the interactive disconnected installation",
			Args:   cobra.ExactArgs(0),
			Hidden: true,
		}, assets: []asset.WritableAsset{
			&workflow.AgentWorkflowInstallInteractiveDisconnected{},
			&image.InteractiveDisconnectedIgnition{},
		},
	}

	agentTargets = []target{
		agentConfigTarget,
		agentManifestsTarget,
		agentImageTarget,
		agentPXEFilesTarget,
		agentConfigImageTarget,
		agentUnconfiguredIgnitionTarget,
		agentInteractiveDisconnectedIgnitionTarget,
	}
)

func newAgentCreateCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating agent installation artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range agentTargets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(ctx, t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}

func newAgentGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Outputs the internal dependency graph for the agent-based installer",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGraphCmd(cmd, args, agentTargets)
		},
	}
	cmd.PersistentFlags().StringVar(&graphOpts.outputFile, "output-file", "", "file where the graph is written, if empty prints the graph to Stdout.")
	return cmd
}
