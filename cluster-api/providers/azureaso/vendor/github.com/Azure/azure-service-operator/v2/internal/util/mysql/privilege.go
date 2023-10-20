/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

type SQLPrivilegeDelta struct {
	AddedPrivileges   set.Set[string]
	DeletedPrivileges set.Set[string]
}

// TODO: ALL is actually not supported in FlexibleServer so there's not a lot of need for this. It doesn't hurt though and makes
// TODO: us more robust if we want to expand to support MySQL users on other server types (standalone server, etc).
const sqlAll = "ALL"

func DiffCurrentAndExpectedSQLPrivileges(currentPrivileges set.Set[string], expectedPrivileges set.Set[string]) SQLPrivilegeDelta {
	result := SQLPrivilegeDelta{
		AddedPrivileges:   set.Make[string](),
		DeletedPrivileges: set.Make[string](),
	}

	for priv := range expectedPrivileges {
		// Escape hatch - if they ask for ALL then we just grant ALL
		// and don't delete any, since the user should have all of
		// them.
		if IsSQLAll(priv) {
			return SQLPrivilegeDelta{
				AddedPrivileges:   set.Make[string](sqlAll),
				DeletedPrivileges: set.Make[string](),
			}
		}

		// If an expected privilege isn't in the current privilege set, we need to add it
		if !currentPrivileges.Contains(priv) {
			result.AddedPrivileges.Add(priv)
		}
	}

	for priv := range currentPrivileges {
		// If a current privilege isn't in the expected set, we need to remove it
		if !expectedPrivileges.Contains(priv) {
			result.DeletedPrivileges.Add(priv)
		}
	}

	return result
}

// IsSQLAll returns whether the string matches the special privilege value ALL.
func IsSQLAll(privilege string) bool {
	return strings.EqualFold(privilege, sqlAll)
}

// GetUserDatabasePrivileges gets the per-database privileges that the
// user has. The user can have different permissions to each
// database. The details of access are returned in the map, keyed by
// database name.
func GetUserDatabasePrivileges(ctx context.Context, db *sql.DB, user string, hostname string) (map[string]set.Set[string], error) {
	hostname = HostnameOrDefault(hostname)

	// Note: This works because we only assign permissions at the DB level, not at the table, column, etc levels -- if we assigned
	// permissions at more levels we would need to do something else here such as join multiple tables or
	// parse SHOW GRANTS with a regex.
	rows, err := db.QueryContext(
		ctx,
		"SELECT TABLE_SCHEMA, PRIVILEGE_TYPE FROM INFORMATION_SCHEMA.SCHEMA_PRIVILEGES WHERE GRANTEE = ?",
		formatUser(user, hostname))
	if err != nil {
		return nil, errors.Wrapf(err, "listing database grants for user %s", user)
	}
	defer rows.Close()

	results := make(map[string]set.Set[string])
	for rows.Next() {
		var database, privilege string
		err := rows.Scan(&database, &privilege)
		if err != nil {
			return nil, errors.Wrapf(err, "extracting privilege row")
		}

		var privileges set.Set[string]
		if existingPrivileges, found := results[database]; found {
			privileges = existingPrivileges
		} else {
			privileges = make(set.Set[string])
			results[database] = privileges
		}
		privileges.Add(privilege)
	}

	if rows.Err() != nil {
		return nil, errors.Wrapf(rows.Err(), "iterating database privileges")
	}

	return results, nil
}

// GetUserServerPrivileges gets the server-level privileges the user has as a set.
func GetUserServerPrivileges(ctx context.Context, db *sql.DB, user string, hostname string) (set.Set[string], error) {
	hostname = HostnameOrDefault(hostname)

	// Note: This works because we only assign permissions at the DB level, not at the table, column, etc levels -- if we assigned
	// permissions at more levels we would need to do something else here such as join multiple tables or
	// parse SHOW GRANTS with a regex.
	// Remove "USAGE" as it's special and we never grant or remove it.
	rows, err := db.QueryContext(
		ctx,
		"SELECT PRIVILEGE_TYPE FROM INFORMATION_SCHEMA.USER_PRIVILEGES WHERE GRANTEE = ? AND PRIVILEGE_TYPE != 'USAGE'",
		formatUser(user, hostname))
	if err != nil {
		return nil, errors.Wrapf(err, "listing grants for user %s", user)
	}
	defer rows.Close()

	result := make(set.Set[string])
	for rows.Next() {
		var row string
		err := rows.Scan(&row)
		if err != nil {
			return nil, errors.Wrapf(err, "extracting privilege field")
		}

		result.Add(row)
	}
	if rows.Err() != nil {
		return nil, errors.Wrapf(rows.Err(), "iterating privileges")
	}

	return result, nil
}

// ReconcileUserServerPrivileges revokes and grants server-level privileges as
// needed so the privileges for the user match those passed in.
func ReconcileUserServerPrivileges(ctx context.Context, db *sql.DB, user string, hostname string, privileges []string) error {
	var errs []error
	desiredPrivileges := set.Make[string](privileges...)

	currentPrivileges, err := GetUserServerPrivileges(ctx, db, user, hostname)
	if err != nil {
		return errors.Wrapf(err, "couldn't get existing privileges for user %s", user)
	}

	privsDiff := DiffCurrentAndExpectedSQLPrivileges(currentPrivileges, desiredPrivileges)
	err = addPrivileges(ctx, db, "", user, privsDiff.AddedPrivileges)
	if err != nil {
		errs = append(errs, err)
	}
	err = deletePrivileges(ctx, db, "", user, privsDiff.DeletedPrivileges)
	if err != nil {
		errs = append(errs, err)
	}

	err = kerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	return nil
}

// ReconcileUserDatabasePrivileges revokes and grants database privileges as needed
// so they match the ones passed in. If there's an error applying
// privileges for one database it will still continue to apply
// privileges for subsequent databases (before reporting all errors).
func ReconcileUserDatabasePrivileges(ctx context.Context, conn *sql.DB, user string, hostname string, dbPrivs map[string][]string) error {
	desiredPrivs := make(map[string]set.Set[string])
	for database, privs := range dbPrivs {
		desiredPrivs[database] = set.Make[string](privs...)
	}

	currentPrivs, err := GetUserDatabasePrivileges(ctx, conn, user, hostname)
	if err != nil {
		return errors.Wrapf(err, "couldn't get existing database privileges for user %s", user)
	}

	allDatabases := make(set.Set[string])
	for db := range desiredPrivs {
		allDatabases.Add(db)
	}
	for db := range currentPrivs {
		allDatabases.Add(db)
	}

	var dbErrors []error
	for db := range allDatabases {
		privsDiff := DiffCurrentAndExpectedSQLPrivileges(
			currentPrivs[db],
			desiredPrivs[db],
		)

		err = addPrivileges(ctx, conn, db, user, privsDiff.AddedPrivileges)
		if err != nil {
			dbErrors = append(dbErrors, errors.Wrap(err, db))
		}
		err = deletePrivileges(ctx, conn, db, user, privsDiff.DeletedPrivileges)
		if err != nil {
			dbErrors = append(dbErrors, errors.Wrap(err, db))
		}
	}

	return kerrors.NewAggregate(dbErrors)
}

func addPrivileges(ctx context.Context, db *sql.DB, database string, user string, privileges set.Set[string]) error {
	if len(privileges) == 0 {
		// Nothing to do
		return nil
	}

	toAdd := strings.Join(privileges.Values(), ",")
	// TODO: Is there a way to just disable G201, which this violates?
	// We say //nolint:gosec below because gosec is trying to tell us this is a dangerous SQL query with a risk of SQL
	// injection. The user effectively has admin access to the DB through the operator already the minute that they can
	// create users with arbitrary permission levels.
	_, err := db.ExecContext(ctx, fmt.Sprintf("GRANT %s ON %s TO ?", toAdd, asGrantTarget(database)), user) //nolint:gosec

	return err
}

func deletePrivileges(ctx context.Context, db *sql.DB, database string, user string, privileges set.Set[string]) error {
	if len(privileges) == 0 {
		// Nothing to do
		return nil
	}

	toDelete := strings.Join(privileges.Values(), ",")
	// TODO: Is there a way to just disable G201, which this violates?
	// We say //nolint:gosec below because gosec is trying to tell us this is a dangerous SQL query with a risk of SQL
	// injection. The user effectively has admin access to the DB through the operator already the minute that they can
	// create users with arbitrary permission levels.
	tsql := fmt.Sprintf("REVOKE %s ON %s FROM ?", toDelete, asGrantTarget(database)) //nolint:gosec
	_, err := db.ExecContext(ctx, tsql, user)

	return err
}

// asGrantTarget formats the database name as a target suitable for a
// grant or revoke statement. If database is empty it returns "*.*"
// for server-level privileges.
func asGrantTarget(database string) string {
	if database == "" {
		return "*.*"
	}
	return fmt.Sprintf("`%s`.*", database)
}

func formatUser(user string, hostname string) string {
	// Wrap the user name in the weird formatting MySQL uses.
	return fmt.Sprintf("'%s'@'%s'", user, hostname)
}
