package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"net"
	"os"
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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	configlisters "github.com/openshift/client-go/config/listers/config/v1"
	operatorclient "github.com/openshift/client-go/operator/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/lbconfig"
	"github.com/openshift/installer/pkg/asset/logging"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	destroybootstrap "github.com/openshift/installer/pkg/destroy/bootstrap"
	timer "github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/vsphere"
	baremetalutils "github.com/openshift/installer/pkg/utils/baremetal"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/route/routeapihelpers"
)

type target struct {
	name    string
	command *cobra.Command
	assets  []asset.WritableAsset
}

const (
	exitCodeInstallConfigError = iota + 3
	exitCodeInfrastructureFailed
	exitCodeBootstrapFailed
	exitCodeInstallFailed
	exitCodeOperatorStabilityFailed
	exitCodeInterrupt

	// coStabilityThreshold is how long a cluster operator must have Progressing=False
	// in order to be considered stable. Measured in seconds.
	coStabilityThreshold float64 = 30
)

var skipPasswordPrintFlag bool

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
			PostRun: func(cmd *cobra.Command, _ []string) {
				// Get the context that was set in newCreateCmd.
				ctx := cmd.Context()

				exitCode, err := clusterCreatePostRun(ctx)
				if err != nil {
					logrus.Fatal(err)
				}
				if exitCode != 0 {
					logrus.Exit(exitCode)
				}
			},
		},
		assets: targetassets.Cluster,
	}

	permissionsTarget = target{
		name: "Permissions",
		command: &cobra.Command{
			Use:   "permissions-policy",
			Short: "Generates a list of required permissions asset",
			// This is internal-only for now
			Hidden: true,
		},
		assets: targetassets.Permissions,
	}

	targets = []target{installConfigTarget, manifestsTarget, ignitionConfigsTarget, clusterTarget, singleNodeIgnitionConfigTarget, permissionsTarget}
)

// clusterCreatePostRun is the main entrypoint for the cluster create command
// it was moved out of the clusterTarget.command.PostRun function to allow cleanup operations to always
// run in a defer statement, given that we had multiple exit points in the function, like logrus.Fatal or logrus.Exit.
//
// Currently this function returns an exit code and an error, we should refactor this to only return an error,
// that can be wrapped if we want a custom exit code.
func clusterCreatePostRun(ctx context.Context) (int, error) {
	cleanup := command.SetupFileHook(command.RootOpts.Dir)
	defer cleanup()

	// FIXME: pulling the kubeconfig and metadata out of the root
	// directory is a bit cludgy when we already have them in memory.
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(command.RootOpts.Dir, "auth", "kubeconfig"))
	if err != nil {
		return 0, errors.Wrap(err, "loading kubeconfig")
	}

	// Handle the case when the API server is not reachable.
	if err := handleUnreachableAPIServer(ctx, config); err != nil {
		logrus.Fatal(fmt.Errorf("unable to handle api server override: %w", err))
	}

	//
	// Wait for the bootstrap to complete.
	//
	timer.StartTimer("Bootstrap Complete")
	if err := waitForBootstrapComplete(ctx, config); err != nil {
		if err := logClusterOperatorConditions(ctx, config); err != nil {
			logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err)
		}
		logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
		logrus.Error(err.Error())
		gatherAndAnalyzeBootstrapLogs(ctx, command.RootOpts.Dir)
		return exitCodeBootstrapFailed, nil
	}
	timer.StopTimer("Bootstrap Complete")

	//
	// Wait for the bootstrap to be destroyed.
	//
	timer.StartTimer("Bootstrap Destroy")
	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP"); ok && oi != "" {
		logrus.Warn("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP is set, not destroying bootstrap resources. " +
			"Warning: this should only be used for debugging purposes, and poses a risk to cluster stability.")
	} else {
		logrus.Info("Destroying the bootstrap resources...")
		err = destroybootstrap.Destroy(ctx, command.RootOpts.Dir)
		if err != nil {
			return 0, err
		}
	}
	timer.StopTimer("Bootstrap Destroy")

	//
	// Wait for the cluster to initialize.
	//
	err = waitForInstallComplete(ctx, config, command.RootOpts.Dir)
	if err != nil {
		if err2 := logClusterOperatorConditions(ctx, config); err2 != nil {
			logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err2)
		}
		logTroubleshootingLink()
		logrus.Error(err)
		return exitCodeInstallFailed, nil
	}
	timer.StopTimer(timer.TotalTimeElapsed)
	timer.LogSummary()
	return 0, nil
}

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

func newCreateCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create part of an OpenShift cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range targets {
		t.command.Args = cobra.ExactArgs(0)
		t.command.Run = runTargetCmd(ctx, t.assets...)
		if t.name == "Cluster" {
			t.command.PersistentFlags().BoolVar(&skipPasswordPrintFlag, "skip-password-print", false, "Do not print the generated user password.")
		}
		cmd.AddCommand(t.command)
	}

	return cmd
}

func runTargetCmd(ctx context.Context, targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) {
	runner := func(directory string) error {
		fetcher := assetstore.NewAssetsFetcher(directory)
		return fetcher.FetchAndPersist(ctx, targets)
	}

	return func(cmd *cobra.Command, args []string) {
		timer.StartTimer(timer.TotalTimeElapsed)

		// Set the context to be used in the PostRun function.
		cmd.SetContext(ctx)

		cleanup := command.SetupFileHook(command.RootOpts.Dir)
		defer cleanup()

		cluster.InstallDir = command.RootOpts.Dir

		err := runner(command.RootOpts.Dir)
		if err != nil {
			if strings.Contains(err.Error(), asset.InstallConfigError) {
				logrus.Error(err)
				logrus.Exit(exitCodeInstallConfigError)
			}
			if strings.Contains(err.Error(), asset.ControlPlaneCreationError) {
				gatherAndAnalyzeBootstrapLogs(ctx, command.RootOpts.Dir)
			}
			if strings.Contains(err.Error(), asset.ClusterCreationError) {
				logrus.Error(err)
				logrus.Exit(exitCodeInfrastructureFailed)
			}
			logrus.Fatal(err)
		}
		switch cmd.Name() {
		case "cluster", "image", "pxe-files":
		default:
			logrus.Infof(logging.LogCreatedFiles(cmd.Name(), command.RootOpts.Dir, targets))
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

	untilTime := time.Now().Add(apiTimeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for the Kubernetes API at %s...",
		apiTimeout, untilTime.Format(time.Kitchen), timezone, config.Host)

	apiContext, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()
	// Poll quickly so we notice changes, but only log when the response
	// changes (because that's interesting) or when we've seen 15 of the
	// same errors in a row (to show we're still alive).
	logDownsample := 15
	silenceRemaining := logDownsample
	previousErrorSuffix := ""
	timer.StartTimer("API")

	if assetStore, err := assetstore.NewStore(command.RootOpts.Dir); err == nil {
		checkIfAgentCommand(assetStore)
	}

	var lastErr error
	err = wait.PollUntilContextCancel(apiContext, 2*time.Second, true, func(_ context.Context) (done bool, err error) {
		version, err := discovery.ServerVersion()
		if err == nil {
			logrus.Infof("API %s up", version)
			timer.StopTimer("API")
			return true, nil
		}

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

		return false, nil
	})
	if err != nil {
		if lastErr != nil {
			return newAPIError(lastErr)
		}
		return newAPIError(err)
	}

	var platformName string

	if assetStore, err := assetstore.NewStore(command.RootOpts.Dir); err == nil {
		if installConfig, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
			platformName = installConfig.(*installconfig.InstallConfig).Config.Platform.Name()
		}
	}

	timeout := 45 * time.Minute

	// Wait longer for baremetal, VSphere due to length of time it takes to boot
	if platformName == baremetal.Name || platformName == vsphere.Name {
		timeout = 60 * time.Minute
	}

	untilTime = time.Now().Add(timeout)
	timezone, _ = untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for bootstrapping to complete...",
		timeout, untilTime.Format(time.Kitchen), timezone)

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if platformName == baremetal.Name {
		if err := baremetalutils.WaitForBaremetalBootstrapControlPlane(waitCtx, config, command.RootOpts.Dir); err != nil {
			return newBootstrapError(err)
		}
		logrus.Infof("  Baremetal control plane finished provisioning.")
	}

	if err := waitForBootstrapConfigMap(waitCtx, client); err != nil {
		return err
	}

	if err := waitForEtcdBootstrapMemberRemoval(ctx, config); err != nil {
		return newBootstrapError(err)
	}

	return nil
}

// waitForBootstrapConfigMap watches the configmaps in the kube-system namespace
// and waits for the bootstrap configmap to report that bootstrapping has
// completed.
func waitForBootstrapConfigMap(ctx context.Context, client *kubernetes.Clientset) *clusterCreateError {
	_, err := clientwatch.UntilWithSync(
		ctx,
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

// If the bootstrap etcd member is cleaned up before it has been removed from the etcd cluster, the
// etcd cluster cannot maintain quorum through the rollout of any single permanent member.
func waitForEtcdBootstrapMemberRemoval(ctx context.Context, config *rest.Config) error {
	timeout := 5 * time.Minute

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for the bootstrap etcd member to be removed...",
		timeout, untilTime.Format(time.Kitchen), timezone)

	client, err := operatorclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating operator client: %w", err)
	}
	// Validate the etcd operator has removed the bootstrap etcd member
	return wait.PollUntilContextTimeout(ctx, 1*time.Second, timeout, true, func(ctx context.Context) (done bool, err error) {
		etcd, err := client.OperatorV1().Etcds().Get(ctx, "cluster", metav1.GetOptions{})
		if err != nil {
			logrus.Debugf("Error getting etcd operator singleton, retrying: %v", err)
			return false, nil
		}

		for _, condition := range etcd.Status.Conditions {
			if condition.Type == "EtcdBootstrapMemberRemoved" {
				return condition.Status == operatorv1.ConditionTrue, nil
			}
		}
		return false, nil
	})
}

// waitForInitializedCluster watches the ClusterVersion waiting for confirmation
// that the cluster has been initialized.
func waitForInitializedCluster(ctx context.Context, config *rest.Config) error {
	// TODO revert this value back to 30 minutes.  It's currently at the end of 4.6 and we're trying to see if the
	timeout := 40 * time.Minute

	// Wait longer for baremetal, due to length of time it takes to boot
	if assetStore, err := assetstore.NewStore(command.RootOpts.Dir); err == nil {
		if installConfig, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
			if installConfig.(*installconfig.InstallConfig).Config.Platform.Name() == baremetal.Name {
				timeout = 60 * time.Minute
			}
		}

		checkIfAgentCommand(assetStore)
	}

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for the cluster at %s to initialize...",
		timeout, untilTime.Format(time.Kitchen), timezone, config.Host)
	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	clusterVersionContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	failing := configv1.ClusterStatusConditionType("Failing")
	timer.StartTimer("Cluster Operators Available")
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
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorAvailable) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, failing) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, configv1.OperatorProgressing) {
					timer.StopTimer("Cluster Operators Available")
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

// waitForStableOperators ensures that each cluster operator is "stable", i.e. the
// operator has not been in a progressing state for at least a certain duration,
// 30 seconds by default. Returns an error if any operator does meet this threshold
// after a deadline, 30 minutes by default.
func waitForStableOperators(ctx context.Context, config *rest.Config) error {
	timer.StartTimer("Cluster Operators Stable")

	stabilityCheckDuration := 30 * time.Minute
	stabilityContext, cancel := context.WithTimeout(ctx, stabilityCheckDuration)
	defer cancel()

	untilTime := time.Now().Add(stabilityCheckDuration)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) to ensure each cluster operator has finished progressing...",
		stabilityCheckDuration, untilTime.Format(time.Kitchen), timezone)

	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	configInformers := configinformers.NewSharedInformerFactory(cc, 0)
	clusterOperatorInformer := configInformers.Config().V1().ClusterOperators().Informer()
	clusterOperatorLister := configInformers.Config().V1().ClusterOperators().Lister()
	configInformers.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), clusterOperatorInformer.HasSynced) {
		return fmt.Errorf("informers never started")
	}

	waitErr := wait.PollUntilContextCancel(stabilityContext, 1*time.Second, true, waitForAllClusterOperators(clusterOperatorLister))
	if waitErr != nil {
		logrus.Errorf("Error checking cluster operator Progressing status: %q", waitErr)
		stableOperators, unstableOperators, err := currentOperatorStability(clusterOperatorLister)
		if err != nil {
			logrus.Errorf("Error checking final cluster operator Progressing status: %q", err)
		}
		logrus.Debugf("These cluster operators were stable: [%s]", strings.Join(sets.List(stableOperators), ", "))
		logrus.Errorf("These cluster operators were not stable: [%s]", strings.Join(sets.List(unstableOperators), ", "))

		logrus.Exit(exitCodeOperatorStabilityFailed)
	}

	timer.StopTimer("Cluster Operators Stable")

	logrus.Info("All cluster operators have completed progressing")

	return nil
}

// getConsole returns the console URL from the route 'console' in namespace openshift-console
func getConsole(ctx context.Context, config *rest.Config) (string, error) {
	url := ""
	// Need to keep these updated if they change
	consoleNamespace := "openshift-console"
	consoleRouteName := "console"
	rc, err := routeclient.NewForConfig(config)
	if err != nil {
		return "", errors.Wrap(err, "creating a route client")
	}

	consoleRouteTimeout := 2 * time.Minute
	logrus.Infof("Checking to see if there is a route at %s/%s...", consoleNamespace, consoleRouteName)
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
		} else if apierrors.IsNotFound(err) {
			logrus.Debug("OpenShift console route does not exist")
			cancel()
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
	pw, err := os.ReadFile(pwFile)
	if err != nil {
		return err
	}
	logrus.Info("Install complete!")
	logrus.Infof("To access the cluster as the system:admin user when using 'oc', run\n    export KUBECONFIG=%s", kubeconfig)
	if consoleURL != "" {
		logrus.Infof("Access the OpenShift web-console here: %s", consoleURL)
		if skipPasswordPrintFlag {
			logrus.Infof("Credentials omitted, if necessary verify the %s file", pwFile)
		} else {
			logrus.Infof("Login to the console with user: %q, and password: %q", "kubeadmin", pw)
		}
	}
	return nil
}

func waitForInstallComplete(ctx context.Context, config *rest.Config, directory string) error {
	if err := waitForInitializedCluster(ctx, config); err != nil {
		return err
	}

	if err := addRouterCAToClusterCA(ctx, config, command.RootOpts.Dir); err != nil {
		return err
	}

	if err := waitForStableOperators(ctx, config); err != nil {
		return err
	}

	consoleURL, err := getConsole(ctx, config)
	if err != nil {
		logrus.Warnf("Cluster does not have a console available: %v", err)
	}

	return logComplete(command.RootOpts.Dir, consoleURL)
}

func logTroubleshootingLink() {
	logrus.Error(`Cluster initialization failed because one or more operators are not functioning properly.
The cluster should be accessible for troubleshooting as detailed in the documentation linked below,
https://docs.openshift.com/container-platform/latest/support/troubleshooting/troubleshooting-installations.html
The 'wait-for install-complete' subcommand can then be used to continue the installation`)
}

func checkIfAgentCommand(assetStore asset.Store) {
	if agentConfig, err := assetStore.Load(&agentconfig.AgentConfig{}); err == nil && agentConfig != nil {
		logrus.Warning("An agent configuration was detected but this command is not the agent wait-for command")
	}
}

func waitForAllClusterOperators(clusterOperatorLister configlisters.ClusterOperatorLister) func(ctx context.Context) (bool, error) {
	previouslyStableOperators := sets.Set[string]{}

	return func(ctx context.Context) (bool, error) {
		stableOperators, unstableOperators, err := currentOperatorStability(clusterOperatorLister)
		if err != nil {
			return false, err
		}
		if newlyStableOperators := stableOperators.Difference(previouslyStableOperators); len(newlyStableOperators) > 0 {
			for _, name := range sets.List(newlyStableOperators) {
				logrus.Debugf("Cluster Operator %s is stable", name)
			}
		}
		if newlyUnstableOperators := previouslyStableOperators.Difference(stableOperators); len(newlyUnstableOperators) > 0 {
			for _, name := range sets.List(newlyUnstableOperators) {
				logrus.Debugf("Cluster Operator %s became unstable", name)
			}
		}
		previouslyStableOperators = stableOperators

		if len(unstableOperators) == 0 {
			return true, nil
		}

		return false, nil
	}
}

func currentOperatorStability(clusterOperatorLister configlisters.ClusterOperatorLister) (sets.Set[string], sets.Set[string], error) {
	clusterOperators, err := clusterOperatorLister.List(labels.Everything())
	if err != nil {
		return nil, nil, err // lister should never fail
	}

	stableOperators := sets.Set[string]{}
	unstableOperators := sets.Set[string]{}
	for _, clusterOperator := range clusterOperators {
		name := clusterOperator.Name
		progressing := cov1helpers.FindStatusCondition(clusterOperator.Status.Conditions, configv1.OperatorProgressing)
		if progressing == nil {
			logrus.Debugf("Cluster Operator %s progressing == nil", name)
			unstableOperators.Insert(name)
			continue
		}
		if meetsStabilityThreshold(progressing) {
			stableOperators.Insert(name)
		} else {
			logrus.Debugf("Cluster Operator %s is Progressing=%s LastTransitionTime=%v DurationSinceTransition=%.fs Reason=%s Message=%s", name, progressing.Status, progressing.LastTransitionTime.Time, time.Since(progressing.LastTransitionTime.Time).Seconds(), progressing.Reason, progressing.Message)
			unstableOperators.Insert(name)
		}
	}

	return stableOperators, unstableOperators, nil
}

func meetsStabilityThreshold(progressing *configv1.ClusterOperatorStatusCondition) bool {
	return progressing.Status == configv1.ConditionFalse && time.Since(progressing.LastTransitionTime.Time).Seconds() > coStabilityThreshold
}

func handleUnreachableAPIServer(ctx context.Context, config *rest.Config) error {
	assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
	if err != nil {
		return fmt.Errorf("failed to create asset store: %w", err)
	}

	// Ensure that the install is expecting the user to provision their own DNS solution.
	installConfig := &installconfig.InstallConfig{}
	if err := assetStore.Fetch(ctx, installConfig); err != nil {
		return fmt.Errorf("failed to fetch %s: %w", installConfig.Name(), err)
	}
	switch installConfig.Config.Platform.Name() { //nolint:gocritic
	case aws.Name:
		if installConfig.Config.AWS.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
			return nil
		}
	case gcp.Name:
		if installConfig.Config.GCP.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
			return nil
		}
	default:
		return nil
	}

	lbConfig := &lbconfig.Config{}
	if err := assetStore.Fetch(ctx, lbConfig); err != nil {
		return fmt.Errorf("failed to fetch %s: %w", lbConfig.Name(), err)
	}

	_, ipAddrs, err := lbConfig.ParseDNSDataFromConfig(lbconfig.PublicLoadBalancer)
	if err != nil {
		return fmt.Errorf("failed to parse lbconfig: %w", err)
	}

	// The kubeconfig handles one ip address
	ipAddr := ""
	if len(ipAddrs) > 0 {
		ipAddr = ipAddrs[0].String()
	}
	if ipAddr == "" {
		return fmt.Errorf("no ip address found in lbconfig")
	}

	dialer := &net.Dialer{
		Timeout:   1 * time.Minute,
		KeepAlive: 1 * time.Minute,
	}
	config.Dial = kubeconfig.CreateDialContext(dialer, ipAddr)

	// The asset is currently saved in <install-dir>/openshift. This directory
	// was consumed during install but this file is generated after that action. This
	// artifact will hang around unless it is purged here.
	if err := asset.DeleteAssetFromDisk(lbConfig, command.RootOpts.Dir); err != nil {
		return fmt.Errorf("failed to delete %s from disk", lbConfig.Name())
	}

	return nil
}
