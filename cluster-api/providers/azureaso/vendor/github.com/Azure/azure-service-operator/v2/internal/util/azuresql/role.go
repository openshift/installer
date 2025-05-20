/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package azuresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

type SQLRoleDelta struct {
	AddedRoles   set.Set[string]
	DeletedRoles set.Set[string]
}

func DiffCurrentAndExpectedSQLRoles(currentRoles set.Set[string], expectedRoles set.Set[string]) SQLRoleDelta {
	result := SQLRoleDelta{
		AddedRoles:   set.Make[string](),
		DeletedRoles: set.Make[string](),
	}

	for role := range expectedRoles {
		// If an expected role isn't in the current role set, we need to add it
		if !currentRoles.Contains(role) {
			result.AddedRoles.Add(role)
		}
	}

	for role := range currentRoles {
		// If a current role isn't in the expected set, we need to remove it
		if !expectedRoles.Contains(role) {
			result.DeletedRoles.Add(role)
		}
	}

	return result
}

// GetUserRoles gets the roles assigned t othe user
func GetUserRoles(ctx context.Context, db *sql.DB, user string) (set.Set[string], error) {
	tsql := `
SELECT r.name role_principal_name
FROM sys.database_role_members rm 
JOIN sys.database_principals r 
    ON rm.role_principal_id = r.principal_id
JOIN sys.database_principals m 
    ON rm.member_principal_id = m.principal_id
where m.name = @user
`

	rows, err := db.QueryContext(
		ctx,
		tsql,
		sql.Named("user", user))
	if err != nil {
		return nil, errors.Wrapf(err, "listing roles for user %s", user)
	}
	defer rows.Close()

	result := make(set.Set[string])
	for rows.Next() {
		var row string
		err := rows.Scan(&row)
		if err != nil {
			return nil, errors.Wrapf(err, "extracting role field")
		}

		result.Add(row)
	}
	if rows.Err() != nil {
		return nil, errors.Wrapf(rows.Err(), "iterating roles")
	}

	return result, nil
}

// ReconcileUserRoles revokes and grants roles as
// needed so the roles for the user match those passed in.
func ReconcileUserRoles(ctx context.Context, db *sql.DB, user string, roles []string) error {
	var errs []error
	desiredRoles := set.Make[string](roles...)

	currentRoles, err := GetUserRoles(ctx, db, user)
	if err != nil {
		return errors.Wrapf(err, "couldn't get existing roles for user %s", user)
	}

	rolesDiff := DiffCurrentAndExpectedSQLRoles(currentRoles, desiredRoles)
	err = addRoles(ctx, db, user, rolesDiff.AddedRoles)
	if err != nil {
		errs = append(errs, err)
	}
	err = deleteRoles(ctx, db, user, rolesDiff.DeletedRoles)
	if err != nil {
		errs = append(errs, err)
	}

	err = kerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	return nil
}

func addRoles(ctx context.Context, db *sql.DB, user string, roles set.Set[string]) error {
	return alterRoles(ctx, db, user, roles, addMember)
}

func deleteRoles(ctx context.Context, db *sql.DB, user string, roles set.Set[string]) error {
	return alterRoles(ctx, db, user, roles, dropMember)
}

type alterRoleMode string

const (
	addMember  = alterRoleMode("ADD")
	dropMember = alterRoleMode("DROP")
)

func alterRoles(ctx context.Context, db *sql.DB, user string, roles set.Set[string], mode alterRoleMode) error {
	if len(roles) == 0 {
		// Nothing to do
		return nil
	}

	// We have a best-effort guard against SQL injection here because the SQL statements for ALTER ROLE don't support
	// named parameters. In any case, there's not much point in breaking this query with SQL injection as (assuming
	// you have permission to use the AzureSQL User CRD) you can just create a user w/ what permissions you want
	// and then log in to that user to attack the server.
	if err := findBadChars(user); err != nil {
		return errors.Wrap(err, "problem found with username")
	}

	builder := strings.Builder{}
	for role := range roles {
		if err := findBadChars(role); err != nil {
			return errors.Wrap(err, "problem found with role")
		}
		_, err := builder.WriteString(fmt.Sprintf("ALTER ROLE %s %s MEMBER %s;\n", role, string(mode), user))
		if err != nil {
			return errors.Wrapf(err, "failed to build T-SQL ALTER ROLE %s statement", mode)
		}
	}

	_, err := db.ExecContext(
		ctx,
		builder.String(),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to execute T-SQL ALTER ROLE %s statement", mode)
	}

	return nil
}
