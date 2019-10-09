package gcp

import (
	"context"
	"time"

	"google.golang.org/api/option"
	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"
)

// API represents the calls made to the API.
type API interface {
	GetNetwork(ctx context.Context, network, project string) (*compute.Network, error)
}

// Client makes calls to the GCP API.
type Client struct{}

// GetNetwork wraps the GCP Compute Service call to get a network by name from a project
func (c *Client) GetNetwork(ctx context.Context, network, project string) (*compute.Network, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := getComputeService(ctx)
	if err != nil {
		return nil, err
	}
	res, err := svc.Networks.Get(project, network).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get network %s", network)
	}
	return res, nil
}

func getComputeService(ctx context.Context) (*compute.Service, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	svc, err := compute.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}
	return svc, nil
}
