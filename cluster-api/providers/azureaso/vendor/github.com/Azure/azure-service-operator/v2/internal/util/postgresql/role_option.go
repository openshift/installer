/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

type SQLRoleOptionDelta struct {
	ChangedRoleOptions set.Set[RoleOption]
}

// RoleOptions PostgreSQL role options but without SuperUser or BypassRLS here,
// because they are only settable with an existing a superuser
// Azure Flexible server does not offer superuser access for customers
type RoleOptions struct {

	// WITH LOGIN or NOLOGIN
	Login bool

	// WITH CREATEROLE or NOCREATEROLE
	CreateRole bool

	// WITH CREATEDB or NOCREATEDB
	CreateDb bool

	// WITH REPLICATION or NOREPLICATION
	Replication bool
}

func DiffCurrentAndExpectedSQLRoleOptions(currentRoleOptions RoleOptions, expectedRoleOptions RoleOptions) SQLRoleOptionDelta {
	result := SQLRoleOptionDelta{
		ChangedRoleOptions: set.Make[RoleOption](),
	}
	if currentRoleOptions.Login != expectedRoleOptions.Login {
		if expectedRoleOptions.Login {
			result.ChangedRoleOptions.Add(Login)
		} else {
			result.ChangedRoleOptions.Add(NoLogin)
		}
	}
	if currentRoleOptions.CreateRole != expectedRoleOptions.CreateRole {
		if expectedRoleOptions.CreateRole {
			result.ChangedRoleOptions.Add(CreateRole)
		} else {
			result.ChangedRoleOptions.Add(NoCreateRole)
		}
	}
	if currentRoleOptions.CreateDb != expectedRoleOptions.CreateDb {
		if expectedRoleOptions.CreateDb {
			result.ChangedRoleOptions.Add(CreateDb)
		} else {
			result.ChangedRoleOptions.Add(NoCreateDb)
		}
	}
	if currentRoleOptions.Replication != expectedRoleOptions.Replication {
		if expectedRoleOptions.Replication {
			result.ChangedRoleOptions.Add(Replication)
		} else {
			result.ChangedRoleOptions.Add(NoReplication)
		}
	}
	return result
}

// GetUserRoleOptions gets the server-level RoleOptions the user has as a set.
func GetUserRoleOptions(ctx context.Context, db *sql.DB, user SQLUser) (*RoleOptions, error) {
	// Note: This query returns always all options which the enabled or disabled key
	// https://www.postgresql.org/docs/current/sql-alterrole.html
	// Exclude roles which start with 'pg_' as these are system roles we don't manage
	rows, err := db.QueryContext(
		ctx,
		"SELECT rolcanlogin, rolcreaterole, rolcreatedb, rolreplication FROM pg_roles WHERE rolname !~ '^pg_' AND rolname = $1",
		user.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "listing grants for user %s", user)
	}
	defer rows.Close()

	result := new(RoleOptions)

	for rows.Next() {
		err := rows.Scan(&result.Login, &result.CreateRole, &result.CreateDb, &result.Replication)
		if err != nil {
			return nil, errors.Wrapf(err, "extracting RoleOption field")
		}
		//No error handling required here, as sql returns already defined constants
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "iterating RoleOptions")
	}

	return result, nil
}

// ReconcileUserRoleOptions revokes and grants server-level role options as
// needed so the role options for the user match those passed in.
func ReconcileUserRoleOptions(ctx context.Context, db *sql.DB, user SQLUser, desiredOptions RoleOptions) error {
	var errs []error
	currentOptions, err := GetUserRoleOptions(ctx, db, user)
	if err != nil {
		return errors.Wrapf(err, "couldn't get existing RoleOptions for user %s", user)
	}

	privsDiff := DiffCurrentAndExpectedSQLRoleOptions(*currentOptions, desiredOptions)
	err = setRoleOptions(ctx, db, user, privsDiff.ChangedRoleOptions)
	if err != nil {
		errs = append(errs, err)
	}

	err = kerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	return nil
}

func setRoleOptions(ctx context.Context, db *sql.DB, user SQLUser, options set.Set[RoleOption]) error {
	if len(options) == 0 {
		// Nothing to do
		return nil
	}
	optionsValues := options.Values()
	values := make([]string, len(optionsValues))
	for i, e := range optionsValues {
		values[i] = string(e)
	}
	toChange := strings.Join(values, " ")
	_, err := db.ExecContext(ctx, fmt.Sprintf("ALTER ROLE \"%s\"  WITH %s", user.Name, toChange))

	return err
}

type RoleOption string

// see https://www.postgresql.org/docs/current/sql-createrole.html
var (
	Login         = RoleOption("LOGIN")
	CreateRole    = RoleOption("CREATEROLE")
	CreateDb      = RoleOption("CREATEDB")
	Replication   = RoleOption("REPLICATION")
	NoLogin       = RoleOption("NOLOGIN")
	NoCreateRole  = RoleOption("NOCREATEROLE")
	NoCreateDb    = RoleOption("NOCREATEDB")
	NoReplication = RoleOption("NOREPLICATION")
)
