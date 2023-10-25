/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package arm

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/Azure/azure-service-operator/v2/internal/config"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

// TODO: I think we will want to pull some of this back into the Generic Controller so that it happens
// TODO: for all resources

type CreateOrUpdateAction string

const (
	CreateOrUpdateActionNoAction        = CreateOrUpdateAction("NoAction")
	CreateOrUpdateActionBeginCreation   = CreateOrUpdateAction("BeginCreateOrUpdate")
	CreateOrUpdateActionMonitorCreation = CreateOrUpdateAction("MonitorCreateOrUpdate")
)

type DeleteAction string

const (
	DeleteActionBeginDelete   = DeleteAction("BeginDelete")
	DeleteActionMonitorDelete = DeleteAction("MonitorDelete")
)

type (
	CreateOrUpdateActionFunc = func(ctx context.Context) (ctrl.Result, error)
	DeleteActionFunc         = func(ctx context.Context) (ctrl.Result, error)
)

var _ genruntime.Reconciler = &AzureDeploymentReconciler{}

type AzureDeploymentReconciler struct {
	reconcilers.ARMOwnedResourceReconcilerCommon
	ARMConnectionFactory ARMConnectionFactory
	KubeClient           kubeclient.Client
	ResourceResolver     *resolver.Resolver
	PositiveConditions   *conditions.PositiveConditionBuilder
	Config               config.Values
	Extension            genruntime.ResourceExtension
}

func NewAzureDeploymentReconciler(
	armConnectionFactory ARMConnectionFactory,
	kubeClient kubeclient.Client,
	resourceResolver *resolver.Resolver,
	positiveConditions *conditions.PositiveConditionBuilder,
	cfg config.Values,
	extension genruntime.ResourceExtension) *AzureDeploymentReconciler {

	return &AzureDeploymentReconciler{
		ARMConnectionFactory: armConnectionFactory,
		KubeClient:           kubeClient,
		ResourceResolver:     resourceResolver,
		PositiveConditions:   positiveConditions,
		Config:               cfg,
		Extension:            extension,
		ARMOwnedResourceReconcilerCommon: reconcilers.ARMOwnedResourceReconcilerCommon{
			ResourceResolver: resourceResolver,
			ReconcilerCommon: reconcilers.ReconcilerCommon{
				KubeClient:         kubeClient,
				PositiveConditions: positiveConditions,
			},
		},
	}
}

func (r *AzureDeploymentReconciler) asARMObj(obj genruntime.MetaObject) (genruntime.ARMMetaObject, error) {
	typedObj, ok := obj.(genruntime.ARMMetaObject)
	if !ok {
		return nil, errors.Errorf("cannot modify resource that is not of type ARMMetaObject. Type is %T", obj)
	}

	return typedObj, nil
}

func (r *AzureDeploymentReconciler) makeInstance(
	ctx context.Context,
	log logr.Logger,
	eventRecorder record.EventRecorder,
	obj genruntime.MetaObject,
) (*azureDeploymentReconcilerInstance, error) {
	typedObj, err := r.asARMObj(obj)
	if err != nil {
		return nil, err
	}
	// Augment Log with ARM specific stuff
	log = log.WithValues("azureName", typedObj.AzureName())

	clientDetails, err := r.ARMConnectionFactory(ctx, typedObj)
	if err != nil {
		return nil, err
	}

	eventRecorder.Eventf(obj, v1.EventTypeNormal, "CredentialFrom", "Using credential from %q", clientDetails.CredentialFrom().String())

	// TODO: The line between AzureDeploymentReconciler and azureDeploymentReconcilerInstance is still pretty blurry
	return newAzureDeploymentReconcilerInstance(typedObj, log, eventRecorder, clientDetails, *r), nil
}

func (r *AzureDeploymentReconciler) CreateOrUpdate(
	ctx context.Context,
	log logr.Logger,
	eventRecorder record.EventRecorder,
	obj genruntime.MetaObject,
) (ctrl.Result, error) {
	instance, err := r.makeInstance(ctx, log, eventRecorder, obj)
	if err != nil {
		return ctrl.Result{}, err
	}
	return instance.CreateOrUpdate(ctx)
}

func (r *AzureDeploymentReconciler) Delete(
	ctx context.Context,
	log logr.Logger,
	eventRecorder record.EventRecorder,
	obj genruntime.MetaObject,
) (ctrl.Result, error) {
	instance, err := r.makeInstance(ctx, log, eventRecorder, obj)
	if err != nil {
		return ctrl.Result{}, err
	}
	return instance.Delete(ctx)
}

func (r *AzureDeploymentReconciler) Claim(
	ctx context.Context,
	log logr.Logger,
	eventRecorder record.EventRecorder,
	obj genruntime.MetaObject,
) error {
	typedObj, err := r.asARMObj(obj)
	if err != nil {
		return err
	}

	claimer := extensions.CreateClaimer(r.Extension, r.ARMOwnedResourceReconcilerCommon.ClaimResource)
	err = claimer(ctx, log, typedObj)
	if err != nil {
		return err
	}

	instance, err := r.makeInstance(ctx, log, eventRecorder, obj)
	if err != nil {
		return err
	}
	// Add ARM specific details
	err = instance.AddInitialResourceState(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to add initial resource state")
	}

	return nil
}

func (r *AzureDeploymentReconciler) UpdateStatus(
	ctx context.Context,
	log logr.Logger,
	eventRecorder record.EventRecorder,
	obj genruntime.MetaObject,
) error {
	instance, err := r.makeInstance(ctx, log, eventRecorder, obj)
	if err != nil {
		return err
	}

	return instance.handleCreateOrUpdateSuccess(ctx, WatchResource)
}
