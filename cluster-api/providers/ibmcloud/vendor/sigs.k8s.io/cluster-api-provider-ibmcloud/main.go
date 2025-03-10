/*
Copyright 2021 The Kubernetes Authors.

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

// main is the main package for the Cluster API IBMCLOUD Provider.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	// +kubebuilder:scaffold:imports
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/flags"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta1"
	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/controllers"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/options"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/component-base/logs/json/register"
)

var (
	watchNamespace       string
	enableLeaderElection bool
	healthAddr           string
	syncPeriod           time.Duration
	managerOptions       = flags.ManagerOptions{}
	logOptions           = logs.NewOptions()
	webhookPort          int
	webhookCertDir       string

	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	klog.InitFlags(nil)

	_ = clientgoscheme.AddToScheme(scheme)

	_ = infrav1beta1.AddToScheme(scheme)
	_ = infrav1beta2.AddToScheme(scheme)
	_ = capiv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func initFlags(fs *pflag.FlagSet) {
	fs.BoolVar(
		&enableLeaderElection,
		"leader-elect",
		false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.",
	)

	fs.StringVar(
		&watchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.",
	)

	fs.StringVar(
		&healthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)

	fs.DurationVar(
		&syncPeriod,
		"sync-period",
		10*time.Minute,
		"The minimum interval at which watched resources are reconciled.",
	)
	fs.StringVar(
		&options.ProviderIDFormat,
		"provider-id-fmt",
		string(options.ProviderIDFormatV2),
		"ProviderID format is used set the Provider ID format for Machine",
	)

	fs.StringVar(
		&endpoints.ServiceEndpointFormat,
		"service-endpoint",
		"",
		"Set custom service endpoint in semi-colon separated format: ${ServiceRegion1}:${ServiceID1}=${URL1},${ServiceID2}=${URL2};${ServiceRegion2}:${ServiceID1}=${URL1}",
	)

	fs.IntVar(&webhookPort,
		"webhook-port",
		9443,
		"The webhook server port the manager will listen on.",
	)

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"The webhook certificate directory, where the server should find the TLS certificate and key.")

	logsv1.AddFlags(logOptions, fs)
	flags.AddManagerOptions(fs, &managerOptions)
}

func validateFlags() error {
	if options.ProviderIDFormatType(options.ProviderIDFormat) == options.ProviderIDFormatV2 {
		setupLog.Info("Using v2 version of ProviderID format")
	} else {
		return fmt.Errorf("invalid value for flag provider-id-fmt: %s, Only supported value is %s", options.ProviderIDFormat, options.ProviderIDFormatV2)
	}

	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		setupLog.Error(err, "unable to validate and apply log options")
		return err
	}

	return nil
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

func main() {
	initFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klog.Background())

	// Parse service endpoints.
	serviceEndpoint, err := endpoints.ParseServiceEndpointFlag(endpoints.ServiceEndpointFormat)
	if err != nil {
		setupLog.Error(err, "unable to parse service endpoint flag", "controller", "cluster")
		os.Exit(1)
	}

	if err := validateFlags(); err != nil {
		setupLog.Error(err, "Flag validation failure")
		os.Exit(1)
	}

	if watchNamespace != "" {
		setupLog.Info("Watching cluster-api objects only in namespace for reconciliation", "namespace", watchNamespace)
	}

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	_, metricsOptions, err := flags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: invalid flags")
		os.Exit(1)
	}

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:           scheme,
		LeaderElection:   enableLeaderElection,
		Metrics:          *metricsOptions,
		LeaderElectionID: "effcf9b8.cluster.x-k8s.io",
		Cache: cache.Options{
			DefaultNamespaces: watchNamespaces,
			SyncPeriod:        &syncPeriod,
		},
		EventBroadcaster:       broadcaster,
		HealthProbeBindAddress: healthAddr,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    webhookPort,
			CertDir: webhookCertDir,
		}),
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					// We want to avoid use of cache for IBMPowerVSCluster as we exclusively depend on IBMPowerVSCluster.Status.[Resource].ControllerCreated
					// to mark resources created by controller.
					&infrav1beta2.IBMPowerVSCluster{},
				},
			},
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("ibmcloud-controller"))

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	setupReconcilers(ctx, mgr, serviceEndpoint)
	setupWebhooks(mgr)
	setupChecks(mgr)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func setupReconcilers(ctx context.Context, mgr ctrl.Manager, serviceEndpoint []endpoints.ServiceEndpoint) {
	if err := (&controllers.IBMVPCClusterReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("IBMVPCCluster"),
		Recorder:        mgr.GetEventRecorderFor("ibmvpccluster-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMVPCCluster")
		os.Exit(1)
	}

	if err := (&controllers.IBMVPCMachineReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("IBMVPCMachine"),
		Recorder:        mgr.GetEventRecorderFor("ibmvpcmachine-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMVPCMachine")
		os.Exit(1)
	}

	if err := (&controllers.IBMPowerVSClusterReconciler{
		Client:          mgr.GetClient(),
		Recorder:        mgr.GetEventRecorderFor("ibmpowervscluster-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSCluster")
		os.Exit(1)
	}

	if err := (&controllers.IBMPowerVSMachineReconciler{
		Client:          mgr.GetClient(),
		Recorder:        mgr.GetEventRecorderFor("ibmpowervsmachine-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSMachine")
		os.Exit(1)
	}

	if err := (&controllers.IBMPowerVSImageReconciler{
		Client:          mgr.GetClient(),
		Recorder:        mgr.GetEventRecorderFor("ibmpowervsimage-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSImage")
		os.Exit(1)
	}

	if err := (&controllers.IBMPowerVSMachineTemplateReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ibmpowervsmachinetemplate")
		os.Exit(1)
	}

	if err := (&controllers.IBMVPCMachineTemplateReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		ServiceEndpoint: serviceEndpoint,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ibmvpcmachinetemplate")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&infrav1beta2.IBMVPCCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCCluster")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMVPCMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCMachine")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMVPCMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCMachineTemplate")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMPowerVSCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSCluster")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMPowerVSMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSMachine")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMPowerVSMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSMachineTemplate")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMPowerVSImage{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSImage")
		os.Exit(1)
	}
	if err := (&infrav1beta2.IBMPowerVSClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSClusterTemplate")
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
