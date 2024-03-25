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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
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
	"github.com/openshift/installer/pkg/gather/service"
	timer "github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/vsphere"
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
				// Setup a context that is canceled when the user presses Ctrl+C,
				// or SIGTERM and SIGINT are received, this allows for a clean shutdown.
				ctx, cancel := context.WithCancel(context.TODO())
				defer cancel()
				logrus.RegisterExitHandler(cancel)

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

	targets = []target{installConfigTarget, manifestsTarget, ignitionConfigsTarget, clusterTarget, singleNodeIgnitionConfigTarget}
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
	if err := handleUnreachableAPIServer(config); err != nil {
		logrus.Fatal(fmt.Errorf("unable to handle api server override: %w", err))
	}

	//
	// Wait for the bootstrap to complete.
	//
	timer.StartTimer("Bootstrap Complete")
	if err := waitForBootstrapComplete(ctx, config); err != nil {
		bundlePath, gatherErr := runGatherBootstrapCmd(command.RootOpts.Dir)
		if gatherErr != nil {
			logrus.Error("Attempted to gather debug logs after installation failure: ", gatherErr)
		}
		if err := command.LogClusterOperatorConditions(ctx, config); err != nil {
			logrus.Error("Attempted to gather ClusterOperator status after installation failure: ", err)
		}
		logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
		logrus.Error(err.Error())
		if gatherErr == nil {
			if err := service.AnalyzeGatherBundle(bundlePath); err != nil {
				logrus.Error("Attempted to analyze the debug logs after installation failure: ", err)
			}
			logrus.Infof("Bootstrap gather logs captured here %q", bundlePath)
		}
		return command.ExitCodeBootstrapFailed, nil
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
	err = command.WaitForInstallComplete(ctx, config, command.RootOpts.Dir)
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
		fetcher := assetstore.NewAssetsFetcher(directory)
		return fetcher.FetchAndPersist(targets)
	}

	return func(cmd *cobra.Command, args []string) {
		timer.StartTimer(timer.TotalTimeElapsed)

		cleanup := command.SetupFileHook(command.RootOpts.Dir)
		defer cleanup()

		cluster.InstallDir = command.RootOpts.Dir

		err := runner(command.RootOpts.Dir)
		if err != nil {
			if strings.Contains(err.Error(), asset.InstallConfigError) {
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallConfigError)
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
			logrus.Infof(logging.LogCreatedFiles(cmd.Name(), command.RootOpts.Dir, targets))
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
		checkIfAgentCommand(assetStore)
	}

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

	if err := waitForBootstrapConfigMap(ctx, client); err != nil {
		return err
	}

	if err := waitForStableSNOBootstrap(ctx, config); err != nil {
		return newBootstrapError(err)
	}

	return nil
}

// waitForBootstrapConfigMap watches the configmaps in the kube-system namespace
// and waits for the bootstrap configmap to report that bootstrapping has
// completed.
func waitForBootstrapConfigMap(ctx context.Context, client *kubernetes.Clientset) *clusterCreateError {
	timeout := 30 * time.Minute

	// Wait longer for baremetal, VSphere due to length of time it takes to boot
	if assetStore, err := assetstore.NewStore(command.RootOpts.Dir); err == nil {
		if installConfig, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
			if installConfig.(*installconfig.InstallConfig).Config.Platform.Name() == baremetal.Name || installConfig.(*installconfig.InstallConfig).Config.Platform.Name() == vsphere.Name {
				timeout = 60 * time.Minute
			}
		}
	}

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for bootstrapping to complete...",
		timeout, untilTime.Format(time.Kitchen), timezone)

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

// When bootstrap on SNO deployments, we should not remove the bootstrap node prematurely,
// here we make sure that the deployment is stable.
// Given the nature of single node we just need to make sure things such as etcd are in the proper state
// before continuing.
func waitForStableSNOBootstrap(ctx context.Context, config *rest.Config) error {
	timeout := 5 * time.Minute

	// If we're not in a single node deployment, bail early
	if isSNO, err := IsSingleNode(); err != nil {
		logrus.Warningf("Can not determine if installing a Single Node cluster, continuing as normal install: %v", err)
		return nil
	} else if !isSNO {
		return nil
	}

	snoBootstrapContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Info("Detected Single Node deployment")
	logrus.Infof("Waiting up to %v (until %v %s) for the bootstrap etcd member to be removed...",
		timeout, untilTime.Format(time.Kitchen), timezone)

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating dynamic client: %w", err)
	}
	gvr := schema.GroupVersionResource{
		Group:    operatorv1.SchemeGroupVersion.Group,
		Version:  operatorv1.SchemeGroupVersion.Version,
		Resource: "etcds",
	}
	resourceClient := client.Resource(gvr)
	// Validate the etcd operator has removed the bootstrap etcd member
	return wait.PollUntilContextCancel(snoBootstrapContext, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		etcdOperator := &operatorv1.Etcd{}
		etcdUnstructured, err := resourceClient.Get(ctx, "cluster", metav1.GetOptions{})
		if err != nil {
			// There might be service disruptions in SNO, we log those here but keep trying with in the time limit
			logrus.Debugf("Error getting ETCD Cluster resource, retrying: %v", err)
			return false, nil
		}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(etcdUnstructured.Object, etcdOperator)
		if err != nil {
			// This error should not happen, if we do, we log the error and keep retrying until we hit the limit
			logrus.Debugf("Error parsing etcds resource, retrying: %v", err)
			return false, nil
		}
		for _, condition := range etcdOperator.Status.Conditions {
			if condition.Type == "EtcdBootstrapMemberRemoved" {
				return configv1.ConditionStatus(condition.Status) == configv1.ConditionTrue, nil
			}
		}
		return false, nil
	})
}

func checkIfAgentCommand(assetStore asset.Store) {
	if agentConfig, err := assetStore.Load(&agentconfig.AgentConfig{}); err == nil && agentConfig != nil {
		logrus.Warning("An agent configuration was detected but this command is not the agent wait-for command")
	}
}

func handleUnreachableAPIServer(config *rest.Config) error {
	assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
	if err != nil {
		return fmt.Errorf("failed to create asset store: %w", err)
	}

	// Ensure that the install is expecting the user to provision their own DNS solution.
	installConfig := &installconfig.InstallConfig{}
	if err := assetStore.Fetch(installConfig); err != nil {
		return fmt.Errorf("failed to fetch %s: %w", installConfig.Name(), err)
	}
	switch installConfig.Config.Platform.Name() { //nolint:gocritic
	case gcp.Name:
		if installConfig.Config.GCP.UserProvisionedDNS != gcp.UserProvisionedDNSEnabled {
			return nil
		}
	default:
		return nil
	}

	lbConfig := &lbconfig.Config{}
	if err := assetStore.Fetch(lbConfig); err != nil {
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

// IsSingleNode determines if we are in a single node configuration based off of the install config
// loaded from the asset store.
func IsSingleNode() (bool, error) {
	assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
	if err != nil {
		return false, fmt.Errorf("error loading asset store: %w", err)
	}
	installConfig, err := assetStore.Load(&installconfig.InstallConfig{})
	if err != nil {
		return false, fmt.Errorf("error loading installConfig: %w", err)
	}
	if installConfig == nil {
		return false, fmt.Errorf("installConfig loaded from asset store was nil")
	}

	config := installConfig.(*installconfig.InstallConfig).Config
	if machinePool := config.ControlPlane; machinePool != nil {
		return *machinePool.Replicas == int64(1), nil
	}
	return false, nil
}
