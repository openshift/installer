/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package postgresql

import (
	"context"
	"database/sql"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	asopostgresql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1"
	dbforpostgressql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1api20221201/storage"
	"github.com/Azure/azure-service-operator/v2/internal/config"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	postgresqlutil "github.com/Azure/azure-service-operator/v2/internal/util/postgresql"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

var _ genruntime.Reconciler = &PostgreSQLUserReconciler{}

type PostgreSQLUserReconciler struct {
	reconcilers.ARMOwnedResourceReconcilerCommon
	ResourceResolver *resolver.Resolver
	Config           config.Values
}

func NewPostgreSQLUserReconciler(
	kubeClient kubeclient.Client,
	resourceResolver *resolver.Resolver,
	positiveConditions *conditions.PositiveConditionBuilder,
	cfg config.Values,
) *PostgreSQLUserReconciler {
	return &PostgreSQLUserReconciler{
		ResourceResolver: resourceResolver,
		Config:           cfg,
		ARMOwnedResourceReconcilerCommon: reconcilers.ARMOwnedResourceReconcilerCommon{
			ResourceResolver: resourceResolver,
			ReconcilerCommon: reconcilers.ReconcilerCommon{
				KubeClient:         kubeClient,
				PositiveConditions: positiveConditions,
			},
		},
	}
}

func (r *PostgreSQLUserReconciler) asUser(obj genruntime.MetaObject) (*asopostgresql.User, error) {
	typedObj, ok := obj.(*asopostgresql.User)
	if !ok {
		return nil, errors.Errorf("cannot modify resource that is not of type *asopostgresql.User. Type is %T", obj)
	}

	return typedObj, nil
}

func (r *PostgreSQLUserReconciler) CreateOrUpdate(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
	user, err := r.asUser(obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Augment Log
	log = log.WithValues("azureName", user.AzureName())

	// Resolve the secrets
	secrets, err := r.ResourceResolver.ResolveResourceSecretReferences(ctx, user)
	if err != nil {
		return ctrl.Result{}, reconcilers.ClassifyResolverError(err)
	}

	db, err := r.connectToDB(ctx, log, user, secrets)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer db.Close()

	log.V(Status).Info("Creating PostgreSql user")

	password, err := secrets.LookupFromPtr(user.Spec.LocalUser.Password)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to look up .spec.localUser.Password")
	}

	// Create or update the user. Note that this updates password if it has changed
	username := user.Spec.AzureName

	sqlUser, err := postgresqlutil.FindUserIfExist(ctx, db, username)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to find user")
	}
	if sqlUser == nil {
		sqlUser, err = postgresqlutil.CreateUser(ctx, db, username, password)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to create user")
		}
	} else {
		err = postgresqlutil.UpdateUser(ctx, db, *sqlUser, password)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to update user")
		}
	}
	// TODO integrate in create and update user?
	if user.Spec.RoleOptions == nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to look up .spec.roleOptions")
	}
	roleOptions := postgresqlutil.RoleOptions(*user.Spec.RoleOptions)
	// Ensure that the user role options are set
	err = postgresqlutil.ReconcileUserRoleOptions(ctx, db, *sqlUser, roleOptions)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "ensuring user role options")
	}

	// Ensure that the roles are set
	err = postgresqlutil.ReconcileUserServerRoles(ctx, db, *sqlUser, user.Spec.Roles)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "ensuring server roles")
	}

	log.V(Status).Info("Successfully reconciled PostgreSqlUser")

	return ctrl.Result{}, nil
}

func (r *PostgreSQLUserReconciler) Delete(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) (ctrl.Result, error) {
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

	secrets, err := r.ResourceResolver.ResolveResourceSecretReferences(ctx, user)
	if err != nil {
		return ctrl.Result{}, err
	}

	db, err := r.connectToDB(ctx, log, user, secrets)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer db.Close()

	// TODO: There's still probably some ways that this user can be deleted but that we don't detect (and
	// TODO: so might cause an error triggering the resource to get stuck).
	// TODO: We check for owner not existing above, but cases where the server is in the process of being
	// TODO: deleted (or all system tables have been wiped?) might also exist...
	err = postgresqlutil.DropUser(ctx, db, user.Spec.AzureName)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PostgreSQLUserReconciler) Claim(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
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

func (r *PostgreSQLUserReconciler) UpdateStatus(ctx context.Context, log logr.Logger, eventRecorder record.EventRecorder, obj genruntime.MetaObject) error {
	user, err := r.asUser(obj)
	if err != nil {
		return err
	}

	secrets, err := r.ResourceResolver.ResolveResourceSecretReferences(ctx, user)
	if err != nil {
		return err
	}

	db, err := r.connectToDB(ctx, log, user, secrets)
	if err != nil {
		return err
	}
	defer db.Close()

	exists, err := postgresqlutil.DoesUserExist(ctx, db, user.Spec.AzureName)
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

func (r *PostgreSQLUserReconciler) connectToDB(ctx context.Context, _ logr.Logger, user *asopostgresql.User, secrets genruntime.Resolved[genruntime.SecretReference, string]) (*sql.DB, error) {
	// Get the owner - at this point it must exist
	ownerDetails, err := r.ResourceResolver.ResolveOwner(ctx, user)
	if err != nil {
		return nil, errors.Wrapf(err, "resolving owner for user %s", user.Name)
	}

	// Note that this is not actually possible for this type because we don't allow ARMID references for these owners,
	// but protecting against it here anyway.
	if !ownerDetails.FoundKubernetesOwner() {
		return nil, errors.Errorf("user owner must exist in Kubernetes for user %s", user.Name)
	}

	flexibleServer, ok := ownerDetails.Owner.(*dbforpostgressql.FlexibleServer)
	if !ok {
		return nil, errors.Errorf("owner was not type FlexibleServer, instead: %T", ownerDetails)
	}
	// Assertion to ensure that this is still the storage type
	// If this doesn't compile, update the version being imported to the new Hub version
	var _ ctrlconversion.Hub = &dbforpostgressql.FlexibleServer{}

	if flexibleServer.Status.FullyQualifiedDomainName == nil {
		// This possibly means that the server hasn't finished deploying yet
		err = errors.Errorf("owning Flexibleserver %q '.status.fullyQualifiedDomainName' not set. Has the server been provisioned successfully?", flexibleServer.Name)
		return nil, conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonWaitingForOwner)
	}
	serverFQDN := *flexibleServer.Status.FullyQualifiedDomainName

	adminPassword, err := secrets.LookupFromPtr(user.Spec.LocalUser.ServerAdminPassword)
	if err != nil {
		err = errors.Wrap(err, "failed to look up .spec.localUser.ServerAdminPassword")
		err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonSecretNotFound)
		return nil, err
	}

	// Admin User
	adminUser := user.Spec.LocalUser.ServerAdminUsername

	// Connect to the DB
	db, err := postgresqlutil.ConnectToDB(ctx, serverFQDN, postgresqlutil.DefaultMaintanenceDatabase, postgresqlutil.PSqlServerPort, adminUser, adminPassword)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to connect database. Server: %s, Database: %s, Port: %d, AdminUser: %s",
			serverFQDN,
			postgresqlutil.DefaultMaintanenceDatabase,
			postgresqlutil.PSqlServerPort,
			adminUser)
	}

	return db, nil
}
