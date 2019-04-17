package ovirt

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// GetConnection is a convenience method to get a connection to ovirt api
// form a Config Object.
func GetConnection(config Config) (*ovirtsdk.Connection, error) {
	con, err := ovirtsdk.NewConnectionBuilder().
		URL(config.URL).
		Username(config.Username).
		Password(config.Password).
		CAFile(config.CAFile).
		Insecure(config.Insecure).
		Build()
	if err != nil {
		return nil, err
	}
	return con, nil
}
