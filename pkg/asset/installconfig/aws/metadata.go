package aws

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *session.Session
	config            *awsv2.Config
	availabilityZones []string
	availableRegions  []string
	edgeZones         []string
	privateSubnets    Subnets
	publicSubnets     Subnets
	edgeSubnets       Subnets
	vpc               string
	instanceTypes     map[string]InstanceType

	Region   string                     `json:"region,omitempty"`
	Subnets  []typesaws.Subnet          `json:"subnets,omitempty"`
	Services []typesaws.ServiceEndpoint `json:"services,omitempty"`

	ec2Client *ec2.Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, subnets []typesaws.Subnet, services []typesaws.ServiceEndpoint) *Metadata {
	return &Metadata{Region: region, Subnets: subnets, Services: services}
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

func (m *Metadata) unlockedConfig(ctx context.Context) (*awsv2.Config, error) {
	if m.config == nil {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(m.Region))
		if err != nil {
			return nil, fmt.Errorf("creating AWS configuration: %w", err)
		}
		m.config = &cfg
	}
	return m.config, nil
}

// EC2Client initiates a new EC2 client when one does not already exist, otherwise the existing client
// is returned.
func (m *Metadata) EC2Client(ctx context.Context) (*ec2.Client, error) {
	if m.ec2Client == nil {
		cfg, err := m.unlockedConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("metadata failed to create config: %w", err)
		}

		optFns := []func(*ec2.Options){}
		for _, service := range m.Services {
			if service.Name == "ec2" {
				optFns = append(optFns, func(o *ec2.Options) {
					o.BaseEndpoint = awssdk.String(service.URL)
				})
				logrus.Warnf("setting ec2 endpoint URL to %s", service.URL)
				break
			}
		}

		m.ec2Client = ec2.NewFromConfig(*cfg, optFns...)
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
	return m.edgeSubnets, nil
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
	return m.privateSubnets, nil
}

// PublicSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be public
// (e.g. with suitable routing for hosting public load balancers).
func (m *Metadata) PublicSubnets(ctx context.Context) (Subnets, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Public Subnets: %w", err)
	}
	return m.publicSubnets, nil
}

// VPC retrieves the VPC ID containing PublicSubnets and PrivateSubnets.
func (m *Metadata) VPC(ctx context.Context) (string, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving VPC: %w", err)
	}
	return m.vpc, nil
}

func (m *Metadata) populateSubnets(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.Subnets) == 0 {
		return errors.New("no subnets configured")
	}

	if m.vpc != "" || len(m.privateSubnets) > 0 || len(m.publicSubnets) > 0 || len(m.edgeSubnets) > 0 {
		// Call to populate subnets has already happened
		return nil
	}

	session, err := m.unlockedSession(ctx)
	if err != nil {
		return err
	}

	subnetIDs := make([]string, len(m.Subnets))
	for i, subnet := range m.Subnets {
		subnetIDs[i] = string(subnet.ID)
	}

	sb, err := subnets(ctx, session, m.Region, subnetIDs)
	m.vpc = sb.VPC
	m.privateSubnets = sb.Private
	m.publicSubnets = sb.Public
	m.edgeSubnets = sb.Edge
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
