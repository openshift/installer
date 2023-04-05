package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	session           *session.Session
	availabilityZones []string
	privateSubnets    map[string]Subnet
	publicSubnets     map[string]Subnet
	edgeSubnets       map[string]Subnet
	vpc               string
	instanceTypes     map[string]InstanceType

	Region   string                     `json:"region,omitempty"`
	Subnets  []string                   `json:"subnets,omitempty"`
	Services []typesaws.ServiceEndpoint `json:"services,omitempty"`

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, subnets []string, services []typesaws.ServiceEndpoint) *Metadata {
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
			return nil, errors.Wrap(err, "creating AWS session")
		}
	}

	return m.session, nil
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
			return nil, errors.Wrap(err, "error retrieving Availability Zones")
		}
	}

	return m.availabilityZones, nil
}

// EdgeSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be edge
// (i.e. Local Zone).
func (m *Metadata) EdgeSubnets(ctx context.Context) (map[string]Subnet, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving Edge Subnets")
	}
	return m.edgeSubnets, nil
}

// PrivateSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be private
// (i.e. not public).
func (m *Metadata) PrivateSubnets(ctx context.Context) (map[string]Subnet, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving Private Subnets")
	}
	return m.privateSubnets, nil
}

// PublicSubnets retrieves subnet metadata indexed by subnet ID, for
// subnets that the cloud-provider logic considers to be public
// (e.g. with suitable routing for hosting public load balancers).
func (m *Metadata) PublicSubnets(ctx context.Context) (map[string]Subnet, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving Public Subnets")
	}
	return m.publicSubnets, nil
}

// VPC retrieves the VPC ID containing PublicSubnets and PrivateSubnets.
func (m *Metadata) VPC(ctx context.Context) (string, error) {
	err := m.populateSubnets(ctx)
	if err != nil {
		return "", errors.Wrap(err, "error retrieving VPC")
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

	sb, err := subnets(ctx, session, m.Region, m.Subnets)
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
			return nil, errors.Wrap(err, "error listing instance types")
		}
	}

	return m.instanceTypes, nil
}
