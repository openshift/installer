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
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	cgscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	// +kubebuilder:scaffold:imports
	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1alpha4 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	eksbootstrapv1alpha3 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha3"
	eksbootstrapv1alpha4 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1alpha4"
	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1beta1"
	eksbootstrapcontrollers "sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/controllers"
	"sigs.k8s.io/cluster-api-provider-aws/controllers"
	ekscontrolplanev1alpha3 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	ekscontrolplanev1alpha4 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha4"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	ekscontrolplanecontrollers "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/controllers"
	expinfrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	expinfrav1alpha4 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/exp/controlleridentitycreator"
	expcontrollers "sigs.k8s.io/cluster-api-provider-aws/exp/controllers"
	"sigs.k8s.io/cluster-api-provider-aws/exp/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/feature"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/version"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = eksbootstrapv1.AddToScheme(scheme)
	_ = eksbootstrapv1alpha3.AddToScheme(scheme)
	_ = eksbootstrapv1alpha4.AddToScheme(scheme)
	_ = cgscheme.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	_ = expclusterv1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	_ = ekscontrolplanev1alpha3.AddToScheme(scheme)
	_ = ekscontrolplanev1alpha4.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	_ = expinfrav1alpha3.AddToScheme(scheme)
	_ = infrav1alpha4.AddToScheme(scheme)
	_ = expinfrav1alpha4.AddToScheme(scheme)
	_ = expinfrav1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

var (
	metricsBindAddr          string
	enableLeaderElection     bool
	leaderElectionNamespace  string
	watchNamespace           string
	watchFilterValue         string
	profilerAddress          string
	awsClusterConcurrency    int
	instanceStateConcurrency int
	awsMachineConcurrency    int
	syncPeriod               time.Duration
	webhookPort              int
	webhookCertDir           string
	healthAddr               string
	serviceEndpoints         string

	// maxEKSSyncPeriod is the maximum allowed duration for the sync-period flag when using EKS. It is set to 10 minutes
	// because during resync it will create a new AWS auth token which can a maximum life of 15 minutes and this ensures
	// the token (and kubeconfig secret) is refreshed before token expiration.
	maxEKSSyncPeriod         = time.Minute * 10
	errMaxSyncPeriodExceeded = errors.New("sync period greater than maximum allowed")
	errEKSInvalidFlags       = errors.New("invalid EKS flag combination")
)

func main() {
	klog.InitFlags(nil)

	rand.Seed(time.Now().UnixNano())
	initFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())

	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	if profilerAddress != "" {
		setupLog.Info("Profiler listening for requests", "profiler-address", profilerAddress)
		go func() {
			setupLog.Error(http.ListenAndServe(profilerAddress, nil), "listen and serve error")
		}()
	}

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	ctx := ctrl.SetupSignalHandler()

	restConfig := ctrl.GetConfigOrDie()
	restConfig.UserAgent = "cluster-api-provider-aws-controller"
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                     scheme,
		MetricsBindAddress:         metricsBindAddr,
		LeaderElection:             enableLeaderElection,
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaderElectionID:           "controller-leader-elect-capa",
		LeaderElectionNamespace:    leaderElectionNamespace,
		SyncPeriod:                 &syncPeriod,
		Namespace:                  watchNamespace,
		EventBroadcaster:           broadcaster,
		Port:                       webhookPort,
		CertDir:                    webhookCertDir,
		HealthProbeBindAddress:     healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("aws-controller"))

	setupLog.V(1).Info(fmt.Sprintf("feature gates: %+v\n", feature.Gates))

	externalResourceGC := false
	if feature.Gates.Enabled(feature.ExternalResourceGC) {
		setupLog.Info("enabling external resource garbage collection")
		externalResourceGC = true
	}

	// Parse service endpoints.
	AWSServiceEndpoints, err := endpoints.ParseFlag(serviceEndpoints)
	if err != nil {
		setupLog.Error(err, "unable to parse service endpoints", "controller", "AWSCluster")
		os.Exit(1)
	}

	if err = (&controllers.AWSMachineReconciler{
		Client:           mgr.GetClient(),
		Log:              ctrl.Log.WithName("controllers").WithName("AWSMachine"),
		Recorder:         mgr.GetEventRecorderFor("awsmachine-controller"),
		Endpoints:        AWSServiceEndpoints,
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsMachineConcurrency, RecoverPanic: true}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AWSMachine")
		os.Exit(1)
	}
	if err = (&controllers.AWSClusterReconciler{
		Client:             mgr.GetClient(),
		Recorder:           mgr.GetEventRecorderFor("awscluster-controller"),
		Endpoints:          AWSServiceEndpoints,
		WatchFilterValue:   watchFilterValue,
		ExternalResourceGC: externalResourceGC,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency, RecoverPanic: true}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AWSCluster")
		os.Exit(1)
	}
	enableGates(ctx, mgr, AWSServiceEndpoints, externalResourceGC)

	if err = (&infrav1.AWSMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachineTemplate")
		os.Exit(1)
	}
	if err = (&infrav1.AWSCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSCluster")
		os.Exit(1)
	}
	if err = (&infrav1.AWSClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSClusterTemplate")
		os.Exit(1)
	}
	if err = (&infrav1.AWSClusterControllerIdentity{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSClusterControllerIdentity")
		os.Exit(1)
	}
	if err = (&infrav1.AWSClusterRoleIdentity{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSClusterRoleIdentity")
		os.Exit(1)
	}
	if err = (&infrav1.AWSClusterStaticIdentity{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSClusterStaticIdentity")
		os.Exit(1)
	}
	if err = (&infrav1.AWSMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachine")
		os.Exit(1)
	}
	if feature.Gates.Enabled(feature.EKS) {
		setupLog.Info("enabling EKS webhooks")
		if err := (&ekscontrolplanev1.AWSManagedControlPlane{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSManagedControlPlane")
			os.Exit(1)
		}
		if feature.Gates.Enabled(feature.EKSFargate) {
			if err = (&expinfrav1.AWSFargateProfile{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "AWSFargateProfile")
				os.Exit(1)
			}
		}
		if feature.Gates.Enabled(feature.MachinePool) {
			if err = (&expinfrav1.AWSManagedMachinePool{}).SetupWebhookWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create webhook", "webhook", "AWSManagedMachinePool")
				os.Exit(1)
			}
		}
	}
	if feature.Gates.Enabled(feature.MachinePool) {
		setupLog.Info("enabling webhook for AWSMachinePool")
		if err = (&expinfrav1.AWSMachinePool{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "AWSMachinePool")
			os.Exit(1)
		}
	}

	// +kubebuilder:scaffold:builder

	if err := mgr.AddReadyzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create ready check")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("webhook", mgr.GetWebhookServer().StartedChecker()); err != nil {
		setupLog.Error(err, "unable to create health check")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "version", version.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func enableGates(ctx context.Context, mgr ctrl.Manager, awsServiceEndpoints []scope.ServiceEndpoint, externalResourceGC bool) {
	if feature.Gates.Enabled(feature.EKS) {
		setupLog.Info("enabling EKS controllers")

		if syncPeriod > maxEKSSyncPeriod {
			setupLog.Error(errMaxSyncPeriodExceeded, "failed to enable EKS", "max-sync-period", maxEKSSyncPeriod, "syn-period", syncPeriod)
			os.Exit(1)
		}

		enableIAM := feature.Gates.Enabled(feature.EKSEnableIAM)
		allowAddRoles := feature.Gates.Enabled(feature.EKSAllowAddRoles)
		setupLog.V(2).Info("EKS IAM role creation", "enabled", enableIAM)
		setupLog.V(2).Info("EKS IAM additional roles", "enabled", allowAddRoles)
		if allowAddRoles && !enableIAM {
			setupLog.Error(errEKSInvalidFlags, "cannot use EKSAllowAddRoles flag without EKSEnableIAM")
			os.Exit(1)
		}

		setupLog.V(2).Info("enabling EKS control plane controller")
		if err := (&ekscontrolplanecontrollers.AWSManagedControlPlaneReconciler{
			Client:               mgr.GetClient(),
			EnableIAM:            enableIAM,
			AllowAdditionalRoles: allowAddRoles,
			Endpoints:            awsServiceEndpoints,
			WatchFilterValue:     watchFilterValue,
			ExternalResourceGC:   externalResourceGC,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency, RecoverPanic: true}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSManagedControlPlane")
			os.Exit(1)
		}

		setupLog.V(2).Info("enabling EKS bootstrap controller")
		if err := (&eksbootstrapcontrollers.EKSConfigReconciler{
			Client:           mgr.GetClient(),
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency, RecoverPanic: true}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "EKSConfig")
			os.Exit(1)
		}

		if feature.Gates.Enabled(feature.EKSFargate) {
			setupLog.V(2).Info("enabling EKS fargate profile controller")
			if err := (&expcontrollers.AWSFargateProfileReconciler{
				Client:           mgr.GetClient(),
				Recorder:         mgr.GetEventRecorderFor("awsfargateprofile-reconciler"),
				EnableIAM:        enableIAM,
				Endpoints:        awsServiceEndpoints,
				WatchFilterValue: watchFilterValue,
			}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency, RecoverPanic: true}); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSFargateProfile")
			}
		}

		if feature.Gates.Enabled(feature.MachinePool) {
			setupLog.V(2).Info("enabling EKS managed machine pool controller")
			if err := (&expcontrollers.AWSManagedMachinePoolReconciler{
				AllowAdditionalRoles: allowAddRoles,
				Client:               mgr.GetClient(),
				EnableIAM:            enableIAM,
				Endpoints:            awsServiceEndpoints,
				Recorder:             mgr.GetEventRecorderFor("awsmanagedmachinepool-reconciler"),
				WatchFilterValue:     watchFilterValue,
			}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: instanceStateConcurrency, RecoverPanic: true}); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "AWSManagedMachinePool")
				os.Exit(1)
			}
		}
	}
	if feature.Gates.Enabled(feature.MachinePool) {
		setupLog.V(2).Info("enabling machine pool controller")
		if err := (&expcontrollers.AWSMachinePoolReconciler{
			Client:           mgr.GetClient(),
			Recorder:         mgr.GetEventRecorderFor("awsmachinepool-controller"),
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: instanceStateConcurrency, RecoverPanic: true}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSMachinePool")
			os.Exit(1)
		}
	}
	if feature.Gates.Enabled(feature.EventBridgeInstanceState) {
		setupLog.Info("EventBridge notifications enabled. enabling AWSInstanceStateController")
		if err := (&instancestate.AwsInstanceStateReconciler{
			Client:           mgr.GetClient(),
			Log:              ctrl.Log.WithName("controllers").WithName("AWSInstanceStateController"),
			Endpoints:        awsServiceEndpoints,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: instanceStateConcurrency, RecoverPanic: true}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSInstanceStateController")
			os.Exit(1)
		}
	}
	if feature.Gates.Enabled(feature.AutoControllerIdentityCreator) {
		setupLog.Info("AutoControllerIdentityCreator enabled")
		if err := (&controlleridentitycreator.AWSControllerIdentityReconciler{
			Client:           mgr.GetClient(),
			Log:              ctrl.Log.WithName("controllers").WithName("AWSControllerIdentity"),
			Endpoints:        awsServiceEndpoints,
			WatchFilterValue: watchFilterValue,
		}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: awsClusterConcurrency, RecoverPanic: true}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "AWSControllerIdentity")
			os.Exit(1)
		}
	}

	if feature.Gates.Enabled(feature.BootstrapFormatIgnition) {
		setupLog.Info("Enabling Ignition support for machine bootstrap data")
	}
}
func initFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&metricsBindAddr,
		"metrics-bind-addr",
		"localhost:8080",
		"The address the metric endpoint binds to.",
	)

	fs.BoolVar(
		&enableLeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
	)

	fs.StringVar(
		&watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	fs.StringVar(
		&leaderElectionNamespace,
		"leader-elect-namespace",
		"",
		"Namespace that the controller performs leader election in. If unspecified, the controller will discover which namespace it is running in.",
	)

	fs.StringVar(
		&profilerAddress,
		"profiler-address",
		"",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)",
	)

	fs.IntVar(&awsClusterConcurrency,
		"awscluster-concurrency",
		5,
		"Number of AWSClusters to process simultaneously",
	)

	fs.IntVar(&instanceStateConcurrency,
		"instance-state-concurrency",
		5,
		"Number of concurrent watches for instance state changes",
	)

	fs.IntVar(&awsMachineConcurrency,
		"awsmachine-concurrency",
		10,
		"Number of AWSMachines to process simultaneously",
	)

	fs.DurationVar(&syncPeriod,
		"sync-period",
		10*time.Minute,
		fmt.Sprintf("The minimum interval at which watched resources are reconciled. If EKS is enabled the maximum allowed is %s", maxEKSSyncPeriod),
	)

	fs.IntVar(&webhookPort,
		"webhook-port",
		9443,
		"Webhook Server port.",
	)

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")

	fs.StringVar(&healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	fs.StringVar(&serviceEndpoints,
		"service-endpoints",
		"",
		"Set custom AWS service endpoins in semi-colon separated format: ${SigningRegion1}:${ServiceID1}=${URL},${ServiceID2}=${URL};${SigningRegion2}...",
	)

	fs.StringVar(
		&watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel),
	)

	feature.MutableGates.AddFlag(fs)
}
