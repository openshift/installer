package packet

import (
	"context"

	packngo "github.com/packethost/packngo"
	"github.com/pkg/errors"
)

const (
	PACKET_CONSUMER_TOKEN = "redhat openshift ipi"
)

//go:generate mockgen -source=./client.go -destination=mock/packetclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	ListProjects(ctx context.Context) ([]packngo.Project, error)
	ListFacilities(ctx context.Context) ([]packngo.Facility, error)
	ListPlans(ctx context.Context) ([]packngo.Plan, error)
}

type Client struct {
	OrganizationID string
	FacilityID     string
	ProjectID      string

	Conn *packngo.Client
}

func (c *Client) ListProjects(ctx context.Context) ([]packngo.Project, error) {
	return nil, nil
}

func (c *Client) ListFacilities(ctx context.Context) ([]packngo.Facility, error) {
	return nil, nil
}

func (c *Client) ListPlans(ctx context.Context) ([]packngo.Plan, error) {
	return nil, nil
}

var _ API = &Client{}

// getConnection is a convenience method to get a Packet API client
// from a Config Object.
func getConnection(c Config) (*packngo.Client, error) {
	return packngo.NewClientWithBaseURL(
		PACKET_CONSUMER_TOKEN, c.APIKey, nil, c.APIURL,
	)
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
