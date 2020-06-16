package equinixmetal

import (
	"context"

	packngo "github.com/packethost/packngo"
	"github.com/pkg/errors"
)

const (
	EQUINIXMETAL_CONSUMER_TOKEN = "redhat openshift ipi"
)

//go:generate mockgen -source=./client.go -destination=mock/equinixmetalclient_generated.go -package=mock

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

// getConnection is a convenience method to get a Equinix Metal API client
// from a Config Object.
func getConnection(c Config) (*packngo.Client, error) {
	return packngo.NewClientWithBaseURL(
		EQUINIXMETAL_CONSUMER_TOKEN, c.APIKey, nil, c.APIURL,
	)
}

// NewConnection returns a new client connection to Equinix Metal's API endpoint.
// It is the responsibility of the caller to close the connection.
func NewConnection() (*packngo.Client, error) {
	equinixmetalConfig, err := NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting Engine configuration")
	}
	con, err := getConnection(equinixmetalConfig)
	if err != nil {
		return nil, errors.Wrap(err, "establishing Engine connection")
	}
	return con, nil
}
