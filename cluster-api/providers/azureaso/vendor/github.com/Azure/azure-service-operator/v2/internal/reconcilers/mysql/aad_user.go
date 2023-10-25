/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package mysql

import (
	"context"
	"database/sql"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	asomysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	mysqlutil "github.com/Azure/azure-service-operator/v2/internal/util/mysql"
)

const Scope = "https://ossrdbms-aad.database.windows.net/.default"

type aadUser struct {
	user               *asomysql.User
	resourceResolver   *resolver.Resolver
	credentialProvider identity.CredentialProvider
	log                logr.Logger
}

var _ Connector = &aadUser{}

func (u *aadUser) CreateOrUpdate(ctx context.Context) error {
	db, err := u.connectToDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	u.log.V(Status).Info("Creating MySQL AAD user")

	username := u.username()
	err = mysqlutil.CreateOrUpdateAADUser(ctx, db, u.user.Spec.AzureName, username)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	// Ensure that the privileges are set
	err = mysqlutil.ReconcileUserServerPrivileges(ctx, db, username, u.user.Spec.Hostname, u.user.Spec.Privileges)
	if err != nil {
		return errors.Wrap(err, "ensuring server roles")
	}

	err = mysqlutil.ReconcileUserDatabasePrivileges(ctx, db, username, u.user.Spec.Hostname, u.user.Spec.DatabasePrivileges)
	if err != nil {
		return errors.Wrap(err, "ensuring database roles")
	}

	u.log.V(Status).Info("Successfully reconciled MySQLUser")
	return nil
}

func (u *aadUser) Delete(ctx context.Context) error {
	db, err := u.connectToDB(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: There's still probably some ways that this user can be deleted but that we don't detect (and
	// TODO: so might cause an error triggering the resource to get stuck).
	// TODO: Cases where the server is in the process of being
	// TODO: deleted (or all system tables have been wiped?) might also exist...
	username := u.username()
	err = mysqlutil.DropUser(ctx, db, username)
	if err != nil {
		return err
	}

	return nil
}

func (u *aadUser) Exists(ctx context.Context) (bool, error) {
	db, err := u.connectToDB(ctx)
	if err != nil {
		return false, err
	}
	defer db.Close()

	username := u.username()
	exists, err := mysqlutil.DoesUserExist(ctx, db, username)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u *aadUser) connectToDB(ctx context.Context) (*sql.DB, error) {
	serverFQDN, err := getServerFQDN(ctx, u.resourceResolver, u.user)
	if err != nil {
		return nil, err
	}

	if u.user.Spec.AADUser == nil {
		return nil, errors.Errorf("AAD User missing $.spec.aadUser field")
	}
	adminUser := u.user.Spec.AADUser.ServerAdminUsername
	if len(adminUser) == 0 {
		return nil, errors.Errorf("AAD User must specify $.spec.aadUser.serverAdminUsernamed")
	}

	return connectToDBAAD(ctx, u.credentialProvider, u.log, u.user, serverFQDN, adminUser)
}

func (u *aadUser) username() string {
	if u.user.Spec.AADUser.Alias != "" {
		return u.user.Spec.AADUser.Alias
	}

	return u.user.Spec.AzureName
}

func connectToDBAAD(ctx context.Context, credentialProvider identity.CredentialProvider, log logr.Logger, user *asomysql.User, fqdn string, adminUser string) (*sql.DB, error) {
	credential, err := credentialProvider.GetCredential(ctx, user)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get credential")
	}
	token, err := credential.TokenCredential().GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{Scope}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token from credential")
	}
	log.V(Verbose).Info("Retrieved token for MySQL", "scope", Scope, "expires", token.ExpiresOn)

	// Connect to the DB
	db, err := mysqlutil.ConnectToDBAAD(ctx, fqdn, mysqlutil.SystemDatabase, mysqlutil.ServerPort, adminUser, token.Token)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"failed to connect database. Server: %s, Database: %s, Port: %d, AdminUser: %s",
			fqdn,
			mysqlutil.SystemDatabase,
			mysqlutil.ServerPort,
			adminUser)
	}

	return db, nil
}
