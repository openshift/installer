package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-log/log/print"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer/pkg/assets"
	destroybootstrap "github.com/openshift/installer/pkg/destroy/bootstrap"
	"github.com/openshift/installer/pkg/installerassets"
	_ "github.com/openshift/installer/pkg/installerassets/aws"
	_ "github.com/openshift/installer/pkg/installerassets/libvirt"
	_ "github.com/openshift/installer/pkg/installerassets/openstack"
	_ "github.com/openshift/installer/pkg/installerassets/tls"
	"github.com/openshift/installer/pkg/terraform"
)

var (
	createAssetsOpts struct {
		prune bool
	}
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create part of an OpenShift cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	assets := &cobra.Command{
		Use:   "assets",
		Short: "Generates installer assets",
		Long:  "Generates installer assets.  Can be run multiple times on the same directory to propagate changes made to any asset through the Merkle tree.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			cleanup, err := setupFileHook(rootOpts.dir)
			if err != nil {
				return errors.Wrap(err, "failed to setup logging hook")
			}
			defer cleanup()

			_, err = syncAssets(ctx, rootOpts.dir, createAssetsOpts.prune)
			return err
		},
	}
	assets.PersistentFlags().BoolVar(&createAssetsOpts.prune, "prune", false, "remove everything except referenced assets from the asset directory")
	cmd.AddCommand(assets)

	cmd.AddCommand(&cobra.Command{
		Use:   "cluster",
		Short: "Creates the cluster",
		Long:  "Generates resources based on the installer assets, launching the cluster.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := context.Background()

			cleanup, err := setupFileHook(rootOpts.dir)
			if err != nil {
				return errors.Wrap(err, "failed to setup logging hook")
			}
			defer cleanup()

			assets, err := syncAssets(ctx, rootOpts.dir, createAssetsOpts.prune)
			if err != nil {
				return err
			}

			err = createCluster(ctx, assets, rootOpts.dir)
			if err != nil {
				return err
			}

			err = destroyBootstrap(ctx, assets)
			if err != nil {
				return err
			}

			return logComplete(rootOpts.dir)
		},
	})

	return cmd
}

func syncAssets(ctx context.Context, directory string, prune bool) (*assets.Assets, error) {
	assets := installerassets.New()
	err := assets.Read(ctx, directory, installerassets.GetDefault, print.New(logrus.StandardLogger()))
	if err != nil {
		return nil, err
	}

	err = assets.Write(ctx, directory, prune)
	return assets, err
}

func createCluster(ctx context.Context, assets *assets.Assets, directory string) error {
	tmpDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	platformAsset, err := assets.GetByName(ctx, "platform")
	if err != nil {
		return errors.Wrapf(err, `retrieve "platform" by name`)
	}
	platform := string(platformAsset.Data)

	for _, filename := range []string{"terraform.tfvars", fmt.Sprintf("%s-terraform.auto.tfvars", platform)} {
		assetName := path.Join("terraform", filename)
		tfVars, err := assets.GetByName(ctx, assetName)
		if err != nil {
			return errors.Wrapf(err, "retrieve %q by name", assetName)
		}

		if err := ioutil.WriteFile(filepath.Join(tmpDir, filename), tfVars.Data, 0600); err != nil {
			return err
		}
	}
	logrus.Info("Using Terraform to create cluster...")
	stateFile, err := terraform.Apply(tmpDir, platform)
	if err != nil {
		err = errors.Wrap(err, "run Terraform")
	}

	if stateFile != "" {
		data, err2 := ioutil.ReadFile(stateFile)
		if err2 == nil {
			err2 = ioutil.WriteFile(filepath.Join(directory, "terraform", "terraform.tfstate"), data, 0666)
		}
		if err == nil {
			err = err2
		} else {
			logrus.Error(errors.Wrap(err2, "read Terraform state"))
		}
	}

	return err
}

func destroyBootstrap(ctx context.Context, assets *assets.Assets) error {
	logrus.Info("Waiting for bootstrap completion...")

	kubeconfig, err := assets.GetByName(ctx, "auth/kubeconfig-admin")
	if err != nil {
		return errors.Wrap(err, `retrieve "auth/kubeconfig-admin" by name`)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig.Data)
	if err != nil {
		return errors.Wrap(err, "loading kubeconfig")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a Kubernetes client")
	}

	discovery := client.Discovery()

	apiContext, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	wait.Until(func() {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Infof("API %s up", version)
			cancel()
		} else {
			logrus.Debugf("API not up yet: %s", err)
		}
	}, 2*time.Second, apiContext.Done())

	events := client.CoreV1().Events(metav1.NamespaceSystem)

	eventContext, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	_, err = Until(
		eventContext,
		"",
		func(sinceResourceVersion string) (watch.Interface, error) {
			for {
				watcher, err := events.Watch(metav1.ListOptions{
					ResourceVersion: sinceResourceVersion,
				})
				if err == nil {
					return watcher, nil
				}
				select {
				case <-eventContext.Done():
					return watcher, err
				default:
					logrus.Warningf("Failed to connect events watcher: %s", err)
					time.Sleep(2 * time.Second)
				}
			}
		},
		func(watchEvent watch.Event) (bool, error) {
			event, ok := watchEvent.Object.(*corev1.Event)
			if !ok {
				return false, nil
			}

			if watchEvent.Type == watch.Error {
				logrus.Debugf("error %s: %s", event.Name, event.Message)
				return false, nil
			}

			if watchEvent.Type != watch.Added {
				return false, nil
			}

			logrus.Debugf("added %s: %s", event.Name, event.Message)
			return event.Name == "bootstrap-complete", nil
		},
	)
	if err != nil {
		return errors.Wrap(err, "waiting for bootstrap-complete")
	}

	logrus.Info("Destroying the bootstrap resources...")
	// FIXME: pulling the metadata out of the root directory is a bit
	// cludgy when we already have it in memory.
	return destroybootstrap.Destroy(rootOpts.dir)
}

// logComplete prints info upon completion
func logComplete(directory string) error {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	kubeconfig := filepath.Join(absDir, "auth", "kubeconfig")
	logrus.Infof("Install complete! Run 'export KUBECONFIG=%s' to manage your cluster.", kubeconfig)
	logrus.Info("After exporting your kubeconfig, run 'oc -h' for a list of OpenShift client commands.")
	return nil
}
