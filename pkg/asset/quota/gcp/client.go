package gcp

import (
	"context"

	"github.com/pkg/errors"
	computev1 "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
)

// MachineTypeGetter returns the machine type info for a type in a zone using GCP API.
type MachineTypeGetter interface {
	GetMachineType(zone string, machineType string) (*computev1.MachineType, error)
}

// Client is GCP client for calculating quota constraint.
type Client struct {
	computeSvc *computev1.Service

	projectID string
}

// NewClient returns Client using the context and session.
func NewClient(ctx context.Context, sess *gcpclient.Session, projectID string) (*Client, error) {
	svc, err := computev1.NewService(ctx, option.WithCredentials(sess.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}
	return &Client{computeSvc: svc, projectID: projectID}, nil
}

// GetMachineType returns the machine type info for a type in a zone using the client.
func (c *Client) GetMachineType(zone string, machineType string) (*computev1.MachineType, error) {
	return c.computeSvc.MachineTypes.Get(c.projectID, zone, machineType).Context(context.TODO()).Do()
}
