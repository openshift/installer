package main

import (
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ibi"
)

func newImageBasedInstallCmd() *cobra.Command {
	ibiCmd := &cobra.Command{
		Use:   "ibi",
		Short: "Commands for supporting cluster installation using the Image-based installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	ibiCmd.AddCommand(newImageBasedInstallCreateCmd())
	return ibiCmd
}

var (
	imageBasedInstallConfigTarget = target{
		name: "Image-based Install Config Template",
		command: &cobra.Command{
			Use:   "ibi-config-template",
			Short: "Generates a template of the Image-based Install config manifest used by the Image-based installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&ibi.ImageBasedInstallConfig{},
		},
	}

	imageBasedInstallImageTarget = target{
		name: "Image-based installation ISO Image",
		command: &cobra.Command{
			Use:   "image",
			Short: "Generates a bootable image containing all the information needed to deploy a cluster",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&ibi.ImageBasedInstallImage{},
			// &kubeconfig.AgentAdminClient{},
			// &password.KubeadminPassword{},
		},
	}

	imageBasedInstallTargets = []target{imageBasedInstallConfigTarget, imageBasedInstallImageTarget}
)

func newImageBasedInstallCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating image-based installation artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range imageBasedInstallTargets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}
