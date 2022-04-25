package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"
)

//go:generate mockgen -source=./client.go -destination=./mock/gcpclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	GetNetwork(ctx context.Context, network, project string) (*compute.Network, error)
	GetMachineType(ctx context.Context, project, zone, machineType string) (*compute.MachineType, error)
	GetPublicDomains(ctx context.Context, project string) ([]string, error)
	GetPublicDNSZone(ctx context.Context, project, baseDomain string) (*dns.ManagedZone, error)
	GetSubnetworks(ctx context.Context, network, project, region string) ([]*compute.Subnetwork, error)
	GetProjects(ctx context.Context) (map[string]string, error)
	GetRecordSets(ctx context.Context, project, zone string) ([]*dns.ResourceRecordSet, error)
	GetZones(ctx context.Context, project, filter string) ([]*compute.Zone, error)
	GetEnabledServices(ctx context.Context, project string) ([]string, error)
}

// Client makes calls to the GCP API.
type Client struct {
	SSN *Session
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
		SSN: ssn,
	}
	return client, nil
}

// GetMachineType uses the GCP Compute Service API to get the specified machine type.
func (c *Client) GetMachineType(ctx context.Context, project, zone, machineType string) (*compute.MachineType, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	req, err := svc.MachineTypes.Get(project, zone, machineType).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
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
	if !strings.HasSuffix(baseDomain, ".") {
		baseDomain = fmt.Sprintf("%s.", baseDomain)
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

// GetRecordSets returns all the records for a DNS zone.
func (c *Client) GetRecordSets(ctx context.Context, project, zone string) ([]*dns.ResourceRecordSet, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	svc, err := c.getDNSService(ctx)
	if err != nil {
		return nil, err
	}

	req := svc.ResourceRecordSets.List(project, zone).Context(ctx)
	var rrSets []*dns.ResourceRecordSet
	if err := req.Pages(ctx, func(page *dns.ResourceRecordSetsListResponse) error {
		rrSets = append(rrSets, page.Rrsets...)
		return nil
	}); err != nil {
		return nil, err
	}
	return rrSets, nil
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
		res = append(res, page.Items...)
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) getComputeService(ctx context.Context) (*compute.Service, error) {
	svc, err := compute.NewService(ctx, option.WithCredentials(c.SSN.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compute service")
	}
	return svc, nil
}

func (c *Client) getDNSService(ctx context.Context) (*dns.Service, error) {
	svc, err := dns.NewService(ctx, option.WithCredentials(c.SSN.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dns service")
	}
	return svc, nil
}

// GetProjects gets the list of project names and ids associated with the current user in the form
// of a map whose keys are ids and values are names.
func (c *Client) GetProjects(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getCloudResourceService(ctx)
	if err != nil {
		return nil, err
	}

	req := svc.Projects.List()
	projects := make(map[string]string)
	if err := req.Pages(ctx, func(page *cloudresourcemanager.ListProjectsResponse) error {
		for _, project := range page.Projects {
			projects[project.ProjectId] = project.Name
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return projects, nil
}

// GetZones uses the GCP Compute Service API to get a list of zones from a project.
func (c *Client) GetZones(ctx context.Context, project, filter string) ([]*compute.Zone, error) {
	zones := []*compute.Zone{}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	req := svc.Zones.List(project)
	if filter != "" {
		req = req.Filter(filter)
	}

	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to get zones from project %s", project)
	}

	return zones, nil
}

func (c *Client) getCloudResourceService(ctx context.Context) (*cloudresourcemanager.Service, error) {
	svc, err := cloudresourcemanager.NewService(ctx, option.WithCredentials(c.SSN.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud resource service")
	}
	return svc, nil
}

// GetEnabledServices gets the list of enabled services for a project.
func (c *Client) GetEnabledServices(ctx context.Context, project string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	svc, err := c.getServiceUsageService(ctx)
	if err != nil {
		return nil, err
	}

	// List accepts a parent, which includes the type of resource with the id.
	parent := fmt.Sprintf("projects/%s", project)
	req := svc.Services.List(parent).Filter("state:ENABLED")
	var services []string
	if err := req.Pages(ctx, func(page *serviceusage.ListServicesResponse) error {
		for _, service := range page.Services {
			//services are listed in the form of project/services/serviceName
			index := strings.LastIndex(service.Name, "/")
			services = append(services, service.Name[index+1:])
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return services, nil
}

func (c *Client) getServiceUsageService(ctx context.Context) (*serviceusage.Service, error) {
	svc, err := serviceusage.NewService(ctx, option.WithCredentials(c.SSN.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create service usage service")
	}
	return svc, nil
}
