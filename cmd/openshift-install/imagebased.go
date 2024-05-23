package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/imagebased/configimage"
	"github.com/openshift/installer/pkg/asset/imagebased/image"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/password"
)

func newImageBasedCmd(ctx context.Context) *cobra.Command {
	imagebasedCmd := &cobra.Command{
		Use:   "image-based",
		Short: "Commands for supporting cluster installation using the image-based installer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	imagebasedCmd.AddCommand(newImageBasedCreateCmd(ctx))
	return imagebasedCmd
}

var (
	imageBasedInstallationConfigTemplateTarget = target{
		name: "Image-based Installation ISO Configuration template",
		command: &cobra.Command{
			Use:   "image-config-template",
			Short: "Generates a template of the Image-based Installation ISO config manifest used by the Image-based installer",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&image.ImageBasedInstallationConfig{},
		},
	}

	imageBasedInstallationImageTarget = target{
		name: "Image-based Installation ISO Image",
		command: &cobra.Command{
			Use:   "image",
			Short: "Generates a bootable ISO image containing all the information needed to deploy a cluster",
			Args:  cobra.ExactArgs(0),
		},
		assets: []asset.WritableAsset{
			&image.Image{},
		},
	}

	imageBasedConfigTemplateTarget = target{
		name: "Image-based Installer Config ISO Configuration Template",
		command: &cobra.Command{
			Use:   "config-template",
			Short: "Generates a template of the Image-based Config ISO config manifest used by the image-based installer",
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

	imageBasedTargets = []target{
		imageBasedInstallationConfigTemplateTarget,
		imageBasedInstallationImageTarget,
		imageBasedConfigTemplateTarget,
		imageBasedConfigImageTarget,
	}
)

func newImageBasedCreateCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating image-based installer artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range imageBasedTargets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(ctx, t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}
