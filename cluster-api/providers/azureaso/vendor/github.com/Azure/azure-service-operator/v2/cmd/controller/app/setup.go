/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package app

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/benbjohnson/clock"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/Azure/azure-service-operator/v2/internal/config"
	"github.com/Azure/azure-service-operator/v2/internal/controllers"
	"github.com/Azure/azure-service-operator/v2/internal/crdmanagement"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	asometrics "github.com/Azure/azure-service-operator/v2/internal/metrics"
	armreconciler "github.com/Azure/azure-service-operator/v2/internal/reconcilers/arm"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers/generic"
	asocel "github.com/Azure/azure-service-operator/v2/internal/util/cel"
	"github.com/Azure/azure-service-operator/v2/internal/util/interval"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/internal/util/lockedrand"
	common "github.com/Azure/azure-service-operator/v2/pkg/common/config"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

type Runnable struct {
	mgr manager.Manager

	// toStart must not block
	toStart []func()
}

func (r *Runnable) Start(ctx context.Context) error {
	for _, f := range r.toStart {
		f()
	}

	return r.mgr.Start(ctx)
}

func SetupControllerManager(ctx context.Context, setupLog logr.Logger, flgs *Flags) *Runnable {
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
		Scheme:           scheme,
		NewCache:         cacheFunc,
		LeaderElection:   flgs.EnableLeaderElection,
		LeaderElectionID: "controllers-leader-election-azinfra-generated",
		// It's only safe to set LeaderElectionReleaseOnCancel to true if the manager binary ends
		// when the manager exits. This is the case with us today, so we set this to true whenever
		// flgs.EnableLeaderElection is true.
		LeaderElectionReleaseOnCancel: flgs.EnableLeaderElection,
		HealthProbeBindAddress:        flgs.HealthAddr,
		Metrics:                       getMetricsOpts(flgs),
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

	switch flgs.CRDManagementMode {
	case "auto":
		var goalCRDs []apiextensions.CustomResourceDefinition
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

			// Note that this step will restart the pod when it succeeds
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

	// There are 3 possibilities once we reach here:
	// 1. Webhooks mode + crd-management-mode=auto: existingCRDs will be up to date (upgraded, crd-pattern applied, etc)
	//    by the time we get here as the pod will keep exiting until it is so (see crdManager.ApplyCRDs above).
	// 2. Non-webhooks mode + auto: As outlined in https://azure.github.io/azure-service-operator/guide/authentication/multitenant-deployment/#upgrading
	//    the webhooks mode pod must be upgraded first, so there's not really much practical difference between this and
	//    crd-management-mode=none (see below).
	// 3. crd-management-mode=none: existingCRDs is the set of CRDs that are installed and we can't do anything else but
	//    trust that they are correct.
	//    TODO: This is not quite true as if we wanted we could still read the CRDs from the filesystem and
	//    TODO: just exit if what we see remotely doesn't match what we have locally, the downside of this is we pay
	//    TODO: the nontrivial startup cost of reading the local copy of CRDs into memory. Since "none" is
	//    TODO: us approximating the standard operator experience we don't perform this assertion currently as most
	//    TODO: operators don't.
	readyResources := crdmanagement.MakeCRDMap(existingCRDs)

	if cfg.OperatorMode.IncludesWatchers() {
		//nolint:contextcheck
		err = initializeWatchers(readyResources, cfg, mgr, clients)
		if err != nil {
			setupLog.Error(err, "failed to initialize watchers")
			os.Exit(1)
		}
	}

	if cfg.OperatorMode.IncludesWebhooks() {
		objs := controllers.GetKnownTypes()

		objs, err = crdmanagement.FilterKnownTypesByReadyCRDs(clients.log, scheme, readyResources, objs)
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

	return &Runnable{
		mgr: mgr,
		toStart: []func(){
			// Starts the expression caches. Note that we don't need to stop these we'll
			// let process teardown stop them
			clients.expressionEvaluator.Start,
		},
	}
}

func getMetricsOpts(flags *Flags) server.Options {
	var metricsOptions server.Options

	if flags.SecureMetrics {
		metricsOptions = server.Options{
			BindAddress:    flags.MetricsAddr,
			SecureServing:  true,
			FilterProvider: filters.WithAuthenticationAndAuthorization,
		}
		// Note that pprof endpoints are meant to be sensitive and shouldn't be exposed publicly.
		if flags.ProfilingMetrics {
			metricsOptions.ExtraHandlers = map[string]http.Handler{
				"/debug/pprof/":        http.HandlerFunc(pprof.Index),
				"/debug/pprof/cmdline": http.HandlerFunc(pprof.Cmdline),
				"/debug/pprof/profile": http.HandlerFunc(pprof.Profile),
				"/debug/pprof/symbol":  http.HandlerFunc(pprof.Symbol),
				"/debug/pprof/trace":   http.HandlerFunc(pprof.Trace),
			}
		}
	} else {
		metricsOptions = server.Options{
			BindAddress: flags.MetricsAddr,
		}
	}

	return metricsOptions
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
	expressionEvaluator  asocel.ExpressionEvaluator
	log                  logr.Logger
	options              generic.Options
}

func initializeClients(cfg config.Values, mgr ctrl.Manager) (*clients, error) {
	armMetrics := asometrics.NewARMClientMetrics()
	celMetrics := asometrics.NewCEL()
	asometrics.RegisterMetrics(armMetrics, celMetrics)

	log := ctrl.Log.WithName("controllers")

	credential, err := getDefaultAzureCredential(cfg, log)
	if err != nil {
		return nil, errors.Wrap(err, "error while fetching default global credential")
	}

	kubeClient := kubeclient.NewClient(mgr.GetClient())
	credentialProvider := identity.NewCredentialProvider(credential, kubeClient, nil)

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

	expressionEvaluator, err := asocel.NewExpressionEvaluator(
		asocel.Metrics(celMetrics),
		asocel.Log(log),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating expression evaluator")
	}
	// Register the evaluator for use by webhooks
	asocel.RegisterEvaluator(expressionEvaluator)

	options := makeControllerOptions(log, cfg)

	return &clients{
		positiveConditions:   positiveConditions,
		armConnectionFactory: connectionFactory,
		credentialProvider:   credentialProvider,
		kubeClient:           kubeClient,
		expressionEvaluator:  expressionEvaluator,
		log:                  log,
		options:              options,
	}, nil
}

func initializeWatchers(readyResources map[string]apiextensions.CustomResourceDefinition, cfg config.Values, mgr ctrl.Manager, clients *clients) error {
	clients.log.V(Status).Info("Configuration details", "config", cfg.String())

	objs, err := controllers.GetKnownStorageTypes(
		mgr,
		clients.armConnectionFactory,
		clients.credentialProvider,
		clients.kubeClient,
		clients.positiveConditions,
		clients.expressionEvaluator,
		clients.options)
	if err != nil {
		return errors.Wrap(err, "failed getting storage types and reconcilers")
	}

	// Filter the types to register
	objs, err = crdmanagement.FilterStorageTypesByReadyCRDs(clients.log, mgr.GetScheme(), readyResources, objs)
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
	var additionalRateLimiters []workqueue.TypedRateLimiter[reconcile.Request]
	if cfg.RateLimit.Mode == config.RateLimitModeBucket {
		additionalRateLimiters = append(
			additionalRateLimiters,
			&workqueue.TypedBucketRateLimiter[reconcile.Request]{
				Limiter: rate.NewLimiter(rate.Limit(cfg.RateLimit.QPS), cfg.RateLimit.BucketSize),
			})
	}

	return generic.Options{
		Config: cfg,
		Options: controller.Options{
			MaxConcurrentReconciles: cfg.MaxConcurrentReconciles,
			LogConstructor: func(req *reconcile.Request) logr.Logger {
				// refer to https://github.com/kubernetes-sigs/controller-runtime/pull/1827/files
				if req == nil {
					return log
				}
				// TODO: do we need GVK here too?
				return log.WithValues("namespace", req.Namespace, "name", req.Name)
			},
			// These rate limits are used for happy-path backoffs (for example polling async operation IDs for PUT/DELETE)
			RateLimiter: generic.NewRateLimiter(1*time.Second, 1*time.Minute, additionalRateLimiters...),
		},
		PanicHandler: func() {},
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
