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
	"flag"
	"fmt"
	"os"
	"reflect"
	goruntime "runtime"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/pflag"
	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1alpha5 "github.com/vmware-tanzu/vm-operator/api/v1alpha5"
	vmoprv1alpha6 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
	"gopkg.in/fsnotify.v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsv1 "k8s.io/component-base/logs/api/v1"
	_ "k8s.io/component-base/logs/json/register"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	controlplanev1 "sigs.k8s.io/cluster-api/api/controlplane/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/controllers/clustercache"
	"sigs.k8s.io/cluster-api/controllers/crdmigrator"
	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/cluster-api/util/apiwarnings"
	capicontrollerutil "sigs.k8s.io/cluster-api/util/controller"
	capiflags "sigs.k8s.io/cluster-api/util/flags"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/controllers"
	"sigs.k8s.io/cluster-api-provider-vsphere/controllers/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	conversionapi "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/version"
	"sigs.k8s.io/cluster-api-provider-vsphere/webhooks"
	vmwarewebhooks "sigs.k8s.io/cluster-api-provider-vsphere/webhooks/vmware"
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
	clusterCacheClientQPS       float32
	clusterCacheClientBurst     int
	syncPeriod                  time.Duration
	webhookOpts                 webhook.Options
	watchNamespace              string
	vmOperatorAPIVersion        string
	featureGates                string

	clusterCacheConcurrency           int
	vSphereClusterConcurrency         int
	vSphereMachineConcurrency         int
	vSphereMachineTemplateConcurrency int
	providerServiceAccountConcurrency int
	serviceDiscoveryConcurrency       int
	vSphereVMConcurrency              int
	vSphereClusterIdentityConcurrency int
	vSphereDeploymentZoneConcurrency  int
	virtualMachineGroupConcurrency    int
	skipCRDMigrationPhases            []string

	managerOptions = capiflags.ManagerOptions{}

	defaultProfilerAddr     = os.Getenv("PROFILER_ADDR")
	defaultSyncPeriod       = manager.DefaultSyncPeriod
	defaultLeaderElectionID = manager.DefaultLeaderElectionID
	defaultPodName          = manager.DefaultPodName

	supportedVMOperatorAPIVersions = []string{vmoprv1alpha2.GroupVersion.Version, vmoprv1alpha5.GroupVersion.Version, vmoprv1alpha6.GroupVersion.Version}
)

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	// Flags specific to CAPV

	fs.StringVar(
		&managerOpts.LeaderElectionID,
		"leader-election-id",
		defaultLeaderElectionID,
		"Name of the config map to use as the locking resource when configuring leader election.")

	fs.IntVar(&clusterCacheConcurrency, "clustercache-concurrency", 100,
		"Number of clusters to process simultaneously")

	fs.IntVar(&vSphereClusterConcurrency, "vspherecluster-concurrency", 50,
		"Number of vSphere clusters to process simultaneously")

	fs.IntVar(&vSphereMachineConcurrency, "vspheremachine-concurrency", 100,
		"Number of vSphere machines to process simultaneously")

	fs.IntVar(&vSphereMachineTemplateConcurrency, "vspheremachinetemplate-concurrency", 10,
		"Number of vSphere machine templates to process simultaneously")

	fs.IntVar(&providerServiceAccountConcurrency, "providerserviceaccount-concurrency", 50,
		"Number of provider service accounts to process simultaneously")

	fs.IntVar(&serviceDiscoveryConcurrency, "servicediscovery-concurrency", 50,
		"Number of vSphere clusters for service discovery to process simultaneously")

	fs.IntVar(&vSphereVMConcurrency, "vspherevm-concurrency", 100,
		"Number of vSphere vms to process simultaneously")

	fs.IntVar(&vSphereClusterIdentityConcurrency, "vsphereclusteridentity-concurrency", 10,
		"Number of vSphere cluster identities to process simultaneously")

	fs.IntVar(&vSphereDeploymentZoneConcurrency, "vspheredeploymentzone-concurrency", 10,
		"Number of vSphere deployment zones to process simultaneously")

	fs.IntVar(&virtualMachineGroupConcurrency, "virtualmachinegroup-concurrency", 50,
		"Number of virtual machine group to process simultaneously")

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

	fs.StringVar(
		&managerOpts.NetworkProvider,
		"network-provider",
		"",
		"network provider to be used by Supervisor based clusters.",
	)

	fs.StringVar(
		&vmOperatorAPIVersion,
		"vm-operator-api-version",
		vmoprv1alpha5.GroupVersion.Version,
		fmt.Sprintf("the API version to use when reading and writing VM Operator resources in supervisor mode. Valid values are: %s", strings.Join(supportedVMOperatorAPIVersions, ", ")),
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

	fs.StringSliceVar(&skipCRDMigrationPhases, "skip-crd-migration-phases", []string{},
		"List of CRD migration phases to skip. Valid values are: StorageVersionMigration, CleanupManagedFields.")

	fs.DurationVar(&syncPeriod, "sync-period", defaultSyncPeriod,
		"The minimum interval at which watched resources are reconciled (e.g. 15m)")

	fs.Float32Var(&restConfigQPS, "kube-api-qps", 100,
		"Maximum queries per second from the controller client to the Kubernetes API server.")

	fs.IntVar(&restConfigBurst, "kube-api-burst", 200,
		"Maximum number of queries that should be allowed in one burst from the controller client to the Kubernetes API server.")

	fs.Float32Var(&clusterCacheClientQPS, "clustercache-client-qps", 20,
		"Maximum queries per second from the cluster cache clients to the Kubernetes API server of workload clusters.")

	fs.IntVar(&clusterCacheClientBurst, "clustercache-client-burst", 30,
		"Maximum number of queries that should be allowed in one burst from the cluster cache clients to the Kubernetes API server of workload clusters.")

	fs.IntVar(&webhookOpts.Port, "webhook-port", 9443,
		"Webhook Server port.")

	fs.StringVar(&webhookOpts.CertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir.")

	fs.StringVar(&webhookOpts.CertName, "webhook-cert-name", "tls.crt",
		"Webhook cert name.")

	fs.StringVar(&webhookOpts.KeyName, "webhook-key-name", "tls.key",
		"Webhook key name.")

	fs.StringVar(&managerOpts.HealthProbeBindAddress, "health-addr", ":9440",
		"The address the health endpoint binds to.",
	)

	capiflags.AddManagerOptions(fs, &managerOptions)
	feature.AddFlag(fs, &featureGates, supportedVMOperatorAPIVersions)
}

// Add RBAC for the authorized diagnostics endpoint.
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create
// ADD CRD RBAC for CRD Migrator.
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch
// govmomi
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions;customresourcedefinitions/status,verbs=update;patch,resourceNames=vsphereclusters.infrastructure.cluster.x-k8s.io;vsphereclustertemplates.infrastructure.cluster.x-k8s.io;vspheremachines.infrastructure.cluster.x-k8s.io;vspheremachinetemplates.infrastructure.cluster.x-k8s.io;vspherevms.infrastructure.cluster.x-k8s.io;vsphereclusteridentities.infrastructure.cluster.x-k8s.io;vspheredeploymentzones.infrastructure.cluster.x-k8s.io;vspherefailuredomains.infrastructure.cluster.x-k8s.io
// supervisor
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions;customresourcedefinitions/status,verbs=update;patch,resourceNames=vsphereclusters.vmware.infrastructure.cluster.x-k8s.io;vsphereclustertemplates.vmware.infrastructure.cluster.x-k8s.io;vspheremachines.vmware.infrastructure.cluster.x-k8s.io;vspheremachinetemplates.vmware.infrastructure.cluster.x-k8s.io;providerserviceaccounts.vmware.infrastructure.cluster.x-k8s.io
// govmomi CRs
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates;vsphereclustertemplates,verbs=get;list;watch;patch;update

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	// Set log level 2 as default.
	if err := pflag.CommandLine.Set("v", "2"); err != nil {
		fmt.Printf("Failed to set default log level: %v\n", err)
		os.Exit(1)
	}
	pflag.Parse()

	if err := logsv1.ValidateAndApply(logOptions, nil); err != nil {
		fmt.Printf("Unable to start manager: %v\n", err)
		os.Exit(1)
	}

	pflag.CommandLine.VisitAll(func(flag *pflag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	// Note: setupLog can only be used after ctrl.SetLogger was called
	setupLog.Info(fmt.Sprintf("Version: %s (git commit: %s)", version.Get().String(), version.Get().GitCommit))

	// Set up the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	managerOpts.KubeConfig = ctrl.GetConfigOrDie()
	managerOpts.KubeConfig.QPS = restConfigQPS
	managerOpts.KubeConfig.Burst = restConfigBurst
	managerOpts.KubeConfig.UserAgent = remote.DefaultClusterAPIUserAgent(controllerName)
	managerOpts.KubeConfig.WarningHandler = apiwarnings.DefaultHandler(klog.Background().WithName("API Server Warning"))

	// Detect which CRDs are loaded.
	httpClient, err := rest.HTTPClientFor(managerOpts.KubeConfig)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: unable to create http client to detect CRDs")
		os.Exit(1)
	}
	restMapper, err := apiutil.NewDynamicRESTMapper(managerOpts.KubeConfig, httpClient)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: unable to create RESTMapper to detect CRDs")
		os.Exit(1)
	}

	govmomiGVR := infrav1.GroupVersion.WithResource(reflect.TypeOf(&infrav1.VSphereCluster{}).Elem().Name())
	supervisorGVR := vmwarev1.GroupVersion.WithResource(reflect.TypeOf(&vmwarev1.VSphereCluster{}).Elem().Name())

	var isSupervisorCRDLoaded, isGovmomiCRDLoaded bool
	var errGovmomi, errSupervisor error
	if err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 30*time.Second, true, func(_ context.Context) (bool, error) {
		// Check for non-supervisor VSphereCluster and start controller if found
		isGovmomiCRDLoaded, errGovmomi = isCRDDeployed(restMapper, govmomiGVR)

		// Check for supervisor VSphereCluster and start controller if found
		isSupervisorCRDLoaded, errSupervisor = isCRDDeployed(restMapper, supervisorGVR)

		// One of govmomi/supervisor mode should be detected, otherwise keep trying until timeout to handle
		// race conditions during controllers startup right after install or upgrades, when also CRDs
		// are installed in a short time frame.
		if (isGovmomiCRDLoaded && errGovmomi == nil) || (isSupervisorCRDLoaded && errSupervisor == nil) {
			return true, nil
		}
		return false, nil
	}); err != nil {
		// Continuing startup does not make sense without one of govmomi/supervisor mode detected.
		// The Pod goes in CrashLoopBack and eventually recover, but failing to detect CRD after 30s is usually
		// a signal of some problem.
		err := fmt.Errorf("neither supervisor nor govmomi CRDs detected: %w", kerrors.NewAggregate([]error{err, errGovmomi, errSupervisor}))
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}
	if isGovmomiCRDLoaded && isSupervisorCRDLoaded {
		err := fmt.Errorf("both supervisor and govmomi CRDs detected (unsupported configuration)")
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	if isGovmomiCRDLoaded {
		if err := feature.SetGovmomiGates(featureGates); err != nil {
			setupLog.Error(err, "invalid argument: --feature-gates")
			os.Exit(1)
		}
	}

	var vm runtime.Object
	var converter *conversion.Converter
	if isSupervisorCRDLoaded {
		if !sets.New(supportedVMOperatorAPIVersions...).Has(vmOperatorAPIVersion) {
			setupLog.Info(fmt.Sprintf("Invalid argument: --vm-operator-api-version must be one of : %s\n", strings.Join(supportedVMOperatorAPIVersions, ", ")))
			os.Exit(1)
		}
		setupLog.Info(fmt.Sprintf("Target API Version for group %s: %s", vmoprvhub.GroupVersion.Group, vmOperatorAPIVersion))

		if err := feature.SetSupervisorGates(vmOperatorAPIVersion, featureGates); err != nil {
			setupLog.Error(err, "invalid argument: --feature-gates")
			os.Exit(1)
		}

		converter = conversionapi.DefaultConverterFor(
			schema.GroupVersion{Group: vmoprvhub.GroupVersion.Group, Version: vmOperatorAPIVersion},
		)

		// Get vm-operator native types in the preferred version for cache filters.
		vmGVK, err := converter.SpokeGroupVersionKindFor(&vmoprvhub.VirtualMachine{})
		if err != nil {
			setupLog.Error(err, "Unable to start manager; failed to get object for VirtualMachine cache filter")
			os.Exit(1)
		}

		scheme := runtime.NewScheme()
		if err := vmoprv1alpha2.AddToScheme(scheme); err != nil {
			setupLog.Error(err, "Unable to start manager; failed register v1alpha2 version for vm-operator API types")
			os.Exit(1)
		}
		if err := vmoprv1alpha5.AddToScheme(scheme); err != nil {
			setupLog.Error(err, "Unable to start manager; failed register v1alpha5 version for vm-operator API types")
			os.Exit(1)
		}
		if err := vmoprv1alpha6.AddToScheme(scheme); err != nil {
			setupLog.Error(err, "Unable to start manager; failed register v1alpha6 version for vm-operator API types")
			os.Exit(1)
		}

		vm, err = scheme.New(vmGVK)
		if err != nil {
			setupLog.Error(err, "Unable to start manager; failed to create object for VirtualMachine cache filter")
			os.Exit(1)
		}
	}

	req, _ := labels.NewRequirement(clusterv1.ClusterNameLabel, selection.Exists, nil)
	clusterSecretCacheSelector := labels.NewSelector().Add(*req)

	req, _ = labels.NewRequirement(vmoperator.ClusterSelectorKey, selection.Exists, nil)
	virtualMachineCacheSelector := labels.NewSelector().Add(*req)

	managerOpts.Scheme = runtime.NewScheme()
	managerOpts.Cache = cache.Options{
		DefaultNamespaces: watchNamespaces,
		SyncPeriod:        &syncPeriod,
		DefaultTransform:  cache.TransformStripManagedFields(),
		ByObject: func() map[client.Object]cache.ByObject {
			byObject := map[client.Object]cache.ByObject{}
			// Optimize the cache for supervisor mode. govmomi mode might have different caching requirements
			if isSupervisorCRDLoaded {
				// Note: Only VirtualMachines with the cluster name label are cached (vmopmachine.go sets this label).
				// NOTE: use vm-operator native types for cache filters (the reconciler uses the internal hub version).
				byObject[vm.(client.Object)] = cache.ByObject{
					Label: virtualMachineCacheSelector,
				}
				byObject[&corev1.Secret{}] = cache.ByObject{
					Label: clusterSecretCacheSelector,
					// Drop data of secrets that we don't use.
					Transform: func(in any) (any, error) {
						if s, ok := in.(*corev1.Secret); ok {
							s.SetManagedFields(nil)
							if !strings.HasSuffix(s.Name, "-kubeconfig") && s.Type != corev1.SecretTypeServiceAccountToken {
								s.Data = nil
							}
						}
						return in, nil
					},
				}
				byObject[&corev1.ConfigMap{}] = cache.ByObject{
					Namespaces: map[string]cache.Config{
						metav1.NamespacePublic: {},
						util.NCPNamespace:      {},
					},
				}
				byObject[&clusterv1.MachineSet{}] = cache.ByObject{
					// Drop data of MachineSets as we only use ownerRefs of MachineSets in clog.AddOwners.
					Transform: func(in any) (any, error) {
						if m, ok := in.(*clusterv1.MachineSet); ok {
							m.SetManagedFields(nil)
							m.Spec = clusterv1.MachineSetSpec{}
							m.Status = clusterv1.MachineSetStatus{}
						}
						return in, nil
					},
				}
				if feature.Gates.Enabled(feature.IPv6DualStack) {
					byObject[&controlplanev1.KubeadmControlPlane{}] = cache.ByObject{
						Transform: func(in any) (any, error) {
							kcp, ok := in.(*controlplanev1.KubeadmControlPlane)
							if !ok {
								return in, nil
							}

							kcp.SetManagedFields(nil)

							kcp.Spec = controlplanev1.KubeadmControlPlaneSpec{
								KubeadmConfigSpec: bootstrapv1.KubeadmConfigSpec{
									ClusterConfiguration: bootstrapv1.ClusterConfiguration{
										APIServer: bootstrapv1.APIServer{
											CertSANs: kcp.Spec.KubeadmConfigSpec.ClusterConfiguration.APIServer.CertSANs,
										},
									},
								},
							}

							kcp.Status = controlplanev1.KubeadmControlPlaneStatus{
								ObservedGeneration: kcp.Status.ObservedGeneration,
								Conditions:         kcp.Status.Conditions,
							}

							return kcp, nil
						},
					}
				}
			}
			return byObject
		}(),
		NewInformer: capicontrollerutil.NewInformerFunc(managerOpts.Scheme, controllerName),
	}
	managerOpts.Client = func() client.Options {
		// Optimize the cache for supervisor mode. govmomi mode might have different caching requirements
		if isSupervisorCRDLoaded {
			return client.Options{
				Cache: &client.CacheOptions{
					DisableFor: []client.Object{
						// We are configuring the mgr.GetClient() to not use the cache for secrets, so that if we miss some
						// Secret Get calls they are not hitting a cache with partial data.
						// If Secrets from the cache should be accessed, we should use the secretCachingClient instead.
						&corev1.Secret{},
					},
				},
			}
		}
		return client.Options{}
	}()

	if enableContentionProfiling {
		goruntime.SetBlockProfileRate(1)
	}

	managerOpts.LeaseDuration = &leaderElectionLeaseDuration
	managerOpts.RenewDeadline = &leaderElectionRenewDeadline
	managerOpts.RetryPeriod = &leaderElectionRetryPeriod

	// Create a function that adds all the controllers and webhooks to the manager.
	addToManager := func(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager) error {
		secretCachingClient, err := client.New(mgr.GetConfig(), client.Options{
			HTTPClient: mgr.GetHTTPClient(),
			Cache: &client.CacheOptions{
				Reader: mgr.GetCache(),
			},
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "unable to create secret caching client")
		}

		clusterCache, err := setupClusterCache(ctx, mgr, secretCachingClient, isSupervisorCRDLoaded)
		if err != nil {
			return pkgerrors.Wrapf(err, "unable to create remote cluster cache tracker")
		}

		if isGovmomiCRDLoaded {
			if err := setupVAPIControllers(ctx, controllerCtx, mgr, clusterCache); err != nil {
				return fmt.Errorf("setupVAPIControllers: %w", err)
			}
		} else {
			setupLog.Info(fmt.Sprintf("CRD for %s not loaded, skipping.", govmomiGVR.String()))
		}

		if isSupervisorCRDLoaded {
			if err := setupSupervisorControllers(ctx, controllerCtx, mgr, clusterCache, secretCachingClient); err != nil {
				return fmt.Errorf("setupSupervisorControllers: %w", err)
			}
		} else {
			setupLog.Info(fmt.Sprintf("CRD for %s not loaded, skipping.", supervisorGVR.String()))
		}

		// Note: The kubebuilder RBAC markers above has to be kept in sync
		// with the CRDs that should be migrated by this provider.
		crdMigratorConfig := map[client.Object]crdmigrator.ByObjectConfig{}
		if isGovmomiCRDLoaded {
			crdMigratorConfig[&infrav1.VSphereCluster{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&infrav1.VSphereClusterTemplate{}] = crdmigrator.ByObjectConfig{UseCache: false}
			crdMigratorConfig[&infrav1.VSphereMachine{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&infrav1.VSphereMachineTemplate{}] = crdmigrator.ByObjectConfig{UseCache: false}
			crdMigratorConfig[&infrav1.VSphereVM{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&infrav1.VSphereClusterIdentity{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&infrav1.VSphereDeploymentZone{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&infrav1.VSphereFailureDomain{}] = crdmigrator.ByObjectConfig{UseCache: false}
		}
		if isSupervisorCRDLoaded {
			crdMigratorConfig[&vmwarev1.VSphereCluster{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&vmwarev1.VSphereClusterTemplate{}] = crdmigrator.ByObjectConfig{UseCache: false}
			crdMigratorConfig[&vmwarev1.VSphereMachine{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&vmwarev1.VSphereMachineTemplate{}] = crdmigrator.ByObjectConfig{UseCache: false, UseStatusForStorageVersionMigration: true}
			crdMigratorConfig[&vmwarev1.ProviderServiceAccount{}] = crdmigrator.ByObjectConfig{UseCache: false}
		}

		crdMigratorSkipPhases := []crdmigrator.Phase{}
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
		}).SetupWithManager(ctx, mgr, concurrency(1)); err != nil {
			setupLog.Error(err, "Unable to create controller", "controller", "CRDMigrator")
			os.Exit(1)
		}

		return nil
	}

	tlsOptions, metricsOptions, err := capiflags.GetManagerOptions(managerOptions)
	if err != nil {
		setupLog.Error(err, "Unable to start manager: invalid flags")
		os.Exit(1)
	}
	webhookOpts.TLSOpts = tlsOptions
	managerOpts.WebhookServer = webhook.NewServer(webhookOpts)
	managerOpts.AddToManager = addToManager
	managerOpts.Metrics = *metricsOptions
	managerOpts.Controller = config.Controller{
		UsePriorityQueue: ptr.To[bool](feature.Gates.Enabled(feature.PriorityQueue)),
		// Give the manager more time to sync the caches during startup. This is required
		// in high scale environments when they are more objects in the system (default is 3m).
		CacheSyncTimeout: 5 * time.Minute,
	}
	managerOpts.Converter = converter

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

func setupVAPIControllers(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager, clusterCache clustercache.ClusterCache) error {
	if err := (&webhooks.VSphereCluster{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&webhooks.VSphereClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&webhooks.VSphereClusterIdentity{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := (&webhooks.VSphereMachine{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereVM{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereDeploymentZone{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := (&webhooks.VSphereFailureDomain{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}

	if err := controllers.AddClusterControllerToManager(ctx, controllerCtx, mgr, false, concurrency(vSphereClusterConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddMachineControllerToManager(ctx, controllerCtx, mgr, false, concurrency(vSphereMachineConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddVMControllerToManager(ctx, controllerCtx, mgr, clusterCache, concurrency(vSphereVMConcurrency)); err != nil {
		return err
	}
	if err := controllers.AddVsphereClusterIdentityControllerToManager(ctx, controllerCtx, mgr, concurrency(vSphereClusterIdentityConcurrency)); err != nil {
		return err
	}

	return controllers.AddVSphereDeploymentZoneControllerToManager(ctx, controllerCtx, mgr, concurrency(vSphereDeploymentZoneConcurrency))
}

func setupSupervisorControllers(ctx context.Context, controllerCtx *capvcontext.ControllerManagerContext, mgr ctrlmgr.Manager, clusterCache clustercache.ClusterCache, secretCachingClient client.Client) error {
	if err := (&vmwarewebhooks.VSphereMachineTemplate{}).SetupWebhookWithManager(mgr, controllerCtx.NetworkProvider); err != nil {
		return err
	}
	if err := (&vmwarewebhooks.VSphereMachine{}).SetupWebhookWithManager(mgr, controllerCtx.NetworkProvider); err != nil {
		return err
	}
	if err := (&vmwarewebhooks.VSphereClusterTemplate{}).SetupWebhookWithManager(mgr, controllerCtx.NetworkProvider); err != nil {
		return err
	}
	if err := (&vmwarewebhooks.VSphereCluster{}).SetupWebhookWithManager(mgr, controllerCtx.NetworkProvider); err != nil {
		return err
	}
	if err := (&vmwarewebhooks.ProviderServiceAccount{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	if err := controllers.AddClusterControllerToManager(ctx, controllerCtx, mgr, true, concurrency(vSphereClusterConcurrency)); err != nil {
		return err
	}

	if err := controllers.AddMachineControllerToManager(ctx, controllerCtx, mgr, true, concurrency(vSphereMachineConcurrency)); err != nil {
		return err
	}

	if err := vmware.AddVSphereMachineTemplateControllerToManager(ctx, controllerCtx, mgr, concurrency(vSphereMachineTemplateConcurrency)); err != nil {
		return err
	}

	if err := vmware.AddServiceAccountProviderControllerToManager(ctx, controllerCtx, mgr, clusterCache, secretCachingClient, concurrency(providerServiceAccountConcurrency)); err != nil {
		return err
	}

	if feature.Gates.Enabled(feature.NamespaceScopedZones) && feature.Gates.Enabled(feature.NodeAutoPlacement) {
		if err := vmware.AddVirtualMachineGroupControllerToManager(ctx, controllerCtx, mgr, concurrency(virtualMachineGroupConcurrency)); err != nil {
			return err
		}
	}

	return vmware.AddServiceDiscoveryControllerToManager(ctx, controllerCtx, mgr, clusterCache, concurrency(serviceDiscoveryConcurrency))
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

func isCRDDeployed(restMapper meta.RESTMapper, gvr schema.GroupVersionResource) (bool, error) {
	_, err := restMapper.KindFor(gvr)
	if err != nil {
		if meta.IsNoMatchError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func concurrency(c int) controller.Options {
	return controller.Options{MaxConcurrentReconciles: c}
}

func setupClusterCache(ctx context.Context, mgr ctrlmgr.Manager, secretCachingClient client.Client, isSupervisorCRDLoaded bool) (clustercache.ClusterCache, error) {
	clusterCache, err := clustercache.SetupWithManager(ctx, mgr, clustercache.Options{
		SecretClient: secretCachingClient,
		Cache: func() clustercache.CacheOptions {
			if isSupervisorCRDLoaded {
				return clustercache.CacheOptions{
					ByObject: map[client.Object]cache.ByObject{
						&corev1.Service{}: {
							Namespaces: map[string]cache.Config{
								vmwarev1.SupervisorHeadlessSvcNamespace: {
									Transform: func(in any) (any, error) {
										if s, ok := in.(*corev1.Service); ok {
											s.SetManagedFields(nil)
											if s.Name != vmwarev1.SupervisorHeadlessSvcName {
												s.Spec = corev1.ServiceSpec{}
												s.Status = corev1.ServiceStatus{}
											}
										}
										return in, nil
									},
								},
							},
						},
						&corev1.Endpoints{}: {
							Namespaces: map[string]cache.Config{
								vmwarev1.SupervisorHeadlessSvcNamespace: {
									Transform: func(in any) (any, error) {
										if s, ok := in.(*corev1.Endpoints); ok {
											s.SetManagedFields(nil)
											if s.Name != vmwarev1.SupervisorHeadlessSvcName {
												s.Subsets = nil
											}
										}
										return in, nil
									},
								},
							},
						},
					},
				}
			}
			return clustercache.CacheOptions{}
		}(),
		Client: clustercache.ClientOptions{
			QPS:       clusterCacheClientQPS,
			Burst:     clusterCacheClientBurst,
			UserAgent: remote.DefaultClusterAPIUserAgent(controllerName),
			Cache: clustercache.ClientCacheOptions{
				DisableFor: []client.Object{
					// Don't cache ConfigMaps & Secrets.
					&corev1.ConfigMap{},
					&corev1.Secret{},
					// Don't cache Namespaces (used in ProviderServiceAccount controller).
					&corev1.Namespace{},
				},
			},
		},
		WatchFilterValue: managerOpts.WatchFilterValue,
	}, concurrency(clusterCacheConcurrency))
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "Unable to create ClusterCache")
	}

	return clusterCache, nil
}
