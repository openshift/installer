package packet

import (
	packngo "github.com/packethost/packngo"
	"github.com/pkg/errors"
)

// getConnection is a convenience method to get a connection to packet api
// form a Config Object.
func getConnection(_ Config) (*packngo.Client, error) {
	// TODO(displague) NewClientWith...
	con, err := packngo.NewClient()
	if err != nil {
		return nil, err
	}
	return con, nil
}

// NewConnection returns a new client connection to Packet's API endpoint.
// It is the responsibility of the caller to close the connection.
func NewConnection() (*packngo.Client, error) {
	packetConfig, err := NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting Engine configuration")
	}
	con, err := getConnection(packetConfig)
	if err != nil {
		return nil, errors.Wrap(err, "establishing Engine connection")
	}
	return con, nil
}
