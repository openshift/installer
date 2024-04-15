/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package main is the main package for CAPV.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	goruntime "runtime"
	"time"

	perrors "github.com/pkg/errors"
	"github.com/spf13/pflag"
	"gopkg.in/fsnotify.v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	_ "k8s.io/component-base/logs/json/register"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	capiflags "sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/controllers"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/constants"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/version"
)

var (
	setupLog       = ctrl.Log.WithName("setup")
	logOptions     = logs.NewOptions()
	controllerName = "cluster-api-vsphere-manager"

	enableContentionProfiling   bool
	leaderElectionLeaseDuration time.Duration
	leaderElectionRenewDeadline time.Duration
	leaderElectionRetryPeriod   time.Duration
	managerOpts                 manager.Options
	restConfigBurst             int
	restConfigQPS               float32
	syncPeriod                  time.Duration
	webhookOpts                 webhook.Options
	watchNamespace              string

	clusterCacheTrackerConcurrency    int
	vSphereClusterConcurrency         int
	vSphereMachineConcurrency         int
	providerServiceAccountConcurrency int
	serviceDiscoveryConcurrency       int
	vSphereVMConcurrency              int
	vSphereClusterIdentityConcurrency int
	vSphereDeploymentZoneConcurrency  int

	tlsOptions         = capiflags.TLSOptions{}
	diagnosticsOptions = capiflags.DiagnosticsOptions{}

	defaultProfilerAddr      = os.Getenv("PROFILER_ADDR")
	defaultSyncPeriod        = manager.DefaultSyncPeriod
	defaultLeaderElectionID  = manager.DefaultLeaderElectionID
	defaultPodName           = manager.DefaultPodName
	defaultWebhookPort       = manager.DefaultWebhookServiceContainerPort
	defaultEnableKeepAlive   = constants.DefaultEnableKeepAlive
	defaultKeepAliveDuration = constants.DefaultKeepAliveDuration
)

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	// Flags specific to CAPV

	fs.StringVar(
		&managerOpts.LeaderElectionID,
		"leader-election-id",
		defaultLeaderElectionID,
		"Name of the config map to use as the locking resource when configuring leader election.")

	fs.IntVar(&clusterCacheTrackerConcurrency, "clustercachetracker-concurrency", 10,
		"Number of clusters to process simultaneously")

	fs.IntVar(&vSphereClusterConcurrency, "vspherecluster-concurrency", 10,
		"Number of vSphere clusters to process simultaneously")

	fs.IntVar(&vSphereMachineConcurrency, "vspheremachine-concurrency", 10,
		"Number of vSphere machines to process simultaneously")

	fs.IntVar(&providerServiceAccountConcurrency, "providerserviceaccount-concurrency", 10,
		"Number of provider service accounts to process simultaneously")

	fs.IntVar(&serviceDiscoveryConcurrency, "servicediscovery-concurrency", 10,
		"Number of vSphere clusters for service discovery to process simultaneously")

	fs.IntVar(&vSphereVMConcurrency, "vspherevm-concurrency", 10,
		"Number of vSphere vms to process simultaneously")

	fs.IntVar(&vSphereClusterIdentityConcurrency, "vsphereclusteridentity-concurrency", 10,
		"Number of vSphere cluster identities to process simultaneously")

	fs.IntVar(&vSphereDeploymentZoneConcurrency, "vspheredeploymentzone-concurrency", 10,
		"Number of vSphere deployment zones to process simultaneously")

	fs.StringVar(
		&managerOpts.PodName,
		"pod-name",
		defaultPodName,
		"The name of the pod running the controller manager.")

	fs.StringVar(
		&managerOpts.CredentialsFile,
		"credentials-file",
		"/etc/capv/credentials.yaml",
		"path to CAPV's credentials file",
	)
	fs.BoolVar(
		&managerOpts.EnableKeepAlive,
		"enable-keep-alive",
		defaultEnableKeepAlive,
		"feature to enable keep alive handler in vsphere sessions. This functionality is enabled by default.")
	fs.DurationVar(
		&managerOpts.KeepAliveDuration,
		"keep-alive-duration",
		defaultKeepAliveDuration,
		"idle time interval(minutes) in between send() requests in keepalive handler",
	)
	fs.StringVar(
		&managerOpts.NetworkProvider,
		"network-provider",
		"",
		"network provider to be used by Supervisor based clusters.",
	)

	// Flags common between CAPI and CAPV

	logsv1.AddFlags(logOptions, fs)

	fs.BoolVar(&managerOpts.LeaderElection, "leader-elect", true,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	fs.DurationVar(&leaderElectionLeaseDuration, "leader-elect-lease-duration", 15*time.Second,
		"Interval at which non-leader candidates will wait to force acquire leadership (duration string)")

	fs.DurationVar(&leaderElectionRenewDeadline, "leader-elect-renew-deadline", 10*time.Second,
		"Duration that the leading controller manager will retry refreshing leadership before giving up (duration string)")

	fs.DurationVar(&leaderElectionRetryPeriod, "leader-elect-retry-period", 2*time.Second,
		"Duration the LeaderElector clients should wait between tries of actions (duration string)")

	fs.StringVar(&watchNamespace, "namespace", "",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")

	fs.StringVar(&managerOpts.WatchFilterValue, "watch-filter", "",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel))

	fs.StringVar(&managerOpts.PprofBindAddress, "profiler-address", defaultProfilerAddr,
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")

	fs.BoolVar(&enableContentionProfiling, "contention-profiling", false,
		"Enable block profiling.")

	fs.DurationVar(&syncPeriod, "sync-period", defaultSyncPeriod,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)")

	fs.Float32Var(&restConfigQPS, "kube-api-qps", 20,
		"Maximum queries per second from the controller client to the Kubernetes API server. Defaults to 20")

	fs.IntVar(&restConfigBurst, "kube-api-burst", 30,
		"Maximum number of queries that should be allowed in one burst from the controller client to the Kubernetes API server. Default 30")

	fs.IntVar(&webhookOpts.Port, "webhook-port", defaultWebhookPort,
		"Webhook Server port")

	fs.StringVar(&webhookOpts.CertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")

	fs.StringVar(&managerOpts.HealthProbeBindAddress, "health-addr", ":9440",
		"The address the health endpoint binds to.",
	)

	capiflags.AddTLSOptions(fs, &tlsOptions)
	capiflags.AddDiagnosticsOptions(fs, &diagnosticsOptions)
	feature.MutableGates.AddFlag(fs)
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	// Set log level 2 as default.
	if err := pflag.CommandLine.Set("v", "2"); err != nil {
		setupLog.Error(err, "failed to set default log level")
		os.Exit(1)
	}
	pflag.Parse()

	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	managerOpts.KubeConfig = ctrl.GetConfigOrDie()
	managerOpts.KubeConfig.QPS = restConfigQPS
	managerOpts.KubeConfig.Burst = restConfigBurst
	managerOpts.KubeConfig.UserAgent = remote.DefaultClusterAPIUserAgent(controllerName)

	if watchNamespace != "" {
		managerOpts.Cache.DefaultNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	if enableContentionProfiling {
		goruntime.SetBlockProfileRate(1)
	}

	setupLog.Info(fmt.Sprintf("Feature gates: %+v\n", feature.Gates))

	managerOpts.Cache.SyncPeriod = &syncPeriod
	managerOpts.LeaseDuration = &leaderElectionLeaseDuration
	managerOpts.RenewDeadline = &leaderElectionRenewDeadline
	managerOpts.RetryPeriod = &leaderElectionRetryPeriod

	// Create a function that adds all the controllers and webhooks to the manager.
	addToManager := func(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager) error {
		tracker, err := setupRemoteClusterCacheTracker(ctx, mgr)
		if err != nil {
			return perrors.Wrapf(err, "unable to create remote cluster cache tracker")
		}

		// Check for non-supervisor VSphereCluster and start controller if found
		gvr := infrav1.GroupVersion.WithResource(reflect.TypeOf(&infrav1.VSphereCluster{}).Elem().Name())
		isNonSupervisorCRDLoaded, err := isCRDDeployed(mgr, gvr)
		if err != nil {
			return err
		}
		if isNonSupervisorCRDLoaded {
			if err := setupVAPIControllers(ctx, controllerCtx, mgr, tracker); err != nil {
				return fmt.Errorf("setupVAPIControllers: %w", err)
			}
		} else {
			setupLog.Info(fmt.Sprintf("CRD for %s not loaded, skipping.", gvr.String()))
		}

		// Check for supervisor VSphereCluster and start controller if found
		gvr = vmwarev1.GroupVersion.WithResource(reflect.TypeOf(&vmwarev1.VSphereCluster{}).Elem().Name())
		isSupervisorCRDLoaded, err := isCRDDeployed(mgr, gvr)
		if err != nil {
			return err
		}
		if isSupervisorCRDLoaded {
			if err := setupSupervisorControllers(ctx, controllerCtx, mgr, tracker); err != nil {
				return fmt.Errorf("setupSupervisorControllers: %w", err)
			}
		} else {
			setupLog.Info(fmt.Sprintf("CRD for %s not loaded, skipping.", gvr.String()))
		}

		// Continuing startup does not make sense without having managers added.
		if !isSupervisorCRDLoaded && !isNonSupervisorCRDLoaded {
			return errors.New("neither supervisor nor non-supervisor CRDs detected")
		}

		return nil
	}

	tlsOptionOverrides, err := capiflags.GetTLSOptionOverrideFuncs(tlsOptions)
	if err != nil {
		setupLog.Error(err, "unable to add TLS settings to the webhook server")
		os.Exit(1)
	}
	webhookOpts.TLSOpts = tlsOptionOverrides
	managerOpts.WebhookServer = webhook.NewServer(webhookOpts)
	managerOpts.AddToManager = addToManager
	managerOpts.Metrics = capiflags.GetDiagnosticsOptions(diagnosticsOptions)

	// Set up the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	mgr, err := manager.New(ctx, managerOpts)
	if err != nil {
		setupLog.Error(err, "Error creating manager")
		os.Exit(1)
	}

	setupChecks(mgr)

	setupLog.Info("Starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "Error starting manager")
		os.Exit(1)
	}

	// initialize notifier for capv-manager-bootstrap-credentials
	watch, err := manager.InitializeWatch(mgr.GetControllerManagerContext(), &managerOpts)
	if err != nil {
		setupLog.Error(err, "failed to initialize watch on CAPV credentials file")
		os.Exit(1)
	}
	defer func(watch *fsnotify.Watcher) {
		_ = watch.Close()
	}(watch)
	defer session.Clear()
}

func setupVAPIControllers(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager, tracker *remote.ClusterCacheTracker) error {
	if err := (&webhooks.VSphereClusterTemplateWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereMachineWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereMachineTemplateWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereVMWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereDeploymentZoneWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereFailureDomainWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := controllers.AddClusterControllerToManager(ctx, controllerCtx, mgr, false, concurrency(vSphereClusterConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddMachineControllerToManager(ctx, controllerCtx, mgr, false, concurrency(vSphereMachineConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddVMControllerToManager(ctx, controllerCtx, mgr, tracker, concurrency(vSphereVMConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddVsphereClusterIdentityControllerToManager(ctx, controllerCtx, mgr, concurrency(vSphereClusterIdentityConcurrency)); err != nil {
		return err
	}

	return controllers.AddVSphereDeploymentZoneControllerToManager(ctx, controllerCtx, mgr, concurrency(vSphereDeploymentZoneConcurrency))
}

func setupSupervisorControllers(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager, tracker *remote.ClusterCacheTracker) error {
	if err := controllers.AddClusterControllerToManager(ctx, controllerCtx, mgr, true, concurrency(vSphereClusterConcurrency)); err != nil {
		return err
	}

	if err := controllers.AddMachineControllerToManager(ctx, controllerCtx, mgr, true, concurrency(vSphereMachineConcurrency)); err != nil {
		return err
	}

	if err := controllers.AddServiceAccountProviderControllerToManager(ctx, controllerCtx, mgr, tracker, concurrency(providerServiceAccountConcurrency)); err != nil {
		return err
	}

	return controllers.AddServiceDiscoveryControllerToManager(ctx, controllerCtx, mgr, tracker, concurrency(serviceDiscoveryConcurrency))
}

func setupChecks(mgr ctrlmgr.Manager) {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}
}

func isCRDDeployed(mgr ctrlmgr.Manager, gvr schema.GroupVersionResource) (bool, error) {
	_, err := mgr.GetRESTMapper().KindFor(gvr)
	if err != nil {
		var discoveryErr *apiutil.ErrResourceDiscoveryFailed
		ok := errors.As(errors.Unwrap(err), &discoveryErr)
		if !ok {
			return false, err
		}
		discoveryErrs := *discoveryErr
		gvrErr, ok := discoveryErrs[gvr.GroupVersion()]
		if !ok {
			return false, err
		}
		if apierrors.IsNotFound(gvrErr) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func concurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}

func setupRemoteClusterCacheTracker(ctx context.Context, mgr ctrlmgr.Manager) (*remote.ClusterCacheTracker, error) {
	secretCachingClient, err := client.New(mgr.GetConfig(), client.Options{
		HTTPClient: mgr.GetHTTPClient(),
		Cache: &client.CacheOptions{
			Reader: mgr.GetCache(),
		},
	})
	if err != nil {
		return nil, perrors.Wrapf(err, "unable to create secret caching client")
	}

	// Set up a ClusterCacheTracker and ClusterCacheReconciler to provide to controllers
	// requiring a connection to a remote cluster
	log := ctrl.Log.WithValues("component", "remote/clustercachetracker")
	tracker, err := remote.NewClusterCacheTracker(
		mgr,
		remote.ClusterCacheTrackerOptions{
			SecretCachingClient: secretCachingClient,
			ControllerName:      controllerName,
			Log:                 &log,
		},
	)
	if err != nil {
		return nil, perrors.Wrapf(err, "unable to create cluster cache tracker")
	}

	if err := (&remote.ClusterCacheReconciler{
		Client:           mgr.GetClient(),
		Tracker:          tracker,
		WatchFilterValue: managerOpts.WatchFilterValue,
	}).SetupWithManager(ctx, mgr, concurrency(clusterCacheTrackerConcurrency)); err != nil {
		return nil, perrors.Wrapf(err, "unable to create ClusterCacheReconciler controller")
	}

	return tracker, nil
}
