package main

import (
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/agent"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
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
	agentManifestsTarget = target{
		name: "Cluster Manifests",
		command: &cobra.Command{
			Use:   "cluster-manifests",
			Short: "Generates the cluster definition manifests used by the agent installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&manifests.AgentManifests{},
		},
	}

	agentImageTarget = target{
		name: "Image",
		command: &cobra.Command{
			Use:   "image",
			Short: "Generates a bootable image containing all the information needed to deploy a cluster",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&image.ISO{},
		},
	}

	agentTargets = []target{agentManifestsTarget, agentImageTarget}
)

func newAgentCreateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating agent installation artifacts",
		Args:  cobra.ExactArgs(0),
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
