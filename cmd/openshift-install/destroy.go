package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/destroy"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
)

func newDestroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy-cluster",
		Short: "Destroy an OpenShift cluster",
		Long:  "",
		RunE:  runDestroyCmd,
	}
}

func runDestroyCmd(cmd *cobra.Command, args []string) error {
	destroyer, err := destroy.New(logrus.StandardLogger(), rootOpts.dir)
	if err != nil {
		return errors.Wrap(err, "Failed while preparing to destroy cluster")
	}
	if err := destroyer.Run(); err != nil {
		return errors.Wrap(err, "Failed to destroy cluster")

	}
	return nil
}
