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
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	_ "k8s.io/component-base/logs/json/register"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	cache "sigs.k8s.io/controller-runtime/pkg/cache"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1alpha6 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha6"
	infrav1alpha7 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/controllers"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/webhooks"
	"sigs.k8s.io/cluster-api-provider-openstack/version"
)

// Constants for TLS versions.
const (
	TLSVersion12 = "TLS12"
	TLSVersion13 = "TLS13"
)

type TLSOptions struct {
	TLSMaxVersion   string
	TLSMinVersion   string
	TLSCipherSuites string
}

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	// flags.
	managerOptions              = flags.ManagerOptions{}
	enableLeaderElection        bool
	leaderElectionLeaseDuration time.Duration
	leaderElectionRenewDeadline time.Duration
	leaderElectionRetryPeriod   time.Duration
	watchNamespace              string
	watchFilterValue            string
	profilerAddress             string
	openStackClusterConcurrency int
	openStackMachineConcurrency int
	syncPeriod                  time.Duration
	restConfigQPS               float32
	restConfigBurst             int
	webhookPort                 int
	webhookCertDir              string
	healthAddr                  string
	lbProvider                  string
	caCertsPath                 string
	showVersion                 bool
	scopeCacheMaxSize           int
	logOptions                  = logs.NewOptions()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = ipamv1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = infrav1alpha6.AddToScheme(scheme)
	_ = infrav1alpha7.AddToScheme(scheme)
	_ = infrav1alpha1.AddToScheme(scheme)
	_ = orcv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme

	metrics.RegisterAPIPrometheusMetrics()
}

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	logs.AddFlags(fs, logs.SkipLoggingConfigurationFlags())
	logsv1.AddFlags(logOptions, fs)

	flags.AddManagerOptions(fs, &managerOptions)

	fs.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	fs.DurationVar(&leaderElectionLeaseDuration, "leader-elect-lease-duration", 15*time.Second,
		"Interval at which non-leader candidates will wait to force acquire leadership (duration string)")

	fs.DurationVar(&leaderElectionRenewDeadline, "leader-elect-renew-deadline", 10*time.Second,
		"Duration that the leading controller manager will retry refreshing leadership before giving up (duration string)")

	fs.DurationVar(&leaderElectionRetryPeriod, "leader-elect-retry-period", 2*time.Second,
		"Duration the LeaderElector clients should wait between tries of actions (duration string)")

	fs.StringVar(&watchNamespace, "namespace", "",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")

	fs.StringVar(&watchFilterValue, "watch-filter", "",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel))

	fs.StringVar(&profilerAddress, "profiler-address", "",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")

	fs.IntVar(&openStackClusterConcurrency, "openstackcluster-concurrency", 10,
		"Number of OpenStackClusters to process simultaneously")

	fs.IntVar(&openStackMachineConcurrency, "openstackmachine-concurrency", 10,
		"Number of OpenStackMachines to process simultaneously")

	fs.DurationVar(&syncPeriod, "sync-period", 10*time.Minute,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)")

	fs.Float32Var(&restConfigQPS, "kube-api-qps", 20,
		"Maximum queries per second from the controller client to the Kubernetes API server. Defaults to 20")

	fs.IntVar(&restConfigBurst, "kube-api-burst", 30,
		"Maximum number of queries that should be allowed in one burst from the controller client to the Kubernetes API server. Defaults to 30")

	fs.IntVar(&webhookPort, "webhook-port", 9443,
		"Webhook Server port")

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")

	fs.StringVar(&healthAddr, "health-addr", ":9440",
		"The address the health endpoint binds to.")

	fs.StringVar(&lbProvider, "lb-provider", "amphora",
		"The name of the load balancer provider (amphora or ovn) to use (defaults to amphora).")

	fs.StringVar(&caCertsPath, "ca-certs", "", "The path to a PEM-encoded CA Certificate file to supply as default for each request.")

	fs.IntVar(&scopeCacheMaxSize, "scope-cache-max-size", 10, "The maximum credentials count the operator should keep in cache. Setting this value to 0 means no cache.")

	fs.BoolVar(&showVersion, "version", false, "Show current version and exit.")
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if showVersion {
		fmt.Println(version.Get().String()) //nolint:forbidigo
		os.Exit(0)
	}

	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	if profilerAddress != "" {
		klog.Infof("Profiler listening for requests at %s", profilerAddress)
		go func() {
			klog.Info(http.ListenAndServe(profilerAddress, nil)) //nolint:gosec
		}()
	}

	cfg, err := config.GetConfigWithContext(os.Getenv("KUBECONTEXT"))
	if err != nil {
		setupLog.Error(err, "unable to get kubeconfig")
		os.Exit(1)
	}
	cfg.QPS = restConfigQPS
	cfg.Burst = restConfigBurst

	var caCerts []byte
	if caCertsPath != "" {
		caCerts, err = os.ReadFile(caCertsPath)
		if err != nil {
			setupLog.Error(err, "unable to read provided ca certificates file")
			os.Exit(1)
		}
	}

	tlsOpts, metricsOpts, err := flags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "unable to get manager options")
		os.Exit(1)
	}

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		Metrics:          *metricsOpts,
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "controller-leader-election-capo",
		LeaseDuration:    &leaderElectionLeaseDuration,
		RenewDeadline:    &leaderElectionRenewDeadline,
		RetryPeriod:      &leaderElectionRetryPeriod,
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
		WebhookServer: webhook.NewServer(
			webhook.Options{
				Port:    webhookPort,
				CertDir: webhookCertDir,
				TLSOpts: tlsOpts,
			},
		),
		HealthProbeBindAddress:        healthAddr,
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("openstack-controller"))

	scopeFactory := scope.NewFactory(scopeCacheMaxSize)

	setupChecks(mgr)
	setupReconcilers(ctx, mgr, caCerts, scopeFactory)
	setupWebhooks(mgr)
	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupChecks(mgr ctrl.Manager) {
	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}
}

func setupReconcilers(ctx context.Context, mgr ctrl.Manager, caCerts []byte, scopeFactory scope.Factory) {
	if err := (&controllers.OpenStackClusterReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("openstackcluster-controller"),
		WatchFilterValue: watchFilterValue,
		ScopeFactory:     scopeFactory,
		CaCertificates:   caCerts,
	}).SetupWithManager(ctx, mgr, concurrency(openStackClusterConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpenStackCluster")
		os.Exit(1)
	}
	if err := (&controllers.OpenStackMachineReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("openstackmachine-controller"),
		WatchFilterValue: watchFilterValue,
		ScopeFactory:     scopeFactory,
		CaCertificates:   caCerts,
	}).SetupWithManager(ctx, mgr, concurrency(openStackMachineConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpenStackMachine")
		os.Exit(1)
	}
	if err := (&controllers.OpenStackFloatingIPPoolReconciler{
		Client:         mgr.GetClient(),
		Recorder:       mgr.GetEventRecorderFor("floatingippool-controller"),
		ScopeFactory:   scopeFactory,
		Scheme:         mgr.GetScheme(),
		CaCertificates: caCerts,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FloatingIPPool")
		os.Exit(1)
	}
	if err := (&controllers.OpenStackServerReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("openstackserver-controller"),
		WatchFilterValue: watchFilterValue,
		ScopeFactory:     scopeFactory,
		CaCertificates:   caCerts,
		Scheme:           mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, concurrency(openStackMachineConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpenStackServer")
		os.Exit(1)
	}
	if err := controllers.ImageController(
		mgr.GetClient(),
		mgr.GetEventRecorderFor("orc-image-controller"),
		watchFilterValue,
		scopeFactory,
		caCerts,
	).SetupWithManager(ctx, mgr, concurrency(openStackMachineConcurrency)); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ORCImage")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	errs := webhooks.RegisterAllWithManager(mgr)
	if len(errs) > 0 {
		for i := range errs {
			setupLog.Error(errs[i], "unable to register webhook")
		}
		os.Exit(1)
	}
}

func concurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}
