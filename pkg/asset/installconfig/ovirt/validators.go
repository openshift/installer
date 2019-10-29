package ovirt

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
)

// authenticated takes an ovirt platform and validates
// its connection to the API by establishing
// the connection and authenticating successfully.
// The API connection is closed in the end and must leak
// or be reused in any way.
func authenticated(c *Config) survey.Validator {
	return func(val interface{}) error {
		connection, err := ovirtsdk.NewConnectionBuilder().
			URL(c.URL).
			Username(c.Username).
			Password(fmt.Sprint(val)).
			CAFile(c.CAFile).
			Insecure(c.Insecure).
			Build()

		if err != nil {
			return errors.Errorf("failed to construct connection to oVirt platform %s", err)
		}

		defer connection.Close()

		err = connection.Test()
		if err != nil {
			return errors.Errorf("failed to connect to oVirt platform %s", err)
		}
		return nil
	}

}
