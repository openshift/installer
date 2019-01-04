package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/destroy/bootstrap"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
	_ "github.com/openshift/installer/pkg/destroy/openstack"
)

func newDestroyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy part of an OpenShift cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newDestroyBootstrapCmd())
	cmd.AddCommand(newDestroyClusterCmd())
	return cmd
}

func newDestroyClusterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster",
		Short: "Destroy an OpenShift cluster",
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			err := runDestroyCmd(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func runDestroyCmd(directory string) error {
	destroyer, err := destroy.New(logrus.StandardLogger(), directory)
	if err != nil {
		return errors.Wrap(err, "Failed while preparing to destroy cluster")
	}
	if err := destroyer.Run(); err != nil {
		return errors.Wrap(err, "Failed to destroy cluster")
	}

	store, err := asset.NewStore(directory)
	if err != nil {
		return errors.Wrapf(err, "failed to create asset store")
	}
	for _, asset := range clusterTarget.assets {
		if err := store.Destroy(asset); err != nil {
			return errors.Wrapf(err, "failed to destroy asset %q", asset.Name())
		}
	}
	return nil
}

func newDestroyBootstrapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Destroy the bootstrap resources",
		Run: func(cmd *cobra.Command, args []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			err := bootstrap.Destroy(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}
