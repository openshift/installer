package aws

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *session.Session
	availabilityZones []string
	availableRegions  []string
	edgeZones         []string
	subnets           SubnetGroups
	vpcSubnets        SubnetGroups
	vpc               VPC
	instanceTypes     map[string]InstanceType

	Hosts           map[string]Host
	Region          string                     `json:"region,omitempty"`
	ProvidedSubnets []typesaws.Subnet          `json:"subnets,omitempty"`
	Services        []typesaws.ServiceEndpoint `json:"services,omitempty"`

	ec2Client *ec2.Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, subnets []typesaws.Subnet, services []typesaws.ServiceEndpoint) *Metadata {
	return &Metadata{Region: region, ProvidedSubnets: subnets, Services: services}
}

// Session holds an AWS session which can be used for AWS API calls
// during asset generation.
func (m *Metadata) Session(ctx context.Context) (*session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.unlockedSession(ctx)
}

func (m *Metadata) unlockedSession(ctx context.Context) (*session.Session, error) {
	if m.session == nil {
		var err error
		m.session, err = GetSessionWithOptions(WithRegion(m.Region), WithServiceEndpoints(m.Region, m.Services))
		if err != nil {
			return nil, fmt.Errorf("creating AWS session: %w", err)
		}
	}

	return m.session, nil
}

// EC2Client initiates a new EC2 client when one does not already exist, otherwise the existing client
// is returned.
func (m *Metadata) EC2Client(ctx context.Context) (*ec2.Client, error) {
	if m.ec2Client == nil {
		ec2Client, err := NewEC2Client(ctx, EndpointOptions{
			Region:    m.Region,
			Endpoints: m.Services,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create EC2 client: %w", err)
		}
		m.ec2Client = ec2Client
	}
	return m.ec2Client, nil
}

// AvailabilityZones retrieves a list of availability zones for the configured region.
func (m *Metadata) AvailabilityZones(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.availabilityZones) == 0 {
		session, err := m.unlockedSession(ctx)
		if err != nil {
			return nil, err
		}
		m.availabilityZones, err = availabilityZones(ctx, session, m.Region)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Availability Zones: %w", err)
		}
	}

	return m.availabilityZones, nil
}

// Regions retrieves a list of all regions.
func (m *Metadata) Regions(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.availableRegions) == 0 {
		client, err := m.EC2Client(ctx)
		if err != nil {
			return nil, err
		}

		output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{AllRegions: aws.Bool(false)})
		if err != nil {
			return nil, fmt.Errorf("failed to get all regions: %w", err)
		}

		for _, region := range output.Regions {
			m.availableRegions = append(m.availableRegions, *region.RegionName)
		}
	}

	return m.availableRegions, nil
}

// EdgeZones retrieves a list of Local and Wavelength zones for the configured region.
func (m *Metadata) EdgeZones(ctx context.Context) ([]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.edgeZones) == 0 {
		session, err := m.unlockedSession(ctx)
		if err != nil {
			return nil, err
		}

		m.edgeZones, err = edgeZones(ctx, session, m.Region)
		if err != nil {
			return nil, fmt.Errorf("getting Local Zones: %w", err)
		}
	}

	return m.edgeZones, nil
}

// EdgeSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be edge
// (i.e. Local Zone).
func (m *Metadata) EdgeSubnets(ctx context.Context) (Subnets, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Edge Subnets: %w", err)
	}
	return m.subnets.Edge, nil
}

// SetZoneAttributes retrieves AWS Zone attributes and update required fields in zones.
func (m *Metadata) SetZoneAttributes(ctx context.Context, zoneNames []string, zones Zones) error {
	sess, err := m.Session(ctx)
	if err != nil {
		return fmt.Errorf("unable to get aws session to populate zone details: %w", err)
	}
	azs, err := describeFilteredZones(ctx, sess, m.Region, zoneNames)
	if err != nil {
		return fmt.Errorf("unable to filter zones: %w", err)
	}

	for _, az := range azs {
		zoneName := awssdk.StringValue(az.ZoneName)
		if _, ok := zones[zoneName]; !ok {
			zones[zoneName] = &Zone{Name: zoneName}
		}
		if zones[zoneName].GroupName == "" {
			zones[zoneName].GroupName = awssdk.StringValue(az.GroupName)
		}
		if zones[zoneName].Type == "" {
			zones[zoneName].Type = awssdk.StringValue(az.ZoneType)
		}
		if az.ParentZoneName != nil {
			zones[zoneName].ParentZoneName = awssdk.StringValue(az.ParentZoneName)
		}
	}
	return nil
}

// AllZones return all the zones and it's attributes available on the region.
func (m *Metadata) AllZones(ctx context.Context) (Zones, error) {
	sess, err := m.Session(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get aws session to populate zone details: %w", err)
	}
	azs, err := describeAvailabilityZones(ctx, sess, m.Region, []string{})
	if err != nil {
		return nil, fmt.Errorf("unable to gather availability zones: %w", err)
	}
	zoneDesc := make(Zones, len(azs))
	for _, az := range azs {
		zoneName := awssdk.StringValue(az.ZoneName)
		zoneDesc[zoneName] = &Zone{
			Name:      zoneName,
			GroupName: awssdk.StringValue(az.GroupName),
			Type:      awssdk.StringValue(az.ZoneType),
		}
		if az.ParentZoneName != nil {
			zoneDesc[zoneName].ParentZoneName = awssdk.StringValue(az.ParentZoneName)
		}
	}
	return zoneDesc, nil
}

// PrivateSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be private
// (i.e. not public).
func (m *Metadata) PrivateSubnets(ctx context.Context) (Subnets, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Private Subnets: %w", err)
	}
	return m.subnets.Private, nil
}

// PublicSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be public
// (e.g. with suitable routing for hosting public load balancers).
func (m *Metadata) PublicSubnets(ctx context.Context) (Subnets, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Public Subnets: %w", err)
	}
	return m.subnets.Public, nil
}

// Subnets retrieves a group of subnet metadata that is indexed by subnet ID for all provided subnets.
// This includes private, public and edge subnets.
func (m *Metadata) Subnets(ctx context.Context) (SubnetGroups, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return m.subnets, fmt.Errorf("error retrieving all Subnets: %w", err)
	}
	return m.subnets, nil
}

// VPCSubnets retrieves a group of all subnet metadata that is indexed by subnet ID in the VPC of the provided subnets.
// These include cluster subnets (i.e. provided in the installconfig) and potentially other non-cluster subnets in the VPC.
//
// This func is only used for validations. Use func Subnets to select only cluster subnets.
func (m *Metadata) VPCSubnets(ctx context.Context) (SubnetGroups, error) {
	err := m.populateVPCSubnets(ctx)
	if err != nil {
		return m.vpcSubnets, fmt.Errorf("error retrieving Subnets in VPC: %w", err)
	}
	return m.vpcSubnets, nil
}

// VPC retrieves the VPC containing provided subnets.
func (m *Metadata) VPC(ctx context.Context) (VPC, error) {
	err := m.populateVPC(ctx)
	if err != nil {
		return m.vpc, fmt.Errorf("error retrieving VPC: %w", err)
	}
	return m.vpc, nil
}

// VPCID retrieves the ID of the VPC containing provided subnets.
func (m *Metadata) VPCID(ctx context.Context) (string, error) {
	err := m.populateVPC(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving VPC: %w", err)
	}
	return m.vpc.ID, nil
}

// SubnetByID retrieves subnet metadata for a subnet ID.
func (m *Metadata) SubnetByID(ctx context.Context, subnetID string) (subnet Subnet, err error) {
	err = m.populateSubnets(ctx)
	if err != nil {
		return subnet, fmt.Errorf("error retrieving subnet for ID %s: %w", subnetID, err)
	}

	if subnet, ok := m.subnets.Private[subnetID]; ok {
		return subnet, nil
	}

	if subnet, ok := m.subnets.Public[subnetID]; ok {
		return subnet, nil
	}

	if subnet, ok := m.subnets.Edge[subnetID]; ok {
		return subnet, nil
	}

	return subnet, fmt.Errorf("no subnet found for ID %s", subnetID)
}

// populateSubnets retrieves metadata for provided subnets.
func (m *Metadata) populateSubnets(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.ProvidedSubnets) == 0 {
		return errors.New("no subnets configured")
	}

	subnetGroups := m.subnets
	if subnetGroups.VpcID != "" || len(subnetGroups.Private) > 0 || len(subnetGroups.Public) > 0 || len(subnetGroups.Edge) > 0 {
		// Call to populate subnets has already happened
		return nil
	}

	client, err := m.EC2Client(ctx)
	if err != nil {
		return err
	}

	subnetIDs := make([]string, len(m.ProvidedSubnets))
	for i, subnet := range m.ProvidedSubnets {
		subnetIDs[i] = string(subnet.ID)
	}

	sb, err := subnets(ctx, client, subnetIDs, "")
	m.subnets = sb
	return err
}

// populateVPCSubnets retrieves metadata for all subnets in the VPC of provided subnets.
func (m *Metadata) populateVPCSubnets(ctx context.Context) error {
	// we need to populate provided subnets to get the VPC ID.
	if err := m.populateVPC(ctx); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	vpcSubnetGroups := m.vpcSubnets
	if len(vpcSubnetGroups.Private) > 0 || len(vpcSubnetGroups.Public) > 0 || len(vpcSubnetGroups.Edge) > 0 {
		// Call to populate subnets has already happened
		return nil
	}

	client, err := m.EC2Client(ctx)
	if err != nil {
		return err
	}

	sb, err := subnets(ctx, client, nil, m.vpc.ID)
	m.vpcSubnets = sb
	return err
}

// populateVPC retrieves metadata for the VPC of provided subnets.
func (m *Metadata) populateVPC(ctx context.Context) error {
	// we need to populate provided subnets to get the VPC ID.
	if err := m.populateSubnets(ctx); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.vpc.ID != "" {
		// Call to populate vpc has already happened
		return nil
	}

	client, err := m.EC2Client(ctx)
	if err != nil {
		return err
	}

	vpc, err := vpc(ctx, client, m.subnets.VpcID)
	m.vpc = vpc
	return err
}

// InstanceTypes retrieves instance type metadata indexed by InstanceType for the configured region.
func (m *Metadata) InstanceTypes(ctx context.Context) (map[string]InstanceType, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.instanceTypes) == 0 {
		session, err := m.unlockedSession(ctx)
		if err != nil {
			return nil, err
		}

		m.instanceTypes, err = instanceTypes(ctx, session, m.Region)
		if err != nil {
			return nil, fmt.Errorf("error listing instance types: %w", err)
		}
	}

	return m.instanceTypes, nil
}

// DedicatedHosts retrieves all hosts available for use to verify against this installation for configured region.
func (m *Metadata) DedicatedHosts(ctx context.Context) (map[string]Host, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.Hosts) == 0 {
		awsSession, err := m.unlockedSession(ctx)
		if err != nil {
			return nil, err
		}

		m.Hosts, err = dedicatedHosts(ctx, awsSession, m.Region)
		if err != nil {
			return nil, fmt.Errorf("error listing dedicated hosts: %w", err)
		}
	}

	return m.Hosts, nil
}
