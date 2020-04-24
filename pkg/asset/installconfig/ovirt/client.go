package ovirt

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
)

// getConnection is a convenience method to get a connection to ovirt api
// form a Config Object.
func getConnection(ovirtConfig Config) (*ovirtsdk.Connection, error) {
	con, err := ovirtsdk.NewConnectionBuilder().
		URL(ovirtConfig.URL).
		Username(ovirtConfig.Username).
		Password(ovirtConfig.Password).
		CAFile(ovirtConfig.CAFile).
		Insecure(ovirtConfig.Insecure).
		Build()
	if err != nil {
		return nil, err
	}
	return con, nil
}

// NewConnection returns a new client connection to oVirt's API endpoint.
// It is the responsibility of the caller to close the connection.
func NewConnection() (*ovirtsdk.Connection, error) {
	ovirtConfig, err := NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting ovirt configuration")
	}
	con, err := getConnection(ovirtConfig)
	if err != nil {
		return nil, errors.Wrap(err, "establishing ovirt connection")
	}
	return con, nil
}
