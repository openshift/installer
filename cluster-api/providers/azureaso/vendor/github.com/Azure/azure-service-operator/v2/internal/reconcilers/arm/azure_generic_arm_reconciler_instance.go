/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers/arm/errorclassification"
	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	"github.com/Azure/azure-service-operator/v2/pkg/common/labels"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/merger"
)

type azureDeploymentReconcilerInstance struct {
	reconcilers.ARMOwnedResourceReconcilerCommon
	Obj           genruntime.ARMMetaObject
	Log           logr.Logger
	Recorder      record.EventRecorder
	Extension     genruntime.ResourceExtension
	ARMConnection Connection
}

func newAzureDeploymentReconcilerInstance(
	metaObj genruntime.ARMMetaObject,
	log logr.Logger,
	recorder record.EventRecorder,
	connection Connection,
	reconciler AzureDeploymentReconciler,
) *azureDeploymentReconcilerInstance {
	return &azureDeploymentReconcilerInstance{
		Obj:                              metaObj,
		Log:                              log,
		Recorder:                         recorder,
		ARMConnection:                    connection,
		Extension:                        reconciler.Extension,
		ARMOwnedResourceReconcilerCommon: reconciler.ARMOwnedResourceReconcilerCommon,
	}
}

func (r *azureDeploymentReconcilerInstance) CreateOrUpdate(ctx context.Context) (ctrl.Result, error) {
	action, actionFunc, err := r.DetermineCreateOrUpdateAction()
	if err != nil {
		r.Log.Error(err, "error determining create or update action")
		r.Recorder.Event(r.Obj, v1.EventTypeWarning, "DetermineCreateOrUpdateActionError", err.Error())

		return ctrl.Result{}, err
	}

	r.Log.V(Verbose).Info("Determined CreateOrUpdate action", "action", action)

	result, err := actionFunc(ctx)
	if err != nil {
		r.Recorder.Event(r.Obj, v1.EventTypeWarning, "CreateOrUpdateActionError", err.Error())

		return ctrl.Result{}, err
	}

	return result, nil
}

func (r *azureDeploymentReconcilerInstance) Delete(ctx context.Context) (ctrl.Result, error) {
	action, actionFunc, err := r.DetermineDeleteAction()
	if err != nil {
		r.Recorder.Event(r.Obj, v1.EventTypeWarning, "DetermineDeleteActionError", err.Error())

		return ctrl.Result{}, err
	}

	r.Log.V(Verbose).Info("Determined Delete action", "action", action)

	result, err := actionFunc(ctx)
	if err != nil {
		r.Recorder.Event(r.Obj, v1.EventTypeWarning, "DeleteActionError", err.Error())

		return ctrl.Result{}, err
	}

	return result, nil
}

func (r *azureDeploymentReconcilerInstance) MakeReadyConditionImpactingErrorFromError(azureErr error) error {
	apiVersion, verr := genruntime.GetAPIVersion(r.Obj, r.ResourceResolver.Scheme())
	if verr != nil {
		return eris.Wrapf(verr, "error getting api version for resource %s while making Ready condition", r.Obj.GetName())
	}
	classifier := extensions.CreateErrorClassifier(r.Extension, errorclassification.ClassifyCloudError, apiVersion, r.Log)
	return errorclassification.MakeReadyConditionImpactingErrorFromError(azureErr, classifier)
}

func (r *azureDeploymentReconcilerInstance) AddInitialResourceState(ctx context.Context) error {
	armResource, err := r.ConvertResourceToARMResource(ctx)
	if err != nil {
		return err
	}
	genruntime.SetResourceID(r.Obj, armResource.GetID())
	labels.SetOwnerNameLabel(r.Log, r.Obj)
	labels.SetOwnerGroupKindLabel(r.Log, r.Obj)
	labels.SetOwnerUIDLabel(r.Obj)
	return nil
}

func (r *azureDeploymentReconcilerInstance) DetermineDeleteAction() (DeleteAction, DeleteActionFunc, error) {
	pollerID, _, hasPollerResumeToken := GetPollerResumeToken(r.Obj)

	if hasPollerResumeToken && pollerID == genericarmclient.DeletePollerID {
		return DeleteActionMonitorDelete, r.MonitorDelete, nil
	}

	if !genruntime.ResourceOperationDelete.IsSupportedBy(r.Obj) {
		// Resource doesn't support delete; we'll end up returning an actionable error
		return DeleteActionNotPossibleInAzure, r.DeleteNotPossibleInAzure, nil
	}

	return DeleteActionBeginDelete, r.StartDeleteOfResource, nil
}

func (r *azureDeploymentReconcilerInstance) DetermineCreateOrUpdateAction() (CreateOrUpdateAction, CreateOrUpdateActionFunc, error) {
	ready := genruntime.GetReadyCondition(r.Obj)
	_, _, hasPollerResumeToken := GetPollerResumeToken(r.Obj)

	if ready != nil && ready.Reason == conditions.ReasonDeleting.Name {
		return CreateOrUpdateActionNoAction, NoAction, eris.Errorf("resource is currently deleting; it can not be applied")
	}

	if hasPollerResumeToken {
		return CreateOrUpdateActionMonitorCreation, r.MonitorResourceCreation, nil
	}

	return CreateOrUpdateActionBeginCreation, r.BeginCreateOrUpdateResource, nil
}

//////////////////////////////////////////
// Actions
//////////////////////////////////////////

func NoAction(_ context.Context) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// StartDeleteOfResource will begin deletion of a resource by telling Azure to start deleting it. The resource will be
// marked with the provisioning state of "Deleting".
func (r *azureDeploymentReconcilerInstance) StartDeleteOfResource(ctx context.Context) (ctrl.Result, error) {
	msg := "Starting delete of resource"
	r.Log.V(Status).Info(msg)
	r.Recorder.Event(r.Obj, v1.EventTypeNormal, string(DeleteActionBeginDelete), msg)

	deleter := extensions.CreateDeleter(r.Extension, r.deleteResource)
	result, err := deleter(ctx, r.Log, r.ResourceResolver, r.ARMConnection.Client(), r.Obj)
	return result, err
}

// MonitorDelete will call Azure to check if the resource still exists. If so, it will requeue, else,
// the finalizer will be removed.
func (r *azureDeploymentReconcilerInstance) MonitorDelete(ctx context.Context) (ctrl.Result, error) {
	msg := "Continue monitoring deletion"
	r.Log.V(Verbose).Info(msg)
	r.Recorder.Event(r.Obj, v1.EventTypeNormal, string(DeleteActionMonitorDelete), msg)
	//
	//// Technically we don't need the resource ID anymore to monitor delete
	//_, hasResourceID := genruntime.GetResourceID(r.Obj)
	//if !hasResourceID {
	//	return ctrl.Result{}, errors.Errorf("can't MonitorDelete a resource without a resource ID")
	//}

	pollerID, pollerResumeToken, hasToken := GetPollerResumeToken(r.Obj)
	if !hasToken {
		return ctrl.Result{}, eris.New("cannot MonitorResourceCreation with empty pollerResumeToken or pollerID")
	}

	if pollerID != genericarmclient.DeletePollerID {
		return ctrl.Result{}, eris.Errorf("cannot MonitorResourceCreation with pollerID=%s", pollerID)
	}

	poller := r.ARMConnection.Client().ResumeDeletePoller(pollerID)
	err := poller.Resume(ctx, r.ARMConnection.Client(), pollerResumeToken)
	if err != nil {
		return ctrl.Result{}, r.handleDeleteFailed(err)
	}
	if poller.Poller.Done() {
		// The resource was deleted
		return ctrl.Result{}, nil
	}

	retryAfter := genericarmclient.GetRetryAfter(poller.RawResponse)
	r.Log.V(Verbose).Info("Found resource: continuing to wait for deletion...")
	// Normally don't need to set both of these fields but because retryAfter can be 0 we do
	return ctrl.Result{Requeue: true, RequeueAfter: retryAfter}, nil
}

// DeleteNotPossibleInAzure is used when the underlying Azure resource doesn't support direct
// deletion, so we return an error unless the resource has already gone.
func (r *azureDeploymentReconcilerInstance) DeleteNotPossibleInAzure(ctx context.Context) (ctrl.Result, error) {
	resourceID, hasResourceID := genruntime.GetResourceID(r.Obj)
	if !hasResourceID {
		// No resource ID means nothing to delete
		return ctrl.Result{}, nil
	}

	_, _, err := r.getStatus(ctx, r.Obj, resourceID)
	if err != nil && genericarmclient.IsNotFoundError(err) {
		// Resource no longer exists
		return ctrl.Result{}, nil
	}

	msg := fmt.Sprintf(
		"Resource does not support deletion in Azure; set annotation '%s: %s' to permit deletion in Kubernetes",
		annotations.ReconcilePolicy,
		annotations.ReconcilePolicyDetachOnDelete)
	r.Log.V(Verbose).Info(msg)
	r.Recorder.Event(r.Obj, v1.EventTypeNormal, string(DeleteActionNotPossibleInAzure), msg)

	// Return a meaningful error so that the Ready condition is updated to show the user why the resource can't yet be deleted.
	if err == nil {
		err = eris.New(msg)
	} else {
		err = eris.Wrap(err, msg)
	}

	return ctrl.Result{},
		conditions.NewReadyConditionImpactingError(
			err,
			conditions.ConditionSeverityWarning,
			conditions.ReasonDeletionNotSupported)
}

func (r *azureDeploymentReconcilerInstance) BeginCreateOrUpdateResource(
	ctx context.Context,
) (ctrl.Result, error) {
	if r.Obj.AzureName() == "" {
		err := eris.Errorf(
			"AzureName was not set on %s. A webhook should default this to .metadata.name if it was omitted. Is the ASO webhook service running?",
			r.Obj.GetType())

		return ctrl.Result{},
			conditions.NewReadyConditionImpactingError(
				err, conditions.ConditionSeverityError, conditions.ReasonFailed)
	}

	// We want to set the latest reconciled generation annotation to keep a track of reconciles per generation.
	SetLatestReconciledGeneration(r.Obj)

	check, err := r.preReconciliationCheck(ctx)
	if err != nil {
		// Failed to do the pre-reconciliation check, this is a serious but non-fatal error
		// Make sure we return a ReadyConditionImpactingError so that the Ready condition is updated for the user
		// Ideally any implementation of the checker should return a ReadyConditionImpactingError, but we can't
		// guarantee that, so we wrap as required
		impactingError, ok := conditions.AsReadyConditionImpactingError(err)
		if !ok {
			impactingError = conditions.NewReadyConditionImpactingError(
				err,
				conditions.ConditionSeverityWarning,
				conditions.ReasonFailed)
		}

		return ctrl.Result{}, impactingError
	}

	// If the check says we're postponing reconcile, we're done for now as there's nothing to do.
	if check.PostponeReconciliation() {
		r.Log.V(Status).Info("Extension recommended postponing reconciliation")
		return ctrl.Result{}, nil
	}

	// If the check says we're blocking reconcile, we return ReadyConditionImpactingError here to update the Ready
	// condition is updated so the user can see why we're not reconciling right now, and to trigger a retry in a bit.
	if check.BlockReconciliation() {
		r.Log.V(Status).Info("Extension recommended blocking reconciliation", "message", check.Message())
		return ctrl.Result{}, check.CreateConditionError()
	}

	resourceID := genruntime.GetResourceIDOrDefault(r.Obj)
	if resourceID != "" {
		err = r.checkSubscription(resourceID)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	armResource, err := r.ConvertResourceToARMResource(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Use conditions.SetConditionReasonAware here to override any Warning conditions set earlier in the reconciliation process.
	// Note that this call should be done after all validation has passed and all that is left to do is send the payload to ARM.
	conditions.SetConditionReasonAware(r.Obj, r.PositiveConditions.Ready.Reconciling(r.Obj.GetGeneration()))

	r.Log.V(Status).Info("About to send resource to Azure")

	// Try to create the resource
	spec := armResource.Spec()
	pollerResp, err := r.ARMConnection.Client().BeginCreateOrUpdateByID(ctx, armResource.GetID(), spec.GetAPIVersion(), spec)
	if err != nil {
		return ctrl.Result{}, r.handleCreateOrUpdateFailed(err)
	}

	r.Log.V(Status).Info("Successfully sent resource to Azure", "id", armResource.GetID())
	r.Recorder.Eventf(r.Obj, v1.EventTypeNormal, string(CreateOrUpdateActionBeginCreation), "Successfully sent resource to Azure with ID %q", armResource.GetID())

	// If we are done here it means the deployment succeeded immediately. It can't have failed because if it did
	// we would have taken the error path above.
	if pollerResp.Poller.Done() {
		return ctrl.Result{}, r.handleCreateOrUpdateSuccess(ctx, ManageResource)
	}

	resumeToken, err := pollerResp.Poller.ResumeToken()
	if err != nil {
		return ctrl.Result{},
			eris.Wrapf(err, "couldn't create PUT resume token for resource %q", armResource.GetID())
	}

	SetPollerResumeToken(r.Obj, pollerResp.ID, resumeToken)
	return ctrl.Result{Requeue: true}, nil
}

func (r *azureDeploymentReconcilerInstance) preReconciliationCheck(
	ctx context.Context,
) (extensions.PreReconcileCheckResult, error) {
	// Check to see which extensions are available
	checker, extensionFound := extensions.CreatePreReconciliationChecker(r.Extension)
	ownerChecker, ownerExtensionFound := extensions.CreatePreReconciliationOwnerChecker(r.Extension)

	if !extensionFound && !ownerExtensionFound {
		// No extensions found, nothing to do
		return extensions.ProceedWithReconcile(), nil
	}

	// Load owner details so it has an an up-to-date status
	ownerDetails, ownerErr := r.ResourceResolver.ResolveOwner(ctx, r.Obj)
	if ownerErr != nil {
		// We can't obtain the owner, so we can't run either extension
		return extensions.PreReconcileCheckResult{}, ownerErr
	}

	// Run the PreReconciliationOwnerChecker if we have one
	if ownerExtensionFound {
		var ownerObj genruntime.ARMMetaObject

		if ownerDetails.Owner != nil {
			// Owner is a Kubernetes resource - update its status
			// Note that this update is not committed back to api-server currently, so if we come back around here
			// and run through this extension again, we will re-fetch the owner status from Azure again.
			err := r.updateStatus(ctx, ownerDetails.Owner)
			if err != nil {
				return extensions.PreReconcileCheckResult{}, err
			}
			ownerObj = ownerDetails.Owner
		} else if ownerDetails.Result == resolver.OwnerFoundARM {
			// Owner is an ARM ID - create a temporary object and fetch status

			// The version retrieved here is always the latest storage version of the owner type
			// (not necessarily the same version as the child resource).
			ownerGroup, ownerKind := genruntime.LookupOwnerGroupKind(r.Obj.GetSpec())
			groupKind := schema.GroupKind{Group: ownerGroup, Kind: ownerKind}
			ownerGVK, err := r.ResourceResolver.FindGVKForGroupKind(groupKind)
			if err != nil {
				return extensions.PreReconcileCheckResult{}, eris.Wrapf(err, "finding GVK for owner")
			}

			// Fetch the owner status from ARM
			ownerObj, err = r.getStatusFromARMID(ctx, ownerDetails.ARMID, ownerGVK)
			if err != nil {
				return extensions.PreReconcileCheckResult{}, eris.Wrapf(err, "getting status for ARM owner %s", ownerDetails.ARMID)
			}
		}

		if ownerObj != nil {
			check, checkErr := ownerChecker(ctx, ownerObj, r.ResourceResolver, r.ARMConnection.Client(), r.Log)
			if checkErr != nil {
				// Something went wrong running the check.
				return extensions.PreReconcileCheckResult{}, checkErr
			}

			// If the check says we're postponing (because a reconcile is not needed) or blocking (because we can't
			// reconcile at all), we're done for now as there's nothing to do.
			if check.PostponeReconciliation() || check.BlockReconciliation() {
				return check, nil
			}
		}
	}

	// Load resource details so it also has an up-to-date status
	// We defer this until after we've done the owner check as if that says postpone
	// we don't need to do this work.
	// Plus, this avoids errors if the owner is in a state where going a GET on the resource will fail.
	statusErr := r.updateStatus(ctx, r.Obj)
	if statusErr != nil && !genericarmclient.IsNotFoundError(statusErr) {
		// We have an error, and it's not because the resource doesn't exist yet
		return extensions.PreReconcileCheckResult{}, statusErr
	}

	// Run our pre-reconciliation checker
	check, checkErr := checker(ctx, r.Obj, r.ResourceResolver, r.ARMConnection.Client(), r.Log)
	if checkErr != nil {
		// Something went wrong running the check.
		return extensions.PreReconcileCheckResult{}, checkErr
	}

	return check, nil
}

// checkSubscription checks if subscription on resource matches with credentials used while creating a resource.
// Which prevents users to modify subscription in their credential.
func (r *azureDeploymentReconcilerInstance) checkSubscription(resourceID string) error {
	// Some resources like '/providers/Microsoft.Subscription/aliases' do not have subscriptionID,
	// so we need to make sure subscriptionID exists before we check.
	parsedRID, err := arm.ParseResourceID(resourceID)
	if err != nil {
		// We never expect the resource ID to be invalid, so return an error here to avoid someone
		// mangling a resource annotation and bypassing the check
		err = eris.Wrapf(err, "parsing resource ID %q", resourceID)

		return conditions.NewReadyConditionImpactingError(
			err, conditions.ConditionSeverityError, conditions.ReasonFailed)
	}

	if !genruntime.CheckARMIDMatchesSubscription(r.ARMConnection.SubscriptionID(), parsedRID) {
		err = eris.Errorf(
			"SubscriptionID %q for %q resource does not match with Client Credential: %q",
			parsedRID.SubscriptionID,
			resourceID,
			r.ARMConnection.SubscriptionID())

		return conditions.NewReadyConditionImpactingError(
			err, conditions.ConditionSeverityError, conditions.ReasonSubscriptionMismatch)
	}

	return nil
}

func (r *azureDeploymentReconcilerInstance) handleCreateOrUpdateFailed(err error) error {
	r.Log.V(Debug).Info(
		"Resource creation/update failure",
		"resourceID", genruntime.GetResourceIDOrDefault(r.Obj),
		"error", err.Error())

	err = r.MakeReadyConditionImpactingErrorFromError(err)
	ClearPollerResumeToken(r.Obj)

	return err
}

func (r *azureDeploymentReconcilerInstance) handleDeleteFailed(err error) error {
	r.Log.V(Debug).Info(
		"Resource deletion failure",
		"resourceID", genruntime.GetResourceIDOrDefault(r.Obj),
		"error", err.Error())

	err = r.MakeReadyConditionImpactingErrorFromError(err)
	// Force all delete errors to have severity Warning, as we don't want to block the deletion of the resource
	// and there's no good way for users to restart a stopped deletion, as the deletionTimestamp can only be set once.
	// Without this, deletes that hit an intermittent 400 will get stuck forever
	if readyConditionImpactingErr, ok := conditions.AsReadyConditionImpactingError(err); ok {
		readyConditionImpactingErr.Severity = conditions.ConditionSeverityWarning
	}
	ClearPollerResumeToken(r.Obj)

	return err
}

type CreateOrUpdateSuccessMode string

const (
	WatchResource  = CreateOrUpdateSuccessMode("watch")
	ManageResource = CreateOrUpdateSuccessMode("manage")
)

func (r *azureDeploymentReconcilerInstance) handleCreateOrUpdateSuccess(ctx context.Context, mode CreateOrUpdateSuccessMode) error {
	r.Log.V(Status).Info(
		"Resource successfully created/updated",
		"resourceID", genruntime.GetResourceIDOrDefault(r.Obj))

	err := r.updateStatus(ctx, r.Obj)
	if err != nil {
		if mode == WatchResource {
			if genericarmclient.IsNotFoundError(err) {
				err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonAzureResourceNotFound)
			}
			return err
		} else {
			if genericarmclient.IsNotFoundError(err) {
				// If we're getting NotFound here there must be an RP bug, as poller said success. If that happens we want
				// to make sure that we don't get stuck, so we clear the poller URL.
				ClearPollerResumeToken(r.Obj)
			}
			return eris.Wrapf(err, "error updating status")
		}
	}

	check, err := r.postReconciliationCheck(ctx)
	if err != nil {
		impactingError, ok := conditions.AsReadyConditionImpactingError(err)
		if !ok {
			impactingError = conditions.NewReadyConditionImpactingError(
				err,
				conditions.ConditionSeverityWarning,
				conditions.ReasonFailed)
		}

		return impactingError
	}

	// If post reconcile check is failed, we return ReadyConditionImpactingError here to update the Ready
	// condition is updated so the user can see why we're not setting ready condition now.
	if check.ReconciliationFailed() {
		r.Log.V(Status).Info("Extension post-reconcile check failure", "message", check.Message())
		return check.CreateConditionError()
	}

	err = r.saveAssociatedKubernetesResources(ctx)
	if err != nil {
		if _, ok := core.AsNotOwnedError(err); ok {
			err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonAdditionalKubernetesObjWriteFailure)
		}

		return err
	}

	onSuccess := extensions.CreateSuccessfulCreationHandler(r.Extension, r.Log)
	err = onSuccess(r.Obj)
	if err != nil {
		return err
	}

	ClearPollerResumeToken(r.Obj)
	return nil
}

func (r *azureDeploymentReconcilerInstance) postReconciliationCheck(ctx context.Context) (extensions.PostReconcileCheckResult, error) {
	// Create a checker for access to the extension point, if required
	checker, extensionFound := extensions.CreatePostReconciliationChecker(r.Extension)
	if !extensionFound {
		// No extension found, nothing to do
		return extensions.PostReconcileCheckResultSuccess(), nil
	}

	// We also need to have our owner, it too with an up-to-date status
	ownerDetails, ownerErr := r.ResourceResolver.ResolveOwner(ctx, r.Obj)
	if ownerErr != nil {
		// We can't obtain the owner, so we can't run the extension
		return extensions.PostReconcileCheckResult{}, ownerErr
	}

	// Run our post-reconciliation checker
	check, checkErr := checker(ctx, r.Obj, ownerDetails.Owner, r.ResourceResolver, r.ARMConnection.Client(), r.Log)
	if checkErr != nil {
		// Something went wrong running the check.
		return extensions.PostReconcileCheckResult{}, checkErr
	}

	return check, nil
}

func (r *azureDeploymentReconcilerInstance) MonitorResourceCreation(ctx context.Context) (ctrl.Result, error) {
	pollerID, pollerResumeToken, hasToken := GetPollerResumeToken(r.Obj)
	if !hasToken {
		return ctrl.Result{}, eris.New("cannot MonitorResourceCreation with empty pollerResumeToken or pollerID")
	}

	if pollerID != genericarmclient.CreatePollerID {
		return ctrl.Result{}, eris.Errorf("cannot MonitorResourceCreation with pollerID=%s", pollerID)
	}

	poller := r.ARMConnection.Client().ResumeCreatePoller(pollerID)
	err := poller.Resume(ctx, r.ARMConnection.Client(), pollerResumeToken)
	if err != nil {
		return r.resultBasedOnGenerationCount(), r.handleCreateOrUpdateFailed(err)
	}

	if poller.Poller.Done() {
		return r.resultBasedOnGenerationCount(), r.handleCreateOrUpdateSuccess(ctx, ManageResource)
	}

	// Requeue to check again later
	retryAfter := genericarmclient.GetRetryAfter(poller.RawResponse)
	r.Log.V(Debug).Info("Resource not created/updated yet, will check again", "requeueAfter", retryAfter)
	return ctrl.Result{Requeue: true, RequeueAfter: retryAfter}, nil
}

func (r *azureDeploymentReconcilerInstance) resultBasedOnGenerationCount() ctrl.Result {
	// Once poller is done or run into error, we need to check if there was another event while resource had a ResumePollerToken.
	// We do it here by checking the latest-reconciled-generation to make sure that we have sent the latest changes to the RP.
	// If there's a mismatch in number of generations we reconciled and generations on spec, we requeue the resource to make sure its in sync.
	generation, hasGenerationAnnotation := GetLatestReconciledGeneration(r.Obj)
	if hasGenerationAnnotation && r.Obj.GetGeneration() != generation {
		r.Log.V(Debug).Info(
			"Generation mismatch detected, requeue-ing the resource",
			"resourceID",
			genruntime.GetResourceIDOrDefault(r.Obj))

		return ctrl.Result{Requeue: true}
	}
	return ctrl.Result{}
}

//////////////////////////////////////////
// Other helpers
//////////////////////////////////////////

var zeroDuration time.Duration = 0

func (r *azureDeploymentReconcilerInstance) getStatus(
	ctx context.Context,
	obj genruntime.ARMMetaObject,
	id string,
) (genruntime.ConvertibleStatus, time.Duration, error) { //nolint:unparam
	armStatus, err := genruntime.NewEmptyARMStatus(obj, r.ResourceResolver.Scheme())
	if err != nil {
		return nil, zeroDuration, eris.Wrapf(err, "constructing ARM status for resource: %q", id)
	}

	apiVersion, verr := genruntime.GetAPIVersion(obj, r.ResourceResolver.Scheme())
	if verr != nil {
		return nil, zeroDuration, eris.Wrapf(verr, "error getting api version for resource %s while getting status", obj.GetName())
	}

	// Get the resource
	if genruntime.ResourceOperationGet.IsSupportedBy(obj) {
		var retryAfter time.Duration
		retryAfter, err = r.ARMConnection.Client().GetByID(ctx, id, apiVersion, armStatus)
		if err != nil {
			return nil, retryAfter, eris.Wrapf(err, "getting resource with ID: %q", id)
		}

		if r.Log.V(Debug).Enabled() {
			statusBytes, marshalErr := json.Marshal(armStatus)
			if marshalErr != nil {
				return nil, zeroDuration, eris.Wrapf(marshalErr, "serializing ARM status to JSON for debugging")
			}

			r.Log.V(Debug).Info("Got ARM status", "status", string(statusBytes))
		}
	} else if genruntime.ResourceOperationHead.IsSupportedBy(obj) {
		var retryAfter time.Duration
		var exists bool
		exists, retryAfter, err = r.ARMConnection.Client().CheckExistenceByID(ctx, id, apiVersion)
		if err != nil {
			return nil, retryAfter, eris.Wrapf(err, "getting resource with ID: %q", id)
		}

		// We expect the resource to exist
		if !exists {
			return nil, retryAfter, eris.Wrapf(err, "getting resource with ID: %q", id)
		}
	} else {
		return nil, zeroDuration, eris.Errorf("resource must support one of GET or HEAD, but it supports neither")
	}

	// Convert the ARM shape to the Kube shape
	status, err := genruntime.NewEmptyVersionedStatus(obj, r.ResourceResolver.Scheme())
	if err != nil {
		return nil, zeroDuration, eris.Wrapf(err, "constructing Kube status object for resource: %q", id)
	}

	// Create an owner reference
	owner := obj.Owner()
	var knownOwner genruntime.ArbitraryOwnerReference
	if owner != nil {
		knownOwner = genruntime.ArbitraryOwnerReference{
			Name:  owner.Name,
			Group: owner.Group,
			Kind:  owner.Kind,
		}
	}

	// Fill the kube status with the results from the arm status
	// TODO: The owner parameter here should be optional
	if s, ok := status.(genruntime.FromARMConverter); ok {
		err = s.PopulateFromARM(knownOwner, reflecthelpers.ValueOfPtr(armStatus)) // TODO: PopulateFromArm expects a value... ick
		if err != nil {
			return nil, zeroDuration, eris.Wrapf(err, "converting ARM status to Kubernetes status")
		}
	} else {
		return nil, zeroDuration, eris.Errorf("expected status %T to implement genruntime.FromARMConverter", s)
	}

	return status, zeroDuration, nil
}

func (r *azureDeploymentReconcilerInstance) setStatus(obj genruntime.ARMMetaObject, status genruntime.ConvertibleStatus) error {
	// Modifications that impact status have to happen after this because this performs a full
	// replace of status
	if status != nil {
		// SetStatus() takes care of any required conversion to the right version
		err := obj.SetStatus(status)
		if err != nil {
			return eris.Wrapf(err, "setting status on %s", obj.GetObjectKind().GroupVersionKind())
		}
	}

	return nil
}

func (r *azureDeploymentReconcilerInstance) updateStatus(ctx context.Context, obj genruntime.ARMMetaObject) error {
	resourceID, hasResourceID := genruntime.GetResourceID(obj)
	if !hasResourceID {
		return eris.Errorf("resource has no resource id")
	}

	status, _, err := r.getStatus(ctx, obj, resourceID)
	if err != nil {
		return eris.Wrapf(err, "error getting status for resource ID %q", resourceID)
	}

	if err = r.setStatus(obj, status); err != nil {
		return err
	}

	return nil
}

// getStatusFromARMID creates a temporary ARMMetaObject for the given GVK, sets the ARM ID on it,
// and fetches its status from ARM. This is used when the owner is specified as an ARM ID rather than
// a Kubernetes resource reference to populate its status details.
func (r *azureDeploymentReconcilerInstance) getStatusFromARMID(
	ctx context.Context,
	armID string,
	gvk schema.GroupVersionKind,
) (genruntime.ARMMetaObject, error) {
	// Create a new empty object of the appropriate type
	obj, err := r.ResourceResolver.Scheme().New(gvk)
	if err != nil {
		return nil, eris.Wrapf(err, "creating new object for GVK %s", gvk)
	}

	// Ensure GVK is set on the object
	obj.GetObjectKind().SetGroupVersionKind(gvk)

	// Cast to ARMMetaObject
	metaObj, ok := obj.(genruntime.ARMMetaObject)
	if !ok {
		return nil, eris.Errorf("object of type %T does not implement genruntime.ARMMetaObject", obj)
	}

	// Set the resource ID annotation so updateStatus can find it
	genruntime.SetResourceID(metaObj, armID)
	// Set the spec.OriginalVersion to the latest version
	err = reflecthelpers.SetProperty(metaObj.GetSpec(), "OriginalVersion", strings.TrimSuffix(gvk.Version, "storage")) // This is real hacky
	if err != nil {
		return nil, eris.Wrapf(err, "setting Spec.OriginalVersion for ARM ID %s", armID)
	}

	// Fetch and populate status from ARM
	err = r.updateStatus(ctx, metaObj)
	if err != nil {
		return nil, eris.Wrapf(err, "updating status for ARM ID %s", armID)
	}

	return metaObj, nil
}

// saveAssociatedKubernetesResources retrieves Kubernetes resources to create and saves them to Kubernetes.
// If there are no resources to save this method is a no-op.
func (r *azureDeploymentReconcilerInstance) saveAssociatedKubernetesResources(ctx context.Context) error {
	originalVersion, err := genruntime.ObjAsOriginalVersion(r.Obj, r.ResourceResolver.Scheme())
	if err != nil {
		return err
	}

	var resources []client.Object

	// Special case, because we need to find what secrets the secretExpressionExporter needs and get those secrets
	additionalSecrets, err := findRequiredSecrets(r.ExpressionEvaluator, r.Obj, originalVersion)
	if err != nil {
		return eris.Wrapf(err, "error finding required secrets")
	}
	secretExporter := &kubernetesSecretExporter{
		obj:               r.Obj,
		connection:        r.ARMConnection,
		log:               r.Log,
		extension:         r.Extension,
		additionalSecrets: additionalSecrets,
	}

	var additionalResources []client.Object
	additionalResources, err = secretExporter.Export(ctx)
	if err != nil {
		return err
	}
	resources = append(resources, additionalResources...)

	exporters := []kubernetesResourceExporter{
		&autoGeneratedConfigExporter{
			versionedObj: originalVersion,
			log:          r.Log,
			connection:   r.ARMConnection,
		},
		&manualConfigExporter{
			obj:        r.Obj,
			extension:  r.Extension,
			log:        r.Log,
			connection: r.ARMConnection,
		},
		&configMapExpressionExporter{
			obj:                 r.Obj,
			versionedObj:        originalVersion,
			expressionEvaluator: r.ExpressionEvaluator,
		},
		&secretExpressionExporter{
			obj:                 r.Obj,
			versionedObj:        originalVersion,
			expressionEvaluator: r.ExpressionEvaluator,
			rawSecrets:          secretExporter.rawSecrets,
		},
	}

	for _, exporter := range exporters {
		additionalResources, err = exporter.Export(ctx)
		if err != nil {
			return err
		}
		resources = append(resources, additionalResources...)
	}

	// We do a bit of duplicate work here because each handler also does merging of its own, but then we
	// have to merge the merges. Technically we could allow each handler to just export a list of secrets (with
	// duplicate entries) and then merge once.
	merged, err := merger.MergeObjects(resources)
	if err != nil {
		return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonAdditionalKubernetesObjWriteFailure)
	}

	results, err := genruntime.ApplyObjsAndEnsureOwner(ctx, r.KubeClient, r.Obj, merged)
	if err != nil {
		return err
	}

	if len(results) != len(merged) {
		return eris.Errorf("unexpected results len %d not equal to Kuberentes resources length %d", len(results), len(resources))
	}

	for i := 0; i < len(merged); i++ {
		resource := merged[i]
		result := results[i]

		r.Log.V(Debug).Info("Successfully created resource",
			"namespace", resource.GetNamespace(),
			"name", resource.GetName(),
			"type", fmt.Sprintf("%T", resource),
			"action", result)
	}

	return nil
}

// ConvertResourceToARMResource converts a genruntime.ARMMetaObject (a Kubernetes representation of a resource) into
// a genruntime.ARMResourceSpec - a specification which can be submitted to Azure for deployment
func (r *azureDeploymentReconcilerInstance) ConvertResourceToARMResource(ctx context.Context) (genruntime.ARMResource, error) {
	metaObject := r.Obj

	result, err := ConvertToARMResourceImpl(ctx, metaObject, r.ResourceResolver, r.ARMConnection.SubscriptionID())
	if err != nil {
		return nil, err
	}

	// Run any resource-specific extensions
	modifier := extensions.CreateARMResourceModifier(r.Extension, r.ARMConnection.Client(), r.KubeClient, r.ResourceResolver, r.Log)
	return modifier(ctx, metaObject, result)
}

// ConvertToARMResourceImpl factored out of AzureDeploymentReconciler.ConvertResourceToARMResource to allow for testing
func ConvertToARMResourceImpl(
	ctx context.Context,
	metaObject genruntime.ARMMetaObject,
	resolver *resolver.Resolver,
	subscriptionID string,
) (genruntime.ARMResource, error) {
	// This calls ObjAsOriginalVersion which technically is doing more work than strictly needed, as it converts the spec,
	// metadata, and status. We need the spec and metadata but don't strictly need the status converted at this point.
	// We could look to do some future optimizations here to avoid conversion of status if it ever becomes an issue.
	versionedMeta, err := genruntime.ObjAsOriginalVersion(metaObject, resolver.Scheme())
	if err != nil {
		return nil, err
	}

	spec := versionedMeta.GetSpec()
	armTransformer, ok := spec.(genruntime.ARMTransformer)
	if !ok {
		return nil, eris.Errorf("spec was of type %T which doesn't implement genruntime.ArmTransformer", spec)
	}

	resourceHierarchy, resolvedDetails, err := resolver.ResolveAll(ctx, versionedMeta)
	if err != nil {
		return nil, reconcilers.ClassifyResolverError(err)
	}

	armSpec, err := armTransformer.ConvertToARM(resolvedDetails)
	if err != nil {
		return nil, eris.Wrapf(err, "transforming resource %s to ARM", metaObject.GetName())
	}

	typedArmSpec, ok := armSpec.(genruntime.ARMResourceSpec)
	if !ok {
		return nil, eris.Errorf("casting armSpec of type %T to genruntime.ARMResourceSpec", armSpec)
	}

	armID, err := resourceHierarchy.FullyQualifiedARMID(subscriptionID)
	if err != nil {
		return nil, reconcilers.ClassifyResolverError(err)
	}

	result := genruntime.NewARMResource(typedArmSpec, nil, armID)
	return result, nil
}

// skipDeletionPrecheck is a set of resource groups for which we skip the pre-deletion existence check.
// This is to bypass the need to re-record every test in one go - we enable the extra check group by group.
var skipDeletionPrecheck = sets.NewString(
	"alertsmanagement.azure.com",
	"apimanagement.azure.com",
	"app.azure.com",
	"appconfiguration.azure.com",
	"cache.azure.com",
	"cdn.azure.com",
	"cognitiveservices.azure.com",
	"compute.azure.com",
	"containerinstance.azure.com",
	"containerregistry.azure.com",
	"containerservice.azure.com",
	"datafactory.azure.com",
	"dataprotection.azure.com",
	"dbformariadb.azure.com",
	"dbforpostgresql.azure.com",
	"devices.azure.com",
	"documentdb.azure.com",
	"eventgrid.azure.com",
	"eventhub.azure.com",
	"insights.azure.com",
	"keyvault.azure.com",
	"kubernetesconfiguration.azure.com",
	"kusto.azure.com",
	"machinelearningservices.azure.com",
	"managedidentity.azure.com",
	"monitor.azure.com",
	"network.azure.com",
	"network.frontdoor.azure.com",
	"notificationhubs.azure.com",
	"operationalinsights.azure.com",
	"quota.azure.com",
	"redhatopenshift.azure.com",
	"resources.azure.com",
	"search.azure.com",
	"servicebus.azure.com",
	"signalrservice.azure.com",
	"sql.azure.com",
	"storage.azure.com",
	"subscription.azure.com",
	"synapse.azure.com",
	"web.azure.com",
)

// deleteResource deletes a resource in ARM. This function is used as the default deletion handler and can
// have its behavior modified by resources implementing the genruntime.Deleter extension
func (r *azureDeploymentReconcilerInstance) deleteResource(
	ctx context.Context,
	log logr.Logger,
	resolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	obj genruntime.ARMMetaObject,
) (ctrl.Result, error) {
	// If we have no resourceID to begin with, the Azure resource was never created
	resourceID := genruntime.GetResourceIDOrDefault(obj)
	if resourceID == "" {
		log.V(Status).Info("Not issuing ARM delete as resource had no ResourceID annotation")
		return ctrl.Result{}, nil
	}

	err := r.checkSubscription(resourceID)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check to see if the resource has already been deleted from Azure - if so, we're done.
	// But, first check to see if this resource is in a deny group, and skip the check if so.
	// This is to allow us to fix up remaining issues one by one instead of all at once.
	group := obj.GetObjectKind().GroupVersionKind().Group
	if !skipDeletionPrecheck.Has(group) {
		if _, _, err := r.getStatus(ctx, obj, resourceID); err != nil {
			if genericarmclient.IsNotFoundError(err) {
				// Resource no longer exists
				log.V(Info).Info("Resource is already gone, skipping issue of DELETE to Azure")
				return ctrl.Result{}, nil
			}
		}
	}

	// Optimizations or complications of this delete path should be undertaken with care.
	// Be especially cautious of relying on the controller-runtime SharedInformer cache
	// as a source of truth about if this resource or its parents have already been deleted, as
	// the SharedInformer cache is not read-through and will return NotFound if it just hasn't been
	// populated yet.
	// Generally speaking the safest thing we can do is just issue the DELETE to Azure.

	// retryAfter = ARM can tell us how long to wait for a DELETE
	originalAPIVersion, err := genruntime.GetAPIVersion(obj, resolver.Scheme())
	if err != nil {
		return ctrl.Result{}, err
	}
	pollerResp, err := armClient.BeginDeleteByID(ctx, resourceID, originalAPIVersion)
	if err != nil {
		if genericarmclient.IsNotFoundError(err) {
			log.V(Info).Info("Successfully issued DELETE to Azure - resource was already gone")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, r.handleDeleteFailed(err)
	}
	log.V(Info).Info("Successfully issued DELETE to Azure")

	// If we are done here it means delete succeeded immediately. It can't have failed because if it did
	// we would have taken the error path, above.
	if pollerResp.Poller.Done() {
		return ctrl.Result{}, nil
	}

	retryAfter := genericarmclient.GetRetryAfter(pollerResp.RawResponse)
	resumeToken, err := pollerResp.Poller.ResumeToken()
	if err != nil {
		return ctrl.Result{}, eris.Wrapf(err, "couldn't create DELETE resume token for resource %q", resourceID)
	}
	SetPollerResumeToken(obj, pollerResp.ID, resumeToken)

	// Normally don't need to set both of these fields but because retryAfter can be 0 we do
	return ctrl.Result{Requeue: true, RequeueAfter: retryAfter}, nil
}
