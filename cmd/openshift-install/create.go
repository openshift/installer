package main

import (
	"context"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/logging"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	destroybootstrap "github.com/openshift/installer/pkg/destroy/bootstrap"
	"github.com/openshift/installer/pkg/gather/service"
	timer "github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types/baremetal"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/route/routeapihelpers"
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
	singleNodeIgnitionConfigTarget = target{
		name: "Single Node Ignition Config",
		command: &cobra.Command{
			Use:   "single-node-ignition-config",
			Short: "Generates the bootstrap-in-place-for-live-iso Ignition Config asset",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
		assets: targetassets.SingleNodeIgnitionConfig,
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

				// FIXME: pulling the kubeconfig and metadata out of the root
				// directory is a bit cludgy when we already have them in memory.
				config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
				if err != nil {
					logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
				}

				timer.StartTimer("Bootstrap Complete")
				if err := waitForBootstrapComplete(ctx, config); err != nil {
					if err2 := logClusterOperatorConditions(ctx, config); err2 != nil {
						logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err2)
					}
					bundlePath, err2 := runGatherBootstrapCmd(rootOpts.dir)
					if err2 != nil {
						logrus.Error("Attempted to gather debug logs after installation failure: ", err2)
					}
					logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
					logrus.Error(err.Error())
					if err2 := service.AnalyzeGatherBundle(bundlePath); err2 != nil {
						logrus.Error("Attempted to analyze the debug logs after installation failure: ", err2)
					}
					logrus.Fatal("Bootstrap failed to complete")
				}
				timer.StopTimer("Bootstrap Complete")
				timer.StartTimer("Bootstrap Destroy")

				if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP"); ok && oi != "" {
					logrus.Warn("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP is set, not destroying bootstrap resources. " +
						"Warning: this should only be used for debugging purposes, and poses a risk to cluster stability.")
				} else {
					logrus.Info("Destroying the bootstrap resources...")
					err = destroybootstrap.Destroy(rootOpts.dir)
					if err != nil {
						logrus.Fatal(err)
					}
				}
				timer.StopTimer("Bootstrap Destroy")

				err = waitForInstallComplete(ctx, config, rootOpts.dir)
				if err != nil {
					if err2 := logClusterOperatorConditions(ctx, config); err2 != nil {
						logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err2)
					}
					logTroubleshootingLink()
					logrus.Fatal(err)
				}
				timer.StopTimer(timer.TotalTimeElapsed)
				timer.LogSummary()
			},
		},
		assets: targetassets.Cluster,
	}

	targets = []target{installConfigTarget, manifestsTarget, ignitionConfigsTarget, clusterTarget, singleNodeIgnitionConfigTarget}
)

// clusterCreateError defines a custom error type that would help identify where the error occurs
// during the bootstrap phase of the installation process. This would help identify whether the error
// comes either from the Kubernetes API failure, the bootstrap failure or a general kubernetes client
// creation error. In the event of any error, this interface packages the error message and a custom
// log message that must be neatly presented to the user before termination of the project.
type clusterCreateError struct {
	wrappedError error
	logMessage   string
}

// Unwrap provides the actual stored error that occured during installation.
func (ce *clusterCreateError) Unwrap() error {
	return ce.wrappedError
}

// Error provides the actual stored error that occured during installation.
func (ce *clusterCreateError) Error() string {
	return ce.logMessage
}

// newAPIError creates a clusterCreateError object with a default error message specific to the API failure.
func newAPIError(errorInfo error) *clusterCreateError {
	return &clusterCreateError{
		wrappedError: errorInfo,
		logMessage: "Failed waiting for Kubernetes API. This error usually happens when there " +
			"is a problem on the bootstrap host that prevents creating a temporary control plane.",
	}
}

// newBootstrapError creates a clusterCreateError object with a default error message specific to the
// bootstrap failure.
func newBootstrapError(errorInfo error) *clusterCreateError {
	return &clusterCreateError{
		wrappedError: errorInfo,
		logMessage: "Failed to wait for bootstrapping to complete. This error usually " +
			"happens when there is a problem with control plane hosts that prevents " +
			"the control plane operators from creating the control plane.",
	}
}

// newClientError creates a clusterCreateError object with a default error message specific to the
// kubernetes client creation failure.
func newClientError(errorInfo error) *clusterCreateError {
	return &clusterCreateError{
		wrappedError: errorInfo,
		logMessage:   "Failed to create a kubernetes client.",
	}
}

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
			logrus.Debugf("bluejay target: %s", a.Name())
			err := assetStore.Fetch(a, targets...)
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
		timer.StartTimer(timer.TotalTimeElapsed)

		cleanup := setupFileHook(rootOpts.dir)
		defer cleanup()

		err := runner(rootOpts.dir)
		if err != nil {
			logrus.Fatal(err)
		}
		if cmd.Name() != "cluster" {
			logrus.Infof(logging.LogCreatedFiles(cmd.Name(), rootOpts.dir, targets))
		}

	}
}

// addRouterCAToClusterCA adds router CA to cluster CA in kubeconfig
func addRouterCAToClusterCA(ctx context.Context, config *rest.Config, directory string) (err error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a Kubernetes client")
	}

	// Configmap may not exist. log and accept not-found errors with configmap.
	caConfigMap, err := client.CoreV1().ConfigMaps("openshift-config-managed").Get(ctx, "default-ingress-cert", metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "fetching default-ingress-cert configmap from openshift-config-managed namespace")
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
			return errors.New("ca-bundle.crt from default-ingress-cert configmap not valid PEM format")
		}

		newCA := append(routerCrtBytes, clusterCABytes...)
		c.CertificateAuthorityData = newCA
	}
	if err := clientcmd.WriteToFile(*kconfig, kubeconfig); err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}
	return nil
}

func waitForBootstrapComplete(ctx context.Context, config *rest.Config) *clusterCreateError {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return newClientError(errors.Wrap(err, "creating a Kubernetes client"))
	}

	discovery := client.Discovery()

	apiTimeout := 20 * time.Minute
	logrus.Infof("Waiting up to %v for the Kubernetes API at %s...", apiTimeout, config.Host)

	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	// Poll quickly so we notice changes, but only log when the response
	// changes (because that's interesting) or when we've seen 15 of the
	// same errors in a row (to show we're still alive).
	logDownsample := 15
	silenceRemaining := logDownsample
	previousErrorSuffix := ""
	timer.StartTimer("API")
	var lastErr error
	wait.Until(func() {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Infof("API %s up", version)
			timer.StopTimer("API")
			cancel()
		} else {
			lastErr = err
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
		if lastErr != nil {
			return newAPIError(lastErr)
		}
		return newAPIError(err)
	}

	return waitForBootstrapConfigMap(ctx, client)
}

// waitForBootstrapConfigMap watches the configmaps in the kube-system namespace
// and waits for the bootstrap configmap to report that bootstrapping has
// completed.
func waitForBootstrapConfigMap(ctx context.Context, client *kubernetes.Clientset) *clusterCreateError {
	timeout := 30 * time.Minute
	logrus.Infof("Waiting up to %v for bootstrapping to complete...", timeout)

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := clientwatch.UntilWithSync(
		waitCtx,
		cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "configmaps", "kube-system", fields.OneTermEqualSelector("metadata.name", "bootstrap")),
		&corev1.ConfigMap{},
		nil,
		func(event watch.Event) (bool, error) {
			switch event.Type {
			case watch.Added, watch.Modified:
			default:
				return false, nil
			}
			cm, ok := event.Object.(*corev1.ConfigMap)
			if !ok {
				logrus.Warnf("Expected a core/v1.ConfigMap object but got a %q object instead", event.Object.GetObjectKind().GroupVersionKind())
				return false, nil
			}
			status, ok := cm.Data["status"]
			if !ok {
				logrus.Debugf("No status found in bootstrap configmap")
				return false, nil
			}
			logrus.Debugf("Bootstrap status: %v", status)
			return status == "complete", nil
		},
	)
	if err != nil {
		return newBootstrapError(err)
	}
	return nil
}

// waitForInitializedCluster watches the ClusterVersion waiting for confirmation
// that the cluster has been initialized.
func waitForInitializedCluster(ctx context.Context, config *rest.Config) error {
	// TODO revert this value back to 30 minutes.  It's currently at the end of 4.6 and we're trying to see if the
	timeout := 40 * time.Minute

	// Wait longer for baremetal, due to length of time it takes to boot
	if assetStore, err := assetstore.NewStore(rootOpts.dir); err == nil {
		if installConfig, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
			if installConfig.(*installconfig.InstallConfig).Config.Platform.Name() == baremetal.Name {
				timeout = 60 * time.Minute
			}
		}
	}

	logrus.Infof("Waiting up to %v for the cluster at %s to initialize...", timeout, config.Host)
	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	clusterVersionContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	failing := configv1.ClusterStatusConditionType("Failing")
	timer.StartTimer("Cluster Operators")
	var lastError string
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
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorAvailable) {
					timer.StopTimer("Cluster Operators")
					return true, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, failing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, failing).Message
				} else if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorProgressing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, configv1.OperatorProgressing).Message
				}
				logrus.Debugf("Still waiting for the cluster to initialize: %s", lastError)
				return false, nil
			}
			logrus.Debug("Still waiting for the cluster to initialize...")
			return false, nil
		},
	)

	if err == nil {
		logrus.Debug("Cluster is initialized")
		return nil
	}

	if lastError != "" {
		if err == wait.ErrWaitTimeout {
			return errors.Errorf("failed to initialize the cluster: %s", lastError)
		}

		return errors.Wrapf(err, "failed to initialize the cluster: %s", lastError)
	}

	return errors.Wrap(err, "failed to initialize the cluster")
}

// waitForConsole returns the console URL from the route 'console' in namespace openshift-console
func waitForConsole(ctx context.Context, config *rest.Config) (string, error) {
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
	timer.StartTimer("Console")
	wait.Until(func() {
		route, err := rc.RouteV1().Routes(consoleNamespace).Get(ctx, consoleRouteName, metav1.GetOptions{})
		if err == nil {
			logrus.Debugf("Route found in openshift-console namespace: %s", consoleRouteName)
			if uri, _, err2 := routeapihelpers.IngressURI(route, ""); err2 == nil {
				url = uri.String()
				logrus.Debug("OpenShift console route is admitted")
				cancel()
			} else {
				err = err2
			}
		}
		if err != nil {
			silenceRemaining--
			if silenceRemaining == 0 {
				logrus.Debugf("Still waiting for the console route: %v", err)
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
	timer.StopTimer("Console")
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
	logrus.Infof("To access the cluster as the system:admin user when using 'oc', run 'export KUBECONFIG=%s'", kubeconfig)
	logrus.Infof("Access the OpenShift web-console here: %s", consoleURL)
	logrus.Infof("Login to the console with user: %q, and password: %q", "kubeadmin", pw)
	return nil
}

func waitForInstallComplete(ctx context.Context, config *rest.Config, directory string) error {
	if err := waitForInitializedCluster(ctx, config); err != nil {
		return err
	}

	consoleURL, err := waitForConsole(ctx, config)
	if err != nil {
		return err
	}

	if err = addRouterCAToClusterCA(ctx, config, rootOpts.dir); err != nil {
		return err
	}

	return logComplete(rootOpts.dir, consoleURL)
}

func logTroubleshootingLink() {
	logrus.Error(`Cluster initialization failed because one or more operators are not functioning properly.
The cluster should be accessible for troubleshooting as detailed in the documentation linked below,
https://docs.openshift.com/container-platform/latest/support/troubleshooting/troubleshooting-installations.html
The 'wait-for install-complete' subcommand can then be used to continue the installation`)
}
