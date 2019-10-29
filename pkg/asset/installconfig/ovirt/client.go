package ovirt

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// GetConnection is a convenience method to get a connection to ovirt api
// form a Config Object.
func GetConnection(ovirtConfig Config) (*ovirtsdk.Connection, error) {
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
