/*
Copyright 2022 Nutanix

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
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	capiflags "sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	"github.com/nutanix-cloud-native/cluster-api-provider-nutanix/controllers"
	//+kubebuilder:scaffold:imports
)

var scheme = runtime.NewScheme()

// gitCommitHash is the git commit hash of the code that is running.
var gitCommitHash string

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(capiv1.AddToScheme(scheme))
	utilruntime.Must(bootstrapv1.AddToScheme(scheme))
	utilruntime.Must(infrav1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

const (
	// DefaultMaxConcurrentReconciles is the default maximum number of concurrent reconciles
	defaultMaxConcurrentReconciles = 10
)

type options struct {
	enableLeaderElection    bool
	healthProbeAddr         string
	maxConcurrentReconciles int

	rateLimiterBaseDelay  time.Duration
	rateLimiterMaxDelay   time.Duration
	rateLimiterBucketSize int
	rateLimiterQPS        int

	managerOptions capiflags.ManagerOptions
	zapOptions     zap.Options
}

type managerConfig struct {
	enableLeaderElection               bool
	healthProbeAddr                    string
	concurrentReconcilesNutanixCluster int
	concurrentReconcilesNutanixMachine int
	metricsServerOpts                  server.Options
	skipNameValidation                 bool

	logger      logr.Logger
	restConfig  *rest.Config
	rateLimiter workqueue.TypedRateLimiter[reconcile.Request]
}

// compositeRateLimiter will build a limiter similar to the default from DefaultControllerRateLimiter but with custom values.
func compositeRateLimiter(baseDelay, maxDelay time.Duration, bucketSize, qps int) (workqueue.TypedRateLimiter[reconcile.Request], error) {
	// Validate the rate limiter configuration
	if err := validateRateLimiterConfig(baseDelay, maxDelay, bucketSize, qps); err != nil {
		return nil, err
	}
	exponentialBackoffLimiter := workqueue.NewTypedItemExponentialFailureRateLimiter[reconcile.Request](baseDelay, maxDelay)
	bucketLimiter := &workqueue.TypedBucketRateLimiter[reconcile.Request]{Limiter: rate.NewLimiter(rate.Limit(qps), bucketSize)}
	return workqueue.NewTypedMaxOfRateLimiter[reconcile.Request](exponentialBackoffLimiter, bucketLimiter), nil
}

// validateRateLimiterConfig validates the rate limiter configuration parameters
func validateRateLimiterConfig(baseDelay, maxDelay time.Duration, bucketSize, qps int) error {
	// Check if baseDelay is a non-negative value
	if baseDelay < 0 {
		return errors.New("baseDelay cannot be negative")
	}

	// Check if maxDelay is non-negative and greater than or equal to baseDelay
	if maxDelay < 0 {
		return errors.New("maxDelay cannot be negative")
	}

	if maxDelay < baseDelay {
		return errors.New("maxDelay should be greater than or equal to baseDelay")
	}

	// Check if bucketSize is a positive number
	if bucketSize <= 0 {
		return errors.New("bucketSize must be positive")
	}

	// Check if qps is a positive number
	if qps <= 0 {
		return errors.New("minimum QPS must be positive")
	}

	// Check if bucketSize is at least as large as the QPS
	if bucketSize < qps {
		return errors.New("bucketSize must be at least as large as the QPS to handle bursts effectively")
	}

	return nil
}

func initializeFlags() *options {
	opts := &options{}

	// Add the controller-runtime flags to the standard library FlagSet.
	ctrl.RegisterFlags(flag.CommandLine)

	// Add the Cluster API flags to the pflag FlagSet.
	capiflags.AddManagerOptions(pflag.CommandLine, &opts.managerOptions)

	// Add zap flags to the standard libary FlagSet.
	opts.zapOptions.BindFlags(flag.CommandLine)

	// Add our own flags to the pflag FlagSet.
	pflag.StringVar(&opts.healthProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	pflag.BoolVar(&opts.enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")

	pflag.IntVar(&opts.maxConcurrentReconciles, "max-concurrent-reconciles", defaultMaxConcurrentReconciles,
		"The maximum number of allowed, concurrent reconciles.")

	pflag.DurationVar(&opts.rateLimiterBaseDelay, "rate-limiter-base-delay", 500*time.Millisecond, "The base delay for the rate limiter.")
	pflag.DurationVar(&opts.rateLimiterMaxDelay, "rate-limiter-max-delay", 15*time.Minute, "The maximum delay for the rate limiter.")
	pflag.IntVar(&opts.rateLimiterBucketSize, "rate-limiter-bucket-size", 100, "The bucket size for the rate limiter.")
	pflag.IntVar(&opts.rateLimiterQPS, "rate-limiter-qps", 10, "The QPS for the rate limiter.")

	// At this point, we should be done adding flags to the standard library FlagSet, flag.CommandLine.
	// So we can include the flags that third-party libraries, e.g. controller-runtime, and zap,
	// have added to the standard library FlagSet, we merge it into the pflag FlagSet.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Parse flags.
	pflag.Parse()

	return opts
}

func initializeConfig(opts *options) (*managerConfig, error) {
	config := &managerConfig{
		enableLeaderElection: opts.enableLeaderElection,
		healthProbeAddr:      opts.healthProbeAddr,
	}

	_, metricsServerOpts, err := capiflags.GetManagerOptions(opts.managerOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get metrics server options: %w", err)
	}
	if metricsServerOpts == nil {
		return nil, errors.New("parsed manager options are nil")
	}
	metricsServerOpts.FilterProvider = filters.WithAuthenticationAndAuthorization
	config.metricsServerOpts = *metricsServerOpts

	config.concurrentReconcilesNutanixCluster = opts.maxConcurrentReconciles
	config.concurrentReconcilesNutanixMachine = opts.maxConcurrentReconciles

	rateLimiter, err := compositeRateLimiter(opts.rateLimiterBaseDelay, opts.rateLimiterMaxDelay, opts.rateLimiterBucketSize, opts.rateLimiterQPS)
	if err != nil {
		return nil, fmt.Errorf("unable to create composite rate limiter: %w", err)
	}
	config.rateLimiter = rateLimiter

	zapOptions := opts.zapOptions
	zapOptions.TimeEncoder = zapcore.RFC3339TimeEncoder
	config.logger = zap.New(zap.UseFlagOptions(&zapOptions))

	// Configure controller-runtime logger before using calling any controller-runtime functions.
	// Otherwise, the user will not see warnings and errors logged by these functions.
	ctrl.SetLogger(config.logger)

	// Before calling GetConfigOrDie, we have parsed flags, because the function reads value of
	// the--kubeconfig flag.
	config.restConfig, err = ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return config, nil
}

func setupLogger() logr.Logger {
	return ctrl.Log.WithName("setup")
}

func addHealthChecks(mgr manager.Manager) error {
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	return nil
}

func createInformers(ctx context.Context, mgr manager.Manager) (coreinformers.SecretInformer, coreinformers.ConfigMapInformer, error) {
	// Create a secret informer for the Nutanix client
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create clientset for management cluster: %w", err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute)
	secretInformer := informerFactory.Core().V1().Secrets()
	informer := secretInformer.Informer()
	go informer.Run(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), informer.HasSynced)

	configMapInformer := informerFactory.Core().V1().ConfigMaps()
	cmInformer := configMapInformer.Informer()
	go cmInformer.Run(ctx.Done())
	cache.WaitForCacheSync(ctx.Done(), cmInformer.HasSynced)

	return secretInformer, configMapInformer, nil
}

func setupNutanixClusterController(ctx context.Context, mgr manager.Manager, secretInformer coreinformers.SecretInformer,
	configMapInformer coreinformers.ConfigMapInformer, opts ...controllers.ControllerConfigOpts,
) error {
	clusterCtrl, err := controllers.NewNutanixClusterReconciler(
		mgr.GetClient(),
		secretInformer,
		configMapInformer,
		mgr.GetScheme(),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("unable to create NutanixCluster controller: %w", err)
	}

	if err := clusterCtrl.SetupWithManager(ctx, mgr); err != nil {
		return fmt.Errorf("unable to setup NutanixCluster controller with manager: %w", err)
	}

	return nil
}

func setupNutanixMachineController(ctx context.Context, mgr manager.Manager, secretInformer coreinformers.SecretInformer,
	configMapInformer coreinformers.ConfigMapInformer, opts ...controllers.ControllerConfigOpts,
) error {
	machineCtrl, err := controllers.NewNutanixMachineReconciler(
		mgr.GetClient(),
		secretInformer,
		configMapInformer,
		mgr.GetScheme(),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("unable to create NutanixMachine controller: %w", err)
	}

	if err := machineCtrl.SetupWithManager(ctx, mgr); err != nil {
		return fmt.Errorf("unable to setup NutanixMachine controller with manager: %w", err)
	}

	return nil
}

func setupNutanixFailureDomainController(ctx context.Context, mgr manager.Manager, secretInformer coreinformers.SecretInformer,
	configMapInformer coreinformers.ConfigMapInformer, opts ...controllers.ControllerConfigOpts,
) error {
	machineCtrl, err := controllers.NewNutanixFailureDomainReconciler(
		mgr.GetClient(),
		secretInformer,
		configMapInformer,
		mgr.GetScheme(),
		opts...,
	)
	if err != nil {
		return fmt.Errorf("unable to create NutanixFailureDomain controller: %w", err)
	}

	if err := machineCtrl.SetupWithManager(ctx, mgr); err != nil {
		return fmt.Errorf("unable to setup NutanixFailureDomain controller with manager: %w", err)
	}

	return nil
}

func runManager(ctx context.Context, mgr manager.Manager, config *managerConfig) error {
	secretInformer, configMapInformer, err := createInformers(ctx, mgr)
	if err != nil {
		return fmt.Errorf("unable to create informers: %w", err)
	}

	clusterControllerOpts := []controllers.ControllerConfigOpts{
		controllers.WithMaxConcurrentReconciles(config.concurrentReconcilesNutanixCluster),
		controllers.WithRateLimiter(workqueue.NewTypedMaxOfRateLimiter(workqueue.NewTypedItemExponentialFailureRateLimiter[reconcile.Request](1*time.Millisecond, 1000*time.Second), &workqueue.TypedBucketRateLimiter[reconcile.Request]{Limiter: rate.NewLimiter(rate.Limit(10), 100)})),
	}

	// Enable SkipNameValidation in test environments
	if config.skipNameValidation {
		clusterControllerOpts = append(clusterControllerOpts, controllers.WithSkipNameValidation(true))
	}

	if err := setupNutanixClusterController(ctx, mgr, secretInformer, configMapInformer, clusterControllerOpts...); err != nil {
		return fmt.Errorf("unable to setup controllers: %w", err)
	}

	machineControllerOpts := []controllers.ControllerConfigOpts{
		controllers.WithMaxConcurrentReconciles(config.concurrentReconcilesNutanixMachine),
		controllers.WithRateLimiter(workqueue.NewTypedMaxOfRateLimiter(workqueue.NewTypedItemExponentialFailureRateLimiter[reconcile.Request](1*time.Millisecond, 1000*time.Second), &workqueue.TypedBucketRateLimiter[reconcile.Request]{Limiter: rate.NewLimiter(rate.Limit(10), 100)})),
	}

	// Enable SkipNameValidation in test environments for machine controllers
	if config.skipNameValidation {
		machineControllerOpts = append(machineControllerOpts, controllers.WithSkipNameValidation(true))
	}

	if err := setupNutanixMachineController(ctx, mgr, secretInformer, configMapInformer, machineControllerOpts...); err != nil {
		return fmt.Errorf("unable to setup controllers: %w", err)
	}

	// Use the same opts for failure domain controller as machine controller
	if err := setupNutanixFailureDomainController(ctx, mgr, secretInformer, configMapInformer, machineControllerOpts...); err != nil {
		return fmt.Errorf("unable to setup controllers: %w", err)
	}

	config.logger.Info("starting CAPX Controller Manager")
	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}

func initializeManager(config *managerConfig) (manager.Manager, error) {
	mgr, err := ctrl.NewManager(config.restConfig, ctrl.Options{
		Scheme:                 scheme,
		Metrics:                config.metricsServerOpts,
		HealthProbeBindAddress: config.healthProbeAddr,
		LeaderElection:         config.enableLeaderElection,
		LeaderElectionID:       "f265110d.cluster.x-k8s.io",
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create manager: %w", err)
	}

	if err := addHealthChecks(mgr); err != nil {
		return nil, fmt.Errorf("unable to add health checks to manager: %w", err)
	}

	return mgr, nil
}

func main() {
	logger := setupLogger()

	logger.Info("Initializing Nutanix Cluster API Infrastructure Provider", "Git Hash", gitCommitHash)

	opts := initializeFlags()
	// After this point, we must not add flags to either the pflag, or the standard library FlagSets.

	config, err := initializeConfig(opts)
	if err != nil {
		logger.Error(err, "unable to configure manager")
		os.Exit(1)
	}

	mgr, err := initializeManager(config)
	if err != nil {
		logger.Error(err, "unable to create manager")
		os.Exit(1)
	}

	// Set up the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()
	if err := runManager(ctx, mgr, config); err != nil {
		logger.Error(err, "problem running manager")
		os.Exit(1)
	}
}
