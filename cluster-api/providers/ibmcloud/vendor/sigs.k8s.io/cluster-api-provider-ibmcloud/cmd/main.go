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
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"

	// +kubebuilder:scaffold:imports
	"github.com/spf13/pflag"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cgrecord "k8s.io/client-go/tools/record"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/crdmigrator"
	"sigs.k8s.io/cluster-api/util/flags"

	powervsinfrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta2"
	powervsinfrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	vpcinfrav1beta1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta1"
	vpcinfrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta2"
	powervscontroller "sigs.k8s.io/cluster-api-provider-ibmcloud/internal/controllers/powervs"
	vpccontroller "sigs.k8s.io/cluster-api-provider-ibmcloud/internal/controllers/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	powervs "sigs.k8s.io/cluster-api-provider-ibmcloud/webhooks/powervs/admission"
	vpc "sigs.k8s.io/cluster-api-provider-ibmcloud/webhooks/vpc/admission"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/component-base/logs/json/register"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	// flags.
	watchNamespace         string
	enableLeaderElection   bool
	healthAddr             string
	syncPeriod             time.Duration
	managerOptions         = flags.ManagerOptions{}
	logOptions             = logs.NewOptions()
	webhookPort            int
	webhookCertDir         string
	watchFilterValue       string
	disableHTTP2           bool
	skipCRDMigrationPhases []string
)

func init() {
	klog.InitFlags(nil)

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(powervsinfrav1beta2.AddToScheme(scheme))
	utilruntime.Must(powervsinfrav1.AddToScheme(scheme))
	utilruntime.Must(vpcinfrav1beta1.AddToScheme(scheme))
	utilruntime.Must(vpcinfrav1.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func initFlags(fs *pflag.FlagSet) {
	logsv1.AddFlags(logOptions, fs)

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

	fs.StringVar(&webhookCertDir,
		"webhook-cert-dir",
		"/tmp/k8s-webhook-server/serving-certs/",
		"The webhook certificate directory, where the server should find the TLS certificate and key.")

	fs.StringVar(&watchFilterValue,
		"watch-filter",
		"",
		fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel))

	fs.BoolVar(&disableHTTP2,
		"disable-http2",
		true,
		"http/2 should be disabled due to its vulnerabilities. More specifically, disabling http/2 will prevent from being vulnerable to the HTTP/2 Stream Cancellation and Rapid Reset CVEs.")

	fs.StringSliceVar(&skipCRDMigrationPhases,
		"skip-crd-migration-phases",
		[]string{},
		"List of CRD migration phases to skip. Valid values are: StorageVersionMigration, CleanupManagedFields.")

	flags.AddManagerOptions(fs, &managerOptions)
}

func validateFlags() error {
	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		setupLog.Error(err, "unable to validate and apply log options")
		return err
	}

	return nil
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
// ADD CRD RBAC for CRD Migrator.
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions;customresourcedefinitions/status,verbs=update;patch,resourceNames=ibmpowervsclusters.infrastructure.cluster.x-k8s.io;ibmpowervsclustertemplates.infrastructure.cluster.x-k8s.io;ibmpowervsmachines.infrastructure.cluster.x-k8s.io;ibmpowervsmachinetemplates.infrastructure.cluster.x-k8s.io;ibmpowervsimages.infrastructure.cluster.x-k8s.io
// ADD CR RBAC for CRD Migrator.
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclustertemplates,verbs=get;list;watch;patch;update

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

	tlsOptions, metricsOptions, err := flags.GetManagerOptions(managerOptions)
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

	if disableHTTP2 {
		metricsOptions.TLSOpts = append(metricsOptions.TLSOpts, func(c *tls.Config) {
			setupLog.Info("disabling http/2")
			c.NextProtos = []string{"http/1.1"}
		})
	} else {
		setupLog.Info("WARNING: It is not recommended to enable http/2 due to https://github.com/kubernetes/kubernetes/issues/121197")
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
			TLSOpts: tlsOptions,
		}),
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					// We want to avoid use of cache for IBMPowerVSCluster as we exclusively depend on IBMPowerVSCluster.Status.[Resource].ControllerCreated
					// to mark resources created by controller.
					&powervsinfrav1.IBMPowerVSCluster{},
				},
			},
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	setupWebhooks(mgr)
	setupChecks(mgr)
	setupReconcilers(ctx, mgr, serviceEndpoint, watchFilterValue, skipCRDMigrationPhases)

	// +kubebuilder:scaffold:builder
	setupLog.Info("starting manager")
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

func setupReconcilers(ctx context.Context, mgr ctrl.Manager, serviceEndpoint []endpoints.ServiceEndpoint, watchFilterValue string, skipCRDMigrationPhases []string) {
	// Note: The kubebuilder RBAC markers above have to be kept in sync
	// with the CRDs that should be migrated by this provider.
	crdMigratorConfig := map[client.Object]crdmigrator.ByObjectConfig{
		&powervsinfrav1.IBMPowerVSCluster{}:         {UseCache: true, UseStatusForStorageVersionMigration: true},
		&powervsinfrav1.IBMPowerVSClusterTemplate{}: {UseCache: false},
		&powervsinfrav1.IBMPowerVSMachine{}:         {UseCache: true, UseStatusForStorageVersionMigration: true},
		&powervsinfrav1.IBMPowerVSMachineTemplate{}: {UseCache: false},
		&powervsinfrav1.IBMPowerVSImage{}:           {UseCache: true, UseStatusForStorageVersionMigration: true},
	}

	crdMigratorSkipPhases := make([]crdmigrator.Phase, 0, len(skipCRDMigrationPhases))
	for _, p := range skipCRDMigrationPhases {
		crdMigratorSkipPhases = append(crdMigratorSkipPhases, crdmigrator.Phase(p))
	}
	if err := (&crdmigrator.CRDMigrator{
		Client:                 mgr.GetClient(),
		APIReader:              mgr.GetAPIReader(),
		SkipCRDMigrationPhases: crdMigratorSkipPhases,
		Config:                 crdMigratorConfig,
		// The CRDMigrator is run with only concurrency 1 to ensure we don't overwhelm the apiserver by patching a
		// lot of CRs concurrently.
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: 1}); err != nil {
		setupLog.Error(err, "Unable to create controller", "controller", "CRDMigrator")
		os.Exit(1)
	}

	if err := (&vpccontroller.IBMVPCClusterReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("IBMVPCCluster"),
		Recorder:        mgr.GetEventRecorderFor("ibmvpccluster-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMVPCCluster")
		os.Exit(1)
	}

	if err := (&vpccontroller.IBMVPCMachineReconciler{
		Client:          mgr.GetClient(),
		Log:             ctrl.Log.WithName("controllers").WithName("IBMVPCMachine"),
		Recorder:        mgr.GetEventRecorderFor("ibmvpcmachine-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMVPCMachine")
		os.Exit(1)
	}

	if err := (&vpccontroller.IBMVPCMachineTemplateReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		ServiceEndpoint: serviceEndpoint,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ibmvpcmachinetemplate")
		os.Exit(1)
	}

	if err := (&powervscontroller.IBMPowerVSClusterReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("ibmpowervscluster-controller"),
		ServiceEndpoint:  serviceEndpoint,
		Scheme:           mgr.GetScheme(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSCluster")
		os.Exit(1)
	}

	if err := (&powervscontroller.IBMPowerVSMachineReconciler{
		Client:           mgr.GetClient(),
		Recorder:         mgr.GetEventRecorderFor("ibmpowervsmachine-controller"),
		ServiceEndpoint:  serviceEndpoint,
		Scheme:           mgr.GetScheme(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSMachine")
		os.Exit(1)
	}

	if err := (&powervscontroller.IBMPowerVSMachineTemplateReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ibmpowervsmachinetemplate")
		os.Exit(1)
	}

	if err := (&powervscontroller.IBMPowerVSImageReconciler{
		Client:          mgr.GetClient(),
		Recorder:        mgr.GetEventRecorderFor("ibmpowervsimage-controller"),
		ServiceEndpoint: serviceEndpoint,
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IBMPowerVSImage")
		os.Exit(1)
	}
}

func setupWebhooks(mgr ctrl.Manager) {
	if err := (&vpc.IBMVPCCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCCluster")
		os.Exit(1)
	}
	if err := (&vpc.IBMVPCMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCMachine")
		os.Exit(1)
	}
	if err := (&vpc.IBMVPCMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMVPCMachineTemplate")
		os.Exit(1)
	}
	if err := (&powervs.IBMPowerVSCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSCluster")
		os.Exit(1)
	}
	if err := (&powervs.IBMPowerVSMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSMachine")
		os.Exit(1)
	}
	if err := (&powervs.IBMPowerVSMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSMachineTemplate")
		os.Exit(1)
	}
	if err := (&powervs.IBMPowerVSImage{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSImage")
		os.Exit(1)
	}
	if err := (&powervs.IBMPowerVSClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "IBMPowerVSClusterTemplate")
		os.Exit(1)
	}
}
