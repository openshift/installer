package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	agentcmd "github.com/openshift/installer/pkg/agent"
)

// NewCreateCmd create the agent commands for generating the manifests and the bootable image.
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Commands for generating agent installer based artifacts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newCreateManifestsCmd())
	cmd.AddCommand(newCreateImageCmd())
	return cmd
}

func newCreateImageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image",
		Short: "Generates a bootable image containing all the information needed to deploy a cluster",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runCreateImageCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func runCreateImageCmd() error {
	return agentcmd.BuildImage()
}

func newCreateManifestsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "manifests",
		Short: "Generates intermediate files required by the agent installer",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runCreateManifestsCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func runCreateManifestsCmd() error {
	return agentcmd.CreateManifests()
}
