/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package generic

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/Azure/azure-service-operator/v2/internal/config"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
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
}

func RegisterWebhooks(mgr ctrl.Manager, objs []client.Object) error {
	var errs []error

	for _, obj := range objs {
		if err := registerWebhook(mgr, obj); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

func registerWebhook(mgr ctrl.Manager, obj client.Object) error {
	_, err := conversion.EnforcePtr(obj)
	if err != nil {
		return errors.Wrap(err, "obj was expected to be ptr but was not")
	}

	return ctrl.NewWebhookManagedBy(mgr).For(obj).Complete()
}

func RegisterAll(
	mgr ctrl.Manager,
	fieldIndexer client.FieldIndexer,
	kubeClient kubeclient.Client,
	positiveConditions *conditions.PositiveConditionBuilder,
	objs []*registration.StorageType,
	options Options) error {

	// pre-register any indexes we need
	for _, obj := range objs {
		for _, indexer := range obj.Indexes {
			options.LogConstructor(nil).V(Info).Info("Registering indexer for type", "type", fmt.Sprintf("%T", obj.Obj), "key", indexer.Key)
			err := fieldIndexer.IndexField(context.Background(), obj.Obj, indexer.Key, indexer.Func)
			if err != nil {
				return errors.Wrapf(err, "failed to register indexer for %T, Key: %q", obj.Obj, indexer.Key)
			}
		}
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
	options Options) error {

	// Use the provided GVK to construct a new runtime object of the desired concrete type.
	gvk, err := apiutil.GVKForObject(info.Obj, mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "creating GVK for obj %T", info)
	}

	loggerFactory := func(mo genruntime.MetaObject) logr.Logger {
		result := options.LogConstructor(nil)
		if options.LoggerFactory != nil {
			if factoryResult := options.LoggerFactory(mo); factoryResult != (logr.Logger{}) && factoryResult != logr.Discard() {
				result = factoryResult
			}
		}

		return result.WithName(info.Name)
	}
	eventRecorder := mgr.GetEventRecorderFor(info.Name)

	options.LogConstructor(nil).V(Status).Info("Registering", "GVK", gvk)

	reconciler := &GenericReconciler{
		Reconciler:                info.Reconciler,
		KubeClient:                kubeClient,
		Config:                    options.Config,
		LoggerFactory:             loggerFactory,
		Recorder:                  eventRecorder,
		GVK:                       gvk,
		PositiveConditions:        positiveConditions,
		RequeueIntervalCalculator: options.RequeueIntervalCalculator,
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		For(info.Obj, ctrlbuilder.WithPredicates(info.Predicate)).
		WithOptions(options.Options)

	for _, watch := range info.Watches {
		builder = builder.Watches(watch.Type, watch.MakeEventHandler(kubeClient, options.LogConstructor(nil).WithName(info.Name)))
	}

	err = builder.Complete(reconciler)
	if err != nil {
		return errors.Wrap(err, "unable to build controllers / reconciler")
	}

	return nil
}
