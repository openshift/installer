/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package mysql

import (
	"context"

	"github.com/go-logr/logr"
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	asomysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	"github.com/Azure/azure-service-operator/v2/internal/config"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

var _ genruntime.Reconciler = &MySQLUserReconciler{}

type MySQLUserReconciler struct {
	reconcilers.ARMOwnedResourceReconcilerCommon
	ResourceResolver   *resolver.Resolver
	CredentialProvider identity.CredentialProvider
	Config             config.Values
}

func NewMySQLUserReconciler(
	kubeClient kubeclient.Client,
	resourceResolver *resolver.Resolver,
	positiveConditions *conditions.PositiveConditionBuilder,
	credentialProvider identity.CredentialProvider,
	cfg config.Values) *MySQLUserReconciler {

	return &MySQLUserReconciler{
		ResourceResolver:   resourceResolver,
		CredentialProvider: credentialProvider,
		Config:             cfg,
		ARMOwnedResourceReconcilerCommon: reconcilers.ARMOwnedResourceReconcilerCommon{
			ResourceResolver: resourceResolver,
			ReconcilerCommon: reconcilers.ReconcilerCommon{
				KubeClient:         kubeClient,
				PositiveConditions: positiveConditions,
			},
		},
	}
}

func (r *MySQLUserReconciler) asUser(obj genruntime.MetaObject) (*asomysql.User, error) {
	typedObj, ok := obj.(*asomysql.User)
	if !ok {
		return nil, errors.Errorf("cannot modify resource that is not of type *asomysql.User. Type is %T", obj)
	}

	return typedObj, nil
}

func (r *MySQLUserReconciler) CreateOrUpdate(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
	user, err := r.asUser(obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Augment Log
	log = log.WithValues("azureName", user.AzureName())
	connector, err := r.newDBConnector(log, user)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = connector.CreateOrUpdate(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MySQLUserReconciler) Delete(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
	user, err := r.asUser(obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Augment Log
	log = log.WithValues("azureName", user.AzureName())

	log.V(Status).Info("Starting delete of resource")

	// Check that this objects owner still exists
	// This is an optimization to avoid excess requests to Azure.
	_, err = r.ResourceResolver.ResolveOwner(ctx, user)
	if err != nil {
		var typedErr *core.ReferenceNotFound
		if errors.As(err, &typedErr) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	connector, err := r.newDBConnector(log, user)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = connector.Delete(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MySQLUserReconciler) Claim(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
	user, err := r.asUser(obj)
	if err != nil {
		return err
	}

	err = r.ARMOwnedResourceReconcilerCommon.ClaimResource(ctx, log, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLUserReconciler) UpdateStatus(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
	user, err := r.asUser(obj)
	if err != nil {
		return err
	}

	connector, err := r.newDBConnector(log, user)
	if err != nil {
		return err
	}

	exists, err := connector.Exists(ctx)
	if err != nil {
		return err
	}

	if !exists {
		err = errors.Errorf("user %s does not exist", user.Spec.AzureName)
		err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonAzureResourceNotFound)
		return err
	}

	return nil
}

func (r *MySQLUserReconciler) newDBConnector(log logr.Logger, user *asomysql.User) (Connector, error) {
	if user.Spec.LocalUser != nil {
		return &localUser{
			user:               user,
			resourceResolver:   r.ResourceResolver,
			credentialProvider: r.CredentialProvider,
			log:                log,
		}, nil
	}

	if user.Spec.AADUser != nil {
		return &aadUser{
			user:               user,
			resourceResolver:   r.ResourceResolver,
			credentialProvider: r.CredentialProvider,
			log:                log,
		}, nil
	}

	// This is also enforced with a webhook
	err := errors.Errorf("unknown user type, user must be LocalUser or AADUser")
	return nil, conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonFailed)
}
