/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package sql

import (
	"context"
	"database/sql"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	asosql "github.com/Azure/azure-service-operator/v2/api/sql/v1"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/reconcilers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	azuresqlutil "github.com/Azure/azure-service-operator/v2/internal/util/azuresql"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

type localUser struct {
	user               *asosql.User
	resourceResolver   *resolver.Resolver
	credentialProvider identity.CredentialProvider
	log                logr.Logger
}

var _ Connector = &localUser{}

func (u *localUser) CreateOrUpdate(ctx context.Context) error {
	secrets, err := u.resourceResolver.ResolveResourceSecretReferences(ctx, u.user)
	if err != nil {
		return reconcilers.ClassifyResolverError(err)
	}

	db, err := u.connectToDB(ctx, secrets)
	if err != nil {
		return err
	}
	defer db.Close()

	u.log.V(Status).Info("Creating/updating AzureSQL local user")

	password, err := secrets.LookupFromPtr(u.user.Spec.LocalUser.Password)
	if err != nil {
		return errors.Wrap(err, "failed to look up .spec.localUser.Password")
	}

	// Create or update the user. Note that this updates password if it has changed
	username := u.user.Spec.AzureName
	err = azuresqlutil.CreateOrUpdateUser(ctx, db, username, password)
	if err != nil {
		return errors.Wrap(err, "failed to create/update user")
	}

	// Ensure that the roles are set
	err = azuresqlutil.ReconcileUserRoles(ctx, db, username, u.user.Spec.Roles)
	if err != nil {
		return errors.Wrap(err, "ensuring server roles")
	}

	u.log.V(Status).Info("Successfully reconciled AzureSQLUser")

	return nil
}

func (u *localUser) Delete(ctx context.Context) error {
	secrets, err := u.resourceResolver.ResolveResourceSecretReferences(ctx, u.user)
	if err != nil {
		return err
	}

	db, err := u.connectToDB(ctx, secrets)
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: There's still probably some ways that this user can be deleted but that we don't detect (and
	// TODO: so might cause an error triggering the resource to get stuck).
	// TODO: Cases where the server is in the process of being
	// TODO: deleted (or all system tables have been wiped?) might also exist...
	err = azuresqlutil.DropUser(ctx, db, u.user.Spec.AzureName)
	if err != nil {
		return err
	}

	return nil
}

func (u *localUser) Exists(ctx context.Context) (bool, error) {
	secrets, err := u.resourceResolver.ResolveResourceSecretReferences(ctx, u.user)
	if err != nil {
		return false, err
	}

	db, err := u.connectToDB(ctx, secrets)
	if err != nil {
		return false, err
	}
	defer db.Close()

	exists, err := azuresqlutil.DoesUserExist(ctx, db, u.user.Spec.AzureName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *localUser) connectToDB(ctx context.Context, secrets genruntime.Resolved[genruntime.SecretReference, string]) (*sql.DB, error) {
	details, err := getOwnerDetails(ctx, u.resourceResolver, u.user)
	if err != nil {
		return nil, err
	}

	// Admin User
	adminUser := u.user.Spec.LocalUser.ServerAdminUsername

	if u.user.Spec.LocalUser.ServerAdminPassword != nil {
		var db *sql.DB
		var adminPassword string
		adminPassword, err = secrets.LookupFromPtr(u.user.Spec.LocalUser.ServerAdminPassword)
		if err != nil {
			err = errors.Wrap(err, "failed to look up .spec.localUser.ServerAdminPassword")
			err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonSecretNotFound)
			return nil, err
		}

		// Connect to the DB
		db, err = azuresqlutil.ConnectToDB(ctx, details.fqdn, details.database, azuresqlutil.ServerPort, adminUser, adminPassword)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"failed to connect database. Server: %s, Database: %s, Port: %d, AdminUser: %s",
				details.fqdn,
				details.database,
				azuresqlutil.ServerPort,
				adminUser)
		}

		return db, nil
	}

	return nil, errors.Errorf("local user does not have password")
}
