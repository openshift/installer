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

type managerConfig struct {
	enableLeaderElection               bool
	probeAddr                          string
	concurrentReconcilesNutanixCluster int
	concurrentReconcilesNutanixMachine int
	managerOptions                     capiflags.ManagerOptions

	logger      logr.Logger
	restConfig  *rest.Config
	rateLimiter workqueue.RateLimiter
}

// compositeRateLimiter will build a limiter similar to the default from DefaultControllerRateLimiter but with custom values.
func compositeRateLimiter(baseDelay, maxDelay time.Duration, bucketSize, qps int) (workqueue.RateLimiter, error) {
	// Validate the rate limiter configuration
	if err := validateRateLimiterConfig(baseDelay, maxDelay, bucketSize, qps); err != nil {
		return nil, err
	}
	exponentialBackoffLimiter := workqueue.NewItemExponentialFailureRateLimiter(baseDelay, maxDelay)
	bucketLimiter := &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(qps), bucketSize)}
	return workqueue.NewMaxOfRateLimiter(exponentialBackoffLimiter, bucketLimiter), nil
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

func parseFlags(config *managerConfig) {
	capiflags.AddManagerOptions(pflag.CommandLine, &config.managerOptions)
	pflag.StringVar(&config.probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	pflag.BoolVar(&config.enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	var maxConcurrentReconciles int
	pflag.IntVar(&maxConcurrentReconciles, "max-concurrent-reconciles", defaultMaxConcurrentReconciles,
		"The maximum number of allowed, concurrent reconciles.")

	var baseDelay, maxDelay time.Duration
	var bucketSize, qps int
	pflag.DurationVar(&baseDelay, "rate-limiter-base-delay", 500*time.Millisecond, "The base delay for the rate limiter.")
	pflag.DurationVar(&maxDelay, "rate-limiter-max-delay", 15*time.Minute, "The maximum delay for the rate limiter.")
	pflag.IntVar(&bucketSize, "rate-limiter-bucket-size", 100, "The bucket size for the rate limiter.")
	pflag.IntVar(&qps, "rate-limiter-qps", 10, "The QPS for the rate limiter.")

	opts := zap.Options{
		TimeEncoder: zapcore.RFC3339TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)

	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	config.concurrentReconcilesNutanixCluster = maxConcurrentReconciles
	config.concurrentReconcilesNutanixMachine = maxConcurrentReconciles

	rateLimiter, err := compositeRateLimiter(baseDelay, maxDelay, bucketSize, qps)
	if err != nil {
		config.logger.Error(err, "unable to create composite rate limiter")
		os.Exit(1)
	}

	config.rateLimiter = rateLimiter
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

func runManager(ctx context.Context, mgr manager.Manager, config *managerConfig) error {
	secretInformer, configMapInformer, err := createInformers(ctx, mgr)
	if err != nil {
		return fmt.Errorf("unable to create informers: %w", err)
	}

	clusterControllerOpts := []controllers.ControllerConfigOpts{
		controllers.WithMaxConcurrentReconciles(config.concurrentReconcilesNutanixCluster),
		controllers.WithRateLimiter(config.rateLimiter),
	}

	if err := setupNutanixClusterController(ctx, mgr, secretInformer, configMapInformer, clusterControllerOpts...); err != nil {
		return fmt.Errorf("unable to setup controllers: %w", err)
	}

	machineControllerOpts := []controllers.ControllerConfigOpts{
		controllers.WithMaxConcurrentReconciles(config.concurrentReconcilesNutanixMachine),
		controllers.WithRateLimiter(config.rateLimiter),
	}

	if err := setupNutanixMachineController(ctx, mgr, secretInformer, configMapInformer, machineControllerOpts...); err != nil {
		return fmt.Errorf("unable to setup controllers: %w", err)
	}

	config.logger.Info("starting CAPX Controller Manager")
	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}

func initializeManager(config *managerConfig) (manager.Manager, error) {
	_, metricsOpts, err := capiflags.GetManagerOptions(config.managerOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to get manager options: %w", err)
	}

	if metricsOpts == nil {
		return nil, errors.New("parsed manager options are nil")
	}

	mgr, err := ctrl.NewManager(config.restConfig, ctrl.Options{
		Scheme:                 scheme,
		Metrics:                *metricsOpts,
		HealthProbeBindAddress: config.probeAddr,
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

	config := &managerConfig{}
	parseFlags(config)

	// Flags must be parsed before calling GetConfigOrDie, because
	// it reads the value of the--kubeconfig flag.
	config.restConfig = ctrl.GetConfigOrDie()

	config.logger = logger

	logger.Info("Initializing Nutanix Cluster API Infrastructure Provider", "Git Hash", gitCommitHash)
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
