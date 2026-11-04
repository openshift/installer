/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package generic

import (
	"context"
	"fmt"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Azure/azure-service-operator/v2/internal/config"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/util/interval"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/registration"
)

type Options struct {
	controller.Options

	// options specific to our controller
	RequeueIntervalCalculator interval.Calculator
	Config                    config.Values
	LoggerFactory             func(obj metav1.Object) logr.Logger

	PanicHandler func()
}

func RegisterWebhooks(mgr ctrl.Manager, objs []*registration.KnownType) error {
	var errs []error

	for _, obj := range objs {
		if err := registerWebhook(mgr, obj); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func registerWebhook(mgr ctrl.Manager, knownType *registration.KnownType) error {
	_, err := conversion.EnforcePtr(knownType.Obj)
	if err != nil {
		return eris.Wrap(err, "obj was expected to be ptr but was not")
	}

	// Register the webhooks. Note that this is safe to call even if there isn't a defaulter/validator
	// as the NewWebhookManagedBy builder no-ops in the case they're both not set.
	err = ctrl.NewWebhookManagedBy(mgr, knownType.Obj).
		WithCustomDefaulter(knownType.Defaulter).
		WithCustomValidator(knownType.Validator).
		Complete()
	if err != nil {
		return eris.Wrapf(err, "unable to register webhooks for %T", knownType.Obj)
	}

	return nil
}

func RegisterAll(
	mgr ctrl.Manager,
	fieldIndexer client.FieldIndexer,
	kubeClient kubeclient.Client,
	positiveConditions *conditions.PositiveConditionBuilder,
	objs []*registration.StorageType,
	options Options,
) error {
	// pre-register any indexes we need
	for _, obj := range objs {
		for _, indexer := range obj.Indexes {
			mgr.GetLogger().V(Info).Info("Registering indexer for type", "type", fmt.Sprintf("%T", obj.Obj), "key", indexer.Key)
			err := fieldIndexer.IndexField(context.Background(), obj.Obj, indexer.Key, indexer.Func)
			if err != nil {
				return eris.Wrapf(err, "failed to register indexer for %T, Key: %q", obj.Obj, indexer.Key)
			}
		}
	}

	// Start informer for secrets. Not all ASO resources have registered for events from secrets, but we always need the ability
	// to read secrets as at the very least we need to check for the aso-credential secret at the root of the resource namespace.
	// This ensures the cache is populated since ReaderFailOnMissingInformer is true.
	secretGVK := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	}
	mgr.GetLogger().V(Info).Info("Registering informer for type", "type", secretGVK.String())
	// We don't need to block until synced, we just want to make sure the informer is going to start
	if _, err := mgr.GetCache().GetInformerForKind(context.Background(), secretGVK, cache.BlockUntilSynced(false)); err != nil {
		return eris.Wrapf(err, "failed to start informer for secrets")
	}

	// Start informer for configmaps. Not all ASO resources have registered for events from configmaps, but we always need the ability
	// to read configmaps. This ensures the cache is populated since ReaderFailOnMissingInformer is true.
	configMapGVK := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}
	mgr.GetLogger().V(Info).Info("Registering informer for type", "type", configMapGVK.String())
	// We don't need to block until synced, we just want to make sure the informer is going to start
	if _, err := mgr.GetCache().GetInformerForKind(context.Background(), configMapGVK, cache.BlockUntilSynced(false)); err != nil {
		return eris.Wrapf(err, "failed to start informer for configmaps")
	}

	var errs []error
	for _, obj := range objs {
		// TODO: Consider pulling some of the construction of things out of register (gvk, etc), so that we can pass in just
		// TODO: the applicable extensions rather than a map of all of them
		if err := register(mgr, kubeClient, positiveConditions, obj, options); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func register(
	mgr ctrl.Manager,
	kubeClient kubeclient.Client,
	positiveConditions *conditions.PositiveConditionBuilder,
	info *registration.StorageType,
	options Options,
) error {
	// Use the provided GVK to construct a new runtime object of the desired concrete type.
	gvk, err := apiutil.GVKForObject(info.Obj, mgr.GetScheme())
	if err != nil {
		return eris.Wrapf(err, "creating GVK for obj %T", info)
	}

	loggerFactory := func(mo genruntime.MetaObject) logr.Logger {
		result := mgr.GetLogger()
		if options.LoggerFactory != nil {
			if factoryResult := options.LoggerFactory(mo); factoryResult != (logr.Logger{}) && factoryResult != logr.Discard() {
				result = factoryResult
			}
		}

		return result.WithName(info.Name)
	}
	eventRecorder := mgr.GetEventRecorderFor(info.Name)

	mgr.GetLogger().V(Status).Info("Registering", "GVK", gvk)

	reconciler := &GenericReconciler{
		Reconciler:                info.Reconciler,
		KubeClient:                kubeClient,
		Config:                    options.Config,
		LoggerFactory:             loggerFactory,
		Recorder:                  eventRecorder,
		GVK:                       gvk,
		PositiveConditions:        positiveConditions,
		RequeueIntervalCalculator: options.RequeueIntervalCalculator,
		PanicHandler:              options.PanicHandler,
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		For(info.Obj, ctrlbuilder.WithPredicates(info.Predicate)).
		WithOptions(options.Options)
	builder.Named(info.Name)

	// All resources watch namespace for the reconcile-policy annotation
	builder = builder.Watches(
		&corev1.Namespace{},
		handler.EnqueueRequestsFromMapFunc(registerNamespaceWatcher(kubeClient, gvk)),
		ctrlbuilder.WithPredicates(reconcilers.ARMReconcilerAnnotationChangedPredicate()),
	)

	for _, watch := range info.Watches {
		builder = builder.Watches(watch.Type, watch.MakeEventHandler(kubeClient, mgr.GetLogger().WithName(info.Name)))
	}

	err = builder.Complete(reconciler)
	if err != nil {
		return eris.Wrap(err, "unable to build controllers / reconciler")
	}

	return nil
}

// This registers a watcher on the namespace for the watched resource
func registerNamespaceWatcher(kubeclient kubeclient.Client, gvk schema.GroupVersionKind) func(ctx context.Context, obj client.Object) []reconcile.Request {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		reconcileRequests := []reconcile.Request{}
		log := log.FromContext(ctx)

		aList := &unstructured.UnstructuredList{}
		aList.SetGroupVersionKind(gvk)
		if err := kubeclient.List(ctx, aList, &client.ListOptions{Namespace: obj.GetName()}); err != nil {
			return []reconcile.Request{}
		}
		// list the objects for the current kind

		log.V(Verbose).Info("Detected namespace reconcile-policy annotation change",
			"namespace", obj.GetName(),
			"reconcile-policy", obj.GetAnnotations()["serviceoperator.azure.com/reconcile-policy"])
		for _, el := range aList.Items {
			if _, ok := el.GetAnnotations()["serviceoperator.azure.com/reconcile-policy"]; ok {
				// If the annotation is defined for the object, there's no need to reconcile it and we skip the object
				continue
			}
			if _, ok := obj.GetAnnotations()["serviceoperator.azure.com/reconcile-policy"]; ok {
				reconcileRequests = append(reconcileRequests, reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      el.GetName(),
					Namespace: el.GetNamespace(),
				}})
			}
		}
		return reconcileRequests
	}
}
