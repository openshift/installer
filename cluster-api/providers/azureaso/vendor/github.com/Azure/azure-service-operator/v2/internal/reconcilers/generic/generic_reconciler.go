/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package generic

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Azure/azure-service-operator/v2/internal/config"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/util/interval"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

// NamespaceAnnotation defines the annotation name to use when marking
// a resource with the namespace of the managing operator.
const NamespaceAnnotation = "serviceoperator.azure.com/operator-namespace"

type (
	LoggerFactory func(genruntime.MetaObject) logr.Logger
)

// GenericReconciler reconciles resources
type GenericReconciler struct {
	Reconciler                genruntime.Reconciler
	LoggerFactory             LoggerFactory
	KubeClient                kubeclient.Client
	Recorder                  record.EventRecorder
	Config                    config.Values
	GVK                       schema.GroupVersionKind
	PositiveConditions        *conditions.PositiveConditionBuilder
	RequeueIntervalCalculator interval.Calculator
}

var _ reconcile.Reconciler = &GenericReconciler{} // GenericReconciler is a reconcile.Reconciler

// Reconcile will take state in K8s and apply it to Azure
func (gr *GenericReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	metaObj, err := gr.getObjectToReconcile(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if metaObj == nil {
		// This means that the resource doesn't exist
		return ctrl.Result{}, nil
	}

	originalObj := metaObj
	// Ensure that we're always operating on a copy and not on the value returned from the client directly.
	// This is important as it avoids us modifying the cached object.
	metaObj = metaObj.DeepCopyObject().(genruntime.MetaObject)

	log := gr.LoggerFactory(metaObj).WithValues("name", req.Name, "namespace", req.Namespace)
	reconcilers.LogObj(log, Verbose, "Reconcile invoked", metaObj)

	// Ensure the resource is tagged with the operator's namespace.
	ownershipResult, err := gr.takeOwnership(ctx, metaObj)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to take ownership of %s", metaObj.GetName())
	}
	if ownershipResult != nil {
		return *ownershipResult, nil
	}

	var result ctrl.Result
	if !metaObj.GetDeletionTimestamp().IsZero() {
		result, err = gr.delete(ctx, log, metaObj)
	} else {
		result, err = gr.createOrUpdate(ctx, log, metaObj)
	}

	if err != nil {
		err = gr.writeReadyConditionErrorOrDefault(ctx, log, metaObj, err)
		result, err = gr.RequeueIntervalCalculator.NextInterval(req, result, err)
		log.V(Verbose).Info("Encountered error, re-queuing...", "result", result)
		return result, err
	}

	if (result == ctrl.Result{}) {
		// If result is a success, ensure that we note that on Ready condition
		conditions.SetCondition(metaObj, gr.PositiveConditions.Ready.Succeeded(metaObj.GetGeneration()))
	}

	// There are (unfortunately) two ways that an interval can get produced:
	// 1. By this IntervalCalculator
	// 2. By the controller-runtime RateLimiter when a raw ctrl.Result{Requeue: true} is returned, or an error is returned.
	// We used to only have the controller-runtime RateLimiter, but it is quite limited in what information it has access
	// to when generating a backoff. It only has access to the req and a history of how many times that req has been requeued.
	// It doesn't know what error triggered the requeue (or if it was a success).
	// 1m max retry is too aggressive in some cases (see https://github.com/Azure/azure-service-operator/issues/2575),
	// and not aggressive enough in other situations (such as when detecting parent resources have been created, see
	// https://github.com/Azure/azure-service-operator/issues/2556).
	// In order to cater to the above scenarios we calculate some intervals ourselves using this IntervalCalculator and pass others
	// up to the controller-runtime RateLimiter.
	result, err = gr.RequeueIntervalCalculator.NextInterval(req, result, nil)
	if err != nil {
		// This isn't really going to happen but just do it defensively anyway
		return result, err
	}

	// Write the object
	err = gr.CommitUpdate(ctx, log, originalObj, metaObj)
	if err != nil {
		// NotFound is a superfluous error as per https://github.com/kubernetes-sigs/controller-runtime/issues/377
		// The correct handling is just to ignore it and we will get an event shortly with the updated version to patch
		// We do NOT ignore conflict here because it's hard to tell if it's coming from an attempt to update a non-existing resource
		// (see https://github.com/kubernetes/kubernetes/issues/89985), or if it's from an attempt to update a resource which
		// was updated by a user. If we ignore the user-update case, we MIGHT get another event since they changed the resource,
		// but since we don't trigger updates on all changes (some annotations are ignored) we also MIGHT NOT get a fresh event
		// and get stuck. The solution is to let the GET at the top of the controller check for the not-found case and requeue
		// on everything else.
		log.Error(err, "Failed to commit object to etcd")
		return ctrl.Result{}, kubeclient.IgnoreNotFound(err)
	}

	log.V(Verbose).Info("Done with reconcile", "result", result)
	return result, nil
}

func (gr *GenericReconciler) getObjectToReconcile(ctx context.Context, req ctrl.Request) (genruntime.MetaObject, error) {
	obj, err := gr.KubeClient.GetObjectOrDefault(ctx, req.NamespacedName, gr.GVK)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		// This means that the resource doesn't exist
		return nil, nil
	}

	// Always operate on a copy rather than the object from the client, as per
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/controllers.md, which says:
	// Never mutate original objects! Caches are shared across controllers, this means that if you mutate your "copy"
	// (actually a reference or shallow copy) of an object, you'll mess up other controllers (not just your own).
	obj = obj.DeepCopyObject().(client.Object)

	// The Go type for the Kubernetes object must understand how to
	// convert itself to/from the corresponding Azure types.
	metaObj, ok := obj.(genruntime.MetaObject)
	if !ok {
		return nil, errors.Errorf("object is not a genruntime.MetaObject, found type: %T", obj)
	}

	return metaObj, nil
}

func (gr *GenericReconciler) claimResource(ctx context.Context, log logr.Logger, metaObj genruntime.MetaObject) error {
	if !gr.needToAddFinalizer(metaObj) {
		// TODO: This means that if a user messes with some reconciler-specific registration stuff (like owner),
		// TODO: but doesn't remove the finalizer, we won't re-add the reconciler specific stuff. Possibly we should
		// TODO: always re-add that stuff too (it's idempotent)... but then ideally we would avoid a call to Commit
		// TODO: unless it was actually needed?
		return nil
	}

	// Claim the resource
	err := gr.Reconciler.Claim(ctx, log, gr.Recorder, metaObj)
	if err != nil {
		log.Error(err, "Error claiming resource")
		return kubeclient.IgnoreNotFound(err)
	}

	// Adding the finalizer should happen in a reconcile loop prior to the PUT being sent to Azure to avoid situations where
	// we issue a PUT to Azure but the commit of the resource into etcd fails, causing us to have an unset
	// finalizer and have started resource creation in Azure.
	log.V(Info).Info("adding finalizer")
	controllerutil.AddFinalizer(metaObj, genruntime.ReconcilerFinalizer)

	// Passing nil for original here as we know we've made a change and original is only used to determine if the obj
	// has changed to avoid excess commits. In this case, we always need to commit at this stage as adding the finalizer
	// must be persisted to etcd before proceeding.
	err = gr.CommitUpdate(ctx, log, nil, metaObj)
	if err != nil {
		log.Error(err, "Error adding finalizer")
		return kubeclient.IgnoreNotFound(err)
	}

	return nil
}

func (gr *GenericReconciler) needToAddFinalizer(metaObj genruntime.MetaObject) bool {
	unsetFinalizer := !controllerutil.ContainsFinalizer(metaObj, genruntime.ReconcilerFinalizer)
	return unsetFinalizer
}

func (gr *GenericReconciler) createOrUpdate(ctx context.Context, log logr.Logger, metaObj genruntime.MetaObject) (ctrl.Result, error) {
	// Claim the resource
	err := gr.claimResource(ctx, log, metaObj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check the reconcile-policy to ensure we're allowed to issue a CreateOrUpdate
	reconcilePolicy := reconcilers.GetReconcilePolicy(metaObj, log)
	if !reconcilePolicy.AllowsModify() {
		return ctrl.Result{}, gr.handleSkipReconcile(ctx, log, metaObj)
	}

	conditions.SetCondition(metaObj, gr.PositiveConditions.Ready.Reconciling(metaObj.GetGeneration()))

	return gr.Reconciler.CreateOrUpdate(ctx, log, gr.Recorder, metaObj)
}

func (gr *GenericReconciler) delete(ctx context.Context, log logr.Logger, metaObj genruntime.MetaObject) (ctrl.Result, error) {
	// Check the reconcile policy to ensure we're allowed to issue a delete
	reconcilePolicy := reconcilers.GetReconcilePolicy(metaObj, log)
	if !reconcilePolicy.AllowsDelete() {
		log.V(Info).Info("Bypassing delete of resource due to policy", "policy", reconcilePolicy)
		controllerutil.RemoveFinalizer(metaObj, genruntime.ReconcilerFinalizer)
		log.V(Status).Info("Deleted resource")
		return ctrl.Result{}, nil
	}

	// Check if we actually need to issue a delete
	hasFinalizer := controllerutil.ContainsFinalizer(metaObj, genruntime.ReconcilerFinalizer)
	if !hasFinalizer {
		log.Info("Deleted resource")
		return ctrl.Result{}, nil
	}

	result, err := gr.Reconciler.Delete(ctx, log, gr.Recorder, metaObj)
	// If the Delete call had no error and isn't asking us to requeue, then it succeeded and we can remove
	// the finalizer
	if (result == ctrl.Result{} && err == nil) {
		log.V(Info).Info("Delete succeeded, removing finalizer")
		controllerutil.RemoveFinalizer(metaObj, genruntime.ReconcilerFinalizer)
	}

	// TODO: can't set this before the delete call right now due to how ARM resources determine if they need to issue a first delete.
	// TODO: Once I merge a fix to use the async operation for delete polling this can move up to above the Delete call in theory
	conditions.SetCondition(metaObj, gr.PositiveConditions.Ready.Deleting(metaObj.GetGeneration()))

	return result, err
}

// NewRateLimiter creates a new workqueue.Ratelimiter for use controlling the speed of reconciliation.
// It throttles individual requests exponentially and also controls for multiple requests.
func NewRateLimiter(minBackoff time.Duration, maxBackoff time.Duration, limitBurst bool) workqueue.RateLimiter {
	limiters := []workqueue.RateLimiter{
		workqueue.NewItemExponentialFailureRateLimiter(minBackoff, maxBackoff),
	}

	if limitBurst {
		limiters = append(
			limiters,
			// TODO: We could have an azure global (or per subscription) bucket rate limiter to prevent running into subscription
			// TODO: level throttling. For now though just stay with the default that client-go uses.
			// Setting the limiter to 1 every 3 seconds & a burst of 40
			// Based on ARM limits of 1200 puts per hour (20 per minute),
			&workqueue.BucketRateLimiter{
				Limiter: rate.NewLimiter(rate.Limit(0.2), 20),
			})
	}

	return workqueue.NewMaxOfRateLimiter(limiters...)
}

func (gr *GenericReconciler) WriteReadyConditionError(ctx context.Context, log logr.Logger, obj genruntime.MetaObject, err *conditions.ReadyConditionImpactingError) error {
	conditions.SetCondition(obj, gr.PositiveConditions.Ready.ReadyCondition(
		err.Severity,
		obj.GetGeneration(),
		err.Reason,
		err.Cause().Error())) // Don't use err.Error() here because it also includes details about Reason, Severity, which are getting displayed as part of the condition structure
	commitErr := gr.CommitUpdate(ctx, log, nil, obj)
	if commitErr != nil {
		return errors.Wrap(commitErr, "updating resource error")
	}

	return err
}

// takeOwnership marks this resource as owned by this operator. It returns a ctrl.Result ptr to indicate if the result
// should be returned or not. If the result is nil, ownership does not need to be taken
func (gr *GenericReconciler) takeOwnership(ctx context.Context, metaObj genruntime.MetaObject) (*ctrl.Result, error) {
	// Ensure the resource is tagged with the operator's namespace.
	annotations := metaObj.GetAnnotations()
	reconcilerNamespace := annotations[NamespaceAnnotation]
	if reconcilerNamespace != gr.Config.PodNamespace && reconcilerNamespace != "" {
		// We don't want to get into a fight with another operator -
		// so if we see another operator already has this object leave
		// it alone. This will do the right thing in the case of two
		// operators trying to manage the same namespace. It makes
		// moving objects between namespaces or changing which
		// operator owns a namespace fiddlier (since you'd need to
		// remove the annotation) but those operations are likely to
		// be rare.
		message := fmt.Sprintf("Operators in %q and %q are both configured to manage this resource", gr.Config.PodNamespace, reconcilerNamespace)
		gr.Recorder.Event(metaObj, corev1.EventTypeWarning, "Overlap", message)
		return &ctrl.Result{}, nil
	} else if reconcilerNamespace == "" && gr.Config.PodNamespace != "" {
		genruntime.AddAnnotation(metaObj, NamespaceAnnotation, gr.Config.PodNamespace)
		return &ctrl.Result{Requeue: true}, gr.KubeClient.Update(ctx, metaObj)
	}

	return nil, nil
}

func (gr *GenericReconciler) CommitUpdate(ctx context.Context, log logr.Logger, original genruntime.MetaObject, obj genruntime.MetaObject) error {
	if reflect.DeepEqual(original, obj) {
		log.V(Debug).Info("Didn't commit obj as there was no change")
		return nil
	}

	err := gr.KubeClient.CommitObject(ctx, obj)
	if err != nil {
		return err
	}
	reconcilers.LogObj(log, Debug, "updated resource in etcd", obj)
	return nil
}

func (gr *GenericReconciler) handleSkipReconcile(ctx context.Context, log logr.Logger, obj genruntime.MetaObject) error {
	reconcilePolicy := reconcilers.GetReconcilePolicy(obj, log) // TODO: Pull this whole method up here
	log.V(Status).Info(
		"Skipping creation/update of resource due to policy",
		annotations.ReconcilePolicy, reconcilePolicy)

	err := gr.Reconciler.UpdateStatus(ctx, log, gr.Recorder, obj)
	if err != nil {
		return err
	}
	conditions.SetCondition(obj, gr.PositiveConditions.Ready.Succeeded(obj.GetGeneration()))

	return nil
}

func (gr *GenericReconciler) writeReadyConditionErrorOrDefault(ctx context.Context, log logr.Logger, metaObj genruntime.MetaObject, err error) error {
	// If the error in question is NotFound or Conflict from KubeClient just return it right away as there is no reason to wrap it
	if kubeclient.IsNotFoundOrConflict(err) {
		return err
	}

	readyErr, ok := conditions.AsReadyConditionImpactingError(err)
	if !ok {
		// An unknown error, we wrap it as a ready condition error so that the user will always see something, even if
		// the error is generic
		readyErr = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonFailed)
	}

	log.Error(readyErr, "Encountered error impacting Ready condition")
	err = gr.WriteReadyConditionError(ctx, log, metaObj, readyErr)
	return err
}
