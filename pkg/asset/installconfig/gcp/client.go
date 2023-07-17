package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

//go:generate mockgen -source=./client.go -destination=./mock/gcpclient_generated.go -package=mock

const defaultTimeout = 2 * time.Minute

var (
	// RequiredBasePermissions is the list of permissions required for an installation.
	// A list of valid permissions can be found at https://cloud.google.com/iam/docs/understanding-roles.
	RequiredBasePermissions = []string{}
)

// API represents the calls made to the API.
type API interface {
	GetNetwork(ctx context.Context, network, project string) (*compute.Network, error)
	GetMachineType(ctx context.Context, project, zone, machineType string) (*compute.MachineType, error)
	GetMachineTypeWithZones(ctx context.Context, project, region, machineType string) (*compute.MachineType, sets.Set[string], error)
	GetPublicDomains(ctx context.Context, project string) ([]string, error)
	GetDNSZone(ctx context.Context, project, baseDomain string, isPublic bool) (*dns.ManagedZone, error)
	GetDNSZoneByName(ctx context.Context, project, zoneName string) (*dns.ManagedZone, error)
	GetSubnetworks(ctx context.Context, network, project, region string) ([]*compute.Subnetwork, error)
	GetProjects(ctx context.Context) (map[string]string, error)
	GetRegions(ctx context.Context, project string) ([]string, error)
	GetRecordSets(ctx context.Context, project, zone string) ([]*dns.ResourceRecordSet, error)
	GetZones(ctx context.Context, project, filter string) ([]*compute.Zone, error)
	GetEnabledServices(ctx context.Context, project string) ([]string, error)
	GetServiceAccount(ctx context.Context, project, serviceAccount string) (string, error)
	GetCredentials() *googleoauth.Credentials
	GetProjectPermissions(ctx context.Context, project string, permissions []string) (sets.Set[string], error)
	GetProjectByID(ctx context.Context, project string) (*cloudresourcemanager.Project, error)
	ValidateServiceAccountHasPermissions(ctx context.Context, project string, permissions []string) (bool, error)
}

// Client makes calls to the GCP API.
type Client struct {
	ssn *Session
}

// NewClient initializes a client with a session.
func NewClient(ctx context.Context) (*Client, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	client := &Client{
		ssn: ssn,
	}
	return client, nil
}

// GetMachineType uses the GCP Compute Service API to get the specified machine type.
func (c *Client) GetMachineType(ctx context.Context, project, zone, machineType string) (*compute.MachineType, error) {
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req, err := svc.MachineTypes.Get(project, zone, machineType).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return req, nil
}

// GetMachineTypeList retrieves the machine type with the specified fields.
func GetMachineTypeList(ctx context.Context, svc *compute.Service, project, region, machineType, fields string) ([]*compute.MachineType, error) {
	var machines []*compute.MachineType

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	filter := fmt.Sprintf("name = \"%s\" AND zone : %s-*", machineType, region)
	req := svc.MachineTypes.AggregatedList(project).Filter(filter).Context(ctx)
	if len(fields) > 0 {
		req.Fields(googleapi.Field(fields))
	}
	err := req.Pages(ctx, func(page *compute.MachineTypeAggregatedList) error {
		for _, scopedList := range page.Items {
			machines = append(machines, scopedList.MachineTypes...)
		}
		return nil
	})
	if len(machines) == 0 {
		return nil, errors.New("failed to fetch instance type, this error usually occurs if the region or the instance type is not found")
	}

	return machines, err
}

// GetMachineTypeWithZones retrieves the specified machine type and the zones in which it is available.
func (c *Client) GetMachineTypeWithZones(ctx context.Context, project, region, machineType string) (*compute.MachineType, sets.Set[string], error) {
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, nil, err
	}

	machines, err := GetMachineTypeList(ctx, svc, project, region, machineType, "")
	if err != nil {
		return nil, nil, err
	}

	zones := sets.New[string]()
	for _, machine := range machines {
		zones.Insert(machine.Zone)
	}

	// Restrict to zones available in the project
	pz, err := GetZones(ctx, svc, project, fmt.Sprintf("region eq .*%s", region))
	if err != nil {
		return nil, nil, err
	}
	projZones := sets.New[string]()
	for _, zone := range pz {
		projZones.Insert(zone.Name)
	}
	zones = zones.Intersection(projZones)

	return machines[0], zones, nil
}

// GetNetwork uses the GCP Compute Service API to get a network by name from a project.
func (c *Client) GetNetwork(ctx context.Context, network, project string) (*compute.Network, error) {
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	res, err := svc.Networks.Get(project, network).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get network %s", network)
	}
	return res, nil
}

// GetPublicDomains returns all of the domains from among the project's public DNS zones.
func (c *Client) GetPublicDomains(ctx context.Context, project string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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

// GetDNSZoneByName returns a DNS zone matching the `zoneName` if the DNS zone exists
// and can be seen (correct permissions for a private zone) in the project.
func (c *Client) GetDNSZoneByName(ctx context.Context, project, zoneName string) (*dns.ManagedZone, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	svc, err := c.getDNSService(ctx)
	if err != nil {
		return nil, err
	}
	returnedZone, err := svc.ManagedZones.Get(project, zoneName).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DNS Zones")
	}
	return returnedZone, nil
}

// GetDNSZone returns a DNS zone for a basedomain.
func (c *Client) GetDNSZone(ctx context.Context, project, baseDomain string, isPublic bool) (*dns.ManagedZone, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
			if v.Visibility != "private" && isPublic {
				res = page.ManagedZones[idx]
			} else if v.Visibility == "private" && !isPublic {
				res = page.ManagedZones[idx]
			}
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to list DNS Zones")
	}
	if res == nil {
		if isPublic {
			return nil, errors.New("no matching public DNS Zone found")
		}
		// A Private DNS Zone may be created (if the correct permissions exist)
		return nil, nil
	}
	return res, nil
}

// GetRecordSets returns all the records for a DNS zone.
func (c *Client) GetRecordSets(ctx context.Context, project, zone string) ([]*dns.ResourceRecordSet, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("network eq .*%s", network)
	req := svc.Subnetworks.List(project, region).Filter(filter)
	var res []*compute.Subnetwork

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if err := req.Pages(ctx, func(page *compute.SubnetworkList) error {
		res = append(res, page.Items...)
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

// GetProjects gets the list of project names and ids associated with the current user in the form
// of a map whose keys are ids and values are names.
func (c *Client) GetProjects(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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

// GetProjectByID retrieves the project specified by its ID.
func (c *Client) GetProjectByID(ctx context.Context, project string) (*cloudresourcemanager.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	svc, err := c.getCloudResourceService(ctx)
	if err != nil {
		return nil, err
	}

	return svc.Projects.Get(project).Context(ctx).Do()
}

// GetRegions gets the regions that are valid for the project. An error is returned when unsuccessful
func (c *Client) GetRegions(ctx context.Context, project string) ([]string, error) {
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	gcpRegionsList, err := svc.Regions.List(project).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get regions for project")
	}

	computeRegions := make([]string, len(gcpRegionsList.Items))
	for _, region := range gcpRegionsList.Items {
		computeRegions = append(computeRegions, region.Name)
	}

	return computeRegions, nil
}

// GetZones retrieves the available zones for the project.
func GetZones(ctx context.Context, svc *compute.Service, project, filter string) ([]*compute.Zone, error) {
	req := svc.Zones.List(project)
	if filter != "" {
		req = req.Filter(filter)
	}

	zones := []*compute.Zone{}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := req.Pages(ctx, func(page *compute.ZoneList) error {
		zones = append(zones, page.Items...)
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to get zones from project %s", project)
	}

	return zones, nil
}

// GetZones uses the GCP Compute Service API to get a list of zones from a project.
func (c *Client) GetZones(ctx context.Context, project, filter string) ([]*compute.Zone, error) {
	svc, err := c.getComputeService(ctx)
	if err != nil {
		return nil, err
	}

	return GetZones(ctx, svc, project, filter)
}

func (c *Client) getCloudResourceService(ctx context.Context) (*cloudresourcemanager.Service, error) {
	svc, err := cloudresourcemanager.NewService(ctx, option.WithCredentials(c.ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud resource service")
	}
	return svc, nil
}

// GetEnabledServices gets the list of enabled services for a project.
func (c *Client) GetEnabledServices(ctx context.Context, project string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
	svc, err := serviceusage.NewService(ctx, option.WithCredentials(c.ssn.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create service usage service")
	}
	return svc, nil
}

// GetServiceAccount retrieves a service account from a project if it exists.
func (c *Client) GetServiceAccount(ctx context.Context, project, serviceAccount string) (string, error) {
	svc, err := iam.NewService(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "failed create IAM service")
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	fullServiceAccountPath := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, serviceAccount)
	rsp, err := svc.Projects.ServiceAccounts.Get(fullServiceAccountPath).Context(ctx).Do()
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("failed to find resource %s", fullServiceAccountPath))
	}
	return rsp.Name, nil
}

// GetCredentials returns the credentials used to authenticate the GCP session.
func (c *Client) GetCredentials() *googleoauth.Credentials {
	return c.ssn.Credentials
}

func (c *Client) getPermissions(ctx context.Context, project string, permissions []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	service, err := c.getCloudResourceService(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cloud resource manager service")
	}

	projectsService := cloudresourcemanager.NewProjectsService(service)
	rb := &cloudresourcemanager.TestIamPermissionsRequest{Permissions: permissions}
	response, err := projectsService.TestIamPermissions(project, rb).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get Iam permissions")
	}

	return response.Permissions, nil
}

// GetProjectPermissions consumes a set of permissions and returns the set of found permissions for the service
// account (in the provided project). A list of valid permissions can be found at
// https://cloud.google.com/iam/docs/understanding-roles.
func (c *Client) GetProjectPermissions(ctx context.Context, project string, permissions []string) (sets.Set[string], error) {
	validPermissions, err := c.getPermissions(ctx, project, permissions)
	if err != nil {
		return nil, err
	}
	return sets.New[string](validPermissions...), nil
}

// ValidateServiceAccountHasPermissions compares the permissions to the set returned from the GCP API. Returns true
// if all permissions are available to the service account in the project.
func (c *Client) ValidateServiceAccountHasPermissions(ctx context.Context, project string, permissions []string) (bool, error) {
	validPermissions, err := c.GetProjectPermissions(ctx, project, permissions)
	if err != nil {
		return false, err
	}
	return validPermissions.Len() == len(permissions), nil
}
