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
	"flag"
	"fmt"
	"os"
	"time"

	// +kubebuilder:scaffold:imports
	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"

	ctrl "sigs.k8s.io/controller-runtime"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta1"
	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/controllers"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/options"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	watchNamespace       string
	metricsAddr          string
	enableLeaderElection bool
	healthAddr           string
	syncPeriod           time.Duration

	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = infrav1beta1.AddToScheme(scheme)
	_ = infrav1beta2.AddToScheme(scheme)
	_ = capiv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	klog.InitFlags(nil)

	initFlags(pflag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	ctrl.SetLogger(klogr.New())

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

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "effcf9b8.cluster.x-k8s.io",
		SyncPeriod:             &syncPeriod,
		Namespace:              watchNamespace,
		EventBroadcaster:       broadcaster,
		HealthProbeBindAddress: healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("ibmcloud-controller"))

	setupReconcilers(mgr, serviceEndpoint)
	setupWebhooks(mgr)
	setupChecks(mgr)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&metricsAddr,
		"metrics-bind-addr",
		":8080",
		"The address the metric endpoint binds to.",
	)

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

	// Deprecated: Use provider-id-fmt flag go set provider id format for Power VS.
	fs.StringVar(
		&options.PowerVSProviderIDFormat,
		"powervs-provider-id-fmt",
		string(options.PowerVSProviderIDFormatV1),
		"ProviderID format is used set the Provider ID format for Machine",
	)
	_ = fs.MarkDeprecated("powervs-provider-id-fmt", "please use provider-id-fmt flag")

	fs.StringVar(
		&options.ProviderIDFormat,
		"provider-id-fmt",
		string(options.ProviderIDFormatV1),
		"ProviderID format is used set the Provider ID format for Machine",
	)

	fs.StringVar(
		&endpoints.ServiceEndpointFormat,
		"service-endpoint",
		"",
		"Set custom service endpoint in semi-colon separated format: ${ServiceRegion1}:${ServiceID1}=${URL1},${ServiceID2}=${URL2};${ServiceRegion2}:${ServiceID1}=${URL1}",
	)
}

func validateFlags() error {
	switch options.ProviderIDFormatType(options.PowerVSProviderIDFormat) {
	case options.PowerVSProviderIDFormatV1:
		setupLog.Info("Power VS ProviderID format is set to v1 version")
	case options.PowerVSProviderIDFormatV2:
		setupLog.Info("Power VS ProviderID format is set to v2 version")
	default:
		return fmt.Errorf("invalid value for flag powervs-provider-id-fmt: %s, Supported values: v1, v2 ", options.PowerVSProviderIDFormat)
	}
	switch options.ProviderIDFormatType(options.ProviderIDFormat) {
	case options.ProviderIDFormatV1:
		setupLog.Info("Using v1 version of ProviderID format")
	case options.ProviderIDFormatV2:
		setupLog.Info("Using v2 version of ProviderID format")
	default:
		return fmt.Errorf("invalid value for flag provider-id-fmt: %s, Supported values: %s, %s ", options.ProviderIDFormat, options.ProviderIDFormatV1, options.ProviderIDFormatV2)
	}
	return nil
}

func setupReconcilers(mgr ctrl.Manager, serviceEndpoint []endpoints.ServiceEndpoint) {
	if err := (&controllers.IBMVPCClusterReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("IBMVPCCluster"),
		Recorder:        mgr.GetEventRecorderFor("ibmvpccluster-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
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
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSCluster")
		os.Exit(1)
	}

	if err := (&controllers.IBMPowerVSMachineReconciler{
		Client:          mgr.GetClient(),
		Recorder:        mgr.GetEventRecorderFor("ibmpowervsmachine-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
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
