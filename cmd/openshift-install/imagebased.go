package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/imagebased/configimage"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/password"
)

func newImageBasedInstallCmd(ctx context.Context) *cobra.Command {
	imagebasedCmd := &cobra.Command{
		Use:   "imagebased",
		Short: "Commands for supporting cluster installation using the Image-based installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	imagebasedCmd.AddCommand(newImageBasedInstallCreateCmd(ctx))
	return imagebasedCmd
}

var (
	imageBasedConfigTemplateTarget = target{
		name: "Image-based Installer Config ISO Configuration Template",
		command: &cobra.Command{
			Use:   "config-template",
			Short: "Generates a template of the Image-based Config ISO config manifest used by the Image-based installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&configimage.ImageBasedConfig{},
		},
	}

	imageBasedConfigImageTarget = target{
		name: "Image-based Installer Config ISO Image",
		command: &cobra.Command{
			Use:   "config-image",
			Short: "Generates an ISO containing configuration files only",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&configimage.ConfigImage{},
			&kubeconfig.ImageBasedAdminClient{},
			&password.KubeadminPassword{},
		},
	}

	imageBasedInstallTargets = []target{
		imageBasedConfigTemplateTarget,
		imageBasedConfigImageTarget,
	}
)

func newImageBasedInstallCreateCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating image-based installer artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range imageBasedInstallTargets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(ctx, t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}
