package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

//go:generate mockgen -source=./client.go -destination=.mock/gcpclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetNetwork(ctx context.Context, network, project string) (*compute.Network, error)
	GetPublicDomains(ctx context.Context, project string) ([]string, error)
	GetPublicDNSZone(ctx context.Context, baseDomain, project string) (*dns.ManagedZone, error)
	GetSubnetworks(ctx context.Context, network, project, region string) ([]*compute.Subnetwork, error)
}

// Client makes calls to the GCP API.
type Client struct {
	ssn *Session
}

// NewClient initializes a client with a session.
func NewClient(ctx context.Context) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	client := &Client{
		ssn: ssn,
	}
	return client, nil
}

// GetNetwork uses the GCP Compute Service API to get a network by name from a project.
func (c *Client) GetNetwork(ctx context.Context, network, project string) (*compute.Network, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}
	res, err := svc.Networks.Get(project, network).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get network %s", network)
	}
	return res, nil
}

// GetPublicDomains returns all of the domains from among the project's public DNS zones.
func (c *Client) GetPublicDomains(ctx context.Context, project string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	svc, err := c.getDNSService(ctx)
	if err != nil {
		return []string{}, err
	}

	var publicZones []string
	req := svc.ManagedZones.List(project).Context(ctx)
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, v := range page.ManagedZones {
			if v.Visibility != "private" {
				publicZones = append(publicZones, strings.TrimSuffix(v.DnsName, "."))
			}
		}
		return nil
	}); err != nil {
		return publicZones, err
	}
	return publicZones, nil
}

// GetPublicDNSZone returns a public DNS zone for a basedomain.
func (c *Client) GetPublicDNSZone(ctx context.Context, project, baseDomain string) (*dns.ManagedZone, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	svc, err := c.getDNSService(ctx)
	if err != nil {
		return nil, err
	}

	req := svc.ManagedZones.List(project).DnsName(baseDomain).Context(ctx)
	var res *dns.ManagedZone
	if err := req.Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for idx, v := range page.ManagedZones {
			if v.Visibility != "private" {
				res = page.ManagedZones[idx]
			}
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to list DNS Zones")
	}
	if res == nil {
		return nil, errors.New("no matching public DNS Zone found")
	}
	return res, nil
}

// GetSubnetworks uses the GCP Compute Service API to retrieve all subnetworks in a given network.
func (c *Client) GetSubnetworks(ctx context.Context, network, project, region string) ([]*compute.Subnetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("network eq .*%s", network)
	req := svc.Subnetworks.List(project, region).Filter(filter)
	var res []*compute.Subnetwork
	if err := req.Pages(ctx, func(page *compute.SubnetworkList) error {
		for _, subnet := range page.Items {
			res = append(res, subnet)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) getComputeService(ctx context.Context) (*compute.Service, error) {
	svc, err := compute.NewService(ctx, option.WithCredentials(c.ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}
	return svc, nil
}

func (c *Client) getDNSService(ctx context.Context) (*dns.Service, error) {
	svc, err := dns.NewService(ctx, option.WithCredentials(c.ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dns service")
	}
	return svc, nil
}
