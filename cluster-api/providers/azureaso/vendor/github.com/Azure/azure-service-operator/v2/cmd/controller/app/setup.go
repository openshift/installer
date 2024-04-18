/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/benbjohnson/clock"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/Azure/azure-service-operator/v2/api"
	"github.com/Azure/azure-service-operator/v2/internal/config"
	"github.com/Azure/azure-service-operator/v2/internal/controllers"
	"github.com/Azure/azure-service-operator/v2/internal/crdmanagement"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	asometrics "github.com/Azure/azure-service-operator/v2/internal/metrics"
	armreconciler "github.com/Azure/azure-service-operator/v2/internal/reconcilers/arm"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers/generic"
	"github.com/Azure/azure-service-operator/v2/internal/util/interval"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/internal/util/lockedrand"
	common "github.com/Azure/azure-service-operator/v2/pkg/common/config"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

func SetupPreUpgradeCheck(ctx context.Context) error {
	cfg, err := clientconfig.GetConfig()
	if err != nil {
		return errors.Wrap(err, "unable to get client config")
	}

	apiExtClient, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create kubernetes client")
	}

	// Had to list CRDs this way and not with crdManager, since we did not have "serviceoperator.azure.com/version" labels in earlier versions.
	list, err := apiExtClient.CustomResourceDefinitions().List(ctx, v1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list CRDs")
	}

	scheme := api.CreateScheme()
	crdRegexp := regexp.MustCompile(`.*\.azure\.com`)
	var errs []error
	for _, crd := range list.Items {
		crd := crd
		if !crdRegexp.MatchString(crd.Name) {
			continue
		}

		if !scheme.Recognizes(crd.GroupVersionKind()) {
			// Not one of our resources
			continue
		}

		// If this CRD is annotated with "serviceoperator.azure.com/version", it must be >=2.0.0 and so safe
		// as we didn't start using this label until 2.0.0. Same with "app.kubernetes.io/version" which was added in 2.3.0
		// in favor of our custom serviceoperator.azure.com
		_, hasOldLabel := crd.Labels[crdmanagement.ServiceOperatorVersionLabelOld]
		_, hasNewLabel := crd.Labels[crdmanagement.ServiceOperatorVersionLabel]
		if hasOldLabel || hasNewLabel {
			continue
		}

		if policy, ok := crd.Annotations["helm.sh/resource-policy"]; !ok || policy != "keep" {
			err = errors.New(fmt.Sprintf("CRD '%s' does not have annotation for helm keep policy. Make sure the upgrade is from beta.5", crd.Name))
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func SetupControllerManager(ctx context.Context, setupLog logr.Logger, flgs Flags) manager.Manager {
	scheme := controllers.CreateScheme()
	_ = apiextensions.AddToScheme(scheme) // Used for managing CRDs

	cfg, err := config.ReadAndValidate()
	if err != nil {
		setupLog.Error(err, "unable to get env configuration values")
		os.Exit(1)
	}

	var cacheFunc cache.NewCacheFunc
	if cfg.TargetNamespaces != nil && cfg.OperatorMode.IncludesWatchers() {
		cacheFunc = func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = make(map[string]cache.Config, len(cfg.TargetNamespaces))
			for _, ns := range cfg.TargetNamespaces {
				opts.DefaultNamespaces[ns] = cache.Config{}
			}

			return cache.New(config, opts)
		}
	}

	k8sConfig := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(k8sConfig, ctrl.Options{
		Scheme:                 scheme,
		NewCache:               cacheFunc,
		LeaderElection:         flgs.EnableLeaderElection,
		LeaderElectionID:       "controllers-leader-election-azinfra-generated",
		HealthProbeBindAddress: flgs.HealthAddr,
		Metrics: server.Options{
			BindAddress: flgs.MetricsAddr,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    flgs.WebhookPort,
			CertDir: flgs.WebhookCertDir,
		}),
	})
	if err != nil {
		setupLog.Error(err, "unable to create manager")
		os.Exit(1)
	}

	clients, err := initializeClients(cfg, mgr)
	if err != nil {
		setupLog.Error(err, "failed to initialize clients")
		os.Exit(1)
	}

	// TODO: Put all of the CRD stuff into a method?
	crdManager, err := newCRDManager(clients.log, mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "failed to initialize CRD client")
		os.Exit(1)
	}
	existingCRDs, err := crdManager.ListOperatorCRDs(ctx)
	if err != nil {
		setupLog.Error(err, "failed to list current CRDs")
		os.Exit(1)
	}

	// By default, assume the existing CRDs are the goal CRDs. If CRD management is enabled, we will
	// load the goal CRDs from disk and apply them.
	goalCRDs := existingCRDs
	switch flgs.CRDManagementMode {
	case "auto":
		goalCRDs, err = crdManager.LoadOperatorCRDs(crdmanagement.CRDLocation, cfg.PodNamespace)
		if err != nil {
			setupLog.Error(err, "failed to load CRDs from disk")
			os.Exit(1)
		}

		// We only apply CRDs if we're in webhooks mode. No other mode will have CRD CRUD permissions
		if cfg.OperatorMode.IncludesWebhooks() {
			var installationInstructions []*crdmanagement.CRDInstallationInstruction
			installationInstructions, err = crdManager.DetermineCRDsToInstallOrUpgrade(goalCRDs, existingCRDs, flgs.CRDPatterns)
			if err != nil {
				setupLog.Error(err, "failed to determine CRDs to apply")
				os.Exit(1)
			}

			included := crdmanagement.IncludedCRDs(installationInstructions)
			if len(included) == 0 {
				err = errors.New("No existing CRDs in cluster and no --crd-pattern specified")
				setupLog.Error(err, "failed to apply CRDs")
				os.Exit(1)
			}

			err = crdManager.ApplyCRDs(ctx, installationInstructions)
			if err != nil {
				setupLog.Error(err, "failed to apply CRDs")
				os.Exit(1)
			}
		}
	case "none":
		setupLog.Info("CRD management mode was set to 'none', the operator will not manage CRDs and assumes they are already installed and matching the operator version")
	default:
		setupLog.Error(fmt.Errorf("invalid CRD management mode: %s", flgs.CRDManagementMode), "failed to initialize CRD client")
		os.Exit(1)
	}

	// Of all the resources we know of, find any that aren't ready. We will use this collection
	// to skip watching of these not-ready resources.
	nonReadyResources := crdmanagement.GetNonReadyCRDs(cfg, crdManager, goalCRDs, existingCRDs)

	if cfg.OperatorMode.IncludesWatchers() {
		//nolint:contextcheck
		err = initializeWatchers(nonReadyResources, cfg, mgr, clients)
		if err != nil {
			setupLog.Error(err, "failed to initialize watchers")
			os.Exit(1)
		}
	}

	if cfg.OperatorMode.IncludesWebhooks() {
		objs := controllers.GetKnownTypes()

		objs, err = crdmanagement.FilterKnownTypesByReadyCRDs(clients.log, scheme, nonReadyResources, objs)
		if err != nil {
			setupLog.Error(err, "failed to filter known types by ready CRDs")
			os.Exit(1)
		}

		if errs := generic.RegisterWebhooks(mgr, objs); errs != nil {
			setupLog.Error(err, "failed to register webhook for gvks")
			os.Exit(1)
		}
	}

	// Healthz liveness probe endpoint
	err = mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Failed setting up healthz check")
		os.Exit(1)
	}

	// Readyz probe endpoint
	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Failed setting up readyz check")
		os.Exit(1)
	}

	// Readyz probe endpoint
	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		setupLog.Error(err, "Failed setting up readyz check")
		os.Exit(1)
	}
	return mgr
}

func getDefaultAzureCredential(cfg config.Values, setupLog logr.Logger) (*identity.Credential, error) {
	tokenCred, err := getDefaultAzureTokenCredential(cfg, setupLog)
	if err != nil {
		return nil, err
	}
	if tokenCred == nil {
		return nil, nil
	}

	return identity.NewDefaultCredential(
		tokenCred,
		cfg.PodNamespace,
		cfg.SubscriptionID), nil
}

func getDefaultAzureTokenCredential(cfg config.Values, setupLog logr.Logger) (azcore.TokenCredential, error) {
	// If subscriptionID is not supplied, then set default credential to not be used/nil
	if cfg.SubscriptionID == "" {
		setupLog.Info("No global credential configured, continuing without default global credential.")
		return nil, nil
	}

	if cfg.UseWorkloadIdentityAuth {
		credential, err := azidentity.NewWorkloadIdentityCredential(&azidentity.WorkloadIdentityCredentialOptions{
			ClientID:      cfg.ClientID,
			TenantID:      cfg.TenantID,
			TokenFilePath: identity.FederatedTokenFilePath,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get workload identity credential")
		}

		return credential, nil
	}

	if cert := os.Getenv(common.AzureClientCertificate); cert != "" {
		certPassword := os.Getenv(common.AzureClientCertificatePassword)
		credential, err := identity.NewClientCertificateCredential(cfg.TenantID, cfg.ClientID, []byte(cert), []byte(certPassword))
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get client certificate credential")
		}

		return credential, nil
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get default azure credential")
	}

	return credential, err
}

type clients struct {
	positiveConditions   *conditions.PositiveConditionBuilder
	armConnectionFactory armreconciler.ARMConnectionFactory
	credentialProvider   identity.CredentialProvider
	kubeClient           kubeclient.Client
	log                  logr.Logger
	options              generic.Options
}

func initializeClients(cfg config.Values, mgr ctrl.Manager) (*clients, error) {
	armMetrics := asometrics.NewARMClientMetrics()
	asometrics.RegisterMetrics(armMetrics)

	log := ctrl.Log.WithName("controllers")

	credential, err := getDefaultAzureCredential(cfg, log)
	if err != nil {
		return nil, errors.Wrap(err, "error while fetching default global credential")
	}

	kubeClient := kubeclient.NewClient(mgr.GetClient())
	credentialProvider := identity.NewCredentialProvider(credential, kubeClient)

	armClientCache := armreconciler.NewARMClientCache(
		credentialProvider,
		kubeClient,
		cfg.Cloud(),
		nil,
		armMetrics)

	genericarmclient.AddToUserAgent(cfg.UserAgentSuffix)

	var connectionFactory armreconciler.ARMConnectionFactory = func(ctx context.Context, obj genruntime.ARMMetaObject) (armreconciler.Connection, error) {
		return armClientCache.GetConnection(ctx, obj)
	}

	positiveConditions := conditions.NewPositiveConditionBuilder(clock.New())

	options := makeControllerOptions(log, cfg)

	return &clients{
		positiveConditions:   positiveConditions,
		armConnectionFactory: connectionFactory,
		credentialProvider:   credentialProvider,
		kubeClient:           kubeClient,
		log:                  log,
		options:              options,
	}, nil
}

func initializeWatchers(nonReadyResources map[string]apiextensions.CustomResourceDefinition, cfg config.Values, mgr ctrl.Manager, clients *clients) error {
	clients.log.V(Status).Info("Configuration details", "config", cfg.String())

	objs, err := controllers.GetKnownStorageTypes(
		mgr,
		clients.armConnectionFactory,
		clients.credentialProvider,
		clients.kubeClient,
		clients.positiveConditions,
		clients.options)
	if err != nil {
		return errors.Wrap(err, "failed getting storage types and reconcilers")
	}

	// Filter the types to register
	objs, err = crdmanagement.FilterStorageTypesByReadyCRDs(clients.log, mgr.GetScheme(), nonReadyResources, objs)
	if err != nil {
		return errors.Wrap(err, "failed to filter storage types by ready CRDs")
	}

	err = generic.RegisterAll(
		mgr,
		mgr.GetFieldIndexer(),
		clients.kubeClient,
		clients.positiveConditions,
		objs,
		clients.options)
	if err != nil {
		return errors.Wrap(err, "failed to register gvks")
	}

	return nil
}

func makeControllerOptions(log logr.Logger, cfg config.Values) generic.Options {
	return generic.Options{
		Config: cfg,
		Options: controller.Options{
			MaxConcurrentReconciles: 1,
			LogConstructor: func(req *reconcile.Request) logr.Logger {
				// refer to https://github.com/kubernetes-sigs/controller-runtime/pull/1827/files
				if req == nil {
					return log
				}
				// TODO: do we need GVK here too?
				return log.WithValues("namespace", req.Namespace, "name", req.Name)
			},
			// These rate limits are used for happy-path backoffs (for example polling async operation IDs for PUT/DELETE)
			RateLimiter: generic.NewRateLimiter(1*time.Second, 1*time.Minute, true),
		},
		RequeueIntervalCalculator: interval.NewCalculator(
			// These rate limits are primarily for ReadyConditionImpactingError's
			interval.CalculatorParameters{
				//nolint:gosec // do not want cryptographic randomness here
				Rand:              rand.New(lockedrand.NewSource(time.Now().UnixNano())),
				ErrorBaseDelay:    1 * time.Second,
				ErrorMaxFastDelay: 30 * time.Second,
				ErrorMaxSlowDelay: 3 * time.Minute,
				SyncPeriod:        cfg.SyncPeriod,
			}),
	}
}

func newCRDManager(logger logr.Logger, k8sConfig *rest.Config) (*crdmanagement.Manager, error) {
	crdScheme := runtime.NewScheme()
	_ = apiextensions.AddToScheme(crdScheme)
	crdClient, err := client.New(k8sConfig, client.Options{Scheme: crdScheme})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create CRD client")
	}

	crdManager := crdmanagement.NewManager(logger, kubeclient.NewClient(crdClient))
	return crdManager, nil
}
