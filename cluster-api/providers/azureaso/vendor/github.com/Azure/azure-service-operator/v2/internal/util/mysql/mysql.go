/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //mysql drive link
	"github.com/pkg/errors"
)

// ServerPort is the default server port for sql server
const ServerPort = 3306

// DriverName is driver name for psqldb connection
const DriverName = "mysql"

// SystemDatabase is the name of the system database in a MySQL server
// where users and privileges are stored (and which we can always
// assume will exist).
const SystemDatabase = "mysql"

func ConnectToDB(ctx context.Context, serverAddress string, database string, port int, user string, password string) (*sql.DB, error) {
	c := newConfig(serverAddress, database, port, user, password)
	return connectToDB(ctx, c)
}

// ConnectToDBAAD connects to the MySQL database using the specified user.
// user must be an AAD user of the form:
//   - "user@tenant.onmicrosoft.com" (for AAD users)
//   - "my-mi" (for Managed Identities)
//   - "mygroup" (for AAD groups)
func ConnectToDBAAD(ctx context.Context, serverAddress string, database string, port int, user string, token string) (*sql.DB, error) {
	c := newConfig(serverAddress, database, port, user, token)
	c.AllowCleartextPasswords = true // Must be set for AAD auth

	return connectToDB(ctx, c)
}

func HostnameOrDefault(hostname string) string {
	if hostname == "" {
		hostname = "%"
	}

	return hostname
}

func CreateOrUpdateUser(ctx context.Context, db *sql.DB, username string, hostname string, password string) error {
	hostname = HostnameOrDefault(hostname)

	// we call both CREATE and ALTER here so achieve an idempotent operation that also updates the password seamlessly
	// if it has changed

	statement := "CREATE USER IF NOT EXISTS ?@? IDENTIFIED BY ?"
	_, err := db.ExecContext(ctx, statement, username, hostname, password)
	if err != nil {
		return errors.Wrapf(err, "failed to create user %s", username)
	}

	statement = "ALTER USER IF EXISTS ?@? IDENTIFIED BY ?"
	_, err = db.ExecContext(ctx, statement, username, hostname, password)
	if err != nil {
		return errors.Wrapf(err, "failed to alter user %s", username)
	}

	return nil
}

// CreateOrUpdateAADUser creates or updates an AAD user.
// See https://learn.microsoft.com/en-us/azure/mysql/flexible-server/how-to-azure-ad for details on how to create AAD
// users.
// username can be either the actual AAD username (for real AAD users), the group name for groups, or
// the managed identity name for managed identities.
func CreateOrUpdateAADUser(ctx context.Context, db *sql.DB, username string, alias string) error {
	var err error
	if alias != "" {
		statement := "CREATE AADUSER IF NOT EXISTS ? AS ?"
		// Result rowset never returns any rows for CREATE AADUSER, so we just drop it
		_, err = db.ExecContext(ctx, statement, username, alias)
	} else {
		statement := "CREATE AADUSER IF NOT EXISTS ?"
		// Result rowset never returns any rows for CREATE AADUSER, so we just drop it
		_, err = db.ExecContext(ctx, statement, username)
	}
	if err != nil {
		return errors.Wrapf(err, "failed to create user %s", username)
	}

	return nil
}

// DoesUserExist checks if db contains user
func DoesUserExist(ctx context.Context, db *sql.DB, username string) (bool, error) {
	row := db.QueryRowContext(ctx, "SELECT User FROM mysql.user WHERE User = ?", username) // TODO: Username here has to be Alias in AAD case, if there is one
	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't exist
			return false, nil
		}

		// Something else went wrong
		return false, err
	}

	return true, nil
}

// DropUser drops a user from db
func DropUser(ctx context.Context, db *sql.DB, username string) error {
	_, err := db.ExecContext(ctx, "DROP USER IF EXISTS ?", username)
	return err
}

func newConfig(serverAddress string, database string, port int, user string, password string) *mysql.Config {
	c := mysql.NewConfig()
	c.Addr = fmt.Sprintf("%s:%d", serverAddress, port)
	c.DBName = database
	c.User = user
	c.Passwd = password
	c.TLSConfig = "true"

	// Set other options
	c.InterpolateParams = true
	c.Net = "tcp"

	return c
}

func connectToDB(ctx context.Context, c *mysql.Config) (*sql.DB, error) {
	db, err := sql.Open(DriverName, c.FormatDSN())
	if err != nil {
		return db, err
	}
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// We ping here to ensure that the connection is actually viable, as per
	// https://github.com/go-sql-driver/mysql/wiki/Examples#a-word-on-sqlopen
	err = db.PingContext(ctx)
	if err != nil {
		return db, errors.Wrapf(err, "error pinging the mysql db (%s/%s)", c.Addr, c.DBName)
	}

	return db, err
}
