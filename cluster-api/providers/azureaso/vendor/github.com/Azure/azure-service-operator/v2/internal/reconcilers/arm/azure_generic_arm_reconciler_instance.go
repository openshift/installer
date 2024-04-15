/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
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
		r.Log.Error(err, "Error during Delete", "action", action)
		r.Recorder.Event(r.Obj, v1.EventTypeWarning, "DeleteActionError", err.Error())

		return ctrl.Result{}, err
	}

	return result, nil
}

func (r *azureDeploymentReconcilerInstance) MakeReadyConditionImpactingErrorFromError(azureErr error) error {
	var readyConditionError *conditions.ReadyConditionImpactingError
	isReadyConditionImpactingError := errors.As(azureErr, &readyConditionError)
	if isReadyConditionImpactingError {
		// The error has already been classified. This currently only happens in test with the go-vcr injected
		// http client
		return azureErr
	}

	var cloudError *genericarmclient.CloudError
	isCloudErr := errors.As(azureErr, &cloudError)
	if !isCloudErr {
		// This shouldn't happen, as all errors from ARM should be in one of the shapes that CloudError supports. In case
		// we've somehow gotten one that isn't formatted correctly, create a sensible default error
		return conditions.NewReadyConditionImpactingError(
			azureErr,
			conditions.ConditionSeverityWarning,
			conditions.MakeReason(core.UnknownErrorCode))
	}

	apiVersion, verr := r.GetAPIVersion()
	if verr != nil {
		return errors.Wrapf(verr, "error getting api version for resource %s while making Ready condition", r.Obj.GetName())
	}

	classifier := extensions.CreateErrorClassifier(r.Extension, ClassifyCloudError, apiVersion, r.Log)
	details, err := classifier(cloudError)
	if err != nil {
		return errors.Wrapf(
			err,
			"Unable to classify cloud error (%s)",
			cloudError.Error())
	}

	var severity conditions.ConditionSeverity
	switch details.Classification {
	case core.ErrorRetryable:
		severity = conditions.ConditionSeverityWarning
	case core.ErrorFatal:
		severity = conditions.ConditionSeverityError
		// This case purposefully does nothing as the fatal provisioning state was already set above
	default:
		return errors.Errorf(
			"unknown error classification %q while making Ready condition",
			details.Classification)
	}

	// Stick errorDetails.Message into an error so that it will be displayed as the message on the condition
	err = errors.Wrapf(cloudError, details.Message)
	reason := conditions.MakeReason(details.Code)
	result := conditions.NewReadyConditionImpactingError(err, severity, reason)

	return result
}

func (r *azureDeploymentReconcilerInstance) AddInitialResourceState(ctx context.Context) error {
	armResource, err := r.ConvertResourceToARMResource(ctx)
	if err != nil {
		return err
	}
	genruntime.SetResourceID(r.Obj, armResource.GetID())
	labels.SetOwnerNameLabel(r.Obj)
	labels.SetOwnerGroupKindLabel(r.Obj)
	return nil
}

func (r *azureDeploymentReconcilerInstance) DetermineDeleteAction() (DeleteAction, DeleteActionFunc, error) {
	pollerID, _, hasPollerResumeToken := GetPollerResumeToken(r.Obj)

	if hasPollerResumeToken && pollerID == genericarmclient.DeletePollerID {
		return DeleteActionMonitorDelete, r.MonitorDelete, nil
	}

	return DeleteActionBeginDelete, r.StartDeleteOfResource, nil
}

func (r *azureDeploymentReconcilerInstance) DetermineCreateOrUpdateAction() (CreateOrUpdateAction, CreateOrUpdateActionFunc, error) {
	ready := genruntime.GetReadyCondition(r.Obj)
	_, _, hasPollerResumeToken := GetPollerResumeToken(r.Obj)

	if ready != nil && ready.Reason == conditions.ReasonDeleting.Name {
		return CreateOrUpdateActionNoAction, NoAction, errors.Errorf("resource is currently deleting; it can not be applied")
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
		return ctrl.Result{}, errors.New("cannot MonitorResourceCreation with empty pollerResumeToken or pollerID")
	}

	if pollerID != genericarmclient.DeletePollerID {
		return ctrl.Result{}, errors.Errorf("cannot MonitorResourceCreation with pollerID=%s", pollerID)
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

func (r *azureDeploymentReconcilerInstance) BeginCreateOrUpdateResource(
	ctx context.Context,
) (ctrl.Result, error) {
	if r.Obj.AzureName() == "" {
		err := errors.New("AzureName was not set. A webhook should default this to .metadata.name if it was omitted. Is the ASO webhook service running?")
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
			errors.Wrapf(err, "couldn't create PUT resume token for resource %q", armResource.GetID())
	}

	SetPollerResumeToken(r.Obj, pollerResp.ID, resumeToken)
	return ctrl.Result{Requeue: true}, nil
}

func (r *azureDeploymentReconcilerInstance) preReconciliationCheck(ctx context.Context) (extensions.PreReconcileCheckResult, error) {
	// Create a checker for access to the extension point, if required
	checker, extensionFound := extensions.CreatePreReconciliationChecker(r.Extension)
	if !extensionFound {
		// No extension found, nothing to do
		return extensions.ProceedWithReconcile(), nil
	}

	// Having a checker requires our resource to have an up-to-date status
	r.Log.V(Verbose).Info("Refreshing Status of resource")
	statusErr := r.updateStatus(ctx)
	if statusErr != nil && !genericarmclient.IsNotFoundError(statusErr) {
		// We have an error, and it's not because the resource doesn't exist yet
		return extensions.PreReconcileCheckResult{}, statusErr
	}

	// We also need to have our owner, it too with an up-to-date status
	ownerDetails, ownerErr := r.ResourceResolver.ResolveOwner(ctx, r.Obj)
	if ownerErr != nil {
		// We can't obtain the owner, so we can't run the extension
		return extensions.PreReconcileCheckResult{}, ownerErr
	}

	// Run our pre-reconciliation checker
	check, checkErr := checker(ctx, r.Obj, ownerDetails.Owner, r.ResourceResolver, r.ARMConnection.Client(), r.Log)
	if checkErr != nil {
		// Something went wrong running the check.
		return extensions.PreReconcileCheckResult{}, checkErr
	}

	return check, nil
}

// checkSubscription checks if subscription on resource matches with credentials used while creating a resource.
// Which prevents users to modify subscription in their credential.
func (r *azureDeploymentReconcilerInstance) checkSubscription(resourceID string) error {
	parsedRID, err := arm.ParseResourceID(resourceID)
	// Some resources like '/providers/Microsoft.Subscription/aliases' do not have subscriptionID, so we need to make sure subscriptionID exists before we check.
	// TODO: we need a better way?
	if err == nil {
		if parsedRID.ResourceGroupName != "" && parsedRID.SubscriptionID != r.ARMConnection.SubscriptionID() {
			err = errors.Errorf("SubscriptionID %q for %q resource does not match with Client Credential: %q", parsedRID.SubscriptionID, resourceID, r.ARMConnection.SubscriptionID())
			return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonSubscriptionMismatch)
		}
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

	err := r.updateStatus(ctx)
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
			return errors.Wrapf(err, "error updating status")
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
		return ctrl.Result{}, errors.New("cannot MonitorResourceCreation with empty pollerResumeToken or pollerID")
	}

	if pollerID != genericarmclient.CreatePollerID {
		return ctrl.Result{}, errors.Errorf("cannot MonitorResourceCreation with pollerID=%s", pollerID)
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
	r.Log.V(Debug).Info("Resource not created yet, will check again", "requeueAfter", retryAfter)
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

func (r *azureDeploymentReconcilerInstance) getStatus(ctx context.Context, id string) (genruntime.ConvertibleStatus, time.Duration, error) { // nolint:unparam
	armStatus, err := genruntime.NewEmptyARMStatus(r.Obj, r.ResourceResolver.Scheme())
	if err != nil {
		return nil, zeroDuration, errors.Wrapf(err, "constructing ARM status for resource: %q", id)
	}

	apiVersion, verr := r.GetAPIVersion()
	if verr != nil {
		return nil, zeroDuration, errors.Wrapf(verr, "error getting api version for resource %s while getting status", r.Obj.GetName())
	}

	// Get the resource
	if genruntime.ResourceOperationGet.IsSupportedBy(r.Obj) {
		var retryAfter time.Duration
		retryAfter, err = r.ARMConnection.Client().GetByID(ctx, id, apiVersion, armStatus)
		if err != nil {
			return nil, retryAfter, errors.Wrapf(err, "getting resource with ID: %q", id)
		}

		if r.Log.V(Debug).Enabled() {
			statusBytes, marshalErr := json.Marshal(armStatus)
			if marshalErr != nil {
				return nil, zeroDuration, errors.Wrapf(marshalErr, "serializing ARM status to JSON for debugging")
			}

			r.Log.V(Debug).Info("Got ARM status", "status", string(statusBytes))
		}
	} else if genruntime.ResourceOperationHead.IsSupportedBy(r.Obj) {
		var retryAfter time.Duration
		var exists bool
		exists, retryAfter, err = r.ARMConnection.Client().CheckExistenceByID(ctx, id, apiVersion)
		if err != nil {
			return nil, retryAfter, errors.Wrapf(err, "getting resource with ID: %q", id)
		}

		// We expect the resource to exist
		if !exists {
			return nil, retryAfter, errors.Wrapf(err, "getting resource with ID: %q", id)
		}
	} else {
		return nil, zeroDuration, errors.Errorf("resource must support one of GET or HEAD, but it supports neither")
	}

	// Convert the ARM shape to the Kube shape
	status, err := genruntime.NewEmptyVersionedStatus(r.Obj, r.ResourceResolver.Scheme())
	if err != nil {
		return nil, zeroDuration, errors.Wrapf(err, "constructing Kube status object for resource: %q", id)
	}

	// Create an owner reference
	owner := r.Obj.Owner()
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
			return nil, zeroDuration, errors.Wrapf(err, "converting ARM status to Kubernetes status")
		}
	} else {
		return nil, zeroDuration, errors.Errorf("expected status %T to implement genruntime.FromARMConverter", s)
	}

	return status, zeroDuration, nil
}

func (r *azureDeploymentReconcilerInstance) setStatus(status genruntime.ConvertibleStatus) error {
	// Modifications that impact status have to happen after this because this performs a full
	// replace of status
	if status != nil {
		// SetStatus() takes care of any required conversion to the right version
		err := r.Obj.SetStatus(status)
		if err != nil {
			return errors.Wrapf(err, "setting status on %s", r.Obj.GetObjectKind().GroupVersionKind())
		}
	}

	return nil
}

func (r *azureDeploymentReconcilerInstance) updateStatus(ctx context.Context) error {
	resourceID, hasResourceID := genruntime.GetResourceID(r.Obj)
	if !hasResourceID {
		return errors.Errorf("resource has no resource id")
	}

	status, _, err := r.getStatus(ctx, resourceID)
	if err != nil {
		return errors.Wrapf(err, "error getting status for resource ID %q", resourceID)
	}

	if err = r.setStatus(status); err != nil {
		return err
	}

	return nil
}

// saveAssociatedKubernetesResources retrieves Kubernetes resources to create and saves them to Kubernetes.
// If there are no resources to save this method is a no-op.
func (r *azureDeploymentReconcilerInstance) saveAssociatedKubernetesResources(ctx context.Context) error {
	// Check if this resource has a handcrafted extension for exporting
	retriever := extensions.CreateKubernetesExporter(ctx, r.Extension, r.ARMConnection.Client(), r.Log)
	resources, err := retriever(r.Obj)
	if err != nil {
		return errors.Wrap(err, "extension failed to produce resources for export")
	}

	// Also check if the resource itself implements KubernetesExporter
	exporter, ok := r.ObjAsKubernetesExporter()
	if ok {
		var additionalResources []client.Object
		additionalResources, err = exporter.ExportKubernetesResources(ctx, r.Obj, r.ARMConnection.Client(), r.Log)
		if err != nil {
			return errors.Wrap(err, "failed to produce resources for export")
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
		return errors.Errorf("unexpected results len %d not equal to Kuberentes resources length %d", len(results), len(resources))
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

// ObjAsKubernetesExporter returns r.Obj as a genruntime.KubernetesExporter if it supports that interface, respecting
// the original API version used to create the resource.
func (r *azureDeploymentReconcilerInstance) ObjAsKubernetesExporter() (genruntime.KubernetesExporter, bool) {
	resource := r.Obj
	resourceGVK := resource.GetObjectKind().GroupVersionKind()

	desiredGVK := genruntime.GetOriginalGVK(resource)
	if resourceGVK != desiredGVK {
		// Need to convert the hub resource we have to the original version used
		scheme := r.ResourceResolver.Scheme()
		versionedResource, err := genruntime.NewEmptyVersionedResourceFromGVK(scheme, desiredGVK)
		if err != nil {
			r.Log.V(Status).Info(
				"Unable to create expected resource version",
				"have", resourceGVK,
				"desired", desiredGVK)
			return nil, false
		}

		if convertible, ok := versionedResource.(conversion.Convertible); ok {
			hub := resource.(conversion.Hub)
			err := convertible.ConvertFrom(hub)
			if err != nil {
				r.Log.V(Status).Info(
					"Unable to convert resource to expected version",
					"original", resourceGVK,
					"destination", desiredGVK)
				return nil, false
			}
		}

		resource = versionedResource
	}

	// Now test whether we support the interface
	result, ok := resource.(genruntime.KubernetesExporter)
	return result, ok
}

// ConvertResourceToARMResource converts a genruntime.ARMMetaObject (a Kubernetes representation of a resource) into
// a genruntime.ARMResourceSpec - a specification which can be submitted to Azure for deployment
func (r *azureDeploymentReconcilerInstance) ConvertResourceToARMResource(ctx context.Context) (genruntime.ARMResource, error) {
	metaObject := r.Obj
	scheme := r.ResourceResolver.Scheme()

	result, err := ConvertToARMResourceImpl(ctx, metaObject, scheme, r.ResourceResolver, r.ARMConnection.SubscriptionID())
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
	scheme *runtime.Scheme,
	resolver *resolver.Resolver,
	subscriptionID string) (genruntime.ARMResource, error) {
	spec, err := genruntime.GetVersionedSpec(metaObject, scheme)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get spec from %s", metaObject.GetObjectKind().GroupVersionKind())
	}

	armTransformer, ok := spec.(genruntime.ARMTransformer)
	if !ok {
		return nil, errors.Errorf("spec was of type %T which doesn't implement genruntime.ArmTransformer", spec)
	}

	resourceHierarchy, resolvedDetails, err := resolver.ResolveAll(ctx, metaObject)
	if err != nil {
		return nil, reconcilers.ClassifyResolverError(err)
	}

	armSpec, err := armTransformer.ConvertToARM(resolvedDetails)
	if err != nil {
		return nil, errors.Wrapf(err, "transforming resource %s to ARM", metaObject.GetName())
	}

	typedArmSpec, ok := armSpec.(genruntime.ARMResourceSpec)
	if !ok {
		return nil, errors.Errorf("casting armSpec of type %T to genruntime.ARMResourceSpec", armSpec)
	}

	armID, err := resourceHierarchy.FullyQualifiedARMID(subscriptionID)
	if err != nil {
		return nil, reconcilers.ClassifyResolverError(err)
	}

	result := genruntime.NewARMResource(typedArmSpec, nil, armID)
	return result, nil
}

// GetAPIVersion returns the ARM API version for the resource we're reconciling
func (r *azureDeploymentReconcilerInstance) GetAPIVersion() (string, error) {
	metaObject := r.Obj
	scheme := r.ResourceResolver.Scheme()

	return genruntime.GetAPIVersion(metaObject, scheme)
}

// deleteResource deletes a resource in ARM. This function is used as the default deletion handler and can
// have its behavior modified by resources implementing the genruntime.Deleter extension
func (r *azureDeploymentReconcilerInstance) deleteResource(
	ctx context.Context,
	log logr.Logger,
	resolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	obj genruntime.ARMMetaObject) (ctrl.Result, error) {

	// If we have no resourceID to begin with, the Azure resource was never created
	resourceID := genruntime.GetResourceIDOrDefault(obj)
	if resourceID == "" {
		log.V(Status).Info("Not issuing delete as resource had no ResourceID annotation")
		return ctrl.Result{}, nil
	}

	err := r.checkSubscription(resourceID)
	if err != nil {
		return ctrl.Result{}, err
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
		return ctrl.Result{}, errors.Wrapf(err, "deleting resource %q", resourceID)
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
		return ctrl.Result{}, errors.Wrapf(err, "couldn't create DELETE resume token for resource %q", resourceID)
	}
	SetPollerResumeToken(obj, pollerResp.ID, resumeToken)

	// Normally don't need to set both of these fields but because retryAfter can be 0 we do
	return ctrl.Result{Requeue: true, RequeueAfter: retryAfter}, nil
}
