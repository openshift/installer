package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/destroy/bootstrap"
	quotaasset "github.com/openshift/installer/pkg/destroy/quota"
	"github.com/openshift/installer/pkg/metrics/timer"

	_ "github.com/openshift/installer/pkg/destroy/alibabacloud"
	_ "github.com/openshift/installer/pkg/destroy/aws"
	_ "github.com/openshift/installer/pkg/destroy/azure"
	_ "github.com/openshift/installer/pkg/destroy/baremetal"
	_ "github.com/openshift/installer/pkg/destroy/gcp"
	_ "github.com/openshift/installer/pkg/destroy/ibmcloud"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
	_ "github.com/openshift/installer/pkg/destroy/nutanix"
	_ "github.com/openshift/installer/pkg/destroy/openstack"
	_ "github.com/openshift/installer/pkg/destroy/ovirt"
	_ "github.com/openshift/installer/pkg/destroy/powervs"
	_ "github.com/openshift/installer/pkg/destroy/vsphere"
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
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			err := runDestroyCmd(rootOpts.dir, os.Getenv("OPENSHIFT_INSTALL_REPORT_QUOTA_FOOTPRINT") == "true")
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Uninstallation complete!")
		},
	}
}

func runDestroyCmd(directory string, reportQuota bool) error {
	timer.StartTimer(timer.TotalTimeElapsed)
	destroyer, err := destroy.New(logrus.StandardLogger(), directory)
	if err != nil {
		return errors.Wrap(err, "Failed while preparing to destroy cluster")
	}
	quota, err := destroyer.Run()
	if err != nil {
		return errors.Wrap(err, "Failed to destroy cluster")
	}

	if reportQuota {
		if err := quotaasset.WriteQuota(directory, quota); err != nil {
			return errors.Wrap(err, "failed to record quota")
		}
	}

	store, err := assetstore.NewStore(directory)
	if err != nil {
		return errors.Wrap(err, "failed to create asset store")
	}
	for _, asset := range clusterTarget.assets {
		if err := store.Destroy(asset); err != nil {
			return errors.Wrapf(err, "failed to destroy asset %q", asset.Name())
		}
	}

	// delete the state file as well
	err = store.DestroyState()
	if err != nil {
		return errors.Wrap(err, "failed to remove state file")
	}

	// delete terraform files
	tfstateFiles, err := filepath.Glob(filepath.Join(directory, "*.tfstate"))
	if err != nil {
		return errors.Wrap(err, "failed to glob for tfstate files")
	}
	tfvarsFiles, err := filepath.Glob(filepath.Join(directory, "*.tfvars.json"))
	if err != nil {
		return errors.Wrap(err, "failed to glob for tfvars files")
	}
	for _, f := range append(tfstateFiles, tfvarsFiles...) {
		if err := os.Remove(f); err != nil {
			return errors.Wrapf(err, "failed to remove terraform file %q", f)
		}
	}

	timer.StopTimer(timer.TotalTimeElapsed)
	timer.LogSummary()

	return nil
}

func newDestroyBootstrapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Destroy the bootstrap resources",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			timer.StartTimer(timer.TotalTimeElapsed)
			err := bootstrap.Destroy(rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
}
