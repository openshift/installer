package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/command"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/clusterapi"
	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/destroy/bootstrap"
	quotaasset "github.com/openshift/installer/pkg/destroy/quota"
	"github.com/openshift/installer/pkg/metrics/timer"

	_ "github.com/openshift/installer/pkg/destroy/aws"
	_ "github.com/openshift/installer/pkg/destroy/azure"
	_ "github.com/openshift/installer/pkg/destroy/baremetal"
	_ "github.com/openshift/installer/pkg/destroy/gcp"
	_ "github.com/openshift/installer/pkg/destroy/ibmcloud"
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
	destroyCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Destroy an OpenShift cluster",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()
			destroyVolumes, err := cmd.Flags().GetBool("destroy-persistent-volumes")
			if err != nil {
				logrus.Fatal(err)
			}
			kubeConfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				logrus.Fatal(err)
			}

			err = runDestroyCmd(command.RootOpts.Dir, os.Getenv("OPENSHIFT_INSTALL_REPORT_QUOTA_FOOTPRINT") == "true", destroyVolumes, kubeConfig)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Infof("Uninstallation complete!")
		},
	}

	destroyCmd.Flags().Bool("destroy-persistent-volumes", false, "cordons and drains nodes, deletes all persistent volumes.")
	destroyCmd.Flags().String("kubeconfig", "", "path to kubeconfig to allow the installer to destroy persistent volumes.")

	destroyCmd.MarkFlagsRequiredTogether("kubeconfig", "destroy-persistent-volumes")
	return destroyCmd
}

func runDestroyCmd(directory string, reportQuota bool, destroyVolumes bool, kubeConfig string) error {
	timer.StartTimer(timer.TotalTimeElapsed)
	destroyer, err := destroy.New(logrus.StandardLogger(), directory, destroyVolumes, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed while preparing to destroy cluster: %w", err)
	}
	quota, err := destroyer.Run()
	if err != nil {
		return fmt.Errorf("failed to destroy cluster: %w", err)
	}

	if reportQuota {
		if err := quotaasset.WriteQuota(directory, quota); err != nil {
			return fmt.Errorf("failed to record quota: %w", err)
		}
	}

	store, err := assetstore.NewStore(directory)
	if err != nil {
		return fmt.Errorf("failed to create asset store: %w", err)
	}
	for _, asset := range clusterTarget.assets {
		if err := store.Destroy(asset); err != nil {
			return fmt.Errorf("failed to destroy asset %q: %w", asset.Name(), err)
		}
	}

	// delete the state file as well
	err = store.DestroyState()
	if err != nil {
		return fmt.Errorf("failed to remove state file: %w", err)
	}

	// delete terraform files
	tfstateFiles, err := filepath.Glob(filepath.Join(directory, "*.tfstate"))
	if err != nil {
		return fmt.Errorf("failed to glob for tfstate files: %w", err)
	}
	tfvarsFiles, err := filepath.Glob(filepath.Join(directory, "*.tfvars.json"))
	if err != nil {
		return fmt.Errorf("failed to glob for tfvars files: %w", err)
	}
	for _, f := range append(tfstateFiles, tfvarsFiles...) {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("failed to remove terraform file %q: %w", f, err)
		}
	}

	// ensure capi etcd data store and capi artifacts are cleaned up
	capiArtifactsDir := filepath.Join(directory, clusterapi.ArtifactsDir)
	if err := os.RemoveAll(capiArtifactsDir); err != nil {
		logrus.Warnf("failed to remove %s: %v", capiArtifactsDir, err)
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
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			timer.StartTimer(timer.TotalTimeElapsed)
			err := bootstrap.Destroy(context.TODO(), command.RootOpts.Dir)
			if err != nil {
				logrus.Fatal(err)
			}
			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
}
