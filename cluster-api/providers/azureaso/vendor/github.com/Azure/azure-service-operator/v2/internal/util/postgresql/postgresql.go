// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" //the pgx lib
	"github.com/pkg/errors"
)

// PSqlServerPort is the default server port for sql server
const PSqlServerPort = 5432

// PDriverName is driver name for psqldb connection
const PDriverName = "pgx"

// DefaultMaintanenceDatabase is the name of the database in a postgresql server
// where users and roles are stored (and which we can always
// assume will exist).
const DefaultMaintanenceDatabase = "postgres"

// ConnectToDB connects to the PostgreSQL db using the given credentials
func ConnectToDB(ctx context.Context, fullservername string, database string, port int, user string, password string) (*sql.DB, error) {

	connString := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=require connect_timeout=30", fullservername, user, password, port, database)

	db, err := sql.Open(PDriverName, connString)
	if err != nil {
		return db, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return db, err
	}

	return db, err
}

func CreateUser(ctx context.Context, db *sql.DB, username string, password string) (*SQLUser, error) {
	// make an effort to prevent sql injection
	//TODO find better solution to check user and password for SQL Injection
	if err := FindBadChars(username); err != nil {
		return nil, errors.Wrap(err, "problem found with username")
	}
	if err := FindBadChars(password); err != nil {
		return nil, errors.Wrap(err, "problem found with password")
	}
	_, err := db.ExecContext(ctx, fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", username, password))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create user %s", username)
	}
	return &SQLUser{Name: username}, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, user SQLUser, password string) error {
	// make an effort to prevent sql injection
	//TODO find better solution to check password for SQL Injection
	if err := FindBadChars(password); err != nil {
		return errors.Wrap(err, "problem found with password")
	}
	_, err := db.ExecContext(ctx, fmt.Sprintf("ALTER USER \"%s\" WITH PASSWORD '%s'", user.Name, password))
	if err != nil {
		return errors.Wrapf(err, "failed to alter user %s", user.Name)
	}
	return nil
}

func FindUserIfExist(ctx context.Context, db *sql.DB, username string) (*SQLUser, error) {
	res, err := db.ExecContext(ctx, "SELECT usename FROM pg_user WHERE usename = $1", username)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if rows > 0 {
		return &SQLUser{Name: username}, err
	} else {
		return nil, err
	}
}

// DoesUserExist checks if db contains user
func DoesUserExist(ctx context.Context, db *sql.DB, username string) (bool, error) {
	res, err := db.ExecContext(ctx, "SELECT usename FROM pg_user WHERE usename = $1", username)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

// DropUser drops a user from db
func DropUser(ctx context.Context, db *sql.DB, user string) error {
	if err := FindBadChars(user); err != nil {
		return errors.Wrap(err, "problem found with username")
	}

	_, err := db.ExecContext(ctx, fmt.Sprintf("DROP USER IF EXISTS \"%s\"", user))
	return err
}

// DatabaseExists checks if a database exists
func DatabaseExists(ctx context.Context, db *sql.DB, dbName string) (bool, error) {
	res, err := db.ExecContext(ctx, "SELECT datname FROM pg_catalog.pg_database WHERE datname = $1", dbName)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

// RoleExists checks if db contains role
func RoleExists(ctx context.Context, db *sql.DB, rolname string) (bool, error) {
	res, err := db.ExecContext(ctx, "SELECT * FROM pg_roles WHERE rolname = $1", rolname)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

// Use this type only for user, which are already checked
type SQLUser struct {
	Name string
}

// FindBadChars find the bad chars in a postgresql user
func FindBadChars(stack string) error {
	badChars := []string{
		"'",
		"\"",
		";",
		"--",
		"/*",
	}

	for _, s := range badChars {
		if idx := strings.Index(stack, s); idx > -1 {
			return fmt.Errorf("potentially dangerous character sequence found: '%s' at pos: %d", s, idx)
		}
	}
	return nil
}
