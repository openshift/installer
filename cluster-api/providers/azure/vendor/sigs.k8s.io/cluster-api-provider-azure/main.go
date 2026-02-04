/*
Copyright 2020 The Kubernetes Authors.

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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	// +kubebuilder:scaffold:imports
	asocontainerservicev1api20210501 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20210501"
	asocontainerservicev1api20230201 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20230201"
	asocontainerservicev1api20230315preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20230315preview"
	asocontainerservicev1api20231001 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	asocontainerservicev1api20231102preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231102preview"
	asocontainerservicev1api20240402preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240402preview"
	asocontainerservicev1api20240901 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240901"
	asokubernetesconfigurationv1 "github.com/Azure/azure-service-operator/v2/api/kubernetesconfiguration/v1api20230501"
	asonetworkv1api20201101 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	asonetworkv1api20220701 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	asoresourcesv1 "github.com/Azure/azure-service-operator/v2/api/resources/v1api20200601"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/clustercache"
	"sigs.k8s.io/cluster-api/controllers/remote"
	capifeature "sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util/flags"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	infrav1controllersexp "sigs.k8s.io/cluster-api-provider-azure/exp/controllers"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/ot"
	"sigs.k8s.io/cluster-api-provider-azure/util/components"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/version"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = infrav1exp.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = bootstrapv1.AddToScheme(scheme)
	_ = asoresourcesv1.AddToScheme(scheme)
	_ = asocontainerservicev1api20210501.AddToScheme(scheme)
	_ = asocontainerservicev1api20230201.AddToScheme(scheme)
	_ = asocontainerservicev1api20231001.AddToScheme(scheme)
	_ = asonetworkv1api20220701.AddToScheme(scheme)
	_ = asonetworkv1api20201101.AddToScheme(scheme)
	_ = asocontainerservicev1api20230315preview.AddToScheme(scheme)
	_ = asocontainerservicev1api20231102preview.AddToScheme(scheme)
	_ = asocontainerservicev1api20240402preview.AddToScheme(scheme)
	_ = asocontainerservicev1api20240901.AddToScheme(scheme)
	_ = asokubernetesconfigurationv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	enableLeaderElection               bool
	leaderElectionNamespace            string
	leaderElectionLeaseDuration        time.Duration
	leaderElectionRenewDeadline        time.Duration
	leaderElectionRetryPeriod          time.Duration
	watchNamespace                     string
	watchFilterValue                   string
	profilerAddress                    string
	azureClusterConcurrency            int
	azureMachineConcurrency            int
	azureMachinePoolConcurrency        int
	azureMachinePoolMachineConcurrency int
	deprecatedAzureBootrapConfigGVK    string // Deprecated in v1.19.0
	debouncingTimer                    time.Duration
	syncPeriod                         time.Duration
	healthAddr                         string
	webhookPort                        int
	webhookCertDir                     string
	managerOptions                     = flags.ManagerOptions{}
	timeouts                           reconciler.Timeouts
	enableTracing                      bool
	disableControllersOrWebhooks       []string
)

// InitFlags initializes all command-line flags.
func InitFlags(fs *pflag.FlagSet) {
	fs.BoolVar(
		&enableLeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
	)

	flag.StringVar(
		&leaderElectionNamespace,
		"leader-election-namespace",
		"",
		"Namespace that the controller performs leader election in. If unspecified, the controller will discover which namespace it is running in.",
	)

	fs.DurationVar(
		&leaderElectionLeaseDuration,
		"leader-elect-lease-duration",
		15*time.Second,
		"Interval at which non-leader candidates will wait to force acquire leadership (duration string)",
	)

	fs.DurationVar(
		&leaderElectionRenewDeadline,
		"leader-elect-renew-deadline",
		10*time.Second,
		"Duration that the leading controller manager will retry refreshing leadership before giving up (duration string)",
	)

	fs.DurationVar(
		&leaderElectionRetryPeriod,
		"leader-elect-retry-period",
		2*time.Second,
		"Duration the LeaderElector clients should wait between tries of actions (duration string)",
	)

	fs.StringVar(
		&watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	fs.StringVar(
		&watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel),
	)

	fs.StringVar(
		&profilerAddress,
		"profiler-address",
		"",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)",
	)

	fs.IntVar(&azureClusterConcurrency,
		"azurecluster-concurrency",
		10,
		"Number of AzureClusters to process simultaneously",
	)

	fs.IntVar(&azureMachineConcurrency,
		"azuremachine-concurrency",
		10,
		"Number of AzureMachines to process simultaneously",
	)

	fs.IntVar(&azureMachinePoolConcurrency,
		"azuremachinepool-concurrency",
		10,
		"Number of AzureMachinePools to process simultaneously")

	fs.IntVar(&azureMachinePoolMachineConcurrency,
		"azuremachinepoolmachine-concurrency",
		10,
		"Number of AzureMachinePoolMachines to process simultaneously")

	fs.DurationVar(&debouncingTimer,
		"debouncing-timer",
		10*time.Second,
		"The minimum interval the controller should wait after a successful reconciliation of a particular object before reconciling it again",
	)

	fs.DurationVar(&syncPeriod,
		"sync-period",
		10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)",
	)

	fs.StringVar(&healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	fs.IntVar(&webhookPort,
		"webhook-port",
		9443,
		"The webhook server port the manager will listen on.",
	)

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"The webhook certificate directory, where the server should find the TLS certificate and key.")

	fs.DurationVar(&timeouts.Loop,
		"reconcile-timeout",
		reconciler.DefaultLoopTimeout,
		"The maximum duration a reconcile loop can run (e.g. 10m)",
	)

	fs.DurationVar(&timeouts.AzureServiceReconcile,
		"service-reconcile-timeout",
		reconciler.DefaultAzureServiceReconcileTimeout,
		"The maximum duration each Azure service reconcile can run (e.g. 90m)",
	)

	fs.DurationVar(&timeouts.AzureCall,
		"api-call-timeout",
		reconciler.DefaultAzureCallTimeout,
		"The maximum duration CAPZ will wait for each Azure API request before it is considered long running and performed async (e.g. 10s)",
	)

	fs.DurationVar(&timeouts.Requeue,
		"reconciler-requeue",
		reconciler.DefaultReconcilerRequeue,
		"The duration to wait before retrying after a transient reconcile error occurs (e.g. 15s)",
	)

	fs.BoolVar(
		&enableTracing,
		"enable-tracing",
		false,
		"Enable tracing to the opentelemetry-collector service in the same namespace.",
	)

	fs.StringVar(&deprecatedAzureBootrapConfigGVK,
		"bootstrap-config-gvk",
		"",
		"(Deprecated) Provide fully qualified GVK string to override default kubeadm config watch source, in the form of Kind.version.group (default: KubeadmConfig.v1beta1.bootstrap.cluster.x-k8s.io)",
	)

	fs.StringSliceVar(&disableControllersOrWebhooks,
		"disable-controllers-or-webhooks",
		[]string{},
		"Comma-separated list of controllers or webhooks to disable. The list can contain the following values: DisableASOSecretController,DisableAzureJSONMachineController",
	)

	flags.AddManagerOptions(fs, &managerOptions)

	feature.MutableGates.AddFlag(fs)
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

func main() {
	InitFlags(pflag.CommandLine)
	klog.InitFlags(flag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	tlsOptions, metricsOptions, err := flags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: invalid flags")
		os.Exit(1)
	}

	if deprecatedAzureBootrapConfigGVK != "" {
		setupLog.Error(fmt.Errorf("bootstrap-config-gvk argument is deprecated and no longer needed"), "Deprecated argument")
	}

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	// Validate valid disable components were passed in the flag
	if len(disableControllersOrWebhooks) > 0 {
		for _, component := range disableControllersOrWebhooks {
			if ok := components.IsValidDisableComponent(component); !ok {
				setupLog.Error(fmt.Errorf("invalid disable-controllers-or-webhooks value %s", component), "Invalid argument")
				os.Exit(1)
			}
		}
	}

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-azure-manager"
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                     scheme,
		LeaderElection:             enableLeaderElection,
		LeaderElectionID:           "controller-leader-election-capz",
		LeaderElectionNamespace:    leaderElectionNamespace,
		LeaseDuration:              &leaderElectionLeaseDuration,
		RenewDeadline:              &leaderElectionRenewDeadline,
		RetryPeriod:                &leaderElectionRetryPeriod,
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		HealthProbeBindAddress:     healthAddr,
		PprofBindAddress:           profilerAddress,
		Metrics:                    *metricsOptions,
		Cache: cache.Options{
			DefaultNamespaces: watchNamespaces,
			SyncPeriod:        &syncPeriod,
		},
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					&corev1.ConfigMap{},
					&corev1.Secret{},
				},
			},
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    webhookPort,
			CertDir: webhookCertDir,
			TLSOpts: tlsOptions,
		}),
		EventBroadcaster: broadcaster,
	})

	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("azure-controller"))

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	if enableTracing {
		if err := ot.RegisterTracing(ctx, setupLog); err != nil {
			setupLog.Error(err, "unable to initialize tracing")
			os.Exit(1)
		}
	}

	if err := ot.RegisterMetrics(); err != nil {
		setupLog.Error(err, "unable to initialize metrics")
		os.Exit(1)
	}

	registerControllers(ctx, mgr)

	registerWebhooks(mgr)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func registerControllers(ctx context.Context, mgr manager.Manager) {
	credCache := azure.NewCredentialCache()

	machineCache, err := coalescing.NewRequestCache(debouncingTimer)
	if err != nil {
		setupLog.Error(err, "failed to build machineCache ReconcileCache")
	}
	if err := controllers.NewAzureMachineReconciler(mgr.GetClient(),
		mgr.GetEventRecorderFor("azuremachine-reconciler"),
		timeouts,
		watchFilterValue,
		credCache,
	).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureMachineConcurrency}, Cache: machineCache}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AzureMachine")
		os.Exit(1)
	}

	clusterCache, err := coalescing.NewRequestCache(debouncingTimer)
	if err != nil {
		setupLog.Error(err, "failed to build clusterCache ReconcileCache")
	}
	if err := controllers.NewAzureClusterReconciler(
		mgr.GetClient(),
		mgr.GetEventRecorderFor("azurecluster-reconciler"),
		timeouts,
		watchFilterValue,
		credCache,
	).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}, Cache: clusterCache}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AzureCluster")
		os.Exit(1)
	}

	if err := (&controllers.AzureJSONTemplateReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("azurejsontemplate-reconciler"),
		Timeouts:         timeouts,
		WatchFilterValue: watchFilterValue,
		CredentialCache:  credCache,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureMachineConcurrency}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AzureJSONTemplate")
		os.Exit(1)
	}

	if !components.IsComponentDisabled(disableControllersOrWebhooks, infrav1.DisableAzureJSONMachineController) {
		if err := (&controllers.AzureJSONMachineReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("azurejsonmachine-reconciler"),
			Timeouts:         timeouts,
			WatchFilterValue: watchFilterValue,
			CredentialCache:  credCache,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureMachineConcurrency, SkipNameValidation: ptr.To(true)}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureJSONMachine")
			os.Exit(1)
		}
	}

	if !components.IsComponentDisabled(disableControllersOrWebhooks, infrav1.DisableASOSecretController) {
		if err := (&controllers.ASOSecretReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("asosecret-reconciler"),
			Timeouts:         timeouts,
			WatchFilterValue: watchFilterValue,
			CredentialCache:  credCache,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ASOSecret")
			os.Exit(1)
		}
	}

	// just use CAPI MachinePool feature flag rather than create a new one
	setupLog.V(1).Info(fmt.Sprintf("%+v\n", feature.Gates))
	if feature.Gates.Enabled(capifeature.MachinePool) {
		mpCache, err := coalescing.NewRequestCache(debouncingTimer)
		if err != nil {
			setupLog.Error(err, "failed to build mpCache ReconcileCache")
		}

		if err := infrav1controllersexp.NewAzureMachinePoolReconciler(
			mgr.GetClient(),
			mgr.GetEventRecorderFor("azuremachinepool-reconciler"),
			timeouts,
			watchFilterValue,
			credCache,
		).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureMachinePoolConcurrency}, Cache: mpCache}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureMachinePool")
			os.Exit(1)
		}

		mpmCache, err := coalescing.NewRequestCache(debouncingTimer)
		if err != nil {
			setupLog.Error(err, "failed to build mpmCache ReconcileCache")
		}

		if err := infrav1controllersexp.NewAzureMachinePoolMachineController(
			mgr.GetClient(),
			mgr.GetEventRecorderFor("azuremachinepoolmachine-reconciler"),
			timeouts,
			watchFilterValue,
			credCache,
		).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureMachinePoolMachineConcurrency}, Cache: mpmCache}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureMachinePoolMachine")
			os.Exit(1)
		}

		if err := (&controllers.AzureJSONMachinePoolReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("azurejsonmachinepool-reconciler"),
			Timeouts:         timeouts,
			WatchFilterValue: watchFilterValue,
			CredentialCache:  credCache,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureMachinePoolConcurrency, SkipNameValidation: ptr.To(true)}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureJSONMachinePool")
			os.Exit(1)
		}

		mmpmCache, err := coalescing.NewRequestCache(debouncingTimer)
		if err != nil {
			setupLog.Error(err, "failed to build mmpmCache ReconcileCache")
		}

		if err := controllers.NewAzureManagedMachinePoolReconciler(
			mgr.GetClient(),
			mgr.GetEventRecorderFor("azuremanagedmachinepoolmachine-reconciler"),
			timeouts,
			watchFilterValue,
			credCache,
		).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureMachinePoolConcurrency}, Cache: mmpmCache}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureManagedMachinePool")
			os.Exit(1)
		}

		mcCache, err := coalescing.NewRequestCache(debouncingTimer)
		if err != nil {
			setupLog.Error(err, "failed to build mcCache ReconcileCache")
		}

		if err := (&controllers.AzureManagedClusterReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("azuremanagedcluster-reconciler"),
			Timeouts:         timeouts,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}, Cache: mcCache}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureManagedCluster")
			os.Exit(1)
		}

		mcpCache, err := coalescing.NewRequestCache(debouncingTimer)
		if err != nil {
			setupLog.Error(err, "failed to build mcpCache ReconcileCache")
		}

		if err := (&controllers.AzureManagedControlPlaneReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("azuremanagedcontrolplane-reconciler"),
			Timeouts:         timeouts,
			WatchFilterValue: watchFilterValue,
			CredentialCache:  credCache,
		}).SetupWithManager(ctx, mgr, controllers.Options{Options: controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}, Cache: mcpCache}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureManagedControlPlane")
			os.Exit(1)
		}
	}

	if feature.Gates.Enabled(feature.ASOAPI) {
		if err := (&controllers.AzureASOManagedClusterReconciler{
			Client:           mgr.GetClient(),
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureASOManagedCluster")
			os.Exit(1)
		}

		if err := (&controllers.AzureASOManagedControlPlaneReconciler{
			Client:           mgr.GetClient(),
			WatchFilterValue: watchFilterValue,
			CredentialCache:  controllers.NewASOCredentialCache(credCache, mgr.GetClient()),
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureASOManagedControlPlane")
			os.Exit(1)
		}

		// The AzureASOManagedMachinePool controller reads the nodes in clusters to set provider IDs.
		secretCachingClient, err := client.New(mgr.GetConfig(), client.Options{
			HTTPClient: mgr.GetHTTPClient(),
			Cache: &client.CacheOptions{
				Reader: mgr.GetCache(),
			},
		})
		if err != nil {
			setupLog.Error(err, "unable to create secret caching client")
			os.Exit(1)
		}
		clusterCache, err := clustercache.SetupWithManager(ctx, mgr, clustercache.Options{
			SecretClient: secretCachingClient,
			Cache:        clustercache.CacheOptions{},
			Client: clustercache.ClientOptions{
				UserAgent: remote.DefaultClusterAPIUserAgent("azure-controller"),
				Cache: clustercache.ClientCacheOptions{
					DisableFor: []client.Object{
						// Don't cache ConfigMaps or Secrets.
						&corev1.ConfigMap{},
						&corev1.Secret{},
					},
				},
			},
			WatchFilterValue: watchFilterValue,
		}, controller.Options{})
		if err != nil {
			setupLog.Error(err, "unable to create cluster cache")
			os.Exit(1)
		}

		if err := (&controllers.AzureASOManagedMachinePoolReconciler{
			Client:           mgr.GetClient(),
			WatchFilterValue: watchFilterValue,
			Tracker:          clusterCache,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureMachinePoolConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AzureASOManagedMachinePool")
			os.Exit(1)
		}

		if err := (&controllers.ManagedClusterAdoptReconciler{
			Client: mgr.GetClient(),
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureClusterConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ManagedCluster")
			os.Exit(1)
		}

		if err := (&controllers.AgentPoolAdoptReconciler{
			Client: mgr.GetClient(),
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: azureMachinePoolConcurrency}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AgentPool")
			os.Exit(1)
		}
	}
}

func registerWebhooks(mgr manager.Manager) {
	if err := (&infrav1.AzureCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AzureCluster")
		os.Exit(1)
	}

	if err := (&infrav1.AzureClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AzureClusterTemplate")
		os.Exit(1)
	}

	if err := infrav1.SetupAzureMachineWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AzureMachine")
		os.Exit(1)
	}

	if err := (&infrav1.AzureMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AzureMachineTemplate")
		os.Exit(1)
	}

	if err := (&infrav1.AzureClusterIdentity{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AzureClusterIdentity")
		os.Exit(1)
	}

	if feature.Gates.Enabled(capifeature.MachinePool) {
		if err := infrav1exp.SetupAzureMachinePoolWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureMachinePool")
			os.Exit(1)
		}

		if err := (&infrav1exp.AzureMachinePoolMachine{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureMachinePoolMachine")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureManagedMachinePoolWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedMachinePool")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureManagedMachinePoolTemplateWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedMachinePoolTemplate")
			os.Exit(1)
		}

		if err := (&infrav1.AzureManagedCluster{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedCluster")
			os.Exit(1)
		}

		if err := (&infrav1.AzureManagedClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedClusterTemplate")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureManagedControlPlaneWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedControlPlane")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureManagedControlPlaneTemplateWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureManagedControlPlaneTemplate")
			os.Exit(1)
		}
	}

	if feature.Gates.Enabled(feature.ASOAPI) {
		if err := infrav1.SetupAzureASOManagedClusterWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureASOManagedCluster")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureASOManagedControlPlaneWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureASOManagedControlPlane")
			os.Exit(1)
		}

		if err := infrav1.SetupAzureASOManagedMachinePoolWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AzureASOManagedMachinePool")
			os.Exit(1)
		}
	}

	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}
}
