package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/templates"
	destroybootstrap "github.com/openshift/installer/pkg/destroy/bootstrap"
)

type target struct {
	name    string
	command *cobra.Command
	assets  []asset.WritableAsset
}

// each target is a variable to preserve the order when creating subcommands and still
// allow other functions to directly access each target individually.
var (
	installConfigTarget = target{
		name: "Install Config",
		command: &cobra.Command{
			Use:   "install-config",
			Short: "Generates the Install Config asset",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: []asset.WritableAsset{&installconfig.InstallConfig{}},
	}

	manifestsTarget = target{
		name: "Manifests",
		command: &cobra.Command{
			Use:   "manifests",
			Short: "Generates the Kubernetes manifests",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: []asset.WritableAsset{&manifests.Manifests{}, &manifests.Tectonic{}},
	}

	manifestTemplatesTarget = target{
		name: "Manifest templates",
		command: &cobra.Command{
			Use:   "manifest-templates",
			Short: "Generates the unrendered Kubernetes manifest templates",
			Long:  "",
		},
		assets: []asset.WritableAsset{&templates.Templates{}},
	}

	ignitionConfigsTarget = target{
		name: "Ignition Configs",
		command: &cobra.Command{
			Use:   "ignition-configs",
			Short: "Generates the Ignition Config asset",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: []asset.WritableAsset{&bootstrap.Bootstrap{}, &machine.Master{}, &machine.Worker{}},
	}

	clusterTarget = target{
		name: "Cluster",
		command: &cobra.Command{
			Use:   "cluster",
			Short: "Create an OpenShift cluster",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
			PostRunE: func(_ *cobra.Command, _ []string) error {
				err := destroyBootstrap(context.Background(), rootOpts.dir)
				if err != nil {
					return err
				}
				return logComplete(rootOpts.dir)
			},
		},
		assets: []asset.WritableAsset{&cluster.TerraformVariables{}, &kubeconfig.Admin{}, &cluster.Cluster{}},
	}

	targets = []target{installConfigTarget, manifestTemplatesTarget, manifestsTarget, ignitionConfigsTarget, clusterTarget}
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create part of an OpenShift cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range targets {
		t.command.RunE = runTargetCmd(t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}

func runTargetCmd(targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cleanup, err := setupFileHook(rootOpts.dir)
		if err != nil {
			return errors.Wrap(err, "failed to setup logging hook")
		}
		defer cleanup()

		assetStore, err := asset.NewStore(rootOpts.dir)
		if err != nil {
			return errors.Wrapf(err, "failed to create asset store")
		}

		for _, a := range targets {
			err := assetStore.Fetch(a)
			if err != nil {
				if exitError, ok := errors.Cause(err).(*exec.ExitError); ok && len(exitError.Stderr) > 0 {
					logrus.Error(strings.Trim(string(exitError.Stderr), "\n"))
				}
				err = errors.Wrapf(err, "failed to fetch %s", a.Name())
			}

			if err2 := asset.PersistToFile(a, rootOpts.dir); err2 != nil {
				err2 = errors.Wrapf(err2, "failed to write asset (%s) to disk", a.Name())
				if err != nil {
					logrus.Error(err2)
					return err
				}
				return err2
			}

			if err != nil {
				return err
			}
		}
		return nil
	}
}

// FIXME: pulling the kubeconfig and metadata out of the root
// directory is a bit cludgy when we already have them in memory.
func destroyBootstrap(ctx context.Context, directory string) (err error) {
	cleanup, err := setupFileHook(rootOpts.dir)
	if err != nil {
		return errors.Wrap(err, "failed to setup logging hook")
	}
	defer cleanup()

	logrus.Info("Waiting for bootstrap completion...")
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(directory, "auth", "kubeconfig"))
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

	events := client.CoreV1().Events("kube-system")

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
