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

// main is the main package for the Cluster API GCP Provider.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	// +kubebuilder:scaffold:imports
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	expcontrollers "sigs.k8s.io/cluster-api-provider-gcp/exp/controllers"
	"sigs.k8s.io/cluster-api-provider-gcp/feature"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-gcp/version"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/flags"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme         = runtime.NewScheme()
	setupLog       = ctrl.Log.WithName("setup")
	managerOptions = flags.ManagerOptions{}
)

func init() {
	klog.InitFlags(nil)

	_ = clientgoscheme.AddToScheme(scheme)
	_ = infrav1beta1.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = expclusterv1.AddToScheme(scheme)
	_ = infrav1exp.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	enableLeaderElection        bool
	leaderElectionNamespace     string
	watchNamespace              string
	profilerAddress             string
	healthAddr                  string
	watchFilterValue            string
	webhookCertDir              string
	gcpClusterConcurrency       int
	gcpMachineConcurrency       int
	webhookPort                 int
	reconcileTimeout            time.Duration
	syncPeriod                  time.Duration
	leaderElectionLeaseDuration time.Duration
	leaderElectionRenewDeadline time.Duration
	leaderElectionRetryPeriod   time.Duration
)

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

func main() {
	initFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	_, metricsOptions, err := flags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: invalid flags")
	}

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	if profilerAddress != "" {
		setupLog.Info("Profiler listening for requests", "profiler-address", profilerAddress)
		go func() {
			server := &http.Server{
				Addr:              profilerAddress,
				ReadHeaderTimeout: 3 * time.Second,
			}

			err := server.ListenAndServe()
			if err != nil {
				setupLog.Error(err, "listen and serve error")
			}
		}()
	}

	ctrl.SetLogger(klog.Background())

	setupLog.Info(fmt.Sprintf("feature gates: %+v\n", feature.Gates))

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		Metrics:                 *metricsOptions,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "controller-leader-election-capg",
		LeaderElectionNamespace: leaderElectionNamespace,
		LeaseDuration:           &leaderElectionLeaseDuration,
		RenewDeadline:           &leaderElectionRenewDeadline,
		RetryPeriod:             &leaderElectionRetryPeriod,
		Cache: cache.Options{
			DefaultNamespaces: watchNamespaces,
			SyncPeriod:        &syncPeriod,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    webhookPort,
			CertDir: webhookCertDir,
		}),
		HealthProbeBindAddress: healthAddr,
		EventBroadcaster:       broadcaster,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("gcp-controller"))

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	if err := setupReconcilers(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to setup reconcilers")
		os.Exit(1)
	}

	if err := setupWebhooks(mgr); err != nil {
		setupLog.Error(err, "unable to setup webhooks")
		os.Exit(1)
	}

	if err := setupProbes(mgr); err != nil {
		setupLog.Error(err, "unable to setup probes")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager", "version", version.Get().String(), "extended_info", version.Get())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupReconcilers(ctx context.Context, mgr ctrl.Manager) error {
	if err := (&controllers.GCPMachineReconciler{
		Client:           mgr.GetClient(),
		ReconcileTimeout: reconcileTimeout,
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: gcpMachineConcurrency}); err != nil {
		return fmt.Errorf("setting up GCPMachine controller: %w", err)
	}
	if err := (&controllers.GCPClusterReconciler{
		Client:           mgr.GetClient(),
		ReconcileTimeout: reconcileTimeout,
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: gcpClusterConcurrency}); err != nil {
		return fmt.Errorf("setting up GCPCluster controller: %w", err)
	}

	if feature.Gates.Enabled(feature.GKE) {
		setupLog.Info("Enabling GKE reconcilers")

		if err := (&expcontrollers.GCPManagedClusterReconciler{
			Client:           mgr.GetClient(),
			ReconcileTimeout: reconcileTimeout,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: gcpClusterConcurrency}); err != nil {
			return fmt.Errorf("setting up GCPManagedCluster controller: %w", err)
		}

		if err := (&expcontrollers.GCPManagedControlPlaneReconciler{
			Client:           mgr.GetClient(),
			ReconcileTimeout: reconcileTimeout,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: gcpClusterConcurrency}); err != nil {
			return fmt.Errorf("setting up GCPManagedControlPlane controller: %w", err)
		}

		if err := (&expcontrollers.GCPManagedMachinePoolReconciler{
			Client:           mgr.GetClient(),
			ReconcileTimeout: reconcileTimeout,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: gcpMachineConcurrency}); err != nil {
			return fmt.Errorf("setting up GCPManagedMachinePool controller: %w", err)
		}
	}

	return nil
}

func setupWebhooks(mgr ctrl.Manager) error {
	if err := (&infrav1beta1.GCPCluster{}).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("setting up GCPCluster webhook: %w", err)
	}
	if err := (&infrav1beta1.GCPClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("setting up GCPClusterTemplate webhook: %w", err)
	}
	if err := (&infrav1beta1.GCPMachine{}).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("setting up GCPMachine webhook: %w", err)
	}
	if err := (&infrav1beta1.GCPMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		return fmt.Errorf("setting up GCPMachineTemplate webhook: %w", err)
	}

	if feature.Gates.Enabled(feature.GKE) {
		setupLog.Info("Enabling GKE webhooks")

		if err := (&infrav1exp.GCPManagedCluster{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("setting up GCPManagedCluster webhook: %w", err)
		}
		if err := (&infrav1exp.GCPManagedControlPlane{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("setting up GCPManagedControlPlane webhook: %w", err)
		}
		if err := (&infrav1exp.GCPManagedMachinePool{}).SetupWebhookWithManager(mgr); err != nil {
			return fmt.Errorf("setting up GCPManagedMachinePool webhook: %w", err)
		}
	}

	return nil
}

func setupProbes(mgr ctrl.Manager) error {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		return fmt.Errorf("creating ready check: %w", err)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		return fmt.Errorf("creating health check: %w", err)
	}

	return nil
}

func initFlags(fs *pflag.FlagSet) {
	fs.BoolVar(
		&enableLeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
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
		&leaderElectionNamespace,
		"leader-election-namespace",
		"",
		"Namespace that the controller performs leader election in. If unspecified, the controller will discover which namespace it is running in.",
	)

	fs.StringVar(
		&profilerAddress,
		"profiler-address",
		"",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)",
	)

	fs.StringVar(
		&watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel),
	)

	fs.IntVar(&gcpClusterConcurrency,
		"gcpcluster-concurrency",
		10,
		"Number of GCPClusters to process simultaneously",
	)

	fs.IntVar(&gcpMachineConcurrency,
		"gcpmachine-concurrency",
		10,
		"Number of GCPMachines to process simultaneously",
	)

	fs.DurationVar(&syncPeriod,
		"sync-period",
		10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)",
	)

	fs.IntVar(&webhookPort,
		"webhook-port",
		9443,
		"Webhook Server port",
	)

	fs.StringVar(&webhookCertDir,
		"webhook-cert-dir",
		"/tmp/k8s-webhook-server/serving-certs",
		"Webhook Server Certificate Directory, is the directory that contains the server key and certificate",
	)

	fs.StringVar(&healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	fs.DurationVar(&reconcileTimeout,
		"reconcile-timeout",
		reconciler.DefaultLoopTimeout,
		"The maximum duration a reconcile loop can run (e.g. 90m)",
	)

	flags.AddManagerOptions(fs, &managerOptions)

	feature.MutableGates.AddFlag(fs)
}
