package main

import (
	"context"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	operatorv1 "github.com/openshift/api/operator/v1"
	operatorclient "github.com/openshift/client-go/operator/clientset/versioned"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset"
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
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/vsphere"
	baremetalutils "github.com/openshift/installer/pkg/utils/baremetal"
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
		if err := command.LogClusterOperatorConditions(ctx, config); err != nil {
			logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err)
		}
		logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
		logrus.Error(err.Error())
		gatherAndAnalyzeBootstrapLogs(ctx, command.RootOpts.Dir)
		return command.ExitCodeBootstrapFailed, nil
	}
	timer.StopTimer("Bootstrap Complete")

	// In CI, we want to collect an installer log bundle so we can examine bootstrap logs that aren't
	// collectable anywhere else.
	if gatherBootstrap, ok := os.LookupEnv("OPENSHIFT_INSTALL_GATHER_BOOTSTRAP"); ok && gatherBootstrap != "" {
		timer.StartTimer("Bootstrap Gather")
		logrus.Infof("OPENSHIFT_INSTALL_GATHER_BOOTSTRAP is set, will attempt to gather a log bundle")
		bundlePath, gatherErr := runGatherBootstrapCmd(ctx, command.RootOpts.Dir)
		if gatherErr != nil {
			logrus.Warningf("Attempted to gather debug logs, and it failed: %q ", gatherErr)
		} else {
			logrus.Infof("Bootstrap gather logs captured here %q", bundlePath)
		}
		timer.StopTimer("Bootstrap Gather")
	}

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
	assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
	if err != nil {
		logrus.Error(err)
		logrus.Exit(command.ExitCodeInstallFailed)
	}

	err = command.WaitForInstallComplete(ctx, config, assetStore)
	if err != nil {
		if err2 := command.LogClusterOperatorConditions(ctx, config); err2 != nil {
			logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err2)
		}
		command.LogTroubleshootingLink()
		logrus.Error(err)
		return command.ExitCodeInstallFailed, nil
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
			t.command.PersistentFlags().BoolVar(&command.SkipPasswordPrintFlag, "skip-password-print", false, "Do not print the generated user password.")
		}
		cmd.AddCommand(t.command)
	}

	return cmd
}

func runTargetCmd(ctx context.Context, targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) {
	runner := func(directory string, consumeFiles bool) error {
		if !(forcePreserveInputs || forceConsumeInputs) {
			logrus.Info("Pass --preserve to leave input files in place")
		}

		fetcher := assetstore.NewAssetsFetcher(directory, consumeFiles)
		return fetcher.FetchAndPersist(ctx, targets)
	}

	return func(cmd *cobra.Command, args []string) {
		timer.StartTimer(timer.TotalTimeElapsed)

		// Set the context to be used in the PostRun function.
		cmd.SetContext(ctx)

		cleanup := command.SetupFileHook(command.RootOpts.Dir)
		defer cleanup()

		cluster.InstallDir = command.RootOpts.Dir

		err := runner(command.RootOpts.Dir, command.RootOpts.ConsumeFiles)
		if err != nil {
			if strings.Contains(err.Error(), asset.InstallConfigError) {
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallConfigError)
			}
			if strings.Contains(err.Error(), asset.ControlPlaneCreationError) {
				gatherAndAnalyzeBootstrapLogs(ctx, command.RootOpts.Dir)
			}
			if strings.Contains(err.Error(), asset.ClusterCreationError) {
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInfrastructureFailed)
			}
			logrus.Fatal(err)
		}
		switch cmd.Name() {
		case "cluster", "image", "pxe-files":
		default:
			logrus.Infof("%s", logging.LogCreatedFiles(cmd.Name(), command.RootOpts.Dir, targets))
		}
	}
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
		command.CheckIfAgentCommand(assetStore)
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
	logrus.Info("Waiting for the bootstrap etcd member to be removed...")
	client, err := operatorclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating operator client: %w", err)
	}
	// Validate the etcd operator has removed the bootstrap etcd member
	if err := wait.PollUntilContextCancel(ctx, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
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
	}); err != nil {
		logrus.Warnf("Bootstrap etcd member may not have been removed: %v", err)
		return err
	}
	logrus.Info("Bootstrap etcd member has been removed")
	return nil
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
	case azure.Name:
		if installConfig.Config.Azure.UserProvisionedDNS != dns.UserProvisionedDNSEnabled {
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

	lbType := lbconfig.PublicLoadBalancer
	if !installConfig.Config.PublicAPI() {
		lbType = lbconfig.PrivateLoadBalancer
	}

	_, ipAddrs, err := lbConfig.ParseDNSDataFromConfig(lbType)
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
