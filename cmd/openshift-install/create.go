package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	destroybootstrap "github.com/openshift/installer/pkg/destroy/bootstrap"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
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
		assets: targetassets.InstallConfig,
	}

	manifestsTarget = target{
		name: "Manifests",
		command: &cobra.Command{
			Use:   "manifests",
			Short: "Generates the Kubernetes manifests",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: targetassets.Manifests,
	}

	manifestTemplatesTarget = target{
		name: "Manifest templates",
		command: &cobra.Command{
			Use:   "manifest-templates",
			Short: "Generates the unrendered Kubernetes manifest templates",
			Long:  "",
		},
		assets: targetassets.ManifestTemplates,
	}

	ignitionConfigsTarget = target{
		name: "Ignition Configs",
		command: &cobra.Command{
			Use:   "ignition-configs",
			Short: "Generates the Ignition Config asset",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: targetassets.IgnitionConfigs,
	}

	clusterTarget = target{
		name: "Cluster",
		command: &cobra.Command{
			Use:   "cluster",
			Short: "Create an OpenShift cluster",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
			PostRun: func(_ *cobra.Command, _ []string) {
				ctx := context.Background()

				cleanup := setupFileHook(rootOpts.dir)
				defer cleanup()

				config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
				if err != nil {
					logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
				}

				err = destroyBootstrap(ctx, config, rootOpts.dir)
				if err != nil {
					logrus.Fatal(err)
				}

				if err := waitForInitializedCluster(ctx, config); err != nil {
					logrus.Fatal(err)
				}

				consoleURL, err := waitForConsole(ctx, config, rootOpts.dir)
				if err != nil {
					logrus.Fatal(err)
				}

				if err = addRouterCAToClusterCA(config, rootOpts.dir); err != nil {
					logrus.Fatal(err)
				}

				err = logComplete(rootOpts.dir, consoleURL)
				if err != nil {
					logrus.Fatal(err)
				}
			},
		},
		assets: targetassets.Cluster,
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
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}

func runTargetCmd(targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) {
	runner := func(directory string) error {
		assetStore, err := assetstore.NewStore(directory)
		if err != nil {
			return errors.Wrap(err, "failed to create asset store")
		}

		for _, a := range targets {
			err := assetStore.Fetch(a)
			if err != nil {
				err = errors.Wrapf(err, "failed to fetch %s", a.Name())
			}

			if err2 := asset.PersistToFile(a, directory); err2 != nil {
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

	return func(cmd *cobra.Command, args []string) {
		cleanup := setupFileHook(rootOpts.dir)
		defer cleanup()

		err := runner(rootOpts.dir)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

// addRouterCAToClusterCA adds router CA to cluster CA in kubeconfig
func addRouterCAToClusterCA(config *rest.Config, directory string) (err error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a Kubernetes client")
	}

	// Configmap may not exist. log and accept not-found errors with configmap.
	caConfigMap, err := client.CoreV1().ConfigMaps("openshift-config-managed").Get("router-ca", metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logrus.Infof("router-ca resource not found in cluster, perhaps you are not using default router CA")
			return nil
		}
		return errors.Wrap(err, "fetching router-ca configmap from openshift-config-managed namespace")
	}

	routerCrtBytes := []byte(caConfigMap.Data["ca-bundle.crt"])
	kubeconfig := filepath.Join(directory, "auth", "kubeconfig")
	kconfig, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "loading kubeconfig")
	}

	if kconfig == nil || len(kconfig.Clusters) == 0 {
		return errors.New("kubeconfig is missing expected data")
	}

	for _, c := range kconfig.Clusters {
		clusterCABytes := c.CertificateAuthorityData
		if len(clusterCABytes) == 0 {
			return errors.New("kubeconfig CertificateAuthorityData not found")
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(clusterCABytes) {
			return errors.New("cluster CA found in kubeconfig not valid PEM format")
		}
		if !certPool.AppendCertsFromPEM(routerCrtBytes) {
			return errors.New("ca-bundle.crt from router-ca configmap not valid PEM format")
		}

		newCA := append(routerCrtBytes, clusterCABytes...)
		c.CertificateAuthorityData = newCA
	}
	if err := clientcmd.WriteToFile(*kconfig, kubeconfig); err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}
	return nil
}

// FIXME: pulling the kubeconfig and metadata out of the root
// directory is a bit cludgy when we already have them in memory.
func destroyBootstrap(ctx context.Context, config *rest.Config, directory string) (err error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a Kubernetes client")
	}

	discovery := client.Discovery()

	apiTimeout := 30 * time.Minute
	logrus.Infof("Waiting up to %v for the Kubernetes API...", apiTimeout)
	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	// Poll quickly so we notice changes, but only log when the response
	// changes (because that's interesting) or when we've seen 15 of the
	// same errors in a row (to show we're still alive).
	logDownsample := 15
	silenceRemaining := logDownsample
	previousErrorSuffix := ""
	wait.Until(func() {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Infof("API %s up", version)
			cancel()
		} else {
			silenceRemaining--
			chunks := strings.Split(err.Error(), ":")
			errorSuffix := chunks[len(chunks)-1]
			if previousErrorSuffix != errorSuffix {
				logrus.Debugf("Still waiting for the Kubernetes API: %v", err)
				previousErrorSuffix = errorSuffix
				silenceRemaining = logDownsample
			} else if silenceRemaining == 0 {
				logrus.Debugf("Still waiting for the Kubernetes API: %v", err)
				silenceRemaining = logDownsample
			}
		}
	}, 2*time.Second, apiContext.Done())
	err = apiContext.Err()
	if err != nil && err != context.Canceled {
		return errors.Wrap(err, "waiting for Kubernetes API")
	}

	events := client.CoreV1().Events("kube-system")

	eventTimeout := 30 * time.Minute
	logrus.Infof("Waiting up to %v for the bootstrap-complete event...", eventTimeout)
	if err := waitForEvent(ctx, events, "bootstrap-complete", eventTimeout); err != nil {
		return err
	}

	logrus.Info("Destroying the bootstrap resources...")
	return destroybootstrap.Destroy(rootOpts.dir)
}

// waitForEvent watches the events, waits for the event of the given name and prints out all other events on the way.
func waitForEvent(ctx context.Context, client corev1client.EventInterface, name string, timeout time.Duration) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := clientwatch.UntilWithSync(
		waitCtx,
		&cache.ListWatch{ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			ret, err := client.List(options)
			if err != nil {
				logrus.Debugf("Failed to list events: %s", err)
			}
			return ret, err
		}, WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			ret, err := client.Watch(options)
			if err != nil {
				logrus.Debugf("Failed to watch events: %s", err)
			}
			return ret, err
		}},
		&corev1.Event{},
		nil,
		func(event watch.Event) (bool, error) {
			ev, ok := event.Object.(*corev1.Event)
			if !ok {
				logrus.Warnf("Expected an core/v1.Event object but got a %q object instead", event.Object.GetObjectKind().GroupVersionKind())
				return false, nil
			}

			switch event.Type {
			case watch.Added, watch.Modified:
				logrus.Debugf("added %s: %s", ev.Name, ev.Message)
				return ev.Name == name, nil
			case watch.Error:
				logrus.Debugf("error %s: %s", ev.Name, ev.Message)
			}
			return false, nil
		},
	)

	return errors.Wrapf(err, "failed to wait for %s event", name)
}

// waitForInitializedCluster watches the ClusterVersion waiting for confirmation
// that the cluster has been initialized.
func waitForInitializedCluster(ctx context.Context, config *rest.Config) error {
	timeout := 30 * time.Minute
	logrus.Infof("Waiting up to %v for the cluster to initialize...", timeout)
	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	clusterVersionContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var clusterVersion *configv1.ClusterVersion

	_, err = clientwatch.UntilWithSync(
		clusterVersionContext,
		cache.NewListWatchFromClient(cc.ConfigV1().RESTClient(), "clusterversions", "", fields.OneTermEqualSelector("metadata.name", "version")),
		&configv1.ClusterVersion{},
		nil,
		func(event watch.Event) (bool, error) {
			switch event.Type {
			case watch.Added, watch.Modified:
				cv, ok := event.Object.(*configv1.ClusterVersion)
				if !ok {
					logrus.Warnf("Expected a ClusterVersion object but got a %q object instead", event.Object.GetObjectKind().GroupVersionKind())
					return false, nil
				}
				clusterVersion = cv
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorAvailable) {
					logrus.Debug("Cluster is initialized")
					return true, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorFailing) {
					logrus.Debugf("Still waiting for the cluster to initialize: %v",
						cov1helpers.FindStatusCondition(cv.Status.Conditions, configv1.OperatorFailing).Message)
					return false, nil
				}
			}
			logrus.Debug("Still waiting for the cluster to initialize...")
			return false, nil
		},
	)

	// If we timed out and the CVO failed, print out the failure message
	if err != nil && clusterVersion != nil {
		if cov1helpers.IsStatusConditionTrue(clusterVersion.Status.Conditions, configv1.OperatorFailing) {
			err = errors.New(cov1helpers.FindStatusCondition(clusterVersion.Status.Conditions, configv1.OperatorFailing).Message)
		}
	}

	return errors.Wrap(err, "failed to initialize the cluster")
}

// waitForConsole returns the console URL from the route 'console' in namespace openshift-console
func waitForConsole(ctx context.Context, config *rest.Config, directory string) (string, error) {
	url := ""
	// Need to keep these updated if they change
	consoleNamespace := "openshift-console"
	consoleRouteName := "console"
	rc, err := routeclient.NewForConfig(config)
	if err != nil {
		return "", errors.Wrap(err, "creating a route client")
	}

	consoleRouteTimeout := 10 * time.Minute
	logrus.Infof("Waiting up to %v for the openshift-console route to be created...", consoleRouteTimeout)
	consoleRouteContext, cancel := context.WithTimeout(ctx, consoleRouteTimeout)
	defer cancel()
	// Poll quickly but only log when the response
	// when we've seen 15 of the same errors or output of
	// no route in a row (to show we're still alive).
	logDownsample := 15
	silenceRemaining := logDownsample
	wait.Until(func() {
		consoleRoutes, err := rc.RouteV1().Routes(consoleNamespace).List(metav1.ListOptions{})
		if err == nil && len(consoleRoutes.Items) > 0 {
			for _, route := range consoleRoutes.Items {
				logrus.Debugf("Route found in openshift-console namespace: %s", route.Name)
				if route.Name == consoleRouteName {
					url = fmt.Sprintf("https://%s", route.Spec.Host)
				}
			}
			logrus.Debug("OpenShift console route is created")
			cancel()
		} else if err != nil {
			silenceRemaining--
			if silenceRemaining == 0 {
				logrus.Debugf("Still waiting for the console route: %v", err)
				silenceRemaining = logDownsample
			}
		} else if len(consoleRoutes.Items) == 0 {
			silenceRemaining--
			if silenceRemaining == 0 {
				logrus.Debug("Still waiting for the console route...")
				silenceRemaining = logDownsample
			}
		}
	}, 2*time.Second, consoleRouteContext.Done())
	err = consoleRouteContext.Err()
	if err != nil && err != context.Canceled {
		return url, errors.Wrap(err, "waiting for openshift-console URL")
	}
	if url == "" {
		return url, errors.New("could not get openshift-console URL")
	}
	return url, nil
}

// logComplete prints info upon completion
func logComplete(directory, consoleURL string) error {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	kubeconfig := filepath.Join(absDir, "auth", "kubeconfig")
	pwFile := filepath.Join(absDir, "auth", "kubeadmin-password")
	pw, err := ioutil.ReadFile(pwFile)
	if err != nil {
		return err
	}
	logrus.Info("Install complete!")
	logrus.Infof("Run 'export KUBECONFIG=%s' to manage the cluster with 'oc', the OpenShift CLI.", kubeconfig)
	logrus.Infof("The cluster is ready when 'oc login -u kubeadmin -p %s' succeeds (wait a few minutes).", pw)
	logrus.Infof("Access the OpenShift web-console here: %s", consoleURL)
	logrus.Infof("Login to the console with user: kubeadmin, password: %s", pw)
	return nil
}
