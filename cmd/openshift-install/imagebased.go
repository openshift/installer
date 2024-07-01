package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/imagebased/image"
)

func newImageBasedCmd(ctx context.Context) *cobra.Command {
	imagebasedCmd := &cobra.Command{
		Use:   "imagebased",
		Short: "Commands for supporting cluster installation using the Image-based installer",
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
			Use:   "installation-config-template",
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

	imageBasedTargets = []target{
		imageBasedInstallationConfigTemplateTarget,
		imageBasedInstallationImageTarget,
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
