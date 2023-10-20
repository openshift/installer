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

type SQLRoleDelta struct {
	AddedRoles   set.Set[string]
	DeletedRoles set.Set[string]
}

func DiffCurrentAndExpectedSQLRoles(currentRoles set.Set[string], expectedRoles set.Set[string]) SQLRoleDelta {
	result := SQLRoleDelta{
		AddedRoles:   set.Make[string](),
		DeletedRoles: set.Make[string](),
	}

	for priv := range expectedRoles {
		// If an expected Role isn't in the current Role set, we need to add it
		if !currentRoles.Contains(priv) {
			result.AddedRoles.Add(priv)
		}
	}

	for priv := range currentRoles {
		// If a current Role isn't in the expected set, we need to remove it
		if !expectedRoles.Contains(priv) {
			result.DeletedRoles.Add(priv)
		}
	}

	return result
}

// GetUserServerRoles gets the server-level roles the user has as a set.
func GetUserServerRoles(ctx context.Context, db *sql.DB, user SQLUser) (set.Set[string], error) {

	rows, err := db.QueryContext(
		ctx,
		"SELECT c.rolname as role_name FROM pg_roles a INNER JOIN pg_auth_members b on a.oid = b.member INNER JOIN pg_roles c ON b.roleid = c.oid WHERE a.rolname = $1",
		user.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "listing grants for user %s", user)
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

// ReconcileUserServerRoles revokes and grants server-level roles as
// needed so the roles for the user match those passed in.
func ReconcileUserServerRoles(ctx context.Context, db *sql.DB, user SQLUser, roles []string) error {
	var errs []error
	desiredRoles := set.Make[string](roles...)

	currentRoles, err := GetUserServerRoles(ctx, db, user)
	if err != nil {
		return errors.Wrapf(err, "couldn't get existing roles for user %s", user)
	}

	privsDiff := DiffCurrentAndExpectedSQLRoles(currentRoles, desiredRoles)
	err = addRoles(ctx, db, user, privsDiff.AddedRoles)
	if err != nil {
		errs = append(errs, err)
	}
	err = deleteRoles(ctx, db, user, privsDiff.DeletedRoles)
	if err != nil {
		errs = append(errs, err)
	}

	err = kerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	return nil
}

func addRoles(ctx context.Context, db *sql.DB, user SQLUser, roles set.Set[string]) error {
	if len(roles) == 0 {
		// Nothing to do
		return nil
	}
	var errorStrings []string
	for _, role := range roles.Values() {

		if err := FindBadChars(role); err != nil {
			return errors.Wrap(err, "problem found with role")
		}
		roleExists, err := RoleExists(ctx, db, role)
		if err != nil {
			return err
		}
		if !roleExists {
			return errors.Wrap(err, fmt.Sprintf("Role %q does not exists", role))
		}
	}
	toAdd := strings.Join(roles.Values(), ",")
	_, err := db.ExecContext(ctx, fmt.Sprintf("GRANT %s TO \"%s\"", toAdd, user.Name))
	if err != nil {
		errorStrings = append(errorStrings, err.Error())
	}
	if len(errorStrings) != 0 {
		return fmt.Errorf(strings.Join(errorStrings, "\n"))
	}
	return err
}

func deleteRoles(ctx context.Context, db *sql.DB, user SQLUser, roles set.Set[string]) error {
	if len(roles) == 0 {
		// Nothing to do
		return nil
	}

	toDelete := strings.Join(roles.Values(), ",")
	_, err := db.ExecContext(ctx, fmt.Sprintf("REVOKE %s FROM \"%s\"", toDelete, user.Name))

	return err
}
