package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/kubeconfig"

	agentpkg "github.com/openshift/installer/pkg/agent"
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
	agentCmd.AddCommand(newAgentWaitForCmd())
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
			&image.AgentImage{},
			&kubeconfig.AgentAdminClient{},
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

func newAgentWaitForCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait-for",
		Short: "Wait for install-time events",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newAgentWaitForClusterValidationSuccessCmd())
	cmd.AddCommand(newAgentWaitForBootstrapCompleteCmd())
	cmd.AddCommand(newAgentWaitForInstallCompleteCmd())
	return cmd
}

func newAgentWaitForClusterValidationSuccessCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster-validated",
		Short: "Wait until the cluster manifests are validated for install",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runAgentWaitForClusterValidationSuccessCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newAgentWaitForBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until the cluster bootstrap is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runAgentWaitForBootstrapCompleteCmd(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newAgentWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster installation is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runAgentWaitForInstallCompleteCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

}

func runAgentWaitForClusterValidationSuccessCmd() error {
	return agentpkg.WaitForClusterValidationSuccess()
}

func runAgentWaitForBootstrapCompleteCmd(directory string) error {
	return agentpkg.WaitForBootstrapComplete(directory)
}

func runAgentWaitForInstallCompleteCmd() error {
	return agentpkg.WaitForInstallComplete()
}
