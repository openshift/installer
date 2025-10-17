/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package sql

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	asosql "github.com/Azure/azure-service-operator/v2/api/sql/v1"
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

// Here's some useful reading:
// https://azure.microsoft.com/en-us/blog/adding-users-to-your-sql-azure-database/
// https://learn.microsoft.com/en-us/azure/azure-sql/database/logins-create-manage?view=azuresql

// Test logging in as an AAD User:
// https://learn.microsoft.com/en-us/azure/azure-sql/database/authentication-aad-configure?view=azuresql&tabs=azure-powershell#azure-ad-token

// There's LOGIN's and USERs. We only support USERS. LOGINs are server-wide whereas USERs are DB scoped.
// Here are least some built in DB roles:
// # db_owner, db_securityadmin, db_accessadmin, db_backupoperator, db_ddladmin, db_datawriter, db_datareader, db_denydatawriter, db_denydatareader

var _ genruntime.Reconciler = &AzureSQLUserReconciler{}

type AzureSQLUserReconciler struct {
	reconcilers.ARMOwnedResourceReconcilerCommon
	ResourceResolver   *resolver.Resolver
	CredentialProvider identity.CredentialProvider
	Config             config.Values
}

func NewAzureSQLUserReconciler(
	kubeClient kubeclient.Client,
	resourceResolver *resolver.Resolver,
	positiveConditions *conditions.PositiveConditionBuilder,
	credentialProvider identity.CredentialProvider,
	cfg config.Values,
) *AzureSQLUserReconciler {
	return &AzureSQLUserReconciler{
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

func (r *AzureSQLUserReconciler) asUser(obj genruntime.MetaObject) (*asosql.User, error) {
	typedObj, ok := obj.(*asosql.User)
	if !ok {
		return nil, errors.Errorf("cannot modify resource that is not of type *asosql.User. Type is %T", obj)
	}

	return typedObj, nil
}

func (r *AzureSQLUserReconciler) CreateOrUpdate(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
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

func (r *AzureSQLUserReconciler) Delete(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
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

func (r *AzureSQLUserReconciler) Claim(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
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

func (r *AzureSQLUserReconciler) UpdateStatus(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
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

func (r *AzureSQLUserReconciler) newDBConnector(log logr.Logger, user *asosql.User) (Connector, error) {
	if user.Spec.LocalUser != nil {
		return &localUser{
			user:               user,
			resourceResolver:   r.ResourceResolver,
			credentialProvider: r.CredentialProvider,
			log:                log,
		}, nil
	}

	// This is also enforced with a webhook
	err := errors.Errorf("unknown user type, user must be LocalUser")
	return nil, conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonFailed)
}
