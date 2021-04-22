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
	// ListMetros(ctx context.Context) ([]packngo.Metro, error)
	ListPlans(ctx context.Context) ([]packngo.Plan, error)
}

type Client struct {
	OrganizationID string
	FacilityID     string
	MetroID        string
	ProjectID      string

	Conn *packngo.Client
}

func (c *Client) ListProjects(_ context.Context) ([]packngo.Project, error) {
	p, _, err := c.Conn.Projects.List(nil)
	return p, err
}

func (c *Client) ListFacilities(_ context.Context) ([]packngo.Facility, error) {
	f, _, err := c.Conn.Facilities.List(nil)
	return f, err
}

/*
func (c *Client) ListMetros(_ context.Context) ([]packngo.Metro, error) {
	m, _, err := c.Conn.Metros.List(nil)
	return m, err
}
*/

func (c *Client) ListPlans(_ context.Context) ([]packngo.Plan, error) {
	p, _, err := c.Conn.Plans.List(nil)
	return p, err
}

var _ API = &Client{}

// getConnection is a convenience method to get a Equinix Metal API client
// from a Config Object.
func getConnection(c *Config) (*packngo.Client, error) {
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
	if err := equinixmetalConfig.Save(); err != nil {
		return nil, errors.Wrap(err, "saving Engine configuration")
	}

	con, err := getConnection(equinixmetalConfig)
	if err != nil {
		return nil, errors.Wrap(err, "establishing Engine connection")
	}
	return con, nil
}
